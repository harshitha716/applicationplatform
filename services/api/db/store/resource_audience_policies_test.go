package store

import (
	"context"
	"fmt"
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

func TestGetResourcePolicies(t *testing.T) {
	t.Parallel()

	now := time.Now()
	resourceID := uuid.New()
	policyID1 := uuid.New()
	policyID2 := uuid.New()
	audienceID1 := uuid.New()
	audienceID2 := uuid.New()

	tests := []struct {
		name        string
		resourceID  uuid.UUID
		resourceTyp string
		mockSetup   func(sqlmock.Sqlmock)
		want        []models.ResourceAudiencePolicy
		wantErr     bool
	}{
		{
			name:        "get organization policies - success",
			resourceID:  resourceID,
			resourceTyp: "organization",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"resource_audience_policy_id",
					"resource_audience_type",
					"resource_audience_id",
					"privilege",
					"resource_type",
					"resource_id",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					policyID1,
					"user",
					audienceID1,
					models.PrivilegeOrganizationMember,
					"organization",
					resourceID,
					now,
					now,
					nil,
				).AddRow(
					policyID2,
					"group",
					audienceID2,
					models.PrivilegeOrganizationSystemAdmin,
					"organization",
					resourceID,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "resource_audience_policies" WHERE resource_id = \$1 AND resource_type = \$2`).
					WithArgs(resourceID, models.ResourceTypeOrganization).
					WillReturnRows(rows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" IN ($1,$2)`)).
					WithArgs(audienceID1, audienceID2).WillReturnRows(
					sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(audienceID1, "audience1@user.com", "audience1").
						AddRow(audienceID2, "audience2@user.com", "audience2"),
				)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE "flattened_resource_audience_policies"."resource_audience_policy_id" IN ($1,$2)`)).
					WithArgs(policyID1, policyID2).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}))

			},
			want: []models.ResourceAudiencePolicy{
				{
					ID:                   policyID1,
					ResourceAudienceType: "user",
					ResourceAudienceID:   audienceID1,
					Privilege:            models.PrivilegeOrganizationMember,
					ResourceType:         "organization",
					ResourceID:           resourceID,
					User: &models.User{
						ID:    audienceID1,
						Email: "audience1@user.com",
						Name:  "audience1",
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:                   policyID2,
					ResourceAudienceType: "group",
					ResourceAudienceID:   audienceID2,
					User: &models.User{
						ID:    audienceID2,
						Email: "audience2@user.com",
						Name:  "audience2",
					},
					Privilege:    models.PrivilegeOrganizationSystemAdmin,
					ResourceType: "organization",
					ResourceID:   resourceID,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		},
		{
			name:        "get dataset policies - success",
			resourceID:  resourceID,
			resourceTyp: "dataset",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"resource_audience_policy_id",
					"resource_audience_type",
					"resource_audience_id",
					"privilege",
					"resource_type",
					"resource_id",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					policyID1,
					"user",
					audienceID1,
					models.PrivilegeDatasetViewer,
					"dataset",
					resourceID,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "resource_audience_policies" WHERE resource_id = \$1 AND resource_type = \$2`).
					WithArgs(resourceID, "dataset").
					WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(audienceID1).WillReturnRows(
					sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(audienceID1, "audience1@user.com", "audience1"),
				)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE "flattened_resource_audience_policies"."resource_audience_policy_id" = $1`)).
					WithArgs(policyID1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}))
			},
			want: []models.ResourceAudiencePolicy{
				{
					ID:                   policyID1,
					ResourceAudienceType: "user",
					ResourceAudienceID:   audienceID1,
					Privilege:            models.PrivilegeDatasetViewer,
					User: &models.User{
						ID:    audienceID1,
						Email: "audience1@user.com",
						Name:  "audience1",
					},
					ResourceType: "dataset",
					ResourceID:   resourceID,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		},
		{
			name:        "get page policies - success",
			resourceID:  resourceID,
			resourceTyp: "page",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"resource_audience_policy_id",
					"resource_audience_type",
					"resource_audience_id",
					"privilege",
					"resource_type",
					"resource_id",
					"created_at",
					"updated_at",
					"deleted_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "resource_audience_policies" WHERE resource_id = \$1 AND resource_type = \$2`).
					WithArgs(resourceID, "page").
					WillReturnRows(rows)
			},
			want: []models.ResourceAudiencePolicy{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			var got []models.ResourceAudiencePolicy
			var err error

			switch tt.resourceTyp {
			case "organization":
				got, err = store.GetOrganizationPolicies(context.Background(), tt.resourceID)
			case "dataset":
				got, err = store.GetDatasetPolicies(context.Background(), tt.resourceID)
			case "page":
				got, err = store.GetPagesPolicies(context.Background(), tt.resourceID)
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))

			for i := range got {
				assert.Equal(t, tt.want[i].ID, got[i].ID)
				assert.Equal(t, tt.want[i].ResourceAudienceType, got[i].ResourceAudienceType)
				assert.Equal(t, tt.want[i].ResourceAudienceID, got[i].ResourceAudienceID)
				assert.Equal(t, tt.want[i].Privilege, got[i].Privilege)
				assert.Equal(t, tt.want[i].ResourceType, got[i].ResourceType)
				assert.Equal(t, tt.want[i].ResourceID, got[i].ResourceID)
				assert.Equal(t, tt.want[i].CreatedAt.Unix(), got[i].CreatedAt.Unix())
				assert.Equal(t, tt.want[i].UpdatedAt.Unix(), got[i].UpdatedAt.Unix())
				if tt.want[i].DeletedAt == nil {
					assert.Nil(t, got[i].DeletedAt)
				} else {
					assert.Equal(t, tt.want[i].DeletedAt.Unix(), got[i].DeletedAt.Unix())
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateResourceAudiencePolicy(t *testing.T) {
	t.Parallel()

	policyId := uuid.New()

	resourceID := uuid.New()
	audienceID := uuid.New()

	tests := []struct {
		name         string
		resourceID   uuid.UUID
		resourceType models.ResourceType
		audienceID   uuid.UUID
		privilege    models.ResourcePrivilege
		audienceType models.AudienceType
		mockSetup    func(sqlmock.Sqlmock)
		want         *models.ResourceAudiencePolicy
		wantErr      bool
	}{
		{
			name:         "create organization policy - success",
			resourceID:   resourceID,
			audienceID:   audienceID,
			privilege:    models.PrivilegeOrganizationMember,
			audienceType: models.AudienceTypeUser,
			resourceType: models.ResourceTypeOrganization,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resource_audience_policies" ("resource_audience_type","resource_audience_id","privilege","resource_type","resource_id","created_at","updated_at","deleted_at","metadata") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "resource_audience_policy_id"`)).
					WithArgs("user", audienceID, "member", "organization", resourceID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}).AddRow(policyId))
				mock.ExpectCommit()
			},
			want: &models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: models.AudienceTypeUser,
				ResourceAudienceID:   audienceID,
				Privilege:            models.PrivilegeOrganizationMember,
				ResourceType:         models.ResourceTypeOrganization,
				ResourceID:           resourceID,
			},
			wantErr: false,
		},
		{
			name:         "create organization policy - db error",
			resourceID:   resourceID,
			audienceID:   audienceID,
			privilege:    models.PrivilegeDatasetViewer,
			audienceType: models.AudienceTypeUser,
			resourceType: models.ResourceTypeDataset,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resource_audience_policies" ("resource_audience_type","resource_audience_id","privilege","resource_type","resource_id","created_at","updated_at","deleted_at","metadata") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "resource_audience_policy_id"`)).
					WithArgs("user", audienceID, models.PrivilegeDatasetViewer, "dataset", resourceID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectCommit()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			policy, err := store.CreateOrganizationPolicy(context.Background(), tt.resourceID, tt.audienceType, tt.audienceID, tt.privilege)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, policyId)
				assert.Equal(t, tt.want.ResourceAudienceType, policy.ResourceAudienceType)
				assert.Equal(t, tt.want.ResourceAudienceID, policy.ResourceAudienceID)
				assert.Equal(t, tt.want.Privilege, policy.Privilege)
				assert.Equal(t, tt.want.ResourceType, policy.ResourceType)
				assert.Equal(t, tt.want.ResourceID, policy.ResourceID)
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}

