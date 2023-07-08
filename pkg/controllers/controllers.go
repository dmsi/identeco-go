package controllers

import (
	"github.com/dmsi/identeco/pkg/services/keys"
	"github.com/dmsi/identeco/pkg/services/token"
	"github.com/dmsi/identeco/pkg/storage"
	"golang.org/x/exp/slog"
)

type Controller struct {
	Log          *slog.Logger
	UserStorage  storage.UsersStorage
	KeyStorage   storage.KeysStorage
	TokenService token.TokenService
	KeyService   keys.KeyService
}
