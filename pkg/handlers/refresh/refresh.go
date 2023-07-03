package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco/pkg/token"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Doint refresh dont panic!\n")
	val, ok := req.Headers["Authorization"]
	if !ok {
		return events.APIGatewayProxyResponse{}, errors.New("no authorization header")
	}

	refreshToken := strings.Split(val, " ")[1]
	_ = refreshToken

	verified, err := token.VerifyToken(refreshToken)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Printf(">>>>> token verified %v\n", *verified)
	username := (*verified)["username"].(string)
	fmt.Printf("username >>>>> %v\n", username)

	// username := verified.Claims.

	// 1. Decode token and extract username
	// 2. Verify token signature
	// 3. Verify userjjjj
	// 4. Issue new tokens

	// fmt.Printf("HEADERS: %v\n", req.Headers)

	// creds, err := user.GetCredentialsFromString(req.Body)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	// err = user.AddUser(*creds)
	// var ae *dynamodb.ConditionalCheckFailedException
	// if errors.As(err, &ae) {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Body:       fmt.Sprintf("User %v already exists", creds.Username),
	// 	}, nil
	// }

	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	tokens, err := token.IssueTokens(username)
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
