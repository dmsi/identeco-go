package register

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/controllers"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

func wrap(name string, err error) error {
	return e.Wrap("controllers.register."+name, err)
}

type RegisterController struct {
	controllers.Controller
}

func hashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, wrap("hashPassword", err)
	}

	return aws.String(string(hash)), nil
}

func (r *RegisterController) Register(username, password string) (*string, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, wrap("Register", err)
	}

	user := &storage.UserData{
		Username: username,
		Hash:     *hash,
	}

	err = r.UserStorage.WriteUserData(*user)
	if err != nil {
		return nil, wrap("Register", err)
	}

	return nil, nil
}
