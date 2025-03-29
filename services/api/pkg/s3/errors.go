// Package s3 provides error types for the S3 client
package s3

import (
	"fmt"
)

// BucketError represents an error related to bucket operations
type BucketError struct {
	Bucket string
	Err    error
}

func (e *BucketError) Error() string {
	if e.Bucket == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("bucket %s: %v", e.Bucket, e.Err)
}

func (e *BucketError) Unwrap() error {
	return e.Err
}

// ObjectError represents an error related to object operations
type ObjectError struct {
	Bucket string
	Key    string
	Err    error
}

func (e *ObjectError) Error() string {
	if e.Bucket == "" && e.Key == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("object %s/%s: %v", e.Bucket, e.Key, e.Err)
}

func (e *ObjectError) Unwrap() error {
	return e.Err
}

// Common S3 errors
var (
	// ErrInvalidBucket is returned when an invalid bucket name is provided
	ErrInvalidBucket = fmt.Errorf("invalid bucket name")

	// ErrObjectNotFound is returned when the requested object does not exist
	ErrObjectNotFound = fmt.Errorf("object not found")

	// ErrInvalidKey is returned when an invalid object key is provided
	ErrInvalidKey = fmt.Errorf("invalid object key")

	// ErrBucketNotFound is returned when the specified bucket does not exist
	ErrBucketNotFound = fmt.Errorf("bucket not found")

	// ErrAccessDenied is returned when access to a resource is denied
	ErrAccessDenied = fmt.Errorf("access denied to S3 resource")

	// ErrOperationFailed is returned when an S3 operation fails
	ErrOperationFailed = fmt.Errorf("S3 operation failed")
)
