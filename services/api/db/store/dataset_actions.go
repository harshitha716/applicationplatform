package store

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	dataplatfromactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type DatasetActionStore interface {
	CreateDatasetAction(ctx context.Context, organizationId uuid.UUID, params models.CreateDatasetActionParams) error
	GetDatasetActions(ctx context.Context, organizationId uuid.UUID, filters models.DatasetActionFilters) ([]models.DatasetAction, error)
	GetDatasetActionFromActionId(ctx context.Context, actionId string) (*models.DatasetAction, error)
	UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error
	UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error
}

func (s *appStore) CreateDatasetAction(ctx context.Context, organizationId uuid.UUID, params models.CreateDatasetActionParams) error {
	db := s.client.WithContext(ctx)

	var completedAt *time.Time
	if params.IsCompleted {
		now := time.Now()
		completedAt = &now
	}

	config, err := json.Marshal(params.Config)
	if err != nil {
		return err
	}

	return db.Create(&models.DatasetAction{
		ID:             uuid.New(),
		OrganizationId: organizationId,
		ActionId:       params.ActionId,
		ActionType:     params.ActionType,
		DatasetId:      params.DatasetId,
		Status:         params.Status,
		Config:         config,
		ActionBy:       params.ActionBy,
		StartedAt:      time.Now(),
		CompletedAt:    completedAt,
	}).Error
}

func (s *appStore) GetDatasetActions(ctx context.Context, organizationId uuid.UUID, filters models.DatasetActionFilters) ([]models.DatasetAction, error) {
	db := s.client.WithContext(ctx)

	actions := []models.DatasetAction{}
	db = db.Where("organization_id = ?", organizationId)

	if len(filters.DatasetIds) > 0 {
		db = db.Where("dataset_id IN (?)", filters.DatasetIds)
	}

	if len(filters.ActionIds) > 0 {
		db = db.Where("action_id IN (?)", filters.ActionIds)
	}

	if len(filters.ActionType) > 0 {
		db = db.Where("action_type IN (?)", filters.ActionType)
	}

	if len(filters.ActionBy) > 0 {
		db = db.Where("action_by IN (?)", filters.ActionBy)
	}

	if len(filters.Status) > 0 {
		db = db.Where("status IN (?)", filters.Status)
	}

	return actions, db.Find(&actions).Order("started_at DESC").Error
}

func (s *appStore) GetDatasetActionFromActionId(ctx context.Context, actionId string) (*models.DatasetAction, error) {
	db := s.client.WithContext(ctx)

	action := models.DatasetAction{}
	return &action, db.Where("action_id = ?", actionId).First(&action).Error
}

func (s *appStore) UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error {
	db := s.client.WithContext(ctx)

	action := models.DatasetAction{}

	var completedAt *time.Time
	if slices.Contains(dataplatfromactionconstants.ActionTerminationStatuses, dataplatfromactionconstants.ActionStatus(status)) {
		now := time.Now()
		completedAt = &now
	}

	err := db.Model(&action).Where("action_id = ?", actionId).Updates(map[string]interface{}{
		"status":       status,
		"completed_at": completedAt,
	}).Error

	if err != nil {
		return err
	}

	return nil

}

func (s *appStore) UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error {
	db := s.client.WithContext(ctx)

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return db.Model(&models.DatasetAction{}).Where("action_id = ?", actionId).Updates(map[string]interface{}{
		"config": configBytes,
	}).Error
}
