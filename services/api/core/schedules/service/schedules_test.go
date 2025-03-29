package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	temporalsdk "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStore struct {
	mock.Mock
	store.Store
}

type mockTemporalService struct {
	mock.Mock
	temporalsdk.TemporalService
}

func (m *mockStore) CreateSchedules(ctx context.Context, schedules []models.CreateScheduleParams) error {
	args := m.Called(ctx, schedules)
	return args.Error(0)
}

func (m *mockStore) CreateSchedulePolicy(ctx context.Context, scheduleID uuid.UUID, audienceType models.AudienceType, audienceID uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	args := m.Called(ctx, scheduleID, audienceType, audienceID, privilege)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ResourceAudiencePolicy), args.Error(1)
}

func TestCreateSchedules(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()
	configJSON := json.RawMessage(`{"key": "value"}`)

	tests := []struct {
		name      string
		ctx       context.Context
		schedules []models.CreateScheduleParams
		mockSetup func(*mockStore, *mockTemporalService)
		wantErr   bool
	}{
		{
			name: "successful creation",
			ctx: apicontext.AddAuthToContext(
				context.Background(),
				"user",
				userId,
				[]uuid.UUID{orgId},
			),
			schedules: []models.CreateScheduleParams{
				{
					ID:                 uuid.New(),
					Name:               "Test Schedule",
					ScheduleGroup:      "test-group",
					ConnectorID:        uuid.New(),
					ConnectionID:       uuid.New(),
					TemporalWorkflowID: "test-workflow",
					Status:             "active",
					Config:             configJSON,
					CronSchedule:       "0 * * * *",
				},
			},
			mockSetup: func(m *mockStore, mts *mockTemporalService) {
				m.On("CreateSchedules", mock.Anything, mock.Anything).Return(nil)
				mts.On("ExecuteScheduledWorkflow", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mockStore{}
			mockTemporalService := &mockTemporalService{}
			tt.mockSetup(mockStore, mockTemporalService)

			service := NewScheduleService(mockStore, mockTemporalService)
			err := service.CreateSchedules(tt.ctx, tt.schedules, nil)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockStore.AssertExpectations(t)
			}
		})
	}
}
