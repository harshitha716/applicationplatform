package models

type SignedUrlToUpload struct {
	Url        string
	Identifier string
}

type GetDownloadsignedUrlConfigs struct {
	Key   string
	Value interface{}
}
