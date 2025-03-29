package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDestinationFolder(t *testing.T) {
	t.Parallel()

	filePath := "test/test.txt"
	destinationFolder := GetDestinationFolder(filePath)
	assert.Equal(t, "test", destinationFolder)

	filePath = "zamp-file-imports/test/.ajkdfj/nested/test.new.txt"
	destinationFolder = GetDestinationFolder(filePath)
	assert.Equal(t, "zamp-file-imports/test/.ajkdfj/nested", destinationFolder)

	filePath = "/test.new.txt"
	destinationFolder = GetDestinationFolder(filePath)
	assert.Equal(t, "/", destinationFolder)

}

func TestGetRenamedFilePath(t *testing.T) {
	t.Parallel()

	filePath := "test/test.txt"

	renamedFilePath := GetRenamedFilePath(filePath, "new")
	assert.Equal(t, "test/new.txt", renamedFilePath)

	filePath = "zamp-file-imports/test/.ajkdfj/nested/test.new.parquet"
	renamedFilePath = GetRenamedFilePath(filePath, "new")
	assert.Equal(t, "zamp-file-imports/test/.ajkdfj/nested/new.parquet", renamedFilePath)

	filePath = "/test.new.parquet"
	renamedFilePath = GetRenamedFilePath(filePath, "new")
	assert.Equal(t, "/new.parquet", renamedFilePath)

	filePath = "test.new.parquet"

}

func TestExtractBucketNameAndFolderPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		importUrl      string
		expectedBucket string
		expectedPrefix string
	}{
		{
			name:           "valid s3a url",
			importUrl:      "s3a://zamp-prd-file-imports/test/test.txt",
			expectedBucket: "zamp-prd-file-imports",
			expectedPrefix: "test/test.txt/",
		},
		{
			name:           "empty url",
			importUrl:      "",
			expectedBucket: "",
			expectedPrefix: "",
		},
		{
			name:           "invalid url format",
			importUrl:      "not-a-valid-url",
			expectedBucket: "",
			expectedPrefix: "",
		},
		{
			name:           "url with multiple path segments",
			importUrl:      "s3a://my-bucket/path/to/file.txt",
			expectedBucket: "my-bucket",
			expectedPrefix: "path/to/file.txt/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucketName, folderPrefix := ExtractBucketNameAndFolderPrefix(tt.importUrl)
			assert.Equal(t, tt.expectedBucket, bucketName)
			assert.Equal(t, tt.expectedPrefix, folderPrefix)
		})
	}
}
