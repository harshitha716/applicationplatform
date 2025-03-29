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
)

func TestGetTeams(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teamId1 := uuid.New()
	teamId2 := uuid.New()

	tests := []struct {
		name          string
		orgId         uuid.UUID
		mockSetup     func(sqlmock.Sqlmock)
		expectedTeams []models.Team
		wantErr       bool
	}{
		{
			name:  "success",
			orgId: orgId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"team_id", "organization_id"}).
					AddRow(teamId1, orgId).
					AddRow(teamId2, orgId)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE organization_id = $1 AND "teams"."deleted_at" IS NULL ORDER BY name ASC`)).
					WithArgs(orgId).
					WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE "team_memberships"."team_id" IN ($1,$2) AND "team_memberships"."deleted_at" IS NULL`)).
					WithArgs(teamId1, teamId2).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "team_id", "user_id"}).
						AddRow(uuid.New(), teamId1, uuid.New()).
						AddRow(uuid.New(), teamId2, uuid.New()))
			},
			expectedTeams: []models.Team{
				{TeamID: teamId1, OrganizationID: orgId},
				{TeamID: teamId2, OrganizationID: orgId},
			},
		},
		{
			name:  "no teams found",
			orgId: orgId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE organization_id = $1 AND "teams"."deleted_at" IS NULL ORDER BY name ASC`)).
					WithArgs(orgId).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}))
			},
			expectedTeams: []models.Team{},
		},
		{
			name:  "database error",
			orgId: orgId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE organization_id = $1`)).
					WithArgs(orgId).
					WillReturnError(assert.AnError)
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

			teams, err := store.GetTeams(context.Background(), tt.orgId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expectedTeams), len(teams))
			for i := range teams {
				assert.Equal(t, tt.expectedTeams[i].TeamID, teams[i].TeamID)
				assert.Equal(t, tt.expectedTeams[i].OrganizationID, teams[i].OrganizationID)
			}
		})
	}
}

