package store_test

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateAuditLog(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	userId := uuid.New()
	logId := uuid.New()

	payload, _ := json.Marshal(map[string]interface{}{
		"test": "value",
	})

	tests := []struct {
		name      string
		auditLog  models.AuditLog
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful audit log creation",
			auditLog: models.AuditLog{
				Kind:           models.AuditLogKindInfo,
				OrganizationID: orgId,
				ResourceType:   models.ResourceTypeOrganization,
				ResourceID:     orgId,
				EventName:      "test_event",
				Payload:        payload,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock the insert
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id", "resource_type", "user_id"}).
						AddRow(orgId, models.ResourceTypeOrganization, userId))
				mock.ExpectQuery(`INSERT INTO "audit_logs"`).
					WithArgs(
						models.AuditLogKindInfo,
						orgId,
						sqlmock.AnyArg(), // IP Address
						sqlmock.AnyArg(), // User Email
						sqlmock.AnyArg(), // User Agent
						models.ResourceTypeOrganization,
						orgId,
						"test_event",
						payload,
						sqlmock.AnyArg(), // Created At
					).
					WillReturnRows(sqlmock.NewRows([]string{"audit_log_id"}).AddRow(logId.String()))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "user not in organization",
			auditLog: models.AuditLog{
				Kind:           models.AuditLogKindInfo,
				OrganizationID: orgId,
				ResourceType:   models.ResourceTypeOrganization,
				ResourceID:     orgId,
				EventName:      "test_event",
				Payload:        payload,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id", "resource_type", "user_id"}).
						AddRow(orgId, models.ResourceTypeOrganization, userId))
				mock.ExpectQuery(`INSERT INTO "audit_logs"`).
					WithArgs(
						models.AuditLogKindInfo,
						orgId,
						sqlmock.AnyArg(), // IP Address
						sqlmock.AnyArg(), // User Email
						sqlmock.AnyArg(), // User Agent
						models.ResourceTypeOrganization,
						orgId,
						"test_event",
						payload,
						sqlmock.AnyArg(), // Created At
					).
					WillReturnRows(sqlmock.NewRows([]string{"audit_log_id"}).AddRow(logId.String()))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			auditLog: models.AuditLog{
				Kind:           models.AuditLogKindInfo,
				OrganizationID: orgId,
				ResourceType:   models.ResourceTypeOrganization,
				ResourceID:     orgId,
				EventName:      "test_event",
				Payload:        payload,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock the insert with error
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id", "resource_type", "user_id"}).
						AddRow(orgId, models.ResourceTypeOrganization, userId))
				mock.ExpectQuery(`INSERT INTO "audit_logs"`).
					WithArgs(
						models.AuditLogKindInfo,
						orgId,
						sqlmock.AnyArg(), // IP Address
						sqlmock.AnyArg(), // User Email
						sqlmock.AnyArg(), // User Agent
						models.ResourceTypeOrganization,
						orgId,
						"test_event",
						payload,
						sqlmock.AnyArg(), // Created At
					).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			auditLogStore, cleanup := store.NewStore(&pgclient.PostgresClient{DB: gormDB})
			defer cleanup()

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{})

			tt.mockSetup(mock)

			// Execute
			result, err := auditLogStore.CreateAuditLog(ctx, tt.auditLog)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.auditLog.Kind, result.Kind)
				assert.Equal(t, tt.auditLog.OrganizationID, result.OrganizationID)
				assert.Equal(t, tt.auditLog.ResourceType, result.ResourceType)
				assert.Equal(t, tt.auditLog.ResourceID, result.ResourceID)
				assert.Equal(t, tt.auditLog.EventName, result.EventName)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAuditLogsByOrganizationId(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()

	tests := []struct {
		name      string
		orgId     uuid.UUID
		kind      models.AuditLogKind
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		wantCount int
	}{
		{
			name:  "successful retrieval",
			orgId: orgId,
			kind:  models.AuditLogKindInfo,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"audit_log_id", "kind", "organization_id", "ip_address", "user_email",
					"user_agent", "resource_type", "resource_id", "event_name", "payload", "created_at",
				}).
					AddRow(
						uuid.New(), models.AuditLogKindInfo, orgId, "127.0.0.1", "user@example.com",
						"test-agent", models.ResourceTypeOrganization, orgId, "test_event", []byte(`{"test":"value"}`), time.Now(),
					)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "audit_logs" WHERE kind = $1 AND organization_id = $2 ORDER BY created_at DESC`)).
					WithArgs(models.AuditLogKindInfo, orgId).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:  "no logs found",
			orgId: orgId,
			kind:  models.AuditLogKindInfo,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"audit_log_id", "kind", "organization_id", "ip_address", "user_email",
					"user_agent", "resource_type", "resource_id", "event_name", "payload", "created_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "audit_logs" WHERE kind = $1 AND organization_id = $2 ORDER BY created_at DESC`)).
					WithArgs(models.AuditLogKindInfo, orgId).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:  "database error",
			orgId: orgId,
			kind:  models.AuditLogKindInfo,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "audit_logs" WHERE kind = $1 AND organization_id = $2 ORDER BY created_at DESC`)).
					WithArgs(orgId, models.AuditLogKindInfo).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			auditLogStore, cleanup := store.NewStore(&pgclient.PostgresClient{DB: gormDB})
			defer cleanup()

			tt.mockSetup(mock)

			// Execute
			logs, err := auditLogStore.GetAuditLogsByOrganizationId(context.Background(), tt.orgId, tt.kind)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, logs, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, tt.orgId, logs[0].OrganizationID)
				assert.Equal(t, tt.kind, logs[0].Kind)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
