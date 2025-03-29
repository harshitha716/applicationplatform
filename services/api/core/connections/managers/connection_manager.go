package manager

import (
	"github.com/google/uuid"

	"github.com/Zampfi/application-platform/services/api/db/models"
	temporalModels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
)

type ConnectionManager interface {
	CheckConnection() error
	ExtractAndValidateDefaultSchedules(connectionId uuid.UUID, organizationId uuid.UUID, config models.CreateConnectionParams) ([]models.CreateScheduleParams, error)
	GetDefaultWorkflowType() string
	ExtractAndValidateTemporalParams(schedules []models.CreateScheduleParams) ([]temporalModels.ExecuteWorkflowWithScheduleParams, error)
}
