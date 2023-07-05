package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type TestDownloader struct{}

func (t *TestDownloader) Download(object io.WriterAt, input *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	// I can form the object in here
	fmt.Printf("mock => object: %v\n", object)
	return 0, nil
}

func (t *TestDownloader) DownloadWithContext(ctx aws.Context, object io.WriterAt, input *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	return 0, nil
}

func main() {
	// svc := s3.New(session.New())
	// input := &s3.ListBucketsInput{}

	// result, err := svc.ListBuckets(input)
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		default:
	// 			fmt.Println(aerr.Error())
	// 		}
	// 	} else {
	// 		// Print the error, cast err to awserr.Error to get the Code and
	// 		// Message from an error.
	// 		fmt.Println(err.Error())
	// 	}
	// 	return
	// }
	// fmt.Printf("result: %v\n", result)

	// var buf bytes.Buffer
	// jwk := jwks.JWK{
	// 	E: "101010",
	// 	N: "888888",
	// }
	// err = json.NewEncoder(&buf).Encode(jwk)
	// fmt.Printf("err: %v\n", err)

	// s3helper.WriteObject("test-keys-070523", "jwk1.json", &buf)

	// s := s3helper.NewS3Session()
	// s.Downloader = &TestDownloader{}

	// buf1 := aws.NewWriteAtBuffer([]byte{})
	// err := s.ReadObject("test-keys-070523", "jwk1.json", buf1)
	// fmt.Printf("ReadObject ==> buf: %v, err: %v\n", buf1, err)

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

	// // b := x509.MarshalPKCS1PrivateKey(privatekey)
	// // fmt.Printf(">>> pem %v\n", string(b))

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

	// _ = j
	// _ = privatekey
	// _ = publickey

	// k := keys.KeyService{
	// 	S3:                   s3helper.NewS3Session(),
	// 	Bucket:               "identeco-dev-keys",
	// 	JWKSObjectName:       "jwks.json",
	// 	PrivateKeyObjectName: "privatekey.pem",
	// }

	// err := k.RotateKeys()
	// fmt.Printf(">>> err %v\n", err)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
