package pages

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/pages"
	"github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_dataset_service "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_pages_service "github.com/Zampfi/application-platform/services/api/mocks/core/pages"
	mockcache "github.com/Zampfi/application-platform/services/api/mocks/pkg/cache"
	"github.com/Zampfi/application-platform/services/api/server/routes/pages/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(t *testing.T) (*gin.Engine, *mock_pages_service.MockPagesService, *serverconfig.ServerConfig) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := mock_pages_service.NewMockPagesService(t)
	serverCfg := &serverconfig.ServerConfig{}

	return router, mockService, serverCfg
}

func TestRegisterPagesRoutes(t *testing.T) {
	router, _, serverCfg := setupRouter(t)
	group := router.Group("")

	mockDatasetService := mock_dataset_service.NewMockDatasetService(t)
	RegisterPagesRoutes(group, mockDatasetService, serverCfg)
}
func TestHandleGetPagesAll(t *testing.T) {
	validOrgID := uuid.New()
	validUserID := uuid.New()
	description1 := "Description 1"
	description2 := "Description 2"
	validPages := []models.Page{
		{ID: uuid.New(), Name: "Page 1", Description: &description1, OrganizationId: validOrgID, FractionalIndex: 1},
		{ID: uuid.New(), Name: "Page 2", Description: &description2, OrganizationId: validOrgID, FractionalIndex: 2},
	}

	tests := []struct {
		name         string
		setupMock    func(*mock_pages_service.MockPagesService)
		setupContext func(*gin.Context)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "successful retrieval",
			setupContext: func(c *gin.Context) {
				apicontext.AddAuthToGinContext(c, "user", validUserID, []uuid.UUID{validOrgID})
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPagesByOrganizationId(mock.Anything, validOrgID).
					Return(validPages, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"page_id":          validPages[0].ID.String(),
					"name":             validPages[0].Name,
					"description":      description1,
					"created_at":       validPages[0].CreatedAt.Format(time.RFC3339),
					"updated_at":       validPages[0].UpdatedAt.Format(time.RFC3339),
					"organization_id":  validPages[0].OrganizationId.String(),
					"fractional_index": validPages[0].FractionalIndex,
				},
				map[string]interface{}{
					"page_id":          validPages[1].ID.String(),
					"name":             validPages[1].Name,
					"description":      description2,
					"created_at":       validPages[1].CreatedAt.Format(time.RFC3339),
					"updated_at":       validPages[1].UpdatedAt.Format(time.RFC3339),
					"organization_id":  validPages[1].OrganizationId.String(),
					"fractional_index": validPages[1].FractionalIndex,
				},
			},
		},
		{
			name: "missing user ID",
			setupContext: func(c *gin.Context) {
				c.Set("organizationIds", []uuid.UUID{validOrgID})
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				// No mock expectations needed
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{"error": "unauthorized"},
		},
		{
			name: "missing org IDs",
			setupContext: func(c *gin.Context) {
				apicontext.AddAuthToGinContext(c, "user", validUserID, []uuid.UUID{})
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				// No mock expectations needed
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{"error": "unauthorized"},
		},
		{
			name: "service error",
			setupContext: func(c *gin.Context) {
				apicontext.AddAuthToGinContext(c, "user", validUserID, []uuid.UUID{validOrgID})
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPagesByOrganizationId(mock.Anything, validOrgID).
					Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService, _ := setupRouter(t)

			// Setup mock expectations
			tt.setupMock(mockService)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/pages", nil)

			// Setup route
			group := router.Group("")
			pagesGroup := group.Group("/pages")
			pagesGroup.GET("", func(c *gin.Context) {
				// Setup context before handler
				tt.setupContext(c)
				handleGetPagesAll(c, mockService)
			})

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert response body
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestHandleGetPagesByID(t *testing.T) {
	validID := uuid.New()
	validPage := &models.Page{ID: validID, Name: "Test Page"}

	tests := []struct {
		name         string
		pageID       string
		setupMock    func(*mock_pages_service.MockPagesService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:   "successful retrieval",
			pageID: validID.String(),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPageByID(mock.Anything, validID).
					Return(validPage, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: validPage,
		},
		{
			name:   "invalid uuid",
			pageID: "invalid-uuid",
			setupMock: func(m *mock_pages_service.MockPagesService) {
				// No mock expectations needed for invalid UUID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "invalid page id"},
		},
		{
			name:   "service error",
			pageID: validID.String(),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPageByID(mock.Anything, validID).
					Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService, _ := setupRouter(t)

			// Setup mock expectations
			tt.setupMock(mockService)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/pages/"+tt.pageID, nil)

			// Setup route
			group := router.Group("")
			pagesGroup := group.Group("/pages")
			pagesGroup.GET("/:pageId", func(c *gin.Context) {
				handleGetPagesByID(c, mockService)
			})

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert response body
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// For successful cases, compare the actual page
			if tt.expectedCode == http.StatusOK {
				var page models.Page
				err = json.Unmarshal(w.Body.Bytes(), &page)
				assert.NoError(t, err)
				assert.Equal(t, validPage, &page)
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandleGetPageAudiences(t *testing.T) {
	pageId := uuid.New()

	policies := []models.ResourceAudiencePolicy{{ResourceID: pageId, ResourceType: "page"}}
	mockCacheService := mockcache.NewMockCacheClient(t)
	tests := []struct {
		name          string
		path          string
		setupMock     func(*mock_pages_service.MockPagesService)
		expectedCode  int
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful retrieval",
			path: fmt.Sprintf("/pages/%s/audiences", pageId),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPageAudiences(mock.Anything, pageId).
					Return(policies, nil)
			},
			expectedCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []models.ResourceAudiencePolicy
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp, 1)
				assert.Equal(t, pageId, resp[0].ResourceID)
				assert.Equal(t, models.ResourceTypePage, resp[0].ResourceType)
			},
		},
		{
			name: "invalid uuid",
			path: fmt.Sprintf("/pages/%s/audiences", "not-uuid"),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				// No mock expectations needed for invalid UUID
			},
			expectedCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "invalid page id", resp["error"])
			},
		},
		{
			name: "service error",
			path: fmt.Sprintf("/pages/%s/audiences", pageId),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPageAudiences(mock.Anything, pageId).
					Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "something went wrong", resp["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockPageService, serverCfg := setupRouter(t)

			// Setup mock expectations
			tt.setupMock(mockPageService)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.path, nil)

			datasetService := mock_dataset_service.NewMockDatasetService(t)

			// Setup route
			group := router.Group("")
			registerRoutes(group, mockPageService, datasetService, mockCacheService, serverCfg)

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert response body
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, w.Code)

			tt.checkResponse(t, w)
		})
	}
}

func TestAddPageAudiences(t *testing.T) {
	validPageID := uuid.New()
	validAudience := uuid.New()

	tests := []struct {
		name            string
		pageID          string
		getPayloadBytes func() []byte
		setupMock       func(*mock_pages_service.MockPagesService)
		expectedStatus  int
		expectedBody    map[string]interface{}
	}{
		{
			name:   "Success case",
			pageID: validPageID.String(),
			getPayloadBytes: func() []byte {
				payload := map[string]interface{}{
					"audiences": []map[string]interface{}{
						{
							"audience_id":   validAudience.String(),
							"audience_type": "user",
							"privilege":     "read",
						},
					},
				}
				payloadBytes, _ := json.Marshal(payload)
				return payloadBytes
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().BulkAddAudienceToPage(
					mock.Anything,
					validPageID,
					mock.Anything,
				).Return([]*models.ResourceAudiencePolicy{
					{
						ID:         validAudience,
						ResourceID: validPageID,
					},
				}, pages.BulkAddPageAudienceErrors{})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"audiences": []interface{}{
					map[string]interface{}{
						"resource_audience_policy_id": validAudience.String(),
						"resource_audience_id":        "00000000-0000-0000-0000-000000000000",
						"resource_audience_type":      "",
						"resource_id":                 validPageID.String(),
						"resource_type":               "",
						"privilege":                   "",
						"created_at":                  "0001-01-01T00:00:00Z",
						"updated_at":                  "0001-01-01T00:00:00Z",
						"deleted_at":                  nil,
						"metadata":                    nil,
					},
				},
				"audience_errors": nil,
			},
		},
		{
			name:   "Invalid page ID",
			pageID: "invalid-uuid",
			getPayloadBytes: func() []byte {
				payload := map[string]interface{}{
					"audiences": []map[string]interface{}{
						{
							"audience_id":   validAudience.String(),
							"audience_type": "user",
							"privilege":     "viewer",
						},
					},
				}
				payloadBytes, _ := json.Marshal(payload)
				return payloadBytes
			},
			setupMock:      func(m *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid page ID",
			},
		},
		{
			name:   "Service error",
			pageID: validPageID.String(),
			getPayloadBytes: func() []byte {
				payloadBytes, _ := json.Marshal(map[string]interface{}{
					"audiences": []map[string]interface{}{
						{
							"audience_id":   validAudience.String(),
							"audience_type": "user",
							"role":          "viewer",
						},
					},
				})
				return payloadBytes
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.On("BulkAddAudienceToPage",
					mock.Anything,
					validPageID,
					mock.Anything,
				).Return(nil, pages.BulkAddPageAudienceErrors{Error: errors.New("internal error")})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "internal error",
			},
		},
		{
			name:   "Invalid payload structure",
			pageID: validPageID.String(),
			getPayloadBytes: func() []byte {
				return []byte("invalid json payload")
			},
			setupMock:      func(m *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid character 'i' looking for beginning of value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router, mockSvc, serverCfg := setupRouter(t)
			mockDatasetService := mock_dataset_service.NewMockDatasetService(t)
			mockCacheService := mockcache.NewMockCacheClient(t)
			tt.setupMock(mockSvc)
			g := router.Group("/")
			registerRoutes(g, mockSvc, mockDatasetService, mockCacheService, serverCfg)
			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/pages/"+tt.pageID+"/audiences", bytes.NewBuffer(tt.getPayloadBytes()))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedBody != nil {
				assert.Equal(t, tt.expectedBody, response)
			}

			// Verify all expected mock calls were made
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUpdatePageAudience(t *testing.T) {
	audienceId := uuid.New()
	validPageID := uuid.New()

	tests := []struct {
		name            string
		pageId          string
		getPayloadBytes func() []byte
		setupMock       func(*mock_pages_service.MockPagesService)
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:   "successful update",
			pageId: validPageID.String(),
			getPayloadBytes: func() []byte {
				payload := dtos.UpdateAudienceRoleRequest{
					AudiencId: audienceId,
					Role:      "EDITOR",
				}
				payloadBytes, _ := json.Marshal(payload)
				return payloadBytes
			},
			setupMock: func(mockSvc *mock_pages_service.MockPagesService) {
				mockSvc.EXPECT().UpdatePageAudiencePrivilege(
					mock.Anything,
					validPageID,
					audienceId,
					models.ResourcePrivilege("EDITOR"),
				).Return(&models.ResourceAudiencePolicy{
					ResourceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"created_at":"0001-01-01T00:00:00Z","deleted_at":null,"metadata":null,"privilege":"","resource_audience_id":"00000000-0000-0000-0000-000000000000","resource_audience_policy_id":"00000000-0000-0000-0000-000000000000","resource_audience_type":"","resource_id":"123e4567-e89b-12d3-a456-426614174000","resource_type":"","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:   "invalid page ID",
			pageId: "invalid-uuid",
			getPayloadBytes: func() []byte {
				payload := dtos.UpdateAudienceRoleRequest{
					AudiencId: audienceId,
					Role:      "EDITOR",
				}
				payloadBytes, _ := json.Marshal(payload)
				return payloadBytes
			},
			setupMock:      func(mockSvc *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid page ID"}`,
		},
		{
			name:   "invalid payload",
			pageId: validPageID.String(),
			getPayloadBytes: func() []byte {
				return []byte("invalid json payload")
			},
			setupMock:      func(mockSvc *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:   "service error",
			pageId: validPageID.String(),
			getPayloadBytes: func() []byte {
				payload := dtos.UpdateAudienceRoleRequest{
					AudiencId: audienceId,
					Role:      "EDITOR",
				}
				payloadBytes, _ := json.Marshal(payload)
				return payloadBytes
			},
			setupMock: func(mockSvc *mock_pages_service.MockPagesService) {
				mockSvc.EXPECT().UpdatePageAudiencePrivilege(
					mock.Anything,
					validPageID,
					audienceId,
					models.ResourcePrivilege("EDITOR"),
				).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup
			router, mockSvc, serverCfg := setupRouter(t)
			g := router.Group("/")
			mockDatasetService := mock_dataset_service.NewMockDatasetService(t)
			mockCacheService := mockcache.NewMockCacheClient(t)
			// Setup mock expectations
			tt.setupMock(mockSvc)
			registerRoutes(g, mockSvc, mockDatasetService, mockCacheService, serverCfg)

			// Create request
			payloadBytes := tt.getPayloadBytes()
			req, _ := http.NewRequest(http.MethodPatch, "/pages/"+tt.pageId+"/audiences", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestDeletePageAudience(t *testing.T) {
	validPageID := uuid.New()
	validAudienceID := uuid.New()

	tests := []struct {
		name           string
		pageId         string
		payload        interface{}
		setupMock      func(*mock_pages_service.MockPagesService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			pageId: validPageID.String(),
			payload: dtos.DeleteAudienceRoleRequest{
				AudiencId: validAudienceID,
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().RemoveAudienceFromPage(mock.Anything, validPageID, validAudienceID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:   "Invalid page ID",
			pageId: "invalid-uuid",
			payload: dtos.DeleteAudienceRoleRequest{
				AudiencId: validAudienceID,
			},
			setupMock:      func(m *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid page ID"}`,
		},
		{
			name:           "Invalid request body",
			pageId:         validPageID.String(),
			payload:        interface{}("invalid payload"),
			setupMock:      func(m *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"json: cannot unmarshal string into Go value of type dtos.DeleteAudienceRoleRequest"}`,
		},
		{
			name:   "Service error",
			pageId: validPageID.String(),
			payload: dtos.DeleteAudienceRoleRequest{
				AudiencId: validAudienceID,
			},
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().RemoveAudienceFromPage(mock.Anything, validPageID, validAudienceID).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router, mockSvc, serverCfg := setupRouter(t)
			g := router.Group("/")

			mockDatasetService := mock_dataset_service.NewMockDatasetService(t)
			mockCacheService := mockcache.NewMockCacheClient(t)
			// Setup mock expectations
			tt.setupMock(mockSvc)

			registerRoutes(g, mockSvc, mockDatasetService, mockCacheService, serverCfg)

			// Create request
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodDelete, "/pages/"+tt.pageId+"/audiences", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}
func TestHandleGetPagesByOrganizationId(t *testing.T) {
	validOrganizationID := uuid.New()
	validPageID := uuid.New()
	now := time.Now()

	validPages := []models.Page{
		{
			ID:              validPageID,
			Name:            "Test Page",
			Description:     nil,
			CreatedAt:       now,
			UpdatedAt:       now,
			OrganizationId:  validOrganizationID,
			FractionalIndex: 0,
		},
	}

	tests := []struct {
		name           string
		organizationId string
		setupMock      func(*mock_pages_service.MockPagesService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful retrieval",
			organizationId: validOrganizationID.String(),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPagesByOrganizationId(mock.Anything, validOrganizationID).
					Return(validPages, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`[{
				"page_id": "%s",
				"name": "Test Page",
				"description": null,
				"created_at": "%s",
				"updated_at": "%s", 
				"organization_id": "%s",
				"fractional_index": 0
			}]`, validPageID, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), validOrganizationID),
		},
		{
			name:           "invalid organization ID",
			organizationId: "invalid-uuid",
			setupMock:      func(m *mock_pages_service.MockPagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid organization ID"}`,
		},
		{
			name:           "service error",
			organizationId: validOrganizationID.String(),
			setupMock: func(m *mock_pages_service.MockPagesService) {
				m.EXPECT().
					GetPagesByOrganizationId(mock.Anything, validOrganizationID).
					Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router, mockSvc, serverCfg := setupRouter(t)
			g := router.Group("/")

			mockDatasetService := mock_dataset_service.NewMockDatasetService(t)
			mockCacheService := mockcache.NewMockCacheClient(t)
			// Setup mock expectations
			tt.setupMock(mockSvc)

			registerRoutes(g, mockSvc, mockDatasetService, mockCacheService, serverCfg)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/pages/get-pages-by-organization-id?organizationId=%s", tt.organizationId), nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}
