package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco-go/cmd/awslambda"
	"github.com/dmsi/identeco-go/pkg/myhandlers"
)

var lambdaHandler awslambda.LambdaHandler

func init() {
	log.Println("Login cold start")

	c, err := awslambda.NewController()
	if err != nil {
		log.Fatalf("Unable to create controller: %v", err)
	}

	handlerFn := myhandlers.LoginHandler{Controller: *c}.Handle
	lambdaHandler = awslambda.ChiAdapter(http.MethodPost, "/*", handlerFn)
}

func main() {
	lambda.Start(lambdaHandler)
}
