package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CardPaymentApiStackProps struct {
	awscdk.StackProps
}

func NewCardPaymentApiStack(scope constructs.Construct, id string, props *CardPaymentApiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CardPaymentApiQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	return stack
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

type ApiStackProps struct {
	awscdk.StackProps
	PaymentsTable awsdynamodb.Table
	TokensTable   awsdynamodb.Table
	ChargesTable  awsdynamodb.Table
}

// DatabaseStackProps defines the properties for the DatabaseStack.
type DatabaseStackProps struct {
	awscdk.StackProps
}

// The DatabaseStack stack is used to create the DynamoDB table that will store the payments.
type DatabaseStack struct {
	awscdk.Stack
	PaymentsTable awsdynamodb.Table
	TokensTable   awsdynamodb.Table
	ChargesTable  awsdynamodb.Table
}

// NewDatabaseStack creates a new instance of the DatabaseStack class.
func newDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// DynamoDB table
	paymentsTable := awsdynamodb.NewTable(stack, jsii.String("PaymentsTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("paymentId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	tokensTable := awsdynamodb.NewTable(stack, jsii.String("TokensTable"), &awsdynamodb.TableProps{
		TableName: jsii.String("Tokens"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("tokenId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	chargeTable := awsdynamodb.NewTable(stack, jsii.String("ChargesTable"), &awsdynamodb.TableProps{
		TableName: jsii.String("charges"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	return &DatabaseStack{
		Stack:         stack,
		PaymentsTable: paymentsTable,
		TokensTable:   tokensTable,
		ChargesTable:  chargeTable,
	}
}

func NewApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	currentDir, _ := os.Getwd()

	bankMockUrl := "https://c893-157-100-87-172.ngrok-free.app"

	// Token Function
	tokenFunctionPath := jsii.String(currentDir + "/cmd/token")
	tokenFn := awslambda.NewFunction(stack, jsii.String("TokenFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String(*tokenFunctionPath), nil),
		Environment: &map[string]*string{
			"TABLE_NAME": props.TokensTable.TableName(),
		},
	})

	// Charge Function
	chargeFunctionPath := jsii.String(currentDir + "/cmd/charge")
	chargeFn := awslambda.NewFunction(stack, jsii.String("ChargeFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String(*chargeFunctionPath), nil),
		Environment: &map[string]*string{
			"CHARGE_TABLE_NAME": props.ChargesTable.TableName(),
			"BANK_URL":          jsii.String(bankMockUrl),
			"TOKEN_TABLE_NAME":  props.TokensTable.TableName(),
			"TABLE_NAME":        props.TokensTable.TableName(),
		},
	})

	// get charge lambda function
	getChargeFunctionPath := jsii.String(currentDir + "/cmd/get-charge")
	getChargeFn := awslambda.NewFunction(stack, jsii.String("GetChargeFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String(*getChargeFunctionPath), nil),
		Environment: &map[string]*string{
			"CHARGE_TABLE_NAME": props.ChargesTable.TableName(),
		},
	})

	// Refund Function
	refundFunctionPath := jsii.String(currentDir + "/cmd/refund")
	refundFn := awslambda.NewFunction(stack, jsii.String("RefundFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String(*refundFunctionPath), nil),
		Environment: &map[string]*string{
			"CHARGE_TABLE_NAME": props.ChargesTable.TableName(),
			"BANK_URL":          jsii.String(bankMockUrl),
		},
	})

	// Grant the lambda function permission to read and write to the DynamoDB table
	props.TokensTable.GrantReadWriteData(tokenFn)
	props.ChargesTable.GrantReadWriteData(chargeFn)
	props.TokensTable.GrantReadData(chargeFn)
	props.ChargesTable.GrantReadData(getChargeFn)
	props.ChargesTable.GrantReadWriteData(refundFn)

	// API Gateway
	api := awsapigateway.NewRestApi(stack, jsii.String("PaymentPlatformApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("Payment Platform API"),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
		},
	})

	// Token
	token := api.Root().AddResource(jsii.String("token"), nil)
	token.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(tokenFn, nil), nil)

	// Charge
	charge := api.Root().AddResource(jsii.String("charge"), nil)
	charge.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(chargeFn, nil), nil)

	// Get Charge
	getCharge := charge.AddResource(jsii.String("{id}"), nil)
	getCharge.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getChargeFn, nil), nil)
	// Refund
	refund := charge.AddResource(jsii.String("refund"), nil)
	refund.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(refundFn, nil), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	dbStack := newDatabaseStack(app, "PaymentPlatformDatabaseStack", &DatabaseStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	NewApiStack(app, "PaymentPlatformApiStack", &ApiStackProps{
		awscdk.StackProps{
			Env: env(),
		},
		dbStack.PaymentsTable,
		dbStack.TokensTable,
		dbStack.ChargesTable,
	})

	app.Synth(nil)
}
