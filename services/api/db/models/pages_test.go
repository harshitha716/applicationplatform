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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPagesTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return db, mock
}

func TestPagesGetQueryFilters(t *testing.T) {
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pages" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'page' AND frap.resource_id = pages.page_id AND frap.user_id = $1 AND frap.deleted_at IS NULL )`)).
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

			db, mock := setupPagesTestDB(t)
			page := &Page{}
			tt.setupMock(mock)

			baseQuery := db.Model(page)
			query := page.GetQueryFilters(baseQuery, tt.userId, tt.organizationIDs)

			var results []Page
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

func TestStructImplementsBaseModel_Pages(t *testing.T) {
	var _ pgclient.BaseModel = &Page{}
}

func TestPage_BeforeCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mock sqlmock.Sqlmock, userId uuid.UUID, orgId uuid.UUID)
		setupCtx  func() (context.Context, uuid.UUID)
		orgId     uuid.UUID
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful creation with organization access",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			orgId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, orgId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
						"created_at",
						"updated_at",
					}).AddRow(
						uuid.New(),
						"organization",
						orgId,
						userId,
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
			orgId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, orgId uuid.UUID) {
				// No mock expectations needed as it should fail before DB query
			},
			wantErr: true,
			errMsg:  "no user id found in context",
		},
		{
			name: "failure - no organization access",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			orgId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, orgId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
						"created_at",
						"updated_at",
					}))
			},
			wantErr: true,
			errMsg:  "organization access forbidden",
		},
		{
			name: "failure - database error",
			setupCtx: func() (context.Context, uuid.UUID) {
				userId := uuid.New()
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
				return ctx, userId
			},
			orgId: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID, orgId uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
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

			page := &Page{
				ID:             uuid.New(),
				OrganizationId: tt.orgId,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			tt.setupMock(mock, userId, tt.orgId)

			err := page.BeforeCreate(db)

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
