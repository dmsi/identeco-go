package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return aws.String(string(hash)), nil
}

func (c *Controller) Register(username, password string) (*string, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &storage.UserData{
		Username: username,
		Hash:     *hash,
	}

	err = c.UserStorage.WriteUserData(*user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
