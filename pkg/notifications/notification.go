package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"hermes/pkg/common/crud"
	BaseErrors "hermes/pkg/common/errors"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrorInvalidNotificationData   = "invalid  notification data"
	ErrorNotificationAlreadyExists = "Notification already exists"
	ErrorNotificationDoesNotExists = "Notification does not exist"
)

type Rule struct {
	In  string `json:"in"`
	Op  string `json:"op"`
	Val string `json:"val"`
}

type Template struct {
	Rules []Rule `json:"rules"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
type Query struct {
	DataSetId string `json:"datasetId"`
	Query     string `json:"query"`
}

type Notification struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Templates []Template `json:"templates"`
	Query     `json:"query"`
	Inputs    []string `json:"inputs"`
	Tags      []string `json:"tags"`
}

func FetchNotification(id string, repo crud.CrudRepository) (*Notification, error) {
	item := new(Notification)

	_, err := repo.Get(id, item)

	return item, err
}

func FetchNotifications(repo crud.CrudRepository) (*[]Notification, error) {
	item := new([]Notification)

	_, err := repo.List(item)

	return item, err
}

func CreateNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (*Notification, error) {
	var n Notification
	if err := json.Unmarshal([]byte(req.Body), &n); err != nil {
		return nil, errors.New(ErrorInvalidNotificationData)
	}

	_, err := repo.Create(n)

	if err != nil {
		return nil, err
	}
	return &n, nil
}

func UpdateNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (
	*Notification,
	error,
) {
	var n Notification
	if err := json.Unmarshal([]byte(req.Body), &n); err != nil {
		return nil, errors.New(ErrorInvalidNotificationData)
	}

	// Check if dataset exists
	currentNotification, _ := FetchNotification(n.Id, repo)
	if currentNotification != nil && len(currentNotification.Name) == 0 {
		return nil, errors.New(ErrorNotificationAlreadyExists)
	}

	// Save dataset
	_, err := repo.Update(n.Id, n)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func DeleteNotification(req events.APIGatewayProxyRequest, repo crud.CrudRepository) error {
	id := req.PathParameters["id"]

	currentNotification, _ := FetchNotification(id, repo)
	if currentNotification != nil && len(currentNotification.Name) == 0 {
		return nil
	}

	err := repo.Delete(id)
	if err != nil {
		fmt.Println(err)
		return errors.New(BaseErrors.ErrorCouldNotDeleteItem)
	}

	return nil
}
