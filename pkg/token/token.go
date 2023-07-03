package token

import (
	"errors"
	"fmt"
	"os"

	"github.com/dmsi/identeco/pkg/jwks"
	"github.com/dmsi/identeco/pkg/keys"
	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func accessTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "access",
		"iss":       os.Getenv("ISS_CLAIM"),
		// "exp":      time.Now().Add(t.ttl).Unix(),
	}
}

func refreshTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "refresh",
		"iss":       os.Getenv("ISS_CLAIM"),
		// "exp":      time.Now().Add(t.ttl).Unix(),
	}
}

func IssueTokens(username string) (*Tokens, error) {
	k := keys.NewKeysService()
	privateKey, err := k.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	jwk, err := jwks.PublicKeyToJWK(privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	kid := jwk.Kid

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims(username))
	accessToken.Header["kid"] = kid
	signedAccessToken, err := accessToken.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims(username))
	refreshToken.Header["kid"] = kid
	signedRefreshToken, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func RefreshToken() (*Tokens, error) {
	return &Tokens{}, nil
}

func VerifyToken(token string) (*jwt.MapClaims, error) {
	fmt.Printf("About to verify\n")
	// Some reference
	// parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	// 	if token.Method.Alg() != t.signMethod.Alg() {
	// 		return nil, errors.New("invalid signing method")
	// 	}
	// 	switch k := t.key.(type) {
	// 	case *rsa.PrivateKey:
	// 		return &k.PublicKey, nil
	// 	case *ecdsa.PrivateKey:
	// 		return &k.PublicKey, nil
	// 	default:
	// 		return t.key, nil
	// 	}
	// })

	// As of PoC - OK, but it will not work in case of keys rotation
	// The public key needs to be taken from JWKS based on the "kid" header!
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected token signing method")
		}

		// privatekey, err := keys.NewKeysService().GetPrivateKey()
		// fmt.Printf("Sending back a public key! -- v2\n")
		// return &privatekey.PublicKey, err

		j, err := keys.NewKeysService().GetJWKS()
		if err != nil {
			return nil, err
		}

		// idx := slices.IndexFunc(j.Keys, func(k jwks.JWK) bool { return k.Kid == t.Header["kid"] })
		// fmt.Printf("Index is %v\n", idx)

		idx := -1
		for i, v := range j.Keys {
			if v.Kid == t.Header["kid"] {
				idx = i
				break
			}
		}
		if idx == -1 {
			return nil, fmt.Errorf("key with kid %v not found", t.Header["kid"])
		}

		pub, err := jwks.JWKToPublicKey(j.Keys[idx])
		if err != nil {
			return nil, err
		}

		return pub, nil
	})

	fmt.Printf("Verified! err: %v\n", err)

	if err != nil {
		return nil, err
	}

	_ = decodedToken
	// decodedToken.Claims

	fmt.Printf("Is this failing herer!!!!\n")

	// TODO is that automatically check if the token expired?
	claims, ok := decodedToken.Claims.(jwt.MapClaims)
	if !ok || !decodedToken.Valid {
		return nil, errors.New("verify token claims failed")
	}

	val, ok := claims["token_use"]
	if !ok || val.(string) != "refresh" {
		return nil, errors.New("verify token claims failed")
	}

	return &claims, nil
}
