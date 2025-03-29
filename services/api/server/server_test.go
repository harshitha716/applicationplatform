package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	mockdataplatform "github.com/Zampfi/application-platform/services/api/mocks/core/dataplatform"
	mock_datasetservice "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mockruleservice "github.com/Zampfi/application-platform/services/api/mocks/core/rules/service"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mock_querybuilder "github.com/Zampfi/application-platform/services/api/mocks/pkg/querybuilder/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func getDataPlatformMockConfig() *serverconfig.DataPlatformConfig {
	return &serverconfig.DataPlatformConfig{
		DatabricksConfig: serverconfig.DatabricksSetupConfig{
			ZampDatabricksCatalog:        "zamp",
			ZampDatabricksPlatformSchema: "platform",
			MerchantDataProviderIdMapping: map[string]string{"merchant1": "workspace1"},
			DefaultDataProviderId:         "defaultWorkspace",
			DataProviderConfigs: map[string]dataplatformmodels.DatabricksConfig{
				"workspace1": {
					WarehouseId: "warehouse1",
				},
			},
		},
		PinotConfig: serverconfig.PinotSetupConfig{
			MerchantDataProviderIdMapping: map[string]string{"merchant2": "workspace2"},
			DefaultDataProviderId:         "defaultPinotWorkspace",
		},
		ActionsConfig: serverconfig.ActionsConfig{
			CreateMVJobTemplateConfig: serverconfig.CreateMVJobTemplateConfig{
				CreateMVNotebookPath:   "/path/to/create/mv/notebook",
				SideEffectNotebookPath: "/path/to/sideeffect/notebook",
			},
			WebhookConfig: serverconfig.WebhookConfig{
				WebhookId: "webhook_id",
				UserName:  "user_name",
				Password:  "password",
			},
		},
	}
}
func TestSetupRouter(t *testing.T) {
	serverCfg := serverconfig.GetEmptyServerConfig()
	if serverCfg == nil {
		t.Fatal("serverCfg is nil")
		return
	}

	mockDataplatformService := mockdataplatform.NewMockDataPlatformService(t)
	mockDatasetService := mock_datasetservice.NewMockDatasetService(t)
	mockRuleService := mockruleservice.NewMockRuleService(t)
	mockQueryBuilder := mock_querybuilder.NewMockQueryBuilder(t)
	kratosServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	mock_store := mock_store.NewMockStore(t)

	serverCfg.Env.AuthBaseUrl = kratosServer.URL
	serverCfg.Store = mock_store
	serverCfg.Env.PantheonURL = "http://localhost:9090"
	serverCfg.DataPlatformConfig = getDataPlatformMockConfig()

	if serverCfg == nil {
		t.Fatal("serverCfg is nil")
	}
	logger := zap.NewNop()
	if logger == nil {
		t.Fatal("logger is nil")
	}

	gin.SetMode(gin.TestMode)

	router, err := setupRouter(serverCfg, logger, mockQueryBuilder, mockDataplatformService, mockRuleService, mockDatasetService)

	assert.Nil(t, err)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"ok\"}", w.Body.String())
}
