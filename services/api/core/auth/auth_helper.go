package auth

import (
	"fmt"
	"net/http"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/helper"
	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

func kratosIdentityToZampUser(identity *kratos.Identity) (*models.User, error) {

	userId, err := uuid.Parse(identity.Id)
	if err != nil {
		return nil, err
	}

	if identity.Traits == nil {
		return nil, fmt.Errorf("user traits not found")
	}

	userTraits, ok := identity.Traits.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user traits")
	}

	nameRaw, nameOk := userTraits["name"]
	if !nameOk {
		nameRaw = interface{}("")
	}

	emailRaw, ok := userTraits["email"]
	if !ok {
		return nil, fmt.Errorf("email not found in user traits")
	}

	email, ok := emailRaw.(string)
	if !ok {
		return nil, fmt.Errorf("invalid email")
	}

	name, ok := nameRaw.(string)
	if !ok {
		name = ""
	}

	return &models.User{
		ID:    userId,
		Email: email,
		Name:  name,
	}, nil
}

// SelectOrganization selects a single organization ID based on the header and user's organizations
// It follows the same logic as the middleware:
// 1. Gets the organization ID from the X-Zamp-Organization-Id header
// 2. Checks if the organization ID is in the user's allowed organizations
// 3. If found, uses that single organization ID
// 4. If not found, proceeds with an empty list
// 5. If no header is provided and user has organizations, uses the first organization
// 6. If user has no organizations, proceeds with an empty list
func SelectOrganization(header http.Header, organizationIds []uuid.UUID, logger *zap.Logger) []uuid.UUID {
	// Get the organization ID from the header
	orgIdHeader := header.Get(helper.PROXY_WORKSPACE_ID_HEADER)
	var selectedOrgIds []uuid.UUID

	if orgIdHeader != "" && len(organizationIds) > 0 {
		// Try to parse the organization ID from the header
		orgId, err := uuid.Parse(orgIdHeader)
		if err != nil {
			logger.Error("invalid organization id in header", zap.String("error", err.Error()))
			return []uuid.UUID{organizationIds[0]}
		}

		// Check if the organization ID from the header is in the user's allowed organizations
		found := false
		for _, id := range organizationIds {
			if id == orgId {
				found = true
				selectedOrgIds = []uuid.UUID{id}
				break
			}
		}

		// If the organization ID from the header is not found, proceed with empty list
		if !found {
			logger.Info("organization id in header not found in user's organizations, proceeding with empty list", zap.String("org_id", orgId.String()))
			selectedOrgIds = []uuid.UUID{organizationIds[0]}
		}
	} else if len(organizationIds) > 0 {
		// If no header is provided and user has organizations, use the first organization
		selectedOrgIds = []uuid.UUID{organizationIds[0]}
	} else {
		// If user has no organizations, proceed with empty list
		logger.Info("user has no organizations, proceeding with empty list")
		selectedOrgIds = []uuid.UUID{}
	}

	return selectedOrgIds
}
