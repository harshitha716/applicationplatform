package helper

import (
	"net/http"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"

	"github.com/gin-gonic/gin"
)

const PROXY_USER_ID_HEADER = "X-Zamp-User-Id"
const PROXY_WORKSPACE_IDS_HEADER = "X-Zamp-Organization-Ids"
const PROXY_WORKSPACE_ID_HEADER = "X-Zamp-Organization-Id"
const PROXY_ROLE_HEADER = "X-Zamp-Role"

const ADMIN_SECRET_HEADER = "X-Zamp-Admin-Secret"

func AddAuthHeaders(headers http.Header, ginctx *gin.Context) {

	role, userId, organizations := apicontext.GetAuthFromContext(ginctx)

	headers.Add(PROXY_ROLE_HEADER, role)

	if userId != nil {
		headers.Add(PROXY_USER_ID_HEADER, userId.String())
	}

	// Add the single organization ID to the header
	if len(organizations) > 0 {
		headers.Add(PROXY_WORKSPACE_ID_HEADER, organizations[0].String())
	}

	// For backward compatibility, also add all organizations to the plural header
	for _, organization := range organizations {
		headers.Add(PROXY_WORKSPACE_IDS_HEADER, organization.String())
	}
}

func GetAdminSecretFromHeader(header http.Header) string {
	return header.Get(ADMIN_SECRET_HEADER)
}
