package auth

import (
	"net/http"
	"slices"

	auth "github.com/Zampfi/application-platform/services/api/core/auth"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/helper"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/server/routes/auth/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func handleKratosProxyRequest(c *gin.Context, authService auth.AuthService, proxyPath string) {

	logger := apicontext.GetLoggerFromCtx(c)

	// expose only self-service routes
	kratosPath := c.Param(proxyPath)
	if !authService.IsUserExposedKratosPath(kratosPath) {
		logger.Error("tried accessing internal route of kratos", zap.String("path", kratosPath))
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}

	c.Request.URL.Path = kratosPath

	kratosProxy, url := authService.GetKratosProxy()

	// Update the request to reflect the target scheme and host
	c.Request.URL.Scheme = url.Scheme
	c.Request.URL.Host = url.Host
	c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
	c.Request.Host = url.Host

	kratosProxy.ServeHTTP(c.Writer, c.Request)
}

func getUserInfoWithOrganizations(c *gin.Context, authService auth.AuthService) {

	// TODO: Dry this up -- same code used in middleware
	logger := apicontext.GetLoggerFromCtx(c)

	cookie := c.Request.Header.Get("cookie")

	sessionInfo, _, kerr := authService.ResolveSessionCookie(c, cookie)
	if kerr != nil {
		logger.Error("error resolving session cookie", zap.String("error", kerr.Message))
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userIdStr := sessionInfo.Identity.Id
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		logger.Error("error parsing user id", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	// get email from session traits
	traits := sessionInfo.Identity.Traits.(map[string]interface{})
	if traits == nil {
		logger.Error("error getting user info", zap.String("error", "traits not found in session"))
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userEmail := traits["email"].(string)
	if userEmail == "" {
		logger.Error("error getting user info", zap.String("error", "email not found in session traits"))
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	ctxWithUserId := apicontext.AddAuthToContext(c, "user", userId, []uuid.UUID{})

	orgs, err := authService.GetUserOrganizations(ctxWithUserId, userId)
	if err != nil {
		logger.Error("error getting user info", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	// Convert the organizations to UUIDs
	var orgIds []uuid.UUID
	for _, org := range orgs {
		orgIds = append(orgIds, org.ID)
	}

	// Use the helper function to select the organization
	selectedOrgIds := auth.SelectOrganization(c.Request.Header, orgIds, logger)

	filteredOrgs := []models.Organization{}
	for _, org := range orgs {
		if slices.Contains(selectedOrgIds, org.ID) {
			filteredOrgs = append(filteredOrgs, org)
		}
	}

	if len(filteredOrgs) == 0 && len(orgs) > 0 {
		filteredOrgs = orgs
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userIdStr,
		"user_email": userEmail,
		"orgs":       filteredOrgs,
	})
}

func getLoginFlow(c *gin.Context, authService auth.AuthService) {

	request := dtos.GetLoginFlowRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginFlow, resp, err := authService.GetAuthFlowForUser(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if resp != nil {
		helper.ForwardResponseHeaders(resp, c)
	}

	c.JSON(http.StatusOK, loginFlow)

}

func handleKratosAfterRegistrationWebhookEvent(c *gin.Context, authService auth.AuthService) {

	request := dtos.AfterRegistrationWebhookRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	adminSecret := helper.GetAdminSecretFromHeader(c.Request.Header)

	err = authService.HandleNewUserCreated(c, adminSecret, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})

}
