package awslambda

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
	"github.com/dmsi/identeco/pkg/services/keys"
	"github.com/dmsi/identeco/pkg/services/token"
	"github.com/dmsi/identeco/pkg/storage"
	"github.com/dmsi/identeco/pkg/storage/dynamodb/usersdynamodb"
	"github.com/dmsi/identeco/pkg/storage/s3/keyss3"
	"golang.org/x/exp/slog"
)

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
	stage := os.Getenv("IDO_DEPLOYMENT_STAGE")
	if stage != "prod" {
		lvl = slog.LevelDebug
		src = true
	}

	lg := slog.New(slog.NewJSONHandler(
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
	bits, err := strconv.Atoi(os.Getenv("IDO_PRIVATE_KEY_LENGTH"))
	if err != nil {
		return nil, err
	}

	return &keys.KeyService{
		PrivateKeyBits: bits,
	}, nil
}

func newTokenService(lg *slog.Logger) (*token.TokenService, error) {
	accessTokenLifetime, err := time.ParseDuration(os.Getenv("IDO_ACCESS_TOKEN_LIFETIME"))
	if err != nil {
		return nil, err
	}

	refreshTokenLifetime, err := time.ParseDuration(os.Getenv("IDO_REFRESH_TOKEN_LIFETIME"))
	if err != nil {
		return nil, err
	}

	k, err := newKeyService(lg)
	if err != nil {
		return nil, err
	}

	return &token.TokenService{
		KeyService:           *k,
		Iss:                  os.Getenv("IDO_CLAIM_ISS"),
		AccessTokenLifetime:  accessTokenLifetime,
		RefreshTokenLifetime: refreshTokenLifetime,
	}, nil
}

func newKeyStorage(lg *slog.Logger) (storage.KeysStorage, error) {
	return keyss3.New(
		lg,
		os.Getenv("IDO_BUCKET_NAME"),
		os.Getenv("IDO_PRIVATE_KEY_NAME"),
		os.Getenv("IDO_JWKS_NAME"),
	), nil
}

func newUserStorage(lg *slog.Logger) (storage.UsersStorage, error) {
	return usersdynamodb.New(lg, os.Getenv("IDO_TABLE_NAME")), nil
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

func NewJWKSetsHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &Handler{
		lg:      c.Log,
		jwksets: &jwksets.JWKSetsController{Controller: *c},
	}, nil
}

func NewRegisterHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &Handler{
		lg:       c.Log,
		register: &register.RegisterController{Controller: *c},
	}, nil
}

func NewLoginHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &Handler{
		lg:    c.Log,
		login: &login.LoginController{Controller: *c},
	}, nil
}

func NewRefreshHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &Handler{
		lg:      c.Log,
		refresh: &refresh.RefreshController{Controller: *c},
	}, nil
}

func NewRotateKeysHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, err
	}

	return &Handler{
		lg:         c.Log,
		rotatekeys: &rotatekeys.RotateController{Controller: *c},
	}, nil
}
