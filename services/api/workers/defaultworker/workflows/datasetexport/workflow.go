package datasetexport

import (
	"fmt"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	models "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	datasetService "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	"github.com/Zampfi/application-platform/services/api/workers/common/workflowutil"
	"github.com/Zampfi/application-platform/services/api/workers/defaultworker/constants"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/activity"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

type datasetExportWorkflow struct {
	datasetService datasetService.DatasetService
}

func InitDatasetExport(serverConfig *serverconfig.ServerConfig, datasetService datasetService.DatasetService) datasetExportWorkflow {

	return datasetExportWorkflow{
		datasetService: datasetService,
	}
}

func (d datasetExportWorkflow) ApplicationPlatformDatasetExportWorkflowExecute(wtx workflow.Context, params models.DatasetExportParams, datasetId uuid.UUID, userId uuid.UUID, orgIds []uuid.UUID, workflowId string) (interface{}, error) {

	defer workflowutil.PanicRecoveryHook(wtx)

	if userId == uuid.Nil {
		return nil, fmt.Errorf("user id is required")
	}
	if len(orgIds) == 0 {
		return nil, fmt.Errorf("org ids are required")
	}

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
	}
	paramsBase := baseParams{
		userId: userId,
		orgIds: orgIds,
	}

	wtx = workflowutil.AddAccessControlParamsToWorkflowCtx(wtx, &paramsBase)
	wtx = workflow.WithActivityOptions(wtx, activityOpts)

	var exportPath string
	err := workflow.ExecuteActivity(wtx, d.datasetService.DatasetExportTemporalActivity, params, datasetId, userId, orgIds, workflowId).Get(wtx, &exportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to export dataset: %w", err)
	}

	return exportPath, nil
}

func (d datasetExportWorkflow) GetActivities() []activity.Activity {
	return []activity.Activity{
		{
			Function: d.datasetService.DatasetExportTemporalActivity,
			RegisterOptions: activity.RegisterActivityOptions{
				Name: constants.DatasetExportTemporalActivity,
			},
		},
	}
}
