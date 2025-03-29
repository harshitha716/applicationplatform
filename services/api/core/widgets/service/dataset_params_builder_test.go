package widgets

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	datasetconstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetconstants "github.com/Zampfi/application-platform/services/api/core/widgets/constants"
	"github.com/Zampfi/application-platform/services/api/core/widgets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	querybuilderconstants "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
)

func stringPtr(s string) *string {
	return &s
}

func TestNewDatasetParamsBuilder(t *testing.T) {
	tests := []struct {
		name          string
		widgetType    string
		expectError   bool
		expectedError string
	}{
		{
			name:        "valid bar chart type",
			widgetType:  "bar_chart",
			expectError: false,
		},
		{
			name:          "invalid widget type",
			widgetType:    "invalid_type",
			expectError:   true,
			expectedError: "unsupported widget type: invalid_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder, err := NewDatasetParamsBuilder(tt.widgetType)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, builder)
			}
		})
	}
}

func TestBasicChartStrategy_ToDatasetParams(t *testing.T) {
	instanceID := uuid.New()
	tests := []struct {
		name        string
		instance    *widgetmodels.WidgetInstance
		filters     map[string]widgetmodels.WidgetFilters
		currency    *string
		timeColumns map[string]string
		periodicity *string
		want        map[string]widgetmodels.GetDataByDatasetIDParams
	}{
		{
			name: "basic chart with x and y axis and default filters",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,

				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.XAxisField: {{Column: "category"}},
								widgetconstants.YAxisField: {{Column: "sales", Aggregation: "sum"}},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions:      []datasetmodels.Filter{{Column: "status", Operator: "eq", Value: "active"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []datasetmodels.Filter{{Column: "region", Operator: "eq", Value: "EU"}},
				}},
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "status", Operator: "eq", Value: "active"},
								{Column: "region", Operator: "eq", Value: "EU"},
							},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
						},
					},
				},
			},
		},
		{
			name: "basic chart with expression in x-axis",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.XAxisField: {{Column: "date", Expression: "DATE_TRUNC('month', date)"}},
								widgetconstants.YAxisField: {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {},
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "DATE_TRUNC('month', date)", Alias: stringPtr("date")},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
						},
					},
				},
			},
		},
		{
			name: "basic chart with group by fields",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.XAxisField:   {{Column: "category"}},
								widgetconstants.YAxisField:   {{Column: "sales", Aggregation: "sum"}},
								widgetconstants.GroupByField: {{Column: "region"}, {Column: "product"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {},
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
							{Column: "region", Alias: stringPtr("region")},
							{Column: "product", Alias: stringPtr("product")},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
							{Column: "region", Order: "ASC", Alias: stringPtr("region")},
							{Column: "product", Order: "ASC", Alias: stringPtr("product")},
						},
					},
				},
			},
		},
		{
			name: "basic chart with time column",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.XAxisField: {{Column: "date"}},
								widgetconstants.YAxisField: {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			timeColumns: map[string]string{
				"dataset1": "date",
			},
			periodicity: func() *string { s := "month"; return &s }(),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "date_trunc('month', date)", Alias: stringPtr("date")},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
						},
					},
				},
			},
		},
	}

	strategy := BasicChartStrategy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strategy.ToDatasetParams(tt.instance, widgetmodels.DatasetBuilderParams{
				Filters:     tt.filters,
				Currency:    tt.currency,
				TimeColumns: tt.timeColumns,
				Periodicity: tt.periodicity,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPieChartStrategy_ToDatasetParams(t *testing.T) {
	instanceID := uuid.New()
	tests := []struct {
		name        string
		instance    *widgetmodels.WidgetInstance
		filters     map[string]widgetmodels.WidgetFilters
		currency    *string
		timeColumns map[string]string
		periodicity *string
		want        map[string]widgetmodels.GetDataByDatasetIDParams
	}{
		{
			name: "pie chart with dimension and measure, and default filters",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,

				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.SlicesField: {{Column: "category"}},
								widgetconstants.ValuesField: {{Column: "sales", Aggregation: "sum"}},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "status", Operator: "eq", Value: "active"},
									{Column: "region", Operator: "eq", Value: "EU"},
								},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []datasetmodels.Filter{{Column: "region", Operator: "eq", Value: "EU"}},
				}},
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "status", Operator: "eq", Value: "active"},
								{Column: "region", Operator: "eq", Value: "EU"},
								{Column: "region", Operator: "eq", Value: "EU"},
							},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
						},
					},
				},
			},
		},
		{
			name: "pie chart with expression in dimension",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							Ref:       "ref1",
							DatasetID: "dataset1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.SlicesField: {{Column: "product"}},
								widgetconstants.ValuesField: {{Column: "revenue", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {},
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "revenue", Function: "sum", Alias: "revenue"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "product", Alias: stringPtr("product")},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "product", Order: "ASC", Alias: stringPtr("product")},
						},
					},
				},
			},
		},
		{
			name: "pie chart with group by fields",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.SlicesField: {{Column: "product"}},
								widgetconstants.ValuesField: {{Column: "revenue", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {},
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "revenue", Function: "sum", Alias: "revenue"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "product", Alias: stringPtr("product")},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "product", Order: "ASC", Alias: stringPtr("product")},
						},
					},
				},
			},
		},
	}

	strategy := PieChartStrategy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strategy.ToDatasetParams(tt.instance, widgetmodels.DatasetBuilderParams{
				Filters:     tt.filters,
				Currency:    tt.currency,
				TimeColumns: tt.timeColumns,
				Periodicity: tt.periodicity,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPivotTableStrategy_ToDatasetParams(t *testing.T) {
	instanceID := uuid.New()
	tests := []struct {
		name        string
		instance    *widgetmodels.WidgetInstance
		filters     map[string]widgetmodels.WidgetFilters
		currency    *string
		want        map[string]widgetmodels.GetDataByDatasetIDParams
		timeColumns map[string]string
		periodicity *string
	}{
		{
			name: "pivot table with rows, columns and values",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.RowsField:    {{Column: "category"}},
								widgetconstants.ColumnsField: {{Column: "region"}},
								widgetconstants.ValuesField:  {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {},
			},
			currency: nil,
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
							{Column: "region", Alias: stringPtr("region")},
						},
						Filters:    datasetmodels.FilterModel{},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
							{Column: "region", Order: "ASC", Alias: stringPtr("region")},
						},
					},
				},
			},
		},
		{
			name: "pivot table with rows, columns, values and custom currency",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.RowsField:    {{Column: "category"}},
								widgetconstants.ColumnsField: {{Column: "region"}},
								widgetconstants.ValuesField:  {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {},
			},
			currency: stringPtr("JPY"),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
							{Column: "region", Alias: stringPtr("region")},
						},
						Filters:    datasetmodels.FilterModel{},
						FxCurrency: stringPtr("JPY"),
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
							{Column: "region", Order: "ASC", Alias: stringPtr("region")},
						},
					},
				},
			},
		},
		{
			name: "Complex pivot table with multiple mappings and sort by",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							Ref:       "Opening Balance",
							DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{
									Column:      "balance_value",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "first",
									Alias:       "value",
								}},
								widgetconstants.RowsField: {
									{
										Column:                  "account_type",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
									{
										Column:                  "account_number",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
								},
								widgetconstants.ColumnsField: {{
									Column:                  "time_stamp_local",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "time_stamp_local", Order: "ASC", Alias: "date"},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "opening"},
								},
							},
						},
						{
							Ref:       "Entity Cashflow",
							DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.RowsField: {
									{
										Column:                  "entity_name",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
									{
										Column:                  "account_number",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
								},
								widgetconstants.ValuesField: {{
									Column:      "effective_amount",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "sum",
									Alias:       "value",
								}},
								widgetconstants.ColumnsField: {{
									Column:                  "posted_time_stamp_local",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "posted_time_stamp_local", Order: "ASC", Alias: "date"},
							},
						},
						{
							Ref:       "Closing Balance",
							DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{
									Column:      "balance_value",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "first",
									Alias:       "value",
								}},
								widgetconstants.RowsField: {
									{
										Column:                  "account_type",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
									{
										Column:                  "account_number",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
								},
								widgetconstants.ColumnsField: {{
									Column:                  "time_stamp_local",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "time_stamp_local", Order: "DESC", Alias: "date"},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "closing"},
								},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"60539e25-c3df-4084-9ea2-95036c2612c3": widgetmodels.WidgetFilters{
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
						},
					},
				},
				"eb516bd6-4e8d-46c9-855f-8a57657072a2": widgetmodels.WidgetFilters{
					DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "posted_time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
						},
					},
				},
			},
			timeColumns: map[string]string{
				"60539e25-c3df-4084-9ea2-95036c2612c3": "time_stamp_local",
				"eb516bd6-4e8d-46c9-855f-8a57657072a2": "posted_time_stamp_local",
			},
			periodicity: func() *string { s := "month"; return &s }(),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"Opening Balance": {
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "account_type", Alias: stringPtr("account_type")},
							{Column: "account_number", Alias: stringPtr("account_number")},
							{Column: "date_trunc('month', time_stamp_local)", Alias: stringPtr("date")},
						},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "balance_value", Function: "sum", Alias: "value"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
							{Column: "account_type", Order: "ASC", Alias: stringPtr("account_type")},
							{Column: "account_number", Order: "ASC", Alias: stringPtr("account_number")},
						},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: nil,
							Windows: []datasetmodels.WindowConfig{
								{
									Function: "ROW_NUMBER()",

									PartitionBy: []datasetmodels.ColumnConfig{
										{Column: "account_type"},
										{Column: "account_number"},
										{Column: "date_trunc('month', time_stamp_local)"},
									},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "time_stamp_local", Order: "ASC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "opening"},
									{Column: "time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false},
								},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "balance_value"},
								{Column: "account_type"},
								{Column: "account_number"},
								{Column: "time_stamp_local"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
					},
				},
				"Entity Cashflow": {
					DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "effective_amount", Function: "sum", Alias: "value"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "entity_name", Alias: stringPtr("entity_name")},
							{Column: "account_number", Alias: stringPtr("account_number")},
							{Column: "date_trunc('month', posted_time_stamp_local)", Alias: stringPtr("date")},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "posted_time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
							},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
							{Column: "entity_name", Order: "ASC", Alias: stringPtr("entity_name")},
							{Column: "account_number", Order: "ASC", Alias: stringPtr("account_number")},
						},
					},
				},
				"Closing Balance": {
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "account_type", Alias: stringPtr("account_type")},
							{Column: "account_number", Alias: stringPtr("account_number")},
							{Column: "date_trunc('month', time_stamp_local)", Alias: stringPtr("date")},
						},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "balance_value", Function: "sum", Alias: "value"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "DESC", Alias: stringPtr("date")},
							{Column: "account_type", Order: "ASC", Alias: stringPtr("account_type")},
							{Column: "account_number", Order: "ASC", Alias: stringPtr("account_number")},
						},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: nil,
							Windows: []datasetmodels.WindowConfig{
								{
									Function: "ROW_NUMBER()",
									PartitionBy: []datasetmodels.ColumnConfig{
										{Column: "account_type"},
										{Column: "account_number"},
										{Column: "date_trunc('month', time_stamp_local)"},
									},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "time_stamp_local", Order: "DESC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "closing"},
									{Column: "time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false},
								},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "balance_value"},
								{Column: "account_type"},
								{Column: "account_number"},
								{Column: "time_stamp_local"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
					},
				},
			},
		},
		{
			name: "Cashflow entity with sort by",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							Ref:       "Opening Balance",
							DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{
									Column:      "balance_value",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "first",
									Alias:       "value",
									SortBy: []widgetmodels.SortBy{
										{Column: "time_stamp_local", Order: "ASC"},
										{Column: "account_type", Order: "ASC"},
									},
								}},
								widgetconstants.RowsField: {
									{
										Column:                  "account_type",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
									{
										Column:                  "account_number",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
								},
								widgetconstants.ColumnsField: {{
									Column:                  "time_stamp_local",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "time_stamp_local", Order: "ASC", Alias: "date"},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "opening"},
								},
							},
						},
						{
							Ref:       "Entity Cashflow",
							DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.RowsField: {
									{
										Column:                  "entity_name",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
									{
										Column:                  "account_number",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
								},
								widgetconstants.ValuesField: {{
									Column:      "effective_amount",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "sum",
									Alias:       "value",
								}},
								widgetconstants.ColumnsField: {{
									Column:                  "posted_time_stamp_local",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "posted_time_stamp_local", Order: "ASC", Alias: "date"},
							},
						},
						{
							Ref:       "Closing Balance",
							DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{
									Column:      "balance_value",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "first",
									Alias:       "value",
									SortBy: []widgetmodels.SortBy{
										{Column: "time_stamp_local", Order: "DESC"},
										{Column: "account_type", Order: "DESC"},
									},
								}},
								widgetconstants.RowsField: {
									{
										Column:                  "account_type",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
									{
										Column:                  "account_number",
										Type:                    "string",
										FieldType:               "dimension",
										DrilldownFilterType:     "multi-select",
										DrilldownFilterOperator: "in",
									},
								},
								widgetconstants.ColumnsField: {{
									Column:                  "time_stamp_local",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "time_stamp_local", Order: "DESC", Alias: "date"},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "closing"},
								},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"60539e25-c3df-4084-9ea2-95036c2612c3": widgetmodels.WidgetFilters{
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
						},
					},
				},
				"eb516bd6-4e8d-46c9-855f-8a57657072a2": widgetmodels.WidgetFilters{
					DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "posted_time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
						},
					},
				},
			},
			timeColumns: map[string]string{
				"60539e25-c3df-4084-9ea2-95036c2612c3": "time_stamp_local",
				"eb516bd6-4e8d-46c9-855f-8a57657072a2": "posted_time_stamp_local",
			},
			periodicity: func() *string { s := "month"; return &s }(),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"Opening Balance": {
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "account_type", Alias: stringPtr("account_type")},
							{Column: "account_number", Alias: stringPtr("account_number")},
							{Column: "date_trunc('month', time_stamp_local)", Alias: stringPtr("date")},
						},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "balance_value", Function: "sum", Alias: "value"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
							{Column: "account_type", Order: "ASC", Alias: stringPtr("account_type")},
							{Column: "account_number", Order: "ASC", Alias: stringPtr("account_number")},
						},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: nil,
							Windows: []datasetmodels.WindowConfig{
								{
									Function: "ROW_NUMBER()",

									PartitionBy: []datasetmodels.ColumnConfig{
										{Column: "account_type"},
										{Column: "account_number"},
										{Column: "date_trunc('month', time_stamp_local)"},
									},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "time_stamp_local", Order: "ASC"},
										{Column: "account_type", Order: "ASC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "opening"},
									{Column: "time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false},
								},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "balance_value"},
								{Column: "account_type"},
								{Column: "account_number"},
								{Column: "time_stamp_local"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
					},
				},
				"Entity Cashflow": {
					DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "effective_amount", Function: "sum", Alias: "value"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "entity_name", Alias: stringPtr("entity_name")},
							{Column: "account_number", Alias: stringPtr("account_number")},
							{Column: "date_trunc('month', posted_time_stamp_local)", Alias: stringPtr("date")},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "posted_time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
							},
						},
						FxCurrency: nil,
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
							{Column: "entity_name", Order: "ASC", Alias: stringPtr("entity_name")},
							{Column: "account_number", Order: "ASC", Alias: stringPtr("account_number")},
						},
					},
				},
				"Closing Balance": {
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "account_type", Alias: stringPtr("account_type")},
							{Column: "account_number", Alias: stringPtr("account_number")},
							{Column: "date_trunc('month', time_stamp_local)", Alias: stringPtr("date")},
						},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "balance_value", Function: "sum", Alias: "value"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "DESC", Alias: stringPtr("date")},
							{Column: "account_type", Order: "ASC", Alias: stringPtr("account_type")},
							{Column: "account_number", Order: "ASC", Alias: stringPtr("account_number")},
						},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: nil,
							Windows: []datasetmodels.WindowConfig{
								{
									Function: "ROW_NUMBER()",
									PartitionBy: []datasetmodels.ColumnConfig{
										{Column: "account_type"},
										{Column: "account_number"},
										{Column: "date_trunc('month', time_stamp_local)"},
									},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "time_stamp_local", Order: "DESC"},
										{Column: "account_type", Order: "DESC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "closing"},
									{Column: "time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false},
								},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "balance_value"},
								{Column: "account_type"},
								{Column: "account_number"},
								{Column: "time_stamp_local"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
					},
				},
			},
		},
		{
			name: "Forecasting pivot with sort by",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							Ref:       "Forecasting opening balance",
							DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{
									Column:      "_previous_closing_balance",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "first",
									Alias:       "value",
									SortBy: []widgetmodels.SortBy{
										{Column: "date", Order: "ASC"},
										{Column: "tags", Order: "ASC"},
									},
								}},
								widgetconstants.RowsField: {},
								widgetconstants.ColumnsField: {{
									Column:                  "date",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
						},
						// {
						// 	Ref:       "Forecasting entity cashflow",
						// 	DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
						// 	Fields: map[string][]widgetmodels.Field{
						// 		widgetconstants.RowsField: {
						// 			{
						// 				Column:                  "tags",
						// 				Type:                    "tag",
						// 				FieldType:               "measure",
						// 				DrilldownFilterType:     "search",
						// 				DrilldownFilterOperator: "startswith",
						// 			},
						// 		},
						// 		widgetconstants.ValuesField: {{
						// 			Column:      "amount",
						// 			Type:        "number",
						// 			FieldType:   "measure",
						// 			Aggregation: "sum",
						// 			Alias:       "value",
						// 		}},
						// 		widgetconstants.ColumnsField: {{
						// 			Column:                  "date",
						// 			Type:                    "timestamp",
						// 			Alias:                   "date",
						// 			FieldType:               "dimension",
						// 			DrilldownFilterType:     "date-range",
						// 			DrilldownFilterOperator: "inbetween",
						// 		}},
						// 	},
						// 	SortBy: []widgetmodels.SortBy{
						// 		{Column: "date", Order: "ASC", Alias: "date"},
						// 	},
						// },
						{
							Ref:       "Forecasting closing balance",
							DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{
									Column:      "_closing_balance",
									Type:        "number",
									FieldType:   "measure",
									Aggregation: "first",
									Alias:       "value",
									SortBy: []widgetmodels.SortBy{
										{Column: "date", Order: "ASC"},
										{Column: "tags", Order: "ASC"},
									},
								}},
								widgetconstants.RowsField: {},
								widgetconstants.ColumnsField: {{
									Column:                  "date",
									Type:                    "timestamp",
									Alias:                   "date",
									FieldType:               "dimension",
									DrilldownFilterType:     "date-range",
									DrilldownFilterOperator: "inbetween",
								}},
							},
							SortBy: []widgetmodels.SortBy{
								{Column: "date", Order: "ASC", Alias: "date"},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"60539e25-c3df-4084-9ea2-95036c2612c3": widgetmodels.WidgetFilters{
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "date", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
						},
					},
				},
				"eb516bd6-4e8d-46c9-855f-8a57657072a2": widgetmodels.WidgetFilters{
					DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "date", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
						},
					},
				},
			},
			timeColumns: map[string]string{
				"60539e25-c3df-4084-9ea2-95036c2612c3": "date",
				"eb516bd6-4e8d-46c9-855f-8a57657072a2": "date",
			},
			periodicity: func() *string { s := "month"; return &s }(),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"Forecasting opening balance": {
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "date_trunc('month', date)", Alias: stringPtr("date")},
						},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "_previous_closing_balance", Function: "sum", Alias: "value"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
						},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: nil,
							Windows: []datasetmodels.WindowConfig{
								{
									Function: "ROW_NUMBER()",

									PartitionBy: []datasetmodels.ColumnConfig{
										{Column: "date_trunc('month', date)"},
									},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "date", Order: "ASC"},
										{Column: "tags", Order: "ASC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "date", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false},
								},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "_previous_closing_balance"},
								{Column: "date"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
					},
				},
				// "Entity Cashflow": {
				// 	DatasetID: "eb516bd6-4e8d-46c9-855f-8a57657072a2",
				// 	Params: datasetmodels.DatasetParams{
				// 		Columns: []datasetmodels.ColumnConfig{},
				// 		Aggregations: []datasetmodels.Aggregation{
				// 			{Column: "effective_amount", Function: "sum", Alias: "value"},
				// 		},
				// 		GroupBy: []datasetmodels.GroupBy{
				// 			{Column: "entity_name", Alias: stringPtr("entity_name")},
				// 			{Column: "account_number", Alias: stringPtr("account_number")},
				// 			{Column: "date_trunc('month', posted_time_stamp_local)", Alias: stringPtr("date")},
				// 		},
				// 		Filters: datasetmodels.FilterModel{
				// 			LogicalOperator: "AND",
				// 			Conditions: []datasetmodels.Filter{
				// 				{Column: "posted_time_stamp_local", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
				// 			},
				// 		},
				// 		FxCurrency: nil,
				// 		OrderBy: []datasetmodels.OrderBy{
				// 			{Column: "date", Order: "ASC"},
				// 			{Column: "entity_name", Order: "ASC"},
				// 			{Column: "account_number", Order: "ASC"},
				// 		},
				// 	},
				// },
				"Forecasting closing balance": {
					DatasetID: "60539e25-c3df-4084-9ea2-95036c2612c3",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "date_trunc('month', date)", Alias: stringPtr("date")},
						},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "_closing_balance", Function: "sum", Alias: "value"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
						},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: nil,
							Windows: []datasetmodels.WindowConfig{
								{
									Function: "ROW_NUMBER()",
									PartitionBy: []datasetmodels.ColumnConfig{
										{Column: "date_trunc('month', date)"},
									},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "date", Order: "ASC"},
										{Column: "tags", Order: "ASC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "date", Operator: "inbetween", Value: []string{"2024-01-01", "2024-01-31"}},
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false},
								},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "_closing_balance"},
								{Column: "date"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
					},
				},
			},
		},
	}

	strategy := PivotTableStrategy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strategy.ToDatasetParams(tt.instance, widgetmodels.DatasetBuilderParams{
				Filters:     tt.filters,
				Currency:    tt.currency,
				TimeColumns: tt.timeColumns,
				Periodicity: tt.periodicity,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKPIStrategy_ToDatasetParams(t *testing.T) {
	instanceID := uuid.New()
	tests := []struct {
		name     string
		instance *widgetmodels.WidgetInstance
		filters  map[string]widgetmodels.WidgetFilters
		currency *string
		want     map[string]widgetmodels.GetDataByDatasetIDParams
	}{
		{
			name: "kpi with primary value and filter",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.PrimaryValueField: {{Column: "revenue", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []datasetmodels.Filter{{Column: "year", Operator: "eq", Value: "2024"}},
				}},
			},
			currency: nil,
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "revenue", Function: "sum", Alias: "revenue"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions:      []datasetmodels.Filter{{Column: "year", Operator: "eq", Value: "2024"}},
						},
						FxCurrency: nil,
					},
				},
			},
		},
		{
			name: "kpi with primary value, filter and custom currency",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.PrimaryValueField: {{Column: "revenue", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []datasetmodels.Filter{{Column: "year", Operator: "eq", Value: "2024"}},
				}},
			},
			currency: stringPtr("AUD"),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "revenue", Function: "sum", Alias: "revenue"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions:      []datasetmodels.Filter{{Column: "year", Operator: "eq", Value: "2024"}},
						},
						FxCurrency: stringPtr("AUD"),
					},
				},
			},
		},
		{
			name: "kpi with a `first` aggregation, default filter, and currency ",
			instance: &widgetmodels.WidgetInstance{
				ID: instanceID,
				DataMappings: widgetmodels.DataMappings{
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.PrimaryValueField: {{Column: "revenue", Aggregation: "first", Alias: "revenue", SortBy: []widgetmodels.SortBy{{Column: "date", Order: "ASC"}}}},
							},
							DefaultFilters: &datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions: []datasetmodels.Filter{
									{Column: "balance_type", Operator: "eq", Value: "asset"},
								},
							},
						},
					},
				},
			},
			filters: map[string]widgetmodels.WidgetFilters{
				"dataset1": {Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []datasetmodels.Filter{{Column: "year", Operator: "eq", Value: "2024"}},
				}},
			},
			currency: stringPtr("AUD"),
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "revenue", Function: "sum", Alias: "revenue"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Subquery: &datasetmodels.DatasetParams{
							FxCurrency: stringPtr("AUD"),
							Windows: []datasetmodels.WindowConfig{
								{
									Function:    "ROW_NUMBER()",
									PartitionBy: []datasetmodels.ColumnConfig{},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "date", Order: "ASC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "AND",
								Conditions:      []datasetmodels.Filter{{Column: "balance_type", Operator: "eq", Value: "asset"}, {Column: "year", Operator: "eq", Value: "2024"}, {Column: datasetconstants.ZampIsDeletedColumn, Operator: querybuilderconstants.EqualOperator, Value: false}},
							},
							Columns: []datasetmodels.ColumnConfig{
								{Column: "revenue"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
						},

						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions:      []datasetmodels.Filter{{Column: "rn", Operator: "eq", Value: 1}},
						},
					},
				},
			},
		},
	}

	strategy := KPIStrategy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strategy.ToDatasetParams(tt.instance, widgetmodels.DatasetBuilderParams{
				Filters:  tt.filters,
				Currency: tt.currency,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeFilters(t *testing.T) {
	andOp := datasetmodels.LogicalOperator(querybuilderconstants.LogicalOperatorAnd)

	tests := []struct {
		name           string
		defaultFilters datasetmodels.FilterModel
		sheetFilters   datasetmodels.FilterModel
		want           datasetmodels.FilterModel
	}{
		// Test 1: Both filters empty
		{
			name:           "Empty filters",
			defaultFilters: datasetmodels.FilterModel{},
			sheetFilters:   datasetmodels.FilterModel{},
			want:           datasetmodels.FilterModel{},
		},
		// Test 2: Default union conflicting sheet filter
		{
			name: "Default union sheet filter",
			defaultFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
				},
			},
			sheetFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "Canada"},
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
			want: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
					{Column: "country", Operator: "=", Value: "Canada"},
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
		},
		// Test 3: Non-conflicting filters merged
		{
			name: "Non-conflicting filters merged",
			defaultFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
				},
			},
			sheetFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
			want: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
		},
		{
			name: "Sheet and Default Union",
			defaultFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
			sheetFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "Canada"},
					{Column: "year", Operator: "=", Value: 2023},
				},
			},
			want: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
					{Column: "year", Operator: "=", Value: 2024},
					{Column: "country", Operator: "=", Value: "Canada"},
					{Column: "year", Operator: "=", Value: 2023},
				},
			},
		},
		// Test 5: Default empty, return sheet
		{
			name:           "Default empty, return sheet",
			defaultFilters: datasetmodels.FilterModel{},
			sheetFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
			want: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "year", Operator: "=", Value: 2024},
				},
			},
		},
		// Test 6: Sheet empty, return default
		{
			name: "Sheet empty, return default",
			defaultFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
				},
			},
			sheetFilters: datasetmodels.FilterModel{},
			want: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{Column: "country", Operator: "=", Value: "USA"},
				},
			},
		},
		{
			name: "Nested conditions with different operators",
			defaultFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{
						Column:   "recon_status",
						Operator: "neq",
						Value:    []string{"settled"},
					},
				},
			},
			sheetFilters: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{
						Column:   "pg_record_date",
						Operator: "inbetween",
						Value:    []string{"2024-01-01", "2024-01-31"},
					},
					{
						Column:   "country_code",
						Operator: "in",
						Value:    []string{"US", "CA"},
					},
				},
			},
			want: datasetmodels.FilterModel{
				LogicalOperator: andOp,
				Conditions: []datasetmodels.Filter{
					{
						Column:   "recon_status",
						Operator: "neq",
						Value:    []string{"settled"},
					},
					{
						Column:   "pg_record_date",
						Operator: "inbetween",
						Value:    []string{"2024-01-01", "2024-01-31"},
					},
					{
						Column:   "country_code",
						Operator: "in",
						Value:    []string{"US", "CA"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.MergeFilters(&tt.defaultFilters, &tt.sheetFilters)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Test '%s' failed:\n%s", tt.name, diff)
			}
		})
	}
}

