package datasets

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	_ "github.com/lib/pq"
)

type Connection struct {
	Credentials string `json:"credentials"`
	Type        string `json:"type"`
	Provider    string `json:"provider"`
}

var (
	ErrorInvalidConnectionData             = "invalid connection data"
	ErrorInvalidConnectionCredentials      = "invalid connection credentials"
	ErrorUnableToPing                      = "unable to ping connect. check credentials and access"
	ErrorCouldNotSecureRetrieveCredentials = "could not retrieve securely credentials from SSM"
)

func EnsureConnection(req events.APIGatewayProxyRequest, ssmClient *ssm.SSM) error {
	var c Connection
	if err := json.Unmarshal([]byte(req.Body), &c); err != nil {
		return errors.New(ErrorInvalidConnectionData)
	}
	parsedCrendentials := c.Credentials

	if c.Provider == "ssm" {
		credentials, err := ssmClient.GetParameter(&ssm.GetParameterInput{Name: aws.String(c.Credentials), WithDecryption: aws.Bool(true)})
		if err != nil {
			return errors.New(ErrorCouldNotSecureRetrieveCredentials)
		}
		parsedCrendentials = *credentials.Parameter.Value
	}

	db, err := sql.Open("postgres", parsedCrendentials)
	if err != nil {
		return errors.New(ErrorInvalidConnectionCredentials)
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorUnableToPing)
	}

	return nil
}
