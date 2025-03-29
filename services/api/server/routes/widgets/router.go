package widgets

import (
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	widgetservice "github.com/Zampfi/application-platform/services/api/core/widgets/service"
	"github.com/gin-gonic/gin"
)

func RegisterWidgetRoutes(e *gin.RouterGroup, serverCfg *serverconfig.ServerConfig, datasetService datasetservice.DatasetService) error {
	widgetService := widgetservice.NewWidgetsService(serverCfg.Store, datasetService)

	widgetGroup := e.Group("/widgets")
	{
		widgetGroup.GET("/:widgetInstanceId/data", func(c *gin.Context) {
			GetWidgetInstanceData(c, widgetService)
		})

		widgetGroup.GET("/:widgetInstanceId/instance", func(c *gin.Context) {
			GetWidgetInstance(c, widgetService)
		})
	}
	return nil
}