func TestAddTimeColumn(t *testing.T) {
	tests := []struct {
		name        string
		params      *datasetmodels.DatasetParams
		timeColumns map[string]string
		periodicity *string
		want        datasetmodels.DatasetParams
	}{
		{
			name: "with time column and periodicity",
			params: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "created_at"},
					{Column: "sales"},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "created_at"},
				},
				OrderBy: []datasetmodels.OrderBy{
					{Column: "created_at", Alias: stringPtr("created_at")},
				},
			},
			timeColumns: map[string]string{
				"dataset1": "created_at",
			},
			periodicity: func() *string { s := "month"; return &s }(),
			want: datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "date_trunc('month', created_at)"},
					{Column: "sales"},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "date_trunc('month', created_at)"},
				},
				OrderBy: []datasetmodels.OrderBy{
					{Column: "date_trunc('month', created_at)", Alias: stringPtr("created_at")},
				},
			},
		},
		{
			name: "with nil time column",
			params: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "created_at"},
					{Column: "sales"},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "created_at"},
				},
			},
			timeColumns: nil,
			periodicity: func() *string { s := "month"; return &s }(),
			want: datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "created_at"},
					{Column: "sales"},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "created_at"},
				},
			},
		},
		{
			name: "with nil periodicity",
			params: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "created_at"},
					{Column: "sales"},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "created_at"},
				},
			},
			timeColumns: map[string]string{
				"dataset1": "created_at",
			},
			periodicity: nil,
			want: datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "created_at"},
					{Column: "sales"},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "created_at"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			baseStrategy.AddTimeColumn(tt.params, tt.timeColumns, tt.periodicity, "dataset1")
			assert.Equal(t, tt.want, *tt.params)
		})
	}
}

