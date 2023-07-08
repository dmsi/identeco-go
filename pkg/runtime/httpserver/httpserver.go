package httpserver

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dmsi/identeco/pkg/controllers"
	"github.com/dmsi/identeco/pkg/controllers/jwksets"
	"github.com/dmsi/identeco/pkg/controllers/login"
	"github.com/dmsi/identeco/pkg/controllers/refresh"
	"github.com/dmsi/identeco/pkg/controllers/register"
	"github.com/dmsi/identeco/pkg/controllers/rotatekeys"
	"github.com/dmsi/identeco/pkg/runtime/httpserver/handlers"
	"github.com/dmsi/identeco/pkg/runtime/httpserver/router"
	"github.com/dmsi/identeco/pkg/services/keys"
	"github.com/dmsi/identeco/pkg/services/token"
	"github.com/dmsi/identeco/pkg/storage"
	"github.com/dmsi/identeco/pkg/storage/mongodb/keysmongodb"
	"github.com/dmsi/identeco/pkg/storage/mongodb/usersmongodb"
	"golang.org/x/exp/slog"
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

type Runtime struct {
	Log    *slog.Logger
	Router router.Router
}

func newHandler() (*handlers.Handler, error) {
	return nil, nil
}

func newRouter() (*router.Router, error) {
	return nil, nil
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

func newKeyService(lg *slog.Logger) (*keys.KeyService, error) {
	bits, err := strconv.Atoi(os.Getenv(envPrivateKeyLength))
	if err != nil {
		return nil, err
	}

	return &keys.KeyService{
		PrivateKeyBits: bits,
	}, nil
}

func newTokenService(lg *slog.Logger) (*token.TokenService, error) {
	accessTokenLifetime, err := time.ParseDuration(os.Getenv(envAccessTokenLifetime))
	if err != nil {
		return nil, err
	}

	refreshTokenLifetime, err := time.ParseDuration(os.Getenv(envRefreshTokenLifetime))
	if err != nil {
		return nil, err
	}

	k, err := newKeyService(lg)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return k, nil
}

func newUserStorage(lg *slog.Logger) (storage.UsersStorage, error) {
	u, err := usersmongodb.New(lg, os.Getenv(envMongoURL), mongoDatabase, mongoUsersCollection)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func newController() (*controllers.Controller, error) {
	lg := newLogger()

	userStorage, err := newUserStorage(lg)
	if err != nil {
		return nil, err
	}

	keyStorage, err := newKeyStorage(lg)
	if err != nil {
		return nil, err
	}

	tokenService, err := newTokenService(lg)
	if err != nil {
		return nil, err
	}

	keyService, err := newKeyService(lg)
	if err != nil {
		return nil, err
	}

	return &controllers.Controller{
		Log:          lg,
		UserStorage:  userStorage,
		KeyStorage:   keyStorage,
		TokenService: *tokenService,
		KeyService:   *keyService,
	}, nil
}

func NewRuntime(mount string) (*Runtime, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	h := handlers.Handler{
		Log:        c.Log,
		JWKSets:    &jwksets.JWKSetsController{Controller: *c},
		Register:   &register.RegisterController{Controller: *c},
		Login:      &login.LoginController{Controller: *c},
		Refresh:    &refresh.RefreshController{Controller: *c},
		RotateKeys: &rotatekeys.RotateController{Controller: *c},
	}

	r, err := router.New(mount, h)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Log:    c.Log,
		Router: *r,
	}, nil
}

/*
func CreateLoginController() (*login.LoginController, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &login.LoginController{
		Controller: *c,
	}, nil
}

func CreateRefreshController() (*refresh.RefreshController, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &refresh.RefreshController{
		Controller: *c,
	}, nil
}

func CreateRegisterController() (*register.RegisterController, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &register.RegisterController{
		Controller: *c,
	}, nil
}

func CreateRotateKeysController() (*rotatekeys.RotateController, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &rotatekeys.RotateController{
		Controller: *c,
	}, nil
}

func CreateJwkSetsController() (*jwksets.JWKSetsController, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &jwksets.JWKSetsController{
		Controller: *c,
	}, nil
}
*/
