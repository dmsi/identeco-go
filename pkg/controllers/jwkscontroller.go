package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
)

func (c *Controller) GetJWKS() (*string, error) {
	jwkSetsData, err := c.KeyStorage.ReadJWKSets()
	if err != nil {
		return nil, err
	}

	return aws.String(string(jwkSetsData.Data)), nil
}
