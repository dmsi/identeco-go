package awslambda

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

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
