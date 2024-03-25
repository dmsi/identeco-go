package keys

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"

	"github.com/dmsi/identeco-go/pkg/lib/util"
)

type JWK struct {
	E   string `json:"e"`
	N   string `json:"n"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
}

type JWKSets struct {
	Keys []JWK `json:"keys"`
}

func encodeBigInt(n big.Int) (string, error) {
	return base64.RawURLEncoding.EncodeToString(n.Bytes()), nil
}

func encodeExponent(e int) (string, error) {
	be := new(big.Int).SetInt64(int64(e))
	return encodeBigInt(*be)
}

func encodeModulus(n big.Int) (string, error) {
	return encodeBigInt(n)
}

func PublicKeyToJWK(publicKey *rsa.PublicKey) (*JWK, error) {
	kid := util.KeyID(publicKey)

	e, err := encodeExponent(publicKey.E)
	if err != nil {
		return nil, err
	}

	n, err := encodeModulus(*publicKey.N)
	if err != nil {
		return nil, err
	}

	return &JWK{
		E:   e,
		N:   n,
		Kid: kid,
		Kty: "RSA",
		Alg: "RS256",
		Use: "sig",
	}, nil
}
