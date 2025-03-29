package fileimport

import (
	"context"
	"fmt"
	"testing"

	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	mock_datasetservice "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_fileimports "github.com/Zampfi/application-platform/services/api/mocks/core/fileimports"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/Zampfi/application-platform/services/api/workers/defaultworker/constants"
	mock_temporal "github.com/Zampfi/workflow-sdk-go/mocks/workflowmanagers/temporal"
	temporalmodels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func TestApplicationPlatformDatasetFileImportWorkflowExecute(t *testing.T) {

	testCases := []struct {
		name          string
		input         models.FileImportWorkflowInitPayload
		setupMocks    func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore)
		expectedError bool
	}{
		{
			name: "Missing User ID",
			input: models.FileImportWorkflowInitPayload{
				OrganizationId: uuid.New(),
			},
			expectedError: true,
		},
		{
			name: "Missing Organization ID",
			input: models.FileImportWorkflowInitPayload{
				UserId: uuid.New(),
			},
			expectedError: true,
		},
		{
			name: "Successful Execution",
			input: models.FileImportWorkflowInitPayload{
				UserId:         uuid.New(),
				OrganizationId: uuid.New(),
				FileUploadId:   uuid.New(),
				DatasetId:      uuid.New(),
			},
			setupMocks: func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore) {
				fileUpload := dbmodels.FileUpload{
					StorageFilePath: "test/path",
					StorageBucket:   "test-bucket",
				}

				mockFileUploadsService.EXPECT().
					GetFileUploadByIds(mock.Anything, mock.Anything).
					Return([]dbmodels.FileUpload{fileUpload}, nil).
					Once()

				mockDatasetService.EXPECT().
					GetDatasetImportPath(mock.Anything, mock.Anything, mock.Anything).
					Return(&models.FileImportConfig{
						BronzeSourcePath:   "bronze/path",
						BronzeSourceBucket: "bronze-bucket",
					}, nil).
					Once()

				workflowResponse := map[string]interface{}{
					"transformed_data_path":   "transformed/path",
					"transformed_data_bucket": "transformed-bucket",
					"column_mapping":          map[string]interface{}{},
					"data_preview":            dbmodels.DatasetPreview{},
					"extracted_metadata":      map[string]interface{}{},
				}

				mockTemporalSdk.EXPECT().
					ExecuteSyncWorkflow(mock.Anything, mock.Anything).
					Return(temporalmodels.WorkflowResponse{}, nil).
					Run(func(ctx context.Context, params temporalmodels.ExecuteWorkflowParams) {
						*params.ResultPtr.(*map[string]interface{}) = workflowResponse
					}).
					Once()

				mockDatasetService.EXPECT().
					UpdateDatasetActionStatus(mock.Anything, mock.Anything, constants.DatasetActionStatusSuccessful).
					Return(nil).
					Once()

				mockDatasetService.EXPECT().
					UpdateDatasetFileUploadStatus(mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		s := testsuite.WorkflowTestSuite{}
		env := s.NewTestWorkflowEnvironment()

		mockDatasetService := mock_datasetservice.NewMockDatasetService(t)
		mockFileUploadsService := mock_fileimports.NewMockFileImportService(t)
		mockTemporalSdk := mock_temporal.NewMockTemporalService(t)
		mockDatasetStore := mock_store.NewMockDatasetStore(t)

		workflow := initFileImport(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
			}

			env.RegisterActivity(workflow.ExecuteAINormalizationActivity)
			env.RegisterActivity(workflow.RegisterExitPayloadInDB)
			env.ExecuteWorkflow(workflow.ApplicationPlatformDatasetFileImportWorkflowExecute, tc.input)

			result := ""
			err := env.GetWorkflowResult(&result)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestExecuteAINormalizationActivity(t *testing.T) {

	testCases := []struct {
		name          string
		input         models.FileImportWorkflowInitPayload
		setupMocks    func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore)
		expectedError bool
	}{
		{
			name: "File Upload Not Found",
			input: models.FileImportWorkflowInitPayload{
				FileUploadId: uuid.New(),
			},
			setupMocks: func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore) {
				mockFileUploadsService.EXPECT().
					GetFileUploadByIds(mock.Anything, mock.Anything).
					Return([]dbmodels.FileUpload{}, nil)
			},
			expectedError: true,
		},
		{
			name: "Successful Execution",
			input: models.FileImportWorkflowInitPayload{
				FileUploadId: uuid.New(),
				DatasetId:    uuid.New(),
			},
			setupMocks: func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore) {
				fileUpload := dbmodels.FileUpload{
					StorageFilePath: "test/path",
					StorageBucket:   "test-bucket",
				}

				mockFileUploadsService.EXPECT().
					GetFileUploadByIds(mock.Anything, mock.Anything).
					Return([]dbmodels.FileUpload{fileUpload}, nil)

				mockDatasetService.EXPECT().
					GetDatasetImportPath(mock.Anything, mock.Anything, mock.Anything).
					Return(&models.FileImportConfig{
						BronzeSourcePath:   "bronze/path",
						BronzeSourceBucket: "bronze-bucket",
					}, nil)

				workflowResponse := map[string]interface{}{
					"transformed_data_path":   "transformed/path",
					"transformed_data_bucket": "transformed-bucket",
					"column_mapping":          map[string]interface{}{},
					"data_preview":            dbmodels.DatasetPreview{},
					"extracted_metadata":      map[string]interface{}{},
				}

				mockTemporalSdk.EXPECT().
					ExecuteSyncWorkflow(mock.Anything, mock.Anything).
					Return(temporalmodels.WorkflowResponse{}, nil).
					Run(func(ctx context.Context, params temporalmodels.ExecuteWorkflowParams) {
						*params.ResultPtr.(*map[string]interface{}) = workflowResponse
					})
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDatasetService := mock_datasetservice.NewMockDatasetService(t)
			mockFileUploadsService := mock_fileimports.NewMockFileImportService(t)
			mockTemporalSdk := mock_temporal.NewMockTemporalService(t)
			mockDatasetStore := mock_store.NewMockDatasetStore(t)

			workflow := initFileImport(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
			if tc.setupMocks != nil {
				tc.setupMocks(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
			}

			result, err := workflow.ExecuteAINormalizationActivity(context.Background(), tc.input)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestRegisterExitPayloadInDB(t *testing.T) {

	testCases := []struct {
		name          string
		initPayload   models.FileImportWorkflowInitPayload
		exitPayload   models.FileImportWorkflowExitPayload
		setupMocks    func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore)
		expectedError bool
	}{
		{
			name: "Successful Registration",
			initPayload: models.FileImportWorkflowInitPayload{
				DatasetActionId:     uuid.New(),
				DatasetFileUploadId: uuid.New(),
			},
			exitPayload: models.FileImportWorkflowExitPayload{
				TransformedFilePath: "transformed/path",
				NormalizationResult: models.AINormalizationWorkflowResponse{
					TransformedDataPath:   "transformed/path",
					TransformedDataBucket: "transformed-bucket",
				},
			},
			setupMocks: func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore) {
				mockDatasetService.EXPECT().
					UpdateDatasetActionStatus(mock.Anything, mock.Anything, constants.DatasetActionStatusSuccessful).
					Return(nil)

				mockDatasetService.EXPECT().
					UpdateDatasetFileUploadStatus(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Failed to Update Dataset Action Status",
			initPayload: models.FileImportWorkflowInitPayload{
				DatasetActionId:     uuid.New(),
				DatasetFileUploadId: uuid.New(),
			},
			exitPayload: models.FileImportWorkflowExitPayload{
				TransformedFilePath: "transformed/path",
			},
			setupMocks: func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore) {
				mockDatasetService.EXPECT().
					UpdateDatasetActionStatus(mock.Anything, mock.Anything, constants.DatasetActionStatusSuccessful).
					Return(fmt.Errorf("failed to update dataset action status"))
			},
			expectedError: true,
		},
		{
			name: "Failed to Update Dataset File Upload Status",
			initPayload: models.FileImportWorkflowInitPayload{
				DatasetActionId:     uuid.New(),
				DatasetFileUploadId: uuid.New(),
			},
			exitPayload: models.FileImportWorkflowExitPayload{
				TransformedFilePath: "transformed/path",
				NormalizationResult: models.AINormalizationWorkflowResponse{
					TransformedDataPath:   "transformed/path",
					TransformedDataBucket: "transformed-bucket",
				},
			},
			setupMocks: func(mockDatasetService *mock_datasetservice.MockDatasetService, mockFileUploadsService *mock_fileimports.MockFileImportService, mockTemporalSdk *mock_temporal.MockTemporalService, mockDatasetStore *mock_store.MockDatasetStore) {
				mockDatasetService.EXPECT().
					UpdateDatasetActionStatus(mock.Anything, mock.Anything, constants.DatasetActionStatusSuccessful).
					Return(nil)

				mockDatasetService.EXPECT().
					UpdateDatasetFileUploadStatus(mock.Anything, mock.Anything, mock.Anything).
					Return(fmt.Errorf("failed to update dataset file upload status"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDatasetService := mock_datasetservice.NewMockDatasetService(t)
			mockFileUploadsService := mock_fileimports.NewMockFileImportService(t)
			mockTemporalSdk := mock_temporal.NewMockTemporalService(t)
			mockDatasetStore := mock_store.NewMockDatasetStore(t)

			workflow := initFileImport(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
			if tc.setupMocks != nil {
				tc.setupMocks(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
			}

			result, err := workflow.RegisterExitPayloadInDB(context.Background(), tc.initPayload, tc.exitPayload)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.exitPayload.TransformedFilePath, result)
			}
		})
	}
}

// TestInitFileImportWorkflow tests the InitFileImportWorkflow function
// Skipping this test as it requires a fully configured server config
// which is difficult to mock properly in a test environment
func TestInitFileImportWorkflow(t *testing.T) {
	t.Skip("Skipping test that requires a fully configured server config")
}

// TestGetActivities tests the GetActivities function
func TestGetActivities(t *testing.T) {
	// Create mocks
	mockDatasetService := mock_datasetservice.NewMockDatasetService(t)
	mockFileUploadsService := mock_fileimports.NewMockFileImportService(t)
	mockTemporalSdk := mock_temporal.NewMockTemporalService(t)
	mockDatasetStore := mock_store.NewMockDatasetStore(t)
	
	// Create a workflow instance
	workflow := initFileImport(mockDatasetService, mockFileUploadsService, mockTemporalSdk, mockDatasetStore)
	
	// Call the function
	activities := workflow.GetActivities()
	
	// Verify that the activities are returned correctly
	assert.NotNil(t, activities)
	assert.Len(t, activities, 2) // Should have two activities
	
	// Verify the activity functions
	for _, activity := range activities {
		assert.NotNil(t, activity.Function)
		assert.True(t, activity.RegisterOptions.DisableAlreadyRegisteredCheck)
	}
}
