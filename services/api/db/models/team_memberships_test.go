package models

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTeamMembership_TableName(t *testing.T) {
	tm := &TeamMembership{}
	assert.Equal(t, "team_memberships", tm.TableName())
}

func TestTeamMembership_GetQueryFilters(t *testing.T) {
	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New(), uuid.New()}

	db, mock := setupTestDB(t)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap, "app"."teams" teams WHERE frap.resource_type = $1 AND frap.resource_id = teams.organization_id AND teams.team_id = team_memberships.team_id AND frap.user_id = $2 ) AND "team_memberships"."deleted_at" IS NULL`)

	mock.ExpectQuery(expectedSQL).
		WithArgs(ResourceTypeOrganization, userId).
		WillReturnRows(sqlmock.NewRows([]string{"team_membership_id"}))

	tm := &TeamMembership{}
	filteredDB := tm.GetQueryFilters(db, userId, orgIds)

	var result []TeamMembership
	err := filteredDB.Find(&result).Error
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTeamMembership_ValidateContext(t *testing.T) {
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
			tm := &TeamMembership{}
			userId, orgIds, err := tm.validateContext(tt.ctx)

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

func TestTeamMembership_BeforeCreate(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setup       func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setup: func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				fmt.Println(userId, teamId, orgId)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id"}).AddRow(teamId))

				return &TeamMembership{TeamID: teamId}, db, mock
			},
		},
		{
			name: "error - no teams found",
			setup: func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id"}))

				return &TeamMembership{TeamID: teamId}, db, mock
			},
			wantErr:     true,
			errContains: "no access to modify team memberships",
		},
		{
			name: "db error",
			setup: func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnError(fmt.Errorf("db error"))

				return &TeamMembership{TeamID: teamId}, db, mock
			},
			wantErr:     true,
			errContains: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, db, _ := tt.setup()
			err := tm.BeforeCreate(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTeamMembership_BeforeUpdate(t *testing.T) {
	tm := &TeamMembership{}
	err := tm.BeforeUpdate(&gorm.DB{})
	assert.EqualError(t, err, "update forbidden")
}

func TestTeamMembership_BeforeDelete(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()

	tests := []struct {
		name        string
		setup       func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setup: func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id"}).AddRow(teamId))

				return &TeamMembership{TeamID: teamId}, db, mock
			},
		},
		{
			name: "error - no teams found",
			setup: func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id"}))

				return &TeamMembership{TeamID: teamId}, db, mock
			},
			wantErr:     true,
			errContains: "no access to modify team memberships",
		},
		{
			name: "db error",
			setup: func() (*TeamMembership, *gorm.DB, sqlmock.Sqlmock) {
				db, mock := setupTestDB(t)
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()})
				db = db.WithContext(ctx)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).
					WillReturnError(fmt.Errorf("db error"))

				return &TeamMembership{TeamID: teamId}, db, mock
			},
			wantErr:     true,
			errContains: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, db, _ := tt.setup()
			err := tm.BeforeDelete(db)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
