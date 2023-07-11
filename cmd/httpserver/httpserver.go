package main

import (
	"net/http"

	"github.com/dmsi/identeco-go/pkg/runtime/httpserver"
	_ "github.com/joho/godotenv/autoload"
)

const (
	address = ":3000"
)

func main() {
	router, err := httpserver.NewRouter("/ido")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(address, router.Mux)
	if err != nil {
		panic(err)
	}
}
