package auditlogs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/auditlogs/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mockstore "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupTestService(t *testing.T) (*auditLogService, *mockstore.MockStore) {
	t.Helper()
	mockStore := mockstore.NewMockStore(t)
	service := &auditLogService{
		store: mockStore,
	}
	return service, mockStore
}

func TestNewAuditLogService(t *testing.T) {
	mockStore := new(mockstore.MockStore)

	// Create a mock server config
	serverConfig := &serverconfig.ServerConfig{
		Store: mockStore,
	}

	service := NewAuditLogService(serverConfig)

	assert.NotNil(t, service)
	assert.IsType(t, &auditLogService{}, service)
}

func TestGetAuditLogsByOrganizationId(t *testing.T) {
	tests := []struct {
		name          string
		orgId         uuid.UUID
		kind          dbmodels.AuditLogKind
		mockAuditLogs []dbmodels.AuditLog
		mockError     error
		wantError     bool
	}{
		{
			name:  "Success case",
			orgId: uuid.New(),
			kind:  dbmodels.AuditLogKindInfo,
			mockAuditLogs: []dbmodels.AuditLog{
				{
					ID:             uuid.New(),
					Kind:           dbmodels.AuditLogKindInfo,
					OrganizationID: uuid.New(),
					ResourceType:   dbmodels.ResourceTypeOrganization,
					ResourceID:     uuid.New(),
					EventName:      "test_event",
					Payload:        []byte(`{"test": "value"}`),
				},
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:          "Store returns error",
			orgId:         uuid.New(),
			kind:          dbmodels.AuditLogKindInfo,
			mockAuditLogs: nil,
			mockError:     errors.New("database error"),
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore := setupTestService(t)

			ctx := context.Background()
			logger, _ := zap.NewDevelopment()
			ctx = apicontext.AddLoggerToContext(ctx, logger)

			mockStore.On("GetAuditLogsByOrganizationId", mock.Anything, tt.orgId, tt.kind).Return(tt.mockAuditLogs, tt.mockError)

			logs, err := service.GetAuditLogsByOrganizationId(ctx, tt.orgId, tt.kind)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, logs)
			} else {
				assert.NoError(t, err)
				assert.Len(t, logs, len(tt.mockAuditLogs))
				if len(logs) > 0 {
					assert.Equal(t, tt.mockAuditLogs[0].OrganizationID, logs[0].OrganizationID)
				}
			}
			mockStore.AssertExpectations(t)
		})
	}
}

