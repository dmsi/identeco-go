package controllers

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/services/keys"
)

type RefreshController struct {
	Controller
}

func (r *RefreshController) Refresh(refreshToken string) (*string, error) {
	// Read data
	jwkSetsData, err := r.KeyStorage.ReadJWKSets()
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	jwkSets := keys.JWKSets{}
	err = json.Unmarshal(jwkSetsData.Data, &jwkSets)
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	keyData, err := r.KeyStorage.ReadPrivateKey()
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	privateKey, err := r.KeyService.PrivateKeyDecodePEM(keyData.Data)
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	// Logic
	username, err := r.TokenService.VerifyRefreshToken(refreshToken, jwkSets)
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	_, err = r.UserStorage.ReadUserData(*username)
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	tokens, err := r.TokenService.IssueTokens(*username, privateKey)
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, wrap("Refresh", err)
	}

	return aws.String(string(body)), nil
}
