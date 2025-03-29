package service

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/core/datasets/actions/errors"
	"github.com/Zampfi/application-platform/services/api/core/datasets/actions/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DatasetActionService interface {
	CreateDatasetAction(ctx context.Context, organizationId uuid.UUID, params dbmodels.CreateDatasetActionParams) error
	GetDatasetActions(ctx context.Context, organizationId uuid.UUID, filters dbmodels.DatasetActionFilters) ([]models.DatasetAction, error)
	GetDatasetActionFromActionId(ctx context.Context, actionId string) (*models.DatasetAction, error)
	UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error
	UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error
}

type datasetActionService struct {
	store store.DatasetActionStore
}

func NewDatasetActionService(store store.DatasetActionStore) DatasetActionService {
	return &datasetActionService{
		store: store,
	}
}

func (s *datasetActionService) CreateDatasetAction(ctx context.Context, organizationId uuid.UUID, params dbmodels.CreateDatasetActionParams) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	if !s.validateDatasetActionStatus(params.Status) {
		logger.Error("Invalid dataset action status", zap.Any("datasetId", params.DatasetId), zap.Any("actionId", params.ActionId), zap.Any("status", params.Status))
		return errors.ErrInvalidDatasetActionStatus
	}

	if !s.validateDatasetActionType(params.ActionType) {
		logger.Error("Invalid dataset action type", zap.Any("datasetId", params.DatasetId), zap.Any("actionId", params.ActionId), zap.Any("actionType", params.ActionType))
		return errors.ErrInvalidDatasetActionType
	}

	err := s.store.CreateDatasetAction(ctx, organizationId, params)
	if err != nil {
		logger.Error("Error creating dataset action", zap.Error(err))
		return err
	}

	logger.Info("Dataset action created", zap.Any("datasetId", params.DatasetId), zap.Any("actionId", params.ActionId))
	return nil
}

func (s *datasetActionService) GetDatasetActions(ctx context.Context, organizationId uuid.UUID, filters dbmodels.DatasetActionFilters) ([]models.DatasetAction, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actions, err := s.store.GetDatasetActions(ctx, organizationId, filters)
	if err != nil {
		logger.Error("Error getting dataset actions", zap.Error(err))
		return nil, err
	}

	datasetActions := []models.DatasetAction{}
	for _, action := range actions {
		datasetAction := models.DatasetAction{}
		datasetAction.FromSchema(action)
		datasetActions = append(datasetActions, datasetAction)
	}
	return datasetActions, nil
}

func (s *datasetActionService) GetDatasetActionFromActionId(ctx context.Context, actionId string) (*models.DatasetAction, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	action, err := s.store.GetDatasetActionFromActionId(ctx, actionId)
	if err != nil {
		logger.Error("Error getting dataset action", zap.Error(err))
		return nil, err
	}

	datasetAction := models.DatasetAction{}
	datasetAction.FromSchema(*action)
	return &datasetAction, nil
}

func (s *datasetActionService) UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	if !s.validateDatasetActionStatus(status) {
		logger.Error("Invalid dataset action status", zap.Any("actionId", actionId), zap.Any("status", status))
		return errors.ErrInvalidDatasetActionStatus
	}

	err := s.store.UpdateDatasetActionStatus(ctx, actionId, status)
	if err != nil {
		logger.Error("Error updating dataset action status", zap.Error(err))
		return err
	}

	logger.Info("Dataset action status updated", zap.Any("actionId", actionId), zap.Any("status", status))
	return nil
}

func (s *datasetActionService) UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := s.store.UpdateDatasetActionConfig(ctx, actionId, config)
	if err != nil {
		logger.Error("Error updating dataset action config", zap.Error(err))
		return err
	}

	logger.Info("Dataset action config updated", zap.Any("actionId", actionId), zap.Any("config", config))
	return nil
}
