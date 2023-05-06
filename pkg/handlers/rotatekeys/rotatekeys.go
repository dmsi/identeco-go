package main

import (
	"fmt"

	"github.com/dmsi/identeco/pkg/keys"
)

func main() {
	jwks, _ := keys.GetJWKS()
	fmt.Printf("RotateKeys: %s\n", jwks)
}
