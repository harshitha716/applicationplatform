package connections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mock_temporal "github.com/Zampfi/workflow-sdk-go/mocks/workflowmanagers/temporal"

	client "go.temporal.io/sdk/client"

	temporalsdkmodels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
)

type CustomResponseRecorder struct {
	*httptest.ResponseRecorder
}

func (r *CustomResponseRecorder) CloseNotify() <-chan bool {
	return make(<-chan bool)
}

type testCase struct {
	name          string
	method        string
	path          string
	skip          bool
	inputPayload  map[string]interface{}
	outputPayload string
	statusCode    int
	initServerCfg func() *serverconfig.ServerConfig
}

func TestCreateConnection(t *testing.T) {
	connectorId := uuid.New()
	connectionId := uuid.New()
	userId := uuid.New()
	policyId := uuid.New()

	testCases := []testCase{
		{
			name:   "Create connection success",
			method: "POST",
			path:   "/connections/",
			skip:   false,
			inputPayload: map[string]interface{}{
				"connector_id":   connectorId.String(),
				"display_name":   "Test Connection",
				"connector_name": "gcs",
				"connection_config": map[string]interface{}{"key": "value", "bucket_name": "bucket", "schedules": []map[string]string{
					{
						"cron_expression": "0 * * * *",
						"glob_pattern":    "test/*.csv",
						"file_format":     "csv",
					},
				}},
			},
			outputPayload: `{"connection_id":"` + connectionId.String() + `"}`,
			statusCode:    http.StatusCreated,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				// Create mock store
				mockStore := mock_store.NewMockStore(t)

				mockTemporalService := mock_temporal.NewMockTemporalService(t)
				expectedResponse := temporalsdkmodels.ScheduledWorkflowResponse{
					ScheduleID: "test-schedule",
					// Fill other required fields
				}

				mockTemporalService.On("ExecuteScheduledWorkflow", mock.Anything, mock.Anything).Return(expectedResponse, nil)

				mockStore.On("WithTx", mock.Anything, mock.AnythingOfType("func(store.Store) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.Store) error)
						fn(mockStore)
					}).
					Return(nil)

				// Setup CreateConnection mock
				mockStore.On("CreateConnection", mock.Anything, mock.MatchedBy(func(params *models.CreateConnectionParams) bool {
					return params.ConnectorID == connectorId && params.Name == "Test Connection"
				})).Return(connectionId, nil)

				// Setup CreateConnectionPolicy mock
				mockStore.On("CreateConnectionPolicy",
					mock.Anything,
					connectionId,
					models.AudienceTypeUser,
					userId,
					models.PrivilegeConnectionAdmin,
				).Return(&models.ResourceAudiencePolicy{
					ID:                   policyId,
					ResourceType:         models.ResourceTypeConnection,
					ResourceID:           connectionId,
					ResourceAudienceType: models.AudienceTypeUser,
					ResourceAudienceID:   userId,
					Privilege:            models.PrivilegeConnectionAdmin,
				}, nil)

				mockStore.On("GetOrganizationsAll", mock.Anything, mock.Anything).Return([]models.Organization{
					{
						ID:   uuid.New(),
						Name: "Test Organization",
					},
				}, nil)

				mockStore.On("CreateSchedules", mock.Anything, mock.AnythingOfType("[]models.CreateScheduleParams"), mock.Anything).Return(nil)
				mockTemporalService.On("ExecuteScheduledWorkflow", mock.Anything, mock.Anything).Return(nil)

				svCfg.Store = mockStore
				svCfg.TemporalSdk = mockTemporalService

				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			// Initialize server config
			svCfg := tc.initServerCfg()

			// Create test router
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				// Add auth context for successful test case
				if tc.statusCode == http.StatusCreated {
					// Add auth context variables directly
					ctxVars := map[string]interface{}{
						"user_id":            userId,
						"user_organizations": []uuid.UUID{},
						"user_role":          "user",
					}
					c.Set("context_variables", ctxVars)
				}
			})

			// Register routes
			apiGroup := router.Group("/")
			err := RegisterConnectionRoutes(apiGroup, svCfg)
			assert.NoError(t, err)

			// Create request
			jsonData, err := json.Marshal(tc.inputPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest(tc.method, tc.path, bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := &CustomResponseRecorder{httptest.NewRecorder()}

			// Serve request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tc.statusCode, w.Code)
			assert.JSONEq(t, tc.outputPayload, w.Body.String())

			// Verify mock expectations for successful case
			if tc.statusCode == http.StatusCreated {
				mockStore := svCfg.Store.(*mock_store.MockStore)
				mockStore.AssertExpectations(t)
			}
		})
	}
}

