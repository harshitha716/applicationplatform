// Package s3 provides types and interfaces for S3 operations
package s3

import (
	"context"
	"io"
	"time"
)

// Client defines the interface for S3 operations
type Client interface {
	// GenerateUploadURL generates a pre-signed URL for uploading a file
	GenerateUploadURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)

	// ListFiles lists files in a bucket with an optional prefix
	ListFiles(ctx context.Context, bucket, prefix string) ([]FileInfo, error)

	// ReadFile reads a file from a bucket and returns its contents
	ReadFile(ctx context.Context, bucket, key string) (io.ReadCloser, error)

	// RenameFile renames a file within a bucket
	RenameFile(ctx context.Context, bucket, oldKey, newKey string) error

	// CopyFile copies a file from one bucket to another
	CopyFile(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) error

	// DeleteFile deletes a file from a bucket
	DeleteFile(ctx context.Context, bucket, key string) error

	// GetFileInfo gets metadata about a file without downloading it
	GetFileInfo(ctx context.Context, bucket, key string) (*FileInfo, error)
}

// FileInfo represents metadata about a file in S3
type FileInfo struct {
	// Key is the full path/name of the file in the bucket
	Key string

	// Size is the size of the file in bytes
	Size int64

	// LastModified is when the file was last modified
	LastModified time.Time

	// ContentType is the MIME type of the file
	ContentType string

	// ETag is the entity tag of the object
	ETag string
}