func TestCreateDatasetPolicy(t *testing.T) {
	t.Parallel()

	policyId := uuid.New()

	resourceID := uuid.New()

	audienceId := uuid.New()

	tests := []struct {
		name         string
		resourceID   uuid.UUID
		audienceID   uuid.UUID
		privilege    models.ResourcePrivilege
		audienceType models.AudienceType
		mockSetup    func(sqlmock.Sqlmock)
		want         *models.ResourceAudiencePolicy
		wantErr      bool
	}{
		{
			name:         "create dataset policy - success",
			resourceID:   resourceID,
			privilege:    models.PrivilegeDatasetViewer,
			audienceType: models.AudienceTypeUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resource_audience_policies" ("resource_audience_type","resource_audience_id","privilege","resource_type","resource_id","created_at","updated_at","deleted_at","metadata") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "resource_audience_policy_id"`)).
					WithArgs("user", audienceId, models.PrivilegeDatasetViewer, "dataset", resourceID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}).AddRow(policyId))
				mock.ExpectCommit()
			},
			want: &models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: models.AudienceTypeUser,
				ResourceAudienceID:   audienceId,
				Privilege:            models.PrivilegeDatasetViewer,
				ResourceType:         models.ResourceTypeDataset,
				ResourceID:           resourceID,
			},
			wantErr: false,
		},
		{
			name:         "create dataset policy - db error",
			resourceID:   resourceID,
			privilege:    models.PrivilegeDatasetViewer,
			audienceType: models.AudienceTypeUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resource_audience_policies" ("resource_audience_type","resource_audience_id","privilege","resource_type","resource_id","created_at","updated_at","deleted_at","metadata") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "resource_audience_policy_id"`)).
					WithArgs("user", uuid.New(), models.PrivilegeDatasetViewer, "dataset", resourceID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			policy, err := store.CreateDatasetPolicy(context.Background(), tt.resourceID, tt.audienceType, audienceId, tt.privilege)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, policy.ID)
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}

