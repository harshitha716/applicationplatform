package snowflake

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/core/connections/constants"
	"github.com/Zampfi/application-platform/services/api/db/models"
	temporalModels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"github.com/google/uuid"
	temporalio "go.temporal.io/sdk/client"
)

type snowflake_connection_manager struct {
}

func InitSnowflakeConnectionManager() *snowflake_connection_manager {
	return &snowflake_connection_manager{}
}

func (s *snowflake_connection_manager) CheckConnection() error {
	return nil
}

func (s *snowflake_connection_manager) ExtractAndValidateDefaultSchedules(connectionId uuid.UUID, organizationId uuid.UUID, config models.CreateConnectionParams) ([]models.CreateScheduleParams, error) {
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
			temporalWorkflowID := fmt.Sprintf("snowflake_workflow:%s", id)
			scheduleMap[constants.OrganizationID] = organizationId.String()
			scheduleMap[constants.TemporalWorkflowID] = temporalWorkflowID

			scheduleMap[constants.SnowflakeDatabase] = config.Config[constants.SnowflakeDatabase].(string)
			scheduleMap[constants.SnowflakeSchema] = config.Config[constants.SnowflakeSchema].(string)
			scheduleMap[constants.SnowflakeTable] = config.Config[constants.SnowflakeTable].(string)
			scheduleMap[constants.SnowflakeFilterColumnName] = config.Config[constants.SnowflakeFilterColumnName].(string)
			scheduleMap[constants.S3Bucket] = config.Config[constants.S3Bucket].(string)
			scheduleMap[constants.S3DestinationPath] = config.Config[constants.S3DestinationPath].(string)
			scheduleMap[constants.S3StorageIntegration] = config.Config[constants.S3StorageIntegration].(string)

			if startOffsetDays, ok := config.Config[constants.StartOffsetDays]; ok {
				scheduleMap[constants.StartOffsetDays] = startOffsetDays.(int)
			}
			if endOffsetDays, ok := config.Config[constants.EndOffsetDays]; ok {
				scheduleMap[constants.EndOffsetDays] = endOffsetDays.(int)
			}

			scheduleModel, err := models.NewScheduleParams(
				uuid.New(),
				fmt.Sprintf("Schedule-snowflake_workflow%d", i),
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

func (s *snowflake_connection_manager) ExtractAndValidateTemporalParams(schedules []models.CreateScheduleParams) ([]temporalModels.ExecuteWorkflowWithScheduleParams, error) {
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
					Workflow: s.GetDefaultWorkflowType(),
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

func (s *snowflake_connection_manager) GetDefaultWorkflowType() string {
	return "snowflake-connector"
}
