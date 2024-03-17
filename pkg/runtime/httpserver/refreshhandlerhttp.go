package httpserver

import (
	"errors"
	"net/http"
	"strings"
)

func (h *handler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		h.errResponse(errors.New("no authorization header"), http.StatusUnauthorized, w, r)
		return
	}

	refreshToken := strings.Split(auth, " ")[1]
	res, err := h.refresh.Refresh(refreshToken)
	if err != nil {
		h.errResponse(err, http.StatusUnauthorized, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}
