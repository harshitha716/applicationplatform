package templates

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

const (
	registerJobRunName = "register_job"
	registerJobTaskKey = "register_job_task"
)

func GetRegisterJobJobTemplate(config *serverconfig.DataPlatformConfig, merchantId string, datasetId string) *jobs.SubmitRun {
	return &jobs.SubmitRun{
		RunName: fmt.Sprintf("%s_%s", registerJobRunName, datasetId),
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
				TaskKey: registerJobTaskKey,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.RegisterJobJobTemplateConfig.RegisterJobNotebookPath,
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
