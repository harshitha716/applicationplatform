package helpers

import (
	"testing"

	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/errors"
	"github.com/stretchr/testify/assert"
)

func TestExtractBucketAndPath(t *testing.T) {
	tests := []struct {
		name         string
		fileURL      string
		wantBucket   string
		wantFilePath string
		wantErr      error
	}{
		{
			name:         "valid gs url",
			fileURL:      "gs://bucket-name/path/to/file.txt",
			wantBucket:   "bucket-name",
			wantFilePath: "path/to/file.txt",
			wantErr:      nil,
		},
		{
			name:         "missing gs:// prefix",
			fileURL:      "bucket-name/file.txt",
			wantBucket:   "",
			wantFilePath: "",
			wantErr:      errors.ErrDefaultGcsBucket,
		},
		{
			name:         "invalid format",
			fileURL:      "gs://invalid-format",
			wantBucket:   "",
			wantFilePath: "",
			wantErr:      errors.ErrInvalidGcsFileUrl,
		},
		{
			name:         "empty url",
			fileURL:      "",
			wantBucket:   "",
			wantFilePath: "",
			wantErr:      errors.ErrDefaultGcsBucket,
		},
		{
			name:         "empty bucket",
			fileURL:      "gs:///file.txt",
			wantBucket:   "",
			wantFilePath: "file.txt",
			wantErr:      nil,
		},
		{
			name:         "empty file path",
			fileURL:      "gs://bucket-name/",
			wantBucket:   "bucket-name",
			wantFilePath: "",
			wantErr:      nil,
		},
		{
			name:         "invalid format no slash",
			fileURL:      "gs://invalid-format",
			wantBucket:   "",
			wantFilePath: "",
			wantErr:      errors.ErrInvalidGcsFileUrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket, path, err := ExtractBucketAndPath(tt.fileURL)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.wantBucket, bucket)
				assert.Equal(t, tt.wantFilePath, path)
			}
		})
	}
}
