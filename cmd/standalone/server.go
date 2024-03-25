package main

import (
	"log"
	"net/http"

	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/controllers"
	"github.com/dmsi/identeco-go/pkg/myhandlers"
	"github.com/dmsi/identeco-go/pkg/mylog"
	"github.com/dmsi/identeco-go/pkg/storageselector"
	"github.com/dmsi/identeco-go/pkg/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	port = ":3000"
)

func newController() (*controllers.Controller, error) {
	lg := mylog.Lg

	userStorage, err := storageselector.NewUsersStorage(lg)
	if err != nil {
		return nil, err
	}

	keyStorage, err := storageselector.NewKeysStorage(lg)
	if err != nil {
		return nil, err
	}

	return &controllers.Controller{
		Log:         lg,
		UserStorage: userStorage,
		KeyStorage:  keyStorage,
		TokenIssuer: token.TokenIssuer{
			Iss:                  config.Cfg.TokenIssClaim,
			AccessTokenLifetime:  config.Cfg.TokenAccessDuration,
			RefreshTokenLifetime: config.Cfg.TokenRefreshDuration,
		},
	}, nil
}

func main() {
	lg := mylog.Lg
	lg.Info("Config", "config", config.Cfg)

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
