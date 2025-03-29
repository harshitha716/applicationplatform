package sheets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/db/models"
	mock_datasetsService "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_sheets "github.com/Zampfi/application-platform/services/api/mocks/core/sheets"
	mockcache "github.com/Zampfi/application-platform/services/api/mocks/pkg/cache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(t *testing.T) (*gin.Engine, *mock_sheets.MockSheetsService, *serverconfig.ServerConfig) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := mock_sheets.NewMockSheetsService(t)
	serverCfg := &serverconfig.ServerConfig{}

	return router, mockService, serverCfg
}

func TestRegisterSheetsRoutes(t *testing.T) {
	router, _, serverCfg := setupRouter(t)
	group := router.Group("")

	// Verify that registration doesn't panic
	assert.NotPanics(t, func() {
		mockDatasetService := mock_datasetsService.NewMockDatasetService(t)
		mockCacheService := mockcache.NewMockCacheClient(t)
		RegisterSheetsRoutes(group, mockDatasetService, mockCacheService, serverCfg)
	})
}

func TestHandleGetSheetsAll(t *testing.T) {
	sheet1ID := uuid.MustParse("2e6bc1b8-a52a-4ea2-8761-5d1c5833d4d5")
	sheet2ID := uuid.MustParse("a821f150-6d28-47ab-aa51-15a55e7630f2")

	validSheets := []models.Sheet{
		{ID: sheet1ID, Name: "Sheet 1", SheetConfig: json.RawMessage(`{"test":"test"}`)},
		{ID: sheet2ID, Name: "Sheet 2", SheetConfig: json.RawMessage(`{"test":"test"}`)},
	}

	pageId := uuid.New()

	tests := []struct {
		name         string
		setupMock    func(*mock_sheets.MockSheetsService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mock_sheets.MockSheetsService) {
				m.EXPECT().
					GetSheetsByPageId(mock.Anything, pageId).
					Return(validSheets, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: validSheets,
		},
		{
			name: "service error",
			setupMock: func(m *mock_sheets.MockSheetsService) {
				m.EXPECT().
					GetSheetsByPageId(mock.Anything, pageId).
					Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService, _ := setupRouter(t)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/sheets/", pageId.String()), nil)

			sheetsGroup := router.Group("/:pageId/sheets")
			sheetsGroup.GET("/", func(c *gin.Context) {
				handleGetSheetsAll(c, mockService)
			})

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var actualSheets []models.Sheet
				err := json.Unmarshal(w.Body.Bytes(), &actualSheets)
				assert.NoError(t, err)

				// Compare each field separately to avoid JSON formatting issues
				assert.Equal(t, len(validSheets), len(actualSheets))
				for i := range validSheets {
					assert.Equal(t, validSheets[i].ID, actualSheets[i].ID)
					assert.Equal(t, validSheets[i].Name, actualSheets[i].Name)

					// Compare the parsed JSON instead of raw bytes
					var expectedConfig, actualConfig map[string]interface{}
					err = json.Unmarshal(validSheets[i].SheetConfig, &expectedConfig)
					assert.NoError(t, err)
					err = json.Unmarshal(actualSheets[i].SheetConfig, &actualConfig)
					assert.NoError(t, err)
					assert.Equal(t, expectedConfig, actualConfig)
				}
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandleGetSheetByID(t *testing.T) {

	pageId := uuid.New()
	validID := uuid.New()
	validSheet := &models.Sheet{
		ID:              uuid.New(),
		Name:            "Sheet 1",
		WidgetInstances: []models.WidgetInstance{{ID: uuid.New(), DataMappings: nil}},
		SheetConfig: json.RawMessage(`{
			"version": "1.0",
			"native_filter_config": [],
			"sheet_layout": [
				{
					"id": "123e4567-e89b-12d3-a456-426614174000",
					"type": "widget",
					"name": "Widget 1",
					"layout": {
						"x": 0,
						"y": 0,
						"w": 12,
						"h": 8
					}
				}
			]
		}`),
	}

	tests := []struct {
		name         string
		sheetID      string
		setupMock    func(*mock_sheets.MockSheetsService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:    "successful retrieval",
			sheetID: validID.String(),
			setupMock: func(m *mock_sheets.MockSheetsService) {
				m.EXPECT().
					GetSheetById(mock.Anything, validID).
					Return(validSheet, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: validSheet,
		},
		// {
		// 	name:    "invalid uuid",
		// 	sheetID: "invalid-uuid",
		// 	setupMock: func(m *mock_sheets.MockSheetsService) {
		// 		// No mock expectations needed for invalid UUID
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// 	expectedBody: map[string]interface{}{"error": "invalid sheet id"},
		// },
		// {
		// 	name:    "service error",
		// 	sheetID: validID.String(),
		// 	setupMock: func(m *mock_sheets.MockSheetsService) {
		// 		m.EXPECT().
		// 			GetSheetById(mock.Anything, validID).
		// 			Return(nil, errors.New("service error"))
		// 	},
		// 	expectedCode: http.StatusInternalServerError,
		// 	expectedBody: map[string]interface{}{"error": "internal server error"},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService, _ := setupRouter(t)

			// Setup mock expectations
			tt.setupMock(mockService)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/sheets/%s", pageId.String(), tt.sheetID), nil)

			// Setup route
			sheetsGroup := router.Group("/:pageId/sheets")
			sheetsGroup.GET("/:sheetId", func(c *gin.Context) {
				handleGetSheetByID(c, mockService)
			})

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert response body
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedCode == http.StatusOK {
				var sheet models.Sheet
				err = json.Unmarshal(w.Body.Bytes(), &sheet)
				assert.NoError(t, err)

				assert.Equal(t, sheet.ID, validSheet.ID)
				assert.Equal(t, sheet.Name, validSheet.Name)
				assert.Equal(t, len(sheet.WidgetInstances), len(validSheet.WidgetInstances))
				assert.Equal(t, sheet.CreatedAt, validSheet.CreatedAt)
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
