package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAuditLogTableName(t *testing.T) {
	auditLog := AuditLog{}
	assert.Equal(t, "audit_logs", auditLog.TableName())
}

func TestAuditLogGetQueryFilters(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	auditLog := AuditLog{
		ResourceType: ResourceTypeOrganization,
		ResourceID:   uuid.New(),
	}

	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New(), uuid.New()}

	// Set up expectations for the SQL query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "audit_logs" WHERE EXISTS ( SELECT 1 FROM app.flattened_resource_audience_policies frap WHERE frap.resource_type = audit_logs.resource_type AND frap.resource_id = audit_logs.resource_id AND frap.user_id = $1 AND frap.deleted_at IS NULL )`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	result := auditLog.GetQueryFilters(gormDB, userId, orgIds)
	result.Find(&[]AuditLog{})

	// Verify the query was executed as expected
	assert.NotNil(t, result)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err, "there were unfulfilled expectations")
}

func TestAuditLogBeforeCreate(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name          string
		setupContext  func() *gorm.DB
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError string
	}{
		{
			name: "Success - user has access",
			setupContext: func() *gorm.DB {
				ctx := apicontext.AddAuthToContext(gormDB.Statement.Context, "test@example.com", userId, []uuid.UUID{})
				return gormDB.WithContext(ctx)
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id", "resource_type", "user_id"}).
						AddRow(orgId, ResourceTypeOrganization, userId))
			},
			expectedError: "",
		},
		{
			name: "Error - no user ID in context",
			setupContext: func() *gorm.DB {
				return gormDB.WithContext(context.Background())
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedError: "no user id found in context",
		},
		{
			name: "Error - user does not have access",
			setupContext: func() *gorm.DB {
				ctx := apicontext.AddAuthToContext(gormDB.Statement.Context, "test@example.com", userId, []uuid.UUID{})
				return gormDB.WithContext(ctx)
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id", "resource_type", "user_id"}))
			},
			expectedError: "user does not have access to the resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditLog := AuditLog{
				OrganizationID: orgId,
			}

			tt.mockSetup(mock)
			db := tt.setupContext()

			err := auditLog.BeforeCreate(db)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestAuditLogBeforeUpdate(t *testing.T) {
	auditLog := AuditLog{}
	err := auditLog.BeforeUpdate(nil)
	assert.Error(t, err)
	assert.Equal(t, "forbidden: audit logs cannot be updated", err.Error())
}

func TestAuditLogBeforeDelete(t *testing.T) {
	auditLog := AuditLog{}
	err := auditLog.BeforeDelete(nil)
	assert.Error(t, err)
	assert.Equal(t, "forbidden: audit logs cannot be deleted", err.Error())
}