func TestCreateConnectionSnowflake(t *testing.T) {

	connectorId := uuid.New()
	connectionId := uuid.New()
	userId := uuid.New()
	policyId := uuid.New()

	testCases := []testCase{
		{
			name:   "Create connection success",
			method: "POST",
			path:   "/connections/",
			skip:   false,
			inputPayload: map[string]interface{}{
				"connector_id":   connectorId.String(),
				"display_name":   "Test Snowflake Connection",
				"connector_name": "snowflake",
				"connection_config": map[string]interface{}{"s3_bucket": "zamp-prd-file-imports",
					"snowflake_table":              "sample_events",
					"snowflake_schema":             "STREAM_TEST",
					"snowflake_database":           "dev",
					"s3_destination_path":          "test-data",
					"s3_storage_integration":       "s3_snowflake_integration",
					"snowflake_filter_column_name": "event_time",
					"schedules": []map[string]string{
						{
							"cron_expression": "0 * * * *",
							"glob_pattern":    "test/*.csv",
							"file_format":     "csv",
						},
					}},
			},
			outputPayload: `{"connection_id":"` + connectionId.String() + `"}`,
			statusCode:    http.StatusCreated,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				// Create mock store
				mockStore := mock_store.NewMockStore(t)

				mockTemporalService := mock_temporal.NewMockTemporalService(t)
				expectedResponse := temporalsdkmodels.ScheduledWorkflowResponse{
					ScheduleID: "test-schedule",
					// Fill other required fields
				}

				mockTemporalService.On("ExecuteScheduledWorkflow", mock.Anything, mock.Anything).Return(expectedResponse, nil)

				mockStore.On("WithTx", mock.Anything, mock.AnythingOfType("func(store.Store) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.Store) error)
						fn(mockStore)
					}).
					Return(nil)

				// Setup CreateConnection mock
				mockStore.On("CreateConnection", mock.Anything, mock.MatchedBy(func(params *models.CreateConnectionParams) bool {
					return params.ConnectorID == connectorId && params.Name == "Test Snowflake Connection"
				})).Return(connectionId, nil)

				// Setup CreateConnectionPolicy mock
				mockStore.On("CreateConnectionPolicy",
					mock.Anything,
					connectionId,
					models.AudienceTypeUser,
					userId,
					models.PrivilegeConnectionAdmin,
				).Return(&models.ResourceAudiencePolicy{
					ID:                   policyId,
					ResourceType:         models.ResourceTypeConnection,
					ResourceID:           connectionId,
					ResourceAudienceType: models.AudienceTypeUser,
					ResourceAudienceID:   userId,
					Privilege:            models.PrivilegeConnectionAdmin,
				}, nil)

				mockStore.On("GetOrganizationsAll", mock.Anything, mock.Anything).Return([]models.Organization{
					{
						ID:   uuid.New(),
						Name: "Test Organization",
					},
				}, nil)

				mockStore.On("CreateSchedules", mock.Anything, mock.AnythingOfType("[]models.CreateScheduleParams"), mock.Anything).Return(nil)
				mockTemporalService.On("ExecuteScheduledWorkflow", mock.Anything, mock.Anything).Return(nil)

				svCfg.Store = mockStore
				svCfg.TemporalSdk = mockTemporalService

				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			// Initialize server config
			svCfg := tc.initServerCfg()

			// Create test router
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				// Add auth context for successful test case
				if tc.statusCode == http.StatusCreated {
					// Add auth context variables directly
					ctxVars := map[string]interface{}{
						"user_id":            userId,
						"user_organizations": []uuid.UUID{},
						"user_role":          "user",
					}
					c.Set("context_variables", ctxVars)
				}
			})

			// Register routes
			apiGroup := router.Group("/")
			err := RegisterConnectionRoutes(apiGroup, svCfg)
			assert.NoError(t, err)

			// Create request
			jsonData, err := json.Marshal(tc.inputPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest(tc.method, tc.path, bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := &CustomResponseRecorder{httptest.NewRecorder()}

			// Serve request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tc.statusCode, w.Code)
			assert.JSONEq(t, tc.outputPayload, w.Body.String())

			// Verify mock expectations for successful case
			if tc.statusCode == http.StatusCreated {
				mockStore := svCfg.Store.(*mock_store.MockStore)
				mockStore.AssertExpectations(t)
			}
		})
	}

}