func TestAddSortBy(t *testing.T) {
	tests := []struct {
		name   string
		params *datasetmodels.DatasetParams
		sortBy []widgetmodels.SortBy
		want   []datasetmodels.OrderBy
	}{
		{
			name: "empty sort by",
			params: &datasetmodels.DatasetParams{
				OrderBy: make([]datasetmodels.OrderBy, 0),
			},
			sortBy: []widgetmodels.SortBy{},
			want:   []datasetmodels.OrderBy{},
		},
		{
			name:   "single sort by",
			params: &datasetmodels.DatasetParams{},
			sortBy: []widgetmodels.SortBy{
				{Column: "revenue", Order: "DESC"},
			},
			want: []datasetmodels.OrderBy{
				{Column: "revenue", Order: "DESC", Alias: stringPtr("revenue")},
			},
		},
		{
			name: "multiple sort by with existing column",
			params: &datasetmodels.DatasetParams{
				OrderBy: []datasetmodels.OrderBy{
					{Column: "existing_col", Order: "ASC", Alias: stringPtr("existing_col")},
				},
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "revenue", Order: "DESC"},
				{Column: "date", Order: "ASC"},
			},
			want: []datasetmodels.OrderBy{
				{Column: "existing_col", Order: "ASC", Alias: stringPtr("existing_col")},
				{Column: "revenue", Order: "DESC", Alias: stringPtr("revenue")},
				{Column: "date", Order: "ASC", Alias: stringPtr("date")},
			},
		},
		{
			name: "with group by columns",
			params: &datasetmodels.DatasetParams{
				GroupBy: []datasetmodels.GroupBy{
					{Column: "category", Alias: ptr("category_alias")},
				},
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "revenue", Order: "DESC"},
			},
			want: []datasetmodels.OrderBy{
				{Column: "revenue", Order: "DESC", Alias: stringPtr("revenue")},
				{Column: "category_alias", Order: "ASC", Alias: stringPtr("category_alias")},
			},
		},
		{
			name: "with aggregation columns",
			params: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{
					{Column: "sales", Function: "sum", Alias: "total_sales"},
				},
			},
			sortBy: []widgetmodels.SortBy{},
			want:   nil,
		},
		{
			name: "with duplicate columns",
			params: &datasetmodels.DatasetParams{
				GroupBy: []datasetmodels.GroupBy{
					{Column: "category", Alias: ptr("category_alias")},
				},
				Aggregations: []datasetmodels.Aggregation{
					{Column: "category", Function: "count", Alias: "category_alias"},
				},
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "category_alias", Order: "DESC"},
			},
			want: []datasetmodels.OrderBy{
				{Column: "category_alias", Order: "DESC", Alias: stringPtr("category_alias")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			baseStrategy.AddSortBy(tt.params, tt.sortBy)
			assert.Equal(t, tt.want, tt.params.OrderBy)
		})
	}
}

