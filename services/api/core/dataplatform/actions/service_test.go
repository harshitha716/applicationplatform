package actions

import (
	"context"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	serviceconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	dataserviceconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataservicemodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	helper "github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	mockdataservice "github.com/Zampfi/application-platform/services/api/mocks/core/dataplatform/data"
	mockdatabricksservice "github.com/Zampfi/application-platform/services/api/mocks/pkg/dataplatform/providers/databricks"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ActionServiceTestSuite struct {
	suite.Suite
	mockDataService       *mockdataservice.MockDataService
	service               *actionService
	mockDatabricksService *mockdatabricksservice.MockDatabricksService
}

func getDataPlatformMockConfig() *serverconfig.DataPlatformConfig {
	return &serverconfig.DataPlatformConfig{
		DatabricksConfig: serverconfig.DatabricksSetupConfig{
			MerchantDataProviderIdMapping: map[string]string{"merchant1": "workspace1"},
			DefaultDataProviderId:         "defaultWorkspace",
			DataProviderConfigs: map[string]dataplatformmodels.DatabricksConfig{
				"workspace1": {
					WarehouseId: "warehouse1",
				},
			},
			ZampDatabricksCatalog:        "zamp",
			ZampDatabricksPlatformSchema: "platform",
		},
		PinotConfig: serverconfig.PinotSetupConfig{
			MerchantDataProviderIdMapping: map[string]string{"merchant2": "workspace2"},
			DefaultDataProviderId:         "defaultPinotWorkspace",
		},
		ActionsConfig: serverconfig.ActionsConfig{
			CreateMVJobTemplateConfig: serverconfig.CreateMVJobTemplateConfig{
				SideEffectNotebookPath: "/path/to/sideeffect/notebook",
				CreateMVNotebookPath:   "/path/to/create/mv/notebook",
			},
			RegisterDatasetJobTemplateConfig: serverconfig.RegisterDatasetJobTemplateConfig{
				RegisterDatasetNotebookPath: "/path/to/register/dataset/notebook",
			},
			RegisterJobJobTemplateConfig: serverconfig.RegisterJobJobTemplateConfig{
				RegisterJobNotebookPath: "/path/to/register/job/notebook",
			},
			WebhookConfig: serverconfig.WebhookConfig{
				WebhookId: "webhook_id",
				UserName:  "user_name",
				Password:  "password",
			},
			UpsertTemplateJobTemplateConfig: serverconfig.UpsertTemplateJobTemplateConfig{
				UpsertTemplateNotebookPath: "/path/to/upsert/template/notebook",
			},
			UpdateDatasetJobTemplateConfig: serverconfig.UpdateDatasetJobTemplateConfig{
				UpdateDatasetNotebookPath: "/path/to/update/dataset/notebook",
			},
			CopyDatasetJobTemplateConfig: serverconfig.CopyDatasetJobTemplateConfig{
				CopyDatasetNotebookPath: "/path/to/copy/dataset/notebook",
			},
			RunAsUserName:              "test_user",
			PinotIngestionNotebookPath: "/path/to/pinot/ingestion/notebook",
			PinotClusterId:             "1111",
			DataPlatformModulesSrc:     "test_src",
		},
	}
}

func TestActionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ActionServiceTestSuite))
}

func (s *ActionServiceTestSuite) SetupTest() {
	s.mockDataService = new(mockdataservice.MockDataService)
	s.mockDatabricksService = new(mockdatabricksservice.MockDatabricksService)
	s.service = &actionService{
		dataService: s.mockDataService,
	}
}

