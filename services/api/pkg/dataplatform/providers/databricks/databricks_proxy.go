package databricks

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/logger"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"go.uber.org/zap"
)

type DatabricksSDKProxy interface {
	SubmitOneTimeJob(ctx context.Context, jobConfig *jobs.SubmitRun) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse], error)
	CreateJob(ctx context.Context, jobConfig *jobs.CreateJob) (*jobs.CreateResponse, error)
	GetRunDetails(ctx context.Context, runId int64) (*jobs.Run, error)
	RunNow(ctx context.Context, params jobs.RunNow) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse], error)
}

func InitDatabricksSDKProxy(configs models.DatabricksConfig) (*databricks.WorkspaceClient, error) {
	ws, err := databricks.NewWorkspaceClient(&databricks.Config{
		Host:  configs.ServerHostname,
		Token: configs.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (db *databricksService) SubmitOneTimeJob(ctx context.Context, jobConfig *jobs.SubmitRun) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse], error) {
	logger := logger.GetLoggerFromCtx(ctx)
	run, err := db.ws.Jobs.Submit(ctx, *jobConfig)
	if err != nil {
		logger.Error("ERR_SUBMIT_ONE_TIME_JOB", zap.Error(err))
		return nil, err
	}
	return run, nil
}

func (db *databricksService) CreateJob(ctx context.Context, jobConfig *jobs.CreateJob) (*jobs.CreateResponse, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	createResponse, err := db.ws.Jobs.Create(ctx, *jobConfig)
	if err != nil {
		logger.Error("ERR_CREATE_JOB", zap.Error(err))
		return nil, err
	}
	return createResponse, nil
}

func (db *databricksService) GetRunDetails(ctx context.Context, runId int64) (*jobs.Run, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	runDetails, err := db.ws.Jobs.GetRun(ctx, jobs.GetRunRequest{
		RunId: runId,
	})
	if err != nil {
		logger.Error("ERR_GET_RUN_DETAILS", zap.Error(err))
		return nil, err
	}
	return runDetails, nil
}

func (db *databricksService) RunNow(ctx context.Context, params jobs.RunNow) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse], error) {
	logger := logger.GetLoggerFromCtx(ctx)
	run, err := db.ws.Jobs.RunNow(ctx, params)
	if err != nil {
		logger.Error("ERR_RUN_NOW", zap.Error(err))
		return nil, err
	}
	return run, nil
}
