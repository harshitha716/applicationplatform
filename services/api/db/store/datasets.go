package store

import (
	"context"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatasetStore interface {
	GetDatasetById(ctx context.Context, datasetId string) (*models.Dataset, error)
	GetDatasetsAll(ctx context.Context, filters models.DatasetFilters) ([]models.Dataset, error)
	GetDatasetCount(ctx context.Context, filters models.DatasetFilters) (int64, error)
	CreateDataset(ctx context.Context, dataset models.Dataset) (uuid.UUID, error)
	UpdateDataset(ctx context.Context, dataset models.Dataset) (uuid.UUID, error)
	datasetPoliciesStore
	WithDatasetTransaction(ctx context.Context, fn func(DatasetStore) error) error
	DeleteDataset(ctx context.Context, dataset models.Dataset) error
}

func (s *appStore) GetDatasetById(ctx context.Context, datasetId string) (*models.Dataset, error) {
	var dataset models.Dataset
	err := s.client.WithContext(ctx).Model(dataset).Where("dataset_id = ?", datasetId).First(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

func (s *appStore) GetDatasetsAll(ctx context.Context, filters models.DatasetFilters) ([]models.Dataset, error) {
	db := s.client.WithContext(ctx)

	datasets := []models.Dataset{}
	if len(filters.OrganizationIds) > 0 {
		db = db.Where("organization_id IN (?)", filters.OrganizationIds)
	}

	if len(filters.DatasetIds) > 0 {
		db = db.Where("dataset_id IN (?)", filters.DatasetIds)
	}

	if len(filters.CreatedBy) > 0 {
		db = db.Where("created_by IN (?)", filters.CreatedBy)
	}

	if len(filters.Type) > 0 {
		db = db.Where("type IN (?)", filters.Type)
	}

	db = db.Where("deleted_at IS NULL")

	if filters.Page > 0 {
		db = db.Offset((filters.Page - 1) * filters.Limit)
	}

	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}

	if len(filters.SortParams) > 0 {
		for _, sort := range filters.SortParams {
			db = db.Order(clause.OrderByColumn{
				Column: clause.Column{Name: sort.Column},
				Desc:   sort.Desc,
			})
		}
	}

	result := db.Find(&datasets)
	if result.Error != nil {
		return nil, result.Error
	}

	return datasets, nil

}

func (s *appStore) GetDatasetCount(ctx context.Context, filters models.DatasetFilters) (int64, error) {
	db := s.client.WithContext(ctx)

	var count int64
	if len(filters.OrganizationIds) > 0 {
		db = db.Where("organization_id IN (?)", filters.OrganizationIds)
	}

	if len(filters.DatasetIds) > 0 {
		db = db.Where("dataset_id IN (?)", filters.DatasetIds)
	}

	if len(filters.CreatedBy) > 0 {
		db = db.Where("created_by IN (?)", filters.CreatedBy)
	}

	if len(filters.Type) > 0 {
		db = db.Where("type IN (?)", filters.Type)
	}

	db = db.Where("deleted_at IS NULL")

	result := db.Model(&models.Dataset{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (s *appStore) CreateDataset(ctx context.Context, dataset models.Dataset) (uuid.UUID, error) {
	err := s.client.WithContext(ctx).Create(&dataset).Error
	if err != nil {
		return uuid.Nil, err
	}
	return dataset.ID, nil
}

func (s *appStore) UpdateDataset(ctx context.Context, dataset models.Dataset) (uuid.UUID, error) {
	err := s.client.WithContext(ctx).
		Model(&dataset).
		Where("dataset_id = ?", dataset.ID).
		Updates(&dataset).Error
	if err != nil {
		return uuid.Nil, err
	}
	return dataset.ID, nil
}

func (s *appStore) DeleteDataset(ctx context.Context, dataset models.Dataset) error {
	now := time.Now()
	err := s.client.WithContext(ctx).
		Model(&models.Dataset{}).
		Where("dataset_id = ?", dataset.ID).
		Update("deleted_at", &now).Error
	return err
}

func (s *appStore) WithDatasetTransaction(ctx context.Context, fn func(DatasetStore) error) error {
	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txClient := pgclient.PostgresClient{DB: tx}
		return fn(&appStore{client: &txClient})
	})
}
