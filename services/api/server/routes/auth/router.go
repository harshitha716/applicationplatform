package auth

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/auth"
	"github.com/gin-gonic/gin"
)

// accepts a gin engine and registers all the endpoitns for the auth service at /auth
func RegisterAuthRoutes(e *gin.Engine, serverConfig *serverconfig.ServerConfig) error {

	const kratosPath = "kratosPath"

	authService, err := auth.NewAuthService(serverConfig.Env.ControlPlaneAdminSecrets, serverConfig.Env.AuthBaseUrl, serverConfig.Store, serverConfig.Env.Environment)
	if err != nil {
		return fmt.Errorf("failed to initialize auth service: %w", err)
	}

	registerRoutes(e, authService)

	return nil
}

func registerRoutes(e *gin.Engine, authService auth.AuthService) {
	authGroup := e.Group("/auth")
	{
		authGroup.GET("/whoami", func(c *gin.Context) {
			getUserInfoWithOrganizations(c, authService)
		})

		authGroup.POST("/login/flow/create", func(c *gin.Context) {
			getLoginFlow(c, authService)
		})

		authGroup.POST("/internal/webhook", func(c *gin.Context) {
			handleKratosAfterRegistrationWebhookEvent(c, authService)
		})

		authGroup.Any("relay/*kratosPath", func(c *gin.Context) {
			handleKratosProxyRequest(c, authService, "kratosPath")
		})
	}
}
