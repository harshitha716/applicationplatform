package organizations

import (
	"net/http"

	"github.com/Zampfi/application-platform/services/api/core/organizations"
	"github.com/Zampfi/application-platform/services/api/core/organizations/teams"
	"github.com/Zampfi/application-platform/services/api/db/models"
	dtos "github.com/Zampfi/application-platform/services/api/server/routes/organizations/dtos"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func getOrganizations(c *gin.Context, svc organizations.OrganizationService) {
	orgs, err := svc.GetOrganizations(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func getOrganizationAudiences(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	audiences, err := svc.GetOrganizationAudiences(c, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, audiences)
}
func updateMemberRole(c *gin.Context, svc organizations.OrganizationService) {

	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var requestBody dtos.UpdateMemberRoleRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	policy, err := svc.UpdateMemberRole(c, orgId, requestBody.UserId, models.ResourcePrivilege(requestBody.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, policy)
}

func bulkInviteMembers(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var requestBody dtos.BulkInvitationPayload
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	bulkInvitationPayload := organizations.BulkInvitationPayload{}

	for _, invitation := range requestBody.Invitations {
		bulkInvitationPayload.Invitations = append(bulkInvitationPayload.Invitations, organizations.InvitationPayload{
			Privilege: invitation.Role,
			Email:     invitation.Email,
		})
	}

	invitations, invitationError := svc.BulkInviteMembers(c, orgId, bulkInvitationPayload)
	if invitationError.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invitations": invitations, "invitation_errors": invitationError.Invitations})
}

func getOrganizationInvitations(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	invitations, err := svc.GetAllOrganizationInvitations(c, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, invitations)
}

func removeOrganizationMember(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var deleteMemberRequest dtos.DeleteMemberRequest
	if err := c.BindJSON(&deleteMemberRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = svc.RemoveOrganizationMember(c, orgId, deleteMemberRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}

func getOrganizationMembershipRequests(c *gin.Context, svc organizations.OrganizationService) {
	requests, err := svc.GetOrganizationMembershipRequestsAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func getOrganizationMembershipRequestsByOrganizationId(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	requests, err := svc.GetOrganizationMembershipRequestsByOrganizationId(c, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func approveOrganizationMembershipRequest(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var requestBody dtos.ApproveOrganizationMembershipRequestRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	request, err := svc.ApprovePendingOrganizationMembershipRequest(c, orgId, requestBody.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, request)
}

func getTeams(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	teams, err := svc.TeamService().GetTeamsByOrganizationID(c, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func createTeam(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var requestBody teams.CreateTeamPayload
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	team, err := svc.TeamService().CreateTeam(c, orgId, requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, team)
}

func getTeam(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")
	teamIdStr := c.Param("teamId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	team, err := svc.TeamService().GetTeamById(c, orgId, teamId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func addUserToTeam(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")
	teamIdStr := c.Param("teamId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	var addUserToTeamRequest teams.AddUserToTeamPayload
	if err := c.BindJSON(&addUserToTeamRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	teamMembership, err := svc.TeamService().AddUserToTeam(c, orgId, teamId, addUserToTeamRequest.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teamMembership)

}

func removeUserFromTeam(c *gin.Context, svc organizations.OrganizationService) {
	orgIdStr := c.Param("orgId")
	teamIdStr := c.Param("teamId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	var removeUserFromTeamRequest teams.RemoveUserFromTeamPayload
	if err := c.BindJSON(&removeUserFromTeamRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = svc.TeamService().RemoveUserFromTeam(c, orgId, teamId, removeUserFromTeamRequest.TeamMembershipID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user removed from team successfully"})
}
