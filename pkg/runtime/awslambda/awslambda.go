package awslambda

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dmsi/identeco-go/pkg/controllers"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/dmsi/identeco-go/pkg/services/token"
	"github.com/dmsi/identeco-go/pkg/storage"
	"github.com/dmsi/identeco-go/pkg/storage/dynamodb/usersdynamodb"
	"github.com/dmsi/identeco-go/pkg/storage/s3/keyss3"
	"golang.org/x/exp/slog"
)

const (
	envStage                = "IDO_DEPLOYMENT_STAGE"
	envPrivateKeyLength     = "IDO_PRIVATE_KEY_LENGTH"
	envAccessTokenLifetime  = "IDO_ACCESS_TOKEN_LIFETIME"
	envRefreshTokenLifetime = "IDO_REFRESH_TOKEN_LIFETIME"
	envIssClaim             = "IDO_CLAIM_ISS"
	envTableName            = "IDO_TABLE_NAME"
	envBucketName           = "IDO_BUCKET_NAME"
	envPrivateKeyName       = "IDO_PRIVATE_KEY_NAME"
	envJWKSetsName          = "IDO_JWKS_NAME"
)

func wrap(name string, err error) error {
	return e.Wrap("runtime.awslambda."+name, err)
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
	return keyss3.New(
		lg,
		os.Getenv(envBucketName),
		os.Getenv(envPrivateKeyName),
		os.Getenv(envJWKSetsName),
	), nil
}

func newUserStorage(lg *slog.Logger) (storage.UsersStorage, error) {
	return usersdynamodb.New(lg, os.Getenv(envTableName)), nil
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

func NewJWKSetsHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, wrap("NewJWKSetsHandler", err)
	}

	return &Handler{
		lg:      c.Log,
		jwksets: &controllers.JWKSetsController{Controller: *c},
	}, nil
}

func NewRegisterHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, wrap("NewRegisterHandler", err)
	}

	return &Handler{
		lg:       c.Log,
		register: &controllers.RegisterController{Controller: *c},
	}, nil
}

func NewLoginHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, wrap("NewLoginHandler", err)
	}

	return &Handler{
		lg:    c.Log,
		login: &controllers.LoginController{Controller: *c},
	}, nil
}

func NewRefreshHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, wrap("NewRefreshHandler", err)
	}

	return &Handler{
		lg:      c.Log,
		refresh: &controllers.RefreshController{Controller: *c},
	}, nil
}

func NewRotateKeysHandler() (*Handler, error) {
	c, err := newController()
	if err != nil {
		return nil, wrap("NewRotateKeysHandler", err)
	}

	return &Handler{
		lg:         c.Log,
		rotatekeys: &controllers.RotateController{Controller: *c},
	}, nil
}
