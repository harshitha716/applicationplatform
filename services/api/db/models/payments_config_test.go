package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPaymentsConfig_TableName(t *testing.T) {
	t.Parallel()
	config := PaymentsConfig{}
	assert.Equal(t, "payments_config", config.TableName())
}

func TestStructImplementsBaseModel_PaymentsConfig(t *testing.T) {
	var _ pgclient.BaseModel = &PaymentsConfig{}
}

func TestPaymentsConfig_GetQueryFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New()}

	tests := []struct {
		name      string
		setupMock func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "adds correct filter",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments_config" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'payments' AND frap.resource_id = payments_config.id AND frap.user_id = $1 AND frap.deleted_at IS NULL )`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupTestDB(t)
			config := &PaymentsConfig{}
			tt.setupMock(mock)

			baseQuery := db.Model(config)
			query := config.GetQueryFilters(baseQuery, userId, orgIds)

			var results []PaymentsConfig
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

func TestPaymentsConfig_BeforeCreate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		setupCtx    func() context.Context
		config      PaymentsConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation with system admin privilege",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             uuid.New(),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("organization", orgID, userID, PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "organization", orgID, userID, PrivilegeOrganizationSystemAdmin))
			},
			wantErr: false,
		},
		{
			name: "no user ID in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             uuid.New(),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No db expectations as function should exit early
			},
			wantErr:     true,
			errContains: "no user id found in context",
		},
		{
			name: "no permission",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             uuid.New(),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("organization", orgID, userID, PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
			},
			wantErr:     true,
			errContains: "organization access forbidden",
		},
		{
			name: "database error",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("organization", orgID, userID, PrivilegeOrganizationSystemAdmin, 1).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock := setupTestDB(t)

			ctx := tt.setupCtx()
			db.Statement.Context = ctx

			tt.setupMock(mock)

			err := tt.config.BeforeCreate(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPaymentsConfig_BeforeUpdate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	orgID := uuid.New()
	paymentConfigID := uuid.New()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		setupCtx    func() context.Context
		config      PaymentsConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update with payments admin privilege",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin))
			},
			wantErr: false,
		},
		{
			name: "no user ID in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No db expectations as function should exit early
			},
			wantErr:     true,
			errContains: "no user id found in context",
		},
		{
			name: "no permission",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
			},
			wantErr:     true,
			errContains: "payments access forbidden",
		},
		{
			name: "database error",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin, 1).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock := setupTestDB(t)

			ctx := tt.setupCtx()
			db.Statement.Context = ctx

			tt.setupMock(mock)

			err := tt.config.BeforeUpdate(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPaymentsConfig_BeforeDelete(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	orgID := uuid.New()
	paymentConfigID := uuid.New()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		setupCtx    func() context.Context
		config      PaymentsConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "successful delete with payments admin privilege",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin))
			},
			wantErr: false,
		},
		{
			name: "no user ID in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No db expectations as function should exit early
			},
			wantErr:     true,
			errContains: "no user id found in context",
		},
		{
			name: "no permission",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}))
			},
			wantErr:     true,
			errContains: "payments access forbidden",
		},
		{
			name: "database error",
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})
			},
			config: PaymentsConfig{
				OrganizationID: orgID,
				ID:             paymentConfigID,
			},

			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs(ResourceTypePayments, paymentConfigID, userID, PrivilegePaymentsAdmin, 1).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock := setupTestDB(t)

			ctx := tt.setupCtx()
			db.Statement.Context = ctx

			tt.setupMock(mock)

			err := tt.config.BeforeDelete(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
