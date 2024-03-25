package myhandlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dmsi/identeco-go/pkg/controllers"
)

type RefreshHandler struct {
	Controller controllers.Controller
}

func (h RefreshHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lg := h.Controller.Log

	auth := r.Header.Get("Authorization")
	if auth == "" {
		renderErr(errors.New("not authorized"), http.StatusUnauthorized, w)
		return
	}

	lg.Info("Refreshing token")

	refreshToken := strings.Split(auth, " ")[1]
	res, err := h.Controller.Refresh(refreshToken)
	if err != nil {
		lg.Info("Unable to refresh token", "err", err)
		renderErr(err, http.StatusUnauthorized, w)
	} else {
		renderOK(res, w)
	}
}
