// internal/aws/s3.go
package aws

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client wraps the AWS S3 client
type S3Client struct {
	Client *s3.Client
}

// NewS3Client initializes and returns a new S3Client
func NewS3Client(region string) *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566", // Localstack URL
					SigningRegion: region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	// Configure the S3 client to use path-style addressing
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{Client: client}
}

// GetObject retrieves an object from S3
func (s *S3Client) GetObject(bucket, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	resp, err := s.Client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// PutObject uploads an object to S3
func (s *S3Client) PutObject(bucket, key string, body io.Reader) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	_, err := s.Client.PutObject(context.TODO(), input)
	return err
}
