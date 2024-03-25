package myhandlers

import (
	"net/http"

	"github.com/dmsi/identeco-go/pkg/controllers"
)

type JWKSHandler struct {
	Controller controllers.Controller
}

func (h JWKSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lg := h.Controller.Log
	lg.Info("Getting JWKS")

	res, err := h.Controller.GetJWKS()
	if err != nil {
		renderErr(err, http.StatusNotFound, w)
	} else {
		renderOK(res, w)
	}
}
