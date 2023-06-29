package main

import (
	"fmt"
	"os"

	"github.com/dmsi/identeco/pkg/keys"
	"github.com/dmsi/identeco/pkg/s3helper"
)

func main() {
	k := keys.KeyService{
		S3:                   s3helper.NewS3Session(),
		Bucket:               os.Getenv("BUCKET_NAME"),
		JWKSObjectName:       os.Getenv("JWKS_JSON_NAME"),
		PrivateKeyObjectName: os.Getenv("PRIVATE_KEY_NAME"),
	}

	err := k.RotateKeys()
	fmt.Printf(">>> err %v\n", err)
}
