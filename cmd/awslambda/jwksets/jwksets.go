package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco-go/pkg/runtime/awslambda"
)

var handler *awslambda.Handler

func init() {
	h, err := awslambda.NewJWKSetsHandler()
	if err != nil {
		panic("can't create handler")
	}

	handler = h
}

func main() {
	lambda.Start(handler.JWKSetsHandler)
}