func TestCreatePagePolicy(t *testing.T) {
	t.Parallel()

	policyId := uuid.New()

	resourceID := uuid.New()

	audienceId := uuid.New()

	tests := []struct {
		name         string
		resourceID   uuid.UUID
		audienceID   uuid.UUID
		privilege    models.ResourcePrivilege
		audienceType models.AudienceType
		mockSetup    func(sqlmock.Sqlmock)
		want         *models.ResourceAudiencePolicy
		wantErr      bool
	}{
		{
			name:         "create page policy - success",
			resourceID:   resourceID,
			privilege:    models.PrivilegePageRead,
			audienceType: models.AudienceTypeUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resource_audience_policies" ("resource_audience_type","resource_audience_id","privilege","resource_type","resource_id","created_at","updated_at","deleted_at","metadata") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "resource_audience_policy_id"`)).
					WithArgs("user", audienceId, models.PrivilegePageRead, "page", resourceID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}).AddRow(policyId))
				mock.ExpectCommit()
			},
			want: &models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: models.AudienceTypeUser,
				ResourceAudienceID:   audienceId,
				Privilege:            models.PrivilegePageRead,
				ResourceType:         models.ResourceTypePage,
				ResourceID:           resourceID,
			},
			wantErr: false,
		},
		{
			name:         "create page policy - db error",
			resourceID:   resourceID,
			privilege:    models.PrivilegePageRead,
			audienceType: models.AudienceTypeUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resource_audience_policies" ("resource_audience_type","resource_audience_id","privilege","resource_type","resource_id","created_at","updated_at","deleted_at","metadata") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "resource_audience_policy_id"`)).
					WithArgs("user", uuid.New(), models.PrivilegePageRead, models.ResourceTypePage, resourceID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			policy, err := store.CreatePagePolicy(context.Background(), tt.resourceID, tt.audienceType, audienceId, tt.privilege)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, policy.ID)
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}

