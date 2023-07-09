package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Mux http.Handler
}

func newRouter(mount string, h handler) (*Router, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/.well-known/jwks.json", h.jwkSetsHandler)
	r.Post("/register", h.registerHandler)
	r.Post("/login", h.loginHandler)
	r.Get("/refresh", h.refreshHandler)
	// TODO instead of a route ->
	// It should check the last rotation data upon start and then start the countdown!
	// For that we need to keep track of the last rotation time, somewhere in DB
	r.Get("/rotatekeys", h.rotateKeysHandler)

	m := chi.NewRouter()
	m.Mount(mount, r)

	return &Router{
		Mux: m,
	}, nil
}
