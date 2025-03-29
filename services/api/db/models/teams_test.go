package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTeam_TableName(t *testing.T) {
	team := &Team{}
	assert.Equal(t, "teams", team.TableName())
}

func TestTeam_GetQueryFilters(t *testing.T) {
	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New(), uuid.New()}

	db, mock := setupTestDB(t)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = $1 AND frap.resource_id = teams.organization_id AND frap.user_id = $2 ) AND "teams"."deleted_at" IS NULL`)

	mock.ExpectQuery(expectedSQL).
		WithArgs(ResourceTypeOrganization, userId).
		WillReturnRows(sqlmock.NewRows([]string{"team_id"}))

	team := &Team{}
	filteredDB := team.GetQueryFilters(db, userId, orgIds)

	var result []Team
	err := filteredDB.Find(&result).Error
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTeam_ValidateContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		wantErr     bool
		errContains string
	}{
		{
			name:        "error - no user in context",
			ctx:         context.Background(),
			wantErr:     true,
			errContains: "no user id found in context",
		},
		{
			name:        "error - no org ids",
			ctx:         apicontext.AddAuthToContext(context.Background(), "user", uuid.New(), []uuid.UUID{}),
			wantErr:     true,
			errContains: "user is not part of the organization",
		},
		{
			name: "success",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", uuid.New(), []uuid.UUID{uuid.New()}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team := &Team{}
			userId, orgIds, err := team.validateContext(tt.ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, userId)
				assert.Nil(t, orgIds)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userId)
				assert.NotEmpty(t, orgIds)
			}
		})
	}
}

func TestTeam_BeforeCreate(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setup       func() (*Team, *gorm.DB, sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				return &Team{OrganizationID: orgId}, db, mock
			},
		},
		{
			name: "error - wrong org id",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)
				return &Team{OrganizationID: uuid.New()}, db, mock
			},
			wantErr:     true,
			errContains: "user is not part of the organization",
		},
		{
			name: "error - no privileges",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				return &Team{OrganizationID: orgId}, db, mock
			},
			wantErr:     true,
			errContains: "user does not have privileges to create a team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, db, _ := tt.setup()
			err := team.BeforeCreate(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTeam_BeforeUpdate(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setup       func() (*Team, *gorm.DB, sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				return &Team{OrganizationID: orgId}, db, mock
			},
		},
		{
			name: "error - wrong org id",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)
				return &Team{OrganizationID: uuid.New()}, db, mock
			},
			wantErr:     true,
			errContains: "user is not part of the organization",
		},
		{
			name: "error - no privileges",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				return &Team{OrganizationID: orgId}, db, mock
			},
			wantErr:     true,
			errContains: "user does not have privileges to update a team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, db, _ := tt.setup()
			err := team.BeforeUpdate(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTeam_BeforeDelete(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setup       func() (*Team, *gorm.DB, sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				return &Team{OrganizationID: orgId}, db, mock
			},
		},
		{
			name: "error - wrong org id",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)
				return &Team{OrganizationID: uuid.New()}, db, mock
			},
			wantErr:     true,
			errContains: "user is not part of the organization",
		},
		{
			name: "error - no privileges",
			setup: func() (*Team, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				return &Team{OrganizationID: orgId}, db, mock
			},
			wantErr:     true,
			errContains: "user does not have privileges to delete a team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, db, _ := tt.setup()
			err := team.BeforeDelete(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
