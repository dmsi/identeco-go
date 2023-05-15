package s3helper

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadObject(t *testing.T) {
	bucket := "test-bucket"
	key := "test-key"
	path := bucket + "/" + key

	svc := NewMockSession()
	downloader := svc.Downloader.(*S3Mock)
	require.NotNil(t, downloader)

	_, ok := downloader.Data[path]
	require.False(t, ok)
	downloader.Data[path] = []byte("Sun/Milky Way")

	buf := aws.NewWriteAtBuffer([]byte{})
	err := svc.ReadObject(bucket, key, buf)
	require.ErrorIs(t, err, nil)
	assert.Equal(t, string(downloader.Data[path]), string(buf.Bytes()))
}

func TestWriteObject(t *testing.T) {
	bucket := "test-bucket"
	key := "test-key"
	path := bucket + "/" + key

	svc := NewMockSession()
	uploader := svc.Uploader.(*S3Mock)
	require.NotNil(t, uploader)

	_, ok := uploader.Data[path]
	require.False(t, ok)

	expected := "Arcturus/Milky Way"
	buf := bytes.Buffer{}
	_, err := buf.WriteString(expected)
	require.ErrorIs(t, err, nil)

	err = svc.WriteObject(bucket, key, &buf)
	require.ErrorIs(t, err, nil)

	assert.Equal(t, expected, string(uploader.Data[path]))
}
