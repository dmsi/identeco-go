package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

type RegisterController struct {
	Controller
}

func hashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, wrap("hashPassword", err)
	}

	return aws.String(string(hash)), nil
}

func (c *RegisterController) Register(username, password string) (*string, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, wrap("Register", err)
	}

	user := &storage.UserData{
		Username: username,
		Hash:     *hash,
	}

	err = c.UserStorage.WriteUserData(*user)
	if err != nil {
		return nil, wrap("Register", err)
	}

	return nil, nil
}
