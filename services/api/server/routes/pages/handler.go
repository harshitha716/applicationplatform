package pages

import (
	"net/http"

	"github.com/Zampfi/application-platform/services/api/core/pages"
	models "github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/server/routes/pages/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleGetPagesAll(c *gin.Context, pageService pages.PagesService) {

	_, userId, orgIds := apicontext.GetAuthFromContext(c)
	if userId == nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	if len(orgIds) == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	pages, err := pageService.GetPagesByOrganizationId(c, orgIds[0])
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, pages)
}

func handleGetPagesByID(c *gin.Context, pageService pages.PagesService) {
	pageIDStr := c.Param("pageId")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid page id"})
		return
	}

	page, err := pageService.GetPageByID(c, pageID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, page)
}

func handleGetPageAudiences(c *gin.Context, pageService pages.PagesService) {
	pageIDStr := c.Param("pageId")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid page id"})
		return
	}

	audiences, err := pageService.GetPageAudiences(c, pageID)
	if err != nil {
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(200, audiences)
}

func addPageAudiences(c *gin.Context, svc pages.PagesService) {
	pageId, err := uuid.Parse(c.Param("pageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page ID"})
		return
	}

	var payload dtos.BulkAddAudienceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bulkAddAudiencePayload pages.BulkAddPageAudiencePayload
	for _, audience := range payload.Audiences {
		bulkAddAudiencePayload.Audiences = append(bulkAddAudiencePayload.Audiences, pages.AddPageAudiencePayload{
			AudienceId:   audience.AudienceId,
			AudienceType: models.AudienceType(audience.AudienceType),
			Privilege:    models.ResourcePrivilege(audience.Role),
		})
	}

	response, errs := svc.BulkAddAudienceToPage(c, pageId, bulkAddAudiencePayload)
	if errs.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audiences": response, "audience_errors": errs.Audiences})
}

func updatePageAudience(c *gin.Context, svc pages.PagesService) {
	pageId, err := uuid.Parse(c.Param("pageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page ID"})
		return
	}

	var payload dtos.UpdateAudienceRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := svc.UpdatePageAudiencePrivilege(c, pageId, payload.AudiencId, models.ResourcePrivilege(payload.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func deletePageAudience(c *gin.Context, svc pages.PagesService) {
	pageId, err := uuid.Parse(c.Param("pageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page ID"})
		return
	}

	var payload dtos.DeleteAudienceRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = svc.RemoveAudienceFromPage(c, pageId, payload.AudiencId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func handleGetPagesByOrganizationId(c *gin.Context, svc pages.PagesService) {

	organizationId, err := uuid.Parse(c.Query("organizationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	pages, err := svc.GetPagesByOrganizationId(c, organizationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}
