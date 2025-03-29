package helper

import (
	"path/filepath"
	"strings"
)

func GetDestinationFolder(filePath string) string {
	return filepath.Dir(filePath)
}

func GetRenamedFilePath(fileWithPath string, newName string) string {
	fileName := filepath.Base(fileWithPath)
	fileExt := filepath.Ext(fileName)

	newFileName := newName + fileExt

	return filepath.Join(GetDestinationFolder(fileWithPath), newFileName)
}

// importUrl -> s3a://zamp-prd-file-imports/../..
func ExtractBucketNameAndFolderPrefix(importUrl string) (string, string) {
	parts := strings.Split(importUrl, "/")
	if len(parts) < 3 {
		return "", ""
	}

	return parts[2], strings.Join(parts[3:], "/") + "/"
}
