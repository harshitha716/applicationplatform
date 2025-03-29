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

func TestGetOrganizationById(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	ownerID := uuid.New()
	now := time.Now()
	description := "Test Description"

	tests := []struct {
		name      string
		orgID     string
		mockSetup func(sqlmock.Sqlmock)
		want      *models.Organization
		wantErr   bool
	}{
		{
			name:  "success",
			orgID: orgID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_id",
					"name",
					"description",
					"created_at",
					"updated_at",
					"deleted_at",
					"owner_id",
				}).AddRow(
					orgID,
					"Test Organization",
					&description,
					now,
					now,
					nil,
					ownerID,
				)

				mock.ExpectQuery(`SELECT \* FROM "organizations" WHERE organization_id = \$1 ORDER BY "organizations"."organization_id" LIMIT \$2`).
					WithArgs(orgID.String(), 1).
					WillReturnRows(rows)
			},
			want: &models.Organization{
				ID:          orgID,
				Name:        "Test Organization",
				Description: &description,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
				OwnerId:     ownerID,
			},
		},
		{
			name:  "not found",
			orgID: uuid.New().String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations" WHERE organization_id = \$1 ORDER BY "organizations"."organization_id" LIMIT \$2`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:  "invalid uuid",
			orgID: "invalid-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations" WHERE organization_id = \$1 ORDER BY "organizations"."organization_id" LIMIT \$2`).
					WithArgs("invalid-uuid", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:  "database error",
			orgID: orgID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations" WHERE organization_id = \$1 ORDER BY "organizations"."organization_id" LIMIT \$2`).
					WithArgs(orgID.String(), 1).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationById(context.Background(), tt.orgID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.OwnerId, got.OwnerId)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
			assert.Equal(t, tt.want.UpdatedAt.Unix(), got.UpdatedAt.Unix())
			if tt.want.DeletedAt == nil {
				assert.Nil(t, got.DeletedAt)
			} else {
				assert.Equal(t, tt.want.DeletedAt.Unix(), got.DeletedAt.Unix())
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetOrganizationsAll(t *testing.T) {
	t.Parallel()

	org1ID := uuid.New()
	org2ID := uuid.New()
	ownerID := uuid.New()
	now := time.Now()
	description := "Test Description"

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []models.Organization
		wantErr   bool
	}{
		{
			name: "success - multiple organizations",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_id",
					"name",
					"description",
					"created_at",
					"updated_at",
					"deleted_at",
					"owner_id",
				}).AddRow(
					org1ID,
					"Test Organization 1",
					&description,
					now,
					now,
					nil,
					ownerID,
				).AddRow(
					org2ID,
					"Test Organization 2",
					nil, // Testing nil description
					now,
					now,
					nil,
					ownerID,
				)

				mock.ExpectQuery(`SELECT \* FROM "organizations"`).
					WillReturnRows(rows)
			},
			want: []models.Organization{
				{
					ID:          org1ID,
					Name:        "Test Organization 1",
					Description: &description,
					CreatedAt:   now,
					UpdatedAt:   now,
					DeletedAt:   nil,
					OwnerId:     ownerID,
				},
				{
					ID:          org2ID,
					Name:        "Test Organization 2",
					Description: nil,
					CreatedAt:   now,
					UpdatedAt:   now,
					DeletedAt:   nil,
					OwnerId:     ownerID,
				},
			},
		},
		{
			name: "success - empty result",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations"`).
					WillReturnRows(sqlmock.NewRows([]string{
						"organization_id",
						"name",
						"description",
						"created_at",
						"updated_at",
						"deleted_at",
						"owner_id",
					}))
			},
			want: []models.Organization{},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations"`).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantErr: true,
		},
		{
			name: "context canceled",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations"`).
					WillReturnError(context.Canceled)
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

			// Execute
			got, err := store.GetOrganizationsAll(context.Background())

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
				assert.Equal(t, tt.want[i].Name, got[i].Name)
				assert.Equal(t, tt.want[i].Description, got[i].Description)
				assert.Equal(t, tt.want[i].OwnerId, got[i].OwnerId)
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

func TestGetOrganizationInvitationsAll(t *testing.T) {
	t.Parallel()

	invitation1ID := uuid.New()
	invitation2ID := uuid.New()
	organizationID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []models.OrganizationInvitation
		wantErr   bool
	}{
		{
			name: "success - multiple invitations",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_invitation_id",
					"organization_id",
					"email",
					"created_at",
					"updated_at",
				}).AddRow(
					invitation1ID,
					organizationID,
					"test1@example.com",
					now,
					now,
				).AddRow(
					invitation2ID,
					organizationID,
					"test2@example.com",
					now,
					now,
				)

				mock.ExpectQuery(`SELECT \* FROM "organization_invitations"`).
					WillReturnRows(rows)
			},
			want: []models.OrganizationInvitation{
				{
					OrganizationInvitationID: invitation1ID,
					OrganizationID:           organizationID,
					TargetEmail:              "test1@example.com",
					CreatedAt:                now,
					UpdatedAt:                now,
				},
				{
					OrganizationInvitationID: invitation2ID,
					OrganizationID:           organizationID,
					TargetEmail:              "test2@example.com",
					CreatedAt:                now,
					UpdatedAt:                now,
				},
			},
		},
		{
			name: "success - empty result",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_invitations"`).
					WillReturnRows(sqlmock.NewRows([]string{
						"organization_invitation_id",
						"organization_id",
						"email",
						"created_at",
						"updated_at",
					}))
			},
			want: []models.OrganizationInvitation{},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_invitations"`).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationInvitationsAll(context.Background())

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetOrganizationInvitationsByOrganizationId(t *testing.T) {
	t.Parallel()

	invitationID := uuid.New()
	organizationID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		orgID     uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      []models.OrganizationInvitation
		wantErr   bool
	}{
		{
			name:  "success",
			orgID: organizationID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_invitation_id",
					"organization_id",
					"email",
					"created_at",
					"updated_at",
				}).AddRow(
					invitationID,
					organizationID,
					"test@example.com",
					now,
					now,
				)

				mock.ExpectQuery(`SELECT \* FROM "organization_invitations" WHERE organization_id = \$1`).
					WithArgs(organizationID.String()).
					WillReturnRows(rows)
			},
			want: []models.OrganizationInvitation{
				{
					OrganizationInvitationID: invitationID,
					OrganizationID:           organizationID,
					TargetEmail:              "test@example.com",
					CreatedAt:                now,
					UpdatedAt:                now,
				},
			},
		},
		{
			name:  "not found",
			orgID: organizationID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_invitations" WHERE organization_id = \$1`).
					WithArgs(organizationID.String()).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			want: []models.OrganizationInvitation{},
		},
		{
			name:  "database error",
			orgID: organizationID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_invitations" WHERE organization_id = \$1`).
					WithArgs(organizationID.String()).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationInvitationsByOrganizationId(context.Background(), tt.orgID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetOrganizationsByMemberId(t *testing.T) {
	t.Parallel()

	memberID := uuid.New()
	org1ID := uuid.New()
	org2ID := uuid.New()
	ownerID := uuid.New()
	now := time.Now()
	description := "Test Description"

	tests := []struct {
		name      string
		memberID  uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      []models.Organization
		wantErr   bool
	}{
		{
			name:     "success - multiple organizations",
			memberID: memberID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_id",
					"name",
					"description",
					"created_at",
					"updated_at",
					"deleted_at",
					"owner_id",
				}).AddRow(
					org1ID,
					"Test Organization 1",
					&description,
					now,
					now,
					nil,
					ownerID,
				).AddRow(
					org2ID,
					"Test Organization 2",
					nil,
					now,
					now,
					nil,
					ownerID,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations`)).
					WillReturnRows(rows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_audience_policies" WHERE "resource_audience_policies"."resource_id" IN ($1,$2) AND (resource_audience_policies.resource_audience_id = $3 AND resource_audience_policies.resource_audience_type = $4 AND resource_audience_policies.resource_type = $5)`)).
					WithArgs(org1ID, org2ID, memberID, models.AudienceTypeUser, models.ResourceTypeOrganization).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_id", "resource_id", "resource_type", "resource_audience_type", "privilege"}))
			},
			want: []models.Organization{
				{
					ID:          org1ID,
					Name:        "Test Organization 1",
					Description: &description,
					CreatedAt:   now,
					UpdatedAt:   now,
					DeletedAt:   nil,
					OwnerId:     ownerID,
				},
				{
					ID:          org2ID,
					Name:        "Test Organization 2",
					Description: nil,
					CreatedAt:   now,
					UpdatedAt:   now,
					DeletedAt:   nil,
					OwnerId:     ownerID,
					ResourceAudiencePolicies: []models.ResourceAudiencePolicy{
						{
							ResourceAudienceID:   memberID,
							ResourceID:           org2ID,
							ResourceType:         models.ResourceTypeOrganization,
							ResourceAudienceType: models.AudienceTypeUser,
							Privilege:            models.PrivilegeOrganizationSystemAdmin,
						},
					},
				},
			},
		},
		{
			name:     "success - no organizations",
			memberID: memberID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations`)).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			want: []models.Organization{},
		},
		{
			name:     "database error",
			memberID: memberID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations`)).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationsByMemberId(context.Background(), tt.memberID)

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
				assert.Equal(t, tt.want[i].Name, got[i].Name)
				assert.Equal(t, tt.want[i].Description, got[i].Description)
				assert.Equal(t, tt.want[i].OwnerId, got[i].OwnerId)
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

