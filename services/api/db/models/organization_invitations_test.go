package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// setupAuthContext creates a context with user authentication
func setupAuthContext(userId uuid.UUID) context.Context {
	ctx := context.Background()
	return apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{})
}

func TestOrganizationInvitation_TableName(t *testing.T) {
	invitation := &OrganizationInvitation{}
	assert.Equal(t, "organization_invitations", invitation.TableName())
}

func TestOrganizationInvitation_GetQueryFilters(t *testing.T) {
	// Setup
	db, mock := setupTestDB(t)

	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New()}
	invitation := &OrganizationInvitation{OrganizationInvitationID: uuid.New()}

	// Mock the query execution
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE ( EXISTS ( SELECT 1 FROM app.flattened_resource_audience_policies frap WHERE frap.resource_type = 'organization' AND frap.resource_id = organization_invitations.organization_id AND frap.user_id = $1 AND frap.privilege = $2 AND frap.deleted_at IS NULL ) OR EXISTS ( SELECT 1 FROM app.users_with_traits uwt WHERE uwt.user_id = $3 AND uwt.email = organization_invitations.email ) ) AND ( organization_invitation_id NOT IN ( SELECT organization_invitation_id FROM app.organization_invitation_statuses ois WHERE ois.organization_invitation_id = organization_invitations.organization_invitation_id ) )`)).
		WithArgs(userId, PrivilegeOrganizationSystemAdmin, userId).
		WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitation.OrganizationInvitationID))

	// Execute test
	query := invitation.GetQueryFilters(db, userId, orgIds)

	// Execute the query to trigger the mock
	var result []OrganizationInvitation
	err := query.Find(&result).Error
	assert.Nil(t, err)
}

func TestOrganizationInvitation_BeforeCreate(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupMock     func(sqlmock.Sqlmock)
		invitation    *OrganizationInvitation
		expectedError bool
	}{
		{
			name: "successful creation with valid permissions",
			setupContext: func() context.Context {
				return setupAuthContext(uuid.New())
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow("organization", uuid.New(), uuid.New(), PrivilegeOrganizationSystemAdmin))
			},
			invitation: &OrganizationInvitation{
				OrganizationID: uuid.New(),
				TargetEmail:    "test@example.com",
				Privilege:      ResourcePrivilege("admin"),
				InvitedBy:      uuid.New(),
			},
			expectedError: false,
		},
		{
			name: "failed creation - no user ID in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func(mock sqlmock.Sqlmock) {},
			invitation: &OrganizationInvitation{
				OrganizationID: uuid.New(),
			},
			expectedError: true,
		},
		{
			name: "failed creation - no organization access",
			setupContext: func() context.Context {
				return setupAuthContext(uuid.New())
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}))
			},
			invitation: &OrganizationInvitation{
				OrganizationID: uuid.New(),
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			db, mock := setupTestDB(t)

			// Setup mock expectations
			tt.setupMock(mock)

			// Set context
			db = db.WithContext(tt.setupContext())

			// Execute test
			err := tt.invitation.BeforeCreate(db)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrganizationInvitation_BeforeUpdate(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupMock     func(sqlmock.Sqlmock)
		invitation    *OrganizationInvitation
		expectedError bool
	}{
		{
			name: "successful update with valid permissions",
			setupContext: func() context.Context {
				return setupAuthContext(uuid.New())
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}).
						AddRow("organization", uuid.New(), uuid.New(), PrivilegeOrganizationSystemAdmin))
			},
			invitation: &OrganizationInvitation{
				OrganizationID: uuid.New(),
			},
			expectedError: false,
		},
		{
			name: "failed update - no user ID in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func(mock sqlmock.Sqlmock) {},
			invitation: &OrganizationInvitation{
				OrganizationID: uuid.New(),
			},
			expectedError: true,
		},
		{
			name: "failed update - no organization access",
			setupContext: func() context.Context {
				return setupAuthContext(uuid.New())
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id", "privilege"}))
			},
			invitation: &OrganizationInvitation{
				OrganizationID: uuid.New(),
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB

			db, mock := setupTestDB(t)

			// Setup mock expectations
			tt.setupMock(mock)

			// Set context
			db = db.WithContext(tt.setupContext())

			// Execute test
			err := tt.invitation.BeforeUpdate(db)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
