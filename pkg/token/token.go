package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token = string

func keepImportedLol() {
	_ = jwt.ClaimStrings{}
}

func IssueToken() (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))

	fmt.Printf("ztoken: %v, err: %v\n", tokenString, err)
	// jwt.SigningMethodRS256.Alg()

	return "<token placeholder>", nil
}

func RefreshToken() (Token, error) {
	return "<token placeholder>", nil
}
