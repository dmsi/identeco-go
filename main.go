package main

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dmsi/identeco/pkg/s3helper"
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

	s := s3helper.NewS3Session()
	s.Downloader = &TestDownloader{}

	buf1 := aws.NewWriteAtBuffer([]byte{})
	err := s.ReadObject("test-keys-070523", "jwk1.json", buf1)
	fmt.Printf("ReadObject ==> buf: %v, err: %v\n", buf1, err)
}
