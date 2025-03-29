package s3

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockS3API struct {
	mock.Mock
}

func (m *mockS3API) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.HeadObjectOutput), args.Error(1)
}

func (m *mockS3API) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func (m *mockS3API) CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.CopyObjectOutput), args.Error(1)
}

type getFileDetailsTestCase struct {
	name           string
	bucket         string
	key            string
	mockOutput     *s3.HeadObjectOutput
	mockError      error
	expectedError  bool
	expectedResult *FileInfo
}

func TestGetFileDetails(t *testing.T) {
	mockAPI := new(mockS3API)
	client := &s3Client{
		client: mockAPI,
	}

	ctx := context.Background()
	now := time.Now()

	tests := []getFileDetailsTestCase{
		{
			name:   "successful file details retrieval",
			bucket: "test-bucket",
			key:    "test-key",
			mockOutput: &s3.HeadObjectOutput{
				ContentLength: aws.Int64(100),
				LastModified:  aws.Time(now),
				ContentType:   aws.String("text/plain"),
				ETag:          aws.String("etag123"),
			},
			mockError: nil,
			expectedResult: &FileInfo{
				Key:          "test-key",
				Size:         100,
				LastModified: now,
				ContentType:  "text/plain",
				ETag:         "etag123",
			},
		},
		{
			name:           "error getting file details",
			bucket:         "test-bucket",
			key:            "test-key",
			mockOutput:     nil,
			mockError:      errors.New("s3 error"),
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.On("HeadObject", ctx, &s3.HeadObjectInput{
				Bucket: aws.String(tc.bucket),
				Key:    aws.String(tc.key),
			}).Return(tc.mockOutput, tc.mockError).Once()

			fileInfo, err := client.GetFileDetails(ctx, tc.bucket, tc.key)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, fileInfo)
				assert.Contains(t, err.Error(), "failed to get file details from S3")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult.Key, fileInfo.Key)
				assert.Equal(t, tc.expectedResult.Size, fileInfo.Size)
				assert.Equal(t, tc.expectedResult.LastModified, fileInfo.LastModified)
				assert.Equal(t, tc.expectedResult.ContentType, fileInfo.ContentType)
				assert.Equal(t, tc.expectedResult.ETag, fileInfo.ETag)
			}
		})
	}
}

type getSampleFilePathTestCase struct {
	name          string
	bucket        string
	prefix        string
	mockOutput    *s3.ListObjectsV2Output
	mockError     error
	expectedPath  string
	expectedError bool
	errorContains string
}

func TestGetSampleFilePathFromFolder(t *testing.T) {
	mockAPI := new(mockS3API)
	client := &s3Client{
		client: mockAPI,
	}

	ctx := context.Background()

	tests := []getSampleFilePathTestCase{
		{
			name:   "successful sample file path retrieval",
			bucket: "test-bucket",
			prefix: "test-prefix/",
			mockOutput: &s3.ListObjectsV2Output{
				Contents: []types.Object{
					{Key: aws.String("test-prefix/file1"), Size: aws.Int64(100)},
					{Key: aws.String("test-prefix/file2"), Size: aws.Int64(50)},
				},
			},
			expectedPath: "test-prefix/file2",
		},
		{
			name:   "empty folder",
			bucket: "test-bucket",
			prefix: "test-prefix/",
			mockOutput: &s3.ListObjectsV2Output{
				Contents: []types.Object{},
			},
			expectedError: true,
			errorContains: "no files found in folder prefix",
		},
		{
			name:   "only zero size files",
			bucket: "test-bucket",
			prefix: "test-prefix/",
			mockOutput: &s3.ListObjectsV2Output{
				Contents: []types.Object{
					{Key: aws.String("test-prefix/file1"), Size: aws.Int64(0)},
				},
			},
			expectedError: true,
			errorContains: "no files with size > 0 found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.On("ListObjectsV2", ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(tc.bucket),
				Prefix: aws.String(tc.prefix),
			}).Return(tc.mockOutput, tc.mockError).Once()

			path, err := client.GetSampleFilePathFromFolder(ctx, tc.bucket, tc.prefix)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Empty(t, path)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPath, path)
			}
		})
	}
}

type copyFileTestCase struct {
	name          string
	srcBucket     string
	srcKey        string
	destBucket    string
	destKey       string
	mockOutput    *s3.CopyObjectOutput
	mockError     error
	expectedError bool
	errorContains string
}

func TestCopyFile(t *testing.T) {
	mockAPI := new(mockS3API)
	client := &s3Client{
		client: mockAPI,
	}

	ctx := context.Background()

	tests := []copyFileTestCase{
		{
			name:          "successful file copy",
			srcBucket:     "source-bucket",
			srcKey:        "source-key",
			destBucket:    "dest-bucket",
			destKey:       "dest-key",
			mockOutput:    &s3.CopyObjectOutput{},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "copy error",
			srcBucket:     "source-bucket",
			srcKey:        "source-key",
			destBucket:    "dest-bucket",
			destKey:       "dest-key",
			mockOutput:    nil,
			mockError:     errors.New("copy error"),
			expectedError: true,
			errorContains: "failed to copy file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.On("CopyObject", ctx, &s3.CopyObjectInput{
				Bucket:     aws.String(tc.destBucket),
				Key:        aws.String(tc.destKey),
				CopySource: aws.String(fmt.Sprintf("%s/%s", tc.srcBucket, tc.srcKey)),
			}).Return(tc.mockOutput, tc.mockError).Once()

			err := client.CopyFile(ctx, tc.srcBucket, tc.srcKey, tc.destBucket, tc.destKey)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type generateUploadURLTestCase struct {
	name          string
	bucket        string
	key           string
	expiry        time.Time
	contentType   string
	expectedError bool
	errorContains string
}

func TestGenerateUploadURL(t *testing.T) {
	client := &s3Client{}
	ctx := context.Background()

	tests := []generateUploadURLTestCase{
		{
			name:          "expiry in past",
			bucket:        "test-bucket",
			key:           "test-key",
			contentType:   "text/plain",
			expiry:        time.Now().Add(-1 * time.Hour),
			expectedError: true,
			errorContains: "expiry time must be in the future",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url, err := client.GenerateUploadURL(ctx, tc.bucket, tc.key, tc.expiry, tc.contentType)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Empty(t, url)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
			}
		})
	}
}
