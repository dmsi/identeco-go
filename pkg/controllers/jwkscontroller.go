package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
)

func (c *Controller) GetJWKS() (*string, error) {
	keysData, err := c.KeyStorage.ReadKeys()
	if err != nil {
		return nil, err
	}

	return aws.String(string(keysData.JWKS)), nil
}
