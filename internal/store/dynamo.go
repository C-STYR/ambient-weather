package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cstyr/weather/internal/ambient"
)

type Store struct {
	client    *dynamodb.Client
	tableName string
}

func New(region, tableName string) (*Store, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	return &Store{
		client:    client,
		tableName: tableName,
	}, nil
}

func (s *Store) WriteRecords(station string, records []ambient.Record) error {
	// for each record...
	for _, r := range records {
		// convert the struct to the ddb format
		item, err := attributevalue.MarshalMap(r)
		if err != nil {
			return err
		}

		// add the station value
		item["station"] = &types.AttributeValueMemberS{Value: station}

		// add to the table
		_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: &s.tableName,
			Item:      item,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
