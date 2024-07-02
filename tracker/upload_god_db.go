package tracker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3Client *s3.Client // The AWS S3 client for interacting with the service
	bucket   string     // The name of the bucket in Cloudflare R2 Storage
}

func (t *Tracker) uploadGodDB(path string) error {
	s3Service, err := NewR2Service()
	if err != nil {
		return fmt.Errorf("failed to create R2 service: %w", err)
	}

	// Open the file to upload
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Upload the file to Cloudflare R2 Storage
	objectPrefix := os.Getenv("GOD_DB_PREFIX")
	uploadPath := filepath.Join(objectPrefix, "goddb.sqlite")
	err = s3Service.UploadFileToR2(context.Background(), uploadPath, fileBytes)
	if err != nil {
		return fmt.Errorf("failed to upload file to R2: %w", err)
	}

	return nil
}

func NewR2Service() (*S3Service, error) {
	// Replace these values with your Cloudflare R2 Storage credentials
	account := os.Getenv("CF_ACCOUNT_ID")
	accessKey := os.Getenv("GOD_DB_ACCESS_KEY")
	secretKey := os.Getenv("GOD_DB_SECRET_KEY")
	bucket := os.Getenv("GOD_DB_BUCKET")

	// Create custom resolver for R2 endpoint
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", account),
		}, nil
	})

	// Load AWS config with custom resolver
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("apac"),
	)
	if err != nil {
		return nil, err
	}

	// Create a new S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Return a new S3Service instance initialized with the S3 client and bucket name
	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (s *S3Service) UploadFileToR2(ctx context.Context, key string, file []byte) error {
	// Create a PutObjectInput with the specified bucket, key, file content, and content type
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("application/x-sqlite3"),
	}

	// Upload the file to Cloudflare R2 Storage
	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
