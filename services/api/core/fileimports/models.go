package fileimports

import (
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type FileImportSignedURL struct {
	UploadURL    string          `json:"upload_url"`
	Key          string          `json:"key"`
	FileUploadID uuid.UUID       `json:"file_upload_id"`
	FileType     models.FileType `json:"file_type"`
	FileName     string          `json:"file_name"`
}

type AckFileImportCompletionPayload struct {
	FileUploadID uuid.UUID `json:"file_upload_id"`
}
