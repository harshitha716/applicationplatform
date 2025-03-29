package middleware

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/db/models"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/Zampfi/application-platform/services/api/pkg/kratosclient"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getMockServerConfig(statusCode int, responseBody string) (*serverconfig.ServerConfig, *httptest.Server) {

	mockKratosServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return a resposne with write status code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		io.WriteString(w, responseBody)
	}))

	serverconfig := serverconfig.GetEmptyServerConfig()
	serverconfig.Env.Port = "8080"
	serverconfig.Env.PantheonURL = "http://localhost:8100"
	serverconfig.Env.AuthBaseUrl = mockKratosServer.URL

	authClient, _ := kratosclient.NewClient(mockKratosServer.URL)

	serverconfig.AuthClient = authClient
	return serverconfig, mockKratosServer
}

type testCase struct {
	name             string
	statusCode       int
	getMockServerCfg func() *serverconfig.ServerConfig
	assertResponse   func(t *testing.T, w *httptest.ResponseRecorder)
	headers          map[string]string
}

const SESSION_MOCK_INVALID_EMAIL = `
	{
	  "id": "fc59625f-8ad6-491a-b3e0-f5676a9232f3",
	  "active": true,
	  "expires_at": "2024-08-02T21:13:10.021766508Z",
	  "authenticated_at": "2024-08-01T21:13:10.021766508Z",
	  "authenticator_assurance_level": "aal1",
	  "authentication_methods": [
		  {
			  "method": "password",
			  "aal": "aal1",
			  "completed_at": "2024-08-01T21:13:10.021757758Z"
		  }
	  ],
	  "issued_at": "2024-08-01T21:13:10.021766508Z",
	  "identity": {
		  "id": "d93e9fb4-2451-4eb5-aa86-aaf017c74c39",
		  "schema_id": "default",
		  "schema_url": "http://auth:4433/schemas/ZGVmYXVsdA",
		  "state": "active",
		  "state_changed_at": "2024-08-01T19:59:14.935789Z",
		  "traits": {
			  "email": "admin@gmail.com"
			}
		}
	}
`

const SESSION_MOCK = `
	{
	  "id": "fc59625f-8ad6-491a-b3e0-f5676a9232f3",
	  "active": true,
	  "expires_at": "2024-08-02T21:13:10.021766508Z",
	  "authenticated_at": "2024-08-01T21:13:10.021766508Z",
	  "authenticator_assurance_level": "aal1",
	  "authentication_methods": [
		  {
			  "method": "password",
			  "aal": "aal1",
			  "completed_at": "2024-08-01T21:13:10.021757758Z"
		  }
	  ],
	  "issued_at": "2024-08-01T21:13:10.021766508Z",
	  "identity": {
		  "id": "d93e9fb4-2451-4eb5-aa86-aaf017c74c39",
		  "schema_id": "default",
		  "schema_url": "http://auth:4433/schemas/ZGVmYXVsdA",
		  "state": "active",
		  "state_changed_at": "2024-08-01T19:59:14.935789Z",
		  "traits": {
			  "email": "admin@zamp.ai"
			}
		}
	}
`

