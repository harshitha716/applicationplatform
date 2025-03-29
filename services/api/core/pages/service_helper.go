package pages

import (
	"context"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

func ensureCurrentUsersAdminAccess(ctx context.Context, policies []models.ResourceAudiencePolicy) error {
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)

	if currentUserId == nil {
		return fmt.Errorf("no user ID found in the context")
	}

	for _, policy := range policies {
		for _, userPolicy := range policy.UserPolicies {
			// user has admin access directly
			if userPolicy.UserId == *currentUserId && userPolicy.Privilege == models.PrivilegePageAdmin {
				return nil
			}
		}
	}

	return fmt.Errorf("current user does not have access to change permissions on the page")
}

func ensureAudienceNotAlreadyAdded(audienceType models.AudienceType, audienceId uuid.UUID, policies []models.ResourceAudiencePolicy) error {
	for _, policy := range policies {
		if policy.ResourceAudienceType == audienceType && policy.ResourceAudienceID == audienceId {
			return fmt.Errorf("audience already exists on the page")
		}
	}

	return nil
}

func ensureUserIsNotChangingTheirOwnAdminPolicy(policyToBeUpdated models.ResourceAudiencePolicy, allPolicies []models.ResourceAudiencePolicy, currentUserId uuid.UUID) error {
	if policyToBeUpdated.ResourceAudienceType == models.AudienceTypeUser && policyToBeUpdated.ResourceAudienceID == currentUserId && policyToBeUpdated.Privilege == models.PrivilegePageAdmin {
		return fmt.Errorf("you cannot change own permissions")
	}

	isCurrentUserSeparatelyAdded := false
	if policyToBeUpdated.ResourceAudienceType == models.AudienceTypeOrganization || policyToBeUpdated.ResourceAudienceType == models.AudienceTypeTeam {
		for _, policy := range allPolicies {
			for _, userPolicy := range policy.UserPolicies {
				if policyToBeUpdated.ID != userPolicy.ResourceAudiencePolicyId && userPolicy.UserId == currentUserId && userPolicy.Privilege == models.PrivilegePageAdmin {
					isCurrentUserSeparatelyAdded = true
					return nil
				}
			}
		}

		if !isCurrentUserSeparatelyAdded {
			return fmt.Errorf("you cannot change own permissions")
		}
	}

	return nil
}
