package pages

import (
	"context"
	"testing"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// String constants for audience types
const (
	stringUser         = "user"
	stringOrganization = "organization"
	stringTeam         = "team"
)

func TestEnsureCurrentUsersAdminAccess(t *testing.T) {
	tests := []struct {
		name          string
		currentUserId *uuid.UUID
		policies      []models.ResourceAudiencePolicy
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "Error - No user ID in context",
			currentUserId: nil,
			policies:      []models.ResourceAudiencePolicy{},
			wantErr:       true,
			expectedErr:   "no user ID found in the context",
		},
		{
			name: "Error - User does not have admin access",
			currentUserId: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			policies: []models.ResourceAudiencePolicy{
				{
					ResourceAudienceType: models.AudienceTypeUser,
					ResourceAudienceID:   uuid.New(),
					Privilege:            "admin",
					UserPolicies: []models.FlattenedResourceAudiencePolicy{
						{
							UserId:    uuid.New(),
							Privilege: "admin",
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: "current user does not have access to change permissions on the page",
		},
		{
			name: "Success - User has admin access through UserPolicies",
			currentUserId: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			policies: []models.ResourceAudiencePolicy{
				{
					ResourceAudienceType: models.AudienceTypeUser,
					ResourceAudienceID:   uuid.New(),
					Privilege:            "admin",
					UserPolicies: []models.FlattenedResourceAudiencePolicy{
						{
							UserId:    uuid.New(),
							Privilege: "admin",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of currentUserId to use in policies
			var currentId uuid.UUID
			if tt.currentUserId != nil {
				currentId = *tt.currentUserId
				// Add a policy for the current user with admin privilege
				if !tt.wantErr {
					for i := range tt.policies {
						tt.policies[i].UserPolicies = append(tt.policies[i].UserPolicies, models.FlattenedResourceAudiencePolicy{
							UserId:    currentId,
							Privilege: models.PrivilegePageAdmin,
						})
					}
				}
			}

			ctx := context.Background()
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", currentId, []uuid.UUID{})
			}

			err := ensureCurrentUsersAdminAccess(ctx, tt.policies)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEnsureAudienceNotAlreadyAdded(t *testing.T) {
	audienceId := uuid.New()

	tests := []struct {
		name         string
		audienceType models.AudienceType
		audienceId   uuid.UUID
		policies     []models.ResourceAudiencePolicy
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "Error - Audience already exists",
			audienceType: models.AudienceTypeUser,
			audienceId:   audienceId,
			policies: []models.ResourceAudiencePolicy{
				{
					ResourceAudienceType: stringUser,
					ResourceAudienceID:   audienceId,
				},
			},
			wantErr:     true,
			expectedErr: "audience already exists on the page",
		},
		{
			name:         "Success - Audience does not exist",
			audienceType: models.AudienceTypeUser,
			audienceId:   audienceId,
			policies: []models.ResourceAudiencePolicy{
				{
					ResourceAudienceType: models.AudienceTypeTeam,
					ResourceAudienceID:   uuid.New(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ensureAudienceNotAlreadyAdded(tt.audienceType, tt.audienceId, tt.policies)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEnsureUserIsNotChangingTheirOwnAdminPolicy(t *testing.T) {
	currentUserId := uuid.New()
	policyId := uuid.New()

	tests := []struct {
		name              string
		policyToBeUpdated models.ResourceAudiencePolicy
		allPolicies       []models.ResourceAudiencePolicy
		currentUserId     uuid.UUID
		wantErr           bool
		expectedErr       string
	}{
		{
			name: "Error - User changing their own direct admin policy",
			policyToBeUpdated: models.ResourceAudiencePolicy{
				ResourceAudienceType: stringUser,
				ResourceAudienceID:   currentUserId,
				Privilege:            models.PrivilegePageAdmin,
			},
			allPolicies:   []models.ResourceAudiencePolicy{},
			currentUserId: currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
		{
			name: "Error - User changing org admin policy without separate admin access",
			policyToBeUpdated: models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: stringOrganization,
				Privilege:            models.PrivilegePageAdmin,
				UserPolicies: []models.FlattenedResourceAudiencePolicy{
					{
						ResourceAudiencePolicyId: policyId,
						UserId:                   currentUserId,
						Privilege:                models.PrivilegePageAdmin,
					},
				},
			},
			allPolicies: []models.ResourceAudiencePolicy{
				{
					ID:                   policyId,
					ResourceAudienceType: stringOrganization,
					Privilege:            models.PrivilegePageAdmin,
					UserPolicies: []models.FlattenedResourceAudiencePolicy{
						{
							ResourceAudiencePolicyId: policyId,
							UserId:                   currentUserId,
							Privilege:                models.PrivilegePageAdmin,
						},
					},
				},
			},
			currentUserId: currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
		{
			name: "Success - User changing org admin policy with separate admin access",
			policyToBeUpdated: models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: stringOrganization,
				Privilege:            models.PrivilegePageAdmin,
				UserPolicies: []models.FlattenedResourceAudiencePolicy{
					{
						ResourceAudiencePolicyId: policyId,
						UserId:                   currentUserId,
						Privilege:                models.PrivilegePageAdmin,
					},
				},
			},
			allPolicies: []models.ResourceAudiencePolicy{
				{
					ID:                   uuid.New(),
					ResourceAudienceType: stringOrganization,
					Privilege:            models.PrivilegePageAdmin,
					UserPolicies: []models.FlattenedResourceAudiencePolicy{
						{
							ResourceAudiencePolicyId: uuid.New(),
							UserId:                   currentUserId,
							Privilege:                models.PrivilegePageAdmin,
						},
					},
				},
			},
			currentUserId: currentUserId,
			wantErr:       false,
		},
		{
			name: "Success - User changing team admin policy with separate admin access",
			policyToBeUpdated: models.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: stringTeam,
				Privilege:            models.PrivilegePageAdmin,
				UserPolicies: []models.FlattenedResourceAudiencePolicy{
					{
						ResourceAudiencePolicyId: policyId,
						UserId:                   currentUserId,
						Privilege:                models.PrivilegePageAdmin,
					},
				},
			},
			allPolicies: []models.ResourceAudiencePolicy{
				{
					ID:                   uuid.New(),
					ResourceAudienceType: stringTeam,
					Privilege:            models.PrivilegePageAdmin,
					UserPolicies: []models.FlattenedResourceAudiencePolicy{
						{
							ResourceAudiencePolicyId: uuid.New(),
							UserId:                   currentUserId,
							Privilege:                models.PrivilegePageAdmin,
						},
					},
				},
			},
			currentUserId: currentUserId,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ensureUserIsNotChangingTheirOwnAdminPolicy(tt.policyToBeUpdated, tt.allPolicies, tt.currentUserId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}