func TestGetOrganizationInvitationById(t *testing.T) {
	t.Parallel()

	invitationID := uuid.New()
	organizationID := uuid.New()
	invitedEmail := "test@example.com"
	now := time.Now()

	tests := []struct {
		name      string
		invID     uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      *models.OrganizationInvitation
		wantErr   bool
	}{
		{
			name:  "success",
			invID: invitationID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_invitation_id",
					"organization_id",
					"email",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					invitationID,
					organizationID,
					invitedEmail,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "organization_invitations" WHERE organization_invitation_id = \$1 ORDER BY "organization_invitations"."organization_invitation_id" LIMIT \$2`).
					WithArgs(invitationID.String(), 1).
					WillReturnRows(rows)
			},
			want: &models.OrganizationInvitation{
				OrganizationInvitationID: invitationID,
				OrganizationID:           organizationID,
				TargetEmail:              invitedEmail,
				CreatedAt:                now,
				UpdatedAt:                now,
				DeletedAt:                nil,
			},
		},
		{
			name:  "not found",
			invID: uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_invitations" WHERE organization_invitation_id = \$1 ORDER BY "organization_invitations"."organization_invitation_id" LIMIT \$2`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:  "database error",
			invID: invitationID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_invitations" WHERE organization_invitation_id = \$1 ORDER BY "organization_invitations"."organization_invitation_id" LIMIT \$2`).
					WithArgs(invitationID.String(), 1).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationInvitationById(context.Background(), tt.invID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.TargetEmail, got.TargetEmail)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
			assert.Equal(t, tt.want.UpdatedAt.Unix(), got.UpdatedAt.Unix())
			if tt.want.DeletedAt == nil {
				assert.Nil(t, got.DeletedAt)
			} else {
				assert.Equal(t, tt.want.DeletedAt.Unix(), got.DeletedAt.Unix())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateOrganizationInvitation(t *testing.T) {
	t.Parallel()

	invitationId := uuid.New()
	organizationID := uuid.New()
	invitedEmail := "test@example.com"
	invitedBy := uuid.New()
	tests := []struct {
		name      string
		ctxSetup  func() context.Context
		mockSetup func(sqlmock.Sqlmock)
		want      *models.OrganizationInvitation
		wantErr   bool
	}{
		{
			name: "success",
			ctxSetup: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", invitedBy, []uuid.UUID{})
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow("organization", organizationID, invitedBy, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_invitations" ("organization_id","email","privilege","invited_by","email_retry_count","email_sent_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "organization_invitation_id","created_at","updated_at"`)).
					WithArgs(organizationID, invitedEmail, models.PrivilegeOrganizationSystemAdmin, invitedBy, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationId))
				mock.ExpectCommit()
			},
			want: &models.OrganizationInvitation{
				OrganizationInvitationID: invitationId,
				OrganizationID:           organizationID,
				TargetEmail:              invitedEmail,
				InvitedBy:                invitedBy,
			},
		},
		{
			name: "missing user in context",
			ctxSetup: func() context.Context {
				return context.Background()
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
		},
		{
			name: "database error",
			ctxSetup: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", invitedBy, []uuid.UUID{})
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow("organization", organizationID, invitedBy, models.PrivilegeOrganizationSystemAdmin))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_invitations" ("organization_id","email","privilege","invited_by","email_retry_count","email_sent_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "organization_invitation_id","created_at","updated_at"`)).
					WithArgs(organizationID, invitedEmail, models.PrivilegeOrganizationSystemAdmin, invitedBy, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := tt.ctxSetup()
			tt.mockSetup(mock)

			// Execute
			got, err := store.CreateOrganizationInvitation(ctx, organizationID, invitedEmail, models.PrivilegeOrganizationSystemAdmin)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.TargetEmail, got.TargetEmail)
			assert.Equal(t, tt.want.InvitedBy, got.InvitedBy)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateOrganizationInvitationStatus(t *testing.T) {
	t.Parallel()

	invitationID := uuid.New()
	userID := uuid.New()
	newStatus := models.InvitationStatusAccepted

	tests := []struct {
		name      string
		ctxSetup  func() context.Context
		mockSetup func(sqlmock.Sqlmock)
		want      *models.OrganizationInvitationStatus
		wantErr   bool
	}{
		{
			name: "success",
			ctxSetup: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", userID, []uuid.UUID{})
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationID, userID).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationID))
				mock.ExpectQuery(`INSERT INTO \"organization_invitation_statuses\"`).
					WithArgs(invitationID, newStatus).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationID))
				mock.ExpectCommit()
			},
			want: &models.OrganizationInvitationStatus{
				OrganizationInvitationID: invitationID,
				Status:                   newStatus,
			},
		},
		{
			name: "missing user in context",
			ctxSetup: func() context.Context {
				return context.Background()
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
		},
		{
			name: "database error",
			ctxSetup: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", userID, []uuid.UUID{})

			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationID, userID).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationID))
				mock.ExpectQuery(`INSERT INTO \"organization_invitation_statuses\"`).
					WithArgs(invitationID, newStatus).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := tt.ctxSetup()
			tt.mockSetup(mock)

			// Execute
			got, err := store.UpdateOrganizationInvitationStatus(ctx, invitationID, newStatus)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.OrganizationInvitationID, got.OrganizationInvitationID)
			assert.Equal(t, tt.want.Status, got.Status)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteOrganizationInvitation(t *testing.T) {
	t.Parallel()

	invitationID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		ctxSetup  func() context.Context
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			ctxSetup: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", userID, []uuid.UUID{})
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`
	 SELECT "organization_invitations"."organization_invitation_id","organization_invitations"."organization_id","organization_invitations"."email","organization_invitations"."privilege","organization_invitations"."created_at","organization_invitations"."invited_by","organization_invitations"."email_retry_count","organization_invitations"."email_sent_at","organization_invitations"."updated_at","organization_invitations"."deleted_at" FROM "organization_invitations" JOIN flattened_resource_audience_policies frap ON frap.resource_type = $1 frap.resource_id = organization_invitations.organization_id AND organization_invitations.organization_invitation_id = $2 AND frap.user_id = $3 AND frap.privilege = $4 AND frap.deleted_at IS NULL LIMIT $5
				`)).WithArgs("organization", invitationID, userID, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "organization_invitations" WHERE "organization_invitations"."organization_invitation_id" = $1`)).
					WithArgs(invitationID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "missing user in context",
			ctxSetup: func() context.Context {
				return context.Background()
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
		},
		{
			name: "database error",
			ctxSetup: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", userID, []uuid.UUID{})
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`
	 SELECT "organization_invitations"."organization_invitation_id","organization_invitations"."organization_id","organization_invitations"."email","organization_invitations"."privilege","organization_invitations"."created_at","organization_invitations"."invited_by","organization_invitations"."email_retry_count","organization_invitations"."email_sent_at","organization_invitations"."updated_at","organization_invitations"."deleted_at" FROM "organization_invitations" JOIN flattened_resource_audience_policies frap ON frap.resource_type = $1 frap.resource_id = organization_invitations.organization_id AND organization_invitations.organization_invitation_id = $2 AND frap.user_id = $3 AND frap.privilege = $4 AND frap.deleted_at IS NULL LIMIT $5
				`)).WithArgs("organization", invitationID, userID, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "organization_invitations" WHERE "organization_invitations"."organization_invitation_id" = $1`)).
					WithArgs(invitationID).
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
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := tt.ctxSetup()
			tt.mockSetup(mock)

			// Execute
			err := store.DeleteOrganizationInvitation(ctx, invitationID)

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