func TestGetTeam(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teamId := uuid.New()

	tests := []struct {
		name      string
		orgId     uuid.UUID
		teamId    uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantTeam  *models.Team
		wantErr   bool
	}{
		{
			name:   "success",
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"team_id", "organization_id"}).
					AddRow(teamId, orgId)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE (team_id = $1 AND organization_id = $2) AND "teams"."deleted_at" IS NULL ORDER BY "teams"."team_id" LIMIT $3`)).
					WithArgs(teamId, orgId, 1).
					WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE "team_memberships"."team_id" = $1 AND "team_memberships"."deleted_at" IS NULL`)).
					WithArgs(teamId).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "team_id", "user_id"}).
						AddRow(uuid.New(), teamId, uuid.New()))
			},
			wantTeam: &models.Team{TeamID: teamId, OrganizationID: orgId},
		},
		{
			name:   "team not found",
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE (team_id = $1 AND organization_id = $2) AND "teams"."deleted_at" IS NULL ORDER BY "teams"."team_id" LIMIT $3`)).
					WithArgs(teamId, orgId, 1).
					WillReturnError(assert.AnError)
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

			team, err := store.GetTeam(context.Background(), tt.orgId, tt.teamId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTeam.TeamID, team.TeamID)
			assert.Equal(t, tt.wantTeam.OrganizationID, team.OrganizationID)
		})
	}
}

func TestGetTeamByName(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teamId := uuid.New()
	teamName := "test-team"

	tests := []struct {
		name      string
		orgId     uuid.UUID
		teamName  string
		mockSetup func(sqlmock.Sqlmock)
		wantTeams []models.Team
		wantErr   bool
	}{
		{
			name:     "success",
			orgId:    orgId,
			teamName: teamName,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE (organization_id = $1 AND name = $2) AND "teams"."deleted_at" IS NULL LIMIT $3`)).
					WithArgs(orgId, teamName, 1).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id", "name"}).
						AddRow(teamId, orgId, teamName))
			},
			wantTeams: []models.Team{
				{
					TeamID:         teamId,
					OrganizationID: orgId,
					Name:           teamName,
				},
			},
		},
		{
			name:     "team not found",
			orgId:    orgId,
			teamName: "non-existent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE (organization_id = $1 AND name = $2) AND "teams"."deleted_at" IS NULL LIMIT $3`)).
					WithArgs(orgId, "non-existent", 1).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id", "name"}))
			},
			wantTeams: []models.Team{},
		},
		{
			name:     "database error",
			orgId:    orgId,
			teamName: teamName,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE (organization_id = $1 AND name = $2) AND "teams"."deleted_at" IS NULL LIMIT $3`)).
					WithArgs(orgId, teamName, 1).
					WillReturnError(assert.AnError)
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

			teams, err := store.GetTeamByName(context.Background(), tt.orgId, tt.teamName)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTeams, teams)
		})
	}
}

func TestCreateOrganizationTeam(t *testing.T) {
	t.Parallel()

	teamId := uuid.New()
	orgId := uuid.New()
	userId := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		team      models.Team
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			team: models.Team{OrganizationID: orgId, Name: "test", Description: "test", CreatedBy: userId},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teams" ("organization_id","name","description","deleted_at","created_by") VALUES ($1,$2,$3,$4,$5) RETURNING "team_id","created_at","updated_at","metadata"`)).
					WithArgs(orgId, "test", "test", sqlmock.AnyArg(), userId).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "created_at", "updated_at", "metadata"}).
						AddRow(teamId, now, now, json.RawMessage(`{}`)))
				mock.ExpectCommit()
			},
		},
		{
			name: "no user ID in context",
			team: models.Team{OrganizationID: orgId, Name: "test", Description: "test"},
			mockSetup: func(mock sqlmock.Sqlmock) {
			},
			wantErr: true,
		},
		{
			name: "database error",
			team: models.Team{OrganizationID: orgId, Name: "test", Description: "test", CreatedBy: userId},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teams" ("organization_id","name","description","deleted_at","created_by") VALUES ($1,$2,$3,$4,$5) RETURNING "team_id","created_at","updated_at","metadata"`)).
					WithArgs(orgId, "test", "test", sqlmock.AnyArg(), userId).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			team, err := store.CreateOrganizationTeam(ctx, orgId, tt.team)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, team)
			assert.Equal(t, teamId, team.TeamID)
			assert.Equal(t, orgId, team.OrganizationID)
			assert.Equal(t, "test", team.Name)
			assert.Equal(t, "test", team.Description)
			assert.Equal(t, now, team.CreatedAt)
			assert.Equal(t, now, team.UpdatedAt)
			assert.Nil(t, team.DeletedAt)
		})
	}
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teamId := uuid.New()
	userId := uuid.New()
	now := time.Now()
	tests := []struct {
		name      string
		team      models.Team
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			team: models.Team{TeamID: teamId, OrganizationID: orgId, Name: "test", Description: "test"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "teams" SET "organization_id"=$1,"name"=$2,"description"=$3,"updated_at"=$4 WHERE "teams"."deleted_at" IS NULL AND "team_id" = $5 RETURNING *`)).
					WithArgs(orgId, "test", "test", sqlmock.AnyArg(), teamId).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id", "name", "description", "created_at", "updated_at", "deleted_at"}).
						AddRow(teamId, orgId, "test", "test", now, now, nil))
				mock.ExpectCommit()
			},
		},
		{
			name: "database error",
			team: models.Team{TeamID: teamId},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "teams" SET "organization_id"=$1,"name"=$2,"description"=$3,"updated_at"=$4 WHERE "teams"."deleted_at" IS NULL AND "team_id" = $5 RETURNING *`)).
					WithArgs(orgId, "test", "test", sqlmock.AnyArg(), teamId).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			team, err := store.UpdateTeam(ctx, orgId, tt.team)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, team)
		})
	}
}

