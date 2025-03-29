package service

import (
	"context"
	"sort"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"go.uber.org/zap"

	temporalsdk "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal"
	temporalModels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
)

type ScheduleServiceStore interface {
	store.ScheduleStore
}

type ScheduleService interface {
	CreateSchedules(ctx context.Context, params []models.CreateScheduleParams, tx store.ScheduleStore) error
	CreateSchedulesInTemporal(ctx context.Context, schedules []temporalModels.ExecuteWorkflowWithScheduleParams) error
	GetLastSyncedAt(ctx context.Context, connectionId uuid.UUID) (time.Time, error)
	GetSchedules(ctx context.Context, connectionId uuid.UUID) ([]models.Schedule, error)
	GetScheduleDetailsFromTemporal(ctx context.Context, scheduleId uuid.UUID) (temporalModels.QueryScheduleResponse, error)
}

type scheduleService struct {
	store       store.ScheduleStore
	temporalSdk temporalsdk.TemporalService
}

func NewScheduleService(appStore store.Store, temporalSdk temporalsdk.TemporalService) *scheduleService {
	return &scheduleService{
		store:       appStore,
		temporalSdk: temporalSdk,
	}
}

func (s *scheduleService) CreateSchedules(ctx context.Context, createScheduleParams []models.CreateScheduleParams, tx store.ScheduleStore) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	if tx == nil {
		tx = s.store
	}

	err := tx.CreateSchedules(ctx, createScheduleParams)
	if err != nil {
		logger.Error("Failed to create schedules", zap.Error(err))
		return err
	}

	return nil
}

func (s *scheduleService) GetSchedules(ctx context.Context, connectionId uuid.UUID) ([]models.Schedule, error) {
	return s.store.GetSchedulesByConnectionID(ctx, connectionId)
}

func (s *scheduleService) CreateSchedulesInTemporal(ctx context.Context, schedules []temporalModels.ExecuteWorkflowWithScheduleParams) error {
	for _, schedule := range schedules {
		_, err := s.temporalSdk.ExecuteScheduledWorkflow(
			ctx,
			schedule,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *scheduleService) GetLastSyncedAt(ctx context.Context, connectionId uuid.UUID) (time.Time, error) {
	schedules, err := s.store.GetSchedulesByConnectionID(ctx, connectionId)
	if err != nil {
		return time.Time{}, err
	}

	lastRuns := []time.Time{}
	for _, schedule := range schedules {
		temporalSchedule, err := s.temporalSdk.QuerySchedule(ctx, temporalModels.QueryScheduleParams{
			ScheduleID: schedule.ID.String(),
		})

		if err != nil {
			return time.Time{}, err
		}

		lastRun, err := s.getLastRunDetails(ctx, temporalSchedule)

		if err != nil {
			return time.Time{}, err
		}

		lastRuns = append(
			lastRuns,
			lastRun,
		)
	}

	if len(lastRuns) > 0 {
		sort.Slice(lastRuns, func(i, j int) bool {
			return lastRuns[i].After(lastRuns[j])
		})
	}
	return lastRuns[0], nil
}

func (s *scheduleService) getLastRunDetails(ctx context.Context, temporalSchedule temporalModels.QueryScheduleResponse) (time.Time, error) {
	for index := range temporalSchedule.Info.RecentActions {
		scheduleAction := temporalSchedule.Info.RecentActions[len(temporalSchedule.Info.RecentActions)-1-index]
		workflowID := scheduleAction.StartWorkflowResult.WorkflowID
		workflowDetails, err := s.temporalSdk.GetWorkflowDetails(ctx, temporalModels.GetWorkflowDetailsParams{
			WorkflowID: workflowID,
		})
		if err != nil {
			return time.Time{}, err
		}

		if workflowDetails.Details.Status != "RUNNING" {
			return time.Unix(workflowDetails.Details.StartTime, 0), nil
		}
	}
	return time.Time{}, nil
}

func (s *scheduleService) GetScheduleDetailsFromTemporal(ctx context.Context, scheduleId uuid.UUID) (temporalModels.QueryScheduleResponse, error) {
	temporalSchedule, err := s.temporalSdk.QuerySchedule(ctx, temporalModels.QueryScheduleParams{
		ScheduleID: scheduleId.String(),
	})
	if err != nil {
		return temporalModels.QueryScheduleResponse{}, err
	}

	return temporalSchedule, nil
}
