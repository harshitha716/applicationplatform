package createdataset

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform"
	dpactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	"github.com/Zampfi/application-platform/services/api/core/datasets/service"
	datasetService "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	fileimportsservice "github.com/Zampfi/application-platform/services/api/core/fileimports"
	rulesservice "github.com/Zampfi/application-platform/services/api/core/rules/service"
	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	cloudservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/service"
	querybuilder "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/activity"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

type CreateDatasetWorkflow struct {
	datasetService datasetService.DatasetService
}

func InitCreateDatasetWorkflow(serverConfig *serverconfig.ServerConfig) CreateDatasetWorkflow {
	queryBuilderService := querybuilder.NewQueryBuilder()
	dataplatformService, err := dataplatform.InitDataPlatformService(serverConfig.DataPlatformConfig)
	if err != nil {
		panic(err)
	}
	ruleService := rulesservice.NewRuleService(serverConfig.Store)
	fileImportService := fileimportsservice.NewFileImportService(serverConfig.DefaultS3Client, serverConfig.Store, serverConfig.Env.AWSDefaultBucketName)
	cloudService, err := cloudservice.NewCloudService("GCP", *serverConfig.Env)
	if err != nil {
		panic(err)
	}

	return CreateDatasetWorkflow{
		datasetService: service.NewDatasetService(
			serverConfig.Store,
			queryBuilderService,
			dataplatformService,
			ruleService,
			fileImportService,
			serverConfig.TemporalSdk,
			cloudService,
			serverConfig.DefaultS3Client,
			*serverConfig.DatasetConfig,
			serverConfig.CacheClient,
		),
	}
}

func (w *CreateDatasetWorkflow) ApplicationPlatformCreateDatasetWorkflowExecute(wtx workflow.Context, params CreateDatasetWorkflowInitPayload) (interface{}, error) {

	var registerDatasetResponse *RegisterDatasetResponse
	err := workflow.ExecuteActivity(wtx, w.RegisterDatasetActivity, params).Get(wtx, &registerDatasetResponse)
	if err != nil {
		return nil, err
	}

	logger := workflow.GetLogger(wtx)
	logger.Info("Dataset created", zap.Any("dataset_id", registerDatasetResponse.DatasetId))
	logger.Info("Waiting for dataset to be ready")

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Minute,
	}

	wtx = workflow.WithActivityOptions(wtx, activityOpts)

	var datasetAction *datasetmodels.DatasetAction
	err = workflow.ExecuteActivity(wtx, w.WaitForDatasetToBeReady, params.UserId, params.OrganizationId, registerDatasetResponse.ActionId).Get(wtx, &datasetAction)
	if err != nil {
		return nil, err
	}

	logger.Info("Dataset ready", zap.Any("dataset_id", datasetAction.DatasetId))

	return CreateDatasetWorkflowExitPayload{
		DatasetId: registerDatasetResponse.DatasetId,
	}, nil
}

func (w *CreateDatasetWorkflow) RegisterDatasetActivity(ctx context.Context, params CreateDatasetWorkflowInitPayload) (*RegisterDatasetResponse, error) {

	userId, orgIds := params.GetAccessControlParams()
	if userId == nil {
		return nil, errors.New("user id is required")
	}

	if len(orgIds) == 0 {
		return nil, errors.New("organization id is required")
	}

	ctxWithAuth := apicontext.AddAuthToContext(ctx, "user", *userId, orgIds)

	actionId, datasetId, err := w.datasetService.RegisterDataset(ctxWithAuth, orgIds[0], *userId, params.RegisterDatasetPayload)
	if err != nil {
		return nil, err
	}
	return &RegisterDatasetResponse{
		ActionId:  actionId,
		DatasetId: datasetId,
	}, nil
}

func (w *CreateDatasetWorkflow) WaitForDatasetToBeReady(ctx context.Context, userId uuid.UUID, organizationId uuid.UUID, actionId string) (*datasetmodels.DatasetAction, error) {

	ctxWithAuth := apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{organizationId})

	action := &datasetmodels.DatasetAction{}

	for !slices.Contains(dpactionconstants.ActionTerminationStatuses, dpactionconstants.ActionStatus(action.Status)) {
		datasetActions, err := w.datasetService.GetDatasetActions(ctxWithAuth, organizationId, storemodels.DatasetActionFilters{
			ActionIds: []string{actionId},
		})
		if err != nil {
			return nil, fmt.Errorf("error getting dataset actions: %w", err)
		}

		if len(datasetActions) == 0 {
			return nil, errors.New("no dataset actions found")
		}

		action = &datasetActions[0]
		time.Sleep(3 * time.Second)
	}

	return action, nil
}

func (w *CreateDatasetWorkflow) GetActivities() []activity.Activity {
	return []activity.Activity{
		{
			Function: w.RegisterDatasetActivity,
			RegisterOptions: activity.RegisterActivityOptions{
				DisableAlreadyRegisteredCheck: true,
			},
		},
		{
			Function: w.WaitForDatasetToBeReady,
			RegisterOptions: activity.RegisterActivityOptions{
				DisableAlreadyRegisteredCheck: true,
			},
		},
	}
}