func TestDeleteTeam(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teamId := uuid.New()
	userId := uuid.New()
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teams" SET "deleted_at"=$1 WHERE team_id = $2 AND organization_id = $3 AND "teams"."team_id" = $4 AND "teams"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), teamId, orgId, teamId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teams" SET "deleted_at"=$1 WHERE team_id = $2 AND organization_id = $3 AND "teams"."team_id" = $4 AND "teams"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), teamId, orgId, teamId).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			err := store.DeleteTeam(ctx, orgId, teamId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestAppStore_GetTeamMemberships(t *testing.T) {
	teamId := uuid.New()
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE team_id = $1 AND "team_memberships"."deleted_at" IS NULL`)).
					WithArgs(teamId).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "team_id", "user_id"}).
						AddRow(uuid.New(), teamId, userId))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams"`)).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}).
						AddRow(teamId, orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).
						AddRow(userId))
			},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE team_id = $1 AND "team_memberships"."deleted_at" IS NULL`)).
					WithArgs(teamId).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			memberships, err := store.GetTeamMemberships(ctx, teamId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, memberships)
		})
	}
}

func TestAppStore_GetTeamMembershipById(t *testing.T) {
	membershipId := uuid.New()
	userId := uuid.New()
	teamId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE team_membership_id = $1 AND "team_memberships"."deleted_at" IS NULL ORDER BY "team_memberships"."team_membership_id" LIMIT $2`)).
					WithArgs(membershipId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "team_id", "user_id"}).
						AddRow(membershipId, teamId, userId))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams"`)).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}).
						AddRow(teamId, orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).
						AddRow(userId))
			},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE team_membership_id = $1 AND "team_memberships"."deleted_at" IS NULL ORDER BY "team_memberships"."team_membership_id" LIMIT $2`)).
					WithArgs(membershipId, 1).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			membership, err := store.GetTeamMembershipById(ctx, membershipId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, membership)
		})
	}
}

func TestAppStore_GetTeamMembershipByUserIdTeamId(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE (user_id = $1 AND team_id = $2) AND "team_memberships"."deleted_at" IS NULL ORDER BY "team_memberships"."team_membership_id" LIMIT $3`)).
					WithArgs(userId, teamId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "team_id", "user_id"}).
						AddRow(uuid.New(), teamId, userId))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams"`)).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}).
						AddRow(teamId, orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).
						AddRow(userId))
			},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "team_memberships" WHERE (user_id = $1 AND team_id = $2) AND "team_memberships"."deleted_at" IS NULL ORDER BY "team_memberships"."team_membership_id" LIMIT $3`)).
					WithArgs(userId, teamId, 1).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			membership, err := store.GetTeamMembershipByUserIdTeamId(ctx, userId, teamId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, membership)
		})
	}
}

func TestAppStore_CreateTeamMembership(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	orgId := uuid.New()
	now := time.Now()

	currentUserId := uuid.New()

	tests := []struct {
		name      string
		teamId    uuid.UUID
		userId    uuid.UUID
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:   "success",
			teamId: teamId,
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, models.ResourceTypeOrganization, currentUserId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}).AddRow(teamId, orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "team_memberships" ("team_id","user_id","deleted_at","created_by") VALUES ($1,$2,$3,$4) RETURNING "team_membership_id","created_at","updated_at"`)).
					WithArgs(teamId, userId, nil, currentUserId).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "created_at", "updated_at"}).AddRow(uuid.New(), now, now))
				mock.ExpectCommit()
			},
		},
		{
			name:   "no user ID in context",
			teamId: teamId,
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
			},
			wantErr: true,
		},
		{
			name:   "database error",
			teamId: teamId,
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, models.ResourceTypeOrganization, currentUserId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}).AddRow(teamId, orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "team_memberships" ("team_id","user_id","deleted_at","created_by") VALUES ($1,$2,$3,$4) RETURNING "team_membership_id","created_at","updated_at"`)).
					WithArgs(teamId, userId, nil, currentUserId).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", currentUserId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			membership, err := store.CreateTeamMembership(ctx, tt.teamId, tt.userId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, membership)
		})
	}
}

func TestAppStore_DeleteTeamMembershipByUserIdTeamId(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	orgId := uuid.New()
	teamMembershipId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(teamId, models.ResourceTypeOrganization, userId, models.PrivilegeOrganizationSystemAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"team_id", "organization_id"}).AddRow(teamId, orgId))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "team_memberships" SET "deleted_at"=$1 WHERE team_id = $2 AND "team_memberships"."team_membership_id" = $3 AND "team_memberships"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), teamId, teamMembershipId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teams" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE teams.team_id = $1 AND frap.resource_type = $2 AND frap.resource_id = teams.organization_id AND frap.user_id = $3 AND frap.deleted_at IS NULL AND frap.privilege = $4 ) AND "teams"."deleted_at" IS NULL`)).
					WithArgs(userId, teamId).
					WillReturnRows(sqlmock.NewRows([]string{"team_membership_id", "team_id", "user_id"}).
						AddRow(uuid.New(), teamId, userId))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "team_memberships" SET "deleted_at"=$1 WHERE team_id = $2 AND "team_memberships"."team_membership_id" = $3 AND "team_memberships"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), teamId, teamMembershipId).
					WillReturnError(assert.AnError)
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

			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			err := store.DeleteTeamMembership(ctx, teamId, teamMembershipId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
