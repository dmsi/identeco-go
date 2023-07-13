package httpserver

import (
	"encoding/json"
	"net/http"
)

func (h *handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user)
	if err != nil {
		h.errResponse(err, http.StatusBadRequest, w, r)
		return
	}

	res, err := h.register.Register(user.Username, user.Password)
	if err != nil {
		h.errResponse(err, http.StatusBadRequest, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}
