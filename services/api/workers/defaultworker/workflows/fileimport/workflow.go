package fileimport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform"

	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	datasetService "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	fileimportsvc "github.com/Zampfi/application-platform/services/api/core/fileimports"
	rulesservice "github.com/Zampfi/application-platform/services/api/core/rules/service"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	cloudservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/service"
	"github.com/Zampfi/application-platform/services/api/pkg/errorreporting"
	querybuilder "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	"github.com/Zampfi/application-platform/services/api/workers/defaultworker/constants"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/activity"
	temporalmodels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

type fileImportWorkflow struct {
	datasetService     datasetService.DatasetService
	fileUploadsService fileimportsvc.FileImportService
	temporalSdk        temporal.TemporalService
	datasetStore       store.DatasetStore
}

func InitFileImportWorkflow(serverConfig *serverconfig.ServerConfig) (fileImportWorkflow, error) {

	queryBuilderService := querybuilder.NewQueryBuilder()
	dataplatformService, err := dataplatform.InitDataPlatformService(serverConfig.DataPlatformConfig)
	if err != nil {
		return fileImportWorkflow{}, err
	}
	ruleService := rulesservice.NewRuleService(serverConfig.Store)
	fileUploadsService := fileimportsvc.NewFileImportService(serverConfig.DefaultS3Client, serverConfig.Store, serverConfig.Env.AWSDefaultBucketName)
	cloudService, err := cloudservice.NewCloudService("GCP", *serverConfig.Env)
	if err != nil {
		return fileImportWorkflow{}, err
	}

	datasetService := datasetService.NewDatasetService(serverConfig.Store, queryBuilderService, dataplatformService, ruleService, fileUploadsService, serverConfig.TemporalSdk, cloudService, serverConfig.DefaultS3Client, *serverConfig.DatasetConfig, serverConfig.CacheClient)

	return initFileImport(datasetService, fileUploadsService, serverConfig.TemporalSdk, serverConfig.Store), nil
}

func initFileImport(datasetService datasetService.DatasetService, fileUploadsService fileimportsvc.FileImportService, temporalSdk temporal.TemporalService, datasetStore store.DatasetStore) fileImportWorkflow {
	return fileImportWorkflow{
		datasetService:     datasetService,
		fileUploadsService: fileUploadsService,
		temporalSdk:        temporalSdk,
		datasetStore:       datasetStore,
	}
}

func (f fileImportWorkflow) ApplicationPlatformDatasetFileImportWorkflowExecute(wtx workflow.Context, params models.FileImportWorkflowInitPayload) (interface{}, error) {
	logger := workflow.GetLogger(wtx)

	if params.UserId == uuid.Nil {
		logger.Error("user id is required in params")
		return nil, fmt.Errorf("user id is required in params")
	}

	if params.OrganizationId == uuid.Nil {
		logger.Error("organization id is required in params")
		return nil, fmt.Errorf("organization id is required in params")
	}

	helloWorldActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}

	wtx = workflow.WithActivityOptions(wtx, helloWorldActivityOptions)

	logger.Info("executing ai normalization activity")
	var fileImportExitPayload models.FileImportWorkflowExitPayload
	err := workflow.ExecuteActivity(wtx, f.ExecuteAINormalizationActivity, params).Get(wtx, &fileImportExitPayload)
	if err != nil {
		fileImportExitPayload.Error = fmt.Sprintf("failed to normalize file into a table: %s", err.Error())
		logger.Error("failed to normalize file into a table", zap.Error(err))
		return fileImportExitPayload, err
	}

	logger.Info("registering exit payload in DB")
	var transformedFilePath string
	err = workflow.ExecuteActivity(wtx, f.RegisterExitPayloadInDB, params, fileImportExitPayload).Get(wtx, &transformedFilePath)
	if err != nil {
		fileImportExitPayload.Error = fmt.Sprintf("failed to register exit payload in DB: %s", err.Error())
		logger.Error("failed to register exit payload in DB", zap.Error(err))
		return fileImportExitPayload, err
	}

	logger.Info("returning transformed file path")
	return transformedFilePath, nil
}