func TestGetOrganizationSSOConfigsByOrganizationId(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	configID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []models.OrganizationSSOConfig
		wantErr   bool
	}{
		{
			name: "success - configs found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_sso_config_id",
					"organization_id",
					"sso_provider_id",
					"sso_provider_name",
					"sso_config",
					"email_domain",
					"created_at",
					"updated_at",
				}).AddRow(
					configID,
					orgID,
					"provider1",
					"Provider 1",
					[]byte(`{"key": "value"}`),
					"example.com",
					now,
					now,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_sso_configs" WHERE organization_id = $1`)).
					WithArgs(orgID).
					WillReturnRows(rows)
			},
			want: []models.OrganizationSSOConfig{
				{
					OrganizationSSOConfigID: configID,
					OrganizationID:          orgID,
					SSOProviderID:           "provider1",
					SSOProviderName:         "Provider 1",
					SSOConfig:               json.RawMessage(`{"key": "value"}`),
					EmailDomain:             "example.com",
					CreatedAt:               now,
					UpdatedAt:               now,
				},
			},
		},
		{
			name: "success - no configs found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_sso_configs" WHERE organization_id = $1`)).
					WithArgs(orgID).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			want: []models.OrganizationSSOConfig{},
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_sso_configs" WHERE organization_id = $1`)).
					WithArgs(orgID).
					WillReturnError(assert.AnError)
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
			ctx := context.Background()
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationSSOConfigsByOrganizationId(ctx, orgID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_GetSSOConfigByDomain(t *testing.T) {
	domain := "test.com"
	orgID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		want      *models.OrganizationSSOConfig
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_sso_config_id",
					"organization_id",
					"email_domain",
					"is_primary",
					"created_at",
					"updated_at",
				}).AddRow(
					uuid.New(),
					orgID,
					domain,
					true,
					now,
					now,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_sso_configs" WHERE email_domain = $1 ORDER BY is_primary DESC,"organization_sso_configs"."organization_sso_config_id" LIMIT $2`)).
					WithArgs(domain, 1).
					WillReturnRows(rows)
			},
			want: &models.OrganizationSSOConfig{
				OrganizationID: orgID,
				EmailDomain:    domain,
				IsPrimary:      true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "error - not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_sso_configs" WHERE email_domain = $1 ORDER BY is_primary DESC,"organization_sso_configs"."organization_sso_config_id" LIMIT $2`)).
					WithArgs(domain, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_sso_configs" WHERE email_domain = $1 ORDER BY is_primary DESC,"organization_sso_configs"."organization_sso_config_id" LIMIT $2`)).
					WithArgs(domain, 1).
					WillReturnError(assert.AnError)
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
			ctx := context.Background()
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetSSOConfigByDomain(ctx, domain)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.EmailDomain, got.EmailDomain)
			assert.Equal(t, tt.want.IsPrimary, got.IsPrimary)
			assert.Equal(t, tt.want.CreatedAt, got.CreatedAt)
			assert.Equal(t, tt.want.UpdatedAt, got.UpdatedAt)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_GetPrimarySSOConfigByDomain(t *testing.T) {
	domain := "test.com"
	orgID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		want      *models.OrganizationSSOConfig
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"organization_sso_config_id",
					"organization_id",
					"sso_provider_id",
					"sso_provider_name",
					"sso_config",
					"email_domain",
					"is_primary",
					"created_at",
					"updated_at",
				}).AddRow(
					uuid.New(),
					orgID,
					"test-provider",
					"Test Provider",
					json.RawMessage(`{}`),
					domain,
					true,
					now,
					now,
				)

				mock.ExpectQuery(`SELECT \* FROM "organization_sso_configs" WHERE email_domain = \$1 AND is_primary = \$2 ORDER BY "organization_sso_configs"."organization_sso_config_id" LIMIT \$3`).
					WithArgs(domain, true, 1).
					WillReturnRows(rows)
			},
			want: &models.OrganizationSSOConfig{
				OrganizationID:  orgID,
				SSOProviderID:   "test-provider",
				SSOProviderName: "Test Provider",
				SSOConfig:       json.RawMessage(`{}`),
				EmailDomain:     domain,
				IsPrimary:       true,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			wantErr: false,
		},
		{
			name: "not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organization_sso_configs" WHERE email_domain = \$1 AND is_primary = \$2 ORDER BY "organization_sso_configs"."organization_sso_config_id" LIMIT \$3`).
					WithArgs(domain, true, 1).
					WillReturnError(gorm.ErrRecordNotFound)
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
			ctx := context.Background()
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetPrimarySSOConfigByDomain(ctx, domain)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.SSOProviderID, got.SSOProviderID)
			assert.Equal(t, tt.want.SSOProviderName, got.SSOProviderName)
			assert.Equal(t, tt.want.EmailDomain, got.EmailDomain)
			assert.Equal(t, tt.want.IsPrimary, got.IsPrimary)
			assert.Equal(t, tt.want.CreatedAt, got.CreatedAt)
			assert.Equal(t, tt.want.UpdatedAt, got.UpdatedAt)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateSSOConfig(t *testing.T) {
	orgID := uuid.New()
	ssoProviderID := "test-provider"
	ssoProviderName := "Test Provider"
	ssoConfig := json.RawMessage(`{"key": "value"}`)
	emailDomain := "test.com"
	userId := uuid.New()
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      *models.OrganizationSSOConfig
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(models.ResourceTypeOrganization, orgID, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgID))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_sso_configs" ("organization_id","sso_provider_id","sso_provider_name","sso_config","email_domain","is_primary","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "organization_sso_config_id"`)).
					WithArgs(orgID, ssoProviderID, ssoProviderName, ssoConfig, emailDomain, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"organization_sso_config_id"}).AddRow(uuid.New()))
				mock.ExpectCommit()
			},
			want: &models.OrganizationSSOConfig{
				OrganizationID:  orgID,
				SSOProviderID:   ssoProviderID,
				SSOProviderName: ssoProviderName,
				SSOConfig:       ssoConfig,
				EmailDomain:     emailDomain,
				IsPrimary:       false,
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3`)).
					WithArgs(models.ResourceTypeOrganization, orgID, userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgID))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_sso_configs" ("organization_id","sso_provider_id","sso_provider_name","sso_config","email_domain","is_primary","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "organization_sso_config_id"`)).
					WithArgs(orgID, ssoProviderID, ssoProviderName, ssoConfig, emailDomain, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgID})
			tt.mockSetup(mock)

			// Execute
			got, err := store.CreateSSOConfig(ctx, orgID, ssoProviderID, ssoProviderName, ssoConfig, emailDomain)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.SSOProviderID, got.SSOProviderID)
			assert.Equal(t, tt.want.SSOProviderName, got.SSOProviderName)
			assert.Equal(t, tt.want.EmailDomain, got.EmailDomain)
			assert.Equal(t, tt.want.IsPrimary, got.IsPrimary)
			assert.JSONEq(t, string(tt.want.SSOConfig), string(got.SSOConfig))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_CreateOrganizationInvitationStatus(t *testing.T) {
	invitationId := uuid.New()
	status := models.InvitationStatusAccepted

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		want      *models.OrganizationInvitationStatus
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationId, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationId))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_invitation_statuses" ("organization_invitation_id","status") VALUES ($1,$2) RETURNING "created_at","updated_at"`)).
					WithArgs(invitationId, status).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationId))
				mock.ExpectCommit()
			},
			want: &models.OrganizationInvitationStatus{
				OrganizationInvitationID: invitationId,
				Status:                   status,
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationId, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationId))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_invitation_statuses" ("organization_invitation_id","status") VALUES ($1,$2) RETURNING "created_at","updated_at"`)).
					WithArgs(invitationId, status).
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

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := apicontext.AddAuthToContext(context.Background(), "user", uuid.New(), []uuid.UUID{})
			tt.mockSetup(mock)

			// Execute
			got, err := store.CreateOrganizationInvitationStatus(ctx, invitationId, status)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationInvitationID, got.OrganizationInvitationID)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_GetOrganizationMembershipRequestsByOrganizationId(t *testing.T) {
	orgId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []models.OrganizationMembershipRequest
		wantErr   bool
	}{
		{
			name: "success - returns requests",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"organization_id", "user_id", "status"}).
					AddRow(orgId, uuid.New(), models.OrgMembershipStatusPending).
					AddRow(orgId, uuid.New(), models.OrgMembershipStatusPending)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests" WHERE organization_id = $1 AND status = $2`)).
					WithArgs(orgId, models.OrgMembershipStatusPending).
					WillReturnRows(rows)
			},
			want: []models.OrganizationMembershipRequest{
				{
					OrganizationID: orgId,
					Status:         models.OrgMembershipStatusPending,
				},
				{
					OrganizationID: orgId,
					Status:         models.OrgMembershipStatusPending,
				},
			},
			wantErr: false,
		},
		{
			name: "success - no requests found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests" WHERE organization_id = $1 AND status = $2`)).
					WithArgs(orgId, models.OrgMembershipStatusPending).
					WillReturnRows(sqlmock.NewRows([]string{"organization_id", "user_id", "status"}))
			},
			want:    []models.OrganizationMembershipRequest{},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests" WHERE organization_id = $1 AND status = $2`)).
					WithArgs(orgId, models.OrgMembershipStatusPending).
					WillReturnError(assert.AnError)
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
			ctx := context.Background()
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationMembershipRequestsByOrganizationId(ctx, orgId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, got, len(tt.want))
			for i := range got {
				assert.Equal(t, tt.want[i].OrganizationID, got[i].OrganizationID)
				assert.Equal(t, tt.want[i].Status, got[i].Status)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_GetOrganizationMembershipRequestsAll(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	userId := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []models.OrganizationMembershipRequest
		wantErr   bool
	}{
		{
			name: "success - multiple requests",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"organization_id",
					"user_id",
					"created_at",
					"status",
					"updated_at",
					"deleted_at",
				}).
					AddRow(uuid.New(), orgId, userId, now, models.OrgMembershipStatusPending, now, nil).
					AddRow(uuid.New(), orgId, userId, now, models.OrgMembershipStatusPending, now, nil)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests"`)).
					WillReturnRows(rows)
			},
			want: []models.OrganizationMembershipRequest{
				{
					OrganizationID: orgId,
					UserID:         userId,
					CreatedAt:      now,
					Status:         models.OrgMembershipStatusPending,
					UpdatedAt:      now,
				},
				{
					OrganizationID: orgId,
					UserID:         userId,
					CreatedAt:      now,
					Status:         models.OrgMembershipStatusPending,
					UpdatedAt:      now,
				},
			},
			wantErr: false,
		},
		{
			name: "success - empty result",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "organization_id", "user_id", "status"}))
			},
			want:    []models.OrganizationMembershipRequest{},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests"`)).
					WillReturnError(assert.AnError)
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
			ctx := context.Background()
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationMembershipRequestsAll(ctx)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, got, len(tt.want))
			for i := range got {
				assert.Equal(t, tt.want[i].OrganizationID, got[i].OrganizationID)
				assert.Equal(t, tt.want[i].UserID, got[i].UserID)
				assert.Equal(t, tt.want[i].Status, got[i].Status)
				assert.Equal(t, tt.want[i].CreatedAt.Unix(), got[i].CreatedAt.Unix())
				assert.Equal(t, tt.want[i].UpdatedAt.Unix(), got[i].UpdatedAt.Unix())
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_CreateOrganizationMembershipRequest(t *testing.T) {
	orgId := uuid.New()
	userId := uuid.New()
	status := models.OrgMembershipStatusPending

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      *models.OrganizationMembershipRequest
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_membership_requests" ("organization_id","user_id","status","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "id","created_at","updated_at"`)).
					WithArgs(orgId, userId, status, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(uuid.New(), time.Now(), time.Now()))

				mock.ExpectCommit()
			},
			want: &models.OrganizationMembershipRequest{
				OrganizationID: orgId,
				UserID:         userId,
				Status:         status,
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organization_membership_requests" ("organization_id","user_id","status","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "id","created_at","updated_at"`)).
					WithArgs(orgId, userId, status, sqlmock.AnyArg()).
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

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			// Execute
			got, err := store.CreateOrganizationMembershipRequest(ctx, orgId, userId, status)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.UserID, got.UserID)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_UpdatePendingOrganizationMembershipRequest(t *testing.T) {
	orgId := uuid.New()
	userId := uuid.New()
	status := models.OrgMembershipStatusApproved

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		want      *models.OrganizationMembershipRequest
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 LIMIT $5`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "organization_membership_requests" SET "status"=$1,"updated_at"=$2 WHERE organization_id = $3 AND user_id = $4 AND status = $5 RETURNING *`)).
					WithArgs(status, sqlmock.AnyArg(), orgId, userId, models.OrgMembershipStatusPending).
					WillReturnRows(sqlmock.NewRows([]string{"organization_id", "user_id", "status"}).AddRow(orgId, userId, status))
				mock.ExpectCommit()
			},
			want: &models.OrganizationMembershipRequest{
				OrganizationID: orgId,
				UserID:         userId,
				Status:         status,
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 LIMIT $5`)).
					WithArgs(models.ResourceTypeOrganization, orgId, userId, models.PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgId))
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "organization_membership_requests" SET "status"=$1,"updated_at"=$2 WHERE organization_id = $3 AND user_id = $4 AND status = $5 RETURNING *`)).
					WithArgs(status, sqlmock.AnyArg(), orgId, userId, models.OrgMembershipStatusPending).
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

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			// Execute
			got, err := store.UpdatePendingOrganizationMembershipRequest(ctx, orgId, userId, status)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.want.UserID, got.UserID)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_CreateOrganization(t *testing.T) {
	orgId := uuid.New()
	userId := uuid.New()
	name := "Test Organization"
	description := "Test Description"

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		want      *models.Organization
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organizations" ("name","description","deleted_at","owner_id") VALUES ($1,$2,$3,$4) RETURNING "organization_id","created_at","updated_at"`)).
					WithArgs(name, &description, sqlmock.AnyArg(), userId).
					WillReturnRows(sqlmock.NewRows([]string{"organization_id"}).AddRow(orgId))
				mock.ExpectCommit()
			},
			want: &models.Organization{
				ID:          orgId,
				Name:        name,
				Description: &description,
				OwnerId:     userId,
			},
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organizations" ("name","description","deleted_at","owner_id") VALUES ($1,$2,$3,$4) RETURNING "organization_id","created_at","updated_at"`)).
					WithArgs(name, &description, sqlmock.AnyArg(), userId).
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

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := apicontext.AddAuthToContext(context.Background(), "admin", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			// Execute
			got, err := store.CreateOrganization(ctx, name, &description, userId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.OwnerId, got.OwnerId)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestAppStore_GetOrganizationInvitationsAndMembershipRequests(t *testing.T) {
	orgId := uuid.New()
	userId := uuid.New()
	name := "Test Organization"
	description := "Test Description"
	t.Parallel()

	now := time.Now()

	invitationStatus := models.OrganizationInvitationStatus{
		OrganizationInvitationID: uuid.New(),
		Status:                   models.InvitationStatusAccepted,
		UpdatedAt:                now,
	}

	invitation := models.OrganizationInvitation{
		OrganizationInvitationID: invitationStatus.OrganizationInvitationID,
		OrganizationID:           orgId,
		TargetEmail:              "test@example.com",
		InvitedBy:                userId,
		Privilege:                models.PrivilegeOrganizationSystemAdmin,
		InvitationStatuses: []models.OrganizationInvitationStatus{
			invitationStatus,
		},
	}

	user := models.User{
		ID:    userId,
		Email: "test@example.com",
		Name:  "Test User",
	}

	membershipRequest := models.OrganizationMembershipRequest{
		OrganizationID: orgId,
		UserID:         userId,
		Status:         models.OrgMembershipStatusPending,
		User:           user,
	}

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      *models.Organization
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				orgRows := sqlmock.NewRows([]string{
					"organization_id",
					"name",
					"description",
					"created_at",
					"updated_at",
					"deleted_at",
					"owner_id",
				}).AddRow(
					orgId,
					name,
					&description,
					now,
					now,
					nil,
					userId,
				)

				invitationRows := sqlmock.NewRows([]string{
					"organization_invitation_id",
					"organization_id",
					"target_email",
					"invited_by",
					"privilege",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					invitation.OrganizationInvitationID,
					invitation.OrganizationID,
					invitation.TargetEmail,
					invitation.InvitedBy,
					invitation.Privilege,
					now,
					now,
					nil,
				)

				invitationStatusRows := sqlmock.NewRows([]string{
					"organization_invitation_status_id",
					"organization_invitation_id",
					"status",
					"created_at",
					"updated_at",
				}).AddRow(
					invitationStatus.OrganizationInvitationID,
					invitationStatus.OrganizationInvitationID,
					invitationStatus.Status,
					invitationStatus.UpdatedAt,
					invitationStatus.UpdatedAt,
				)

				membershipRequestRows := sqlmock.NewRows([]string{
					"organization_membership_request_id",
					"organization_id",
					"user_id",
					"status",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					uuid.New(),
					membershipRequest.OrganizationID,
					membershipRequest.UserID,
					membershipRequest.Status,
					now,
					now,
					nil,
				)

				userRows := sqlmock.NewRows([]string{
					"user_id",
					"email",
					"name",
				}).AddRow(
					user.ID,
					user.Email,
					user.Name,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE organization_id = $1 ORDER BY "organizations"."organization_id" LIMIT $2`)).
					WithArgs(orgId, 1).
					WillReturnRows(orgRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE "organization_invitations"."organization_id" = $1`)).
					WithArgs(orgId).
					WillReturnRows(invitationRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitation_statuses" WHERE "organization_invitation_statuses"."organization_invitation_id" = $1`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(invitationStatusRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests" WHERE "organization_membership_requests"."organization_id" = $1`)).
					WithArgs(orgId).
					WillReturnRows(membershipRequestRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(userId).
					WillReturnRows(userRows)
			},
			want: &models.Organization{
				ID:                 orgId,
				Name:               name,
				Description:        &description,
				OwnerId:            userId,
				Invitations:        []models.OrganizationInvitation{invitation},
				MembershipRequests: []models.OrganizationMembershipRequest{membershipRequest},
			},
			wantErr: false,
		},
		{
			name: "error - organization not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE organization_id = $1 ORDER BY "organizations"."organization_id" LIMIT 1`)).
					WithArgs(orgId).
					WillReturnError(gorm.ErrRecordNotFound)
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
			ctx := context.Background()
			tt.mockSetup(mock)

			// Execute
			got, err := store.GetOrganizationInvitationsAndMembershipRequests(ctx, orgId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.OwnerId, got.OwnerId)
			assert.Len(t, got.Invitations, len(tt.want.Invitations))
			assert.Len(t, got.MembershipRequests, len(tt.want.MembershipRequests))
			assert.Equal(t, tt.want.MembershipRequests[0].User, got.MembershipRequests[0].User)
			assert.Equal(t, tt.want.Invitations[0].InvitationStatuses[0].Status, got.Invitations[0].InvitationStatuses[0].Status)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
