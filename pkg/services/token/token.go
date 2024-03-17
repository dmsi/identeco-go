package token

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	KeyService           keys.KeyService
	Iss                  string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

type Tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func wrap(name string, err error) error {
	return e.Wrap("services.token."+name, err)
}

func (t *TokenService) accessTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "access",
		"iss":       t.Iss,
		"exp":       time.Now().Add(t.AccessTokenLifetime).Unix(),
	}
}

func (t *TokenService) refreshTokenClaims(username string) jwt.MapClaims {
	return jwt.MapClaims{
		"username":  username,
		"token_use": "refresh",
		"iss":       t.Iss,
		"exp":       time.Now().Add(t.RefreshTokenLifetime).Unix(),
	}
}

func (t *TokenService) signToken(claims jwt.MapClaims, privateKey rsa.PrivateKey) (*string, error) {
	jwk, err := t.KeyService.PublicKeyToJWK(privateKey.PublicKey)
	if err != nil {
		return nil, wrap("signToken", err)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessToken.Header["kid"] = jwk.Kid

	signedAccessToken, err := accessToken.SignedString(&privateKey)
	if err != nil {
		return nil, wrap("signToken", err)
	}

	return aws.String(signedAccessToken), nil
}

func (t *TokenService) IssueTokens(username string, privateKey *rsa.PrivateKey) (*Tokens, error) {
	accessToken, err := t.signToken(t.accessTokenClaims(username), *privateKey)
	if err != nil {
		return nil, wrap("IssueTokens", err)
	}

	refreshToken, err := t.signToken(t.refreshTokenClaims(username), *privateKey)
	if err != nil {
		return nil, wrap("IssueTokens", err)
	}

	return &Tokens{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (t *TokenService) VerifyRefreshToken(token string, jwkSets keys.JWKSets) (*string, error) {
	getKey := func(tn *jwt.Token) (interface{}, error) {
		if _, ok := tn.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, wrap("VerifyRefreshToken", errors.New("unexpected signing method"))
		}

		k := keys.KeyService{}
		publicKey, err := k.JWKSetsToPublicKey(jwkSets, tn.Header["kid"].(string))
		if err != nil {
			return nil, wrap("VerifyRefreshToken", err)
		}

		return publicKey, nil
	}

	parsed, err := jwt.Parse(token, getKey)
	if err != nil {
		return nil, wrap("VerifyRefreshToken", err)
	}

	_ = parsed

	if !parsed.Valid {
		return nil, wrap("VerifyRefreshToken", errors.New("invalid token claims"))
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, wrap("VerifyRefreshToken", errors.New("invalid token claims"))
	}

	use, ok := claims["token_use"].(string)
	if !ok || use != "refresh" {
		return nil, wrap("VerifyRefreshToken", errors.New("invalid token claims"))
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, wrap("VerifyRefreshToken", errors.New("invalid token claims"))
	}

	return aws.String(username), nil
}
