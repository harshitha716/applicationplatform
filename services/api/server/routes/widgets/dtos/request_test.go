package dtos

import (
	"testing"

	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	"github.com/stretchr/testify/assert"
)

func TestWidgetQueryParamsToModels(t *testing.T) {
	// Test with empty filters and time columns
	t.Run("empty params", func(t *testing.T) {
		t.Parallel()
		params := WidgetQueryParams{
			Filters:     []WidgetFilters{},
			TimeColumns: []ColumnMapping{},
		}
		
		result := params.ToModels()
		
		assert.Empty(t, result.Filters)
		assert.Empty(t, result.TimeColumns)
		assert.Nil(t, result.Periodicity)
		assert.Nil(t, result.Currency)
	})
	
	// Test with populated filters and time columns
	t.Run("populated params", func(t *testing.T) {
		t.Parallel()
		periodicity := "monthly"
		currency := "USD"
		
		params := WidgetQueryParams{
			Filters: []WidgetFilters{
				{
					DatasetID: "dataset1",
					Filters: datasetmodels.FilterModel{
						Conditions: []datasetmodels.Filter{
							{
								Column:   "col1",
								Operator: "eq",
								Value:    "val1",
							},
						},
					},
					Pagination: &PaginationParams{
						Page:     1,
						PageSize: 10,
					},
				},
			},
			TimeColumns: []ColumnMapping{
				{
					DatasetID: "dataset1",
					Column:    "date_col",
				},
			},
			Periodicity: &periodicity,
			Currency:    &currency,
		}
		
		result := params.ToModels()
		
		assert.Len(t, result.Filters, 1)
		assert.Equal(t, "dataset1", result.Filters[0].DatasetID)
		assert.Len(t, result.Filters[0].Filters.Conditions, 1)
		assert.Equal(t, "col1", result.Filters[0].Filters.Conditions[0].Column)
		assert.NotNil(t, result.Filters[0].Pagination)
		assert.Equal(t, 1, result.Filters[0].Pagination.Page)
		assert.Equal(t, 10, result.Filters[0].Pagination.PageSize)
		
		assert.Len(t, result.TimeColumns, 1)
		assert.Equal(t, "dataset1", result.TimeColumns[0].DatasetID)
		assert.Equal(t, "date_col", result.TimeColumns[0].Column)
		
		assert.Equal(t, &periodicity, result.Periodicity)
		assert.Equal(t, &currency, result.Currency)
	})
	
	// Test with nil pagination
	t.Run("nil pagination", func(t *testing.T) {
		t.Parallel()
		params := WidgetQueryParams{
			Filters: []WidgetFilters{
				{
					DatasetID: "dataset1",
					Filters: datasetmodels.FilterModel{
						Conditions: []datasetmodels.Filter{},
					},
					Pagination: nil,
				},
			},
		}
		
		result := params.ToModels()
		
		assert.Len(t, result.Filters, 1)
		assert.Equal(t, "dataset1", result.Filters[0].DatasetID)
		assert.Nil(t, result.Filters[0].Pagination)
	})
}
