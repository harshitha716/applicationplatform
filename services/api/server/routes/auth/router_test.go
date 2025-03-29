package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/db/models"
	mock_auth "github.com/Zampfi/application-platform/services/api/mocks/core/auth"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CustomResponseRecorder struct {
	*httptest.ResponseRecorder
}

// CloseNotify implements http.CloseNotifier
func (r *CustomResponseRecorder) CloseNotify() <-chan bool {
	return make(<-chan bool)
}

func TestRegisterRoutes_Error(t *testing.T) {
	svCfg := serverconfig.GetEmptyServerConfig()
	e := gin.New()
	err := RegisterAuthRoutes(e, svCfg)
	assert.NotNil(t, err)
}

func TestRegisterRoutes_Success(t *testing.T) {
	svCfg := serverconfig.GetEmptyServerConfig()
	svCfg.Env.AuthBaseUrl = "http://localhost:4433"
	e := gin.New()
	err := RegisterAuthRoutes(e, svCfg)
	assert.Nil(t, err)
}

type testCase struct {
	name          string
	method        string
	path          string
	skip          bool
	outputPayload string
	statusCode    int
	headers       map[string]string
	initServerCfg func() *serverconfig.ServerConfig
}

func TestHandleKratosProxyRequest(t *testing.T) {

	testCases := []testCase{
		{
			name:          "Accessing restricted route returns 404",
			method:        "GET",
			path:          "/auth/relay/admin/identities",
			outputPayload: `{"message":"not found"}`,
			statusCode:    404,
			headers:       nil,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()
				svCfg.Env.AuthBaseUrl = "http://localhost:4433"
				return svCfg
			},
		},
		{
			name:          "Proxies correctly",
			method:        "GET",
			path:          "/auth/relay/self-service/login",
			outputPayload: `{"message":"success"}`,
			statusCode:    200,
			headers:       nil,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					w.Write([]byte(`{"message":"success"}`))
				}))

				svCfg.Env.AuthBaseUrl = authServer.URL
				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Arrange
			svCfg := tc.initServerCfg()
			// Act

			// register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			err := RegisterAuthRoutes(e, svCfg)
			assert.Nil(t, err)

			// Create a test request
			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// fire the request
			w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, tc.outputPayload, w.Body.String())
		})
	}
}

const sessionMock = `
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
              "email": "rishichandra1@gmail.com"
          },
          "metadata_public": null,
          "created_at": "2024-08-01T19:59:14.937324Z",
          "updated_at": "2024-08-01T19:59:14.937324Z",
          "organization_id": null
      },
      "devices": [
          {
              "id": "3da58b8b-1f62-47aa-9a0e-e506470efc20",
              "ip_address": "192.168.65.1:58625",
              "user_agent": "PostmanRuntime/7.40.0",
              "location": ""
          }
      ]
  }
`

func TestGetUserInfoWithOrganizations(t *testing.T) {

	testCases := []testCase{
		{
			name:          "Unauthorized request returns 401",
			skip:          true,
			method:        "GET",
			path:          "/auth/whoami",
			outputPayload: `{"message":"unauthorized"}`,
			statusCode:    401,
			headers:       nil,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(401)
					w.Write([]byte(`{"message":"unauthorized"}`))
				}))

				svCfg.Env.AuthBaseUrl = authServer.URL
				return svCfg
			},
		},
		{
			name:          "Authorized request returns 200",
			method:        "GET",
			skip:          false,
			path:          "/auth/whoami",
			outputPayload: `{"orgs":[{"organization_id":"357623aa-4945-4310-8714-523acd54d4b7","name":"test","description":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","owner_id":"00000000-0000-0000-0000-000000000000","resource_audience_policies":null,"invitations":null,"membership_requests":null,"sso_configs":null}],"user_email":"rishichandra1@gmail.com","user_id":"d93e9fb4-2451-4eb5-aa86-aaf017c74c39"}`,
			statusCode:    200,
			headers:       nil,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				orgId, _ := uuid.Parse("357623aa-4945-4310-8714-523acd54d4b7")

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetOrganizationsByMemberId(mock.Anything, mock.Anything).Return([]models.Organization{{
					ID:          orgId,
					Name:        "test",
					Description: nil,
				}}, nil)

				authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Header()["Content-Type"] = []string{"application/json"}
					w.WriteHeader(200)
					w.Write([]byte(sessionMock))
				}))

				svCfg.Env.AuthBaseUrl = authServer.URL
				svCfg.Store = mockStore
				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.skip {
				t.Skip()
			}

			// Arrange
			svCfg := tc.initServerCfg()
			// Act

			// register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			err := RegisterAuthRoutes(e, svCfg)
			assert.Nil(t, err)

			// Create a test request
			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// fire the request
			w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, tc.outputPayload, w.Body.String())
		})
	}
}

