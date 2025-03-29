package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/auth"
	manager "github.com/Zampfi/application-platform/services/api/core/connections/managers"
	"github.com/Zampfi/application-platform/services/api/core/connections/service"
	catalogservice "github.com/Zampfi/application-platform/services/api/core/connectivity_catalog/service"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform"
	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	fileimportsservice "github.com/Zampfi/application-platform/services/api/core/fileimports"
	"github.com/Zampfi/application-platform/services/api/core/organizations"
	orgservice "github.com/Zampfi/application-platform/services/api/core/organizations"
	"github.com/Zampfi/application-platform/services/api/core/pages"
	rules "github.com/Zampfi/application-platform/services/api/core/rules/service"
	schedules "github.com/Zampfi/application-platform/services/api/core/schedules/service"
	"github.com/Zampfi/application-platform/services/api/core/sheets"
	sheetmodels "github.com/Zampfi/application-platform/services/api/core/sheets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	widgets "github.com/Zampfi/application-platform/services/api/core/widgets/service"
	models "github.com/Zampfi/application-platform/services/api/db/models"
	store "github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/helper"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	cloudservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/service"
	querybuilder "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	"github.com/Zampfi/application-platform/services/api/server/routes/admin/templates"
	connections "github.com/Zampfi/application-platform/services/api/server/routes/connections"
	"github.com/Zampfi/application-platform/services/api/server/routes/datasets/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RegisterAdminRoutes(router *gin.Engine, serverCfg *serverconfig.ServerConfig) {

	templateLoader := templates.InitTemplateLoader(serverCfg.Env.AdminHTMLTemplatesPath)
	authSvc, err := auth.NewAuthService(serverCfg.Env.ControlPlaneAdminSecrets, serverCfg.Env.AuthBaseUrl, serverCfg.Store, serverCfg.Env.Environment)
	if err != nil {
		panic(err)
	}

	registerRoutes(router, serverCfg, templateLoader, authSvc)
}

type adminRouteController struct {
	routeGroup     *gin.RouterGroup
	templateLoader templates.TemplateLoader
	authSvc        auth.AuthService
	serverCfg      *serverconfig.ServerConfig
	routes         []struct {
		Title string `json:"title"`
		Path  string `json:"path"`
	}
}

func validateProvider(provider dataplatformdataconstants.Provider) error {
	switch provider {
	case dataplatformdataconstants.ProviderDatabricks, dataplatformdataconstants.ProviderPinot:
		return nil
	default:
		return errors.New("invalid provider")
	}
}

func (c *adminRouteController) registerAdminRoute(path string, routeName string, data any, handler gin.HandlerFunc) {
	c.routeGroup.Any(path, func(ctx *gin.Context) {
		if ctx.Request.Method == "GET" {
			templateName := "base-form.html"
			if path == "/login" {
				templateName = "admin-home.html"
			}
			c.templateLoader.ExecuteTemplate(ctx.Writer, templateName, templates.TemplateData{
				Title:       routeName,
				Environment: c.serverCfg.Env.Environment,
				Route:       fmt.Sprintf("%s%s", "/admin", path),
				Data:        data,
			})
			ctx.Status(200)
			ctx.Writer.Flush()
			ctx.Abort()
			return
		} else if ctx.Request.Method == "POST" {
			role, userId, orgIds := c.authSvc.ResolveAdminInfo(ctx, ctx.Request.Header)
			if role != "admin" {
				ctx.JSON(403, gin.H{"error": "forbidden"})
				return
			}
			apicontext.AddAuthToGinContext(ctx, role, userId, orgIds)
			handler(ctx)
		} else {
			ctx.JSON(405, gin.H{"error": "method not allowed"})
			return
		}
	})
	c.routes = append(c.routes, struct {
		Title string `json:"title"`
		Path  string `json:"path"`
	}{
		Title: routeName,
		Path:  fmt.Sprintf("%s%s", "/admin", path),
	})
}

