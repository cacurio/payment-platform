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
	dbClient := adapters.NewTokenDynamoRepository(client)
	tokenUSeCase := usecases.NewTokenUseCase(dbClient)

	if request.Body == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	// parse request body
	var tokenRequest dtos.TokenDTO
	err = json.Unmarshal([]byte(request.Body), &tokenRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	// create token
	token, err := tokenUSeCase.CreateToken(tokenRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, err
	}

	response := map[string]string{
		"token": token,
	}
	jsonResponse, err := json.Marshal(response)
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
