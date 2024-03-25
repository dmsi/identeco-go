package myhandlers

import (
	"net/http"
)

func renderErr(err error, status int, w http.ResponseWriter) {
	_ = err
	w.WriteHeader(status)
}

func renderOK(body *string, w http.ResponseWriter) {
	if body != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(*body))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
