package token

import (
	"crypto/rsa"
	"time"

	"github.com/dmsi/identeco-go/pkg/lib/util"
	"github.com/golang-jwt/jwt/v5"
)

type TokenIssuer struct {
	Iss                  string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

type Tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (t TokenIssuer) accessTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "access",
		"iss":       t.Iss,
		"exp":       time.Now().Add(t.AccessTokenLifetime).Unix(),
	}
}

func (t TokenIssuer) refreshTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "refresh",
		"iss":       t.Iss,
		"exp":       time.Now().Add(t.RefreshTokenLifetime).Unix(),
	}
}

func (t TokenIssuer) signToken(claims jwt.MapClaims, privateKey rsa.PrivateKey) (*string, error) {
	kid := util.KeyID(&privateKey.PublicKey)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessToken.Header["kid"] = kid

	signedAccessToken, err := accessToken.SignedString(&privateKey)
	if err != nil {
		return nil, err
	}

	return &signedAccessToken, nil
}

func (t TokenIssuer) IssueTokens(username string, privateKey *rsa.PrivateKey) (*Tokens, error) {
	accessToken, err := t.signToken(t.accessTokenClaims(username), *privateKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.signToken(t.refreshTokenClaims(username), *privateKey)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
