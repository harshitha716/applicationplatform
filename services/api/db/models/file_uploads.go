package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileUploadStatus string

const (
	FileUploadStatusPending   FileUploadStatus = "pending"
	FileUploadStatusCompleted FileUploadStatus = "completed"
	FileUploadStatusFailed    FileUploadStatus = "failed"
)

type FileType string

const (
	FileTypeCSV     FileType = "csv"
	FileTypeXLSX    FileType = "xlsx"
	FileTypeXLS     FileType = "xls"
	FileTypeParquet FileType = "parquet"
)

type StorageType string

const (
	StorageTypeS3  StorageType = "s3"
	StorageTypeGCS StorageType = "gcs"
)

type FileUpload struct {
	ID               uuid.UUID        `json:"id" gorm:"primaryKey; default:gen_random_uuid()"`
	CreatedAt        time.Time        `json:"created_at" gorm:"not null; default:now()"`
	Name             string           `json:"name" gorm:"not null"`
	FileType         FileType         `json:"file_type" gorm:"not null"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"not null; default:now()"`
	OrganizationID   uuid.UUID        `json:"organization_id" gorm:"not null"`
	UploadedByUserID uuid.UUID        `json:"uploaded_by_user_id" gorm:"not null"`
	UploadedByUser   *User            `json:"user" gorm:"foreignKey:UploadedByUserID;references:ID"`
	PresignedURL     string           `json:"presigned_url"`
	Expiry           time.Time        `json:"expiry"`
	StorageProvider  StorageType      `json:"storage_provider"`
	StorageBucket    string           `json:"storage_bucket"`
	StorageFilePath  string           `json:"storage_file_path"`
	Status           FileUploadStatus `json:"status"`
	DeletedAt        *time.Time       `json:"deleted_at"`
}

func (f *FileUpload) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	// TODO: OR dataset access filters
	return db.Where("uploaded_by_user_id = ? AND deleted_at IS NULL", userId)
}

func (f *FileUpload) BeforeCreate(db *gorm.DB) (err error) {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	f.UploadedByUserID = *userId

	fraps := []FlattenedResourceAudiencePolicy{}
	err = db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", ResourceTypeOrganization, f.OrganizationID, userId).Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil
}

func (f *FileUpload) BeforeUpdate(db *gorm.DB) (err error) {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	if f.UploadedByUserID != *userId {
		return fmt.Errorf("user does not have permission to update this file upload")
	}

	return nil
}
