package pages

import (
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	"github.com/Zampfi/application-platform/services/api/core/pages"
	"github.com/Zampfi/application-platform/services/api/pkg/cache"
	"github.com/Zampfi/application-platform/services/api/server/routes/pages/sheets"
	"github.com/gin-gonic/gin"
)

func RegisterPagesRoutes(e *gin.RouterGroup, datasetService datasetservice.DatasetService, serverCfg *serverconfig.ServerConfig) {
	pageService := pages.NewPagesService(serverCfg.Store)

	registerRoutes(e, pageService, datasetService, serverCfg.CacheClient, serverCfg)

}
func registerRoutes(e *gin.RouterGroup, pageService pages.PagesService, datasetService datasetservice.DatasetService, cacheService cache.CacheClient, serverCfg *serverconfig.ServerConfig) {
	pagesGroup := e.Group("/pages")
	{
		pagesGroup.GET("/get-pages", func(c *gin.Context) {
			handleGetPagesAll(c, pageService)
		})

		pagesGroup.GET("/get-pages-by-organization-id", func(c *gin.Context) {
			handleGetPagesByOrganizationId(c, pageService)
		})

		pagesGroup.GET("/:pageId", func(c *gin.Context) {
			handleGetPagesByID(c, pageService)
		})

		pagesGroup.GET("/:pageId/audiences", func(c *gin.Context) {
			handleGetPageAudiences(c, pageService)
		})

		pagesGroup.POST("/:pageId/audiences", func(c *gin.Context) {
			addPageAudiences(c, pageService)
		})

		pagesGroup.PATCH("/:pageId/audiences", func(c *gin.Context) {
			updatePageAudience(c, pageService)
		})

		pagesGroup.DELETE("/:pageId/audiences", func(c *gin.Context) {
			deletePageAudience(c, pageService)
		})

		sheetsGroup := pagesGroup.Group("/:pageId/sheets")

		sheets.RegisterSheetsRoutes(sheetsGroup, datasetService, cacheService, serverCfg)

	}
}
