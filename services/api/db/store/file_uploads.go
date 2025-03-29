package store

import (
	"context"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

type FileUploadStore interface {
	CreateFileUpload(ctx context.Context, fileUpload *models.FileUpload) (*models.FileUpload, error)
	GetFileUploadByIds(ctx context.Context, ids []uuid.UUID) ([]models.FileUpload, error)
	GetAllFileUploads(ctx context.Context) ([]*models.FileUpload, error)
	UpdateFileUploadStatus(ctx context.Context, fileUploadId uuid.UUID, status models.FileUploadStatus) (*models.FileUpload, error)
}

func (s *appStore) CreateFileUpload(ctx context.Context, fileUpload *models.FileUpload) (*models.FileUpload, error) {
	err := s.client.WithContext(ctx).Create(fileUpload).Error
	if err != nil {
		return nil, err
	}
	return fileUpload, nil
}

func (s *appStore) GetFileUploadByIds(ctx context.Context, ids []uuid.UUID) ([]models.FileUpload, error) {
	var fileUploads []models.FileUpload
	err := s.client.WithContext(ctx).Where("id IN (?)", ids).Preload("UploadedByUser").Find(&fileUploads).Error
	if err != nil {
		return nil, err
	}
	return fileUploads, nil
}

func (s *appStore) GetAllFileUploads(ctx context.Context) ([]*models.FileUpload, error) {
	var fileUploads []*models.FileUpload
	err := s.client.WithContext(ctx).Order("created_at DESC").Find(&fileUploads).Error
	if err != nil {
		return nil, err
	}
	return fileUploads, nil
}

func (s *appStore) UpdateFileUploadStatus(ctx context.Context, fileUploadId uuid.UUID, status models.FileUploadStatus) (*models.FileUpload, error) {

	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user found in context")
	}

	var fileUpload models.FileUpload
	fileUpload.ID = fileUploadId
	fileUpload.UploadedByUserID = *currentUserId

	err := s.client.WithContext(ctx).Model(&fileUpload).Update("status", status).Error
	if err != nil {
		return nil, err
	}
	return &fileUpload, nil
}
