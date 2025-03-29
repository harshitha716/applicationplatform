package connections

import (
	config "github.com/Zampfi/application-platform/services/api/config"
	manager "github.com/Zampfi/application-platform/services/api/core/connections/managers"
	"github.com/Zampfi/application-platform/services/api/core/connections/service"
	"github.com/Zampfi/application-platform/services/api/core/organizations"
	schedules "github.com/Zampfi/application-platform/services/api/core/schedules/service"
	"github.com/gin-gonic/gin"
)

func RegisterConnectionRoutes(e *gin.RouterGroup, serverConfig *config.ServerConfig) error {
	connectionService := service.NewConnectionService(serverConfig.Store)
	scheduleService := schedules.NewScheduleService(serverConfig.Store, serverConfig.TemporalSdk)
	organizationService := organizations.NewOrganizationService(serverConfig)

	connectionManagerRegistry := manager.NewConnectionManagerRegistry()
	connectionGroup := e.Group("/connections")
	{
		connectionGroup.POST("/", func(c *gin.Context) {
			HandleCreateConnection(
				c,
				connectionService,
				serverConfig.Store,
				connectionManagerRegistry,
				scheduleService,
				organizationService,
			)
		})

		connectionGroup.GET("", func(c *gin.Context) {
			handleGetConnections(
				c,
				connectionService,
				scheduleService,
				serverConfig.Store,
			)
		})

		connectionGroup.GET("/:ConnectionID/schedules", func(c *gin.Context) {
			handleGetSchedules(
				c,
				connectionService,
				scheduleService,
			)
		})

		return nil
	}
}
