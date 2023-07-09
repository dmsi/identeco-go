package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco/pkg/runtime/awslambda"
)

var handler *awslambda.Handler

func init() {
	h, err := awslambda.NewRotateKeysHandler()
	if err != nil {
		panic("can't create handler")
	}

	handler = h
}

func main() {
	lambda.Start(handler.RotateKeysHandler)
}