func TestEmitAuditLog(t *testing.T) {
	tests := []struct {
		name           string
		params         models.AuditLogEmitParams
		contextEmail   string
		contextIP      string
		contextAgent   string
		mockSetup      func(*mockstore.MockStore)
		wantError      bool
		errorMsg       string
		expectedEmail  string
		expectedIP     string
		expectedAgent  string
	}{
		{
			name: "Success case with params values",
			params: models.AuditLogEmitParams{
				Kind:           dbmodels.AuditLogKindInfo,
				ResourceType:   dbmodels.ResourceTypeOrganization,
				ResourceID:     uuid.New(),
				EventName:      "test_event",
				Payload:        map[string]interface{}{"test": "value"},
				OrganizationID: uuid.New(),
				UserEmail:      "params@example.com",
				IPAddress:      "10.0.0.1",
				UserAgent:      "ParamsUserAgent",
			},
			contextEmail:  "context@example.com",
			contextIP:     "192.168.1.1",
			contextAgent:  "ContextUserAgent",
			expectedEmail: "params@example.com",
			expectedIP:    "10.0.0.1",
			expectedAgent: "ParamsUserAgent",
			mockSetup: func(mockStore *mockstore.MockStore) {
				mockStore.On("CreateAuditLog", mock.Anything, mock.MatchedBy(func(log dbmodels.AuditLog) bool {
					var payload map[string]interface{}
					err := json.Unmarshal(log.Payload, &payload)
					return err == nil && 
						log.Kind == dbmodels.AuditLogKindInfo &&
						log.ResourceType == dbmodels.ResourceTypeOrganization &&
						log.UserEmail == "params@example.com" &&
						log.IPAddress == "10.0.0.1" &&
						log.UserAgent == "ParamsUserAgent"
				})).Return(&dbmodels.AuditLog{}, nil)
			},
			wantError: false,
		},
		{
			name: "Success case with context values",
			params: models.AuditLogEmitParams{
				Kind:           dbmodels.AuditLogKindInfo,
				ResourceType:   dbmodels.ResourceTypeOrganization,
				ResourceID:     uuid.New(),
				EventName:      "test_event",
				Payload:        map[string]interface{}{"test": "value"},
				OrganizationID: uuid.New(),
				// No email, IP, or user agent in params
			},
			contextEmail:  "context@example.com",
			contextIP:     "192.168.1.1",
			contextAgent:  "ContextUserAgent",
			expectedEmail: "context@example.com",
			expectedIP:    "192.168.1.1",
			expectedAgent: "ContextUserAgent",
			mockSetup: func(mockStore *mockstore.MockStore) {
				mockStore.On("CreateAuditLog", mock.Anything, mock.MatchedBy(func(log dbmodels.AuditLog) bool {
					var payload map[string]interface{}
					err := json.Unmarshal(log.Payload, &payload)
					return err == nil && 
						log.Kind == dbmodels.AuditLogKindInfo &&
						log.ResourceType == dbmodels.ResourceTypeOrganization &&
						log.UserEmail == "context@example.com" &&
						log.IPAddress == "192.168.1.1" &&
						log.UserAgent == "ContextUserAgent"
				})).Return(&dbmodels.AuditLog{}, nil)
			},
			wantError: false,
		},
		{
			name: "JSON marshal error",
			params: models.AuditLogEmitParams{
				Kind:         dbmodels.AuditLogKindInfo,
				ResourceType: dbmodels.ResourceTypeOrganization,
				ResourceID:   uuid.New(),
				EventName:    "test_event",
				Payload: map[string]interface{}{
					"test": make(chan int),
				},
				OrganizationID: uuid.New(),
			},
			contextEmail: "context@example.com",
			contextIP:    "192.168.1.1",
			contextAgent: "ContextUserAgent",
			mockSetup:    func(mockStore *mockstore.MockStore) {},
			wantError:    true,
			errorMsg:     "json",
		},
		{
			name: "Store error",
			params: models.AuditLogEmitParams{
				Kind:           dbmodels.AuditLogKindInfo,
				ResourceType:   dbmodels.ResourceTypeOrganization,
				ResourceID:     uuid.New(),
				EventName:      "test_event",
				Payload:        map[string]interface{}{"test": "value"},
				OrganizationID: uuid.New(),
			},
			contextEmail: "context@example.com",
			contextIP:    "192.168.1.1",
			contextAgent: "ContextUserAgent",
			mockSetup: func(mockStore *mockstore.MockStore) {
				mockStore.On("CreateAuditLog", mock.Anything, mock.Anything).Return(nil, errors.New("store error"))
			},
			wantError: true,
			errorMsg:  "store error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore := setupTestService(t)

			ctx := context.Background()
			logger, _ := zap.NewDevelopment()
			ctx = apicontext.AddLoggerToContext(ctx, logger)
			
			// Add audit info to context
			ctx = apicontext.AddAuditInfoToContext(ctx, tt.contextEmail, tt.contextIP, tt.contextAgent)

			tt.mockSetup(mockStore)

			err := service.EmitAuditLog(ctx, tt.params)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
			mockStore.AssertExpectations(t)
		})
	}
}

func TestWithResource(t *testing.T) {
	service, _ := setupTestService(t)

	resourceService := service.WithResource("test_resource")

	assert.NotNil(t, resourceService)
	assert.IsType(t, &auditLogServiceWithResource{}, resourceService)
}

