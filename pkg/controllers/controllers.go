package controllers

import (
	"log/slog"

	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/dmsi/identeco-go/pkg/services/token"
	"github.com/dmsi/identeco-go/pkg/storage"
)

type Controller struct {
	Log          *slog.Logger
	UserStorage  storage.UsersStorage
	KeyStorage   storage.KeysStorage
	TokenService token.TokenService
	KeyService   keys.KeyService
}

func wrap(name string, err error) error {
	return e.Wrap("controllers."+name, err)
}
