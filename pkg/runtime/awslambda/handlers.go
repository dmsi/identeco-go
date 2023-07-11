package awslambda

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dmsi/identeco-go/pkg/controllers/jwksets"
	"github.com/dmsi/identeco-go/pkg/controllers/login"
	"github.com/dmsi/identeco-go/pkg/controllers/refresh"
	"github.com/dmsi/identeco-go/pkg/controllers/register"
	"github.com/dmsi/identeco-go/pkg/controllers/rotatekeys"
	"golang.org/x/exp/slog"
)

type Handler struct {
	lg         *slog.Logger
	jwksets    *jwksets.JWKSetsController
	register   *register.RegisterController
	login      *login.LoginController
	refresh    *refresh.RefreshController
	rotatekeys *rotatekeys.RotateController
}

func (h *Handler) errResponse(err error, status int) (events.APIGatewayProxyResponse, error) {
	h.lg.Error("oops", "error", err)

	return events.APIGatewayProxyResponse{
		StatusCode: status,
	}, nil
}

func (h *Handler) okResponse(body *string) (events.APIGatewayProxyResponse, error) {
	if body == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: *body,
	}, nil
}

func (h *Handler) JWKSetsHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body, err := h.jwksets.GetJWKSets()
	if err != nil {
		return h.errResponse(err, http.StatusNotFound)
	}

	return h.okResponse(body)
}

func (h *Handler) RegisterHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	creds := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal([]byte(req.Body), creds)
	if err != nil {
		return h.errResponse(err, http.StatusBadRequest)
	}

	body, err := h.register.Register(creds.Username, creds.Password)
	if err != nil {
		return h.errResponse(err, http.StatusBadRequest)
	}

	return h.okResponse(body)
}

func (h *Handler) LoginHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	creds := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal([]byte(req.Body), creds)
	if err != nil {
		return h.errResponse(err, http.StatusUnauthorized)
	}

	body, err := h.login.Login(creds.Username, creds.Password)
	if err != nil {
		return h.errResponse(err, http.StatusUnauthorized)
	}

	return h.okResponse(body)
}

func (h *Handler) RefreshHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	val, ok := req.Headers["Authorization"]
	if !ok {
		return h.errResponse(errors.New("no authorization header"), http.StatusUnauthorized)
	}

	refreshToken := strings.Split(val, " ")[1]

	body, err := h.refresh.Refresh(refreshToken)
	if err != nil {
		return h.errResponse(err, http.StatusUnauthorized)
	}

	return h.okResponse(body)
}

func (h *Handler) RotateKeysHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := h.rotatekeys.RotateKeys()
	if err != nil {
		return h.errResponse(err, http.StatusInternalServerError)
	}

	return h.okResponse(nil)
}
