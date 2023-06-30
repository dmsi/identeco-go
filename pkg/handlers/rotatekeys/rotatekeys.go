package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco/pkg/keys"
	"github.com/dmsi/identeco/pkg/s3helper"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	k := keys.KeyService{
		S3:                   s3helper.NewS3Session(),
		Bucket:               os.Getenv("BUCKET_NAME"),
		JWKSObjectName:       os.Getenv("JWKS_JSON_NAME"),
		PrivateKeyObjectName: os.Getenv("PRIVATE_KEY_NAME"),
	}

	err := k.RotateKeys()
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
