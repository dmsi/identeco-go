package awslambda

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

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
