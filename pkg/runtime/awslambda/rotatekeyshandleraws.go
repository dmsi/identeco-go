package awslambda

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) RotateKeysHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := h.rotatekeys.RotateKeys()
	if err != nil {
		return h.errResponse(err, http.StatusInternalServerError)
	}

	return h.okResponse(nil)
}