func (s *ActionServiceTestSuite) TestAddJobPayloadToTemplate() {
	tests := []struct {
		name           string
		jobTemplate    *jobs.SubmitRun
		jobPayload     map[string]string
		expectedResult *jobs.SubmitRun
	}{
		{
			name: "Add payload to notebook task",
			jobTemplate: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{
					{
						NotebookTask: &jobs.NotebookTask{},
					},
				},
			},
			jobPayload: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			expectedResult: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{
					{
						NotebookTask: &jobs.NotebookTask{
							BaseParameters: map[string]string{
								"param1":                                "value1",
								"param2":                                "value2",
								serviceconstants.DataPlatformModulesSrc: "test_src",
							},
						},
					},
				},
			},
		},
		{
			name: "Add payload to SQL task",
			jobTemplate: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{
					{
						SqlTask: &jobs.SqlTask{},
					},
				},
			},
			jobPayload: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			expectedResult: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{
					{
						SqlTask: &jobs.SqlTask{
							Parameters: map[string]string{
								"param1":                                "value1",
								"param2":                                "value2",
								serviceconstants.DataPlatformModulesSrc: "test_src",
							},
						},
					},
				},
			},
		},
		{
			name: "No tasks in job template",
			jobTemplate: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{},
			},
			jobPayload: map[string]string{
				"param1": "value1",
			},
			expectedResult: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig()).Once()
			result := s.service.addJobPayloadToTemplate(tt.jobTemplate, tt.jobPayload)
			s.Equal(tt.expectedResult, result)
		})
	}
}

