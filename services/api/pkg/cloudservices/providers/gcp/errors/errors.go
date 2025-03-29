package errors

import "errors"

const (
	ErrCodeDefaultGcsBucket = "DEFAULT_GCS_BUCKET"
	ErrMsgDefaultGcsBucket  = "Default GCS bucket"

	ErrCodeInvalidGcsFileUrl = "INVALID_GCS_FILE_URL"
	ErrMsgInvalidGcsFileUrl  = "Invalid GCS file url"

	ErrCodeInvalidDownloadParams = "INVALID_DOWNLOAD_PARAMS"
	ErrMsgInvalidDownloadParams  = "Invalid download params"
)

var (
	ErrDefaultGcsBucket      = errors.New(ErrCodeDefaultGcsBucket)
	ErrInvalidGcsFileUrl     = errors.New(ErrCodeInvalidGcsFileUrl)
	ErrInvalidDownloadParams = errors.New(ErrCodeInvalidDownloadParams)
)