func TestGetConnections(t *testing.T) {
	userId := uuid.New()
	connectorId := uuid.New()
	connectionId := uuid.New()
	scheduleId := uuid.New()
	lastSyncedAt := time.Time{}.UTC()
	createdAt := time.Now().UTC()

	testCases := []testCase{
		{
			name:   "Get connections success",
			method: "GET",
			path:   "/connections",
			skip:   false,
			outputPayload: fmt.Sprintf(`[{
				"id": "%s",
				"name": "Test Connection",
				"icon_url": "test-logo.png",
				"last_synced_at": "%s",
				"created_at": "%s"
			}]`, connectionId.String(), lastSyncedAt.Format(time.RFC3339), createdAt.Format(time.RFC3339)),
			statusCode: http.StatusOK,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockTemporalService := mock_temporal.NewMockTemporalService(t)

				// Setup GetConnections mock
				mockStore.On("GetConnections", mock.Anything, mock.Anything).Return([]models.Connection{
					{
						ID:          connectionId,
						ConnectorID: connectorId,
						Name:        "Test Connection",
						Status:      "active",
						CreatedAt:   createdAt,
						Connector: models.Connector{
							LogoURL: "test-logo.png",
						},
						Schedules: []models.Schedule{
							{
								ID:                 scheduleId,
								ConnectionID:       connectionId,
								TemporalWorkflowID: "test-workflow",
								ScheduleGroup:      "test-group",
								CronSchedule:       "0 * * * *",
								Config:             json.RawMessage(`{"key": "value"}`),
								ConnectorID:        connectorId,
								OrganizationID:     uuid.New(),
								CreatedAt:          createdAt,
								UpdatedAt:          createdAt,
							},
						},
					},
				}, nil)

				// Setup GetLastSyncedAt mock
				mockTemporalService.On("QuerySchedule", mock.Anything, temporalsdkmodels.QueryScheduleParams{
					ScheduleID: scheduleId.String(),
				}).Return(temporalsdkmodels.QueryScheduleResponse{
					ScheduleDescription: client.ScheduleDescription{
						Info: client.ScheduleInfo{
							RecentActions: []client.ScheduleActionResult{
								{
									ScheduleTime: lastSyncedAt,
									ActualTime:   lastSyncedAt,
									StartWorkflowResult: &client.ScheduleWorkflowExecution{
										WorkflowID: "test-workflow",
									},
								},
							},
						},
					},
				}, nil)

				mockTemporalService.On("GetWorkflowDetails", mock.Anything, temporalsdkmodels.GetWorkflowDetailsParams{
					WorkflowID: "test-workflow",
				}).Return(temporalsdkmodels.WorkflowDetailsResponse{
					Details: temporalsdkmodels.WorkflowExecutionDetails{
						Status: "RUNNING",
					},
				}, nil)

				mockStore.On("GetSchedulesByConnectionID", mock.Anything, mock.Anything).Return([]models.Schedule{
					{
						ID:                 scheduleId,
						ConnectionID:       connectionId,
						Status:             "active",
						CreatedAt:          createdAt,
						UpdatedAt:          createdAt,
						CronSchedule:       "0 * * * *",
						Config:             json.RawMessage(`{"key": "value"}`),
						ConnectorID:        connectorId,
						OrganizationID:     uuid.New(),
						ScheduleGroup:      "test-group",
						TemporalWorkflowID: "test-workflow",
					},
				}, nil)

				mockTemporalService.On("QuerySchedule", mock.Anything, temporalsdkmodels.QueryScheduleParams{
					ScheduleID: scheduleId.String(),
				}).Return(temporalsdkmodels.QueryScheduleResponse{
					ScheduleDescription: client.ScheduleDescription{
						Info: client.ScheduleInfo{
							NextActionTimes: []time.Time{lastSyncedAt},
							RecentActions: []client.ScheduleActionResult{
								{
									ScheduleTime: lastSyncedAt,
									ActualTime:   lastSyncedAt,
								},
							},
							CreatedAt:    lastSyncedAt,
							LastUpdateAt: lastSyncedAt,
						},
					},
				}, nil)

				mockTemporalService.On("GetWorkflowDetails", mock.Anything, temporalsdkmodels.GetWorkflowDetailsParams{
					WorkflowID: "test-workflow",
				}).Return(temporalsdkmodels.WorkflowDetailsResponse{
					Details: temporalsdkmodels.WorkflowExecutionDetails{
						Status: "COMPLETED",
					},
				}, nil)

				svCfg.Store = mockStore
				svCfg.TemporalSdk = mockTemporalService

				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			// Initialize server config
			svCfg := tc.initServerCfg()

			// Create test router
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				if tc.statusCode == http.StatusOK {
					ctxVars := map[string]interface{}{
						"user_id":            userId,
						"user_organizations": []uuid.UUID{},
						"user_role":          "user",
					}
					c.Set("context_variables", ctxVars)
				}
			})

			// Register routes
			apiGroup := router.Group("/")
			err := RegisterConnectionRoutes(apiGroup, svCfg)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest(tc.method, tc.path, nil)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := &CustomResponseRecorder{httptest.NewRecorder()}

			// Serve request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tc.statusCode, w.Code)
			assert.JSONEq(t, tc.outputPayload, w.Body.String())

			// Verify mock expectations
			if tc.statusCode == http.StatusOK {
				mockStore := svCfg.Store.(*mock_store.MockStore)
				mockStore.AssertExpectations(t)
			}
		})
	}
}

