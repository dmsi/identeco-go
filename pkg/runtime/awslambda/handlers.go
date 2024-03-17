package awslambda

import (
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dmsi/identeco-go/pkg/controllers"
)

type Handler struct {
	lg         *slog.Logger
	jwksets    *controllers.JWKSetsController
	register   *controllers.RegisterController
	login      *controllers.LoginController
	refresh    *controllers.RefreshController
	rotatekeys *controllers.RotateController
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
