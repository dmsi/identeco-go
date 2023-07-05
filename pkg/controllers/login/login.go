package login

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco/pkg/controllers"
	"golang.org/x/crypto/bcrypt"

	e "github.com/dmsi/identeco/pkg/lib/err"
)

type LoginController struct {
	controllers.Controller
}

func op(name string) string {
	return "controllers.login." + name
}

func comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (l *LoginController) Login(username, password string) (*string, error) {
	lg := l.Log.With("user", username)

	// Read data
	keyData, err := l.KeyStorage.ReadPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, err := l.KeyService.PrivateKeyDecodePEM(keyData.Data)
	if err != nil {
		return nil, err
	}

	// TODO user not found -> return error
	user, err := l.UserStorage.ReadUserData(username)
	if err != nil {
		lg.Error("read user", "error", err)
		return nil, err
	}

	// Logic
	if !comparePassword(password, user.Hash) {
		lg.Info("invalid password")
		return nil, e.Wrap(op("Login"), errors.New("invalid password"))
	}

	// TODO: store the refresh token
	tokens, err := l.TokenService.IssueTokens(username, privateKey)
	if err != nil {
		lg.Error("issue tokens", "error", err)
		return nil, err
	}

	body, err := json.Marshal(tokens)
	if err != nil {
		return nil, e.Wrap(op("Login"), err)
	}

	return aws.String(string(body)), nil
}
