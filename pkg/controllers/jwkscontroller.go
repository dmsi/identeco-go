package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
)

type JWKSetsController struct {
	Controller
}

func (c *JWKSetsController) GetJWKSets() (*string, error) {
	jwkSetsData, err := c.KeyStorage.ReadJWKSets()
	if err != nil {
		return nil, wrap("GetJWKSets", err)
	}

	return aws.String(string(jwkSetsData.Data)), nil
}
