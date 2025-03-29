package datasets

// This file is intentionally kept minimal to avoid conflicts with display_config_test.go
// All display config tests have been moved to display_config_test.go

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

	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	dsMock "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_fileimports "github.com/Zampfi/application-platform/services/api/mocks/core/fileimports"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/Zampfi/application-platform/services/api/server/middleware"
	"github.com/Zampfi/application-platform/services/api/server/routes/datasets/dtos"
)

func TestUpdateDataset(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	datasetId := uuid.New()
	datasetIdStr := datasetId.String()
	merchantId := uuid.New()
	actionId := uuid.New()

	validDatasetConfig := &dataplatformDataModels.DatasetConfig{
		Columns:            map[string]dataplatformDataModels.DatasetColumnConfig{},
		CustomColumnGroups: []dataplatformDataModels.CustomColumnGroup{},
		Rules:              map[string][]dataplatformDataModels.Rule{},
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
			name: "Success case - valid update",
			path: fmt.Sprintf("/datasets/%s/update", datasetIdStr),
			body: dtos.UpdateDatasetRequest{
				Title:         stringPtr("Updated Title"),
				Description:   stringPtr("Updated Description"),
				DatasetConfig: validDatasetConfig,
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
			path: fmt.Sprintf("/datasets/%s/update", datasetIdStr),
			body: "invalid json",
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				// No mock expectations needed as the handler should return early
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid character 'i' looking for beginning of value"},
		},
		{
			name: "Invalid dataset config",
			path: fmt.Sprintf("/datasets/%s/update", datasetIdStr),
			body: dtos.UpdateDatasetRequest{
				Title:         stringPtr("Updated Title"),
				Description:   stringPtr("Updated Description"),
				DatasetConfig: &dataplatformDataModels.DatasetConfig{
					// Missing required fields
				},
			},
			setupMock: func(m *dsMock.MockDatasetService, mockStore *mock_store.MockStore) {
				// No mock expectations needed as the handler should return early
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "dataset config must include columns, custom column groups, and rules",
				"validation": map[string]bool{
					"has_columns":       false,
					"has_custom_groups": false,
					"has_rules":         false,
				},
			},
		},
		{
			name: "Service error",
			path: fmt.Sprintf("/datasets/%s/update", datasetIdStr),
			body: dtos.UpdateDatasetRequest{
				Title:         stringPtr("Updated Title"),
				Description:   stringPtr("Updated Description"),
				DatasetConfig: validDatasetConfig,
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
			// Register routes
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

			// Create test request
			var bodyReader *bytes.Reader
			if jsonBody, ok := tt.body.(dtos.UpdateDatasetRequest); ok {
				bodyJSON, _ := json.Marshal(jsonBody)
				bodyReader = bytes.NewReader(bodyJSON)
			} else if stringBody, ok := tt.body.(string); ok {
				bodyReader = bytes.NewReader([]byte(stringBody))
			}

			req, err := http.NewRequest("PATCH", tt.path, bodyReader)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Fire the request
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)
				assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			}
		})
	}
}

func TestConfirmDatasetImport(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	merchantId := uuid.New()
	fileUploadId := uuid.New()
	datasetId := uuid.New()

	tests := []struct {
		name         string
		path         string
		body         interface{}
		setupMock    func(*dsMock.MockDatasetService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "Success case",
			path: fmt.Sprintf("/datasets/file-imports/%s/confirm", fileUploadId),
			body: dtos.ConfirmDatasetImportRequest{
				DatasetId: datasetId,
			},
			setupMock: func(m *dsMock.MockDatasetService) {
				m.EXPECT().ImportDataFromFile(
					mock.Anything,
					merchantId,
					datasetId,
					fileUploadId,
				).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"message": "success"},
		},
		{
			name: "Invalid request body",
			path: fmt.Sprintf("/datasets/file-imports/%s/confirm", fileUploadId),
			body: "invalid json",
			setupMock: func(m *dsMock.MockDatasetService) {
				// No mock expectations needed as handler should return early
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid request"},
		},
		{
			name: "Service error",
			path: fmt.Sprintf("/datasets/file-imports/%s/confirm", fileUploadId),
			body: dtos.ConfirmDatasetImportRequest{
				DatasetId: datasetId,
			},
			setupMock: func(m *dsMock.MockDatasetService) {
				m.EXPECT().ImportDataFromFile(
					mock.Anything,
					merchantId,
					datasetId,
					fileUploadId,
				).Return(errors.New("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: gin.H{"error": "internal error"},
		},
		{
			name: "Invalid file upload ID",
			path: "/datasets/file-imports/invalid-uuid/confirm",
			body: dtos.ConfirmDatasetImportRequest{
				DatasetId: datasetId,
			},
			setupMock: func(m *dsMock.MockDatasetService) {
				// No mock expectations needed as handler should return early
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid file upload id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			g := e.Group("/")

			mockDatasetService := dsMock.NewMockDatasetService(t)
			mockStore := mock_store.NewMockStore(t)
			mockFileUploadService := mock_fileimports.NewMockFileImportService(t)
			tt.setupMock(mockDatasetService)

			g.Use(func(c *gin.Context) {
				apicontext.AddAuthToGinContext(c, "user", uuid.New(), []uuid.UUID{merchantId})
				c.Next()
			})

			registerRoutes(g, mockDatasetService, mockStore, mockFileUploadService)

			// Create test request
			var bodyReader *bytes.Reader
			if jsonBody, ok := tt.body.(dtos.ConfirmDatasetImportRequest); ok {
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

			// Fire the request
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)
				assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			}
		})
	}
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
