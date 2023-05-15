package keys

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco/pkg/jwks"
	"github.com/dmsi/identeco/pkg/s3helper"
)

// I/O with S3 bucket
type KeyService struct {
	S3         *s3helper.S3Session
	Bucket     string
	JWKSJson   string
	PrivateKey string
}

func (k *KeyService) RotateKeys() error {
	return nil
}

func (k *KeyService) GetJWKS() (jwks.JWKS, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	err := k.S3.ReadObject(k.Bucket, k.JWKSJson, buf)
	if err != nil {
		return jwks.JWKS{}, err
	}

	j := jwks.JWKS{}
	err = json.Unmarshal(buf.Bytes(), &j)
	if err != nil {
		return jwks.JWKS{}, err
	}

	return j, nil
}

// Reads pem from S3 and return as rsa.PrivateKey
func (k *KeyService) GetPrivateKey() (string, error) {
	return "", nil
}
