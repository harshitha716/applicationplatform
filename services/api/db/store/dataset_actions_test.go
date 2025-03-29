package store

import (
	"context"
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

func TestCreateDatasetAction(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	datasetID := uuid.New()
	actionID := uuid.New().String()
	actorID := uuid.New()

	tests := []struct {
		name      string
		orgID     uuid.UUID
		params    models.CreateDatasetActionParams
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:  "success - not completed",
			orgID: orgID,
			params: models.CreateDatasetActionParams{
				ActionId:    actionID,
				DatasetId:   datasetID,
				ActionType:  "PROCESS",
				Status:      "PENDING",
				Config:      map[string]interface{}{"key": "value"},
				ActionBy:    actorID,
				IsCompleted: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO "dataset_actions"`).
					WithArgs(
						sqlmock.AnyArg(),
						actionID,
						"PROCESS",
						datasetID,
						orgID,
						"PENDING",
						[]byte(`{"key":"value"}`),
						actorID,
						sqlmock.AnyArg(),
						nil,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "success - completed",
			orgID: orgID,
			params: models.CreateDatasetActionParams{
				ActionId:    actionID,
				DatasetId:   datasetID,
				ActionType:  "PROCESS",
				Status:      "COMPLETED",
				Config:      map[string]interface{}{"key": "value"},
				ActionBy:    actorID,
				IsCompleted: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO "dataset_actions"`).
					WithArgs(
						sqlmock.AnyArg(),
						actionID,
						"PROCESS",
						datasetID,
						orgID,
						"COMPLETED",
						[]byte(`{"key":"value"}`),
						actorID,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "error - invalid dataset id",
			orgID: orgID,
			params: models.CreateDatasetActionParams{
				ActionId:    actionID,
				DatasetId:   uuid.Nil,
				ActionType:  "PROCESS",
				Status:      "PENDING",
				Config:      map[string]interface{}{"key": "value"},
				ActionBy:    actorID,
				IsCompleted: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO "dataset_actions"`).
					WithArgs(
						sqlmock.AnyArg(),
						actionID,
						"PROCESS",
						uuid.Nil,
						orgID,
						"PENDING",
						[]byte(`{"key": "value"}`),
						actorID,
						sqlmock.AnyArg(),
						nil,
					).
					WillReturnError(gorm.ErrInvalidField)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}

			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
				WithArgs("dataset", sqlmock.AnyArg(), actorID, "admin", 1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
					AddRow(uuid.New(), "dataset", datasetID, actorID, "admin"))
			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "role", actorID, []uuid.UUID{})

			err := store.CreateDatasetAction(ctx, tt.orgID, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetDatasetActions(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	datasetID := uuid.New()
	actionID := uuid.New().String()
	actorID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		orgID     uuid.UUID
		filters   models.DatasetActionFilters
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - with filters",
			orgID: orgID,
			filters: models.DatasetActionFilters{
				DatasetIds: []uuid.UUID{datasetID},
				ActionIds:  []string{actionID},
				ActionType: []string{"PROCESS"},
				ActionBy:   []uuid.UUID{actorID},
				Status:     []string{"PENDING"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "organization_id", "action_id", "action_type", "dataset_id",
					"status", "config", "action_by", "started_at", "completed_at",
				}).AddRow(
					uuid.New(), orgID, actionID, "PROCESS", datasetID,
					"PENDING", []byte(`{"key":"value"}`), actorID, now, nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "dataset_actions" WHERE organization_id = \$1`).
					WithArgs(
						orgID,
						datasetID,
						actionID,
						"PROCESS",
						actorID,
						"PENDING",
					).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:    "success - no filters",
			orgID:   orgID,
			filters: models.DatasetActionFilters{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "dataset_actions" WHERE organization_id = \$1`).
					WithArgs(orgID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "organization_id", "action_id", "action_type", "dataset_id",
						"status", "config", "action_by", "started_at", "completed_at",
					}))
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			actions, err := store.GetDatasetActions(context.Background(), tt.orgID, tt.filters)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, actions, tt.wantCount)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateDatasetActionConfig(t *testing.T) {
	t.Parallel()

	actionID := uuid.New().String()
	config := map[string]interface{}{
		"key":    "value",
		"number": 42,
	}

	tests := []struct {
		name      string
		actionID  string
		config    map[string]interface{}
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:     "success",
			actionID: actionID,
			config:   config,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dataset_actions" SET "config"=$1 WHERE action_id = $2`)).
					WithArgs(
						[]byte(`{"key":"value","number":42}`),
						actionID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "error - no rows affected",
			actionID: actionID,
			config:   config,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dataset_actions" SET "config"=$1 WHERE action_id = $2`)).
					WithArgs(
						[]byte(`{"key":"value","number":42}`),
						actionID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "error - invalid action id",
			actionID: "",
			config:   config,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dataset_actions" SET "config"=$1 WHERE action_id = $2`)).
					WithArgs(
						[]byte(`{"key":"value","number":42}`),
						"",
					).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			err := store.UpdateDatasetActionConfig(context.Background(), tt.actionID, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
