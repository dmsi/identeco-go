package controllers

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/keys"
	"golang.org/x/crypto/bcrypt"
)

func comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (c *Controller) Login(username, password string) (*string, error) {
	// Read data
	keysData, err := c.KeyStorage.ReadKeys()
	if err != nil {
		return nil, err
	}

	privateKey, err := keys.DecodePEM(keysData.PrivateKey)
	if err != nil {
		return nil, err
	}

	// TODO user not found -> return error
	user, err := c.UserStorage.ReadUserData(username)
	if err != nil {
		return nil, err
	}

	// Logic
	if !comparePassword(password, user.Hash) {
		return nil, errors.New("password is not correct")
	}

	tokens, err := c.TokenIssuer.IssueTokens(username, privateKey)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}

	return aws.String(string(body)), nil
}
