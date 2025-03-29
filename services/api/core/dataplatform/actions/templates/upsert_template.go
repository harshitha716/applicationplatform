package templates

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

const (
	upsertTemplateRunName = "upsert_template"
	upsertTemplateTaskKey = "upsert_template_task"
)

func GetUpsertTemplateJobTemplate(config *serverconfig.DataPlatformConfig, merchantId string, templateId string) *jobs.SubmitRun {
	return &jobs.SubmitRun{
		RunName: fmt.Sprintf("%s_%s", upsertTemplateRunName, templateId),
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
				TaskKey: upsertTemplateTaskKey,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.UpsertTemplateJobTemplateConfig.UpsertTemplateNotebookPath,
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
