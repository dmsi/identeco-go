package keys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco/pkg/jwks"
	"github.com/dmsi/identeco/pkg/s3helper"
)

// I/O with S3 bucket
type KeyService struct {
	S3                   *s3helper.S3Session
	Bucket               string
	JWKSObjectName       string
	PrivateKeyObjectName string
}

type keys struct {
	privatePem []byte
	jwk        *jwks.JWK
}

func (k *KeyService) generateKeys(bits int) (*keys, error) {
	// Generate keypair
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	// Encode private key as PEM
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
		},
	)

	// Convert public key to JWK
	jwk, err := jwks.PublicKeyToJWK(privatekey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &keys{privatePem: pemdata, jwk: jwk}, nil
}

func (k *KeyService) RotateKeys() error {
	keys, err := k.generateKeys(2048)
	if err != nil {
		return err
	}

	// TODO atomic - write all or nothing

	// Write Private Key PEM
	buf := bytes.Buffer{}
	_, err = buf.Write(keys.privatePem)
	if err != nil {
		return err
	}

	err = k.S3.WriteObject(k.Bucket, k.PrivateKeyObjectName, &buf)
	if err != nil {
		return err
	}

	// Write JWKS
	// Make JWKS set of a previous and current key or just the current key in case
	// of first rotation

	// Generate keys
	// fmt.Printf("Tests here!\n")
	// privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	// if err != nil {
	// 	fmt.Printf("Cannot generate RSA key\n")
	// 	os.Exit(1)
	// }
	// publickey := &privatekey.PublicKey

	// fmt.Printf(">>> private %v\n", privatekey)
	// fmt.Printf(">>> public %v\n", publickey)

	// // GetJWKS from public key
	// j, err := jwks.PublicKeyToJWK(*publickey)
	// fmt.Printf(">>> jwks %v, err %v\n", j, err)

	// // Encode private key as PEM
	// pemdata := pem.EncodeToMemory(
	// 	&pem.Block{
	// 		Type:  "RSA PRIVATE KEY",
	// 		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	// 	},
	// )

	// fmt.Printf(">>> pem %v\n", string(pemdata))

	// p, _ := pem.Decode(pemdata)
	// priv, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	// fmt.Printf(">>> frompem %v, err :%v\n", priv, err)
	return nil
}

func (k *KeyService) GetJWKS() (jwks.JWKS, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	err := k.S3.ReadObject(k.Bucket, k.JWKSObjectName, buf)
	if err != nil {
		return jwks.JWKS{}, err
	}

	j := jwks.JWKS{}
	err = json.Unmarshal(buf.Bytes(), &j)
	if err != nil {
		return jwks.JWKS{}, err
	}

	return j, nil
}

// Reads pem from S3 and return as rsa.PrivateKey
func (k *KeyService) GetPrivateKey() (string, error) {
	return "", nil
}
