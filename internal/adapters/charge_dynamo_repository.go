package adapters

import (
	"card-payment-api/internal/domain"
	"card-payment-api/internal/domain/entities"
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ChargeDynamoRepository struct {
	DynamoClient *dynamodb.Client
}

func NewChargeDynamoRepository(
	dynamoClient *dynamodb.Client,
) *ChargeDynamoRepository {
	return &ChargeDynamoRepository{
		DynamoClient: dynamoClient,
	}
}

func (c *ChargeDynamoRepository) Save(charge *entities.ChargeEntity) error {
	tableName := os.Getenv("CHARGE_TABLE_NAME")
	item, err := attributevalue.MarshalMap(charge)
	if err != nil {
		return domain.ErrNoCreatedCharge
	}
	_, err = c.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return domain.ErrNoCreatedCharge
	}

	return nil
}

func (c *ChargeDynamoRepository) Get(id string) (*entities.ChargeEntity, error) {
	tableName := os.Getenv("CHARGE_TABLE_NAME")
	result, err := c.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, domain.ErrNoFoundCharge
	}
	var charge entities.ChargeEntity
	err = attributevalue.UnmarshalMap(result.Item, &charge)
	if err != nil {
		return nil, domain.ErrNoFoundCharge
	}

	return &charge, nil
}