func TestGetResourceAuditLogMiddleware(t *testing.T) {
	service, _ := setupTestService(t)

	middleware := service.GetResourceAuditLogMiddleware("test_resource", "id")

	assert.NotNil(t, middleware)
	assert.IsType(t, gin.HandlerFunc(nil), middleware)

	// Test the middleware execution
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = []gin.Param{{Key: "id", Value: "test_id"}}

	// Mock request
	req, _ := http.NewRequest("POST", "/test", nil)
	c.Request = req

	// Create a test handler that will be called after middleware
	testHandler := func(c *gin.Context) {
		// In Phase 2, we'll implement the actual middleware functionality
		// For now, we're just testing that the middleware doesn't break the request chain
	}

	// Execute middleware and handler
	middleware(c)
	testHandler(c)

	// Since the middleware is a placeholder for Phase 2, we're just verifying it exists
	// and doesn't break the request chain
	assert.True(t, true)
}
func TestAuditLogServiceWithResourceEmitAuditLog(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   dbmodels.ResourceType
		userRole       string
		userEmail      string
		contextEmail   string
		contextIP      string
		contextAgent   string
		eventName      string
		payload        map[string]interface{}
		getUserError   error
		createLogError error
		wantError      bool
		useContextInfo bool
	}{
		{
			name:         "successful audit log creation with user info",
			resourceType: dbmodels.ResourceTypeOrganization,
			userRole:     "admin",
			userEmail:    "test@example.com",
			eventName:    "test_event",
			payload: map[string]interface{}{
				"test": "value",
			},
			useContextInfo: false,
		},
		{
			name:         "successful audit log creation with context info",
			resourceType: dbmodels.ResourceTypeOrganization,
			userRole:     "admin",
			userEmail:    "test@example.com",
			contextEmail: "context@example.com",
			contextIP:    "192.168.1.1",
			contextAgent: "ContextUserAgent",
			eventName:    "test_event",
			payload: map[string]interface{}{
				"test": "value",
			},
			useContextInfo: true,
		},
		{
			name:         "get user error with no context info",
			resourceType: dbmodels.ResourceTypeOrganization,
			userRole:     "admin",
			userEmail:    "test@example.com",
			eventName:    "test_event",
			payload: map[string]interface{}{
				"test": "value",
			},
			getUserError:   fmt.Errorf("user not found"),
			useContextInfo: false,
		},
		{
			name:         "get user error with context info",
			resourceType: dbmodels.ResourceTypeOrganization,
			userRole:     "admin",
			userEmail:    "test@example.com",
			contextEmail: "context@example.com",
			contextIP:    "192.168.1.1",
			contextAgent: "ContextUserAgent",
			eventName:    "test_event",
			payload: map[string]interface{}{
				"test": "value",
			},
			getUserError:   fmt.Errorf("user not found"),
			useContextInfo: true,
		},
		{
			name:         "create log error",
			resourceType: dbmodels.ResourceTypeOrganization,
			userRole:     "admin",
			userEmail:    "test@example.com",
			contextEmail: "context@example.com",
			contextIP:    "192.168.1.1",
			contextAgent: "ContextUserAgent",
			eventName:    "test_event",
			payload: map[string]interface{}{
				"test": "value",
			},
			createLogError: fmt.Errorf("failed to create log"),
			wantError:      true,
			useContextInfo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore := setupTestService(t)

			ctx := context.Background()
			logger, _ := zap.NewDevelopment()
			ctx = apicontext.AddLoggerToContext(ctx, logger)

			// Mock auth context
			userId := uuid.New()
			orgId := uuid.New()
			ctx = apicontext.AddAuthToContext(ctx, tt.userRole, userId, []uuid.UUID{orgId})
			
			// Add audit info to context if test case uses it
			if tt.useContextInfo {
				ctx = apicontext.AddAuditInfoToContext(ctx, tt.contextEmail, tt.contextIP, tt.contextAgent)
			}

			// Mock user if needed
			if !tt.useContextInfo || tt.contextEmail == "" {
				user := &dbmodels.User{
					Email: tt.userEmail,
				}
				mockStore.On("GetUserById", mock.Anything, userId.String()).Return(user, tt.getUserError)
			}

			if tt.getUserError == nil || tt.useContextInfo {
				// Mock CreateAuditLog
				mockStore.On("CreateAuditLog", mock.Anything, mock.MatchedBy(func(log dbmodels.AuditLog) bool {
					expectedEmail := tt.userEmail
					if tt.useContextInfo && tt.contextEmail != "" {
						expectedEmail = tt.contextEmail
					}
					
					match := log.ResourceType == tt.resourceType &&
						log.ResourceID == orgId &&
						log.EventName == tt.eventName &&
						log.Kind == dbmodels.AuditLogKindInfo
					
					if tt.useContextInfo {
						match = match && 
							log.UserEmail == expectedEmail &&
							log.IPAddress == tt.contextIP &&
							log.UserAgent == tt.contextAgent
					} else {
						match = match && log.UserEmail == expectedEmail
					}
					
					return match
				})).Return(&dbmodels.AuditLog{}, tt.createLogError)
			}

			resourceService := service.WithResource(tt.resourceType)

			err := resourceService.EmitAuditLog(ctx, orgId, dbmodels.AuditLogKindInfo, tt.eventName, tt.payload)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuditLogServiceWithResourceEmitAuditLogNoAuth(t *testing.T) {
	service, mockStore := setupTestService(t)

	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	ctx = apicontext.AddLoggerToContext(ctx, logger)

	orgId := uuid.New()

	// No auth context added

	resourceService := service.WithResource("test_resource")

	err := resourceService.EmitAuditLog(ctx, orgId, dbmodels.AuditLogKindInfo, "test_event", map[string]interface{}{"test": "value"})

	// Should not return error even if auth is missing
	assert.NoError(t, err)
	// No calls to CreateAuditLog should be made
	mockStore.AssertNotCalled(t, "CreateAuditLog")
}
