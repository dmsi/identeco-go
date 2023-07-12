package awslambda

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) JWKSetsHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body, err := h.jwksets.GetJWKSets()
	if err != nil {
		return h.errResponse(err, http.StatusNotFound)
	}

	return h.okResponse(body)
}
