package controllers

import (
	"log/slog"

	"github.com/dmsi/identeco-go/pkg/storage"
	"github.com/dmsi/identeco-go/pkg/token"
)

type Controller struct {
	Log         *slog.Logger
	UserStorage storage.UsersStorage
	KeyStorage  storage.KeysStorage
	TokenIssuer token.TokenIssuer
}
