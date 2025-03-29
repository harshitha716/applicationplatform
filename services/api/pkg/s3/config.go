// Package s3 provides configuration types for the S3 client
package s3

import "fmt"

// Config holds the configuration for the S3 client
type Config struct {
	// Region specifies the AWS region for S3 operations
	Region string

	// Endpoint is an optional field for S3-compatible services
	// If not specified, the default AWS S3 endpoint will be used
	Endpoint string

	// DefaultBucket is the default bucket to use when not explicitly specified
	DefaultBucket string

	// AllowedBuckets is a list of buckets that the client is allowed to access
	// If empty, all buckets are allowed
	AllowedBuckets []string

	// MaxRetries specifies the maximum number of retries for failed operations
	MaxRetries int

	// Timeout specifies the default timeout for operations
	// If zero, a default timeout will be used
	Timeout int64
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Region == "" {
		return fmt.Errorf("invalid configuration: %w: region is required", ErrInvalidBucket)
	}

	if c.DefaultBucket == "" {
		return fmt.Errorf("invalid configuration: %w: default bucket is required", ErrInvalidBucket)
	}

	return nil
}
