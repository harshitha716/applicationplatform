package widgets

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	datasetsmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_datasets "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_widgets "github.com/Zampfi/application-platform/services/api/mocks/core/widgets/service"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegisterWidgetRoutes tests the RegisterWidgetRoutes function
func TestRegisterWidgetRoutes(t *testing.T) {
	t.Run("register routes successfully", func(t *testing.T) {
		// Setup
		gin.SetMode(gin.TestMode)
		router := gin.New()
		routerGroup := router.Group("/api")
		
		// Create mock dependencies
		mockStore := mock_store.NewMockStore(t)
		mockDatasetService := mock_datasets.NewMockDatasetService(t)
		
		// Create server config with mock store
		serverCfg := &serverconfig.ServerConfig{
			Store: mockStore,
		}
		
		// Register routes
		err := RegisterWidgetRoutes(routerGroup, serverCfg, mockDatasetService)
		
		// Verify routes were registered successfully
		assert.NoError(t, err)
		
		// We'll just verify that the routes were registered without making actual requests
		// since that would require mocking the service implementation
	})
}

// TestGetWidgetInstance tests the GetWidgetInstance handler
func TestGetWidgetInstance(t *testing.T) {
	tests := []struct {
		name           string
		widgetID       string
		setupMock      func(*mock_widgets.MockWidgetsService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "success",
			widgetID: uuid.New().String(),
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				m.On("GetWidgetInstance", mock.Anything, mock.Anything).Return(widgetmodels.WidgetInstance{
					ID:         uuid.New(),
					Title:      "Test Widget",
					WidgetType: "bar_chart",
					DataMappings: widgetmodels.DataMappings{
						Version:  widgetmodels.DataMappingVersion1,
						Mappings: []widgetmodels.DataMappingFields{},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid widget ID",
			widgetID: "invalid-uuid",
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid widget instance ID"},
		},
		{
			name:     "service error",
			widgetID: uuid.New().String(),
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				m.On("GetWidgetInstance", mock.Anything, mock.Anything).Return(widgetmodels.WidgetInstance{}, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "failed to get widget instance"},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := mock_widgets.NewMockWidgetsService(t)
			tt.setupMock(mockService)
			
			router := gin.New()
			router.GET("/widgets/:widgetInstanceId/instance", func(c *gin.Context) {
				GetWidgetInstance(c, mockService)
			})
			
			// Execute request
			req, _ := http.NewRequest("GET", "/widgets/"+tt.widgetID+"/instance", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedBody != nil {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				expectedBody, _ := tt.expectedBody.(gin.H)
				assert.Equal(t, expectedBody["error"], response["error"])
			}
		})
	}
}

// TestGetWidgetInstanceData tests the GetWidgetInstanceData handler
func TestGetWidgetInstanceData(t *testing.T) {
	orgID := uuid.New()
	
	tests := []struct {
		name           string
		widgetID       string
		queryParams    map[string]string
		skipAuth       bool
		setupMock      func(*mock_widgets.MockWidgetsService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "success with no filters",
			widgetID: uuid.New().String(),
			queryParams: map[string]string{
				"periodicity": "monthly",
				"currency":    "USD",
			},
			skipAuth: false,
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// Mock dataset data response
				datasetData := []datasetsmodels.DatasetData{
					{
						TotalCount:  nil,
						Title:       "Test Dataset",
						Description: nil,
						DatasetConfig: datasetsmodels.DatasetConfig{
							IsDrilldownEnabled: false,
						},
						Metadata: datasetsmodels.DatasetMetadataConfig{},
					},
				}
				
				m.On("GetWidgetInstanceData", mock.Anything, orgID, mock.Anything, mock.Anything).Return(datasetData, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "no organization context",
			widgetID:    uuid.New().String(),
			queryParams: map[string]string{},
			skipAuth:    true, // Skip auth for this test
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "no organization ids found"},
		},
		{
			name:        "invalid widget ID",
			widgetID:    "invalid-uuid",
			queryParams: map[string]string{},
			skipAuth:    false,
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid widget instance ID"},
		},
		{
			name:     "invalid filters format",
			widgetID: uuid.New().String(),
			queryParams: map[string]string{
				"filters": "invalid-json",
			},
			skipAuth: false,
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid filters
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid filters format"},
		},
		{
			name:     "invalid time columns format",
			widgetID: uuid.New().String(),
			queryParams: map[string]string{
				"time_columns": "invalid-json",
			},
			skipAuth: false,
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid time columns
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid time columns format"},
		},
		{
			name:     "service error",
			widgetID: uuid.New().String(),
			queryParams: map[string]string{},
			skipAuth: false,
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				m.On("GetWidgetInstanceData", mock.Anything, orgID, mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "failed to get widget instance data"},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// Setup
			gin.SetMode(gin.TestMode)
			mockService := mock_widgets.NewMockWidgetsService(t)
			tt.setupMock(mockService)
			
			router := gin.New()
			
			// Add auth middleware for this specific test
			if !tt.skipAuth {
				router.Use(func(c *gin.Context) {
					apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
					c.Next()
				})
			}
			
			router.GET("/widgets/:widgetInstanceId/data", func(c *gin.Context) {
				GetWidgetInstanceData(c, mockService)
			})
			
			// Create request with query parameters
			req, _ := http.NewRequest("GET", "/widgets/"+tt.widgetID+"/data", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			
			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedBody != nil {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				expectedBody, _ := tt.expectedBody.(gin.H)
				assert.Equal(t, expectedBody["error"], response["error"])
			}
		})
	}
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
