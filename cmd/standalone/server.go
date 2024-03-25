package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dmsi/identeco-go/pkg/controllers"
	"github.com/dmsi/identeco-go/pkg/myhandlers"
	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/dmsi/identeco-go/pkg/services/token"
	"github.com/dmsi/identeco-go/pkg/storage"
	"github.com/dmsi/identeco-go/pkg/storage/mongodb/keysmongodb"
	"github.com/dmsi/identeco-go/pkg/storage/mongodb/usersmongodb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
)

const (
	port                    = ":3000"
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

func newKeyService() (*keys.KeyService, error) {
	bits, err := strconv.Atoi(os.Getenv(envPrivateKeyLength))
	if err != nil {
		return nil, err
	}

	return &keys.KeyService{
		PrivateKeyBits: bits,
	}, nil
}

func newTokenService() (*token.TokenService, error) {
	accessTokenLifetime, err := time.ParseDuration(os.Getenv(envAccessTokenLifetime))
	if err != nil {
		return nil, err
	}

	refreshTokenLifetime, err := time.ParseDuration(os.Getenv(envRefreshTokenLifetime))
	if err != nil {
		return nil, err
	}

	k, err := newKeyService()
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

	tokenService, err := newTokenService()
	if err != nil {
		return nil, err
	}

	keyService, err := newKeyService()
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

func main() {
	controller, err := newController()
	if err != nil {
		log.Fatalf("Unable to create controller: %s", err)
	}
	api := chi.NewRouter()

	api.Get("/.well-known/jwks.json", myhandlers.JWKSHandler{Controller: *controller}.Handle)
	api.Post("/register", myhandlers.RegisterHandler{Controller: *controller}.Handle)
	api.Post("/login", myhandlers.LoginHandler{Controller: *controller}.Handle)
	api.Get("/refresh", myhandlers.RefreshHandler{Controller: *controller}.Handle)
	api.Get("/rotatekeys", myhandlers.RotateKeysHandler{Controller: *controller}.Handle)

	root := chi.NewRouter()
	root.Use(middleware.RequestID)
	root.Use(middleware.RealIP)
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)
	root.Mount("/api", api)

	err = http.ListenAndServe(port, root)
	if err != nil {
		log.Fatalf("Unable to start server: %s", err)
	}

	log.Printf("Started server at [::]:%s", port)
}
