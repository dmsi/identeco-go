package myhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmsi/identeco-go/pkg/controllers"
)

type LoginHandler struct {
	Controller controllers.Controller
}

func (h LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lg := h.Controller.Log

	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user)
	if err != nil {
		renderErr(err, http.StatusUnauthorized, w)
		return
	}

	lg.Info("Logging in", "username", user.Username)

	res, err := h.Controller.Login(user.Username, user.Password)
	if err != nil {
		renderErr(err, http.StatusBadRequest, w)
	} else {
		renderOK(res, w)
	}
}
