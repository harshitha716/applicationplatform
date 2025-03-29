package templates

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

const (
	updateDatasetRunName  = "update_dataset"
	updateDatasetTaskKey  = "update_dataset_task"
	pinotIngestionTaskKey = "pinot_ingestion_task"
)

func GetUpdateDatasetJobTemplate(config *serverconfig.DataPlatformConfig, merchantId string, datasetId string) *jobs.SubmitRun {
	return &jobs.SubmitRun{
		RunName: fmt.Sprintf("%s_%s", updateDatasetRunName, datasetId),
		WebhookNotifications: &jobs.WebhookNotifications{
			OnStart: []jobs.Webhook{
				{
					Id: config.ActionsConfig.WebhookConfig.WebhookId,
				},
			},
			OnSuccess: []jobs.Webhook{
				{
					Id: config.ActionsConfig.WebhookConfig.WebhookId,
				},
			},
			OnFailure: []jobs.Webhook{
				{
					Id: config.ActionsConfig.WebhookConfig.WebhookId,
				},
			},
		},
		Tasks: []jobs.SubmitTask{
			{
				TaskKey: updateDatasetTaskKey,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.UpdateDatasetJobTemplateConfig.UpdateDatasetNotebookPath,
					Source:       jobs.SourceWorkspace,
				},
				TimeoutSeconds:       0,
				EmailNotifications:   &jobs.JobEmailNotifications{},
				WebhookNotifications: &jobs.WebhookNotifications{},
			},
			{
				TaskKey: pinotIngestionTaskKey,
				DependsOn: []jobs.TaskDependency{
					{
						TaskKey: updateDatasetTaskKey,
					},
				},
				RunIf: jobs.RunIfAllSuccess,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.PinotIngestionNotebookPath,
					Source:       jobs.SourceWorkspace,
				},
				ExistingClusterId:    config.ActionsConfig.PinotClusterId,
				TimeoutSeconds:       0,
				EmailNotifications:   &jobs.JobEmailNotifications{},
				WebhookNotifications: &jobs.WebhookNotifications{},
			},
		},
		Queue: &jobs.QueueSettings{
			Enabled: true,
		},
		RunAs: &jobs.JobRunAs{
			UserName: config.ActionsConfig.RunAsUserName,
		},
	}
}