func (s *ActionServiceTestSuite) TestGetCreateMVJobTemplate() {

	tests := []struct {
		name                  string
		providerId            string
		payload               models.CreateActionPayload
		mockWarehouseId       string
		mockError             error
		processParamsForQuery dataservicemodels.QueryMetadata
		expectedJobTemplate   *jobs.SubmitRun
		expectError           bool
	}{
		{
			name:       "Successful job template creation",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeCreateMV,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.CreateMVActionPayload{
					Query:            "SELECT * FROM {{.zamp_table_name_1}} RIGHT JOIN {{.zamp_table_name_2}} ON {{.zamp_table_name_1}}.id = {{.zamp_table_name_2}}.id",
					QueryParams:      map[string]string{"zamp_table_name_1": "1", "zamp_table_name_2": "2"},
					MVDatasetId:      "mv_dataset_id",
					ParentDatasetIds: []string{"1", "2"},
					DedupColumns:     []string{"id"},
					OrderByColumn:    "id",
				},
			},
			mockError: nil,
			processParamsForQuery: dataservicemodels.QueryMetadata{
				TableNames: []string{"1", "2"},
				Params:     map[string]string{"zamp_table_name_1": "dataset1", "zamp_table_name_2": "dataset2"},
			},
			mockWarehouseId: "warehouse1",
			expectedJobTemplate: &jobs.SubmitRun{
				RunName: "create_mv_merchant1",
				WebhookNotifications: &jobs.WebhookNotifications{
					OnStart: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnSuccess: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnFailure: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
				},
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "create_mv_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/create/mv/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.CreateMVParams:         `{"parent_dataset_ids":["1","2"],"query":"SELECT * FROM dataset1 RIGHT JOIN dataset2 ON dataset1.id = dataset2.id","merchant_id":"merchant1","dataset_id":"mv_dataset_id","dedup_columns":["id"],"order_by_column":"id"}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
								serviceconstants.DatasetIdParam:         "mv_dataset_id",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
					{
						TaskKey: "create_mv_sideeffect",
						DependsOn: []jobs.TaskDependency{
							{
								TaskKey: "create_mv_task",
							},
						},
						RunIf: jobs.RunIfAllSuccess,
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/sideeffect/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.CreateMVParams:         `{"parent_dataset_ids":["1","2"],"query":"SELECT * FROM dataset1 RIGHT JOIN dataset2 ON dataset1.id = dataset2.id","merchant_id":"merchant1","dataset_id":"mv_dataset_id","dedup_columns":["id"],"order_by_column":"id"}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
								serviceconstants.DatasetIdParam:         "mv_dataset_id",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
					{
						TaskKey: "pinot_ingestion_task",
						DependsOn: []jobs.TaskDependency{
							{
								TaskKey: "create_mv_sideeffect",
							},
						},
						RunIf: jobs.RunIfAllSuccess,
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/pinot/ingestion/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.CreateMVParams:         `{"parent_dataset_ids":["1","2"],"query":"SELECT * FROM dataset1 RIGHT JOIN dataset2 ON dataset1.id = dataset2.id","merchant_id":"merchant1","dataset_id":"mv_dataset_id","dedup_columns":["id"],"order_by_column":"id"}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
								serviceconstants.DatasetIdParam:         "mv_dataset_id",
							},
						},
						ExistingClusterId:    "1111",
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
				},
				Queue: &jobs.QueueSettings{
					Enabled: true,
				},
				RunAs: &jobs.JobRunAs{
					UserName: "test_user",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())
			s.mockDataService.On("GetDatabricksWarehouseId", ctx, mock.Anything).Return(tt.mockWarehouseId, nil).Once()
			s.mockDataService.On("ProcessParamsForQuery", ctx, mock.Anything, mock.Anything, mock.Anything).Return(tt.processParamsForQuery, tt.mockError).Once()
			s.mockDataService.On("TranslateQuery", ctx, mock.Anything, mock.Anything).Return("SELECT * FROM dataset1 RIGHT JOIN dataset2 ON dataset1.id = dataset2.id", tt.mockError).Once()

			jobTemplate, err := s.service.getCreateMVJobTemplate(ctx, tt.providerId, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetJobTemplate() {

	tests := []struct {
		name                string
		providerId          string
		payload             models.CreateActionPayload
		mockWarehouseId     string
		mockError           error
		expectedJobTemplate *jobs.SubmitRun
		expectError         bool
	}{
		{
			name:       "Error getting job template",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType:            serviceconstants.ActionTypeCreateMV,
				MerchantID:            "merchant1",
				ActionMetadataPayload: models.CreateActionResponse{},
			},
			mockWarehouseId: "warehouse1",
			mockError:       errors.ErrInvalidActionMetadataPayload,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDatabricksWarehouseId", ctx, mock.Anything).Return(tt.mockWarehouseId, nil).Once()
			jobTemplate, err := s.service.getJobTemplate(ctx, tt.providerId, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
				s.Nil(jobTemplate)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetActionDetails() {
	tests := []struct {
		name           string
		runId          int64
		workspaceId    string
		mockResponse   dataplatformmodels.QueryResult
		mockError      error
		expectedAction models.Action
		expectError    bool
	}{
		{
			name:        "Successful action retrieval",
			runId:       123,
			workspaceId: "workspace1",
			mockResponse: dataplatformmodels.QueryResult{
				Rows: dataplatformmodels.Rows{
					{
						"id":     "action1",
						"status": "SUCCESSFUL",
					},
				},
			},
			mockError: nil,
			expectedAction: models.Action{
				ID:           "action1",
				ActionStatus: serviceconstants.ActionStatusSuccessful,
			},
			expectError: false,
		},
		{
			name:         "Error retrieving action",
			runId:        456,
			workspaceId:  "workspace2",
			mockResponse: dataplatformmodels.QueryResult{},
			mockError:    errors.ErrGettingActionByRunIdFailed,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())
			s.mockDataService.On("GetDatabricksServiceForProvider", ctx, tt.workspaceId).Return(s.mockDatabricksService, nil).Once()
			s.mockDatabricksService.On("Query", ctx, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError).Once()

			action, err := s.service.getActionDetails(ctx, s.mockDatabricksService, tt.runId, tt.workspaceId)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
				s.Empty(action)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedAction, action)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestHandleJobStatusUpdate() {
	tests := []struct {
		name                 string
		runDetails           *jobs.Run
		runId                int64
		workspaceId          string
		mockError            error
		expectedError        error
		expectedActionStatus serviceconstants.ActionStatus
	}{
		{
			name: "Successful job status update",
			runDetails: &jobs.Run{
				State: &jobs.RunState{
					ResultState: jobs.RunResultStateSuccess,
				},
			},
			runId:                123,
			workspaceId:          "workspace1",
			mockError:            nil,
			expectedError:        nil,
			expectedActionStatus: serviceconstants.ActionStatusSuccessful,
		},
		{
			name: "Failed job status update",
			runDetails: &jobs.Run{
				State: &jobs.RunState{
					ResultState: jobs.RunResultStateFailed,
				},
			},
			runId:                456,
			workspaceId:          "workspace2",
			mockError:            nil,
			expectedError:        nil,
			expectedActionStatus: serviceconstants.ActionStatusFailed,
		},
		{
			name: "Job status not success or failed",
			runDetails: &jobs.Run{
				State: &jobs.RunState{
					ResultState: "UNKNOWN",
				},
			},
			runId:                789,
			workspaceId:          "workspace3",
			mockError:            nil,
			expectedError:        errors.ErrJobStatusNotSuccessOrFailed,
			expectedActionStatus: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())
			s.mockDataService.On("GetDatabricksServiceForProvider", ctx, tt.workspaceId).Return(s.mockDatabricksService, nil).Once()

			if tt.expectedError == nil {
				s.mockDatabricksService.On("Query", ctx, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{}, nil).Once()
			} else {
				s.mockDatabricksService.On("Query", ctx, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{}, tt.mockError).Once()
			}

			err := s.service.handleJobStatusUpdate(ctx, s.mockDatabricksService, tt.runDetails, tt.runId, tt.workspaceId)

			if tt.expectedError != nil {
				s.Error(err)
				s.Equal(tt.expectedError, err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestSaveAction() {
	tests := []struct {
		name           string
		payload        models.CreateActionPayload
		mockProviderId string
		mockError      error
		expectError    bool
	}{
		{
			name: "Successful action save",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeCreateMV,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.CreateMVActionPayload{
					Query: "SELECT * FROM table",
				},
			},
			mockProviderId: "workspace1",
			expectError:    false,
		},
		{
			name: "Error saving action",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeCreateMV,
				MerchantID: "merchant1",
			},
			mockProviderId: "workspace1",
			mockError:      errors.ErrQueryingDatabricksFailed,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			actionId := helper.GenerateUUIDWithUnderscores()

			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())

			if !tt.expectError {
				s.mockDatabricksService.On("Query", ctx, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{}, nil).Once()
			} else {
				s.mockDatabricksService.On("Query", ctx, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{}, tt.mockError).Once()
			}

			err := s.service.saveAction(ctx, s.mockDatabricksService, actionId, tt.mockProviderId, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestHandleOneTimeJobAction() {
	tests := []struct {
		name            string
		providerId      string
		payload         models.CreateActionPayload
		mockJobTemplate *jobs.SubmitRun
		mockError       error
		expectedRunId   int64
		expectError     bool
		expectedError   error
	}{
		{
			name:       "Successful job submission",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeUpdateDatasetData,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.CreateMVActionPayload{
					Query: "SELECT * FROM table",
				},
			},
			mockJobTemplate: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "task1",
					},
				},
			},
			expectedRunId: 0,
			expectError:   true,
			expectedError: errors.ErrInvalidActionType,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()

			runId, err := s.service.handleOneTimeJobAction(ctx, s.mockDatabricksService, tt.providerId, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.expectedError, err)
				s.Equal(models.SubmitActionResponse{}, runId)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedRunId, runId)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestHandleJobAction() {
	tests := []struct {
		name            string
		providerId      string
		payload         models.CreateActionPayload
		mockJobTemplate *jobs.SubmitRun
		mockError       error
		expectedRunId   int64
		expectError     bool
		expectedError   error
	}{
		{
			name:       "Successful job submission",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeCreateMV,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.CreateMVActionPayload{
					Query: "SELECT * FROM table",
				},
			},
			mockJobTemplate: &jobs.SubmitRun{
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "task1",
					},
				},
			},
			expectedRunId: 0,
			expectError:   true,
			expectedError: errors.ErrInvalidActionType,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()

			runId, err := s.service.handleJobAction(ctx, s.mockDatabricksService, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.expectedError, err)
				s.Equal(models.SubmitActionResponse{}, runId)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedRunId, runId)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetRegisterDatasetJobTemplate() {
	tests := []struct {
		name                string
		providerId          string
		payload             models.CreateActionPayload
		mockError           error
		expectedJobTemplate *jobs.SubmitRun
		expectError         bool
	}{
		{
			name:       "Successful job template creation",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeRegisterDataset,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.RegisterDatasetActionPayload{
					MerchantId: "merchant1",
					DatasetId:  "dataset1",
					DatasetConfig: dataservicemodels.DatasetConfig{
						Columns: map[string]dataservicemodels.DatasetColumnConfig{
							"column1": {
								CustomType: constants.DatabricksColumnCustomTypeCurrency,
							},
						},
					},
					Provider: dataserviceconstants.ProviderDatabricks,
				},
			},
			mockError: nil,
			expectedJobTemplate: &jobs.SubmitRun{
				RunName: "register_dataset_dataset1",
				WebhookNotifications: &jobs.WebhookNotifications{
					OnStart: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnSuccess: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnFailure: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
				},
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "register_dataset_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/register/dataset/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.RegisterDatasetParams:  `{"merchant_id":"merchant1","dataset_id":"dataset1","dataset_config":{"columns":{"column1":{"custom_type":"currency","custom_type_config":null}},"custom_column_groups":null,"rules":null,"computed_columns":null},"databricks_config":{"dedup_columns":null,"order_by_column":"","partition_columns":null,"cluster_columns":null},"provider":"databricks"}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
				},
				Queue: &jobs.QueueSettings{
					Enabled: true,
				},
				RunAs: &jobs.JobRunAs{
					UserName: "test_user",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())

			jobTemplate, err := s.service.getRegisterDatasetJobTemplate(ctx, tt.payload)
			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetRegisterJobJobTemplate() {

	tests := []struct {
		name                string
		providerId          string
		payload             models.CreateActionPayload
		mockError           error
		expectedJobTemplate *jobs.SubmitRun
		expectError         bool
	}{
		{
			name:       "Successful job template creation",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeRegisterJob,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.RegisterJobActionPayload{
					MerchantId:       "merchant1",
					JobType:          constants.DatabricksJobTypeTransformation,
					SourceType:       constants.DatabricksJobSourceTypeFolder,
					SourceValue:      "folder1",
					DestinationType:  constants.DatabricksJobDestinationTypeDataset,
					DestinationValue: "dataset1",
					TemplateId:       "template1",
				},
			},
			mockError: nil,
			expectedJobTemplate: &jobs.SubmitRun{
				RunName: "register_job_dataset1",
				WebhookNotifications: &jobs.WebhookNotifications{
					OnStart: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnSuccess: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnFailure: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
				},
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "register_job_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/register/job/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.RegisterJobParams:      `{"merchant_id":"merchant1","job_type":"transformation","source_type":"folder","source_value":"folder1","destination_type":"dataset","destination_value":"dataset1","template_id":"template1","quartz_cron_expression":""}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
				},
				Queue: &jobs.QueueSettings{
					Enabled: true,
				},
				RunAs: &jobs.JobRunAs{
					UserName: "test_user",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())

			jobTemplate, err := s.service.getRegisterJobJobTemplate(ctx, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetCopyDatasetJobTemplate() {
	tests := []struct {
		name                string
		payload             models.CreateActionPayload
		mockError           error
		expectedJobTemplate *jobs.SubmitRun
		expectError         bool
	}{
		{
			name: "Successful job template creation",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeCopyDataset,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.CopyDatasetActionPayload{
					OriginalDatasetId: "dataset1",
					NewDatasetId:      "dataset2",
					MerchantId:        "merchant1",
				},
			},
			mockError: nil,
			expectedJobTemplate: &jobs.SubmitRun{
				RunName: "copy_dataset_dataset2",
				WebhookNotifications: &jobs.WebhookNotifications{
					OnStart:   []jobs.Webhook{{Id: "webhook_id"}},
					OnSuccess: []jobs.Webhook{{Id: "webhook_id"}},
					OnFailure: []jobs.Webhook{{Id: "webhook_id"}},
				},
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "copy_dataset_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/copy/dataset/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.CopyDatasetParams:      `{"original_dataset_id":"dataset1","new_dataset_id":"dataset2","merchant_id":"merchant1"}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
				},
				Queue: &jobs.QueueSettings{Enabled: true},
				RunAs: &jobs.JobRunAs{UserName: "test_user"},
			},
			expectError: false,
		},
		{
			name: "Invalid payload type",
			payload: models.CreateActionPayload{
				ActionType:            serviceconstants.ActionTypeCopyDataset,
				MerchantID:            "merchant1",
				ActionMetadataPayload: "invalid_payload",
			},
			mockError:           errors.ErrInvalidActionMetadataPayload,
			expectedJobTemplate: nil,
			expectError:         true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())

			jobTemplate, err := s.service.getCopyDatasetJobTemplate(ctx, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
				s.Nil(jobTemplate)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetUpdateDatasetJobTemplate() {
	tests := []struct {
		name                string
		providerId          string
		payload             models.CreateActionPayload
		mockError           error
		expectedJobTemplate *jobs.SubmitRun
		expectError         bool
	}{
		{
			name:       "Successful job template creation",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeUpdateDataset,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.UpdateDatasetEvent{
					EventType: serviceconstants.UpdateDatasetEventTypeUpdateCustomColumn,
					EventData: models.UpdateDatasetActionPayload{
						DatasetId: "dataset1",
						DatasetConfig: dataservicemodels.DatasetConfig{
							Columns: map[string]dataservicemodels.DatasetColumnConfig{
								"column1": {
									CustomType: constants.DatabricksColumnCustomTypeCurrency,
								},
							},
						},
					},
					EventMetadata: models.UpsertRuleEventMetadata{},
				},
			},
			mockError: nil,
			expectedJobTemplate: &jobs.SubmitRun{
				RunName: "update_dataset_dataset1",
				WebhookNotifications: &jobs.WebhookNotifications{
					OnStart: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnSuccess: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnFailure: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
				},
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "update_dataset_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/update/dataset/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.UpdateDatasetEventParams: `{"event_type":"update_custom_column","event_data":{"dataset_id":"dataset1","dataset_config":{"columns":{"column1":{"custom_type":"currency","custom_type_config":null}},"custom_column_groups":null,"rules":null,"computed_columns":null}},"event_metadata":{"delta_rule_id":"","column":"","type":""}}`,
								serviceconstants.DatasetIdParam:           "dataset1",
								serviceconstants.DataPlatformModulesSrc:   "test_src",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
					{
						TaskKey: "pinot_ingestion_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/pinot/ingestion/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.UpdateDatasetEventParams: `{"event_type":"update_custom_column","event_data":{"dataset_id":"dataset1","dataset_config":{"columns":{"column1":{"custom_type":"currency","custom_type_config":null}},"custom_column_groups":null,"rules":null,"computed_columns":null}},"event_metadata":{"delta_rule_id":"","column":"","type":""}}`,
								serviceconstants.DatasetIdParam:           "dataset1",
								serviceconstants.DataPlatformModulesSrc:   "test_src",
							},
						},
						DependsOn: []jobs.TaskDependency{
							{
								TaskKey: "update_dataset_task",
							},
						},
						RunIf:                jobs.RunIfAllSuccess,
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
						ExistingClusterId:    "1111",
					},
				},
				Queue: &jobs.QueueSettings{
					Enabled: true,
				},
				RunAs: &jobs.JobRunAs{
					UserName: "test_user",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())

			jobTemplate, err := s.service.getUpdateDatasetJobTemplate(ctx, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetUpsertTemplateJobTemplate() {
	tests := []struct {
		name                string
		providerId          string
		payload             models.CreateActionPayload
		mockError           error
		expectedJobTemplate *jobs.SubmitRun
		expectError         bool
	}{
		{
			name:       "Successful job template creation",
			providerId: "workspace1",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeUpsertTemplate,
				MerchantID: "merchant1",
				ActionMetadataPayload: models.UpsertTemplateActionPayload{
					Id:            "template1",
					Name:          "template1",
					Configuration: "1234",
					TemplateType:  "join",
				},
			},
			mockError: nil,
			expectedJobTemplate: &jobs.SubmitRun{
				RunName: "upsert_template_template1",
				WebhookNotifications: &jobs.WebhookNotifications{
					OnStart: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnSuccess: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
					OnFailure: []jobs.Webhook{
						{
							Id: "webhook_id",
						},
					},
				},
				Tasks: []jobs.SubmitTask{
					{
						TaskKey: "upsert_template_task",
						NotebookTask: &jobs.NotebookTask{
							NotebookPath: "/path/to/upsert/template/notebook",
							Source:       jobs.SourceWorkspace,
							BaseParameters: map[string]string{
								serviceconstants.UpsertTemplateParams:   `{"id":"template1","name":"template1","template_type":"join","configuration":"1234"}`,
								serviceconstants.DataPlatformModulesSrc: "test_src",
							},
						},
						TimeoutSeconds:       0,
						EmailNotifications:   &jobs.JobEmailNotifications{},
						WebhookNotifications: &jobs.WebhookNotifications{},
					},
				},
				Queue: &jobs.QueueSettings{
					Enabled: true,
				},
				RunAs: &jobs.JobRunAs{
					UserName: "test_user",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())

			jobTemplate, err := s.service.getUpsertTemplateJobTemplate(ctx, tt.payload)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedJobTemplate, jobTemplate)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestGetActionById() {
	tests := []struct {
		name           string
		merchantId     string
		actionId       string
		mockResponse   dataplatformmodels.QueryResult
		mockError      error
		expectedAction models.Action
		expectError    bool
	}{
		{
			name:       "Successful action retrieval",
			merchantId: "merchant1",
			actionId:   "action1",
			mockResponse: dataplatformmodels.QueryResult{
				Rows: dataplatformmodels.Rows{
					{
						"id":     "action1",
						"status": "SUCCESSFUL",
					},
				},
			},
			mockError: nil,
			expectedAction: models.Action{
				ID:           "action1",
				ActionStatus: serviceconstants.ActionStatusSuccessful,
			},
			expectError: false,
		},
		{
			name:         "Error retrieving action",
			merchantId:   "merchant1",
			actionId:     "action2",
			mockResponse: dataplatformmodels.QueryResult{},
			mockError:    errors.ErrGettingActionByRunIdFailed,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDataPlatformConfig").Return(getDataPlatformMockConfig())
			s.mockDataService.On("GetDatabricksServiceForMerchant", ctx, tt.merchantId).Return(s.mockDatabricksService, nil).Once()
			s.mockDatabricksService.On("Query", ctx, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError).Once()

			action, err := s.service.GetActionById(ctx, tt.merchantId, tt.actionId)

			if tt.expectError {
				s.Error(err)
				s.Equal(tt.mockError, err)
				s.Empty(action)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedAction, action)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestVerifyCountryColumn() {
	tests := []struct {
		name          string
		countryColumn interface{}
		expected      bool
		expectError   bool
		mockError     error
	}{
		{
			name:          "Valid country column",
			countryColumn: "IND",
			expected:      true,
			expectError:   false,
		},
		{
			name:          "Invalid country column",
			countryColumn: "XYZ",
			expected:      false,
			expectError:   true,
			mockError:     errors.ErrInvalidCountryValue,
		},
		{
			name:          "Empty country column",
			countryColumn: "",
			expected:      false,
			expectError:   true,
			mockError:     errors.ErrInvalidCountryValue,
		},
		{
			name:          "Nil country column",
			countryColumn: nil,
			expected:      false,
			expectError:   true,
			mockError:     errors.ErrInvalidCountryValue,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			result := verifyCountryColumn(ctx, tt.countryColumn)
			if tt.expectError {
				s.Error(result)
				s.Equal(tt.mockError, result)
			} else {
				s.NoError(result)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestVerifyCurrencyColumn() {
	tests := []struct {
		name           string
		currencyColumn interface{}
		expected       bool
		expectError    bool
		mockError      error
	}{
		{
			name:           "Valid currency column",
			currencyColumn: "INR",
			expected:       true,
			expectError:    false,
		},
		{
			name:           "Invalid currency column",
			currencyColumn: "XYZ",
			expected:       false,
			expectError:    true,
			mockError:      errors.ErrInvalidCurrencyValue,
		},
		{
			name:           "Empty currency column",
			currencyColumn: "",
			expected:       false,
			expectError:    true,
			mockError:      errors.ErrInvalidCurrencyValue,
		},
		{
			name:           "Nil currency column",
			currencyColumn: nil,
			expected:       false,
			expectError:    true,
			mockError:      errors.ErrInvalidCurrencyValue,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			result := verifyCurrencyColumn(ctx, tt.currencyColumn)
			if tt.expectError {
				s.Error(result)
				s.Equal(tt.mockError, result)
			} else {
				s.NoError(result)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestVerifyBankColumn() {
	tests := []struct {
		name        string
		bankColumn  interface{}
		expected    bool
		expectError bool
		mockError   error
	}{
		{
			name:        "Valid bank column",
			bankColumn:  "CRB",
			expected:    true,
			expectError: false,
		},
		{
			name:        "Invalid bank column",
			bankColumn:  "XYZ",
			expected:    false,
			expectError: true,
			mockError:   errors.ErrInvalidBankValue,
		},
		{
			name:        "Number bank column",
			bankColumn:  123,
			expected:    false,
			expectError: true,
			mockError:   errors.ErrInvalidBankValue,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			result := verifyBankColumn(ctx, tt.bankColumn)
			if tt.expectError {
				s.Error(result)
				s.Equal(tt.mockError, result)
			} else {
				s.NoError(result)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestVerifyAmountColumn() {
	tests := []struct {
		name         string
		amountColumn interface{}
		expected     bool
		expectError  bool
		mockError    error
	}{
		{
			name:         "Valid amount column",
			amountColumn: 1001111,
			expected:     true,
			expectError:  false,
		},
		{
			name:         "Invalid amount column",
			amountColumn: "XYZ",
			expected:     false,
			expectError:  true,
			mockError:    errors.ErrInvalidAmountValue,
		},
		{
			name:         "Empty amount column",
			amountColumn: "",
			expected:     false,
			expectError:  true,
			mockError:    errors.ErrInvalidAmountValue,
		},
		{
			name:         "Nil amount column",
			amountColumn: nil,
			expected:     false,
			expectError:  true,
			mockError:    errors.ErrInvalidAmountValue,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			result := verifyAmountColumn(ctx, tt.amountColumn)
			if tt.expectError {
				s.Error(result)
				s.Equal(tt.mockError, result)
			} else {
				s.NoError(result)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestVerifyTagsColumn() {
	tests := []struct {
		name        string
		tagsColumn  interface{}
		expected    bool
		expectError bool
		mockError   error
	}{
		{
			name:        "Valid tags column",
			tagsColumn:  "tag1.tag2",
			expected:    true,
			expectError: false,
		},
		{
			name:        "Invalid tags column",
			tagsColumn:  1,
			expected:    false,
			expectError: true,
			mockError:   errors.ErrInvalidTagsValue,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			result := verifyTagsColumn(ctx, tt.tagsColumn)
			if tt.expectError {
				s.Error(result)
				s.Equal(tt.mockError, result)
			} else {
				s.NoError(result)
			}
		})
	}
}

func (s *ActionServiceTestSuite) TestHandleValidationsAndSourceUpdates() {
	tests := []struct {
		name        string
		payload     models.CreateActionPayload
		expected    error
		expectError bool
	}{
		{
			name: "Valid payload Handle Custom Column",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeUpdateDatasetData,
				ActionMetadataPayload: models.UpdateDatasetDataActionPayload{
					DatasetId:    "dataset1",
					SqlCondition: "id = 1",
					UpdateValues: map[string]any{
						"bank":       "CB",
						"currency":   "INR",
						"random_col": "y",
					},
				},
			},
			expected:    nil,
			expectError: false,
		},
		{
			name: "Valid payload Handle Normal Column",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeUpdateDatasetData,
				ActionMetadataPayload: models.UpdateDatasetDataActionPayload{
					DatasetId:    "dataset1",
					SqlCondition: "id = 1",
					UpdateValues: map[string]any{
						"random_col_1": "y",
						"random_col_2": 2000,
						"random_col_3": "2025-01-01",
					},
				},
			},
			expected:    nil,
			expectError: false,
		},
		{
			name: "Invalid payload Handle Bank",
			payload: models.CreateActionPayload{
				ActionType: serviceconstants.ActionTypeUpdateDatasetData,
				ActionMetadataPayload: models.UpdateDatasetDataActionPayload{
					DatasetId:    "dataset1",
					SqlCondition: "id = 1",
					UpdateValues: map[string]any{
						"bank":       "CBk",
						"random_col": "y",
					},
				},
			},
			expected:    errors.ErrInvalidBankValue,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockDataService.On("GetDatasetConfig", ctx, mock.Anything, mock.Anything).Return(dataservicemodels.DatasetConfig{
				Columns: map[string]dataservicemodels.DatasetColumnConfig{
					"bank":     {CustomType: constants.DatabricksColumnCustomTypeBank},
					"currency": {CustomType: constants.DatabricksColumnCustomTypeCurrency},
					"amount":   {CustomType: constants.DatabricksColumnCustomTypeAmount},
					"tags":     {CustomType: constants.DatabricksColumnCustomTypeTags},
					"country":  {CustomType: constants.DatabricksColumnCustomTypeCountry},
				},
			}, nil)
			err := s.service.handleValidationsAndSourceUpdates(ctx, tt.payload)
			if tt.expectError {
				s.Error(err)
				s.Equal(tt.expected, err)
			} else {
				s.NoError(err)
			}
		})
	}
}
