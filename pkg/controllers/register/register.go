package register

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/controllers"
	"github.com/dmsi/identeco-go/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

type RegisterController struct {
	controllers.Controller
}

func hashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return aws.String(string(hash)), nil
}

func (r *RegisterController) Register(username, password string) (*string, error) {
	// Read data
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	keyData, err := r.KeyStorage.ReadPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, err := r.KeyService.PrivateKeyDecodePEM(keyData.Data)
	if err != nil {
		return nil, err
	}

	// Logic
	user := &storage.UserData{
		Username: username,
		Hash:     *hash,
	}

	err = r.UserStorage.WriteUserData(*user)
	if err != nil {
		return nil, err
	}

	tokens, err := r.TokenService.IssueTokens(username, privateKey)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}

	return aws.String(string(body)), nil
}
