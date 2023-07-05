package keys

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
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

func decodeBigInt(s string) (*big.Int, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(bytes), nil
}

func encodeExponent(e int) (string, error) {
	be := new(big.Int).SetInt64(int64(e))
	return encodeBigInt(*be)
}

func encodeModulus(n big.Int) (string, error) {
	return encodeBigInt(n)
}

func decodeExponent(s string) (int, error) {
	e, err := decodeBigInt(s)
	if err != nil {
		return 0, err
	}
	return int(e.Int64()), nil
}

func decodeModulus(s string) (*big.Int, error) {
	return decodeBigInt(s)
}

func (k *KeyService) PublicKeyToJWK(pub rsa.PublicKey) (*JWK, error) {
	hash := md5.Sum(pub.N.Bytes())
	kid := hex.EncodeToString(hash[:])

	e, err := encodeExponent(pub.E)
	if err != nil {
		return nil, err
	}

	n, err := encodeModulus(*pub.N)
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

func (k *KeyService) JWKToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	e, err := decodeExponent(jwk.E)
	if err != nil {
		return nil, err
	}

	n, err := decodeModulus(jwk.N)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		E: e,
		N: n,
	}, nil
}

func (k *KeyService) JWKSetsToPublicKey(jwkSets JWKSets, kid string) (*rsa.PublicKey, error) {
	idx := -1
	for i, v := range jwkSets.Keys {
		if v.Kid == kid {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, fmt.Errorf("key with kid %s not found", kid)
	}

	publicKey, err := k.JWKToPublicKey(jwkSets.Keys[idx])
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