func TestGetLoginFlow(t *testing.T) {
	testCases := []struct {
		name          string
		skip          bool
		path          string
		body          string
		statusCode    int
		outputPayload string
		setupMock     func(*mock_auth.MockAuthService)
	}{
		{
			name:          "should return 200 and login flow when valid email",
			path:          "/auth/login/flow/create",
			body:          `{"email": "test@example.com"}`,
			statusCode:    http.StatusOK,
			outputPayload: `{"flow_id":"test-flow-id"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
				mockAuthService.EXPECT().GetAuthFlowForUser(mock.Anything, "test@example.com").Return(&kratos.LoginFlow{Id: "test-flow-id"}, nil, nil)

			},
		},
		{
			name:          "should return 400 when invalid json",
			path:          "/auth/login/flow/create",
			body:          `invalid json`,
			statusCode:    http.StatusBadRequest,
			outputPayload: `{"error":"invalid character 'i' looking for beginning of value"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
			},
		},
		{
			name:          "should return 400 when auth service returns error",
			path:          "/auth/login/flow/create",
			body:          `{"email": "test@example.com"}`,
			statusCode:    http.StatusBadRequest,
			outputPayload: `{"error":"auth service error"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
				mockAuthService.EXPECT().GetAuthFlowForUser(mock.Anything, "test@example.com").Return(nil, nil, errors.New("auth service error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			authService := mock_auth.NewMockAuthService(t)
			tc.setupMock(authService)

			// register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			registerRoutes(e, authService)

			// Create a test request
			req, err := http.NewRequest("POST", tc.path, strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// fire the request
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestHandleKratosAfterRegistrationWebhookEvent(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		body          string
		statusCode    int
		outputPayload string
		setupMock     func(mockAuthService *mock_auth.MockAuthService)
	}{
		{
			name:          "should return 400 when request body is invalid json",
			path:          "/auth/internal/webhook",
			body:          `invalid json`,
			statusCode:    http.StatusBadRequest,
			outputPayload: `{"error":"invalid json"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
			},
		},
		{
			name:          "should return 400 when user id is invalid",
			path:          "/auth/internal/webhook",
			body:          `{"user_id": "invalid-uuid"}`,
			statusCode:    http.StatusBadRequest,
			outputPayload: `{"error":"invalid user id"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
			},
		},
		{
			name:          "should return 400 when auth service returns error",
			path:          "/auth/internal/webhook",
			body:          `{"user_id": "d93e9fb4-2451-4eb5-aa86-aaf017c74c39"}`,
			statusCode:    http.StatusBadRequest,
			outputPayload: `{"error":"unauthorized"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
				userId, _ := uuid.Parse("d93e9fb4-2451-4eb5-aa86-aaf017c74c39")
				mockAuthService.EXPECT().HandleNewUserCreated(mock.Anything, "", userId).Return(errors.New("unauthorized"))
			},
		},
		{
			name:          "should return 200 when request is valid",
			path:          "/auth/internal/webhook",
			body:          `{"user_id": "d93e9fb4-2451-4eb5-aa86-aaf017c74c39"}`,
			statusCode:    http.StatusOK,
			outputPayload: `{"message":"success"}`,
			setupMock: func(mockAuthService *mock_auth.MockAuthService) {
				userId, _ := uuid.Parse("d93e9fb4-2451-4eb5-aa86-aaf017c74c39")
				mockAuthService.EXPECT().HandleNewUserCreated(mock.Anything, "", userId).Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := mock_auth.NewMockAuthService(t)
			tc.setupMock(authService)

			// register routes
			e := gin.New()
			gin.SetMode(gin.TestMode)
			registerRoutes(e, authService)

			// Create a test request
			req, err := http.NewRequest("POST", tc.path, strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// fire the request
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, tc.outputPayload, strings.TrimSpace(w.Body.String()))
		})
	}
}