func (f fileImportWorkflow) ExecuteAINormalizationActivity(ctx context.Context, params models.FileImportWorkflowInitPayload) (*models.FileImportWorkflowExitPayload, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	ctx = apicontext.AddAuthToContext(ctx, "user", params.UserId, []uuid.UUID{params.OrganizationId})

	file, err := f.fileUploadsService.GetFileUploadByIds(ctx, []uuid.UUID{params.FileUploadId})
	if err != nil {
		logger.Error("failed to get file upload by id", zap.Error(err))
		return nil, err
	}

	if len(file) == 0 {
		logger.Error("no file upload found")
		return nil, fmt.Errorf("no file upload found")
	}

	fileUpload := file[0]

	logger.Info("sending path to AI", zap.String("storage_file_path", fileUpload.StorageFilePath), zap.String("storage_bucket", fileUpload.StorageBucket))

	fileImportConfig, err := f.datasetService.GetDatasetImportPath(ctx, params.OrganizationId, params.DatasetId)
	if err != nil {
		logger.Error("failed to get dataset import path", zap.String("error", err.Error()))
		return nil, err
	}

	tableDetectionWorkflowInput := map[string]interface{}{
		"source_file_path":     fileUpload.StorageFilePath,
		"source_bucket":        fileUpload.StorageBucket,
		"output_format_path":   fileImportConfig.BronzeSourcePath,
		"output_format_bucket": fileImportConfig.BronzeSourceBucket,
		"config":               fileImportConfig.BronzeSourceConfig,
	}

	var result map[string]interface{}
	workflowId := uuid.New()

	_, workflowErr := f.temporalSdk.ExecuteSyncWorkflow(ctx, temporalmodels.ExecuteWorkflowParams{
		Options: temporalmodels.StartWorkflowOptions{
			ID:        workflowId.String(),
			TaskQueue: constants.ImportFileNormalizationWorkflowTaskQueue,
		},
		Workflow: constants.ImportFileNormalizationWorkflowName,
		Args: []interface{}{
			tableDetectionWorkflowInput,
		},
		ResultPtr: &result,
	})

	logger.Info("Normalization workflow response:", zap.Any("workflow_result:", result))

	if workflowErr != nil {
		logger.Error("failed to execute AI normalization workflow", zap.Error(workflowErr))
		errorreporting.CaptureException(fmt.Errorf("failed to execute AI normalization workflow: %w", workflowErr), ctx)
		return nil, fmt.Errorf("failed to execute AI normalization workflow: %w", workflowErr)
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		logger.Error("failed to marshal workflow result", zap.Error(err))
		errorreporting.CaptureException(fmt.Errorf("failed to marshal workflow result into bytes: %w", err), ctx)
		return nil, fmt.Errorf("failed to marshal workflow result into bytes: %w", err)
	}

	var normalizationWorkflowResponse models.AINormalizationWorkflowResponse
	if err := json.Unmarshal(resultBytes, &normalizationWorkflowResponse); err != nil {
		logger.Error("failed to unmarshal workflow result", zap.Error(err))
		errorreporting.CaptureException(fmt.Errorf("failed to unmarshal AI workflow result into internal type: %w", err), ctx)
		return nil, fmt.Errorf("failed to unmarshal workflow result: %w", err)
	}

	exitPayload := models.FileImportWorkflowExitPayload{}
	exitPayload.NormalizationResult = normalizationWorkflowResponse
	exitPayload.NormalizedFilePath = normalizationWorkflowResponse.TransformedDataPath
	exitPayload.TransformedFilePath = fileImportConfig.BronzeSourcePath
	exitPayload.SourceFilePath = fileUpload.StorageFilePath
	exitPayload.Error = ""

	return &exitPayload, nil
}

func (f fileImportWorkflow) RegisterExitPayloadInDB(ctx context.Context, initParams models.FileImportWorkflowInitPayload, exitParams models.FileImportWorkflowExitPayload) (string, error) {
	ctx = apicontext.AddAuthToContext(ctx, "user", initParams.UserId, []uuid.UUID{initParams.OrganizationId})
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := f.datasetService.UpdateDatasetActionStatus(ctx, initParams.DatasetActionId.String(), constants.DatasetActionStatusSuccessful)
	if err != nil {
		logger.Error("could not update dataset action config", zap.Error(err))
		return "", fmt.Errorf("could not update dataset action config: %s", err.Error())
	}

	err = f.datasetService.UpdateDatasetFileUploadStatus(ctx, initParams.DatasetFileUploadId, models.UpdateDatasetFileUploadParams{
		FileAllignmentStatus: dbmodels.DatasetFileAllignmentStatusCompleted,
		Metadata: dbmodels.DatasetFileUploadMetadata{
			DataPreview:           exitParams.NormalizationResult.DataPreview,
			TransformedDataBucket: exitParams.NormalizationResult.TransformedDataBucket,
			TransformedDataPath:   exitParams.NormalizationResult.TransformedDataPath,
			ColumnMapping:         exitParams.NormalizationResult.ColumnMapping,
			ExtractedMetadata: dbmodels.ExtractedMetadata{
				Data: exitParams.NormalizationResult.ExtractedMetadata,
			},
		},
	})
	if err != nil {
		logger.Error("failed to update dataset file upload status", zap.Error(err))
		return "", fmt.Errorf("failed to update dataset file upload status: %w", err)
	}

	return exitParams.TransformedFilePath, nil
}

func (f fileImportWorkflow) GetActivities() []activity.Activity {
	return []activity.Activity{
		{
			Function: f.ExecuteAINormalizationActivity,
			RegisterOptions: activity.RegisterActivityOptions{
				DisableAlreadyRegisteredCheck: true,
			},
		},
		{
			Function: f.RegisterExitPayloadInDB,
			RegisterOptions: activity.RegisterActivityOptions{
				DisableAlreadyRegisteredCheck: true,
			},
		},
	}
}
