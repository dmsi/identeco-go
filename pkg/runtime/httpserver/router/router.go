package router

import (
	"net/http"

	"github.com/dmsi/identeco/pkg/runtime/httpserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Mux http.Handler
}

func New(mount string, h handlers.Handler) (*Router, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/.well-known/jwks.json", h.JWKSetsHandler)
	r.Post("/register", h.RegisterHandler)
	r.Post("/login", h.LoginHandler)
	r.Get("/refresh", h.RefreshHandler)
	// TODO instead of a route ->
	// It should check the last rotation data upon start and then start the countdown!
	// For that we need to keep track of the last rotation time, somewhere in DB
	r.Get("/rotatekeys", h.RotateKeysHandler)

	m := chi.NewRouter()
	m.Mount(mount, r)

	return &Router{
		Mux: m,
	}, nil
}
