package main

import (
	"card-payment-api/internal/adapters"
	"card-payment-api/internal/domain/dtos"
	"card-payment-api/internal/usecases"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func handlerRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// inject dependencies
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, err
	}
	client := dynamodb.NewFromConfig(cfg)

	chargeRepository := adapters.NewChargeDynamoRepository(client)
	bankMock := adapters.NewMockBankGateway()
	RefundUseCase := usecases.NewRefundUseCase(bankMock, chargeRepository)

	var chargeRequest dtos.RefundDTO
	err = json.Unmarshal([]byte(request.Body), &chargeRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, err
	}

	refundResponse, err := RefundUseCase.Execute(chargeRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	jsonResponse, err := json.Marshal(refundResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResponse),
	}, nil
}

func main() {
	lambda.Start(handlerRequest)
}
