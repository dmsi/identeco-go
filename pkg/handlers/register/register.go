package main

import (
	"fmt"

	"github.com/dmsi/identeco/pkg/token"
)

func main() {
	token, _ := token.IssueToken()
	fmt.Printf("Register! token: %v\n", token)
}
