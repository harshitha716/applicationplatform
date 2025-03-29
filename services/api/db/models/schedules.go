package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Schedule struct {
	ID                 uuid.UUID       `json:"id" gorm:"column:id"`
	Name               string          `json:"name" gorm:"column:name"`
	ScheduleGroup      string          `json:"schedule_group" gorm:"column:schedule_group"`
	ConnectorID        uuid.UUID       `json:"connector_id" gorm:"column:connector_id"`
	ConnectionID       uuid.UUID       `json:"connection_id" gorm:"column:connection_id"`
	TemporalWorkflowID string          `json:"temporal_workflow_id" gorm:"column:temporal_workflow_id"`
	Status             string          `json:"status" gorm:"column:status"`
	Config             json.RawMessage `json:"config" gorm:"column:config"`
	CronSchedule       string          `json:"cron_schedule" gorm:"column:cron_schedule"`
	CreatedAt          time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt          time.Time       `json:"deleted_at" gorm:"column:deleted_at"`

	OrganizationID uuid.UUID `json:"organization_id" gorm:"column:organization_id"`
}

type CreateScheduleParams struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	ScheduleGroup      string          `json:"schedule_group"`
	ConnectorID        uuid.UUID       `json:"connector_id"`
	ConnectionID       uuid.UUID       `json:"connection_id"`
	TemporalWorkflowID string          `json:"temporal_workflow_id"`
	Status             string          `json:"status"`
	Config             json.RawMessage `json:"config"`
	CronSchedule       string          `json:"cron_schedule"`
}

func NewScheduleParams(id uuid.UUID, name, scheduleGroup string, connectorID uuid.UUID, connectionID uuid.UUID, temporalWorkflowID, status string, config map[string]interface{}, cronSchedule string) (CreateScheduleParams, error) {
	configJson, err := json.Marshal(config)
	if err != nil {
		return CreateScheduleParams{}, err
	}
	return CreateScheduleParams{
		ID:                 id,
		Name:               name,
		ScheduleGroup:      scheduleGroup,
		ConnectorID:        connectorID,
		ConnectionID:       connectionID,
		TemporalWorkflowID: temporalWorkflowID,
		Status:             status,
		Config:             configJson,
		CronSchedule:       cronSchedule,
	}, nil
}

func (s *Schedule) TableName() string {
	return "schedules"
}

func (s *Schedule) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'connection'
			AND frap.resource_id = schedules.connection_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (s *Schedule) BeforeCreate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", "organization", s.OrganizationID, userId).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil

}
