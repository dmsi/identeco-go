package httpserver

import "net/http"

func (h *handler) jwkSetsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.jwksets.GetJWKSets()
	if err != nil {
		h.errResponse(err, http.StatusNotFound, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}
