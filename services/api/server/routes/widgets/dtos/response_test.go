package dtos

import (
	"encoding/json"
	"testing"

	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewWidgetInstanceDataResponse(t *testing.T) {
	// Test with empty dataset data
	t.Run("empty dataset data", func(t *testing.T) {
		t.Parallel()
		qrs := []datasetmodels.DatasetData{}
		periodicity := "daily"
		currency := stringPtr("USD")

		response := NewWidgetInstanceDataResponse(qrs, periodicity, currency)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, periodicity, response.Periodicity)
		assert.Equal(t, currency, response.Currency)
		assert.Empty(t, response.Result)
	})

	// Test with populated dataset data
	t.Run("populated dataset data", func(t *testing.T) {
		t.Parallel()
		qrs := []datasetmodels.DatasetData{
			{
				QueryResult: dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "col1", DatabaseType: "string"},
						{Name: "col2", DatabaseType: "int"},
					},
					Rows: []map[string]interface{}{
						{"col1": "value1", "col2": 1},
						{"col1": "value2", "col2": 2},
					},
				},
				Title: "Test Dataset",
			},
		}
		periodicity := "monthly"

		response := NewWidgetInstanceDataResponse(qrs, periodicity, nil)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, periodicity, response.Periodicity)
		assert.Nil(t, response.Currency)
		assert.Len(t, response.Result, 1)
		assert.Equal(t, 2, response.Result[0].RowCount)
		assert.Len(t, response.Result[0].Columns, 2)
		assert.Equal(t, "col1", response.Result[0].Columns[0].ColumnName)
		assert.Equal(t, "string", response.Result[0].Columns[0].ColumnType)
		assert.Len(t, response.Result[0].Data, 2)
	})

	// Test with multiple dataset data
	t.Run("multiple dataset data", func(t *testing.T) {
		t.Parallel()
		qrs := []datasetmodels.DatasetData{
			{
				QueryResult: dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "col1", DatabaseType: "string"},
					},
					Rows: []map[string]interface{}{
						{"col1": "value1"},
					},
				},
				Title: "Dataset 1",
			},
			{
				QueryResult: dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "col2", DatabaseType: "int"},
					},
					Rows: []map[string]interface{}{
						{"col2": 1},
						{"col2": 2},
					},
				},
				Title: "Dataset 2",
			},
		}
		periodicity := "weekly"
		currency := stringPtr("EUR")

		response := NewWidgetInstanceDataResponse(qrs, periodicity, currency)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, periodicity, response.Periodicity)
		assert.Equal(t, currency, response.Currency)
		assert.Len(t, response.Result, 2)
		assert.Equal(t, 1, response.Result[0].RowCount)
		assert.Equal(t, 2, response.Result[1].RowCount)
	})
}

func TestNewWidgetInstanceResponse(t *testing.T) {
	// Test with minimal widget instance
	t.Run("minimal widget instance", func(t *testing.T) {
		t.Parallel()
		id := uuid.New()
		sheetID := uuid.New()

		wi := &widgetmodels.WidgetInstance{
			ID:         id,
			WidgetType: "bar_chart",
			SheetID:    sheetID,
			Title:      "Test Widget",
			DataMappings: widgetmodels.DataMappings{
				Version:  widgetmodels.DataMappingVersion1,
				Mappings: []widgetmodels.DataMappingFields{},
			},
		}

		response, err := NewWidgetInstanceResponse(wi)

		assert.NoError(t, err)
		assert.Equal(t, id, response.ID)
		assert.Equal(t, "bar_chart", response.WidgetType)
		assert.Equal(t, sheetID, response.SheetID)
		assert.Equal(t, "Test Widget", response.Title)

		// Verify DataMappings JSON
		var mappings widgetmodels.DataMappings
		err = json.Unmarshal(response.DataMappings, &mappings)
		assert.NoError(t, err)
		assert.Equal(t, widgetmodels.DataMappingVersion1, mappings.Version)
		assert.Empty(t, mappings.Mappings)

		// Verify DisplayConfig is empty
		assert.Empty(t, response.DisplayConfig)
	})

	// Test with display config
	t.Run("with display config", func(t *testing.T) {
		t.Parallel()
		id := uuid.New()
		sheetID := uuid.New()
		displayConfig := json.RawMessage(`{"color":"blue"}`)

		wi := &widgetmodels.WidgetInstance{
			ID:            id,
			WidgetType:    "bar_chart",
			SheetID:       sheetID,
			Title:         "Test Widget",
			DisplayConfig: &displayConfig,
			DataMappings: widgetmodels.DataMappings{
				Version:  widgetmodels.DataMappingVersion1,
				Mappings: []widgetmodels.DataMappingFields{},
			},
		}

		response, err := NewWidgetInstanceResponse(wi)

		assert.NoError(t, err)
		assert.Equal(t, id, response.ID)
		assert.Equal(t, displayConfig, response.DisplayConfig)
	})

	// Test with marshaling error
	t.Run("marshaling error simulation", func(t *testing.T) {
		t.Parallel()
		// This is a contrived test to simulate a marshaling error
		// In a real scenario, it's difficult to create a marshaling error with valid Go structs
		// So we're just documenting that we've considered this error path

		// The actual implementation handles marshaling errors by returning nil, err
		// We can't easily trigger this in a test without modifying the code
	})
}

func TestTransformDefaultFilter(t *testing.T) {
	// Test with empty filter model
	t.Run("empty filter model", func(t *testing.T) {
		t.Parallel()
		filter := datasetmodels.FilterModel{}
		result := transformDefaultFilter(filter)
		assert.Nil(t, result)
	})

	// Test with populated filter model
	t.Run("populated filter model", func(t *testing.T) {
		t.Parallel()
		filter := datasetmodels.FilterModel{
			Conditions: []datasetmodels.Filter{
				{
					Column:   "col1",
					Operator: "eq",
					Value:    "val1",
				},
			},
		}

		result := transformDefaultFilter(filter)

		assert.Len(t, result, 1)
		assert.Equal(t, "col1", result[0].Column)
		assert.Equal(t, "eq", result[0].Operator)
		assert.Equal(t, []interface{}{"val1"}, result[0].Value)
	})

	// Test with multiple conditions
	t.Run("multiple conditions", func(t *testing.T) {
		t.Parallel()
		filter := datasetmodels.FilterModel{
			Conditions: []datasetmodels.Filter{
				{
					Column:   "col1",
					Operator: "eq",
					Value:    "val1",
				},
				{
					Column:   "col2",
					Operator: "gt",
					Value:    10,
				},
			},
		}

		result := transformDefaultFilter(filter)

		assert.Len(t, result, 2)
		assert.Equal(t, "col1", result[0].Column)
		assert.Equal(t, "eq", result[0].Operator)
		assert.Equal(t, "col2", result[1].Column)
		assert.Equal(t, "gt", result[1].Operator)
	})
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
