package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dmsi/identeco/pkg/controllers/jwksets"
	"github.com/dmsi/identeco/pkg/controllers/login"
	"github.com/dmsi/identeco/pkg/controllers/refresh"
	"github.com/dmsi/identeco/pkg/controllers/register"
	"github.com/dmsi/identeco/pkg/controllers/rotatekeys"
	"golang.org/x/exp/slog"
)

type Handler struct {
	Log        *slog.Logger
	JWKSets    *jwksets.JWKSetsController
	Register   *register.RegisterController
	Login      *login.LoginController
	Refresh    *refresh.RefreshController
	RotateKeys *rotatekeys.RotateController
}

func (h *Handler) errResponse(err error, status int, w http.ResponseWriter, r *http.Request) {
	h.Log.Error("oops", "error", err)
	w.WriteHeader(status)
}

func (h *Handler) okResponse(body *string, w http.ResponseWriter, r *http.Request) {
	if body != nil {
		w.Write([]byte(*body))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) JWKSetsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.JWKSets.GetJWKSets()
	if err != nil {
		h.errResponse(err, http.StatusNotFound, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.Register.Register(user.Username, user.Password)
	if err != nil {
		h.errResponse(err, http.StatusBadRequest, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.Login.Login(user.Username, user.Password)
	if err != nil {
		h.errResponse(err, http.StatusUnauthorized, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}

func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		h.errResponse(errors.New("no authorization header"), http.StatusUnauthorized, w, r)
		return
	}

	refreshToken := strings.Split(auth, " ")[1]
	res, err := h.Refresh.Refresh(refreshToken)
	if err != nil {
		h.errResponse(err, http.StatusUnauthorized, w, r)
	} else {
		h.okResponse(res, w, r)
	}
}

func (h *Handler) RotateKeysHandler(w http.ResponseWriter, r *http.Request) {
	err := h.RotateKeys.RotateKeys()
	if err != nil {
		h.errResponse(err, http.StatusInternalServerError, w, r)
	} else {
		h.okResponse(nil, w, r)
	}
}
