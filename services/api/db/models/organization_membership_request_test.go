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

func TestOrganizationMembershipRequest_TableName(t *testing.T) {
	o := &OrganizationMembershipRequest{}
	assert.Equal(t, "organization_membership_requests", o.TableName())
}

func TestOrganizationMembershipRequest_GetQueryFilters(t *testing.T) {
	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New(), uuid.New()}

	db, mock := setupTestDB(t)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "organization_membership_requests" WHERE EXISTS ( SELECT 1 FROM app.flattened_resource_audience_policies frap WHERE frap.resource_type = $1 AND frap.resource_id = organization_membership_requests.organization_id AND frap.user_id = $2 AND frap.privilege = $3 ) OR organization_membership_requests.user_id = $4`)

	mock.ExpectQuery(expectedSQL).
		WithArgs(ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin, userId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	o := &OrganizationMembershipRequest{}
	filteredDB := o.GetQueryFilters(db, userId, orgIds)

	// Execute the query to verify the SQL
	var result []OrganizationMembershipRequest
	err := filteredDB.Find(&result).Error
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestOrganizationMembershipRequest_BeforeCreate(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setup       func() (*OrganizationMembershipRequest, context.Context)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - same user",
			setup: func() (*OrganizationMembershipRequest, context.Context) {
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				return &OrganizationMembershipRequest{
					UserID: userId,
				}, ctx
			},
			wantErr: false,
		},
		{
			name: "error - different user",
			setup: func() (*OrganizationMembershipRequest, context.Context) {
				ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
				return &OrganizationMembershipRequest{
					UserID: uuid.New(), // Different user ID
				}, ctx
			},
			wantErr:     true,
			errContains: "user cannot create request for another user",
		},
		{
			name: "error - no user in context",
			setup: func() (*OrganizationMembershipRequest, context.Context) {
				return &OrganizationMembershipRequest{
					UserID: uuid.New(),
				}, context.Background()
			},
			wantErr:     true,
			errContains: "no user id found in context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupTestDB(t)
			req, ctx := tt.setup()
			db = db.WithContext(ctx)

			err := req.BeforeCreate(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrganizationMembershipRequest_BeforeUpdate(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - user has permission",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 LIMIT $5`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(orgId))
			},
			wantErr: false,
		},
		{
			name: "error - no permission",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 LIMIT $5`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_id"}))
			},
			wantErr:     true,
			errContains: "organization access forbidden",
		},
		{
			name: "error - db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 LIMIT $5`)).
					WithArgs(ResourceTypeOrganization, orgId, userId, PrivilegeOrganizationSystemAdmin, 1).
					WillReturnError(assert.AnError)
			},
			wantErr:     true,
			errContains: "failed to get flattened resource audience policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			ctx := apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId})
			db = db.WithContext(ctx)

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			req := &OrganizationMembershipRequest{
				OrganizationID: orgId,
			}

			err := req.BeforeUpdate(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
