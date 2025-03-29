package widgets

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	datasetsmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_widgets "github.com/Zampfi/application-platform/services/api/mocks/core/widgets/service"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWidgetInstanceHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		setupMock      func(*mock_widgets.MockWidgetsService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			setupContext: func(c *gin.Context) {
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				widgetID := uuid.New()
				m.On("GetWidgetInstance", mock.Anything, mock.Anything).Return(widgetmodels.WidgetInstance{
					ID:         widgetID,
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
			name: "invalid widget ID",
			setupContext: func(c *gin.Context) {
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: "invalid-uuid"},
				}
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid widget instance ID"},
		},
		{
			name: "service error",
			setupContext: func(c *gin.Context) {
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
			},
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
			
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			// Setup context
			tt.setupContext(c)
			
			// Setup mock
			mockService := mock_widgets.NewMockWidgetsService(t)
			tt.setupMock(mockService)
			
			// Call handler
			GetWidgetInstance(c, mockService)
			
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

func TestGetWidgetInstanceDataHandler(t *testing.T) {
	orgID := uuid.New()
	
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		setupMock      func(*mock_widgets.MockWidgetsService)
		expectedStatus int
		expectedBody   interface{}
		checkResponse  func(*testing.T, []byte)
	}{
		{
			name: "success with no filters",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
				
				// Set query params
				c.Request = httptest.NewRequest("GET", "/?periodicity=monthly&currency=USD", nil)
			},
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
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.Equal(t, "monthly", response["periodicity"])
				assert.Equal(t, "USD", response["currency"])
			},
		},
		{
			name: "success with valid filters",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				widgetID := uuid.New()
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: widgetID.String()},
				}
				
				// Set valid filters query param
				filtersJSON := `[{"dataset_id":"dataset1","filters":{"conditions":[{"column":"col1","operator":"eq","value":"val1"}]}}]`
				c.Request = httptest.NewRequest("GET", "/?filters="+filtersJSON, nil)
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// Mock dataset data response with filtered data
				datasetData := []datasetsmodels.DatasetData{
					{
						QueryResult: dataplatformmodels.QueryResult{
							Columns: []dataplatformmodels.ColumnMetadata{
								{Name: "col1", DatabaseType: "string"},
							},
							Rows: []map[string]interface{}{
								{"col1": "val1"},
							},
						},
						Title: "Test Dataset",
					},
				}
				
				m.On("GetWidgetInstanceData", mock.Anything, orgID, mock.Anything, mock.Anything).Return(datasetData, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				
				// Check result structure
				results, ok := response["result"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, results, 1)
				
				// Check columns
				result := results[0].(map[string]interface{})
				columns, ok := result["columns"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, columns, 1)
				
				// Check column details
				column := columns[0].(map[string]interface{})
				assert.Equal(t, "col1", column["column_name"])
				assert.Equal(t, "string", column["column_type"])
			},
		},
		{
			name: "success with valid time columns",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				widgetID := uuid.New()
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: widgetID.String()},
				}
				
				// Set valid time columns query param
				timeColumnsJSON := `[{"dataset_id":"dataset1","column":"date_col"}]`
				c.Request = httptest.NewRequest("GET", "/?time_columns="+timeColumnsJSON, nil)
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// Mock dataset data response with time column data
				datasetData := []datasetsmodels.DatasetData{
					{
						QueryResult: dataplatformmodels.QueryResult{
							Columns: []dataplatformmodels.ColumnMetadata{
								{Name: "date_col", DatabaseType: "timestamp"},
							},
							Rows: []map[string]interface{}{
								{"date_col": "2023-01-01"},
							},
						},
						Title: "Test Dataset",
					},
				}
				
				m.On("GetWidgetInstanceData", mock.Anything, orgID, mock.Anything, mock.Anything).Return(datasetData, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				
				// Check result structure
				results, ok := response["result"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, results, 1)
				
				// Check columns
				result := results[0].(map[string]interface{})
				columns, ok := result["columns"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, columns, 1)
				
				// Check column details
				column := columns[0].(map[string]interface{})
				assert.Equal(t, "date_col", column["column_name"])
				assert.Equal(t, "timestamp", column["column_type"])
			},
		},
		{
			name: "success with both filters and time columns",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				widgetID := uuid.New()
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: widgetID.String()},
				}
				
				// Set both filters and time columns query params
				filtersJSON := `[{"dataset_id":"dataset1","filters":{"conditions":[{"column":"category","operator":"eq","value":"A"}]}}]`
				timeColumnsJSON := `[{"dataset_id":"dataset1","column":"date_col"}]`
				c.Request = httptest.NewRequest("GET", "/?filters="+filtersJSON+"&time_columns="+timeColumnsJSON+"&periodicity=weekly", nil)
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// Mock dataset data response with filtered data and time columns
				datasetData := []datasetsmodels.DatasetData{
					{
						QueryResult: dataplatformmodels.QueryResult{
							Columns: []dataplatformmodels.ColumnMetadata{
								{Name: "category", DatabaseType: "string"},
								{Name: "date_col", DatabaseType: "timestamp"},
								{Name: "value", DatabaseType: "int"},
							},
							Rows: []map[string]interface{}{
								{"category": "A", "date_col": "2023-01-01", "value": 100},
								{"category": "A", "date_col": "2023-01-08", "value": 150},
							},
						},
						Title: "Test Dataset",
					},
				}
				
				m.On("GetWidgetInstanceData", mock.Anything, orgID, mock.Anything, mock.Anything).Return(datasetData, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.Equal(t, "weekly", response["periodicity"])
				
				// Check result structure
				results, ok := response["result"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, results, 1)
				
				// Check row count
				result := results[0].(map[string]interface{})
				assert.Equal(t, float64(2), result["rowcount"])
				
				// Check columns
				columns, ok := result["columns"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, columns, 3)
			},
		},
		{
			name: "no organization context",
			setupContext: func(c *gin.Context) {
				// No auth context setup
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "no organization ids found"},
		},
		{
			name: "invalid widget ID",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set invalid widget ID param
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: "invalid-uuid"},
				}
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid widget instance ID"},
		},
		{
			name: "invalid filters format",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
				
				// Set invalid filters query param
				c.Request = httptest.NewRequest("GET", "/?filters=invalid-json", nil)
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid filters
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid filters format"},
		},
		{
			name: "invalid time columns format",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
				
				// Set invalid time columns query param
				c.Request = httptest.NewRequest("GET", "/?time_columns=invalid-json", nil)
			},
			setupMock: func(m *mock_widgets.MockWidgetsService) {
				// No mock setup needed for invalid time columns
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid time columns format"},
		},
		{
			name: "service error",
			setupContext: func(c *gin.Context) {
				// Add auth context
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{orgID})
				
				// Set widget ID param
				c.Params = []gin.Param{
					{Key: "widgetInstanceId", Value: uuid.New().String()},
				}
			},
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
			
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			// Setup context
			tt.setupContext(c)
			
			// Setup mock
			mockService := mock_widgets.NewMockWidgetsService(t)
			tt.setupMock(mockService)
			
			// Call handler
			GetWidgetInstanceData(c, mockService)
			
			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedBody != nil {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				expectedBody, _ := tt.expectedBody.(gin.H)
				assert.Equal(t, expectedBody["error"], response["error"])
			}
			
			// Check response structure if needed
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}
