package myhandlers

import (
	"net/http"

	"github.com/dmsi/identeco-go/pkg/controllers"
)

type RotateKeysHandler struct {
	Controller controllers.Controller
}

func (h RotateKeysHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lg := h.Controller.Log

	lg.Info("Rotating keys")

	err := h.Controller.RotateKeys()
	if err != nil {
		renderErr(err, http.StatusInternalServerError, w)
	} else {
		renderOK(nil, w)
	}
}