func TestUpdateOrganizationPolicy(t *testing.T) {
	t.Parallel()

	policyId := uuid.New()

	currentUserId := uuid.New()

	resourceID := uuid.New()
	audienceID := uuid.New()

	tests := []struct {
		name         string
		resourceID   uuid.UUID
		audienceID   uuid.UUID
		privilege    models.ResourcePrivilege
		audienceType models.AudienceType
		mockSetup    func(sqlmock.Sqlmock)
		want         *models.ResourceAudiencePolicy
		wantErr      bool
	}{
		{
			name:         "update organization policy - success",
			resourceID:   resourceID,
			audienceID:   audienceID,
			privilege:    models.PrivilegeOrganizationMember,
			audienceType: models.AudienceTypeUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
					WithArgs(resourceID, models.ResourceTypeOrganization, currentUserId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(resourceID))
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "resource_audience_policies" SET "privilege"=$1,"updated_at"=$2 WHERE resource_type = $3 AND resource_id = $4 AND resource_audience_id = $5`)).
					WithArgs("member", sqlmock.AnyArg(), models.ResourceTypeOrganization, resourceID, audienceID).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id", "resource_audience_type", "resource_audience_id", "resource_type", "resource_id"}).
						AddRow(policyId, models.AudienceTypeUser, audienceID, models.ResourceTypeOrganization, resourceID))
				mock.ExpectCommit()
			},
			want: &models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: models.AudienceTypeUser,
				ResourceAudienceID:   audienceID,
				Privilege:            models.PrivilegeOrganizationMember,
				ResourceType:         models.ResourceTypeOrganization,
				ResourceID:           resourceID,
			},
		},
		{
			name:         "update organization policy - db error",
			resourceID:   resourceID,
			audienceID:   audienceID,
			privilege:    models.PrivilegeDatasetViewer,
			audienceType: models.AudienceTypeUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
					WithArgs(resourceID, models.ResourceTypeOrganization, currentUserId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(resourceID))
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "resource_audience_policies" SET "privilege"=$1,"updated_at"=$2 WHERE resource_type = $3 AND resource_id = $4 AND resource_audience_id = $5`)).
					WithArgs("member", sqlmock.AnyArg(), models.ResourceTypeOrganization, resourceID, audienceID).
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "user", currentUserId, []uuid.UUID{uuid.New()})

			// Execute
			policy, err := store.UpdateOrganizationPolicy(ctx, tt.resourceID, tt.audienceID, tt.privilege)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, policyId)
				assert.Equal(t, tt.want.ResourceAudienceType, policy.ResourceAudienceType)
				assert.Equal(t, tt.want.ResourceAudienceID, policy.ResourceAudienceID)
				assert.Equal(t, tt.want.Privilege, policy.Privilege)
				assert.Equal(t, tt.want.ResourceType, policy.ResourceType)
				assert.Equal(t, tt.want.ResourceID, policy.ResourceID)
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}

