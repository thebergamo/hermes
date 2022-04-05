package campaings

import (
	"encoding/json"
	"errors"
	"fmt"
	"hermes/pkg/common/crud"
	BaseErrors "hermes/pkg/common/errors"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrorInvalidCampaingData   = "invalid  notification data"
	ErrorCampaingAlreadyExists = "Campaing already exists"
	ErrorCampaingDoesNotExists = "Campaing does not exist"
)

type Agenda struct {
	Cron           string `json:"cron"`
	NotificationId string `json:"notificationId"`
	Seq            int    `json:"seq"`
}

type Campaing struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Filter string   `json:"filter"`
	Agenda []Agenda `json:"agenda"`
	Tags   []string `json:"tags"`
}

func FetchCampaing(id string, repo crud.CrudRepository) (*Campaing, error) {
	item := new(Campaing)

	_, err := repo.Get(id, item)

	return item, err
}

func FetchCampaings(repo crud.CrudRepository) (*[]Campaing, error) {
	item := new([]Campaing)

	_, err := repo.List(item)

	return item, err
}

func CreateCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (*Campaing, error) {
	var n Campaing
	if err := json.Unmarshal([]byte(req.Body), &n); err != nil {
		return nil, errors.New(ErrorInvalidCampaingData)
	}

	_, err := repo.Create(n)

	if err != nil {
		return nil, err
	}
	return &n, nil
}

func UpdateCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) (*Campaing, error) {
	var n Campaing
	if err := json.Unmarshal([]byte(req.Body), &n); err != nil {
		return nil, errors.New(ErrorInvalidCampaingData)
	}

	currentCampaing, _ := FetchCampaing(n.Id, repo)
	if currentCampaing != nil && len(currentCampaing.Name) == 0 {
		return nil, errors.New(ErrorCampaingAlreadyExists)
	}

	_, err := repo.Update(n.Id, n)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func DeleteCampaing(req events.APIGatewayProxyRequest, repo crud.CrudRepository) error {
	id := req.PathParameters["id"]

	currentNotification, _ := FetchCampaing(id, repo)
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
