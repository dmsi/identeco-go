package myhandlers

import (
	"net/http"

	"github.com/dmsi/identeco-go/pkg/mylog"
)

func renderErr(err error, status int, w http.ResponseWriter) {
	mylog.Lg.Error("Handler error", "err", err, "status", status)
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
