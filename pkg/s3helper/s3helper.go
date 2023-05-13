package s3helper

import (
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

// type S3Session interface {
// 	ReadObject(bucket, key string, object io.WriterAt) error
// 	WriteObject(buckt, key string, object io.Reader) error
// }

type API interface {
	Download(w io.WriterAt, input *s3.CopyObjectInput)
	Upload(input *s3manager.UploadInput)
}

// One Way
// type S3API struct {
// 	svc *session.Session
// }

// func (s *S3API) Download()

type S3Session struct {
	Uploader   s3manageriface.UploaderAPI
	Downloader s3manageriface.DownloaderAPI
}

func NewS3Session() *S3Session {
	sess := session.New()
	return &S3Session{
		Uploader:   s3manager.NewUploader(sess),
		Downloader: s3manager.NewDownloader(sess),
	}
}

func (s *S3Session) ReadObject(bucket, key string, object io.WriterAt) error {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	_, err := s.Downloader.Download(object, input)

	return err
}

func (s *S3Session) WriteObject(bucket, key string, object io.Reader) error {
	_, err := s.Uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   object,
	})

	return err
}
