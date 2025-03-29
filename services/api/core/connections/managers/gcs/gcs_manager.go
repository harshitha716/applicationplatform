package gcs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/core/connections/constants"
	"github.com/Zampfi/application-platform/services/api/db/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	temporalModels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"github.com/google/uuid"
	temporalio "go.temporal.io/sdk/client"
)

type gcs_connection_manager struct {
}

func InitGCSConnectionManager() *gcs_connection_manager {
	return &gcs_connection_manager{}
}

func (g *gcs_connection_manager) CheckConnection() error {
	return nil
}

func (g *gcs_connection_manager) ExtractAndValidateDefaultSchedules(connectionId uuid.UUID, organizationId uuid.UUID, config models.CreateConnectionParams) ([]models.CreateScheduleParams, error) {
	schedules, ok := config.Config["schedules"].([]interface{})

	if !ok {
		return nil, errors.New("schedules not found in config")
	}

	if len(schedules) == 0 {
		return nil, errors.New("no schedules found in config")
	}

	schedulesModels := []models.CreateScheduleParams{}
	for i, schedule := range schedules {
		if scheduleMap, ok := schedule.(map[string]interface{}); ok {
			id := uuid.New().String()
			temporalWorkflowID := fmt.Sprintf("gcs_workflow:%s", id)
			scheduleMap[constants.BucketName] = config.Config[constants.BucketName].(string)
			scheduleMap[constants.TemporalWorkflowID] = temporalWorkflowID
			scheduleMap[constants.OrganizationID] = organizationId.String()

			scheduleModel, err := models.NewScheduleParams(
				uuid.New(),
				fmt.Sprintf("Schedule-%d", i),
				constants.DefaultScheduleGroup,
				config.ConnectorID,
				connectionId,
				temporalWorkflowID,
				constants.ActiveScheduleStatus,
				scheduleMap,
				scheduleMap[constants.CronExpression].(string),
			)

			if err != nil {

				return nil, err
			}
			schedulesModels = append(schedulesModels, scheduleModel)
		}
	}

	return schedulesModels, nil

}

func (g *gcs_connection_manager) ExtractAndValidateTemporalParams(schedules []dbmodels.CreateScheduleParams) ([]temporalModels.ExecuteWorkflowWithScheduleParams, error) {
	temporalScheduleParams := []temporalModels.ExecuteWorkflowWithScheduleParams{}

	for _, schedule := range schedules {
		config := map[string]interface{}{}
		err := json.Unmarshal(schedule.Config, &config)
		if err != nil {
			return nil, err
		}
		temporalScheduleParams = append(temporalScheduleParams, temporalModels.ExecuteWorkflowWithScheduleParams{
			ScheduleOptions: temporalModels.ScheduleOptions{
				ID: schedule.ID.String(),
				Spec: temporalio.ScheduleSpec{
					CronExpressions: []string{schedule.CronSchedule},
				},
				Action: &temporalio.ScheduleWorkflowAction{
					ID:       schedule.TemporalWorkflowID,
					Workflow: g.GetDefaultWorkflowType(),
					Args: []interface{}{
						map[string]interface{}{
							"connection_id": schedule.ConnectionID,
							"schedule_id":   schedule.ID,
							"config":        schedule.Config,
						},
					},
					Memo:      config,
					TaskQueue: constants.ConnectivityTemporalTaskQueue,
				},
				PauseOnFailure:     true,
				TriggerImmediately: true,
			},
		})
	}

	return temporalScheduleParams, nil
}

func (g *gcs_connection_manager) GetDefaultWorkflowType() string {
	return "gcs-connector"
}
