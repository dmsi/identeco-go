package keyss3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	e "github.com/dmsi/identeco/pkg/lib/err"
	"github.com/dmsi/identeco/pkg/storage"
	"golang.org/x/exp/slog"
)

// TODO: wrap errors

type s3Session struct {
	uploader   s3manageriface.UploaderAPI
	downloader s3manageriface.DownloaderAPI
}

// TODO: private and jwk sets object names can be part of this struct
type KeysStorage struct {
	lg             *slog.Logger
	client         s3Session
	bucket         string
	privateKeyName string
	jwkSetsName    string
}

func New(lg *slog.Logger, bucket, privateKeyName, jwkSetsName string) *KeysStorage {
	sess := session.New()
	client := s3Session{
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}

	return &KeysStorage{
		lg:             lg,
		client:         client,
		bucket:         bucket,
		privateKeyName: privateKeyName,
		jwkSetsName:    jwkSetsName,
	}
}

func op(name string) string {
	return "storage.s3.keyss3." + name
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

func (k *KeysStorage) ReadPrivateKey() (*storage.PrivateKeyData, error) {
	data, err := k.client.readS3SmallObject(k.bucket, k.privateKeyName)
	if err != nil {
		return nil, e.Wrap(op("ReadPrivateKey"), err)
	}

	return &storage.PrivateKeyData{
		Data: data,
	}, nil
}

func (k *KeysStorage) WritePrivateKey(key storage.PrivateKeyData) error {
	err := k.client.writeS3SmallObject(k.bucket, k.privateKeyName, key.Data)
	if err != nil {
		return e.Wrap(op("WritePrivateKey"), err)
	}

	return nil
}

func (k *KeysStorage) ReadJWKSets() (*storage.JWKSetsData, error) {
	k.lg.Debug("Reading JWKSets", slog.Any(k.bucket, k.jwkSetsName))

	data, err := k.client.readS3SmallObject(k.bucket, k.jwkSetsName)
	if err != nil {
		return nil, e.Wrap(op("ReadJWKSets"), err)
	}

	return &storage.JWKSetsData{
		Data: data,
	}, nil
}

func (k *KeysStorage) WriteJWKSets(jwkSets storage.JWKSetsData) error {
	err := k.client.writeS3SmallObject(k.bucket, k.jwkSetsName, jwkSets.Data)
	if err != nil {
		return e.Wrap(op("WriteJWKSets"), err)
	}

	return nil
}
