package controllers

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/keys"
	"github.com/dmsi/identeco-go/pkg/lib/util"
	"github.com/golang-jwt/jwt/v5"
)

func verifyRefreshToken(tokenString string, publicKey *rsa.PublicKey) (*string, error) {
	getKey := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		tokenKID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("token has invalid headers: kid is missing")
		}

		keyKID := util.KeyID(publicKey)
		if tokenKID != keyKID {
			return nil, fmt.Errorf("token kid not match key kid, token: %s, key: %s", tokenKID, keyKID)
		}

		return publicKey, nil
	}

	parsed, err := jwt.Parse(tokenString, getKey)
	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token has invalid claims: claims are missing")
	}

	use, ok := claims["token_use"].(string)
	if !ok || use != "refresh" {
		return nil, fmt.Errorf("token has invalid claims: use=%s", use)
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("token has invalid claims: username is missing")
	}

	return &username, nil
}

func (c *Controller) Refresh(refreshToken string) (*string, error) {
	// Read data
	keysData, err := c.KeyStorage.ReadKeys()
	if err != nil {
		return nil, err
	}

	jwkSets := keys.JWKSets{}
	err = json.Unmarshal(keysData.JWKS, &jwkSets)
	if err != nil {
		return nil, err
	}

	privateKey, err := keys.DecodePEM(keysData.PrivateKey)
	if err != nil {
		return nil, err
	}

	username, err := verifyRefreshToken(refreshToken, &privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	_, err = c.UserStorage.ReadUserData(*username)
	if err != nil {
		return nil, err
	}

	tokens, err := c.TokenIssuer.IssueTokens(*username, privateKey)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}

	return aws.String(string(body)), nil
}
