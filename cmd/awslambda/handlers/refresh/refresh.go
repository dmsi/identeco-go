package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco/pkg/controllers/refresh"
	"github.com/dmsi/identeco/pkg/runtime/awslambda"
)

var controller *refresh.RefreshController

func init() {
	c, err := awslambda.CreateRefreshController()
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
	val, ok := req.Headers["Authorization"]
	if !ok {
		return errResponse(errors.New("no authorization header"))
	}

	refreshToken := strings.Split(val, " ")[1]

	body, err := controller.Refresh(refreshToken)
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
