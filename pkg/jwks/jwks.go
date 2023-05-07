package jwks

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

// TODO to separate module jwks

// jwks.json example
// {
//     keys: [
//         {
//             e: 'AQAB',
//             kid: 'XooolbD0BPGABjHzSDRfQ4YBg8H87zwTJVmmP8I81OA',
//             kty: 'RSA',
//             n: 'tRXzVqY51HMCh-iK2K0YmGF044P2qM_42MDBZuk6CpqUg1Vm7ylBHLm41QWNIwvzyVtBiibjSPtT_Ua2-_6v5dz2bwZqUzxYU_yq5sacv3yfOpwe8mYej2wyaC0fBcKSigrpFj3nDHTXEUGIiR0Vptd7ja7vjOcj_8raGjaR7zGF_5P42OA-UUDmRmyU1PG_d4fV-bagip1byEcPM4GSxqOnWkJdNX9da82S9QxYSofFq9t8MYH2texM5ImcqZ0FmdUXb8k1DeBXv0dqg1ZbhaDvCzNWfgoMjhPeB5lpnCP0gR-X_3dLJDPI1lU0ddnjepCWuh48WuImxfilaoQCcw',
//             alg: 'RS256',
//             use: 'sig',
//         },
//     ],
// }

type JWK struct {
	E   string `json:"e"`
	N   string `json:"n"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
}

type JWKS struct {
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

func PublicKeyToJWK(pub rsa.PublicKey) (*JWK, error) {
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

func JWKToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
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
