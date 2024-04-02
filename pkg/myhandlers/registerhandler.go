package myhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmsi/identeco-go/pkg/controllers"
)

type RegisterHandler struct {
	Controller controllers.Controller
}

func (h RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lg := h.Controller.Log

	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user)
	if err != nil {
		renderErr(err, http.StatusBadRequest, w)
		return
	}

	lg.Info("Registering user", "username", user.Username)

	res, err := h.Controller.Register(user.Username, user.Password)
	if err != nil {
		renderErr(err, http.StatusBadRequest, w)
	} else {
		renderOK(res, w)
	}
}
