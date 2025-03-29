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

func TestOrganizationSSOConfig_TableName(t *testing.T) {
	config := &OrganizationSSOConfig{}
	assert.Equal(t, "organization_sso_configs", config.TableName())
}

func TestOrganizationSSOConfig_GetQueryFilters(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		wantWhere bool
	}{
		{
			name:      "admin role returns unmodified query",
			role:      "admin",
			wantWhere: false,
		},
		{
			name:      "non-admin role adds where clause",
			role:      "user",
			wantWhere: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupTestDB(t)
			ctx := apicontext.AddAuthToContext(context.Background(), tt.role, uuid.New(), []uuid.UUID{uuid.New()})
			db.Statement.Context = ctx

			config := &OrganizationSSOConfig{}
			userId := uuid.New()
			orgIds := []uuid.UUID{uuid.New()}

			result := config.GetQueryFilters(db, userId, orgIds)

			if tt.wantWhere {
				assert.NotEqual(t, db, result, "Query should be modified for non-admin")
			} else {
				assert.Equal(t, db, result, "Query should be unmodified for admin")
			}
		})
	}
}

func TestOrganizationSSOConfig_BeforeCreate(t *testing.T) {

	userId := uuid.New()
	orgId := uuid.New()
	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - user has access",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgId))
			},
			wantErr: false,
		},
		{
			name: "error - user has no access",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}))
			},
			wantErr:     true,
			errContains: "user does not have access to organization",
		},
		{
			name: "error - db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnError(assert.AnError)
			},
			wantErr:     true,
			errContains: "failed to check if user has access to organization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			db.Statement.Context = ctx

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			config := &OrganizationSSOConfig{
				OrganizationID: orgId,
			}
			err := config.BeforeCreate(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrganizationSSOConfig_BeforeUpdate(t *testing.T) {
	orgId := uuid.New()
	userId := uuid.New()
	tests := []struct {
		name        string
		role        string
		setupMock   func(sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name:        "error - non-admin user",
			role:        "user",
			wantErr:     true,
			errContains: "only admin can update",
		},
		{
			name: "success - admin user with access",
			role: "admin",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock finding policies that grant access
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgId))
			},
			wantErr: false,
		},
		{
			name: "error - admin user without access",
			role: "admin",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock finding no policies
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(ResourceTypeOrganization, orgId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}))
			},
			wantErr:     true,
			errContains: "user does not have access to organization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			ctx := apicontext.AddAuthToContext(context.Background(), tt.role, userId, []uuid.UUID{orgId})
			db.Statement.Context = ctx

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			config := &OrganizationSSOConfig{
				OrganizationID: orgId,
			}
			err := config.BeforeUpdate(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModelImplementsBaseModel(t *testing.T) {
	assert.Implements(t, (*pgclient.BaseModel)(nil), new(OrganizationSSOConfig))
}
