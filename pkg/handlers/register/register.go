package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dmsi/identeco/pkg/token"
	"github.com/dmsi/identeco/pkg/user"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	creds, err := user.GetCredentialsFromString(req.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	err = user.AddUser(*creds)
	var ae *dynamodb.ConditionalCheckFailedException
	if errors.As(err, &ae) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("User %v already exists", creds.Username),
		}, nil
	}

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	tokens, err := token.IssueTokens(creds.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
