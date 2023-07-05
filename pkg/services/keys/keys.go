package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	e "github.com/dmsi/identeco/pkg/lib/err"
)

type KeyService struct {
	PrivateKeyBits int
}

func op(name string) string {
	return "services.keys." + name
}

func (k *KeyService) GenerateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, k.PrivateKeyBits)
	if err != nil {
		return nil, e.Wrap(op("GenerateKey"), err)
	}

	return privateKey, nil
}

func (k *KeyService) PrivateKeyEncodePEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return pemdata, nil
}

func (k *KeyService) PrivateKeyDecodePEM(pemdata []byte) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode(pemdata)
	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, e.Wrap(op("PrivateKeyDecodePEM"), err)
	}

	return privateKey, nil
}
