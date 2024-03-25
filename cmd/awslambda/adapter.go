package awslambda

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/controllers"
	"github.com/dmsi/identeco-go/pkg/mylog"
	"github.com/dmsi/identeco-go/pkg/storageselector"
	"github.com/dmsi/identeco-go/pkg/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewController() (*controllers.Controller, error) {
	lg := mylog.Lg

	userStorage, err := storageselector.NewUsersStorage(lg)
	if err != nil {
		return nil, err
	}

	keyStorage, err := storageselector.NewKeysStorage(lg)
	if err != nil {
		return nil, err
	}

	return &controllers.Controller{
		Log:         lg,
		UserStorage: userStorage,
		KeyStorage:  keyStorage,
		TokenIssuer: token.TokenIssuer{
			Iss:                  config.Cfg.TokenIssClaim,
			AccessTokenLifetime:  config.Cfg.TokenAccessDuration,
			RefreshTokenLifetime: config.Cfg.TokenRefreshDuration,
		},
	}, nil
}

type LambdaHandler func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func ChiAdapter(method string, path string, handlerFn http.HandlerFunc) LambdaHandler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.MethodFunc(method, path, handlerFn)
	adapter := chiadapter.New(r)

	fn := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return adapter.ProxyWithContext(ctx, req)
	}

	return fn
}
