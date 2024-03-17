package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/dmsi/identeco-go/pkg/controllers"
)

type handler struct {
	lg         *slog.Logger
	jwksets    *controllers.JWKSetsController
	register   *controllers.RegisterController
	login      *controllers.LoginController
	refresh    *controllers.RefreshController
	rotatekeys *controllers.RotateController
}

func (h *handler) errResponse(err error, status int, w http.ResponseWriter, _ *http.Request) {
	h.lg.Error("oops", "error", err)
	w.WriteHeader(status)
}

func (h *handler) okResponse(body *string, w http.ResponseWriter, _ *http.Request) {
	if body != nil {
		w.Write([]byte(*body))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
