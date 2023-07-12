package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	e "github.com/dmsi/identeco-go/pkg/lib/err"
)

type KeyService struct {
	PrivateKeyBits int
}

func wrap(name string, err error) error {
	return e.Wrap("services.keys."+name, err)
}

func (k *KeyService) GenerateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, k.PrivateKeyBits)
	if err != nil {
		return nil, wrap("GenerateKey", err)
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
		return nil, wrap("PrivateKeyDecodePEM", err)
	}

	return privateKey, nil
}
