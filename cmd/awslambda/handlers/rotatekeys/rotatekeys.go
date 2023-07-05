package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco/pkg/controllers/rotatekeys"
	"github.com/dmsi/identeco/pkg/runtime/awslambda"
)

var controller *rotatekeys.RotateController

func init() {
	c, err := awslambda.CreateRotateKeysController()
	if err != nil {
		panic("can't create controller")
	}

	controller = c
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := controller.RotateKeys()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func main() {
	lambda.Start(handler)
}
