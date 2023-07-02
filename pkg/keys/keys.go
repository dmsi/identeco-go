package keys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"

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

func NewKeysService() *KeyService {
	return &KeyService{
		S3:                   s3helper.NewS3Session(),
		Bucket:               os.Getenv("BUCKET_NAME"),
		JWKSObjectName:       os.Getenv("JWKS_JSON_NAME"),
		PrivateKeyObjectName: os.Getenv("PRIVATE_KEY_NAME"),
	}
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

	newJWKS := jwks.JWKS{
		Keys: []jwks.JWK{
			*keys.jwk,
		},
	}

	// Write JWKS
	currentJWKS, err := k.GetJWKS()
	if err == nil {
		newJWKS.Keys = append(newJWKS.Keys, currentJWKS.Keys[0])
	}

	buf = bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err = enc.Encode(&newJWKS)
	if err != nil {
		return err
	}

	err = k.S3.WriteObject(k.Bucket, k.JWKSObjectName, &buf)
	if err != nil {
		return err
	}

	return nil
}

func (k *KeyService) readS3Bytes(bucketName, objectName string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	err := k.S3.ReadObject(bucketName, objectName, buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (k *KeyService) GetJWKSBytes() ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	err := k.S3.ReadObject(k.Bucket, k.JWKSObjectName, buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (k *KeyService) GetJWKS() (jwks.JWKS, error) {
	jbytes, err := k.GetJWKSBytes()
	if err != nil {
		return jwks.JWKS{}, err
	}

	j := jwks.JWKS{}
	err = json.Unmarshal(jbytes, &j)
	if err != nil {
		return jwks.JWKS{}, err
	}

	return j, nil
}

// Reads pem from S3 and return as raw pem data
func (k *KeyService) GetPrivateKeyRaw() ([]byte, error) {
	return k.readS3Bytes(k.Bucket, k.PrivateKeyObjectName)
}

// Reads pem from S3 and return as raw rsa.PrivateKey
func (k *KeyService) GetPrivateKey() (*rsa.PrivateKey, error) {
	keydata, err := k.GetPrivateKeyRaw()
	if err != nil {
		return nil, err
	}

	p, _ := pem.Decode(keydata)
	privateKey, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
