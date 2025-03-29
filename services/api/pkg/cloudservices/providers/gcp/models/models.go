package models

import "time"

type DownloadParams struct {
	CustomDownloadName *string
	UrlTtlInMinutes    *time.Duration
}
