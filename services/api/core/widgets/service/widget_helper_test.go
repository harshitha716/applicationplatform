package widgets

import (
	"testing"

	dataplatformdataConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/stretchr/testify/assert"
)

func TestFlattenTags(t *testing.T) {
	service := &widgetsService{}

	tests := []struct {
		name      string
		input     *datasetmodels.DatasetData
		tagColumn string
		want      datasetmodels.DatasetData
	}{
		{
			name: "complex mixed depth tags",
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "tags": "electronics.computers.laptop"},
						{"id": 2, "tags": "electronics.phones"},
						{"id": 3, "tags": "clothing"},
						{"id": 4, "tags": "electronics.computers.desktop.gaming"},
						{"id": 5, "tags": "furniture.outdoor.chairs"},
						{"id": 6, "tags": "electronics"},
						{"id": 7, "tags": "books.fiction.fantasy.young-adult"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "tags", DatabaseType: "string"},
					},
				},
			},
			tagColumn: "tags",
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "__tags_LEVEL_1": "electronics", "__tags_LEVEL_2": "computers", "__tags_LEVEL_3": "laptop"},
						{"id": 2, "__tags_LEVEL_1": "electronics", "__tags_LEVEL_2": "phones"},
						{"id": 3, "__tags_LEVEL_1": "clothing"},
						{"id": 4, "__tags_LEVEL_1": "electronics", "__tags_LEVEL_2": "computers", "__tags_LEVEL_3": "desktop", "__tags_LEVEL_4": "gaming"},
						{"id": 5, "__tags_LEVEL_1": "furniture", "__tags_LEVEL_2": "outdoor", "__tags_LEVEL_3": "chairs"},
						{"id": 6, "__tags_LEVEL_1": "electronics"},
						{"id": 7, "__tags_LEVEL_1": "books", "__tags_LEVEL_2": "fiction", "__tags_LEVEL_3": "fantasy", "__tags_LEVEL_4": "young-adult"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "tags", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						{Name: "__tags_LEVEL_1", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						{Name: "__tags_LEVEL_2", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						{Name: "__tags_LEVEL_3", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						{Name: "__tags_LEVEL_4", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
		{
			name: "with nil tag values",
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "tags": nil},
						{"id": 2, "tags": "electronics"},
						{"id": 3, "tags": "apple"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "tags", DatabaseType: "string"},
					},
				},
			},
			tagColumn: "tags",
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "__tags_LEVEL_1": "__UNTAGGED__"},
						{"id": 2, "__tags_LEVEL_1": "electronics"},
						{"id": 3, "__tags_LEVEL_1": "apple"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "tags", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						{Name: "__tags_LEVEL_1", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.flattenTags(tt.input, tt.tagColumn)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddRefToDataResults(t *testing.T) {
	service := &widgetsService{}

	ref := "test-ref"
	tests := []struct {
		name                     string
		input                    *datasetmodels.DatasetData
		ref                      *string
		want                     datasetmodels.DatasetData
		populateEmptyRowsWithRef bool
	}{
		{
			name: "with nil ref",
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "value": "test1"},
						{"id": 2, "value": "test2"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "value", DatabaseType: "string"},
					},
				},
			},
			ref:                      nil,
			populateEmptyRowsWithRef: true,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "value": "test1"},
						{"id": 2, "value": "test2"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "value", DatabaseType: "string"},
					},
				},
			},
		},
		{
			name:                     "with valid ref",
			populateEmptyRowsWithRef: true,
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "value": "test1"},
						{"id": 2, "value": "test2"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "value", DatabaseType: "string"},
					},
				},
			},
			ref: &ref,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "value": "test1", "__REF": "test-ref"},
						{"id": 2, "value": "test2", "__REF": "test-ref"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "value", DatabaseType: "string"},
						{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
		{
			name:                     "with empty rows",
			populateEmptyRowsWithRef: true,
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "value", DatabaseType: "string"},
					},
				},
			},
			ref: &ref,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{
							"__REF": "test-ref",
						},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "value", DatabaseType: "string"},
						{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
		{
			name:                     "with empty rows and no ref",
			populateEmptyRowsWithRef: true,
			input:                    &datasetmodels.DatasetData{},
			ref:                      &ref,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
					Rows: []map[string]interface{}{
						{
							"__REF": "test-ref",
						},
					},
				},
			},
		},
		{
			name:                     "with empty rows and no ref",
			populateEmptyRowsWithRef: false,
			input:                    &datasetmodels.DatasetData{},
			ref:                      &ref,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
		{
			name:                     "with empty rows and populateEmptyRowsWithRef=false",
			populateEmptyRowsWithRef: false,
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
					},
				},
			},
			ref: &ref,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
		{
			name:                     "with multiple rows",
			populateEmptyRowsWithRef: true,
			input: &datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "name": "Item 1"},
						{"id": 2, "name": "Item 2"},
						{"id": 3, "name": "Item 3"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "name", DatabaseType: "string"},
					},
				},
			},
			ref: &ref,
			want: datasetmodels.DatasetData{
				QueryResult: dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"id": 1, "name": "Item 1", "__REF": "test-ref"},
						{"id": 2, "name": "Item 2", "__REF": "test-ref"},
						{"id": 3, "name": "Item 3", "__REF": "test-ref"},
					},
					Columns: []dataplatformmodels.ColumnMetadata{
						{Name: "id", DatabaseType: "integer"},
						{Name: "name", DatabaseType: "string"},
						{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.addRefToDataResults(tt.input, tt.ref, tt.populateEmptyRowsWithRef)
			assert.Equal(t, tt.want, got)
		})
	}
}
