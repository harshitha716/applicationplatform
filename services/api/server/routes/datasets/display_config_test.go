package datasets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	dsMock "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_fileimports "github.com/Zampfi/application-platform/services/api/mocks/core/fileimports"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/Zampfi/application-platform/services/api/server/middleware"
)

func TestGetDatasetDisplayConfig(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	datasetId := uuid.New()
	datasetIdStr := datasetId.String()
	merchantId := uuid.New()

	displayConfig := []datasetmodels.DisplayConfig{
		{
			Column:     "test_column",
			IsHidden:   false,
			IsEditable: true,
		},
	}

	tests := []struct {
		name         string
		path         string
		setupMock    func(*dsMock.MockDatasetService, *mock_store.MockStore)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "Success case - existing display config",
			path: fmt.Sprintf("/datasets/%s/display-config", datasetIdStr),
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				m.EXPECT().GetDatasetDisplayConfig(mock.Anything, merchantId, datasetId.String()).Return(displayConfig, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{
				"display_config": displayConfig,
			},
		},
		{
			name: "Service error",
			path: fmt.Sprintf("/datasets/%s/display-config", datasetIdStr),
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				m.EXPECT().GetDatasetDisplayConfig(mock.Anything, merchantId, datasetId.String()).Return(nil, errors.New("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: gin.H{"error": "internal error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			g := e.Group("/")

			mockDatasetService := dsMock.NewMockDatasetService(t)
			mockStore := mock_store.NewMockStore(t)
			mockFileUploadService := mock_fileimports.NewMockFileImportService(t)
			tt.setupMock(mockDatasetService, mockStore)

			g.Use(func(c *gin.Context) {
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{merchantId})
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&dbmodels.Dataset{ID: datasetId, Metadata: json.RawMessage(`{}`)}, nil)
				c.Set("datasetContext", middleware.DatasetContext{
					DatasetID:  datasetIdStr,
					MerchantID: merchantId,
				})
				c.Next()
			})

			registerRoutes(g, mockDatasetService, mockStore, mockFileUploadService)

			// Create a test request
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// fire the request
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// For success case, compare the actual data
				if tt.expectedCode == http.StatusOK {
					dataJSON, err := json.Marshal(response["display_config"])
					assert.NoError(t, err)

					var actualDisplayConfig []datasetmodels.DisplayConfig
					err = json.Unmarshal(dataJSON, &actualDisplayConfig)
					assert.NoError(t, err)

					assert.Equal(t, displayConfig, actualDisplayConfig)
				} else {
					// For error cases, compare the error message
					assert.Equal(t, tt.expectedBody.(gin.H)["error"], response["error"])
				}
			}
		})
	}
}

func TestSetDatasetDisplayConfig(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	datasetId := uuid.New()
	datasetIdStr := datasetId.String()
	merchantId := uuid.New()
	actionId := uuid.New()

	displayConfig := []datasetmodels.DisplayConfig{
		{
			Column:     "test_column",
			IsHidden:   false,
			IsEditable: true,
		},
	}

	tests := []struct {
		name          string
		path          string
		body          interface{}
		setupMock     func(*dsMock.MockDatasetService, *mock_store.MockStore)
		expectedCode  int
		expectedBody  interface{}
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success case",
			path: fmt.Sprintf("/datasets/%s/display-config", datasetIdStr),
			body: map[string]interface{}{
				"display_config": displayConfig,
			},
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				m.EXPECT().UpdateDataset(mock.Anything, merchantId, datasetId.String(), mock.Anything).Return(actionId.String(), nil)
			},
			expectedCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, actionId.String(), resp["action_id"])
			},
		},
		{
			name: "Invalid request body",
			path: fmt.Sprintf("/datasets/%s/display-config", datasetIdStr),
			body: "invalid json", // This will cause a JSON unmarshal error
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				// No mock expectations needed as the handler should return early
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid character 'i' looking for beginning of value"},
		},
		{
			name: "Service error",
			path: fmt.Sprintf("/datasets/%s/display-config", datasetIdStr),
			body: map[string]interface{}{
				"display_config": displayConfig,
			},
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				m.EXPECT().UpdateDataset(mock.Anything, merchantId, datasetId.String(), mock.Anything).Return("", errors.New("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: gin.H{"error": "internal error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			g := e.Group("/")

			mockDatasetService := dsMock.NewMockDatasetService(t)
			mockStore := mock_store.NewMockStore(t)
			mockFileUploadService := mock_fileimports.NewMockFileImportService(t)
			tt.setupMock(mockDatasetService, mockStore)

			g.Use(func(c *gin.Context) {
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{merchantId})
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&dbmodels.Dataset{ID: datasetId, Metadata: json.RawMessage(`{}`)}, nil)
				c.Set("datasetContext", middleware.DatasetContext{
					DatasetID:  datasetIdStr,
					MerchantID: merchantId,
				})
				c.Next()
			})

			registerRoutes(g, mockDatasetService, mockStore, mockFileUploadService)

			// Create a test request
			var bodyReader *bytes.Reader
			if jsonBody, ok := tt.body.(map[string]interface{}); ok {
				bodyJSON, _ := json.Marshal(jsonBody)
				bodyReader = bytes.NewReader(bodyJSON)
			} else if stringBody, ok := tt.body.(string); ok {
				bodyReader = bytes.NewReader([]byte(stringBody))
			}

			req, err := http.NewRequest("POST", tt.path, bodyReader)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// fire the request
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			if tt.expectedBody != nil {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(gin.H)["error"], response["error"])
			}
		})
	}
}
