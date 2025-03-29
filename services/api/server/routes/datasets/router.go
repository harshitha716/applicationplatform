package datasets

import (
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	"github.com/Zampfi/application-platform/services/api/core/fileimports"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/server/middleware"
	"github.com/gin-gonic/gin"
)

// accepts a gin engine and registers all the endpoitns for the auth service at /auth
func RegisterDatasetRoutes(e *gin.RouterGroup, serverCfg *serverconfig.ServerConfig, datasetService datasetservice.DatasetService) error {

	fileUploadService := fileimports.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
	err := registerRoutes(e, datasetService, serverCfg.Store, fileUploadService)
	if err != nil {
		return err
	}

	return nil

}

func registerRoutes(e *gin.RouterGroup, datasetService datasetservice.DatasetService, store store.Store, fileUploadService fileimports.FileImportService) error {
	datasetGroup := e.Group("/datasets")
	datasetGroup.Use(middleware.ValidateDatasetAccess(store))
	{
		datasetGroup.GET("/:datasetId/filter-config", func(c *gin.Context) {
			GetFilterConfig(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/data", func(c *gin.Context) {
			GetData(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/drill-down/row/:rowUUID", func(c *gin.Context) {
			GetRowDetails(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/audiences", func(c *gin.Context) {
			GetDatasetAudiences(c, datasetService)
		})
		datasetGroup.POST("/:datasetId/audiences", func(c *gin.Context) {
			addDatasetAudiences(c, datasetService)
		})
		datasetGroup.PATCH("/:datasetId/audiences", func(c *gin.Context) {
			updateDatasetAudience(c, datasetService)
		})
		datasetGroup.DELETE("/:datasetId/audiences", func(c *gin.Context) {
			deleteDatasetAudience(c, datasetService)
		})

		datasetGroup.PATCH("/:datasetId/update", func(c *gin.Context) {
			UpdateDataset(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/actions", func(c *gin.Context) {
			GetDatasetActions(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/export", func(c *gin.Context) {
			CreateDatasetExportAction(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/signed-url/:workflowId", func(c *gin.Context) {
			GetDownloadableDataExportUrl(c, datasetService)
		})

		datasetGroup.GET("/:datasetId/file-imports/history", func(c *gin.Context) {
			GetFileImportsHistory(c, datasetService, fileUploadService)
		})

		datasetGroup.GET("/:datasetId/display-config", func(c *gin.Context) {
			GetDatasetDisplayConfig(c, datasetService)
		})

		datasetGroup.POST("/:datasetId/display-config", func(c *gin.Context) {
			SetDatasetDisplayConfig(c, datasetService)
		})

	}

	datasetCRUDGroup := e.Group("/datasets")
	datasetCRUDGroup.Use(middleware.ValidateUser())
	{
		datasetCRUDGroup.GET("/listing", func(c *gin.Context) {
			GetDatasetListing(c, datasetService)
		})
		datasetCRUDGroup.POST("/register", func(c *gin.Context) {
			RegisterDataset(c, datasetService)
		})
		datasetCRUDGroup.POST("/copy", func(c *gin.Context) {
			CopyDataset(c, datasetService)
		})
		datasetCRUDGroup.POST("/jobs/register", func(c *gin.Context) {
			RegisterDatasetJob(c, datasetService)
		})
		datasetCRUDGroup.POST("/templates/upsert", func(c *gin.Context) {
			UpsertTemplate(c, datasetService)
		})
		datasetCRUDGroup.GET("/rules/listing", func(c *gin.Context) {
			GetRulesByDatasetColumns(c, datasetService)
		})

		datasetCRUDGroup.GET("/rules/ids", func(c *gin.Context) {
			GetRulesByIds(c, datasetService)
		})
		datasetCRUDGroup.POST("/file-imports/init", func(c *gin.Context) {
			InitDatasetFileImport(c, datasetService, fileUploadService)
		})
		datasetCRUDGroup.POST("/file-imports/:fileUploadId", func(c *gin.Context) {
			AckDatasetFileImport(c, datasetService, fileUploadService)
		})
		datasetCRUDGroup.GET("/file-imports/:fileUploadId/preview", func(c *gin.Context) {
			GetDatasetImportPreview(c, datasetService, fileUploadService)
		})
		datasetCRUDGroup.POST("/file-imports/:fileUploadId/confirm", func(c *gin.Context) {
			ConfirmDatasetImport(c, datasetService)
		})

		datasetCRUDGroup.PATCH("/rules/priority", func(c *gin.Context) {
			UpdateRulesPriority(c, datasetService)
		})
	}

	datasetAdminGroup := e.Group("/datasets")
	datasetAdminGroup.Use(middleware.ValidateDatasetAdminAccess(store))
	{
		datasetAdminGroup.POST("/:datasetId/update-data", func(c *gin.Context) {
			UpdateDatasetData(c, datasetService)
		})

		datasetAdminGroup.DELETE("/:datasetId", func(c *gin.Context) {
			DeleteDataset(c, datasetService)
		})
	}

	return nil
}
