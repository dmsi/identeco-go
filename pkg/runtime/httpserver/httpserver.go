package httpserver

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dmsi/identeco-go/pkg/controllers"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/dmsi/identeco-go/pkg/services/token"
	"github.com/dmsi/identeco-go/pkg/storage"
	"github.com/dmsi/identeco-go/pkg/storage/mongodb/keysmongodb"
	"github.com/dmsi/identeco-go/pkg/storage/mongodb/usersmongodb"
)

const (
	mongoDatabase           = "main"
	mongoUsersCollection    = "users"
	mongoKeysCollection     = "keys"
	envMongoURL             = "IDO_MONGODB_URL"
	envStage                = "IDO_DEPLOYMENT_STAGE"
	envPrivateKeyLength     = "IDO_PRIVATE_KEY_LENGTH"
	envAccessTokenLifetime  = "IDO_ACCESS_TOKEN_LIFETIME"
	envRefreshTokenLifetime = "IDO_REFRESH_TOKEN_LIFETIME"
	envIssClaim             = "IDO_CLAIM_ISS"
)

func wrap(name string, err error) error {
	return e.Wrap("runtime.httpserver."+name, err)
}

func newLogger() *slog.Logger {
	// Remove the directory from the source's filename.
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	lvl := slog.LevelInfo
	src := false
	stage := os.Getenv(envStage)
	if stage != "prod" {
		lvl = slog.LevelDebug
		src = true
	}

	lg := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource:   src,
			ReplaceAttr: replace,
			Level:       lvl,
		}),
	)

	return lg
}

func newKeyService(_ *slog.Logger) (*keys.KeyService, error) {
	bits, err := strconv.Atoi(os.Getenv(envPrivateKeyLength))
	if err != nil {
		return nil, wrap("newKeyService", err)
	}

	return &keys.KeyService{
		PrivateKeyBits: bits,
	}, nil
}

func newTokenService(lg *slog.Logger) (*token.TokenService, error) {
	accessTokenLifetime, err := time.ParseDuration(os.Getenv(envAccessTokenLifetime))
	if err != nil {
		return nil, wrap("newTokenService", err)
	}

	refreshTokenLifetime, err := time.ParseDuration(os.Getenv(envRefreshTokenLifetime))
	if err != nil {
		return nil, wrap("newTokenService", err)
	}

	k, err := newKeyService(lg)
	if err != nil {
		return nil, wrap("newTokenService", err)
	}

	return &token.TokenService{
		KeyService:           *k,
		Iss:                  os.Getenv(envIssClaim),
		AccessTokenLifetime:  accessTokenLifetime,
		RefreshTokenLifetime: refreshTokenLifetime,
	}, nil
}

func newKeyStorage(lg *slog.Logger) (storage.KeysStorage, error) {
	k, err := keysmongodb.New(lg, os.Getenv(envMongoURL), "main", "keys")
	if err != nil {
		return nil, wrap("newKeyStorage", err)
	}

	return k, nil
}

func newUserStorage(lg *slog.Logger) (storage.UsersStorage, error) {
	u, err := usersmongodb.New(lg, os.Getenv(envMongoURL), mongoDatabase, mongoUsersCollection)
	if err != nil {
		return nil, wrap("newUserStorage", err)
	}

	return u, nil
}

func newController() (*controllers.Controller, error) {
	lg := newLogger()

	userStorage, err := newUserStorage(lg)
	if err != nil {
		return nil, wrap("newController", err)
	}

	keyStorage, err := newKeyStorage(lg)
	if err != nil {
		return nil, wrap("newController", err)
	}

	tokenService, err := newTokenService(lg)
	if err != nil {
		return nil, wrap("newController", err)
	}

	keyService, err := newKeyService(lg)
	if err != nil {
		return nil, wrap("newController", err)
	}

	return &controllers.Controller{
		Log:          lg,
		UserStorage:  userStorage,
		KeyStorage:   keyStorage,
		TokenService: *tokenService,
		KeyService:   *keyService,
	}, nil
}

func NewRouter(mount string) (*Router, error) {
	c, err := newController()
	if err != nil {
		return nil, wrap("NewRouter", err)
	}

	h := handler{
		lg:         c.Log,
		jwksets:    &controllers.JWKSetsController{Controller: *c},
		register:   &controllers.RegisterController{Controller: *c},
		login:      &controllers.LoginController{Controller: *c},
		refresh:    &controllers.RefreshController{Controller: *c},
		rotatekeys: &controllers.RotateController{Controller: *c},
	}

	r, err := newRouter(mount, h)
	if err != nil {
		return nil, wrap("NewRouter", err)
	}

	return r, nil
}