func TestIsParametrizedString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		// Basic parametrized strings
		{name: "simple parameter", value: "{{.$today}}", want: true},
		{name: "parameter with addDays method", value: "{{.$today.addDays(1)}}", want: true},
		{name: "parameter with addSeconds method", value: "{{.$today.addSeconds(1)}}", want: true},

		// Complex parametrized strings
		{name: "complex parameter with multiple methods", value: "{{.$start_date.addDays(5).addSeconds(30)}}", want: true},
		{name: "parameter with nested parentheses", value: "{{.$today.addDays(1 + (2 * 3))}}", want: true},
		{name: "parameter with special characters", value: "{{.$filter_value-123.addDays(1)}}", want: true},

		// Invalid formats
		{name: "missing $ prefix", value: "{{.today.addDays(1)}}", want: false},
		{name: "unclosed braces", value: "{{.today.addDays(1)}", want: false},
		{name: "missing opening braces", value: ".$today.addDays(1)}}", want: false},
		{name: "completely invalid format", value: "{today}", want: false},
		{name: "non-parametrized string", value: "regular string", want: false},
		{name: "empty string", value: "", want: false},
		{name: "nil value", value: "", want: false},
		{name: "with whitespace", value: "{{ .$today }}", want: false},
		{name: "with multiple parameters", value: "{{.$today}} and {{.$end_date}}", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()

			// First call - should initialize regex
			got := baseStrategy.IsParametrizedString(tt.value)
			assert.Equal(t, tt.want, got)

			// Second call - should use already initialized regex
			got = baseStrategy.IsParametrizedString(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParametrizeDefaultFilters(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tests := []struct {
		name                 string
		defaultFilters       *datasetmodels.FilterModel
		datasetBuilderParams widgetmodels.DatasetBuilderParams
		datasetId            string
		want                 *datasetmodels.FilterModel
	}{
		{
			name: "basic parametrized filter with today",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{"{{.$today}}"}},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters:     map[string]widgetmodels.WidgetFilters{},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{startOfDay.Format(time.DateTime)}},
				},
			},
		},
		{
			name: "parametrized filter with today.addDays(1)",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{"{{.$today.addDays(1)}}"}},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters:     map[string]widgetmodels.WidgetFilters{},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{startOfDay.AddDate(0, 0, 1).Format(time.DateTime)}},
				},
			},
		},
		{
			name: "parametrized filter with start_date from sheet filters",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{"{{.$start_date}}"}},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters: map[string]widgetmodels.WidgetFilters{
					"dataset1": {
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "created_at", Operator: "inbetween", Value: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"}},
							},
						},
					},
				},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{"2024-01-01 00:00:00"}},
				},
			},
		},
		{
			name: "parametrized filter with end_date from sheet filters",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{"{{.$end_date}}"}},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters: map[string]widgetmodels.WidgetFilters{
					"dataset1": {
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "created_at", Operator: "inbetween", Value: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"}},
							},
						},
					},
				},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "created_at", Operator: "=", Value: []string{"2024-01-31 00:00:00"}},
				},
			},
		},
		{
			name: "parametrized filter with multiple conditions",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "start_date", Operator: "=", Value: []string{"{{.$start_date}}"}},
					{Column: "end_date", Operator: "=", Value: []string{"{{.$end_date}}"}},
					{Column: "created_at", Operator: "=", Value: []string{"{{.$today}}"}},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters: map[string]widgetmodels.WidgetFilters{
					"dataset1": {
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "start_date", Operator: "inbetween", Value: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"}},
								{Column: "end_date", Operator: "inbetween", Value: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"}},
								{Column: "created_at", Operator: "inbetween", Value: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"}},
							},
						},
					},
				},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "start_date", Operator: "=", Value: []string{"2024-01-01 00:00:00"}},
					{Column: "end_date", Operator: "=", Value: []string{"2024-01-31 00:00:00"}},
					{Column: "created_at", Operator: "=", Value: []string{startOfDay.Format(time.DateTime)}},
				},
			},
		},
		{
			name: "non-parametrized filter remains unchanged",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "status", Operator: "=", Value: []string{"active"}},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters:     map[string]widgetmodels.WidgetFilters{},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "status", Operator: "=", Value: []string{"active"}},
				},
			},
		},
		{
			name: "test case for operators expecting a float value",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "value", Operator: "gt", Value: 100},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters:     map[string]widgetmodels.WidgetFilters{},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "value", Operator: "gt", Value: 100},
				},
			},
		},
		{
			name: "test case for operators expecting a string value",
			defaultFilters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "value", Operator: "eq", Value: "opening"},
				},
			},
			datasetBuilderParams: widgetmodels.DatasetBuilderParams{
				Filters:     map[string]widgetmodels.WidgetFilters{},
				Currency:    nil,
				TimeColumns: map[string]string{"dataset1": "created_at"},
				Periodicity: stringPtr("month"),
			},
			datasetId: "dataset1",
			want: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "value", Operator: "eq", Value: "opening"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.ParametrizeDefaultFilters(tt.defaultFilters, tt.datasetId, tt.datasetBuilderParams)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessParametrizedFilter(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tests := []struct {
		name                 string
		values               []string
		sheetConditionValues []string
		want                 []string
	}{
		// Basic parametrized values
		{name: "today parameter", values: []string{"{{.$today}}"}, sheetConditionValues: []string{}, want: []string{startOfDay.Format(time.DateTime)}},
		{name: "today with addDays", values: []string{"{{.$today.addDays(10)}}"}, sheetConditionValues: []string{}, want: []string{startOfDay.AddDate(0, 0, 10).Format(time.DateTime)}},
		{name: "end_date with addDays", values: []string{"{{.$end_date.addDays(10)}}"}, sheetConditionValues: []string{"2024-01-01 00:00:00"}, want: []string{time.Date(2024, 1, 1, 0, 0, 0, 0, startOfDay.Location()).AddDate(0, 0, 10).Format(time.DateTime)}},
		{name: "start_date with addDays", values: []string{"{{.$start_date.addDays(10)}}"}, sheetConditionValues: []string{"2024-01-01 00:00:00", "2024-02-10 00:00:00"}, want: []string{time.Date(2024, 1, 1, 0, 0, 0, 0, startOfDay.Location()).AddDate(0, 0, 10).Format(time.DateTime)}},
		{name: "end_date parameter", values: []string{"{{.$end_date}}"}, sheetConditionValues: []string{"2024-01-01 00:00:00", "2024-02-10 00:00:00"}, want: []string{time.Date(2024, 2, 10, 0, 0, 0, 0, startOfDay.Location()).Format(time.DateTime)}},

		// Non-parametrized values
		{name: "non-parametrized value", values: []string{"static-value"}, sheetConditionValues: []string{"2024-01-01 00:00:00"}, want: []string{"static-value"}},
		{name: "empty values array", values: []string{}, sheetConditionValues: []string{"2024-01-01 00:00:00"}, want: []string{}},
		{name: "nil values", values: nil, sheetConditionValues: []string{"2024-01-01 00:00:00"}, want: nil},

		// Mixed parametrized and non-parametrized values
		{name: "mixed values", values: []string{"{{.$today}}", "static-value", "{{.$end_date}}"}, sheetConditionValues: []string{"2024-01-01 00:00:00"},
			want: []string{startOfDay.Format(time.DateTime), "static-value", "2024-01-01 00:00:00"}},

		// Invalid parameter format
		{name: "invalid parameter format", values: []string{"{{.today}}"}, sheetConditionValues: []string{"2024-01-01 00:00:00"}, want: []string{"{{.today}}"}},
		{name: "malformed parameter", values: []string{"{{.$today"}, sheetConditionValues: []string{"2024-01-01 00:00:00"}, want: []string{"{{.$today"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.ProcessParametrizedFilter(tt.values, tt.sheetConditionValues)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPopulateDateParams(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name             string
		value            string
		want             string
		populationValues []string
	}{
		{name: "addDays Positive", value: "{{.$today}}", want: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.DateTime), populationValues: []string{}},
		{name: "addDays Positive Next Day", value: "{{.$today.addDays(1)}}", want: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Format(time.DateTime), populationValues: []string{}},
		{name: "addDays Negative Previous Day", value: "{{.$today.addDays(-1)}}", want: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1).Format(time.DateTime), populationValues: []string{}},
		// {name: "addSeconds Positive", value: "{{.$today.addSeconds(1)}}", want: time.Now().Add(time.Second).Format(time.DateTime)},
		// {name: "addSeconds Negative", value: "{{.$today.addSeconds(-1)}}", want: time.Now().Add(-time.Second).Format(time.DateTime)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.PopulateDateParams(tt.value, tt.populationValues)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetBaseTime(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tests := []struct {
		name             string
		paramName        string
		populationValues []string
		want             time.Time
	}{
		// Today parameter tests
		{name: "today parameter", paramName: widgetconstants.ParameterMethodToday, populationValues: []string{}, want: startOfDay},
		{name: "today parameter with ignored values", paramName: widgetconstants.ParameterMethodToday, populationValues: []string{"2024-01-01 00:00:00"}, want: startOfDay},

		// End day parameter tests
		{name: "end_date with single value", paramName: widgetconstants.ParameterMethodEndDay, populationValues: []string{"2024-01-01 00:00:00"}, want: time.Date(2024, 1, 1, 0, 0, 0, 0, now.Location())},
		{name: "end_date with multiple values", paramName: widgetconstants.ParameterMethodEndDay, populationValues: []string{"2025-01-01 00:00:00", "2024-02-10 00:00:00"}, want: time.Date(2024, 2, 10, 0, 0, 0, 0, now.Location())},
		{name: "end_date with empty population values", paramName: widgetconstants.ParameterMethodEndDay, populationValues: []string{}, want: startOfDay},
		{name: "end_date with nil population values", paramName: widgetconstants.ParameterMethodEndDay, populationValues: nil, want: startOfDay},
		{name: "end_date with invalid date format", paramName: widgetconstants.ParameterMethodEndDay, populationValues: []string{"invalid-date-format"}, want: startOfDay},

		// Start day parameter tests
		{name: "start_date with single value", paramName: widgetconstants.ParameterMethodStartDay, populationValues: []string{"2024-01-01 00:00:00"}, want: time.Date(2024, 1, 1, 0, 0, 0, 0, now.Location())},
		{name: "start_date with multiple values", paramName: widgetconstants.ParameterMethodStartDay, populationValues: []string{"2024-01-01 00:00:00", "2024-02-10 00:00:00"}, want: time.Date(2024, 1, 1, 0, 0, 0, 0, now.Location())},
		{name: "start_date with empty population values", paramName: widgetconstants.ParameterMethodStartDay, populationValues: []string{}, want: startOfDay},
		{name: "start_date with nil population values", paramName: widgetconstants.ParameterMethodStartDay, populationValues: nil, want: startOfDay},
		{name: "start_date with invalid date format", paramName: widgetconstants.ParameterMethodStartDay, populationValues: []string{"invalid-date-format"}, want: startOfDay},

		// Unknown parameter tests
		{name: "unknown parameter name", paramName: "unknown", populationValues: []string{"2024-01-01 00:00:00"}, want: startOfDay},
		{name: "empty parameter name", paramName: "", populationValues: []string{"2024-01-01 00:00:00"}, want: startOfDay},
		{name: "window function as parameter", paramName: widgetconstants.WindowFunctionFirst, populationValues: []string{}, want: startOfDay},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.GetBaseTime(tt.paramName, tt.populationValues)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestApplyMethod(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		baseTime time.Time
		method   string
		args     string
		want     time.Time
	}{
		// addDays tests
		{name: "addDays with positive value", baseTime: now, method: widgetconstants.ParameterMethodAddDays, args: "1", want: now.AddDate(0, 0, 1)},
		{name: "addDays with negative value", baseTime: now, method: widgetconstants.ParameterMethodAddDays, args: "-1", want: now.AddDate(0, 0, -1)},
		{name: "addDays with whitespace", baseTime: now, method: widgetconstants.ParameterMethodAddDays, args: "  5  ", want: now.AddDate(0, 0, 5)},
		{name: "addDays with invalid argument", baseTime: now, method: widgetconstants.ParameterMethodAddDays, args: "invalid", want: now},
		{name: "addDays with empty args", baseTime: now, method: widgetconstants.ParameterMethodAddDays, args: "", want: now},

		// addSeconds tests
		{name: "addSeconds with positive value", baseTime: now, method: widgetconstants.ParameterMethodAddSeconds, args: "60", want: now.Add(60 * time.Second)},
		{name: "addSeconds with negative value", baseTime: now, method: widgetconstants.ParameterMethodAddSeconds, args: "-60", want: now.Add(-60 * time.Second)},
		{name: "addSeconds with whitespace", baseTime: now, method: widgetconstants.ParameterMethodAddSeconds, args: "  30  ", want: now.Add(30 * time.Second)},
		{name: "addSeconds with invalid argument", baseTime: now, method: widgetconstants.ParameterMethodAddSeconds, args: "invalid", want: now},
		{name: "addSeconds with empty args", baseTime: now, method: widgetconstants.ParameterMethodAddSeconds, args: "", want: now},

		// Error cases
		{name: "unknown method", baseTime: now, method: "random", args: "1", want: now},
		{name: "empty method", baseTime: now, method: "", args: "1", want: now},
		{name: "empty method and args", baseTime: now, method: "", args: "", want: now},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.ApplyMethod(tt.baseTime, tt.method, tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper function to create string pointer
func TestParseAndGetValuesFromFilters(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  []string
	}{
		{
			name:  "nil value",
			value: nil,
			want:  nil,
		},
		{
			name:  "string slice",
			value: []string{"value1", "value2"},
			want:  []string{"value1", "value2"},
		},
		{
			name:  "interface slice with strings",
			value: []interface{}{"value1", "value2"},
			want:  []string{"value1", "value2"},
		},
		{
			name:  "interface slice with mixed types",
			value: []interface{}{"value1", 123, true},
			want:  []string{"value1", "", ""},
		},
		{
			name:  "unsupported type",
			value: 123,
			want:  nil,
		},
		{
			name:  "empty slice",
			value: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.ParseAndGetValuesFromFilters(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPopulateParams(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tests := []struct {
		name                 string
		value                string
		sheetConditionValues []string
		want                 string
	}{
		{
			name:                 "today parameter",
			value:                "{{.$today}}",
			sheetConditionValues: []string{},
			want:                 startOfDay.Format(time.DateTime),
		},
		{
			name:                 "today with addDays",
			value:                "{{.$today.addDays(1)}}",
			sheetConditionValues: []string{},
			want:                 startOfDay.AddDate(0, 0, 1).Format(time.DateTime),
		},
		{
			name:                 "end_date parameter",
			value:                "{{.$end_date}}",
			sheetConditionValues: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"},
			want:                 "2024-01-31 00:00:00",
		},
		{
			name:                 "start_date parameter",
			value:                "{{.$start_date}}",
			sheetConditionValues: []string{"2024-01-01 00:00:00", "2024-01-31 00:00:00"},
			want:                 "2024-01-01 00:00:00",
		},
		{
			name:                 "non-parametrized value",
			value:                "static value",
			sheetConditionValues: []string{},
			want:                 "static value",
		},
		{
			name:                 "unsupported parameter",
			value:                "{{.$unsupported}}",
			sheetConditionValues: []string{},
			want:                 "{{.$unsupported}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.PopulateParams(tt.value, tt.sheetConditionValues)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Mock processor for testing
type mockProcessor struct{}

func (p *mockProcessor) CanProcess(paramName string) bool {
	return paramName == "$test_param"
}

func (p *mockProcessor) Process(match string, parts []string, populationValues []string) string {
	return "processed"
}

func TestProcessWithProcessors(t *testing.T) {
	tests := []struct {
		name             string
		value            string
		populationValues []string
		processors       []ParameterProcessor
		want             string
	}{
		{
			name:             "process with mock processor",
			value:            "{{.$test_param}}",
			populationValues: []string{},
			processors:       []ParameterProcessor{&mockProcessor{}},
			want:             "processed",
		},
		{
			name:             "no matching processor",
			value:            "{{.$unknown}}",
			populationValues: []string{},
			processors:       []ParameterProcessor{&mockProcessor{}},
			want:             "{{.$unknown}}",
		},
		{
			name:             "no parameters to process",
			value:            "static value",
			populationValues: []string{},
			processors:       []ParameterProcessor{&mockProcessor{}},
			want:             "static value",
		},
		{
			name:             "multiple parameters",
			value:            "{{.$test_param}} and {{.$unknown}}",
			populationValues: []string{},
			processors:       []ParameterProcessor{&mockProcessor{}},
			want:             "processed and {{.$unknown}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.processWithProcessors(tt.value, tt.populationValues, tt.processors)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHandleAggregation(t *testing.T) {
	tests := []struct {
		name                 string
		params               *datasetmodels.DatasetParams
		field                widgetmodels.Field
		mapping              *widgetmodels.DataMappingFields
		filters              *datasetmodels.FilterModel
		datasetBuilderParams *widgetmodels.DatasetBuilderParams
		want                 *datasetmodels.DatasetParams
		wantErr              bool
	}{
		{
			name: "standard aggregation",
			params: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{},
			},
			field: widgetmodels.Field{
				Column:      "revenue",
				Aggregation: "sum",
				Alias:       "total_revenue",
			},
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			filters:              &datasetmodels.FilterModel{},
			datasetBuilderParams: &widgetmodels.DatasetBuilderParams{},
			want: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{
					{
						Column:   "revenue",
						Function: "sum",
						Alias:    "total_revenue",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "window function first with sort by in field",
			params: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "date", Alias: stringPtr("date")},
				},
			},
			field: widgetmodels.Field{
				Column:      "balance",
				Aggregation: "first",
				Alias:       "opening_balance",
				SortBy: []widgetmodels.SortBy{
					{Column: "date", Order: "ASC"},
				},
			},
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			filters:              &datasetmodels.FilterModel{},
			datasetBuilderParams: &widgetmodels.DatasetBuilderParams{},
			want: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{
					{
						Column:   "balance",
						Function: "sum",
						Alias:    "opening_balance",
					},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "date", Alias: stringPtr("date")},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions: []datasetmodels.Filter{
						{Column: "rn", Operator: "eq", Value: 1},
					},
				},
				Subquery: &datasetmodels.DatasetParams{
					Columns: []datasetmodels.ColumnConfig{
						{Column: "balance"},
						{Column: "date"},
						{Column: datasetconstants.ZampIsDeletedColumn},
					},
					Windows: []datasetmodels.WindowConfig{
						{
							Function: "ROW_NUMBER()",
							PartitionBy: []datasetmodels.ColumnConfig{
								{Column: "date"},
							},
							OrderBy: []datasetmodels.OrderBy{
								{Column: "date", Order: "ASC"},
							},
							Alias: "rn",
						},
					},
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "",
						Conditions: []datasetmodels.Filter{
							{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "window function last with sort by in mapping",
			params: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "date", Alias: stringPtr("date")},
				},
			},
			field: widgetmodels.Field{
				Column:      "balance",
				Aggregation: "last",
				Alias:       "closing_balance",
			},
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
				SortBy: []widgetmodels.SortBy{
					{Column: "date", Order: "DESC"},
				},
			},
			filters:              &datasetmodels.FilterModel{},
			datasetBuilderParams: &widgetmodels.DatasetBuilderParams{},
			want: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{
					{
						Column:   "balance",
						Function: "sum",
						Alias:    "closing_balance",
					},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "date", Alias: stringPtr("date")},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions: []datasetmodels.Filter{
						{Column: "rn", Operator: "eq", Value: 1},
					},
				},
				Subquery: &datasetmodels.DatasetParams{
					Columns: []datasetmodels.ColumnConfig{
						{Column: "balance"},
						{Column: "date"},
						{Column: datasetconstants.ZampIsDeletedColumn},
					},
					Windows: []datasetmodels.WindowConfig{
						{
							Function: "ROW_NUMBER()",
							PartitionBy: []datasetmodels.ColumnConfig{
								{Column: "date"},
							},
							OrderBy: []datasetmodels.OrderBy{
								{Column: "date", Order: "DESC"},
							},
							Alias: "rn",
						},
					},
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "",
						Conditions: []datasetmodels.Filter{
							{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "window function without sort by",
			params: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{},
			},
			field: widgetmodels.Field{
				Column:      "balance",
				Aggregation: "first",
			},
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			filters:              &datasetmodels.FilterModel{},
			datasetBuilderParams: &widgetmodels.DatasetBuilderParams{},
			want:                 nil,
			wantErr:              true,
		},
		{
			name: "window function with time column in partition by",
			params: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "time_stamp", Alias: stringPtr("date")},
				},
			},
			field: widgetmodels.Field{
				Column:      "balance",
				Aggregation: "first",
				SortBy: []widgetmodels.SortBy{
					{Column: "time_stamp", Order: "ASC"},
				},
			},
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			filters: &datasetmodels.FilterModel{},
			datasetBuilderParams: &widgetmodels.DatasetBuilderParams{
				TimeColumns: map[string]string{
					"dataset1": "time_stamp",
				},
				Periodicity: func() *string { s := "month"; return &s }(),
			},
			want: &datasetmodels.DatasetParams{
				Aggregations: []datasetmodels.Aggregation{
					{
						Column:   "balance",
						Function: "sum",
						Alias:    "balance",
					},
				},
				GroupBy: []datasetmodels.GroupBy{
					{Column: "time_stamp", Alias: stringPtr("date")},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions: []datasetmodels.Filter{
						{Column: "rn", Operator: "eq", Value: 1},
					},
				},
				Subquery: &datasetmodels.DatasetParams{
					Columns: []datasetmodels.ColumnConfig{
						{Column: "balance"},
						{Column: "time_stamp"},
						{Column: datasetconstants.ZampIsDeletedColumn},
					},
					Windows: []datasetmodels.WindowConfig{
						{
							Function: "ROW_NUMBER()",
							PartitionBy: []datasetmodels.ColumnConfig{
								{Column: "date_trunc('month', time_stamp)"},
							},
							OrderBy: []datasetmodels.OrderBy{
								{Column: "time_stamp", Order: "ASC"},
							},
							Alias: "rn",
						},
					},
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "",
						Conditions: []datasetmodels.Filter{
							{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			err := baseStrategy.HandleAggregation(tt.params, tt.field, tt.mapping, tt.filters, tt.datasetBuilderParams)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// For window functions, we need to check the subquery structure
			if tt.field.Aggregation == "first" || tt.field.Aggregation == "last" {
				assert.NotNil(t, tt.params.Subquery)
				assert.Equal(t, tt.want.Subquery.Columns, tt.params.Subquery.Columns)
				assert.Equal(t, tt.want.Subquery.Windows[0].Function, tt.params.Subquery.Windows[0].Function)
				assert.Equal(t, tt.want.Subquery.Windows[0].Alias, tt.params.Subquery.Windows[0].Alias)

				// Check partition by columns
				assert.Equal(t, len(tt.want.Subquery.Windows[0].PartitionBy), len(tt.params.Subquery.Windows[0].PartitionBy))
				for i := range tt.want.Subquery.Windows[0].PartitionBy {
					assert.Equal(t, tt.want.Subquery.Windows[0].PartitionBy[i].Column, tt.params.Subquery.Windows[0].PartitionBy[i].Column)
				}

				// Check order by columns
				assert.Equal(t, len(tt.want.Subquery.Windows[0].OrderBy), len(tt.params.Subquery.Windows[0].OrderBy))
				for i := range tt.want.Subquery.Windows[0].OrderBy {
					assert.Equal(t, tt.want.Subquery.Windows[0].OrderBy[i].Column, tt.params.Subquery.Windows[0].OrderBy[i].Column)
					assert.Equal(t, tt.want.Subquery.Windows[0].OrderBy[i].Order, tt.params.Subquery.Windows[0].OrderBy[i].Order)
				}
			}

			// Check aggregations
			assert.Equal(t, len(tt.want.Aggregations), len(tt.params.Aggregations))
			for i := range tt.want.Aggregations {
				assert.Equal(t, tt.want.Aggregations[i].Column, tt.params.Aggregations[i].Column)
				assert.Equal(t, tt.want.Aggregations[i].Function, tt.params.Aggregations[i].Function)
				assert.Equal(t, tt.want.Aggregations[i].Alias, tt.params.Aggregations[i].Alias)
			}

			// Check filters
			if tt.want.Filters.LogicalOperator != "" {
				assert.Equal(t, tt.want.Filters.LogicalOperator, tt.params.Filters.LogicalOperator)
				assert.Equal(t, len(tt.want.Filters.Conditions), len(tt.params.Filters.Conditions))
				for i := range tt.want.Filters.Conditions {
					assert.Equal(t, tt.want.Filters.Conditions[i].Column, tt.params.Filters.Conditions[i].Column)
					assert.Equal(t, tt.want.Filters.Conditions[i].Operator, tt.params.Filters.Conditions[i].Operator)
					assert.Equal(t, tt.want.Filters.Conditions[i].Value, tt.params.Filters.Conditions[i].Value)
				}
			}
		})
	}
}

func TestBuildWindowBasedParams(t *testing.T) {
	tests := []struct {
		name    string
		field   widgetmodels.Field
		sortBy  []widgetmodels.SortBy
		groupBy []datasetmodels.GroupBy
		filters *datasetmodels.FilterModel
		want    *datasetmodels.DatasetParams
		wantErr bool
	}{
		{
			name: "basic window params",
			field: widgetmodels.Field{
				Column: "revenue",
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "date", Order: "ASC"},
			},
			groupBy: []datasetmodels.GroupBy{
				{Column: "category", Alias: stringPtr("category")},
			},
			filters: nil,
			want: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "revenue"},
					{Column: "category"},
					{Column: datasetconstants.ZampIsDeletedColumn},
				},
				Windows: []datasetmodels.WindowConfig{
					{
						Function: "ROW_NUMBER()",
						PartitionBy: []datasetmodels.ColumnConfig{
							{Column: "category"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC"},
						},
						Alias: "rn",
					},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "",
					Conditions: []datasetmodels.Filter{
						{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "window params with existing filters",
			field: widgetmodels.Field{
				Column: "revenue",
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "date", Order: "ASC"},
			},
			groupBy: []datasetmodels.GroupBy{
				{Column: "category", Alias: stringPtr("category")},
			},
			filters: &datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "status", Operator: "eq", Value: "active"},
				},
			},
			want: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "revenue"},
					{Column: "category"},
					{Column: datasetconstants.ZampIsDeletedColumn},
				},
				Windows: []datasetmodels.WindowConfig{
					{
						Function: "ROW_NUMBER()",
						PartitionBy: []datasetmodels.ColumnConfig{
							{Column: "category"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC"},
						},
						Alias: "rn",
					},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions: []datasetmodels.Filter{
						{Column: "status", Operator: "eq", Value: "active"},
						{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "window params with multiple group by",
			field: widgetmodels.Field{
				Column: "revenue",
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "date", Order: "ASC"},
			},
			groupBy: []datasetmodels.GroupBy{
				{Column: "category", Alias: stringPtr("category")},
				{Column: "region", Alias: stringPtr("region")},
			},
			filters: nil,
			want: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "revenue"},
					{Column: "category"},
					{Column: "region"},
					{Column: datasetconstants.ZampIsDeletedColumn},
				},
				Windows: []datasetmodels.WindowConfig{
					{
						Function: "ROW_NUMBER()",
						PartitionBy: []datasetmodels.ColumnConfig{
							{Column: "category"},
							{Column: "region"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC"},
						},
						Alias: "rn",
					},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "",
					Conditions: []datasetmodels.Filter{
						{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "window params with multiple sort by",
			field: widgetmodels.Field{
				Column: "revenue",
			},
			sortBy: []widgetmodels.SortBy{
				{Column: "date", Order: "ASC"},
				{Column: "id", Order: "DESC"},
			},
			groupBy: []datasetmodels.GroupBy{
				{Column: "category", Alias: stringPtr("category")},
			},
			filters: nil,
			want: &datasetmodels.DatasetParams{
				Columns: []datasetmodels.ColumnConfig{
					{Column: "revenue"},
					{Column: "category"},
					{Column: datasetconstants.ZampIsDeletedColumn},
				},
				Windows: []datasetmodels.WindowConfig{
					{
						Function: "ROW_NUMBER()",
						PartitionBy: []datasetmodels.ColumnConfig{
							{Column: "category"},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC"},
							{Column: "id", Order: "DESC"},
						},
						Alias: "rn",
					},
				},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "",
					Conditions: []datasetmodels.Filter{
						{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "error when no sort by provided",
			field:   widgetmodels.Field{Column: "revenue"},
			sortBy:  []widgetmodels.SortBy{},
			groupBy: []datasetmodels.GroupBy{{Column: "category"}},
			filters: nil,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got, err := baseStrategy.BuildWindowBasedParams(tt.field, tt.sortBy, tt.groupBy, tt.filters)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check columns
			assert.Equal(t, len(tt.want.Columns), len(got.Columns))
			for i := range tt.want.Columns {
				assert.Equal(t, tt.want.Columns[i].Column, got.Columns[i].Column)
			}

			// Check windows
			assert.Equal(t, len(tt.want.Windows), len(got.Windows))
			for i := range tt.want.Windows {
				assert.Equal(t, tt.want.Windows[i].Function, got.Windows[i].Function)
				assert.Equal(t, tt.want.Windows[i].Alias, got.Windows[i].Alias)

				// Check partition by
				assert.Equal(t, len(tt.want.Windows[i].PartitionBy), len(got.Windows[i].PartitionBy))
				for j := range tt.want.Windows[i].PartitionBy {
					assert.Equal(t, tt.want.Windows[i].PartitionBy[j].Column, got.Windows[i].PartitionBy[j].Column)
				}

				// Check order by
				assert.Equal(t, len(tt.want.Windows[i].OrderBy), len(got.Windows[i].OrderBy))
				for j := range tt.want.Windows[i].OrderBy {
					assert.Equal(t, tt.want.Windows[i].OrderBy[j].Column, got.Windows[i].OrderBy[j].Column)
					assert.Equal(t, tt.want.Windows[i].OrderBy[j].Order, got.Windows[i].OrderBy[j].Order)
				}
			}

			// Check filters
			assert.Equal(t, tt.want.Filters.LogicalOperator, got.Filters.LogicalOperator)
			assert.Equal(t, len(tt.want.Filters.Conditions), len(got.Filters.Conditions))
			for i := range tt.want.Filters.Conditions {
				assert.Equal(t, tt.want.Filters.Conditions[i].Column, got.Filters.Conditions[i].Column)
				assert.Equal(t, tt.want.Filters.Conditions[i].Operator, got.Filters.Conditions[i].Operator)
				assert.Equal(t, tt.want.Filters.Conditions[i].Value, got.Filters.Conditions[i].Value)
			}
		})
	}
}

func TestAddCurrency(t *testing.T) {
	tests := []struct {
		name     string
		params   *datasetmodels.DatasetParams
		currency *string
		want     *datasetmodels.DatasetParams
	}{
		{
			name: "add currency to params",
			params: &datasetmodels.DatasetParams{
				FxCurrency: nil,
			},
			currency: func() *string { s := "USD"; return &s }(),
			want: &datasetmodels.DatasetParams{
				FxCurrency: func() *string { s := "USD"; return &s }(),
			},
		},
		{
			name: "add currency to params with subquery",
			params: &datasetmodels.DatasetParams{
				Subquery: &datasetmodels.DatasetParams{
					FxCurrency: nil,
				},
			},
			currency: func() *string { s := "EUR"; return &s }(),
			want: &datasetmodels.DatasetParams{
				Subquery: &datasetmodels.DatasetParams{
					FxCurrency: func() *string { s := "EUR"; return &s }(),
				},
			},
		},
		{
			name: "nil currency",
			params: &datasetmodels.DatasetParams{
				FxCurrency: nil,
			},
			currency: nil,
			want: &datasetmodels.DatasetParams{
				FxCurrency: nil,
			},
		},
		{
			name: "overwrite existing currency",
			params: &datasetmodels.DatasetParams{
				FxCurrency: func() *string { s := "GBP"; return &s }(),
			},
			currency: func() *string { s := "USD"; return &s }(),
			want: &datasetmodels.DatasetParams{
				FxCurrency: func() *string { s := "USD"; return &s }(),
			},
		},
		{
			name: "overwrite existing currency in subquery",
			params: &datasetmodels.DatasetParams{
				Subquery: &datasetmodels.DatasetParams{
					FxCurrency: func() *string { s := "GBP"; return &s }(),
				},
			},
			currency: func() *string { s := "EUR"; return &s }(),
			want: &datasetmodels.DatasetParams{
				Subquery: &datasetmodels.DatasetParams{
					FxCurrency: func() *string { s := "EUR"; return &s }(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			baseStrategy.AddCurrency(tt.params, tt.currency)

			if tt.params.Subquery != nil {
				assert.Equal(t, tt.want.Subquery.FxCurrency, tt.params.Subquery.FxCurrency)
			} else {
				assert.Equal(t, tt.want.FxCurrency, tt.params.FxCurrency)
			}
		})
	}
}

func TestInitializeDatasetParams(t *testing.T) {
	tests := []struct {
		name    string
		filters datasetmodels.FilterModel
		want    datasetmodels.DatasetParams
	}{
		{
			name:    "initialize with empty filters",
			filters: datasetmodels.FilterModel{},
			want: datasetmodels.DatasetParams{
				Columns:      []datasetmodels.ColumnConfig{},
				Aggregations: []datasetmodels.Aggregation{},
				GroupBy:      []datasetmodels.GroupBy{},
				Filters:      datasetmodels.FilterModel{},
			},
		},
		{
			name: "initialize with non-empty filters",
			filters: datasetmodels.FilterModel{
				LogicalOperator: "AND",
				Conditions: []datasetmodels.Filter{
					{Column: "status", Operator: "eq", Value: "active"},
				},
			},
			want: datasetmodels.DatasetParams{
				Columns:      []datasetmodels.ColumnConfig{},
				Aggregations: []datasetmodels.Aggregation{},
				GroupBy:      []datasetmodels.GroupBy{},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions: []datasetmodels.Filter{
						{Column: "status", Operator: "eq", Value: "active"},
					},
				},
			},
		},
		{
			name: "initialize with multiple filter conditions",
			filters: datasetmodels.FilterModel{
				LogicalOperator: "OR",
				Conditions: []datasetmodels.Filter{
					{Column: "status", Operator: "eq", Value: "active"},
					{Column: "region", Operator: "eq", Value: "US"},
				},
			},
			want: datasetmodels.DatasetParams{
				Columns:      []datasetmodels.ColumnConfig{},
				Aggregations: []datasetmodels.Aggregation{},
				GroupBy:      []datasetmodels.GroupBy{},
				Filters: datasetmodels.FilterModel{
					LogicalOperator: "OR",
					Conditions: []datasetmodels.Filter{
						{Column: "status", Operator: "eq", Value: "active"},
						{Column: "region", Operator: "eq", Value: "US"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got := baseStrategy.InitializeDatasetParams(tt.filters)

			assert.Equal(t, tt.want.Columns, got.Columns)
			assert.Equal(t, tt.want.Aggregations, got.Aggregations)
			assert.Equal(t, tt.want.GroupBy, got.GroupBy)
			assert.Equal(t, tt.want.Filters, got.Filters)
		})
	}
}

func TestProcessDatasetParams(t *testing.T) {
	tests := []struct {
		name                 string
		mapping              *widgetmodels.DataMappingFields
		datasetbuilderparams widgetmodels.DatasetBuilderParams
		processFields        ProcessFieldsFunc
		want                 widgetmodels.GetDataByDatasetIDParams
		wantErr              bool
	}{
		{
			name: "process dataset params with no errors",
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
				DefaultFilters: &datasetmodels.FilterModel{
					LogicalOperator: "AND",
					Conditions: []datasetmodels.Filter{
						{Column: "status", Operator: "eq", Value: "active"},
					},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{
				Filters: map[string]widgetmodels.WidgetFilters{
					"dataset1": {
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "region", Operator: "eq", Value: "US"},
							},
						},
					},
				},
			},
			processFields: func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
				return nil
			},
			want: widgetmodels.GetDataByDatasetIDParams{
				DatasetID: "dataset1",
				Params: datasetmodels.DatasetParams{
					Columns:      []datasetmodels.ColumnConfig{},
					Aggregations: []datasetmodels.Aggregation{},
					GroupBy:      []datasetmodels.GroupBy{},
					Filters: datasetmodels.FilterModel{
						LogicalOperator: "AND",
						Conditions: []datasetmodels.Filter{
							{Column: "status", Operator: "eq", Value: "active"},
							{Column: "region", Operator: "eq", Value: "US"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "process dataset params with error in processFields",
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			processFields: func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
				return fmt.Errorf("process fields error")
			},
			want:    widgetmodels.GetDataByDatasetIDParams{},
			wantErr: true,
		},
		{
			name: "process dataset params with time column and periodicity",
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{
				TimeColumns: map[string]string{
					"dataset1": "date",
				},
				Periodicity: func() *string { s := "month"; return &s }(),
			},
			processFields: func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
				params.Columns = append(params.Columns, datasetmodels.ColumnConfig{Column: "date"})
				return nil
			},
			want: widgetmodels.GetDataByDatasetIDParams{
				DatasetID: "dataset1",
				Params: datasetmodels.DatasetParams{
					Columns: []datasetmodels.ColumnConfig{
						{Column: "date_trunc('month', date)"},
					},
					Aggregations: []datasetmodels.Aggregation{},
					GroupBy:      []datasetmodels.GroupBy{},
					Filters:      datasetmodels.FilterModel{},
				},
			},
			wantErr: false,
		},
		{
			name: "process dataset params with currency",
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{
				Currency: func() *string { s := "USD"; return &s }(),
			},
			processFields: func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
				return nil
			},
			want: widgetmodels.GetDataByDatasetIDParams{
				DatasetID: "dataset1",
				Params: datasetmodels.DatasetParams{
					Columns:      []datasetmodels.ColumnConfig{},
					Aggregations: []datasetmodels.Aggregation{},
					GroupBy:      []datasetmodels.GroupBy{},
					Filters:      datasetmodels.FilterModel{},
					FxCurrency:   func() *string { s := "USD"; return &s }(),
				},
			},
			wantErr: false,
		},
		{
			name: "process dataset params with sort by",
			mapping: &widgetmodels.DataMappingFields{
				DatasetID: "dataset1",
				SortBy: []widgetmodels.SortBy{
					{Column: "date", Order: "ASC"},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			processFields: func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
				return nil
			},
			want: widgetmodels.GetDataByDatasetIDParams{
				DatasetID: "dataset1",
				Params: datasetmodels.DatasetParams{
					Columns:      []datasetmodels.ColumnConfig{},
					Aggregations: []datasetmodels.Aggregation{},
					GroupBy:      []datasetmodels.GroupBy{},
					Filters:      datasetmodels.FilterModel{},
					OrderBy: []datasetmodels.OrderBy{
						{Column: "date", Order: "ASC"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseStrategy := NewBaseStrategy()
			got, err := baseStrategy.ProcessDatasetParams(tt.mapping, tt.datasetbuilderparams, tt.processFields)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.DatasetID, got.DatasetID)

			// Check columns
			assert.Equal(t, len(tt.want.Params.Columns), len(got.Params.Columns))
			for i := range tt.want.Params.Columns {
				assert.Equal(t, tt.want.Params.Columns[i].Column, got.Params.Columns[i].Column)
			}

			// Check aggregations
			assert.Equal(t, len(tt.want.Params.Aggregations), len(got.Params.Aggregations))

			// Check group by
			assert.Equal(t, len(tt.want.Params.GroupBy), len(got.Params.GroupBy))

			// Check filters
			assert.Equal(t, tt.want.Params.Filters.LogicalOperator, got.Params.Filters.LogicalOperator)
			assert.Equal(t, len(tt.want.Params.Filters.Conditions), len(got.Params.Filters.Conditions))

			// Check currency
			assert.Equal(t, tt.want.Params.FxCurrency, got.Params.FxCurrency)

			// Check order by
			assert.Equal(t, len(tt.want.Params.OrderBy), len(got.Params.OrderBy))
			for i := range tt.want.Params.OrderBy {
				if i < len(got.Params.OrderBy) {
					assert.Equal(t, tt.want.Params.OrderBy[i].Column, got.Params.OrderBy[i].Column)
					assert.Equal(t, tt.want.Params.OrderBy[i].Order, got.Params.OrderBy[i].Order)
				}
			}
		})
	}
}

func TestBasicChartStrategy_ToDatasetParams_EdgeCases(t *testing.T) {
	tests := []struct {
		name                 string
		instance             *widgetmodels.WidgetInstance
		datasetbuilderparams widgetmodels.DatasetBuilderParams
		want                 map[string]widgetmodels.GetDataByDatasetIDParams
		wantErr              bool
	}{
		{
			name: "no mappings",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version:  "1",
					Mappings: []widgetmodels.DataMappingFields{},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			want:                 nil,
			wantErr:              true,
		},
		{
			name: "missing x-axis field",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version: "1",
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.YAxisField: {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing y-axis field",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version: "1",
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.XAxisField: {{Column: "category"}},
							},
						},
					},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns:      []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
						},
						Filters: datasetmodels.FilterModel{},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with time column but no periodicity",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version: "1",
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.XAxisField: {{Column: "date"}},
								widgetconstants.YAxisField: {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{
				TimeColumns: map[string]string{
					"dataset1": "date",
				},
				// No periodicity set
			},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "date", Alias: stringPtr("date")},
						},
						Filters: datasetmodels.FilterModel{},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := BasicChartStrategy{}
			got, err := strategy.ToDatasetParams(tt.instance, tt.datasetbuilderparams)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPieChartStrategy_ToDatasetParams_EdgeCases(t *testing.T) {
	tests := []struct {
		name                 string
		instance             *widgetmodels.WidgetInstance
		datasetbuilderparams widgetmodels.DatasetBuilderParams
		want                 map[string]widgetmodels.GetDataByDatasetIDParams
		wantErr              bool
	}{
		{
			name: "no mappings",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version:  "1",
					Mappings: []widgetmodels.DataMappingFields{},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			want:                 nil,
			wantErr:              true,
		},
		{
			name: "missing slices field",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version: "1",
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.ValuesField: {{Column: "sales", Aggregation: "sum"}},
							},
						},
					},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "sales", Function: "sum", Alias: "sales"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing values field",
			instance: &widgetmodels.WidgetInstance{
				DataMappings: widgetmodels.DataMappings{
					Version: "1",
					Mappings: []widgetmodels.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]widgetmodels.Field{
								widgetconstants.SlicesField: {{Column: "category"}},
							},
						},
					},
				},
			},
			datasetbuilderparams: widgetmodels.DatasetBuilderParams{},
			want: map[string]widgetmodels.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns:      []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{},
						GroupBy: []datasetmodels.GroupBy{
							{Column: "category", Alias: stringPtr("category")},
						},
						Filters: datasetmodels.FilterModel{},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "category", Order: "ASC", Alias: stringPtr("category")},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := PieChartStrategy{}
			got, err := strategy.ToDatasetParams(tt.instance, tt.datasetbuilderparams)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func ptr(s string) *string {
	return &s
}

func TestKPIStrategy_ToDatasetParams_EdgeCases(t *testing.T) {
	tests := []struct {
		name                 string
		instance             *models.WidgetInstance
		datasetbuilderparams models.DatasetBuilderParams
		want                 map[string]models.GetDataByDatasetIDParams
		wantErr              bool
	}{
		{
			name: "no mappings",
			instance: &models.WidgetInstance{
				DataMappings: models.DataMappings{
					Version:  "1",
					Mappings: []models.DataMappingFields{},
				},
			},
			datasetbuilderparams: models.DatasetBuilderParams{},
			want:                 nil,
			wantErr:              true,
		},
		{
			name: "missing primary value field",
			instance: &models.WidgetInstance{
				DataMappings: models.DataMappings{
					Version: "1",
					Mappings: []models.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]models.Field{
								// No primary value field
								widgetconstants.ComparisonValueField: {
									{Column: "previous_revenue", Aggregation: "sum", Type: "number"},
								},
							},
						},
					},
				},
			},
			datasetbuilderparams: models.DatasetBuilderParams{
				Filters: map[string]models.WidgetFilters{},
			},
			want: map[string]models.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns:      []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{},
						GroupBy:      []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "",
							Conditions:      nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with window function in primary value field",
			instance: &models.WidgetInstance{
				DataMappings: models.DataMappings{
					Version: "1",
					Mappings: []models.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							SortBy: []models.SortBy{
								{Column: "date", Order: "ASC"},
							},
							Fields: map[string][]models.Field{
								widgetconstants.PrimaryValueField: {
									{Column: "balance", Aggregation: "first", Type: "number"},
								},
							},
						},
					},
				},
			},
			datasetbuilderparams: models.DatasetBuilderParams{
				Filters: map[string]models.WidgetFilters{},
			},
			want: map[string]models.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "balance", Function: "sum", Alias: "balance"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "AND",
							Conditions: []datasetmodels.Filter{
								{Column: "rn", Operator: "eq", Value: 1},
							},
						},
						OrderBy: []datasetmodels.OrderBy{
							{Column: "date", Order: "ASC", Alias: stringPtr("date")},
						},
						Subquery: &datasetmodels.DatasetParams{
							Columns: []datasetmodels.ColumnConfig{
								{Column: "balance"},
								{Column: datasetconstants.ZampIsDeletedColumn},
							},
							Windows: []datasetmodels.WindowConfig{
								{
									Function:    "ROW_NUMBER()",
									PartitionBy: []datasetmodels.ColumnConfig{},
									OrderBy: []datasetmodels.OrderBy{
										{Column: "date", Order: "ASC"},
									},
									Alias: "rn",
								},
							},
							Filters: datasetmodels.FilterModel{
								LogicalOperator: "",
								Conditions: []datasetmodels.Filter{
									{Column: datasetconstants.ZampIsDeletedColumn, Operator: "eq", Value: false},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with window function but missing sort by",
			instance: &models.WidgetInstance{
				DataMappings: models.DataMappings{
					Version: "1",
					Mappings: []models.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							// No sort by defined
							Fields: map[string][]models.Field{
								widgetconstants.PrimaryValueField: {
									{Column: "balance", Aggregation: "first", Type: "number"},
								},
							},
						},
					},
				},
			},
			datasetbuilderparams: models.DatasetBuilderParams{
				Filters: map[string]models.WidgetFilters{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with both primary and comparison value fields",
			instance: &models.WidgetInstance{
				DataMappings: models.DataMappings{
					Version: "1",
					Mappings: []models.DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]models.Field{
								widgetconstants.PrimaryValueField: {
									{Column: "current_revenue", Aggregation: "sum", Type: "number"},
								},
								widgetconstants.ComparisonValueField: {
									{Column: "previous_revenue", Aggregation: "sum", Type: "number"},
								},
							},
						},
					},
				},
			},
			datasetbuilderparams: models.DatasetBuilderParams{
				Filters:  map[string]models.WidgetFilters{},
				Currency: func() *string { s := "USD"; return &s }(),
			},
			want: map[string]models.GetDataByDatasetIDParams{
				"ref1": {
					DatasetID: "dataset1",
					Params: datasetmodels.DatasetParams{
						Columns: []datasetmodels.ColumnConfig{},
						Aggregations: []datasetmodels.Aggregation{
							{Column: "current_revenue", Function: "sum", Alias: "current_revenue"},
						},
						GroupBy: []datasetmodels.GroupBy{},
						Filters: datasetmodels.FilterModel{
							LogicalOperator: "",
							Conditions:      nil,
						},
						FxCurrency: func() *string { s := "USD"; return &s }(),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := KPIStrategy{BaseStrategy: *NewBaseStrategy()}
			got, err := strategy.ToDatasetParams(tt.instance, tt.datasetbuilderparams)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
