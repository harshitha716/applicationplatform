package organizations

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/organizations"
	"github.com/Zampfi/application-platform/services/api/core/organizations/teams"
	"github.com/Zampfi/application-platform/services/api/db/models"
	mockOrganization "github.com/Zampfi/application-platform/services/api/mocks/core/organizations"
	mock_teams "github.com/Zampfi/application-platform/services/api/mocks/core/organizations/teams"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	dtos "github.com/Zampfi/application-platform/services/api/server/routes/organizations/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type testCase struct {
	name          string
	method        string
	path          string
	skip          bool
	outputPayload string
	statusCode    int
	initServerCfg func() *serverconfig.ServerConfig
}

func TestGetOrganizations(t *testing.T) {
	testCases := []testCase{
		{
			name:          "Get organizations db error",
			method:        "GET",
			skip:          false,
			path:          "/organizations/",
			outputPayload: `{"error":"something went wrong"}`,
			statusCode:    http.StatusInternalServerError,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetOrganizationsAll(mock.Anything).Return([]models.Organization{}, errors.New("something went wrong"))

				svCfg.Store = mockStore

				return svCfg
			},
		},
		{
			name:          "Get organizations success",
			method:        "GET",
			skip:          false,
			path:          "/organizations/",
			outputPayload: `[{"organization_id":"00000000-0000-0000-0000-000000000000","name":"test","description":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","owner_id":"00000000-0000-0000-0000-000000000000","resource_audience_policies":null,"invitations":null,"membership_requests":null,"sso_configs":null}]`,
			statusCode:    http.StatusOK,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				orgId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetOrganizationsAll(mock.Anything).Return([]models.Organization{{ID: orgId, Name: "test", Description: nil}}, nil)

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
			g := e.Group("/")
			RegisterOrganizationRoutes(g, svCfg)

			// Create a test request
			req, err := http.NewRequest(tc.method, tc.path, nil)
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

func TestGetOrganizationAudiences(t *testing.T) {

	orgId := uuid.New()

	policies := []models.ResourceAudiencePolicy{{ResourceID: orgId}}
	policiesRaw, _ := json.Marshal(policies)

	testCases := []testCase{
		{
			name:          "Get organizations audiences db error",
			method:        "GET",
			skip:          false,
			path:          "/organizations/" + orgId.String() + "/audiences",
			outputPayload: `{"error":"something went wrong"}`,
			statusCode:    http.StatusInternalServerError,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetOrganizationPolicies(mock.Anything, orgId).Return([]models.ResourceAudiencePolicy{}, errors.New("db error"))

				svCfg.Store = mockStore

				return svCfg
			},
		},
		{
			name:          "Get organizations bad org ID",
			method:        "GET",
			skip:          false,
			path:          "/organizations/" + "notuuid" + "/audiences",
			outputPayload: `{"error":"invalid organization id"}`,
			statusCode:    http.StatusBadRequest,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)

				svCfg.Store = mockStore

				return svCfg
			},
		},
		{
			name:          "Get organizations success",
			method:        "GET",
			skip:          false,
			path:          "/organizations/" + orgId.String() + "/audiences",
			outputPayload: string(policiesRaw),
			statusCode:    http.StatusOK,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetOrganizationPolicies(mock.Anything, orgId).Return(policies, nil)

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
			g := e.Group("/")
			RegisterOrganizationRoutes(g, svCfg)

			// Create a test request
			req, err := http.NewRequest(tc.method, tc.path, nil)
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

func TestBulkInviteMembers(t *testing.T) {
	validOrgId := uuid.New()

	tests := []struct {
		name           string
		orgId          string
		payload        interface{}
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "successful bulk invite",
			orgId: validOrgId.String(),
			payload: map[string]interface{}{
				"invitations": []map[string]string{
					{"email": "test1@example.com", "role": "MEMBER"},
					{"email": "test2@example.com", "role": "ADMIN"},
				},
			},
			setupMock: func(m *mockOrganization.MockOrganizationService) {
				m.On("BulkInviteMembers",
					mock.Anything,
					validOrgId,
					mock.AnythingOfType("organizations.BulkInvitationPayload"),
				).Return([]models.OrganizationInvitation{
					{TargetEmail: "test1@example.com", Privilege: "MEMBER"},
					{TargetEmail: "test2@example.com", Privilege: "ADMIN"},
				}, organizations.BulkInvitationError{})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"invitations":       []interface{}{},
				"invitation_errors": []interface{}{},
			},
		},
		{
			name:  "invalid organization id",
			orgId: "invalid-uuid",
			payload: map[string]interface{}{
				"invitations": []map[string]string{
					{"email": "test@example.com", "role": "MEMBER"},
				},
			},
			setupMock:      func(m *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "invalid request body",
			orgId:          validOrgId.String(),
			payload:        "invalid json",
			setupMock:      func(m *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request body",
			},
		},
		{
			name:  "service error",
			orgId: validOrgId.String(),
			payload: map[string]interface{}{
				"invitations": []map[string]string{
					{"email": "test@example.com", "role": "MEMBER"},
				},
			},
			setupMock: func(m *mockOrganization.MockOrganizationService) {
				m.On("BulkInviteMembers",
					mock.Anything,
					validOrgId,
					mock.AnythingOfType("organizations.BulkInvitationPayload"),
				).Return([]models.OrganizationInvitation{}, organizations.BulkInvitationError{
					Error: assert.AnError,
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockOrganization.MockOrganizationService)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			// Setup mock if provided
			tt.setupMock(mockService)

			// Prepare request
			var jsonPayload []byte
			var err error
			if str, ok := tt.payload.(string); ok {
				jsonPayload = []byte(str)
			} else {
				jsonPayload, err = json.Marshal(tt.payload)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/organizations/"+tt.orgId+"/audiences/invitations", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check if response contains expected fields
			for key, expectedValue := range tt.expectedBody {
				assert.Contains(t, response, key)
				if key != "invitations" && key != "invitation_errors" {
					assert.Equal(t, expectedValue, response[key])
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetOrganizationInvitations(t *testing.T) {

	tests := []struct {
		name           string
		orgID          string
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "successful retrieval",
			orgID: "123e4567-e89b-12d3-a456-426614174000",
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				expectedInvitations := []models.OrganizationInvitation{
					{
						OrganizationInvitationID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
						TargetEmail:              "test@example.com",
						Privilege:                models.PrivilegeOrganizationMember,
					},
				}
				mockService.EXPECT().GetAllOrganizationInvitations(
					mock.Anything,
					uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				).Return(expectedInvitations, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{ // Changed type here
				map[string]interface{}{
					"organization_invitation_id": "123e4567-e89b-12d3-a456-426614174001",
					"organization_id":            "00000000-0000-0000-0000-000000000000",
					"email":                      "test@example.com",
					"privilege":                  "member",
					"created_at":                 "0001-01-01T00:00:00Z",
					"updated_at":                 "0001-01-01T00:00:00Z",
					"deleted_at":                 nil,
					"invited_by":                 "00000000-0000-0000-0000-000000000000",
					"email_retry_count":          float64(0),
					"email_sent_at":              nil,
				},
			},
		},
		{
			name:  "invalid organization ID",
			orgID: "invalid-uuid",
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{ // Changed from gin.H to map[string]interface{}
				"error": "invalid organization id",
			},
		},
		{
			name:  "service error",
			orgID: "123e4567-e89b-12d3-a456-426614174000",
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().GetAllOrganizationInvitations(
					mock.Anything,
					uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{ // Changed from gin.H to map[string]interface{}
				"error": "something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/organizations/"+tt.orgID+"/audiences/invitations", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestRemoveOrganizationMember(t *testing.T) {
	validOrgID := uuid.New()
	validUserID := uuid.New()

	tests := []struct {
		name           string
		orgID          string
		payload        interface{}
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "success",
			orgID: validOrgID.String(),
			payload: map[string]interface{}{
				"user_id": validUserID,
			},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().RemoveOrganizationMember(
					mock.Anything,
					validOrgID,
					validUserID,
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "member removed successfully",
			},
		},
		{
			name:  "invalid organization id",
			orgID: "invalid-uuid",
			payload: map[string]interface{}{
				"user_id": validUserID,
			},
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "invalid request body",
			orgID:          validOrgID.String(),
			payload:        "invalid json",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request body",
			},
		},
		{
			name:  "service error",
			orgID: validOrgID.String(),
			payload: map[string]interface{}{
				"user_id": validUserID,
			},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().RemoveOrganizationMember(
					mock.Anything,
					validOrgID,
					validUserID,
				).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()

			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("DELETE", "/organizations/"+tt.orgID+"/audiences", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}

}

func TestGetOrganizationMembershipRequests(t *testing.T) {
	validUserID := uuid.New()
	validRequest := models.OrganizationMembershipRequest{
		ID:             uuid.New(),
		OrganizationID: uuid.New(),
		UserID:         validUserID,
		User: models.User{
			ID:    validUserID,
			Email: "test@example.com",
			Name:  "",
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().GetOrganizationMembershipRequestsAll(
					mock.Anything,
				).Return([]models.OrganizationMembershipRequest{validRequest}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":              validRequest.ID.String(),
					"organization_id": validRequest.OrganizationID.String(),
					"user_id":         validRequest.UserID.String(),
					"created_at":      "0001-01-01T00:00:00Z",
					"updated_at":      "0001-01-01T00:00:00Z",
					"deleted_at":      "0001-01-01T00:00:00Z",
					"status":          "",
					"user": map[string]interface{}{
						"user_id": validRequest.User.ID.String(),
						"email":   validRequest.User.Email,
						"name":    validRequest.User.Name,
					},
				},
			},
		},
		{
			name: "service error",
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().GetOrganizationMembershipRequestsAll(
					mock.Anything,
				).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/organizations/membership-requests", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetOrganizationMembershipRequestsByOrganizationId(t *testing.T) {
	validOrgID := uuid.New()
	validRequest := models.OrganizationMembershipRequest{
		ID:             uuid.New(),
		OrganizationID: validOrgID,
		User: models.User{
			ID:    uuid.Nil,
			Email: "",
			Name:  "",
		},
	}

	tests := []struct {
		name           string
		orgID          string
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "success",
			orgID: validOrgID.String(),
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().GetOrganizationMembershipRequestsByOrganizationId(
					mock.Anything,
					validOrgID,
				).Return([]models.OrganizationMembershipRequest{validRequest}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":              validRequest.ID.String(),
					"organization_id": validRequest.OrganizationID.String(),
					"user_id":         "00000000-0000-0000-0000-000000000000",
					"created_at":      "0001-01-01T00:00:00Z",
					"updated_at":      "0001-01-01T00:00:00Z",
					"deleted_at":      "0001-01-01T00:00:00Z",
					"status":          "",
					"user": map[string]interface{}{
						"user_id": validRequest.User.ID.String(),
						"email":   validRequest.User.Email,
						"name":    validRequest.User.Name,
					},
				},
			},
		},
		{
			name:           "invalid org id",
			orgID:          "invalid-uuid",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:  "service error",
			orgID: validOrgID.String(),
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().GetOrganizationMembershipRequestsByOrganizationId(
					mock.Anything,
					validOrgID,
				).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/organizations/"+tt.orgID+"/audiences/requests", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestApproveOrganizationMembershipRequest(t *testing.T) {
	validOrgID := uuid.New()
	validUserID := uuid.New()
	validRequest := &models.OrganizationMembershipRequest{
		ID:             uuid.New(),
		OrganizationID: validOrgID,
		UserID:         validUserID,
		Status:         models.OrgMembershipStatusApproved,
		User: models.User{
			ID:    uuid.Nil,
			Email: "",
			Name:  "",
		},
	}

	tests := []struct {
		name           string
		orgID          string
		requestBody    interface{}
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "success",
			orgID: validOrgID.String(),
			requestBody: dtos.ApproveOrganizationMembershipRequestRequest{
				UserId: validUserID,
			},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().ApprovePendingOrganizationMembershipRequest(
					mock.Anything,
					validOrgID,
					validUserID,
				).Return(validRequest, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":              validRequest.ID.String(),
				"organization_id": validRequest.OrganizationID.String(),
				"user_id":         validRequest.UserID.String(),
				"status":          string(validRequest.Status),
				"created_at":      validRequest.CreatedAt.UTC().Format(time.RFC3339),
				"updated_at":      validRequest.UpdatedAt.UTC().Format(time.RFC3339),
				"deleted_at":      validRequest.DeletedAt.UTC().Format(time.RFC3339),
				"user": map[string]interface{}{
					"user_id": validRequest.User.ID.String(),
					"email":   validRequest.User.Email,
					"name":    validRequest.User.Name,
				},
			},
		},
		{
			name:  "invalid organization id",
			orgID: "invalid-uuid",
			requestBody: dtos.ApproveOrganizationMembershipRequestRequest{
				UserId: validUserID,
			},
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "invalid request body",
			orgID:          validOrgID.String(),
			requestBody:    "invalid json",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request body",
			},
		},
		{
			name:  "service error",
			orgID: validOrgID.String(),
			requestBody: dtos.ApproveOrganizationMembershipRequestRequest{
				UserId: validUserID,
			},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockService.EXPECT().ApprovePendingOrganizationMembershipRequest(
					mock.Anything,
					validOrgID,
					validUserID,
				).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, _ := http.NewRequest("PATCH", "/organizations/"+tt.orgID+"/audiences/requests/approve", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetTeams(t *testing.T) {
	validOrgID := uuid.New()
	validTeam := models.Team{
		TeamID: uuid.New(),
		Name:   "Test Team",
	}

	tests := []struct {
		name           string
		orgID          string
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Invalid org ID",
			orgID:          "invalid-uuid",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:  "Service error",
			orgID: validOrgID.String(),
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().GetTeamsByOrganizationID(mock.Anything, validOrgID).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service error",
			},
		},
		{
			name:  "Success",
			orgID: validOrgID.String(),
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().GetTeamsByOrganizationID(mock.Anything, validOrgID).Return([]models.Team{validTeam}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"team_id":          validTeam.TeamID.String(),
					"organization_id":  validTeam.OrganizationID.String(),
					"name":             validTeam.Name,
					"description":      validTeam.Description,
					"created_at":       validTeam.CreatedAt.Format(time.RFC3339),
					"updated_at":       validTeam.UpdatedAt.Format(time.RFC3339),
					"deleted_at":       nil,
					"metadata":         nil,
					"team_memberships": nil,
					"created_by":       validTeam.CreatedBy.String(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/organizations/"+tt.orgID+"/teams", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateTeam(t *testing.T) {
	validOrgID := uuid.New()
	validTeam := models.Team{
		TeamID: uuid.New(),
		Name:   "Test Team",
	}
	validPayload := teams.CreateTeamPayload{
		Name: "Test Team",
	}

	tests := []struct {
		name           string
		orgID          string
		requestBody    interface{}
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Invalid org ID",
			orgID:          "invalid-uuid",
			requestBody:    validPayload,
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "Invalid request body",
			orgID:          validOrgID.String(),
			requestBody:    "invalid json",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request body",
			},
		},
		{
			name:        "Service error",
			orgID:       validOrgID.String(),
			requestBody: validPayload,
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().CreateTeam(mock.Anything, validOrgID, validPayload).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service error",
			},
		},
		{
			name:        "Success",
			orgID:       validOrgID.String(),
			requestBody: validPayload,
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().CreateTeam(mock.Anything, validOrgID, validPayload).Return(&validTeam, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"team_id":          validTeam.TeamID.String(),
				"organization_id":  validTeam.OrganizationID.String(),
				"name":             validTeam.Name,
				"description":      validTeam.Description,
				"created_at":       validTeam.CreatedAt.Format(time.RFC3339),
				"updated_at":       validTeam.UpdatedAt.Format(time.RFC3339),
				"deleted_at":       nil,
				"metadata":         nil,
				"team_memberships": nil,
				"created_by":       validTeam.CreatedBy.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/organizations/"+tt.orgID+"/teams", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetTeam(t *testing.T) {
	validOrgID := uuid.New()
	validTeamID := uuid.New()
	validTeam := models.Team{
		TeamID: validTeamID,
		Name:   "Test Team",
	}

	tests := []struct {
		name           string
		orgID          string
		teamID         string
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Invalid org ID",
			orgID:          "invalid-uuid",
			teamID:         validTeamID.String(),
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "Invalid team ID",
			orgID:          validOrgID.String(),
			teamID:         "invalid-uuid",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid team id",
			},
		},
		{
			name:   "Service error",
			orgID:  validOrgID.String(),
			teamID: validTeamID.String(),
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().GetTeamById(mock.Anything, validOrgID, validTeamID).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service error",
			},
		},
		{
			name:   "Success",
			orgID:  validOrgID.String(),
			teamID: validTeamID.String(),
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().GetTeamById(mock.Anything, validOrgID, validTeamID).Return(&validTeam, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"team_id":          validTeam.TeamID.String(),
				"organization_id":  validTeam.OrganizationID.String(),
				"name":             validTeam.Name,
				"description":      validTeam.Description,
				"created_at":       validTeam.CreatedAt.Format(time.RFC3339),
				"updated_at":       validTeam.UpdatedAt.Format(time.RFC3339),
				"deleted_at":       nil,
				"metadata":         nil,
				"team_memberships": nil,
				"created_by":       validTeam.CreatedBy.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/organizations/"+tt.orgID+"/teams/"+tt.teamID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestAddUserToTeam(t *testing.T) {
	validOrgID := uuid.New()
	validTeamID := uuid.New()
	validUserID := uuid.New()
	validTeamMembership := models.TeamMembership{
		TeamID: validTeamID,
		UserID: validUserID,
	}

	tests := []struct {
		name           string
		orgID          string
		teamID         string
		requestBody    interface{}
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Invalid org ID",
			orgID:          "invalid-uuid",
			teamID:         validTeamID.String(),
			requestBody:    teams.AddUserToTeamPayload{UserID: validUserID},
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "Invalid team ID",
			orgID:          validOrgID.String(),
			teamID:         "invalid-uuid",
			requestBody:    teams.AddUserToTeamPayload{UserID: validUserID},
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid team id",
			},
		},
		{
			name:           "Invalid request body",
			orgID:          validOrgID.String(),
			teamID:         validTeamID.String(),
			requestBody:    "invalid json",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request body",
			},
		},
		{
			name:        "Service error",
			orgID:       validOrgID.String(),
			teamID:      validTeamID.String(),
			requestBody: teams.AddUserToTeamPayload{UserID: validUserID},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().AddUserToTeam(mock.Anything, validOrgID, validTeamID, validUserID).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service error",
			},
		},
		{
			name:        "Success",
			orgID:       validOrgID.String(),
			teamID:      validTeamID.String(),
			requestBody: teams.AddUserToTeamPayload{UserID: validUserID},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().AddUserToTeam(mock.Anything, validOrgID, validTeamID, validUserID).Return(&validTeamMembership, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"team_membership_id": validTeamMembership.TeamMembershipID.String(),
				"team_id":            validTeamMembership.TeamID.String(),
				"user_id":            validTeamMembership.UserID.String(),
				"created_at":         validTeamMembership.CreatedAt.Format(time.RFC3339),
				"updated_at":         validTeamMembership.UpdatedAt.Format(time.RFC3339),
				"deleted_at":         nil,
				"team":               nil,
				"user":               nil,
				"created_by":         validTeamMembership.CreatedBy.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/organizations/"+tt.orgID+"/teams/"+tt.teamID+"/add", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestRemoveUserFromTeam(t *testing.T) {
	validOrgID := uuid.New()
	validTeamID := uuid.New()
	validTeamMembershipID := uuid.New()

	tests := []struct {
		name           string
		orgID          string
		teamID         string
		requestBody    interface{}
		setupMock      func(*mockOrganization.MockOrganizationService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Invalid org ID",
			orgID:          "invalid-uuid",
			teamID:         validTeamID.String(),
			requestBody:    teams.RemoveUserFromTeamPayload{TeamMembershipID: validTeamMembershipID},
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid organization id",
			},
		},
		{
			name:           "Invalid team ID",
			orgID:          validOrgID.String(),
			teamID:         "invalid-uuid",
			requestBody:    teams.RemoveUserFromTeamPayload{TeamMembershipID: validTeamMembershipID},
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid team id",
			},
		},
		{
			name:           "Invalid request body",
			orgID:          validOrgID.String(),
			teamID:         validTeamID.String(),
			requestBody:    "invalid json",
			setupMock:      func(mockService *mockOrganization.MockOrganizationService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request body",
			},
		},
		{
			name:        "Service error",
			orgID:       validOrgID.String(),
			teamID:      validTeamID.String(),
			requestBody: teams.RemoveUserFromTeamPayload{TeamMembershipID: validTeamMembershipID},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().RemoveUserFromTeam(mock.Anything, validOrgID, validTeamID, validTeamMembershipID).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service error",
			},
		},
		{
			name:        "Success",
			orgID:       validOrgID.String(),
			teamID:      validTeamID.String(),
			requestBody: teams.RemoveUserFromTeamPayload{TeamMembershipID: validTeamMembershipID},
			setupMock: func(mockService *mockOrganization.MockOrganizationService) {
				mockTeamService := mock_teams.NewMockTeamService(t)
				mockService.EXPECT().TeamService().Return(mockTeamService)
				mockTeamService.EXPECT().RemoveUserFromTeam(mock.Anything, validOrgID, validTeamID, validTeamMembershipID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "user removed from team successfully",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mockOrganization.NewMockOrganizationService(t)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			routerGroup := router.Group("/")
			registerRoutes(routerGroup, mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/organizations/"+tt.orgID+"/teams/"+tt.teamID+"/remove", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}
