package register

import (
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
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &storage.UserData{
		Username: username,
		Hash:     *hash,
	}

	err = r.UserStorage.WriteUserData(*user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
