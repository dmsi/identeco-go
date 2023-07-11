package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dmsi/identeco-go/pkg/controllers/jwksets"
	"github.com/dmsi/identeco-go/pkg/controllers/login"
	"github.com/dmsi/identeco-go/pkg/controllers/refresh"
	"github.com/dmsi/identeco-go/pkg/controllers/register"
	"github.com/dmsi/identeco-go/pkg/controllers/rotatekeys"
	"golang.org/x/exp/slog"
)

type handler struct {
	lg         *slog.Logger
	jwksets    *jwksets.JWKSetsController
	register   *register.RegisterController
	login      *login.LoginController
	refresh    *refresh.RefreshController
	rotatekeys *rotatekeys.RotateController
}

func (h *handler) errResponse(err error, status int, w http.ResponseWriter, r *http.Request) {
	h.lg.Error("oops", "error", err)
	w.WriteHeader(status)
}

func (h *handler) okResponse(body *string, w http.ResponseWriter, r *http.Request) {
	if body != nil {
		w.Write([]byte(*body))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *handler) jwkSetsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.jwksets.GetJWKSets()
	if err != nil {
		h.errResponse(err, http.StatusNotFound, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}

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

func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user)
	if err != nil {
		h.errResponse(err, http.StatusUnauthorized, w, r)
		return
	}

	res, err := h.login.Login(user.Username, user.Password)
	if err != nil {
		h.errResponse(err, http.StatusUnauthorized, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}

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

func (h *handler) rotateKeysHandler(w http.ResponseWriter, r *http.Request) {
	err := h.rotatekeys.RotateKeys()
	if err != nil {
		h.errResponse(err, http.StatusInternalServerError, w, r)
	} else {
		h.okResponse(nil, w, r)
	}
}
