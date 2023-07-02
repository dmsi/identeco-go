package token

import (
	"os"

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
	}
}

func refreshTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "refresh",
		"iss":       os.Getenv("ISS_CLAIM"),
	}
}

func IssueTokens(username string) (*Tokens, error) {
	k := keys.NewKeysService()
	privateKey, err := k.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims(username))
	signedAccessToken, err := accessToken.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims(username))
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
