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

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

// Dataset Tests
func TestDataset_TableName(t *testing.T) {
	t.Parallel()
	dataset := Dataset{}
	assert.Equal(t, "datasets", dataset.TableName())
}

func TestDataset_GetQueryFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	org1ID := uuid.New()
	org2ID := uuid.New()
	orgIDs := []uuid.UUID{org1ID, org2ID}
	datasetID := uuid.New()

	tests := []struct {
		name            string
		userId          uuid.UUID
		organizationIDs []uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock)
		wantSQL         string
		wantErr         bool
	}{
		{
			name:            "user and organization access",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"dataset_id"}).AddRow(datasetID)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnRows(rows)
			},
		},
		{
			name:            "empty organization list",
			userId:          userId,
			organizationIDs: []uuid.UUID{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"dataset_id"}))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupTestDB(t)
			dataset := &Dataset{ID: datasetID}
			tt.setupMock(mock)

			query := dataset.GetQueryFilters(db.Model(&Dataset{}), tt.userId, tt.organizationIDs)

			// Verify query executes without error
			var results []Dataset
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

func TestStructImplementsBaseModel_Dataset(t *testing.T) {
	var _ pgclient.BaseModel = &Dataset{}
}

func TestDataset_BeforeCreate(t *testing.T) {
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

			dataset := &Dataset{
				ID:             uuid.New(),
				OrganizationId: tt.orgId,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			tt.setupMock(mock, userId, tt.orgId)

			err := dataset.BeforeCreate(db)

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
