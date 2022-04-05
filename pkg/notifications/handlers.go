package notifications

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

func GetNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	id := req.QueryStringParameters["id"]
	if len(id) > 0 {
		result, err := FetchNotification(id, repo)
		if err != nil {
			return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return handlers.ApiResponse(http.StatusOK, result)
	}

	// Get list of datasets
	result, err := FetchNotifications(repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, result)
}

func NewNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := CreateNotification(req, repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	fmt.Println(result)
	return handlers.ApiResponse(http.StatusCreated, result)
}

func SaveNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := UpdateNotification(req, repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, result)
}

func RemoveNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*events.APIGatewayProxyResponse,
	error,
) {
	err := DeleteNotification(req, repo)
	if err != nil {
		return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return handlers.ApiResponse(http.StatusOK, nil)
}