func registerRoutes(router *gin.Engine, serverCfg *serverconfig.ServerConfig, templateLoader templates.TemplateLoader, authSvc auth.AuthService) {

	adminRouterGroup := router.Group("/admin")

	adminRouteController := &adminRouteController{
		routeGroup:     adminRouterGroup,
		templateLoader: templateLoader,
		authSvc:        authSvc,
		serverCfg:      serverCfg,
		routes: []struct {
			Title string `json:"title"`
			Path  string `json:"path"`
		}{},
	}

	adminRouteController.registerAdminRoute("/login", "Login", struct{}{}, func(ctx *gin.Context) {
		ctx.JSON(405, gin.H{"error": "method not allowed"})
	})

	adminRouterGroup.GET("/api-catalog", func(ctx *gin.Context) {

		role, _, _ := authSvc.ResolveAdminInfo(ctx, ctx.Request.Header)
		if role != "admin" {
			ctx.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		ctx.JSON(200, adminRouteController.routes)
	})

	// creating a page
	var createPagePayload struct {
		PageName        string `form:"label=Page Name"`
		PageDescription string `form:"label=Page Description"`
	}
	adminRouteController.registerAdminRoute("/create-page", "Create Page", createPagePayload, func(ctx *gin.Context) {

		// get page payload from the form
		pagePayload := pages.CreatePagePayload{}
		err := ctx.ShouldBind(&pagePayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		pageSvc := pages.NewPagesService(serverCfg.Store)
		page, err := pageSvc.CreatePage(ctx, pagePayload)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "page created", "data": page})
	})

	// creating a widget
	var createWidgetInstancePayload struct {
		WidgetType    string `form:"label=Widget Type"`
		SheetId       string `form:"label=Sheet ID"`
		Title         string `form:"label=Title"`
		DataMappings  string `form:"label=Data Mappings"`
		DisplayConfig string `form:"label=Display Config"`
	}
	adminRouteController.registerAdminRoute("/create-widget-instance", "Create Widget Instance", createWidgetInstancePayload, func(ctx *gin.Context) {
		widgetPayload := &widgetmodels.CreateWidgetInstancePayload{}
		err := ctx.ShouldBind(&widgetPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		rulesService := rules.NewRuleService(serverCfg.Store)
		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			panic(err)
		}
		if cloudService == nil {
			panic("CloudService is nil")
		}
		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, querybuilder.NewQueryBuilder(), dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)
		widgetSvc := widgets.NewWidgetsService(serverCfg.Store, datasetSvc)

		widgetInstance, err := widgetPayload.ToModel()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		widgetInstance, err = widgetSvc.CreateWidgetInstance(ctx, *widgetInstance)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "widget instance created", "data": widgetInstance})
	})

	// updating a widget
	var updateWidgetInstancePayload struct {
		WidgetInstanceID string  `form:"label=Widget Instance ID"`
		WidgetType       *string `form:"label=Widget Type" form:",omitempty"`
		SheetId          *string `form:"label=Sheet ID" form:",omitempty"`
		Title            *string `form:"label=Title" form:",omitempty"`
		DataMappings     *string `form:"label=Data Mappings" form:",omitempty"`
		DisplayConfig    *string `form:"label=Display Config" form:",omitempty"`
	}
	adminRouteController.registerAdminRoute("/update-widget-instance", "Update Widget Instance", updateWidgetInstancePayload, func(ctx *gin.Context) {
		widgetPayload := &widgetmodels.UpdateWidgetInstancePayload{}
		err := ctx.ShouldBind(&widgetPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		rulesService := rules.NewRuleService(serverCfg.Store)
		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			panic(err)
		}
		if cloudService == nil {
			panic("CloudService is nil")
		}
		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, querybuilder.NewQueryBuilder(), dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)
		widgetSvc := widgets.NewWidgetsService(serverCfg.Store, datasetSvc)

		widgetInstance, err := widgetPayload.ToModel()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		widgetInstance, err = widgetSvc.UpdateWidgetInstance(ctx, *widgetInstance)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "widget instance updated", "data": widgetInstance})
	})

	// create a sheet
	var createSheetPayload struct {
		Name        string  `form:"name"`
		Description *string `form:"description"`
		PageId      string  `form:"page_id"`
		SheetConfig string  `form:"sheet_config"`
	}
	adminRouteController.registerAdminRoute("/create-sheet", "Create Sheet", createSheetPayload, func(ctx *gin.Context) {
		sheetPayload := &sheetmodels.CreateSheetPayload{}
		err := ctx.ShouldBind(&sheetPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		sheet, err := sheetPayload.ToModel()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		rulesService := rules.NewRuleService(serverCfg.Store)
		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			panic(err)
		}
		if cloudService == nil {
			panic("CloudService is nil")
		}
		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, querybuilder.NewQueryBuilder(), dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)
		sheetSvc := sheets.NewSheetsService(serverCfg.Store, datasetSvc, serverCfg.CacheClient)

		sheetRet, err := sheetSvc.CreateSheet(ctx, *sheet)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "sheet created", "data": sheetRet})
	})

	var updateSheetPayload struct {
		SheetId     string  `form:"sheet_id" form:"label=Sheet ID"`
		Name        *string `form:"name" form:"label=Name" form:",omitempty"`
		Description *string `form:"description" form:"label=Description" form:",omitempty"`
		PageId      *string `form:"page_id" form:"label=Page ID" form:",omitempty"`
		SheetConfig *string `form:"sheet_config" form:"label=Sheet Config" form:",omitempty"`
	}
	adminRouteController.registerAdminRoute("/update-sheet", "Update Sheet", updateSheetPayload, func(ctx *gin.Context) {
		sheetPayload := &sheetmodels.UpdateSheetPayload{}
		err := ctx.ShouldBind(&sheetPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		sheet, err := sheetPayload.ToModel()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		rulesService := rules.NewRuleService(serverCfg.Store)
		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			panic(err)
		}
		if cloudService == nil {
			panic("CloudService is nil")
		}
		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, querybuilder.NewQueryBuilder(), dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)
		sheetSvc := sheets.NewSheetsService(serverCfg.Store, datasetSvc, serverCfg.CacheClient)

		sheetRet, err := sheetSvc.UpdateSheet(ctx, sheet)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "sheet updated", "data": sheetRet})
	})

	// register a dataset
	var registerDatasetPayload dtos.RegisterDatasetRequest
	adminRouteController.registerAdminRoute("/register-dataset", "Register Dataset", registerDatasetPayload, func(ctx *gin.Context) {

		// extract the payload from the form
		err := ctx.ShouldBind(&registerDatasetPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		rulesService := rules.NewRuleService(serverCfg.Store)
		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
		queryBuilderService := querybuilder.NewQueryBuilder()
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			panic(err)
		}
		if cloudService == nil {
			panic("CloudService is nil")
		}
		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, queryBuilderService, dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)

		role, userId, orgIds := apicontext.GetAuthFromContext(ctx)
		if role != "admin" {
			ctx.JSON(403, gin.H{"error": "forbidden"})
			return
		}

		if registerDatasetPayload.Provider == "" {
			registerDatasetPayload.Provider = dataplatformdataconstants.ProviderDatabricks
		}

		err = validateProvider(registerDatasetPayload.Provider)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		actionId, datasetId, err := datasetSvc.RegisterDataset(ctx, orgIds[0], *userId, registerDatasetPayload.ToModel())
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "dataset registered", "actionId": actionId, "datasetId": datasetId})

	})

	// copy a dataset
	var copyDatasetPayload datasetmodels.CopyDatasetParams
	adminRouteController.registerAdminRoute("/copy-dataset", "Copy Dataset", copyDatasetPayload, func(ctx *gin.Context) {

		// extract the payload from the form
		err := ctx.ShouldBind(&copyDatasetPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		rulesService := rules.NewRuleService(serverCfg.Store)
		queryBuilderService := querybuilder.NewQueryBuilder()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if cloudService == nil {
			ctx.JSON(500, gin.H{"error": "CloudService is nil"})
			return
		}

		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)

		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, queryBuilderService, dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)

		role, userId, orgIds := apicontext.GetAuthFromContext(ctx)
		if role != "admin" {
			ctx.JSON(403, gin.H{"error": "forbidden"})
			return
		}

		actionId, datasetId, err := datasetSvc.CopyDataset(ctx, orgIds[0], *userId, copyDatasetPayload)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "dataset copied", "actionId": actionId, "datasetId": datasetId})

	})

	// create a new organization
	var createOrganizationPayload struct {
		OrganizationName        string `form:"label=Organization Name"`
		OrganizationDescription string `form:"label=Organization Description"`
		AdminEmail              string `form:"label=Admin Email"`
		AdminPassword           string `form:"label=Password for admin account"`
		SSODomain               string `form:"label=SSO Domain"`
		SSOProviderName         string `form:"label=SSO Provider Name"`
		SSOProviderID           string `form:"label=SSO Provider ID"`
	}
	adminRouteController.registerAdminRoute("/create-organization", "Create Organization", createOrganizationPayload, func(ctx *gin.Context) {

		logger := apicontext.GetLoggerFromCtx(ctx)

		err := ctx.ShouldBind(&createOrganizationPayload)
		if err != nil {
			logger.Error("bad payload", zap.String("error", err.Error()))
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		if !helper.IsValidEmail(createOrganizationPayload.AdminEmail) {
			logger.Error("invalid email", zap.String("email", createOrganizationPayload.AdminEmail))
			ctx.JSON(400, gin.H{"error": "invalid email"})
			return
		}

		domainFromEmail := helper.GetDomainFromEmail(createOrganizationPayload.AdminEmail)
		if domainFromEmail != createOrganizationPayload.SSODomain {
			ctx.JSON(400, gin.H{"error": "invalid domain"})
			return
		}

		// check if the user already exists
		// signup a new user as admin
		user, err := adminRouteController.authSvc.SignupUserAsAdmin(ctx, createOrganizationPayload.AdminEmail, createOrganizationPayload.AdminPassword)
		if err != nil {
			logger.Error("failed to signup user as admin", zap.String("error", err.Error()))
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		logger.Info("user created", zap.Any("user", user))

		logger.Info("creating organization")

		orgService := orgservice.NewOrganizationService(adminRouteController.serverCfg)

		logger.Info("creating organization", zap.String("name", createOrganizationPayload.OrganizationName))
		org, err := orgService.CreateOrganization(ctx, createOrganizationPayload.OrganizationName, &createOrganizationPayload.OrganizationDescription, user.ID)
		if err != nil {
			logger.Error("failed to create organization", zap.String("error", err.Error()))
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		logger.Info("creating sso config")
		_, err = serverCfg.Store.CreateSSOConfig(ctx, org.ID, createOrganizationPayload.SSOProviderID, createOrganizationPayload.SSOProviderName, json.RawMessage("{}"), createOrganizationPayload.SSODomain)
		if err != nil {
			logger.Error("failed to create sso config", zap.String("error", err.Error()))
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		logger.Info("organization created", zap.Any("organization", org))

		ctx.JSON(200, gin.H{"message": "organization created", "organization": org, "user": user})
	})

	// create a connection
	var createConnectionPayload struct {
		DisplayName      string `form:"label=Display Name"`
		ConnectorName    string `form:"label=Connector Name"`
		ConnectorID      string `form:"label=Connector ID"`
		ConnectionConfig string `form:"label=Connection Config (JSON)"`
	}
	adminRouteController.registerAdminRoute("/create-connection", "Create Connection", createConnectionPayload, func(ctx *gin.Context) {
		connectionService := service.NewConnectionService(adminRouteController.serverCfg.Store)
		scheduleService := schedules.NewScheduleService(adminRouteController.serverCfg.Store, adminRouteController.serverCfg.TemporalSdk)
		organizationService := orgservice.NewOrganizationService(adminRouteController.serverCfg)

		connectionManagerRegistry := manager.NewConnectionManagerRegistry()

		err := ctx.ShouldBind(&createConnectionPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		var connectionConfig map[string]interface{}
		err = json.Unmarshal([]byte(createConnectionPayload.ConnectionConfig), &connectionConfig)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid connection config JSON"})
			return
		}

		ctx.Request.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{
			"display_name": "%s",
			"connector_name": "%s", 
			"connector_id": "%s",
			"connection_config": %s
		}`, createConnectionPayload.DisplayName,
			createConnectionPayload.ConnectorName,
			createConnectionPayload.ConnectorID,
			createConnectionPayload.ConnectionConfig)))

		connections.HandleCreateConnection(
			ctx,
			connectionService,
			adminRouteController.serverCfg.Store,
			connectionManagerRegistry,
			scheduleService,
			organizationService,
		)
	})

	// create a connector
	type createConnectorPayload struct {
		Name           string `form:"label=Name"`
		Description    string `form:"label=Description"`
		DisplayName    string `form:"label=Display Name"`
		Documentation  string `form:"label=Documentation"`
		LogoURL        string `form:"label=Logo URL"`
		ConfigTemplate string `form:"label=Config Template (JSON)"`
		Category       string `form:"label=Category"`
		Status         string `form:"label=Status"`
	}
	adminRouteController.registerAdminRoute("/create-connector", "Create Connector", createConnectorPayload{}, func(ctx *gin.Context) {
		var payload createConnectorPayload
		err := ctx.ShouldBind(&payload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}
		// Parse the config template JSON
		var configTemplate json.RawMessage
		err = json.Unmarshal([]byte(payload.ConfigTemplate), &configTemplate)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid config template JSON"})
			return
		}

		// Create the connector
		connector := &models.Connector{
			ID:             uuid.New(),
			Name:           payload.Name,
			Description:    payload.Description,
			DisplayName:    payload.DisplayName,
			Documentation:  payload.Documentation,
			LogoURL:        payload.LogoURL,
			ConfigTemplate: configTemplate,
			Category:       payload.Category,
			Status:         payload.Status,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			IsDeleted:      false,
		}

		var createdConnector *models.Connector
		err = adminRouteController.serverCfg.Store.WithTx(ctx, func(store store.Store) error {
			// Create the connector
			catalogService := catalogservice.NewConnectorService(store)
			err := catalogService.CreateConnector(ctx, connector)
			if err != nil {
				return err
			}
			// Get the created connector
			createdConnector, err = catalogService.GetConnectorByID(ctx, connector.ID)
			return err
		})

		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "connector created", "connector": createdConnector})
	})

	// Approve a membership request
	var acceptMembershipRequestPayload struct {
		OrganizationId string `form:"label=Organization ID"`
		UserId         string `form:"label=User ID"`
	}
	adminRouteController.registerAdminRoute("/approve-membership-request", "Approve Membership Request", acceptMembershipRequestPayload, func(ctx *gin.Context) {

		err := ctx.ShouldBind(&acceptMembershipRequestPayload)
		orgId, err := uuid.Parse(acceptMembershipRequestPayload.OrganizationId)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid organization id"})
			return
		}

		userId, err := uuid.Parse(acceptMembershipRequestPayload.UserId)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid user id"})
			return
		}

		orgService := orgservice.NewOrganizationService(adminRouteController.serverCfg)
		approvedRequest, err := orgService.ApprovePendingOrganizationMembershipRequest(ctx, orgId, userId)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "membership request approved", "data": approvedRequest})
	})

	// Get membership requests for an organization
	var getMembershipRequestsPayload struct {
		OrganizationId string `form:"label=Organization ID"`
	}
	adminRouteController.registerAdminRoute("/get-membership-requests", "Get Membership Request", getMembershipRequestsPayload, func(ctx *gin.Context) {

		err := ctx.ShouldBind(&getMembershipRequestsPayload)
		orgId, err := uuid.Parse(getMembershipRequestsPayload.OrganizationId)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid organization id"})
			return
		}

		orgService := orgservice.NewOrganizationService(adminRouteController.serverCfg)
		requests, err := orgService.GetOrganizationMembershipRequestsByOrganizationId(ctx, orgId)
		ctx.JSON(200, gin.H{"message": "membership requests fetched", "data": requests})
	})

	// Invite users to an organization
	var inviteMembersPayload struct {
		OrganizationId string `form:"label=Organization ID"`
		UserEmails     string `form:"label=User Emails"`
		Privilege      string `form:"label=Role"`
	}
	adminRouteController.registerAdminRoute("/invite-members", "Invite Members", inviteMembersPayload, func(ctx *gin.Context) {

		err := ctx.ShouldBind(&inviteMembersPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		orgId, err := uuid.Parse(inviteMembersPayload.OrganizationId)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid organization id"})
			return
		}

		orgService := orgservice.NewOrganizationService(adminRouteController.serverCfg)
		invitations, ierr := orgService.BulkInviteMembers(ctx, orgId, organizations.BulkInvitationPayload{
			Invitations: []organizations.InvitationPayload{
				{
					Email:     inviteMembersPayload.UserEmails,
					Privilege: inviteMembersPayload.Privilege,
				},
			},
		})
		if ierr.Error != nil {
			ctx.JSON(500, gin.H{"error": ierr.Error.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "processed", "invitations": invitations, "extensions": ierr})
	})

	var updateDatasetTypePayload struct {
		DatasetIDs  string `form:"label=Array of Dataset ID strings"`
		DatasetType string `form:"label=Dataset Type"`
	}
	adminRouteController.registerAdminRoute("/update-dataset-type", "Update Dataset Type", updateDatasetTypePayload, func(ctx *gin.Context) {

		err := ctx.ShouldBind(&updateDatasetTypePayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		datasetIds := []uuid.UUID{}
		err = json.Unmarshal([]byte(updateDatasetTypePayload.DatasetIDs), &datasetIds)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid dataset ids passed. please pass an array of dataset ids as a json array string"})
			return
		}

		var errors []string
		var updatedDatasetIds []uuid.UUID

		for _, datasetId := range datasetIds {
			datasetUpdatePayload := models.Dataset{
				ID:   datasetId,
				Type: models.DatasetType(updateDatasetTypePayload.DatasetType),
			}
			_, err := adminRouteController.serverCfg.Store.UpdateDataset(ctx, datasetUpdatePayload)
			if err != nil {
				errors = append(errors, err.Error())
				continue
			} else {
				updatedDatasetIds = append(updatedDatasetIds, datasetId)
			}
		}

		ctx.JSON(200, gin.H{"message": "processed", "updated_dataset_ids": updatedDatasetIds, "errors": errors})
	})

	// update a dataset
	var updateDatasetPayload struct {
		DatasetId        string `form:"label=Dataset ID"`
		Title            string `form:"label=Title,omitempty"`
		Description      string `form:"label=Description,omitempty"`
		Type             string `form:"label=Type,omitempty"`
		DatasetConfig    string `form:"label=Dataset Config (JSON),omitempty"`
		DisplayConfig    string `form:"label=Display Config (JSON),omitempty"`
		FileImportConfig string `form:"label=File Import Config (JSON),omitempty"`
	}
	adminRouteController.registerAdminRoute("/update-dataset", "Update Dataset", updateDatasetPayload, func(ctx *gin.Context) {
		var request dtos.UpdateDatasetRequest

		if err := ctx.ShouldBind(&updateDatasetPayload); err != nil {
			fmt.Printf("Binding error: %v\n", err)
			ctx.JSON(400, gin.H{"error": "binding error: " + err.Error()})
			return
		}

		if updateDatasetPayload.DatasetConfig != "" {
			var datasetConfig dataplatformDataModels.DatasetConfig
			err := json.Unmarshal([]byte(updateDatasetPayload.DatasetConfig), &datasetConfig)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "invalid dataset config JSON"})
				return
			}

			hasColumns := datasetConfig.Columns != nil
			hasCustomGroups := datasetConfig.CustomColumnGroups != nil
			hasRules := datasetConfig.Rules != nil

			if !hasColumns || !hasCustomGroups || !hasRules {
				ctx.JSON(400, gin.H{
					"error": "dataset config must include columns, custom column groups, and rules",
					"validation": map[string]bool{
						"has_columns":       hasColumns,
						"has_custom_groups": hasCustomGroups,
						"has_rules":         hasRules,
					},
				})
				return
			}
			request.DatasetConfig = &datasetConfig
		}

		if updateDatasetPayload.DisplayConfig != "" {
			var displayConfig []datasetmodels.DisplayConfig
			err := json.Unmarshal([]byte(updateDatasetPayload.DisplayConfig), &displayConfig)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "invalid display config JSON"})
				return
			}
			request.DisplayConfig = &displayConfig
		}

		if updateDatasetPayload.Title != "" {
			request.Title = &updateDatasetPayload.Title
		}

		if updateDatasetPayload.Description != "" {
			request.Description = &updateDatasetPayload.Description
		}

		if updateDatasetPayload.Type != "" {
			request.Type = &updateDatasetPayload.Type
		}

		role, _, orgIds := apicontext.GetAuthFromContext(ctx)
		if role != "admin" || len(orgIds) == 0 {
			ctx.JSON(403, gin.H{"error": "forbidden"})
			return
		}

		dpService, err := dataplatform.InitDataPlatformService(serverCfg.DataPlatformConfig)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		rulesService := rules.NewRuleService(serverCfg.Store)
		queryBuilderService := querybuilder.NewQueryBuilder()
		cloudService, err := cloudservice.NewCloudService("GCP", *serverCfg.Env)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if cloudService == nil {
			ctx.JSON(500, gin.H{"error": "CloudService is nil"})
			return
		}
		fileImportService := fileimportsservice.NewFileImportService(serverCfg.DefaultS3Client, serverCfg.Store, serverCfg.Env.AWSDefaultBucketName)
		datasetSvc := datasetservice.NewDatasetService(serverCfg.Store, queryBuilderService, dpService, rulesService, fileImportService, serverCfg.TemporalSdk, cloudService, serverCfg.DefaultS3Client, *serverCfg.DatasetConfig, serverCfg.CacheClient)

		actionId, err := datasetSvc.UpdateDataset(ctx, orgIds[0], updateDatasetPayload.DatasetId, request.ToModel())
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "dataset updated", "action_id": actionId})
	})

	// add user to organization
	var addUserToOrganizationPayload struct {
		OrganizationId uuid.UUID `form:"label=Organization ID"`
		UserId         uuid.UUID `form:"label=User ID"`
		Privilege      string    `form:"label=Privilege"`
	}

	adminRouteController.registerAdminRoute("/add-user-to-organization", "Add User to Organization", addUserToOrganizationPayload, func(ctx *gin.Context) {
		err := ctx.ShouldBind(&addUserToOrganizationPayload)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "bad payload"})
			return
		}

		if !slices.Contains(models.OrganizationPrivileges, models.ResourcePrivilege(addUserToOrganizationPayload.Privilege)) {
			ctx.JSON(400, gin.H{"error": "invalid privilege"})
			return
		}

		org, err := serverCfg.Store.GetOrganizationById(ctx, addUserToOrganizationPayload.OrganizationId.String())
		if err != nil {
			ctx.JSON(500, gin.H{"error": "error getting organization"})
			return
		}

		if org == nil {
			ctx.JSON(500, gin.H{"error": "empty organization"})
			return
		}

		userPolicy, err := serverCfg.Store.GetOrganizationPolicyByUser(ctx, addUserToOrganizationPayload.OrganizationId, addUserToOrganizationPayload.UserId)
		if userPolicy != nil {
			ctx.JSON(500, gin.H{"error": "user already in the organization"})
			return
		}

		policy, err := serverCfg.Store.CreateOrganizationPolicy(ctx, addUserToOrganizationPayload.OrganizationId, models.AudienceTypeUser, addUserToOrganizationPayload.UserId, models.ResourcePrivilege(addUserToOrganizationPayload.Privilege))
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if policy == nil {
			ctx.JSON(500, gin.H{"error": "failed to create organization policy"})
			return
		}

		ctx.JSON(200, gin.H{"message": "user added to organization", "policy": policy})
	})

}
