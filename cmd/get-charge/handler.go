package main

import (
	"card-payment-api/internal/adapters"
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
	getChargeUseCase := usecases.NewGetChargeUseCase(chargeRepository)

	chargeId := request.PathParameters["id"]
	if chargeId == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad request",
		}, nil
	}

	charge, err := getChargeUseCase.GetCharge(chargeId)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, err
	}

	jsonResponse, err := json.Marshal(charge)
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
