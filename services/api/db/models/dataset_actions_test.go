package models

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDatasetAction_TableName(t *testing.T) {
	t.Parallel()
	datasetAction := DatasetAction{}
	assert.Equal(t, "dataset_actions", datasetAction.TableName())
}

func TestDatasetAction_GetQueryFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	org1ID := uuid.New()
	orgIDs := []uuid.UUID{org1ID}

	tests := []struct {
		name            string
		userId          uuid.UUID
		organizationIDs []uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock)
		wantErr         bool
	}{
		{
			name:            "right filter organization memberships",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dataset_actions" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'dataset' AND frap.resource_id = dataset_actions.dataset_id AND frap.user_id = $1 AND frap.deleted_at IS NULL )`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "dataset_id"}).
						AddRow(uuid.New(), uuid.New()))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupTestDB(t)
			datasetAction := &DatasetAction{}
			tt.setupMock(mock)

			baseQuery := db.Model(datasetAction)
			query := datasetAction.GetQueryFilters(baseQuery, tt.userId, tt.organizationIDs)

			var results []DatasetAction
			err := query.Find(&results).Error

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStructImplementsBaseModel_DatasetActions(t *testing.T) {
	var _ pgclient.BaseModel = &DatasetAction{}
}

func TestDatasetAction_BeforeCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mock sqlmock.Sqlmock, userId uuid.UUID)
		setupCtx  func() (context.Context, uuid.UUID)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful creation with admin privilege",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userId, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userId, "admin"))
			},
			wantErr: false,
		},
		{
			name: "failure - no user ID in context",
			setupCtx: func() (context.Context, uuid.UUID) {
				return context.Background(), uuid.Nil
			},
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
				// No mock expectations needed as it should fail before DB query
			},
			wantErr: true,
			errMsg:  "no user id found in context",
		},
		{
			name: "failure - no admin privilege found",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userId, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
			},
			wantErr: true,
			errMsg:  "dataset write access forbidden",
		},
		{
			name: "failure - database error",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userId, "admin", 1).
					WillReturnError(fmt.Errorf("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock := setupTestDB(t)

			ctx, userId := tt.setupCtx()
			db = db.WithContext(ctx)

			datasetAction := &DatasetAction{
				ID:         uuid.New(),
				DatasetId:  uuid.New(),
				ActionType: "test",
				ActionId:   "test-action",
				Status:     "pending",
				StartedAt:  time.Now(),
				ActionBy:   userId,
			}

			tt.setupMock(mock, userId)

			err := datasetAction.BeforeCreate(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
