package keydatas3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/dmsi/identeco/pkg/lib/e"
	"github.com/dmsi/identeco/pkg/storage"
)

// TODO: wrap errors

type s3Session struct {
	uploader   s3manageriface.UploaderAPI
	downloader s3manageriface.DownloaderAPI
}

// TODO: private and jwk sets object names can be part of this struct
type KeyDataStorageS3 struct {
	client         s3Session
	bucket         string
	privateKeyName string
	jwkSetsName    string
}

func New(bucket, privateKeyName, jwkSetsName string) *KeyDataStorageS3 {
	sess := session.New()
	client := s3Session{
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}

	return &KeyDataStorageS3{
		client:         client,
		bucket:         bucket,
		privateKeyName: privateKeyName,
		jwkSetsName:    jwkSetsName,
	}
}

func op(name string) string {
	return "storage.s3.keydata." + name
}

func (s *s3Session) readS3SmallObject(bucket, key string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	_, err := s.downloader.Download(buf, input)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *s3Session) writeS3SmallObject(bucket, key string, object []byte) error {
	buf := bytes.Buffer{}
	_, err := buf.Write(object)
	if err != nil {
		return err
	}

	input := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   &buf,
	}

	_, err = s.uploader.Upload(input)
	if err != nil {
		return err
	}

	return nil
}

func (k *KeyDataStorageS3) ReadPrivateKey() (*storage.PrivateKeyData, error) {
	data, err := k.client.readS3SmallObject(k.bucket, k.privateKeyName)
	if err != nil {
		return nil, e.Wrap(op("ReadPrivateKey"), err)
	}

	return &storage.PrivateKeyData{
		Data: data,
	}, nil
}

func (k *KeyDataStorageS3) WritePrivateKey(key storage.PrivateKeyData) error {
	err := k.client.writeS3SmallObject(k.bucket, k.privateKeyName, key.Data)
	if err != nil {
		return e.Wrap(op("WritePrivateKey"), err)
	}

	return nil
}

func (k *KeyDataStorageS3) ReadJWKSets() (*storage.JWKSetsData, error) {
	data, err := k.client.readS3SmallObject(k.bucket, k.jwkSetsName)
	if err != nil {
		return nil, e.Wrap(op("ReadJWKSets"), err)
	}

	return &storage.JWKSetsData{
		Data: data,
	}, nil
}

func (k *KeyDataStorageS3) WriteJWKSets(jwkSets storage.JWKSetsData) error {
	err := k.client.writeS3SmallObject(k.bucket, k.jwkSetsName, jwkSets.Data)
	if err != nil {
		return e.Wrap(op("WriteJWKSets"), err)
	}

	return nil
}
