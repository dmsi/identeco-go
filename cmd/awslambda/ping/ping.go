package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dmsi/identeco-go/cmd/awslambda"
	"github.com/go-chi/chi/v5"
)

var lambdaHandler awslambda.LambdaHandler

func handlerFn(w http.ResponseWriter, r *http.Request) {
	c := chi.RouteContext(r.Context())
	id := c.URLParam("id")
	x := c.URLParam("x")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("PONG! id=%s, x=%s", id, x)))
}

func init() {
	log.Println("Ping cold start")
	lambdaHandler = awslambda.ChiAdapter(http.MethodGet, "/{_}/{id}", handlerFn)
}

func main() {
	lambda.Start(lambdaHandler)
}
