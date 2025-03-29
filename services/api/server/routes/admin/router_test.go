package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	mock_auth "github.com/Zampfi/application-platform/services/api/mocks/core/auth"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mock_templates "github.com/Zampfi/application-platform/services/api/mocks/server/routes/admin/templates"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterAdminRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		routeName      string
		data           interface{}
		setupMocks     func(*mock_auth.MockAuthService, *mock_templates.MockTemplateLoader)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "GET request renders template",
			method:    "GET",
			path:      "/test",
			routeName: "Test Route",
			data:      struct{}{},
			setupMocks: func(mockAuth *mock_auth.MockAuthService, mockTemplate *mock_templates.MockTemplateLoader) {
				mockTemplate.EXPECT().ExecuteTemplate(mock.Anything, "base-form.html", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "GET login request renders login template",
			method:    "GET",
			path:      "/login",
			routeName: "Login",
			data:      struct{}{},
			setupMocks: func(mockAuth *mock_auth.MockAuthService, mockTemplate *mock_templates.MockTemplateLoader) {
				mockTemplate.On("ExecuteTemplate", mock.Anything, "admin-home.html", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "POST request with non-admin role returns forbidden",
			method:    "POST",
			path:      "/test",
			routeName: "Test Route",
			data:      struct{}{},
			setupMocks: func(mockAuth *mock_auth.MockAuthService, mockTemplate *mock_templates.MockTemplateLoader) {
				mockAuth.On("ResolveAdminInfo", mock.Anything, mock.Anything).Return("user", uuid.New(), []uuid.UUID{})
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   map[string]interface{}{"error": "forbidden"},
		},
		{
			name:      "Invalid HTTP method returns method not allowed",
			method:    "PUT",
			path:      "/test",
			routeName: "Test Route",
			data:      struct{}{},
			setupMocks: func(mockAuth *mock_auth.MockAuthService, mockTemplate *mock_templates.MockTemplateLoader) {
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   map[string]interface{}{"error": "method not allowed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			mockAuth := mock_auth.NewMockAuthService(t)
			mockTemplate := mock_templates.NewMockTemplateLoader(t)

			tt.setupMocks(mockAuth, mockTemplate)

			controller := &adminRouteController{
				routeGroup:     router.Group("/admin"),
				templateLoader: mockTemplate,
				authSvc:        mockAuth,
				serverCfg:      &serverconfig.ServerConfig{Env: &serverconfig.ConfigVariables{Environment: "test"}},
			}

			controller.registerAdminRoute(tt.path, tt.routeName, tt.data, func(c *gin.Context) {})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "/admin"+tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockAuth.AssertExpectations(t)
			mockTemplate.AssertExpectations(t)
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		payload        interface{}
		setupMocks     func(*mock_auth.MockAuthService, *mock_templates.MockTemplateLoader, *mock_store.MockStore)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Login endpoint POST returns method not allowed",
			method: "POST",
			path:   "/admin/login",
			setupMocks: func(mockAuth *mock_auth.MockAuthService, mockTemplate *mock_templates.MockTemplateLoader, mockStore *mock_store.MockStore) {
				mockAuth.On("ResolveAdminInfo", mock.Anything, mock.Anything).Return("admin", uuid.New(), []uuid.UUID{})
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   map[string]interface{}{"error": "method not allowed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			mockAuth := mock_auth.NewMockAuthService(t)
			mockTemplate := mock_templates.NewMockTemplateLoader(t)
			mockStore := mock_store.NewMockStore(t)

			tt.setupMocks(mockAuth, mockTemplate, mockStore)

			serverCfg := &serverconfig.ServerConfig{
				Store: mockStore,
			}

			registerRoutes(router, serverCfg, mockTemplate, mockAuth)

			w := httptest.NewRecorder()
			var req *http.Request
			if tt.payload != nil {
				payloadBytes, _ := json.Marshal(tt.payload)
				req, _ = http.NewRequest(tt.method, tt.path, bytes.NewBuffer(payloadBytes))
			} else {
				req, _ = http.NewRequest(tt.method, tt.path, nil)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody["message"], response["message"])
			}

			mockAuth.AssertExpectations(t)
			mockTemplate.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}
