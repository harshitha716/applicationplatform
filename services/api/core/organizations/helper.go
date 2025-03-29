package organizations

import (
	"slices"

	"github.com/Zampfi/application-platform/services/api/db/models"
)

func isOrganizationPrivilege(privilege models.ResourcePrivilege) bool {
	return slices.Contains(models.OrganizationPrivileges, privilege)
}
