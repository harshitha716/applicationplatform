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

func TestSheetsGetQueryFilters(t *testing.T) {
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sheets" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'page' AND frap.resource_id = sheets.page_id AND frap.user_id = $1 AND frap.deleted_at IS NULL )`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(userId, "test@example.com", "Test User"))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupTestDB(t)
			sheet := &Sheet{}
			tt.setupMock(mock)

			baseQuery := db.Model(sheet)
			query := sheet.GetQueryFilters(baseQuery, tt.userId, tt.organizationIDs)

			var results []Sheet
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

func TestStructImplementsBaseModel_Sheet(t *testing.T) {
	var _ pgclient.BaseModel = &Sheet{}
}

func TestSheet_BeforeCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mock sqlmock.Sqlmock, userId uuid.UUID, pageId uuid.UUID)
		setupCtx  func() (context.Context, uuid.UUID)
		pageId    uuid.UUID
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
			pageId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, pageId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("page", pageId, userId, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
						"privilege",
						"created_at",
						"updated_at",
					}).AddRow(
						uuid.New(),
						"page",
						pageId,
						userId,
						"admin",
						time.Now(),
						time.Now(),
					))
			},
			wantErr: false,
		},
		{
			name: "failure - no user ID in context",
			setupCtx: func() (context.Context, uuid.UUID) {
				return context.Background(), uuid.Nil
			},
			pageId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, pageId uuid.UUID) {
				// No mock expectations needed as it should fail before DB query
			},
			wantErr: true,
			errMsg:  "no user id found in context",
		},
		{
			name: "failure - no admin privilege",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			pageId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, pageId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("page", pageId, userId, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
						"privilege",
						"created_at",
						"updated_at",
					}))
			},
			wantErr: true,
			errMsg:  "page access forbidden",
		},
		{
			name: "failure - database error",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			pageId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, pageId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("page", pageId, userId, "admin", 1).
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

			sheet := &Sheet{
				ID:        uuid.New(),
				PageId:    tt.pageId,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			tt.setupMock(mock, userId, tt.pageId)

			err := sheet.BeforeCreate(db)

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
