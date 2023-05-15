package s3helper

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type API interface {
	Download(w io.WriterAt, input *s3.CopyObjectInput)
	Upload(input *s3manager.UploadInput)
}

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

// TODO Data -> Storage
type S3Mock struct {
	Data map[string][]byte
}

func (d *S3Mock) Download(object io.WriterAt, input *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	path := *input.Bucket + "/" + *input.Key

	data, ok := d.Data[path]
	if !ok {
		return 0, fmt.Errorf("Path %v not found", path)
	}

	_, err := object.WriteAt(data, 0)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (d *S3Mock) DownloadWithContext(ctx aws.Context, object io.WriterAt, input *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	return 0, nil
}

func (u *S3Mock) Upload(input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	bytes, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}

	path := *input.Bucket + "/" + *input.Key
	u.Data[path] = bytes

	return nil, nil
}

func (u *S3Mock) UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return nil, nil
}

func NewMockSession() *S3Session {
	return &S3Session{
		Uploader: &S3Mock{
			Data: make(map[string][]byte),
		},
		Downloader: &S3Mock{
			Data: make(map[string][]byte),
		},
	}
}
