package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	auth "github.com/Zampfi/application-platform/services/api/core/auth"
	authservice "github.com/Zampfi/application-platform/services/api/core/auth"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	apictx "github.com/Zampfi/application-platform/services/api/helper/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CORS middleware
func GetAuthMiddleware(serverCfg *serverconfig.ServerConfig) (gin.HandlerFunc, error) {

	authSvc, err := authservice.NewAuthService(serverCfg.Env.ControlPlaneAdminSecrets, serverCfg.Env.AuthBaseUrl, serverCfg.Store, serverCfg.Env.Environment)
	if err != nil {
		return nil, err
	}

	return func(c *gin.Context) {
		role, userId, organizationIds := authSvc.ResolveAdminInfo(c, c.Request.Header)

		if role == "admin" {
			apictx.AddAuthToGinContext(c, role, userId, organizationIds)
		} else {
			userMiddleWare(authSvc)(c)
		}

		c.Next()
	}, nil
}

func userMiddleWare(authService authservice.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

		logger := apictx.GetLoggerFromCtx(c)

		reqCookie := c.Request.Header.Get("Cookie")

		// resolve user ID from session
		sessionInfo, _, err := authService.ResolveSessionCookie(c, reqCookie)
		if err != nil {
			logger.Error("failed to resolve session cookie", zap.String("error", err.Message))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		userID, parseErr := uuid.Parse(sessionInfo.Identity.Id)
		if parseErr != nil {
			logger.Error("failed to parse user id", zap.String("error", parseErr.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		// extract traits from session as map[string]string
		traits := sessionInfo.Identity.Traits.(map[string]interface{})
		if traits == nil {
			logger.Error("failed to resolve session cookie", zap.String("error", "traits not found in session"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		// extract email from session traits
		userEmail := traits["email"].(string)
		if userEmail == "" {
			logger.Error("failed to resolve session cookie", zap.String("error", "email not found in session traits"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		ctxWithUserId := apicontext.AddAuthToContext(c, "user", userID, []uuid.UUID{})
		allowedOrganizations, oerr := authService.GetUserOrganizations(ctxWithUserId, userID)
		if oerr != nil {
			logger.Error("failed to get user info", zap.Error(oerr))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		var organizationIds []uuid.UUID
		for _, organization := range allowedOrganizations {
			organizationIds = append(organizationIds, organization.ID)
		}

		// Use the helper function to select the organization
		selectedOrgIds := auth.SelectOrganization(c.Request.Header, organizationIds, logger)

		if len(selectedOrgIds) > 0 {
			logger.Info("adding auth to context", zap.String("user_id", userID.String()), zap.String("organization_id", selectedOrgIds[0].String()))
		} else {
			logger.Info("adding auth to context with empty organization list", zap.String("user_id", userID.String()))
		}
		apictx.AddAuthToGinContext(c, "user", userID, selectedOrgIds)
		
		// Add audit information to context
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()
		apictx.AddAuditInfoToGinContext(c, userEmail, ipAddress, userAgent)

		c.Next()
	}
}

func BasicAuthMiddleware(requiredUsername, requiredPassword string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := apictx.GetLoggerFromCtx(c)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("AUTHORIZATION HEADER REQUIRED")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			logger.Error("INVALID AUTHORIZATION HEADER FORMAT")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			logger.Error("INVALID BASE64 ENCODING")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Base64 encoding"})
			return
		}

		credentials := strings.SplitN(string(decodedBytes), ":", 2)
		if len(credentials) != 2 {
			logger.Error("INVALID CREDENTIALS FORMAT")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			return
		}
		username, password := credentials[0], credentials[1]

		if username != requiredUsername || password != requiredPassword {
			logger.Error("INVALID CREDENTIALS")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		c.Next()
	}
}

func ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := apicontext.GetLoggerFromCtx(c)
		_, userId, _ := apicontext.GetAuthFromContext(c)

		if userId == nil {
			logger.Error("user ID not found", zap.String("error", "user ID not found"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
