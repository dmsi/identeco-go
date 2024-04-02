package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateKey(privateKeyBits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeyBits)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func EncodePEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return pemdata, nil
}

func DecodePEM(pemdata []byte) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode(pemdata)
	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
