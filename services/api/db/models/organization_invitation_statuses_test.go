package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationInvitationStatus_TableName(t *testing.T) {
	invitationStatus := &OrganizationInvitationStatus{}
	assert.Equal(t, "organization_invitation_statuses", invitationStatus.TableName())
}

func TestOrganizationInvitationStatus_GetQueryFilters(t *testing.T) {
	// Setup
	db, mock := setupTestDB(t)

	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New()}
	invitationStatus := &OrganizationInvitationStatus{OrganizationInvitationID: uuid.New()}

	// Mock the query execution
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitation_statuses" WHERE EXISTS ( SELECT 1 FROM app.flattened_resource_audience_policies frap, app.organization_invitations oi WHERE frap.resource_type = $1 AND frap.resource_id = oi.organization_id AND oi.organization_invitation_id = organization_invitation_statuses.organization_invitation_id AND frap.user_id = $2 AND frap.privilege = $3 AND frap.deleted_at IS NULL ) OR EXISTS ( SELECT 1 FROM app.users_with_traits uwt, app.organization_invitations oi WHERE oi.organization_invitation_id = organization_invitation_statuses.organization_invitation_id AND uwt.email = oi.email AND uwt.user_id = $4 )`)).
		WithArgs(ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin, userId).
		WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationStatus.OrganizationInvitationID))

	// Execute test
	query := invitationStatus.GetQueryFilters(db, userId, orgIds)

	// Execute the query to trigger the mock
	var result []OrganizationInvitationStatus
	err := query.Find(&result).Error
	assert.Nil(t, err)

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrganizationInvitationStatus_BeforeCreate(t *testing.T) {

	userId := uuid.New()
	invitationId := uuid.New()

	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupMock     func(sqlmock.Sqlmock)
		invStatus     *OrganizationInvitationStatus
		expectedError bool
	}{
		{
			name: "successful create with valid permissions",
			setupContext: func() context.Context {
				return setupAuthContext(userId)
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}).AddRow(invitationId))
			},
			invStatus: &OrganizationInvitationStatus{
				OrganizationInvitationID: invitationId,
				Status:                   InvitationStatusAccepted,
			},
			expectedError: false,
		},
		{
			name: "failed create - no permissions",
			setupContext: func() context.Context {
				return setupAuthContext(userId)
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"organization_invitation_id"}))
			},
			invStatus: &OrganizationInvitationStatus{
				OrganizationInvitationID: invitationId,
				Status:                   InvitationStatusAccepted,
			},
			expectedError: true,
		},
		{
			name: "failed create - no user ID in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No op
			},
			invStatus: &OrganizationInvitationStatus{
				OrganizationInvitationID: invitationId,
				Status:                   InvitationStatusAccepted,
			},
			expectedError: true,
		},
		{
			name: "failed create - db error",
			setupContext: func() context.Context {
				return setupAuthContext(userId)
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organization_invitations" WHERE EXISTS ( SELECT 1 FROM users_with_traits uwt WHERE organization_invitation_id = $1 AND uwt.email = organization_invitations.email AND uwt.user_id = $2 )`)).
					WithArgs(invitationId, userId).
					WillReturnError(assert.AnError)
			},
			invStatus: &OrganizationInvitationStatus{
				OrganizationInvitationID: invitationId,
				Status:                   InvitationStatusAccepted,
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
			err := tt.invStatus.BeforeCreate(db)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
