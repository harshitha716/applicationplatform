package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	dataplatfromactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	datasetConstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

func (s *datasetService) DatasetExportTemporalActivity(ctx context.Context, params models.DatasetExportParams, datasetId uuid.UUID, userId uuid.UUID, orgIds []uuid.UUID, workflowId string) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	ctx = apicontext.AddAuthToContext(ctx, "user", userId, orgIds)

	exportDatasetQueryConfig := params.QueryConfig
	exportDatasetQueryConfig.Pagination = &models.Pagination{
		Page:     datasetConstants.DatasetExportPage,
		PageSize: datasetConstants.DatasetExportPageSizeLimit,
	}
	exportDatasetQueryConfig.GetDatafromLake = true

	data, err := s.GetDataByDatasetId(ctx, orgIds[0], datasetId.String(), exportDatasetQueryConfig)
	if err != nil {
		logger.Error("failed to get dataset data", zap.String("error", err.Error()))
		return "", fmt.Errorf("failed to get dataset data: %w", err)
	}

	csvData, err := s.createCSVFromQueryResult(data)
	if err != nil {
		logger.Error("failed to create CSV from data", zap.Error(err))
		return "", fmt.Errorf("failed to create CSV from data: %w", err)
	}

	uploadResponse, err := s.cloudService.UploadFileToCloud(ctx, params.ExportPath, csvData.Bytes())
	if err != nil {
		logger.Error("failed to upload CSV to cloud storage", zap.Error(err))
		return "", fmt.Errorf("failed to upload CSV to cloud storage: %w", err)
	}

	err = s.datasetActionService.UpdateDatasetActionStatus(ctx, workflowId, string(dataplatfromactionconstants.ActionStatusSuccessful))
	if err != nil {
		logger.Error("failed to update dataset action", zap.Error(err))
		return "", fmt.Errorf("failed to update dataset action: %w", err)
	}

	return uploadResponse.Url, nil
}
