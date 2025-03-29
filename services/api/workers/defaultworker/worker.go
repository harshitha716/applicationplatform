package main

import (
	"context"
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform"
	datasetService "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	fileimportsservice "github.com/Zampfi/application-platform/services/api/core/fileimports"
	ruleservice "github.com/Zampfi/application-platform/services/api/core/rules/service"
	cloudservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/service"
	"github.com/Zampfi/application-platform/services/api/pkg/logging"
	querybuilderservice "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	"github.com/Zampfi/application-platform/services/api/workers/defaultworker/constants"
	"github.com/Zampfi/application-platform/services/api/workers/defaultworker/workflows/datasetexport"
	"github.com/Zampfi/application-platform/services/api/workers/defaultworker/workflows/fileimport"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/worker"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/workflow"
	temporalworker "go.temporal.io/sdk/worker"
)

// main initializes and runs a Temporal worker that processes workflows and activities
// from a specified task queue. The worker runs until the context is cancelled.
func main() {
	// Initialize server configuration which contains all necessary dependencies
	logger, err := logging.GetLogger()
	if err != nil {
		panic(err)
	}
	serverConfig, cleanup, err := serverconfig.Createserverconfig(logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	queryBuilderService := querybuilderservice.NewQueryBuilder()

	dataPlatformService, err := dataplatform.InitDataPlatformService(serverConfig.DataPlatformConfig)
	if err != nil {
		panic(err)
	}
	if dataPlatformService == nil {
		panic("DataPlatformService is nil")
	}

	ruleService := ruleservice.NewRuleService(serverConfig.Store)
	if ruleService == nil {
		panic("RuleService is nil")
	}

	fileImportService := fileimportsservice.NewFileImportService(serverConfig.DefaultS3Client, serverConfig.Store, serverConfig.Env.AWSDefaultBucketName)
	if fileImportService == nil {
		panic("FileImportService is nil")
	}

	cloudService, err := cloudservice.NewCloudService("GCP", *serverConfig.Env)
	if err != nil {
		panic(err)
	}
	if cloudService == nil {
		panic("CloudService is nil")
	}
	datasetService := datasetService.NewDatasetService(
		serverConfig.Store,
		queryBuilderService,
		dataPlatformService,
		ruleService,
		fileImportService,
		serverConfig.TemporalSdk,
		cloudService,
		serverConfig.DefaultS3Client,
		*serverConfig.DatasetConfig,
		serverConfig.CacheClient,
	)

	// Create root context for the worker
	ctx := context.Background()

	// Create a new worker instance with the server configuration
	defaultWorker := NewDefaultWorker(serverConfig)

	// Start the worker and handle any startup errors
	worker, err := defaultWorker.Run(ctx, datasetService)
	if err != nil {
		panic(err)
	}

	// Keep the worker running until context cancellation
	select {
	case <-ctx.Done():
		fmt.Println("Shutting down worker...")
		worker.Stop(ctx)
		return
	}
}

// defaultWorker represents a Temporal worker implementation that processes workflows
// and activities from a default task queue.
type defaultWorker struct {
	// serverConfig contains all the necessary dependencies and configuration
	serverConfig *serverconfig.ServerConfig
	// workflows is a map of workflow names to their implementations
	workflows map[string]workflow.Workflow
}

// NewDefaultWorker creates a new instance of defaultWorker with the provided server configuration.
func NewDefaultWorker(serverConfig *serverconfig.ServerConfig) *defaultWorker {
	return &defaultWorker{
		serverConfig: serverConfig,
	}
}

// Run initializes and starts the Temporal worker. It registers all workflows and activities
// that this worker should process and begins polling the task queue.
//
// To add a new workflow to this worker:
// 1. Initialize the workflow with its dependencies
// 2. Get the workflow's activities
// 3. Add the workflow and its activities to the worker's registration
func (w *defaultWorker) Run(ctx context.Context, datasetService datasetService.DatasetService) (worker.Worker, error) {
	// Initialize workflows that this worker will process

	// Dataset Export workflow initialization
	datasetExportWorkflow := datasetexport.InitDatasetExport(w.serverConfig, datasetService)
	datasetExportActivities := datasetExportWorkflow.GetActivities()

	fileImportWorkflow, err := fileimport.InitFileImportWorkflow(w.serverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file import workflow: %w", err)
	}
	fileImportActivities := fileImportWorkflow.GetActivities()

	// Combine activities
	allActivities := append(
		datasetExportActivities,
		fileImportActivities...,
	)

	// Initialize the worker with the following configuration:
	// - TaskQueue: The queue this worker will poll for tasks
	// - Workflows: List of workflows this worker can execute
	// - Activities: List of activities this worker can execute
	// - Options: Additional worker options
	worker, err := w.serverConfig.TemporalSdk.GetNewWorker(ctx, models.NewWorkerParams{
		TaskQueue: constants.DefaultTaskQueueName,
		Workflows: []workflow.Workflow{
			{
				Function: datasetExportWorkflow.ApplicationPlatformDatasetExportWorkflowExecute,
			},
			{
				Function: fileImportWorkflow.ApplicationPlatformDatasetFileImportWorkflowExecute,
			},
		},
		Activities: allActivities,
		Options: models.WorkerOptions{
			DisableRegistrationAliasing: true,
		},
		RegisterTasks: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize worker: %w", err)
	}

	// Start the worker with the provided context and interrupt channel
	err = worker.Run(ctx, temporalworker.InterruptCh(), models.RunWorkerParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to start worker: %w", err)
	}

	return worker, nil
}