func TestGetOrganizationPolicyByUser(t *testing.T) {
	t.Parallel()

	now := time.Now()
	resourceID := uuid.New()
	policyID := uuid.New()
	audienceID := uuid.New()

	tests := []struct {
		name       string
		mockSetup  func(sqlmock.Sqlmock)
		want       *models.ResourceAudiencePolicy
		wantErr    bool
		audienceID uuid.UUID
	}{
		{
			name: "get organization policy by user - success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"resource_audience_policy_id",
					"resource_audience_type",
					"resource_audience_id",
					"privilege",
					"resource_type",
					"resource_id",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					policyID,
					"user",
					audienceID,
					models.PrivilegeOrganizationMember,
					"organization",
					resourceID,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_audience_policies" WHERE resource_id = $1 AND resource_type = $2 AND resource_audience_id = $3 ORDER BY "resource_audience_policies"."resource_audience_policy_id" LIMIT $4`)).
					WithArgs(resourceID, "organization", audienceID, 1).
					WillReturnRows(rows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(audienceID).WillReturnRows(
					sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(audienceID, "audience1@example.com", "audience1"),
				)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE "flattened_resource_audience_policies"."resource_audience_policy_id" = $1`)).
					WithArgs(policyID).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}))
			},
			want: &models.ResourceAudiencePolicy{
				ID:                   policyID,
				ResourceAudienceType: "user",
				ResourceAudienceID:   audienceID,
				Privilege:            models.PrivilegeOrganizationMember,
				ResourceType:         "organization",
				ResourceID:           resourceID,
			},
			wantErr:    false,
			audienceID: audienceID,
		},

		{
			name: "get organization policy by user - not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_audience_policies" WHERE resource_id = $1 AND resource_type = $2 AND resource_audience_id = $3 ORDER BY "resource_audience_policies"."resource_audience_policy_id" LIMIT $4`)).
					WithArgs(resourceID, "organization", audienceID, 1).
					WillReturnError(fmt.Errorf("record not found"))
			},
			want:       nil,
			wantErr:    true,
			audienceID: audienceID,
		},
		{
			name: "get organization policy by user - db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_audience_policies" WHERE resource_id = $1 AND resource_type = $2 AND resource_audience_id = $3 ORDER BY "resource_audience_policies"."resource_audience_policy_id" LIMIT $4`)).
					WithArgs(resourceID, "organization", audienceID, 1).
					WillReturnError(fmt.Errorf("db error"))
			},
			want:       nil,
			wantErr:    true,
			audienceID: audienceID,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			policy, err := store.GetOrganizationPolicyByUser(context.Background(), resourceID, tt.audienceID)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, policy.ID)
				assert.Equal(t, tt.want.ResourceAudienceType, policy.ResourceAudienceType)
				assert.Equal(t, tt.want.ResourceAudienceID, policy.ResourceAudienceID)
				assert.Equal(t, tt.want.Privilege, policy.Privilege)
				assert.Equal(t, tt.want.ResourceType, policy.ResourceType)
				assert.Equal(t, tt.want.ResourceID, policy.ResourceID)
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}

func TestGetPolicyByEmail(t *testing.T) {

	t.Parallel()

	types := []models.ResourceType{models.ResourceTypeOrganization, models.ResourceTypeDataset, models.ResourceTypePage}

	for _, resourceType := range types {

		now := time.Now()
		resourceID := uuid.New()
		policyID := uuid.New()
		audienceID := uuid.New()

		var privilege models.ResourcePrivilege
		switch resourceType {
		case models.ResourceTypeOrganization:
			privilege = models.PrivilegeOrganizationMember
		case models.ResourceTypeDataset:
			privilege = models.PrivilegeDatasetViewer
		case models.ResourceTypePage:
			privilege = models.PrivilegePageRead
		}
		tests := []struct {
			name       string
			resourceID uuid.UUID
			email      string
			mockSetup  func(sqlmock.Sqlmock)
			want       []models.ResourceAudiencePolicy
			wantErr    bool
		}{
			{
				name:       fmt.Sprintf("get %s policy by email - success", resourceType),
				resourceID: resourceID,
				email:      "test@example.com",
				mockSetup: func(mock sqlmock.Sqlmock) {

					rows := sqlmock.NewRows([]string{
						"resource_audience_policy_id",
						"resource_audience_type",
						"resource_audience_id",
						"privilege",
						"resource_type",
						"resource_id",
						"created_at",
						"updated_at",
						"deleted_at",
					}).AddRow(
						policyID,
						"user",
						audienceID,
						privilege,
						resourceType,
						resourceID,
						now,
						now,
						nil,
					)

					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "resource_audience_policies"."resource_audience_policy_id","resource_audience_policies"."resource_audience_type","resource_audience_policies"."resource_audience_id","resource_audience_policies"."privilege","resource_audience_policies"."resource_type","resource_audience_policies"."resource_id","resource_audience_policies"."created_at","resource_audience_policies"."updated_at","resource_audience_policies"."deleted_at","resource_audience_policies"."metadata" FROM "resource_audience_policies" JOIN users_with_traits ON resource_audience_id = users_with_traits.user_id AND users_with_traits.email ILIKE $1 WHERE resource_id = $2 AND resource_type = $3 AND resource_audience_type = $4`)).
						WithArgs("test@example.com", resourceID, resourceType, "user").
						WillReturnRows(rows)
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE "flattened_resource_audience_policies"."resource_audience_policy_id" = $1`)).
						WithArgs(policyID).
						WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}))
				},
				want: []models.ResourceAudiencePolicy{{
					ID:                   policyID,
					ResourceAudienceType: "user",
					ResourceAudienceID:   audienceID,
					Privilege:            privilege,
					ResourceType:         resourceType,
					ResourceID:           resourceID,
					CreatedAt:            now,
					UpdatedAt:            now,
				}},
				wantErr: false,
			},
			{
				name:       fmt.Sprintf("get %s policy by email - empty response", resourceType),
				resourceID: resourceID,
				email:      "notfound@example.com",
				mockSetup: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "resource_audience_policies"."resource_audience_policy_id","resource_audience_policies"."resource_audience_type","resource_audience_policies"."resource_audience_id","resource_audience_policies"."privilege","resource_audience_policies"."resource_type","resource_audience_policies"."resource_id","resource_audience_policies"."created_at","resource_audience_policies"."updated_at","resource_audience_policies"."deleted_at","resource_audience_policies"."metadata" FROM "resource_audience_policies" JOIN users_with_traits ON resource_audience_id = users_with_traits.user_id AND users_with_traits.email ILIKE $1 WHERE resource_id = $2 AND resource_type = $3 AND resource_audience_type = $4`)).
						WithArgs("notfound@example.com", resourceID, resourceType, "user").
						WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id", "resource_audience_type", "resource_audience_id", "privilege", "resource_type", "resource_id", "created_at", "updated_at", "deleted_at"}))
				},
				want:    []models.ResourceAudiencePolicy{},
				wantErr: false,
			},
			{
				name:       fmt.Sprintf("get %s policy by email - not found", resourceType),
				resourceID: resourceID,
				email:      "notfound@example.com",
				mockSetup: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "resource_audience_policies"."resource_audience_policy_id","resource_audience_policies"."resource_audience_type","resource_audience_policies"."resource_audience_id","resource_audience_policies"."privilege","resource_audience_policies"."resource_type","resource_audience_policies"."resource_id","resource_audience_policies"."created_at","resource_audience_policies"."updated_at","resource_audience_policies"."deleted_at","resource_audience_policies"."metadata" FROM "resource_audience_policies" JOIN users_with_traits ON resource_audience_id = users_with_traits.user_id AND users_with_traits.email ILIKE $1 WHERE resource_id = $2 AND resource_type = $3 AND resource_audience_type = $4`)).
						WithArgs("notfound@example.com", resourceID, resourceType, "user").
						WillReturnError(fmt.Errorf("record not found"))
				},
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Setup
				gormDB, mock := getMockDB(t)
				store := &appStore{
					client: &pgclient.PostgresClient{DB: gormDB},
				}
				tt.mockSetup(mock)

				var got []models.ResourceAudiencePolicy
				var err error
				// Execute
				switch resourceType {
				case models.ResourceTypeOrganization:
					got, err = store.GetOrganizationPoliciesByEmail(context.Background(), tt.resourceID, tt.email)
				case models.ResourceTypeDataset:
					got, err = store.GetDatasetPoliciesByEmail(context.Background(), tt.resourceID, tt.email)
				case models.ResourceTypePage:
					got, err = store.GetPagePoliciesByEmail(context.Background(), tt.resourceID, tt.email)
				}

				// Assert
				if tt.wantErr {
					assert.Error(t, err)
					assert.Nil(t, got)
					return
				}

				if len(tt.want) == 0 {
					assert.Equal(t, 0, len(got))
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.want[0].ID, got[0].ID)
				assert.Equal(t, tt.want[0].ResourceAudienceType, got[0].ResourceAudienceType)
				assert.Equal(t, tt.want[0].ResourceAudienceID, got[0].ResourceAudienceID)
				assert.Equal(t, tt.want[0].Privilege, got[0].Privilege)
				assert.Equal(t, tt.want[0].ResourceType, got[0].ResourceType)
				assert.Equal(t, tt.want[0].ResourceID, got[0].ResourceID)
				assert.Equal(t, tt.want[0].CreatedAt.Unix(), got[0].CreatedAt.Unix())
				assert.Equal(t, tt.want[0].UpdatedAt.Unix(), got[0].UpdatedAt.Unix())

				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}

	}

}

func TestDeleteResourcePolicy(t *testing.T) {
	t.Parallel()

	resourceID := uuid.New()
	audienceID := uuid.New()

	tests := []struct {
		name         string
		resourceID   uuid.UUID
		audienceID   uuid.UUID
		audienceType models.AudienceType
		resourceType models.ResourceType
		mockSetup    func(sqlmock.Sqlmock)
		wantErr      bool
	}{
		{
			name:         "delete resource policy - success",
			resourceID:   resourceID,
			audienceID:   audienceID,
			audienceType: models.AudienceTypeUser,
			resourceType: models.ResourceTypeOrganization,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
					WithArgs(resourceID, models.ResourceTypeOrganization, audienceID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(resourceID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND resource_audience_type = $3 AND resource_audience_id = $4`)).
					WithArgs("organization", resourceID, "user", audienceID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:         "delete resource policy - db error",
			resourceID:   resourceID,
			audienceID:   audienceID,
			audienceType: models.AudienceTypeUser,
			resourceType: models.ResourceTypeOrganization,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
					WithArgs(resourceID, models.ResourceTypeOrganization, audienceID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(resourceID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND resource_audience_type = $3 AND resource_audience_id = $4`)).
					WithArgs("organization", resourceID, "user", audienceID).
					WillReturnError(fmt.Errorf("db error"))
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "user", audienceID, []uuid.UUID{uuid.New()})

			// Execute
			err := store.deleteResourcePolicy(ctx, tt.resourceType, tt.resourceID, tt.audienceType, tt.audienceID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
