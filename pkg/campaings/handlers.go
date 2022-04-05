package campaings

import (
	"fmt"
	"hermes/pkg/common/crud"
	"hermes/pkg/handlers"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	id := req.QueryStringParameters["id"]
	if len(id) > 0 {
		result, err := FetchCampaing(id, repo)
		if err != nil {
			return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return handlers.ApiResponse(http.StatusOK, result)
	}

	// Get list of datasets
	result, err := FetchCampaings(repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, result)
}

func NewCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := CreateCampaing(req, repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	fmt.Println(result)
	return handlers.ApiResponse(http.StatusCreated, result)
}

func SaveCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := UpdateCampaing(req, repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, result)
}

func RemoveCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	err := DeleteCampaing(req, repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, nil)
}
