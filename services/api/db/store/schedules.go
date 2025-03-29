package store

import (
	"context"
	"fmt"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

type ScheduleStore interface {
	CreateSchedules(ctx context.Context, schedules []models.CreateScheduleParams) error
	GetSchedulesByConnectionID(ctx context.Context, connectionID uuid.UUID) ([]models.Schedule, error)
}

type scheduleStore struct {
	db *pgclient.PostgresClient
}

func NewScheduleStore(db *pgclient.PostgresClient) *scheduleStore {
	return &scheduleStore{
		db: db,
	}
}

func (s *appStore) CreateSchedules(ctx context.Context, schedules []models.CreateScheduleParams) error {
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)

	if len(orgIds) != 1 {
		return fmt.Errorf("organization access forbidden")
	}

	scheduleModels := []models.Schedule{}
	for _, schedule := range schedules {
		scheduleModel := models.Schedule{
			ID:                 schedule.ID,
			Name:               schedule.Name,
			ScheduleGroup:      schedule.ScheduleGroup,
			ConnectorID:        schedule.ConnectorID,
			ConnectionID:       schedule.ConnectionID,
			TemporalWorkflowID: schedule.TemporalWorkflowID,
			Status:             schedule.Status,
			Config:             schedule.Config,
			CronSchedule:       schedule.CronSchedule,
			OrganizationID:     orgIds[0],
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		scheduleModels = append(scheduleModels, scheduleModel)
	}

	return s.client.WithContext(ctx).Create(scheduleModels).Error
}

func (s *appStore) GetSchedulesByConnectionID(ctx context.Context, connectionID uuid.UUID) ([]models.Schedule, error) {
	db := s.client.WithContext(ctx)
	var schedules []models.Schedule

	result := db.Where("connection_id = ?", connectionID.String()).Order("created_at DESC").Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}

	return schedules, nil
}
