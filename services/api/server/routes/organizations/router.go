package organizations

import (
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/organizations"
	"github.com/gin-gonic/gin"
)

// accepts a gin engine and registers all the endpoitns for the auth service at /auth
func RegisterOrganizationRoutes(e *gin.RouterGroup, serverCfg *serverconfig.ServerConfig) {

	orgService := organizations.NewOrganizationService(serverCfg)

	registerRoutes(e, orgService)
}

func registerRoutes(e *gin.RouterGroup, orgService organizations.OrganizationService) {

	orgGroup := e.Group("/organizations")
	{
		orgGroup.GET("/", func(c *gin.Context) {
			getOrganizations(c, orgService)
		})

		orgGroup.GET("/membership-requests", func(c *gin.Context) {
			getOrganizationMembershipRequests(c, orgService)
		})

		orgGroup.GET("/:orgId/audiences", func(c *gin.Context) {
			getOrganizationAudiences(c, orgService)
		})

		orgGroup.PATCH("/:orgId/audiences", func(c *gin.Context) {
			updateMemberRole(c, orgService)
		})

		orgGroup.DELETE("/:orgId/audiences", func(c *gin.Context) {
			removeOrganizationMember(c, orgService)
		})

		orgGroup.POST("/:orgId/audiences/invitations", func(c *gin.Context) {
			bulkInviteMembers(c, orgService)
		})
		orgGroup.GET("/:orgId/audiences/invitations", func(c *gin.Context) {
			getOrganizationInvitations(c, orgService)
		})

		orgGroup.GET("/:orgId/audiences/requests", func(c *gin.Context) {
			getOrganizationMembershipRequestsByOrganizationId(c, orgService)
		})

		orgGroup.PATCH("/:orgId/audiences/requests/approve", func(c *gin.Context) {
			approveOrganizationMembershipRequest(c, orgService)
		})

		orgGroup.GET("/:orgId/teams", func(c *gin.Context) {
			getTeams(c, orgService)
		})

		orgGroup.POST("/:orgId/teams", func(c *gin.Context) {
			createTeam(c, orgService)
		})

		orgGroup.GET("/:orgId/teams/:teamId", func(c *gin.Context) {
			getTeam(c, orgService)
		})

		orgGroup.POST("/:orgId/teams/:teamId/add", func(c *gin.Context) {
			addUserToTeam(c, orgService)
		})

		orgGroup.POST("/:orgId/teams/:teamId/remove", func(c *gin.Context) {
			removeUserFromTeam(c, orgService)
		})

	}
}
