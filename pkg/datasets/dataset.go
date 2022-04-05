package datasets

import (
	"encoding/json"
	"errors"
	"fmt"
	"hermes/pkg/common/crud"
	BaseErrors "hermes/pkg/common/errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	TableName                           = os.Getenv("TABLE_NAME")
	ErrorInvalidDatasetData             = "invalid dataset data"
	ErrorInvalidProvider                = "invalid provider"
	ErrorInvalidType                    = "invalid type. Only SQL is supported"
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

func FetchDataset(id string, repo crud.CrudRepository) (*DataSet, error) {
	item := new(DataSet)
	_, err := repo.Get(id, item)

	return item, err
}

func FetchDatasets(repo crud.CrudRepository) (*[]DataSet, error) {
	item := new([]DataSet)
	_, err := repo.Create(item)

	return item, err
}

func CreateDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository, ssmClient *ssm.SSM) (
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

	_, err := repo.Create(d)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

func UpdateDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository, ssmClient *ssm.SSM) (
	*DataSet,
	error,
) {
	var d DataSet
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDatasetData)
	}

	// Check if dataset exists
	currentDataset, _ := FetchDataset(d.Id, repo)
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
	_, err := repo.Update(d.Id, d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func DeleteDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository, ssmClient *ssm.SSM) error {
	id := req.PathParameters["id"]

	currentDataset, _ := FetchDataset(id, repo)
	if currentDataset != nil && len(currentDataset.Name) == 0 {
		return nil
	}

	if currentDataset.Provider == "ssm" {
		ssmClient.DeleteParameter(&ssm.DeleteParameterInput{Name: aws.String(id)})
	}

	err := repo.Delete(id)
	if err != nil {
		fmt.Println(err)
		return errors.New(BaseErrors.ErrorCouldNotDeleteItem)
	}

	return nil
}
