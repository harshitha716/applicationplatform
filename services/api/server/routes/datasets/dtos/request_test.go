package dtos

import (
	"testing"

	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDatasetRequest_ToModel(t *testing.T) {
	description := "test description"
	alias := "Test Column"
	displayConfig := []datasetmodels.DisplayConfig{
		{
			Column: "test_column",
			Alias:  &alias,
		},
	}
	mvConfig := &datasetmodels.MVConfig{}

	request := RegisterDatasetRequest{
		DatasetTitle:       "Test Dataset",
		DatasetDescription: &description,
		DatasetType:        storemodels.DatasetTypeBronze,
		DatasetConfig: dataplatformDataModels.DatasetConfig{
			Columns: map[string]dataplatformDataModels.DatasetColumnConfig{
				"test_column": {
					CustomType: "string",
				},
			},
		},
		DatabricksConfig: dataplatformDataModels.DatabricksConfig{
			DedupColumns: []string{"test_column"},
		},
		DisplayConfig: displayConfig,
		MVConfig:      mvConfig,
		Provider:      dataplatformdataconstants.ProviderDatabricks,
	}

	result := request.ToModel()

	assert.Equal(t, request.DatasetTitle, result.DatasetTitle)
	assert.Equal(t, request.DatasetDescription, result.DatasetDescription)
	assert.Equal(t, request.DatasetType, result.DatasetType)
	assert.Equal(t, request.DatasetConfig, result.DatasetConfig)
	assert.Equal(t, request.DatabricksConfig, result.DatabricksConfig)
	assert.Equal(t, request.DisplayConfig, result.DisplayConfig)
	assert.Equal(t, request.MVConfig, result.MVConfig)
	assert.Equal(t, request.Provider, result.Provider)
}
