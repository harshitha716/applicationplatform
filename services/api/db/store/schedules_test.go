package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// NewTestClient creates a new test database client
func NewTestClient() *pgclient.PostgresClient {
	return &pgclient.PostgresClient{
		DB: &gorm.DB{},
	}
}

func TestCreateSchedules(t *testing.T) {
	orgId := uuid.New()
	userId := uuid.New()

	configJSON := json.RawMessage(`{"key": "value"}`)
	orgID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		schedules []models.CreateScheduleParams
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful creation",
			ctx: apicontext.AddAuthToContext(
				context.Background(),
				"user",
				uuid.New(),
				[]uuid.UUID{orgID},
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
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId.String(), userId.String(), 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_type", "resource_audience_id", "privilege", "resource_type", "resource_id", "created_at", "updated_at", "deleted_at"}).AddRow("user", userId.String(), "read", "organization", orgId.String(), time.Now(), time.Now(), nil))
				mock.ExpectExec(`INSERT INTO "schedules"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}

			tt.mockSetup(mock)

			ctx := context.Background()

			ctx = apicontext.AddAuthToContext(ctx, "user_id", userId, []uuid.UUID{orgId})

			err := store.CreateSchedules(ctx, tt.schedules)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetSchedulesByConnectionID(t *testing.T) {
	orgId := uuid.New()

	tests := []struct {
		name      string
		connID    uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      []models.Schedule
		wantErr   bool
	}{
		{
			name:   "success",
			connID: uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				scheduleID := uuid.New()
				connectorID := uuid.New()
				now := time.Now()

				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"schedule_group",
					"connector_id",
					"connection_id",
					"temporal_workflow_id",
					"status",
					"config",
					"cron_schedule",
					"organization_id",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					scheduleID,
					"Test Schedule",
					"test-group",
					connectorID,
					uuid.New(),
					"test-workflow",
					"active",
					[]byte(`{"key": "value"}`),
					"0 * * * *",
					orgId,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE connection_id = $1 ORDER BY created_at DESC`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want: []models.Schedule{
				{
					ID:                 uuid.UUID{},
					Name:               "Test Schedule",
					ScheduleGroup:      "test-group",
					ConnectorID:        uuid.UUID{},
					ConnectionID:       uuid.UUID{},
					TemporalWorkflowID: "test-workflow",
					Status:             "active",
					Config:             []byte(`{"key": "value"}`),
					CronSchedule:       "0 * * * *",
					OrganizationID:     orgId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name:   "no schedules found",
			connID: uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"schedule_group",
					"connector_id",
					"connection_id",
					"temporal_workflow_id",
					"status",
					"config",
					"cron_schedule",
					"organization_id",
					"created_at",
					"updated_at",
					"deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE connection_id = $1 ORDER BY created_at DESC`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want:    []models.Schedule{},
			wantErr: false,
		},
		{
			name:   "database error",
			connID: uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE connection_id = $1 ORDER BY created_at DESC`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}

			tt.mockSetup(mock)

			got, err := store.GetSchedulesByConnectionID(context.Background(), tt.connID)
			if (err != nil) != tt.wantErr {
				t.Errorf("appStore.GetSchedulesByConnectionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("appStore.GetSchedulesByConnectionID() returned %d schedules, want %d", len(got), len(tt.want))
			}
		})
	}
}
