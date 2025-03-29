package dataplatformwebhooks

import (
	dataplatformservice "github.com/Zampfi/application-platform/services/api/core/dataplatform"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	"github.com/gin-gonic/gin"
)

func RegisterWebhooksRoutes(e *gin.RouterGroup, dataplatformService dataplatformservice.DataPlatformService, datasetService datasetservice.DatasetService) error {
	e.POST("/databricks/jobs/", func(c *gin.Context) {
		HandleJobStatusUpdate(c, dataplatformService, datasetService)
	})

	return nil
}
