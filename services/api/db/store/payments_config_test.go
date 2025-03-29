package store

import (
	"context"
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

func TestGetPaymentsConfigsByOrganizationId(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	configID := uuid.New()
	datasetID := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name           string
		organizationID string
		mockSetup      func(sqlmock.Sqlmock)
		expectedConfig models.PaymentsConfig
		wantErr        bool
	}{
		{
			name:           "success",
			organizationID: orgID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"organization_id",
					"accounts_dataset_id",
					"mapping_config",
					"status",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					configID,
					orgID,
					datasetID,
					[]byte(`{"key": "value"}`),
					models.PaymentsConfigStatusConnected,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "payments_config" WHERE organization_id = \$1 AND "payments_config"."deleted_at" IS NULL ORDER BY "payments_config"."id" LIMIT \$2`).
					WithArgs(orgID.String(), 1).
					WillReturnRows(rows)
			},
			expectedConfig: models.PaymentsConfig{
				ID:                configID,
				OrganizationID:    orgID,
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			wantErr: false,
		},
		{
			name:           "not found",
			organizationID: uuid.New().String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "payments_config" WHERE organization_id = \$1 AND "payments_config"."deleted_at" IS NULL ORDER BY "payments_config"."id" LIMIT \$2`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedConfig: models.PaymentsConfig{},
			wantErr:        true,
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

			config, err := store.GetPaymentsConfigsByOrganizationId(context.Background(), tt.organizationID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedConfig.ID, config.ID)
			assert.Equal(t, tt.expectedConfig.OrganizationID, config.OrganizationID)
			assert.Equal(t, tt.expectedConfig.AccountsDatasetID, config.AccountsDatasetID)
			assert.Equal(t, tt.expectedConfig.Status, config.Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreatePaymentsConfig(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	orgID := uuid.New()
	datasetID := uuid.New()
	configID := uuid.New()

	tests := []struct {
		name       string
		config     models.PaymentsConfig
		mockSetup  func(sqlmock.Sqlmock)
		wantErr    bool
		expectedID uuid.UUID
	}{
		{
			name: "success",
			config: models.PaymentsConfig{
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("organization", orgID, userID, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "organization", orgID, userID, models.PrivilegeOrganizationSystemAdmin))

				mock.ExpectQuery(`INSERT INTO "payments_config"`).
					WithArgs(
						orgID,                                // organization_id
						datasetID,                            // accounts_dataset_id
						[]byte(`{"key": "value"}`),           // mapping_config
						models.PaymentsConfigStatusConnected, // status
						nil,                                  // deleted_at
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(configID, time.Now(), time.Now()))
				mock.ExpectCommit()
			},
			wantErr:    false,
			expectedID: configID,
		},
		{
			name: "no user ID in context",
			config: models.PaymentsConfig{
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No db expectations as the function should exit early
			},
			wantErr: true,
		},
		{
			name: "permission check fails",
			config: models.PaymentsConfig{
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("organization", orgID, userID, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "database error",
			config: models.PaymentsConfig{
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("organization", orgID, userID, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "organization", orgID, userID, models.PrivilegeOrganizationSystemAdmin))

				mock.ExpectQuery(`INSERT INTO "payments_config"`).
					WithArgs(
						orgID,
						datasetID,
						[]byte(`{"key": "value"}`),
						models.PaymentsConfigStatusConnected,
						nil,
					).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "no organizations in context",
			config: models.PaymentsConfig{
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No db expectations as the function should exit early
			},
			wantErr: true,
		},
		{
			name: "multiple organizations in context",
			config: models.PaymentsConfig{
				AccountsDatasetID: datasetID,
				MappingConfig:     json.RawMessage(`{"key": "value"}`),
				Status:            models.PaymentsConfigStatusConnected,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No db expectations as the function should exit early
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

			var ctx context.Context
			if tt.name == "no user ID in context" {
				ctx = context.Background()
			} else if tt.name == "no organizations in context" {
				ctx = apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{})
			} else if tt.name == "multiple organizations in context" {
				ctx = apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID, uuid.New()})
			} else {
				ctx = apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			}

			result, err := store.CreatePaymentsConfig(ctx, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "no organizations in context" || tt.name == "multiple organizations in context" {
					assert.Contains(t, err.Error(), "organization access forbidden")
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, orgID, result.OrganizationID)
			assert.Equal(t, tt.config.AccountsDatasetID, result.AccountsDatasetID)
			assert.Equal(t, tt.config.MappingConfig, result.MappingConfig)
			assert.Equal(t, tt.config.Status, result.Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePaymentsConfig(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	configID := uuid.New()
	newConfig := json.RawMessage(`{"key": "updated_value"}`)

	tests := []struct {
		name      string
		configID  uuid.UUID
		config    json.RawMessage
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:     "success",
			configID: configID,
			config:   newConfig,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), models.ResourceTypePayments, uuid.New(), userID, models.PrivilegePaymentsAdmin))

				mock.ExpectExec(`UPDATE "payments_config" SET "mapping_config"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "payments_config"."deleted_at" IS NULL AND "id" = \$4`).
					WithArgs(
						[]byte(`{"key": "updated_value"}`),
						sqlmock.AnyArg(),
						configID,
						configID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "no user ID in context",
			configID: configID,
			config:   newConfig,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No db expectations as the function should exit early
			},
			wantErr: true,
		},
		{
			name:     "permission check fails",
			configID: configID,
			config:   newConfig,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:     "database error",
			configID: configID,
			config:   newConfig,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), models.ResourceTypePayments, uuid.New(), userID, models.PrivilegePaymentsAdmin))

				mock.ExpectExec(`UPDATE "payments_config" SET "mapping_config"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "payments_config"."deleted_at" IS NULL AND "id" = \$4`).
					WithArgs(
						[]byte(`{"key": "updated_value"}`),
						sqlmock.AnyArg(),
						configID,
						configID,
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

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			var ctx context.Context
			if tt.name == "no user ID in context" {
				ctx = context.Background()
			} else {
				ctx = apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{})
			}

			result, err := store.UpdatePaymentsConfig(ctx, tt.configID, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.configID, result.ID)
			assert.Equal(t, tt.config, result.MappingConfig)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePaymentsConfigStatus(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	configID := uuid.New()
	newStatus := models.PaymentsConfigStatusConnected

	tests := []struct {
		name      string
		configID  uuid.UUID
		status    models.PaymentsConfigStatus
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:     "success",
			configID: configID,
			status:   newStatus,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), models.ResourceTypePayments, uuid.New(), userID, models.PrivilegePaymentsAdmin))

				mock.ExpectExec(`UPDATE "payments_config" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "payments_config"."deleted_at" IS NULL AND "id" = \$4`).
					WithArgs(
						newStatus,
						sqlmock.AnyArg(),
						configID,
						configID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "no user ID in context",
			configID: configID,
			status:   newStatus,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No db expectations as the function should exit early
			},
			wantErr: true,
		},
		{
			name:     "permission check fails",
			configID: configID,
			status:   newStatus,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:     "database error",
			configID: configID,
			status:   newStatus,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), models.ResourceTypePayments, uuid.New(), userID, models.PrivilegePaymentsAdmin))

				mock.ExpectExec(`UPDATE "payments_config" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "payments_config"."deleted_at" IS NULL AND "id" = \$4`).
					WithArgs(
						newStatus,
						sqlmock.AnyArg(),
						configID,
						configID,
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

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			var ctx context.Context
			if tt.name == "no user ID in context" {
				ctx = context.Background()
			} else {
				ctx = apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{})
			}

			result, err := store.UpdatePaymentsConfigStatus(ctx, tt.configID, tt.status)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.configID, result.ID)
			assert.Equal(t, tt.status, result.Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePaymentsConfigById(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	configID := uuid.New()

	tests := []struct {
		name      string
		configID  uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:     "success",
			configID: configID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), models.ResourceTypePayments, uuid.New(), userID, models.PrivilegePaymentsAdmin))

				mock.ExpectExec(`UPDATE "payments_config" SET "deleted_at"=\$1 WHERE id = \$2 AND "payments_config"."deleted_at" IS NULL`).
					WithArgs(sqlmock.AnyArg(), configID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "no user ID in context",
			configID: configID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No db expectations as the function should exit early
			},
			wantErr: true,
		},
		{
			name:     "permission check fails",
			configID: configID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:     "database error",
			configID: configID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(models.ResourceTypePayments, sqlmock.AnyArg(), userID, models.PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), models.ResourceTypePayments, uuid.New(), userID, models.PrivilegePaymentsAdmin))

				mock.ExpectExec(`UPDATE "payments_config" SET "deleted_at"=\$1 WHERE id = \$2 AND "payments_config"."deleted_at" IS NULL`).
					WithArgs(sqlmock.AnyArg(), configID).
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

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			var ctx context.Context
			if tt.name == "no user ID in context" {
				ctx = context.Background()
			} else {
				ctx = apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{})
			}

			err := store.DeletePaymentsConfigById(ctx, tt.configID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
