package datasets

import (
	"fmt"
	"hermes/pkg/common/crud"
	"hermes/pkg/handlers"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var ErrorMethodNotAllowed = "method Not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	id := req.QueryStringParameters["id"]
	if len(id) > 0 {
		// Get single dataset
		result, err := FetchDataset(id, repo)
		if err != nil {
			return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return handlers.ApiResponse(http.StatusOK, result)
	}

	// Get list of datasets
	result, err := FetchDatasets(repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, result)
}

func NewDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository, ssmClient *ssm.SSM) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := CreateDataset(req, repo, ssmClient)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	fmt.Println(result)
	return handlers.ApiResponse(http.StatusCreated, result)
}

func SaveDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository, ssmClient *ssm.SSM) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := UpdateDataset(req, repo, ssmClient)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, result)
}

func RemoveDataset(req events.APIGatewayProxyRequest, repo crud.CrudRepository, ssmClient *ssm.SSM) (
	*events.APIGatewayProxyResponse,
	error,
) {
	err := DeleteDataset(req, repo, ssmClient)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, nil)
}

func TestConnection(req events.APIGatewayProxyRequest, ssmClient *ssm.SSM) (
	*events.APIGatewayProxyResponse,
	error,
) {
	err := EnsureConnection(req, ssmClient)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, nil)
}
