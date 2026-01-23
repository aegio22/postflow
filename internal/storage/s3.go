package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Bucket  string
	Client  *s3.Client
	Presign *s3.PresignClient
}

func NewS3(cfg aws.Config, bucket string) *S3Client {
	cli := s3.NewFromConfig(cfg)
	return &S3Client{
		Bucket:  bucket,
		Client:  cli,
		Presign: s3.NewPresignClient(cli),
	}
}

func (s *S3Client) PresignUpload(ctx context.Context, key string, expires time.Duration) (string, error) {
	out, err := s.Presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func (s *S3Client) PresignDownload(ctx context.Context, key string, expires time.Duration, filename string) (string, error) {
	out, err := s.Presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		ResponseContentDisposition: aws.String(
			fmt.Sprintf(`attachment; filename="%s"`, filename),
		),
	}, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func (s *S3Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.Client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.Bucket), Key: (aws.String(key))})
	if err != nil {
		return nil, fmt.Errorf("error fetching s3 object: %v", err)
	}
	return obj.Body, nil
}

func (s *S3Client) DeleteObject(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}
