// Package s3 provides a high-level SDK for AWS S3 operations
package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3API defines the interface for AWS S3 operations
type S3API interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
}

// S3Client defines the interface for high-level S3 operations
type S3Client interface {
	GetFileDetails(ctx context.Context, bucket string, key string) (*FileInfo, error)
	GetSampleFilePathFromFolder(ctx context.Context, bucket string, folderPrefix string) (string, error)
	GenerateUploadURL(ctx context.Context, bucket string, key string, expiry time.Time, contentType string) (string, error)
	CopyFile(ctx context.Context, srcBucket string, srcKey string, destBucket string, destKey string) error
}

type s3Client struct {
	client        S3API
	config        *Config
	presignClient *s3.PresignClient
}

func newS3Client(ctx context.Context, cfg *Config) (S3Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config: %w", err)
	}

	client := s3.NewFromConfig(sdkConfig)

	presignClient := s3.NewPresignClient(client)

	return &s3Client{
		client:        client,
		config:        cfg,
		presignClient: presignClient,
	}, nil
}

func NewDefaultS3Client(ctx context.Context) (S3Client, error) {
	cfg := &Config{
		Region: "us-east-1",
	}
	return newS3Client(ctx, cfg)
}

func (c *s3Client) GetFileDetails(ctx context.Context, bucket string, key string) (*FileInfo, error) {

	result, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file details from S3: %w", err)
	}
	return &FileInfo{
		Key:          key,
		Size:         *result.ContentLength,
		LastModified: *result.LastModified,
		ContentType:  *result.ContentType,
		ETag:         *result.ETag,
	}, nil
}

func (c *s3Client) GetSampleFilePathFromFolder(ctx context.Context, bucket string, folderPrefix string) (string, error) {
	result, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folderPrefix),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get file path from folder: %w", err)
	}

	if len(result.Contents) == 0 {
		return "", fmt.Errorf("no files found in folder prefix: %s", folderPrefix)
	}

	smallestFile := result.Contents[0]
	for _, file := range result.Contents {
		if file.Size != nil && *file.Size > 0 {
			smallestFile = file
			break
		}
	}

	if *smallestFile.Size == 0 {
		return "", fmt.Errorf("no files with size > 0 found in folder prefix: %s", folderPrefix)
	}

	// Find the smallest file
	for _, file := range result.Contents {
		if file.Size != nil && *file.Size > 0 && *file.Size < *smallestFile.Size {
			smallestFile = file
		}
	}

	return *smallestFile.Key, nil
}

func (c *s3Client) GenerateUploadURL(ctx context.Context, bucket string, key string, expiry time.Time, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Calculate duration until expiry, but ensure it's not negative
	duration := time.Until(expiry)
	if duration < 0 {
		return "", fmt.Errorf("expiry time must be in the future")
	}

	presignedRequest, err := c.presignClient.PresignPutObject(ctx, input,
		s3.WithPresignExpires(duration),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload URL: %w", err)
	}
	return presignedRequest.URL, nil
}

func (c *s3Client) CopyFile(ctx context.Context, srcBucket string, srcKey string, destBucket string, destKey string) error {
	_, err := c.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(destBucket),
		Key:        aws.String(destKey),
		CopySource: aws.String(fmt.Sprintf("%s/%s", srcBucket, srcKey)),
	})
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}
