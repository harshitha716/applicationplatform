package store

import (
	"context"
	"encoding/json"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type DatasetFileUploadStore interface {
	CreateDatasetFileUpload(ctx context.Context, datasetFileUpload *models.DatasetFileUpload) (*models.DatasetFileUpload, error)
	GetDatasetFileUploadByDatasetId(ctx context.Context, datasetId uuid.UUID) ([]models.DatasetFileUpload, error)
	GetDatasetFileUploadById(ctx context.Context, fileUploadId uuid.UUID) (models.DatasetFileUpload, error)
	UpdateDatasetFileUploadStatus(ctx context.Context, id uuid.UUID, fileAllignmentStatus models.DatasetFileAllignmentStatus, metadata models.DatasetFileUploadMetadata) (*models.DatasetFileUpload, error)
}

func (s *appStore) CreateDatasetFileUpload(ctx context.Context, datasetFileUpload *models.DatasetFileUpload) (*models.DatasetFileUpload, error) {
	err := s.client.WithContext(ctx).Create(datasetFileUpload).Error
	if err != nil {
		return nil, err
	}
	return datasetFileUpload, nil
}

func (s *appStore) GetDatasetFileUploadByDatasetId(ctx context.Context, datasetId uuid.UUID) ([]models.DatasetFileUpload, error) {
	var datasetFileUploads []models.DatasetFileUpload
	err := s.client.WithContext(ctx).Where("dataset_id = ?", datasetId).Find(&datasetFileUploads).Error
	if err != nil {
		return nil, err
	}
	return datasetFileUploads, nil
}

func (s *appStore) GetDatasetFileUploadById(ctx context.Context, fileUploadId uuid.UUID) (models.DatasetFileUpload, error) {
	var datasetFileUpload models.DatasetFileUpload
	err := s.client.WithContext(ctx).
		Where("file_upload_id = ?", fileUploadId).
		First(&datasetFileUpload).Error
	if err != nil {
		return models.DatasetFileUpload{}, err
	}
	return datasetFileUpload, nil
}

func (s *appStore) UpdateDatasetFileUploadStatus(ctx context.Context, id uuid.UUID, fileAllignmentStatus models.DatasetFileAllignmentStatus, metadata models.DatasetFileUploadMetadata) (*models.DatasetFileUpload, error) {
	var datasetFileUpload models.DatasetFileUpload
	err := s.client.WithContext(ctx).Where("id = ?", id).First(&datasetFileUpload).Error
	if err != nil {
		return nil, err
	}

	datasetFileUpload.FileAllignmentStatus = fileAllignmentStatus
	metadataByte, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	datasetFileUpload.Metadata = metadataByte

	err = s.client.WithContext(ctx).Model(&datasetFileUpload).Updates(map[string]interface{}{
		"status":   fileAllignmentStatus,
		"metadata": metadataByte,
	}).Error
	if err != nil {
		return nil, err
	}
	return &datasetFileUpload, nil
}
