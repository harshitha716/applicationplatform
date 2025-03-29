package models

import (
	"encoding/json"
	"fmt"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatasetFileAllignmentStatus string

const (
	DatasetFileAllignmentStatusPending   DatasetFileAllignmentStatus = "allignment_pending"
	DatasetFileAllignmentStatusCompleted DatasetFileAllignmentStatus = "allignment_completed"
	DatasetFileAllignmentStatusFailed    DatasetFileAllignmentStatus = "allignment_failed"
)

type DatasetFileUpload struct {
	ID                   uuid.UUID                   `json:"id" gorm:"primaryKey; default:gen_random_uuid()"`
	DatasetID            uuid.UUID                   `json:"dataset_id" gorm:"not null"`
	FileUploadID         uuid.UUID                   `json:"file_upload_id" gorm:"not null"`
	FileAllignmentStatus DatasetFileAllignmentStatus `json:"status" gorm:"column:status;not null"`
	Metadata             json.RawMessage             `json:"metadata"`
}

// TODO Add error msg accordingly
type DatasetFileUploadMetadata struct {
	DataPreview           DatasetPreview         `json:"data_preview"`
	TransformedDataBucket string                 `json:"transformed_data_bucket"`
	TransformedDataPath   string                 `json:"transformed_data_path"`
	ExtractedMetadata     ExtractedMetadata      `json:"extracted_metadata"`
	ColumnMapping         map[string]interface{} `json:"column_mapping"`
	Error                 string                 `json:"error"`
}

type DatasetPreview struct {
	Columns []string                 `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
}

type ExtractedMetadata struct {
	Data map[string]interface{} `json:"data"`
}

func (dfu *DatasetFileUpload) TableName() string {
	return "dataset_file_uploads"
}

func (dfu *DatasetFileUpload) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
				SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
				WHERE frap.resource_type = 'dataset'
				AND frap.resource_id = dataset_file_uploads.dataset_id
				AND frap.user_id = ?
				AND frap.deleted_at IS NULL
			)`, userId,
	)
}

func (dfu *DatasetFileUpload) BeforeCreate(db *gorm.DB) (err error) {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err = db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", ResourceTypeDataset, dfu.DatasetID, userId).Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("dataset access forbidden")
	}

	return nil
}
