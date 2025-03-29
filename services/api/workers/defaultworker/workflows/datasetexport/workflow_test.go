package datasetexport

import (
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	mock_service "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func TestDatasetExportWorkflow(t *testing.T) {
	tests := []struct {
		name       string
		params     models.DatasetExportParams
		userId     uuid.UUID
		orgIds     []uuid.UUID
		workflowId string
		datasetId  uuid.UUID
		mockSetup  func(*mock_service.MockDatasetService)
		wantResult string
		wantErr    bool
	}{
		{
			name:       "Successfully exports dataset",
			params:     models.DatasetExportParams{},
			userId:     uuid.New(),
			orgIds:     []uuid.UUID{uuid.New()},
			workflowId: "test-workflow-id",
			datasetId:  uuid.New(),
			mockSetup: func(m *mock_service.MockDatasetService) {
				m.On("DatasetExportTemporalActivity", mock.Anything, mock.AnythingOfType("models.DatasetExportParams"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("string")).
					Return("exported/path/file.csv", nil)
			},
			wantResult: "exported/path/file.csv",
			wantErr:    false,
		},
		{
			name:       "Returns error when user ID is nil",
			params:     models.DatasetExportParams{},
			userId:     uuid.Nil,
			orgIds:     []uuid.UUID{uuid.New()},
			workflowId: "test-workflow-id",
			datasetId:  uuid.New(),
			mockSetup:  func(m *mock_service.MockDatasetService) {},
			wantResult: "",
			wantErr:    true,
		},
		{
			name:       "Returns error when org IDs are empty",
			params:     models.DatasetExportParams{},
			userId:     uuid.New(),
			orgIds:     []uuid.UUID{},
			workflowId: "test-workflow-id",
			datasetId:  uuid.New(),
			mockSetup:  func(m *mock_service.MockDatasetService) {},
			wantResult: "",
			wantErr:    true,
		},
		{
			name:       "Returns error when export activity fails",
			params:     models.DatasetExportParams{},
			userId:     uuid.New(),
			orgIds:     []uuid.UUID{uuid.New()},
			workflowId: "test-workflow-id",
			datasetId:  uuid.New(),
			mockSetup: func(m *mock_service.MockDatasetService) {
				m.On("DatasetExportTemporalActivity", mock.Anything, mock.AnythingOfType("models.DatasetExportParams"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("string")).
					Return("", assert.AnError)
			},
			wantResult: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestWorkflowEnvironment()

			mockDatasetService := mock_service.NewMockDatasetService(t)
			workflow := datasetExportWorkflow{
				datasetService: mockDatasetService,
			}

			// Setup mocks
			tt.mockSetup(mockDatasetService)

			env.RegisterActivity(workflow.datasetService.DatasetExportTemporalActivity)

			var result interface{}
			env.ExecuteWorkflow(workflow.ApplicationPlatformDatasetExportWorkflowExecute, tt.params, tt.datasetId, tt.userId, tt.orgIds, tt.workflowId)

			err := env.GetWorkflowResult(&result)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResult, result)
			mockDatasetService.AssertExpectations(t)
		})
	}
}

// TestInitDatasetExport tests the InitDatasetExport function
func TestInitDatasetExport(t *testing.T) {
	// Create a mock server config
	mockServerConfig := &serverconfig.ServerConfig{}
	
	// Create a mock dataset service
	mockDatasetService := mock_service.NewMockDatasetService(t)
	
	// Call the function
	workflow := InitDatasetExport(mockServerConfig, mockDatasetService)
	
	// Verify that the workflow is initialized correctly
	assert.NotNil(t, workflow)
	assert.Equal(t, mockDatasetService, workflow.datasetService)
}

// TestGetActivities tests the GetActivities function
func TestGetActivities(t *testing.T) {
	// Create a mock dataset service
	mockDatasetService := mock_service.NewMockDatasetService(t)
	
	// Create a workflow instance
	workflow := datasetExportWorkflow{
		datasetService: mockDatasetService,
	}
	
	// Call the function
	activities := workflow.GetActivities()
	
	// Verify that the activities are returned correctly
	assert.NotNil(t, activities)
	assert.Len(t, activities, 1)
	
	// Verify the activity function
	assert.NotNil(t, activities[0].Function)
	
	// Verify the register options
	assert.Equal(t, "DatasetExportTemporalActivity", activities[0].RegisterOptions.Name)
}
