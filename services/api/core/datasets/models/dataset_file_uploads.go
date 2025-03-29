package models

import (
	"time"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type DatasetFileUpload struct {
	ID                   uuid.UUID                            `json:"id"`
	DatasetID            uuid.UUID                            `json:"dataset_id"`
	FileID               uuid.UUID                            `json:"file_id"`
	FileName             string                               `json:"file_name"`
	UploadedByUserID     uuid.UUID                            `json:"uploaded_by_user_id"`
	UploadedByUser       *dbmodels.User                       `json:"uploaded_by_user"`
	FileUploadStatus     dbmodels.FileUploadStatus            `json:"file_upload_status"`
	FileUploadCreatedAt  time.Time                            `json:"file_upload_created_at"`
	FileUploadDeletedAt  *time.Time                           `json:"file_upload_deleted_at"`
	FileAllignmentStatus dbmodels.DatasetFileAllignmentStatus `json:"status"`
	Metadata             dbmodels.DatasetFileUploadMetadata   `json:"metadata"`
}

type UpdateDatasetFileUploadParams struct {
	FileAllignmentStatus dbmodels.DatasetFileAllignmentStatus `json:"status"`
	Metadata             dbmodels.DatasetFileUploadMetadata   `json:"metadata"`
}
