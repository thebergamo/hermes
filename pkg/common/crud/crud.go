package crud

import (
	"errors"
	BaseErrors "hermes/pkg/common/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type CrudRepository interface {
	List(item interface{}) (interface{}, error)
	Get(id string, item interface{}) (interface{}, error)
	Create(dto interface{}) (interface{}, error)
	Update(id string, dto interface{}) (interface{}, error)
	Delete(id string) error
}

type DynamoCrud struct {
	dynaClient dynamodbiface.DynamoDBAPI
	tableName  string
}

func (d *DynamoCrud) List(item interface{}) (interface{}, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(d.tableName),
	}

	result, err := d.dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorFailedToFetchRecord)
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func (d *DynamoCrud) Get(id string, item interface{}) (interface{}, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(d.tableName),
	}

	result, err := d.dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorFailedToFetchRecord)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func (d *DynamoCrud) Create(dto interface{}) (interface{}, error) {
	av, err := dynamodbattribute.MarshalMap(dto)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.tableName),
	}

	_, err = d.dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorCouldNotDynamoPutItem)
	}
	return &dto, nil
}

func (d *DynamoCrud) Update(id string, dto interface{}) (interface{}, error) {
	// Save dataset
	av, err := dynamodbattribute.MarshalMap(dto)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.tableName),
	}

	_, err = d.dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(BaseErrors.ErrorCouldNotDynamoPutItem)
	}
	return &dto, nil
}

func (d *DynamoCrud) Delete(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(d.tableName),
	}

	_, err := d.dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(BaseErrors.ErrorCouldNotDeleteItem)
	}

	return nil
}

func InitDynamoDbRepo(t string, d dynamodbiface.DynamoDBAPI) *DynamoCrud {
	return &DynamoCrud{
		dynaClient: d,
		tableName:  t,
	}
}
