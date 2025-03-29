package service

import (
	"context"
	"testing"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
)

func TestEnsureCurrentUsersAdminAccess(t *testing.T) {
	tests := []struct {
		name          string
		currentUserId *uuid.UUID
		policies      []storemodels.ResourceAudiencePolicy
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "Error - No user ID in context",
			currentUserId: nil,
			policies:      []storemodels.ResourceAudiencePolicy{},
			wantErr:       true,
			expectedErr:   "no user ID found in the context",
		},
		{
			name: "Error - User does not have admin access",
			currentUserId: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			policies: []storemodels.ResourceAudiencePolicy{
				{
					ID:                   uuid.New(),
					ResourceAudienceType: storemodels.AudienceTypeUser,
					ResourceAudienceID:   uuid.New(),
					Privilege:            storemodels.PrivilegeDatasetAdmin,
					UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
						{
							UserId:    uuid.New(),
							Privilege: storemodels.PrivilegeDatasetAdmin,
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: "current user does not have access to change permissions on the dataset",
		},
		{
			name: "Success - User has admin access",
			currentUserId: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			policies: []storemodels.ResourceAudiencePolicy{
				{
					ID:                   uuid.New(),
					ResourceAudienceType: storemodels.AudienceTypeUser,
					ResourceAudienceID:   uuid.New(),
					Privilege:            storemodels.PrivilegeDatasetViewer,
					UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
						{
							UserId:    uuid.New(),
							Privilege: storemodels.PrivilegeDatasetViewer,
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
					policyId := uuid.New()
					tt.policies = append(tt.policies, storemodels.ResourceAudiencePolicy{
						ID:                   policyId,
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   uuid.New(),
						Privilege:            storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId,
								UserId:                   currentId,
								Privilege:                storemodels.PrivilegeDatasetAdmin,
							},
						},
					})
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
		audienceType storemodels.AudienceType
		audienceId   uuid.UUID
		policies     []storemodels.ResourceAudiencePolicy
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "Error - Audience already exists",
			audienceType: storemodels.AudienceTypeUser,
			audienceId:   audienceId,
			policies: []storemodels.ResourceAudiencePolicy{
				{
					ResourceAudienceType: storemodels.AudienceTypeUser,
					ResourceAudienceID:   audienceId,
				},
			},
			wantErr:     true,
			expectedErr: "audience already exists on the dataset",
		},
		{
			name:         "Success - Audience does not exist",
			audienceType: storemodels.AudienceTypeUser,
			audienceId:   audienceId,
			policies: []storemodels.ResourceAudiencePolicy{
				{
					ResourceAudienceType: storemodels.AudienceTypeTeam,
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
		policyToBeUpdated storemodels.ResourceAudiencePolicy
		allPolicies       []storemodels.ResourceAudiencePolicy
		currentUserId     uuid.UUID
		wantErr           bool
		expectedErr       string
	}{
		{
			name: "Error - User changing their own direct admin policy",
			policyToBeUpdated: storemodels.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: storemodels.AudienceTypeUser,
				ResourceAudienceID:   currentUserId,
				Privilege:            storemodels.PrivilegeDatasetAdmin,
			},
			allPolicies:   []storemodels.ResourceAudiencePolicy{},
			currentUserId: currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
		{
			name: "Error - User changing org admin policy without separate admin access",
			policyToBeUpdated: storemodels.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: storemodels.AudienceTypeOrganization,
				Privilege:            storemodels.PrivilegeDatasetAdmin,
				UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
					{
						ResourceAudiencePolicyId: policyId,
						UserId:                   currentUserId,
						Privilege:                storemodels.PrivilegeDatasetViewer,
					},
				},
			},
			allPolicies: []storemodels.ResourceAudiencePolicy{
				{
					ID:                   policyId,
					ResourceAudienceType: storemodels.AudienceTypeOrganization,
					Privilege:            storemodels.PrivilegeDatasetViewer,
					UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
						{
							ResourceAudiencePolicyId: policyId,
							UserId:                   currentUserId,
							Privilege:                storemodels.PrivilegeDatasetViewer,
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
			policyToBeUpdated: storemodels.ResourceAudiencePolicy{
				ID:                   policyId,
				ResourceAudienceType: storemodels.AudienceTypeOrganization,
				Privilege:            storemodels.PrivilegeDatasetAdmin,
			},
			allPolicies: []storemodels.ResourceAudiencePolicy{
				{
					ID:                   uuid.New(),
					ResourceAudienceType: storemodels.AudienceTypeOrganization,
					Privilege:            storemodels.PrivilegeDatasetAdmin,
					UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
						{
							ResourceAudiencePolicyId: uuid.New(),
							UserId:                   currentUserId,
							Privilege:                storemodels.PrivilegeDatasetAdmin,
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
