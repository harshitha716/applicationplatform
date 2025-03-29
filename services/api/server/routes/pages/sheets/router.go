package sheets

import (
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	datasetsService "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	"github.com/Zampfi/application-platform/services/api/core/sheets"
	"github.com/Zampfi/application-platform/services/api/pkg/cache"
	"github.com/gin-gonic/gin"
)

func RegisterSheetsRoutes(sheetsGroup *gin.RouterGroup, datasetService datasetsService.DatasetService, cacheService cache.CacheClient, serverCfg *serverconfig.ServerConfig) {

	sheetService := sheets.NewSheetsService(serverCfg.Store, datasetService, cacheService)

	sheetsGroup.GET("/", func(c *gin.Context) {
		handleGetSheetsAll(c, sheetService)
	})

	sheetsGroup.GET("/:sheetId", func(c *gin.Context) {
		handleGetSheetByID(c, sheetService)
	})

	sheetsGroup.GET("/:sheetId/filters", func(c *gin.Context) {
		handleGetSheetFilters(c, sheetService)
	})
}