func TestAuthMiddleware(t *testing.T) {
	testCases := []testCase{
		{
			name:       "Unauthorized",
			statusCode: 401,
			getMockServerCfg: func() *serverconfig.ServerConfig {
				serverCfg, _ := getMockServerConfig(401, `{"message":"unauthorized"}`)

				serverCfg.Store = mock_store.NewMockStore(t)
				return serverCfg
			},
			assertResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, 401, w.Code)
				assert.Equal(t, `{"message":"unauthorized"}`, w.Body.String())
			},
		},
		{
			name:       "Success",
			statusCode: 200,
			getMockServerCfg: func() *serverconfig.ServerConfig {
				serverCfg, _ := getMockServerConfig(200, SESSION_MOCK)
				mockStore := mock_store.NewMockStore(t)
				orgId, _ := uuid.Parse("be166699-eeea-4c8a-a3ec-107764dc3e91")
				mockStore.EXPECT().GetOrganizationsByMemberId(mock.Anything, mock.Anything).Return([]models.Organization{{ID: orgId}}, nil)
				serverCfg.Store = mockStore
				return serverCfg
			},
			assertResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, 200, w.Code)
			},
		},
		{
			name:       "Success with X-Zamp-Organization-Id header",
			statusCode: 200,
			getMockServerCfg: func() *serverconfig.ServerConfig {
				serverCfg, _ := getMockServerConfig(200, SESSION_MOCK)
				mockStore := mock_store.NewMockStore(t)
				orgId, _ := uuid.Parse("be166699-eeea-4c8a-a3ec-107764dc3e91")
				mockStore.EXPECT().GetOrganizationsByMemberId(mock.Anything, mock.Anything).Return([]models.Organization{{ID: orgId}}, nil)
				serverCfg.Store = mockStore
				return serverCfg
			},
			assertResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, 200, w.Code)
			},
			headers: map[string]string{
				"X-Zamp-Organization-Id": "be166699-eeea-4c8a-a3ec-107764dc3e91",
			},
		},
		{
			name:       "Success with no organizations",
			statusCode: 200,
			getMockServerCfg: func() *serverconfig.ServerConfig {
				serverCfg, _ := getMockServerConfig(200, SESSION_MOCK)
				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetOrganizationsByMemberId(mock.Anything, mock.Anything).Return([]models.Organization{}, nil)
				serverCfg.Store = mockStore
				return serverCfg
			},
			assertResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, 200, w.Code)
			},
		},
		{
			name:       "Success with invalid organization ID in header",
			statusCode: 200,
			getMockServerCfg: func() *serverconfig.ServerConfig {
				serverCfg, _ := getMockServerConfig(200, SESSION_MOCK)
				mockStore := mock_store.NewMockStore(t)
				orgId, _ := uuid.Parse("be166699-eeea-4c8a-a3ec-107764dc3e91")
				mockStore.EXPECT().GetOrganizationsByMemberId(mock.Anything, mock.Anything).Return([]models.Organization{{ID: orgId}}, nil)
				serverCfg.Store = mockStore
				return serverCfg
			},
			assertResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, 200, w.Code)
			},
			headers: map[string]string{
				"X-Zamp-Organization-Id": "11111111-1111-1111-1111-111111111111",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			serverCfg := tc.getMockServerCfg()

			gin.SetMode(gin.TestMode)
			r := gin.Default()
			authMiddleWare, err := GetAuthMiddleware(serverCfg)
			assert.Nil(t, err)
			r.Use(authMiddleWare)
			r.GET("/route", func(c *gin.Context) {
				c.JSON(200, gin.H{})
			})

			w := httptest.NewRecorder()

			// create request with invalid json
			req, _ := http.NewRequest("GET", "/route", strings.NewReader(`{"test": "json"}`))
			req.Header.Set("Content-Type", "application/json")

			// Add any additional headers
			if tc.headers != nil {
				for key, value := range tc.headers {
					req.Header.Set(key, value)
				}
			}

			r.ServeHTTP(w, req)

			tc.assertResponse(t, w)

		})
	}
}

func TestBasicAuthMiddleware(t *testing.T) {
	requiredUsername := "testuser"
	requiredPassword := "testpass"

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "Valid credentials",
			authHeader:         "Basic " + base64.StdEncoding.EncodeToString([]byte(requiredUsername+":"+requiredPassword)),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Success",
		},
		{
			name:               "Missing Authorization header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Authorization header required"}`,
		},
		{
			name:               "Invalid Authorization header format",
			authHeader:         "Bearer token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Invalid Authorization header format"}`,
		},
		{
			name:               "Invalid Base64 encoding",
			authHeader:         "Basic invalid_base64",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Invalid Base64 encoding"}`,
		},
		{
			name:               "Invalid credentials",
			authHeader:         "Basic " + base64.StdEncoding.EncodeToString([]byte("wronguser:wrongpass")),
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Invalid username or password"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(BasicAuthMiddleware(requiredUsername, requiredPassword))
			r.GET("/protected", func(c *gin.Context) {
				c.String(http.StatusOK, "Success")
			})

			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
