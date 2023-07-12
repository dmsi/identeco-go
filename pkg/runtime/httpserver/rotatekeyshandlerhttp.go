package httpserver

import "net/http"

func (h *handler) rotateKeysHandler(w http.ResponseWriter, r *http.Request) {
	err := h.rotatekeys.RotateKeys()
	if err != nil {
		h.errResponse(err, http.StatusInternalServerError, w, r)
	} else {
		h.okResponse(nil, w, r)
	}
}
