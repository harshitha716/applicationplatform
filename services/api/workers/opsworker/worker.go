package main

import (
	"context"
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/pkg/logging"
	"github.com/Zampfi/application-platform/services/api/workers/opsworker/constants"
	"github.com/Zampfi/application-platform/services/api/workers/opsworker/workflows/createdataset"
	"github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/activity"
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

	// Create root context for the worker
	ctx := context.Background()

	// Create a new worker instance with the server configuration
	opsWorker := NewOpsWorker(serverConfig)

	// Start the worker and handle any startup errors
	worker, err := opsWorker.Run(ctx)
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

// opsWorker represents a Temporal worker implementation that processes workflows
// and activities from a ops task queue.
type opsWorker struct {
	// serverConfig contains all the necessary dependencies and configuration
	serverConfig *serverconfig.ServerConfig
	// workflows is a map of workflow names to their implementations
	workflows map[string]workflow.Workflow
}

// NewOpsWorker creates a new instance of opsWorker with the provided server configuration.
func NewOpsWorker(serverConfig *serverconfig.ServerConfig) *opsWorker {
	return &opsWorker{
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
func (w *opsWorker) Run(ctx context.Context) (worker.Worker, error) {
	// Initialize workflows that this worker will process

	// Dataset Export workflow initialization
	datasetCreationWorkflow := createdataset.InitCreateDatasetWorkflow(w.serverConfig)
	datasetCreationActivities := datasetCreationWorkflow.GetActivities()

	// Combine activities
	allActivities := []activity.Activity{}
	allActivities = append(
		allActivities,
		datasetCreationActivities...,
	)

	// Initialize the worker with the following configuration:
	// - TaskQueue: The queue this worker will poll for tasks
	// - Workflows: List of workflows this worker can execute
	// - Activities: List of activities this worker can execute
	// - Options: Additional worker options
	worker, err := w.serverConfig.TemporalSdk.GetNewWorker(ctx, models.NewWorkerParams{
		TaskQueue: constants.OpsTaskQueueName,
		Workflows: []workflow.Workflow{
			{
				Function: datasetCreationWorkflow.ApplicationPlatformCreateDatasetWorkflowExecute,
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
