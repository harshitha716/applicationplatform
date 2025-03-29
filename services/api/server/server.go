package server

import (
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/server/middleware"
	"github.com/Zampfi/application-platform/services/api/server/routes/admin"
	authroutes "github.com/Zampfi/application-platform/services/api/server/routes/auth"
	"github.com/Zampfi/application-platform/services/api/server/routes/connections"
	"github.com/Zampfi/application-platform/services/api/server/routes/datasets"
	healthroute "github.com/Zampfi/application-platform/services/api/server/routes/health"
	"github.com/Zampfi/application-platform/services/api/server/routes/organizations"
	"github.com/Zampfi/application-platform/services/api/server/routes/pages"

	// 	organizationRouter "github.com/Zampfi/application-platform/services/api/services/organizations/router"
	dataplatformservice "github.com/Zampfi/application-platform/services/api/core/dataplatform"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	ruleservice "github.com/Zampfi/application-platform/services/api/core/rules/service"
	querybuilderservice "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	connectivity_catalog "github.com/Zampfi/application-platform/services/api/server/routes/connectivity_catalog"
	dataplatformwebhooks "github.com/Zampfi/application-platform/services/api/server/routes/dataplatformwebhooks"
	pantheonroutes "github.com/Zampfi/application-platform/services/api/server/routes/pantheon"
	widgets "github.com/Zampfi/application-platform/services/api/server/routes/widgets"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter(serverCfg *serverconfig.ServerConfig, logger *zap.Logger, queryBuilderService querybuilderservice.QueryBuilder, dataplatformService dataplatformservice.DataPlatformService, ruleService ruleservice.RuleService, datasetService datasetservice.DatasetService) (*gin.Engine, error) {

	r := gin.Default()

	r.Use(middleware.GetPanicRecoveryMiddleware())

	// initialize trace ID generation middleware
	r.Use(middleware.GetTraceGenMiddleware())

	// initialize context logger middleware
	r.Use(middleware.GetContextLoggerMiddleware(logger))

	// initialize request logging middleware
	r.Use(middleware.GetLoggingMiddleware(logger))

	// initialize CORS middleware
	r.Use(middleware.GetCORSMiddleware(serverCfg.Env.AllowedCORSOrigins))

	// healthz route
	healthroute.RegisterHealthcheckRoute(r)

	// register auth routes
	err := authroutes.RegisterAuthRoutes(r, serverCfg)
	if err != nil {
		logger.Error("failed to register auth routes", zap.String("error", err.Error()))
		return nil, err
	}

	// register admin routes
	admin.RegisterAdminRoutes(r, serverCfg)

	// initialize authenticated routes
	authenticatedRoutes := r.Group("/")

	// initialize auth middleware for authenticated routes
	authMiddleware, err := middleware.GetAuthMiddleware(serverCfg)
	if err != nil {
		logger.Error("failed to get auth middleware", zap.String("error", err.Error()))
		return nil, err
	}
	authenticatedRoutes.Use(authMiddleware)

	// register pantheon page routes
	pantheonroutes.RegisterPantheonRoutes(authenticatedRoutes, serverCfg)

	// register organization routes
	organizations.RegisterOrganizationRoutes(authenticatedRoutes, serverCfg)

	// register widget routes

	// register dataset routes
	err = datasets.RegisterDatasetRoutes(authenticatedRoutes, serverCfg, datasetService)
	if err != nil {
		logger.Error("failed to register dataset routes", zap.String("error", err.Error()))
		return nil, err
	}

	err = widgets.RegisterWidgetRoutes(authenticatedRoutes, serverCfg, datasetService)
	if err != nil {
		logger.Error("failed to register widget routes", zap.String("error", err.Error()))
		return nil, err
	}

	// register webhooks routes
	webhooksGroup := r.Group("/webhooks")
	userName := serverCfg.DataPlatformConfig.ActionsConfig.WebhookConfig.UserName
	password := serverCfg.DataPlatformConfig.ActionsConfig.WebhookConfig.Password
	webhooksGroup.Use(middleware.BasicAuthMiddleware(userName, password))
	dataplatformwebhooks.RegisterWebhooksRoutes(webhooksGroup, dataplatformService, datasetService)

	// register pages routes
	pages.RegisterPagesRoutes(authenticatedRoutes, datasetService, serverCfg)

	// register connectors routes
	connectivity_catalog.RegisterCatalogRoutes(authenticatedRoutes, serverCfg)

	// register connections routes
	connections.RegisterConnectionRoutes(authenticatedRoutes, serverCfg)

	return r, nil
}

func RunServer(serverCfg *serverconfig.ServerConfig, logger *zap.Logger, queryBuilderService querybuilderservice.QueryBuilder, dataplatformService dataplatformservice.DataPlatformService, ruleService ruleservice.RuleService, datasetService datasetservice.DatasetService) {
	router, err := setupRouter(serverCfg, logger, queryBuilderService, dataplatformService, ruleService, datasetService)
	if err != nil {
		logger.Error("failed to setup router", zap.String("error", err.Error()))
		panic(err)
	}
	router.Run(fmt.Sprintf(":%s", serverCfg.Env.Port)) // listen and serve on
}
