package main

import (
	"hermes/pkg/campaings"
	"hermes/pkg/common/crud"
	"hermes/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	TableName  = os.Getenv("TABLE_NAME")
	dynaClient dynamodbiface.DynamoDBAPI
	ssmClient  *ssm.SSM
	repo       crud.CrudRepository
)

func getAwsSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	isDev := os.Getenv("IS_DEV")

	if isDev == "true" {
		return session.NewSession(&aws.Config{
			Region:   aws.String(region),
			Endpoint: aws.String("http://host.docker.internal:4566"),
		},
		)
	}

	return session.NewSession(&aws.Config{
		Region: aws.String(region),
	},
	)
}

func main() {
	awsSession, err := getAwsSession()

	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)
	repo = crud.InitDynamoDbRepo(TableName, dynaClient)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return campaings.GetCampaing(req, repo)
	case "POST":
		return campaings.NewCampaing(req, repo)
	case "PUT":
		return campaings.SaveCampaing(req, repo)
	case "DELETE":
		return campaings.RemoveCampaing(req, repo)
	default:
		return handlers.UnhandledMethod()
	}
}
