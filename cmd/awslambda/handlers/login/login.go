package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco/pkg/controllers/login"
	"github.com/dmsi/identeco/pkg/runtime/awslambda"
)

var controller *login.LoginController

func init() {
	c, err := awslambda.CreateLoginController()
	if err != nil {
		panic("can't create controller")
	}

	controller = c
}

func errResponse(err error) (events.APIGatewayProxyResponse, error) {
	controller.Log.Error("login failed", "error", err)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusUnauthorized,
	}, nil
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	creds := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal([]byte(req.Body), creds)
	if err != nil {
		return errResponse(err)
	}

	body, err := controller.Login(creds.Username, creds.Password)
	if err != nil {
		return errResponse(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: *body,
	}, nil
}

func main() {
	lambda.Start(handler)
}
