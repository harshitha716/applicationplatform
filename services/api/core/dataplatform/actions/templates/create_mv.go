package templates

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

const (
	createMVRunName           = "create_mv"
	createMVTaskKey           = "create_mv_task"
	createMVSideEffectTaskKey = "create_mv_sideeffect"
)

func GetMVJobTemplate(config *serverconfig.DataPlatformConfig, merchantId string, warehouseId string) *jobs.SubmitRun {
	return &jobs.SubmitRun{
		RunName: fmt.Sprintf("%s_%s", createMVRunName, merchantId),
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
				TaskKey: createMVTaskKey,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.CreateMVJobTemplateConfig.CreateMVNotebookPath,
					Source:       jobs.SourceWorkspace,
				},
				TimeoutSeconds:       0,
				EmailNotifications:   &jobs.JobEmailNotifications{},
				WebhookNotifications: &jobs.WebhookNotifications{},
			},
			{
				TaskKey: createMVSideEffectTaskKey,
				DependsOn: []jobs.TaskDependency{
					{
						TaskKey: createMVTaskKey,
					},
				},
				RunIf: jobs.RunIfAllSuccess,
				NotebookTask: &jobs.NotebookTask{
					NotebookPath: config.ActionsConfig.CreateMVJobTemplateConfig.SideEffectNotebookPath,
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
						TaskKey: createMVSideEffectTaskKey,
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

type MVCreateQueryResponse struct {
	MVName string
	Query  string
}