func TestGetSchedules(t *testing.T) {
	userId := uuid.New()
	connectorId := uuid.New()
	connectionId := uuid.New()
	scheduleId := uuid.New()
	now := time.Now()

	testCases := []testCase{
		{
			name:   "Get schedules success",
			method: "GET",
			path:   fmt.Sprintf("/connections/%s/schedules", connectionId.String()),
			skip:   false,
			outputPayload: fmt.Sprintf(`[{
				"id": "%s",
				"name": "Test Schedule",
				"status": "active",
				"next_run_at": "%s",
				"last_run_at": "%s",
				"created_at": "%s",
				"updated_at": "%s",
				"logo_url": "test-logo.png"
			}]`, scheduleId.String(), now.Format(time.RFC3339), now.Format(time.RFC3339), now.Format(time.RFC3339), now.Format(time.RFC3339)),
			statusCode: http.StatusOK,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockTemporalService := mock_temporal.NewMockTemporalService(t)

				// Setup GetConnectionByID mock
				mockStore.On("GetConnectionByID", mock.Anything, connectionId).Return(&models.Connection{
					ID:          connectionId,
					ConnectorID: connectorId,
					Name:        "Test Connection",
					Status:      "active",
					Connector: models.Connector{
						LogoURL: "test-logo.png",
					},
				}, nil)

				// Setup GetSchedules mock
				mockStore.On("GetSchedulesByConnectionID", mock.Anything, connectionId).Return([]models.Schedule{
					{
						ID:           scheduleId,
						Name:         "Test Schedule",
						Status:       "active",
						ConnectionID: connectionId,
					},
				}, nil)

				// Setup GetScheduleDetailsFromTemporal mock
				mockTemporalService.On("QuerySchedule", mock.Anything, temporalsdkmodels.QueryScheduleParams{
					ScheduleID: scheduleId.String(),
				}).Return(temporalsdkmodels.QueryScheduleResponse{
					ScheduleDescription: client.ScheduleDescription{
						Info: client.ScheduleInfo{
							NextActionTimes: []time.Time{now},
							RecentActions: []client.ScheduleActionResult{
								{
									ScheduleTime: now,
									ActualTime:   now,
								},
							},
							CreatedAt:    now,
							LastUpdateAt: now,
						},
					},
				}, nil)

				svCfg.Store = mockStore
				svCfg.TemporalSdk = mockTemporalService

				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			// Initialize server config
			svCfg := tc.initServerCfg()

			// Create test router
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				if tc.statusCode == http.StatusOK {
					ctxVars := map[string]interface{}{
						"user_id":            userId,
						"user_organizations": []uuid.UUID{},
						"user_role":          "user",
					}
					c.Set("context_variables", ctxVars)
				}
			})

			// Register routes
			apiGroup := router.Group("/")
			err := RegisterConnectionRoutes(apiGroup, svCfg)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest(tc.method, tc.path, nil)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := &CustomResponseRecorder{httptest.NewRecorder()}

			// Serve request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tc.statusCode, w.Code)
			assert.JSONEq(t, tc.outputPayload, w.Body.String())

			// Verify mock expectations
			if tc.statusCode == http.StatusOK {
				mockStore := svCfg.Store.(*mock_store.MockStore)
				mockStore.AssertExpectations(t)
			}
		})
	}
}
