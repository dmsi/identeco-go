package s3helper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDownloader struct {
	data map[string][]byte
}

type testUploader struct {
	data map[string][]byte
}

func (d *testDownloader) Download(object io.WriterAt, input *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	path := *input.Bucket + "/" + *input.Key

	data, ok := d.data[path]
	if !ok {
		return 0, fmt.Errorf("Path %v not found", path)
	}

	_, err := object.WriteAt(data, 0)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (d *testDownloader) DownloadWithContext(ctx aws.Context, object io.WriterAt, input *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	return 0, nil
}

func (u *testUploader) Upload(input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	bytes, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}

	path := *input.Bucket + "/" + *input.Key
	u.data[path] = bytes

	return nil, nil
}

func (u *testUploader) UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return nil, nil
}

func newMockSession() *S3Session {
	return &S3Session{
		Uploader: &testUploader{
			data: make(map[string][]byte),
		},
		Downloader: &testDownloader{
			data: make(map[string][]byte),
		},
	}
}

func TestReadObject(t *testing.T) {
	bucket := "test-bucket"
	key := "test-key"
	path := bucket + "/" + key

	svc := newMockSession()
	downloader := svc.Downloader.(*testDownloader)
	require.NotNil(t, downloader)

	_, ok := downloader.data[path]
	require.False(t, ok)
	downloader.data[path] = []byte("Sun/Milky Way")

	buf := aws.NewWriteAtBuffer([]byte{})
	err := svc.ReadObject(bucket, key, buf)
	require.ErrorIs(t, err, nil)
	assert.Equal(t, string(downloader.data[path]), string(buf.Bytes()))
}

func TestWriteObject(t *testing.T) {
	bucket := "test-bucket"
	key := "test-key"
	path := bucket + "/" + key

	svc := newMockSession()
	uploader := svc.Uploader.(*testUploader)
	require.NotNil(t, uploader)

	_, ok := uploader.data[path]
	require.False(t, ok)

	expected := "Arcturus/Milky Way"
	buf := bytes.Buffer{}
	_, err := buf.WriteString(expected)
	require.ErrorIs(t, err, nil)

	err = svc.WriteObject(bucket, key, &buf)
	require.ErrorIs(t, err, nil)

	assert.Equal(t, expected, string(uploader.data[path]))
}
