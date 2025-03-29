package templates

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

const (
	copyDatasetRunName = "copy_dataset"
	copyDatasetTaskKey = "copy_dataset_task"
)

func GetCopyDatasetJobTemplate(config *serverconfig.DataPlatformConfig, newDatasetId string) *jobs.SubmitRun {
	return &jobs.SubmitRun{
		RunName: fmt.Sprintf("%s_%s", copyDatasetRunName, newDatasetId),
		WebhookNotifications: &jobs.WebhookNotifications{
			OnStart: []jobs.Webhook{
				{Id: config.ActionsConfig.WebhookConfig.WebhookId},
			},
			OnSuccess: []jobs.Webhook{
				{Id: config.ActionsConfig.WebhookConfig.WebhookId},
			},
			OnFailure: []jobs.Webhook{
				{Id: config.ActionsConfig.WebhookConfig.WebhookId},
			},
		},
		Tasks: []jobs.SubmitTask{
			{
				TaskKey: copyDatasetTaskKey,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.CopyDatasetJobTemplateConfig.CopyDatasetNotebookPath,
					Source:       jobs.SourceWorkspace,
				},
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
