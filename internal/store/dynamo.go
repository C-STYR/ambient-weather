package store

import (
	"context"
	"sync"

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

func (s *Store) ProcessRecords(station string, chunkSize int, records []ambient.Record) error {
	wg := sync.WaitGroup{}

	numChunks := (len(records) + chunkSize - 1) / chunkSize
	errs := make(chan error, numChunks)

	for i := 0; i < len(records); i += chunkSize {
		chunk := records[i:min(i+chunkSize, len(records))]
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.BatchWriteRecords(station, chunk)
			if err != nil {
				errs <- err
			} else {
				errs <- nil
			}
		}()
	}
	wg.Wait()

	close(errs)

	for v := range errs {
		if v != nil {
			return v
		}
	}

	return nil
}

func (s *Store) BatchWriteRecords(station string, records []ambient.Record) error {

	reqs := []types.WriteRequest{}

	for _, r := range records {
		item, err := attributevalue.MarshalMap(r)
		if err != nil {
			return err
		}

		item["station"] = &types.AttributeValueMemberS{Value: station}
		reqs = append(reqs, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	requestMap := make(map[string][]types.WriteRequest)
	requestMap[s.tableName] = reqs

	// unused output has an UnprocessedItems field
	// TODO - implement retry loop with unprocessed records
	_, err := s.client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{RequestItems: requestMap})

	if err != nil {
		return err
	}

	return nil
}

// unused - part of prototyping, do not delete
func (s *Store) WriteRecords(station string, records []ambient.Record) error {

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
