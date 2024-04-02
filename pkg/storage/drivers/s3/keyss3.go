package s3driver

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/storage"
)

type s3Session struct {
	uploader   s3manageriface.UploaderAPI
	downloader s3manageriface.DownloaderAPI
}

type UsersS3Driver struct {
	lg             *slog.Logger
	client         s3Session
	bucket         string
	privateKeyName string
	jwksName       string
}

func NewKeysStorage(lg *slog.Logger) (*UsersS3Driver, error) {
	cfg, ok := config.Cfg.KeysStorageDriver.(config.S3DriverConfig)
	if !ok {
		return nil, fmt.Errorf("keys driver configuration is not provided: s3")
	}

	sess := session.New()
	client := s3Session{
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}

	return &UsersS3Driver{
		lg:             lg,
		client:         client,
		bucket:         cfg.BucketName,
		privateKeyName: cfg.PrivateKeyName,
		jwksName:       cfg.JWKSName,
	}, nil
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

func (s UsersS3Driver) ReadKeys() (*storage.Keys, error) {
	keypair, err := s.client.readS3SmallObject(s.bucket, s.privateKeyName)
	if err != nil {
		return nil, err
	}

	jwks, err := s.client.readS3SmallObject(s.bucket, s.jwksName)
	if err != nil {
		return nil, err
	}

	return &storage.Keys{
		PrivateKey: keypair,
		JWKS:       jwks,
	}, nil
}

func (s UsersS3Driver) WriteKeys(k storage.Keys) error {
	err := s.client.writeS3SmallObject(s.bucket, s.privateKeyName, k.PrivateKey)
	if err != nil {
		return err
	}

	err = s.client.writeS3SmallObject(s.bucket, s.jwksName, k.JWKS)
	if err != nil {
		return err
	}

	return nil
}
