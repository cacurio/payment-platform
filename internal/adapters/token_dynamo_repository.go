package adapters

import (
	"card-payment-api/internal/domain"
	"card-payment-api/internal/domain/entities"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TokenDynamoRepository struct {
	DynamoClient *dynamodb.Client
}

func NewTokenDynamoRepository(
	dynamoClient *dynamodb.Client,
) *TokenDynamoRepository {
	return &TokenDynamoRepository{
		DynamoClient: dynamoClient,
	}
}

func (t *TokenDynamoRepository) Save(token *entities.TokenEntity) error {
	tokenName := os.Getenv("TABLE_NAME")
	item, err := attributevalue.MarshalMap(token)

	fmt.Println(item)

	if err != nil {
		return err
	}

	_, err = t.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tokenName),
		Item:      item,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t *TokenDynamoRepository) Get(token string) (*entities.TokenEntity, error) {
	tokenName := os.Getenv("TABLE_NAME")
	result, err := t.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tokenName),
		Key: map[string]types.AttributeValue{
			"tokenId": &types.AttributeValueMemberS{Value: token},
		},
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	var tokenEntity entities.TokenEntity
	err = attributevalue.UnmarshalMap(result.Item, &tokenEntity)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	return &tokenEntity, nil
}
