package datasets

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	TableName                           = os.Getenv("TABLE_NAME")
	ErrorFailedToUnmarshalRecord        = "failed to unmarshal record"
	ErrorFailedToFetchRecord            = "failed to fetch record"
	ErrorInvalidDatasetData             = "invalid dataset data"
	ErrorInvalidProvider                = "invalid provider"
	ErrorInvalidType                    = "invalid type. Only SQL is supported"
	ErrorCouldNotMarshalItem            = "could not marshal item"
	ErrorCouldNotDeleteItem             = "could not delete item"
	ErrorCouldNotDynamoPutItem          = "could not dynamo put item error"
	ErrorDatasetAlreadyExists           = "dataset.Dataset already exists"
	ErrorDatasetDoesNotExists           = "dataset.Dataset does not exist"
	ErrorCouldNotSecureStoreCredentials = "could not store your credentials securely in SSM"
)

type DataSet struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Credentials string   `json:"credentials"`
	Type        string   `json:"type"`
	Provider    string   `json:"provider"`
	Tags        []string `json:"tags"`
}

func FetchDataset(id string, dynaClient dynamodbiface.DynamoDBAPI) (*DataSet, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(TableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)

	}

	item := new(DataSet)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchDatasets(dynaClient dynamodbiface.DynamoDBAPI) (*[]DataSet, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(TableName),
	}
	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]DataSet)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateDataset(req events.APIGatewayProxyRequest, dynaClient dynamodbiface.DynamoDBAPI, ssmClient *ssm.SSM) (
	*DataSet,
	error,
) {
	var d DataSet
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDatasetData)
	}

	if !IsProviderValid(d.Provider) {
		return nil, errors.New(ErrorInvalidProvider)
	}

	if !IsTypeValid(d.Type) {
		return nil, errors.New(ErrorInvalidType)
	}

	if d.Provider == "ssm" {
		_, err := ssmClient.PutParameter(&ssm.PutParameterInput{DataType: aws.String("text"), Name: aws.String(d.Id), Value: aws.String(d.Credentials), Type: aws.String("SecureString")})
		if err != nil {
			fmt.Println(err)
			return nil, errors.New(ErrorCouldNotSecureStoreCredentials)
		}
		d.Credentials = d.Id
	}
	// Save dataset

	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}

func UpdateDataset(req events.APIGatewayProxyRequest, dynaClient dynamodbiface.DynamoDBAPI, ssmClient *ssm.SSM) (
	*DataSet,
	error,
) {
	var d DataSet
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDatasetData)
	}

	// Check if dataset exists
	currentDataset, _ := FetchDataset(d.Id, dynaClient)
	if currentDataset != nil && len(currentDataset.Name) == 0 {
		return nil, errors.New(ErrorDatasetDoesNotExists)
	}

	if !IsProviderValid(d.Provider) {
		return nil, errors.New(ErrorInvalidProvider)
	}

	if !IsTypeValid(d.Type) {
		return nil, errors.New(ErrorInvalidType)
	}

	if currentDataset.Credentials != d.Credentials && d.Provider == "ssm" {
		_, err := ssmClient.PutParameter(&ssm.PutParameterInput{DataType: aws.String("text"), Name: aws.String(d.Id), Value: aws.String(d.Credentials), Type: aws.String("SecureString"), Overwrite: aws.Bool(true)})
		if err != nil {
			return nil, errors.New(ErrorCouldNotSecureStoreCredentials)
		}
	}

	// Save dataset
	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}

func DeleteDataset(req events.APIGatewayProxyRequest, dynaClient dynamodbiface.DynamoDBAPI, ssmClient *ssm.SSM) error {
	fmt.Println(req.PathParameters)
	id := req.PathParameters["id"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(TableName),
	}

	currentDataset, _ := FetchDataset(id, dynaClient)
	if currentDataset != nil && len(currentDataset.Name) == 0 {
		return nil
	}

	if currentDataset.Provider == "ssm" {
		ssmClient.DeleteParameter(&ssm.DeleteParameterInput{Name: aws.String(id)})
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
