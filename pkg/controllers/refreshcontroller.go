package controllers

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/services/keys"
)

func (c *Controller) Refresh(refreshToken string) (*string, error) {
	// Read data
	jwkSetsData, err := c.KeyStorage.ReadJWKSets()
	if err != nil {
		return nil, err
	}

	jwkSets := keys.JWKSets{}
	err = json.Unmarshal(jwkSetsData.Data, &jwkSets)
	if err != nil {
		return nil, err
	}

	keyData, err := c.KeyStorage.ReadPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, err := c.KeyService.PrivateKeyDecodePEM(keyData.Data)
	if err != nil {
		return nil, err
	}

	// Logic
	username, err := c.TokenService.VerifyRefreshToken(refreshToken, jwkSets)
	if err != nil {
		return nil, err
	}

	_, err = c.UserStorage.ReadUserData(*username)
	if err != nil {
		return nil, err
	}

	tokens, err := c.TokenService.IssueTokens(*username, privateKey)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}

	return aws.String(string(body)), nil
}
