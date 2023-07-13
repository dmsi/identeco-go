package controllers

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	Controller
}

func comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (c *LoginController) Login(username, password string) (*string, error) {
	lg := c.Log.With("user", username)

	// Read data
	keyData, err := c.KeyStorage.ReadPrivateKey()
	if err != nil {
		return nil, wrap("Login", err)
	}

	privateKey, err := c.KeyService.PrivateKeyDecodePEM(keyData.Data)
	if err != nil {
		return nil, wrap("Login", err)
	}

	// TODO user not found -> return error
	user, err := c.UserStorage.ReadUserData(username)
	if err != nil {
		lg.Error("read user", "error", err)
		return nil, wrap("Login", err)
	}

	// Logic
	if !comparePassword(password, user.Hash) {
		lg.Info("invalid password")
		return nil, wrap("Login", errors.New("invalid password"))
	}

	tokens, err := c.TokenService.IssueTokens(username, privateKey)
	if err != nil {
		lg.Error("issue tokens", "error", err)
		return nil, wrap("Login", err)
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, wrap("Login", err)
	}

	return aws.String(string(body)), nil
}
