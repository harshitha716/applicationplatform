package helpers

import (
	"strings"

	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/errors"
)

func ExtractBucketAndPath(fileURL string) (bucket, filePath string, err error) {
	if !strings.HasPrefix(fileURL, "gs://") {
		return "", "", errors.ErrDefaultGcsBucket
	}

	trimmedURL := strings.TrimPrefix(fileURL, "gs://")
	splitURL := strings.SplitN(trimmedURL, "/", 2)

	if len(splitURL) != 2 {
		return "", "", errors.ErrInvalidGcsFileUrl
	}

	bucket = splitURL[0]
	filePath = splitURL[1]
	return bucket, filePath, nil
}
