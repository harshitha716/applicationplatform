package connectivity_catalog

import (
	config "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/connectivity_catalog/service"
	"github.com/gin-gonic/gin"
)

func RegisterCatalogRoutes(e *gin.RouterGroup, serverConfig *config.ServerConfig) error {
	connectorService := service.NewConnectorService(serverConfig.Store)

	catalogGroup := e.Group("/connectors")

	catalogGroup.GET("", func(c *gin.Context) {
		handleGetConnectors(c, connectorService)
	})

	catalogGroup.GET("/:connectorId", func(c *gin.Context) {
		handleGetConnectorByID(c, connectorService)
	})

	return nil
}
