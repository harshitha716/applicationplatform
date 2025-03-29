package templates

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

const (
	registerDatasetRunName = "register_dataset"
	registerDatasetTaskKey = "register_dataset_task"
)

func GetRegisterDatasetJobTemplate(config *serverconfig.DataPlatformConfig, merchantId string, datasetId string) *jobs.SubmitRun {
	return &jobs.SubmitRun{
		RunName: fmt.Sprintf("%s_%s", registerDatasetRunName, datasetId),
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
				TaskKey: registerDatasetTaskKey,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.RegisterDatasetJobTemplateConfig.RegisterDatasetNotebookPath,
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
