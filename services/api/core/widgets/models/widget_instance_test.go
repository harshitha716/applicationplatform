package models

import (
	"encoding/json"
	"testing"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestField_GetAlias(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		expected string
	}{
		{
			name:     "returns alias when set",
			field:    Field{Column: "col1", Alias: "alias1"},
			expected: "alias1",
		},
		{
			name:     "returns column when alias not set",
			field:    Field{Column: "col1"},
			expected: "col1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.GetAlias()
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestField_GetExpression(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		expected string
	}{
		{
			name:     "returns expression when set",
			field:    Field{Column: "col1", Expression: "SUM(col1)"},
			expected: "SUM(col1)",
		},
		{
			name:     "returns column when expression not set",
			field:    Field{Column: "col1"},
			expected: "col1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.GetExpression()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWidgetInstance_FromDB(t *testing.T) {
	sheetID := uuid.New()
	instanceID := uuid.New()

	dataMappings := DataMappings{
		Version: DataMappingVersion1,
		Mappings: []DataMappingFields{
			{
				DatasetID: "dataset1",
				Fields: map[string][]Field{
					"metrics": {{
						Column:      "revenue",
						Aggregation: "sum",
						Type:        "number",
						FieldType:   "metric",
					}},
				},
			},
		},
	}

	mappingsJSON, err := json.Marshal(dataMappings)
	require.NoError(t, err)

	require.NoError(t, err)

	dbModel := &dbmodels.WidgetInstance{
		ID:           instanceID,
		SheetID:      sheetID,
		Title:        "Test Widget",
		DataMappings: mappingsJSON,
		WidgetType:   "chart",
	}

	wi := &WidgetInstance{}
	err = wi.FromDB(dbModel)
	require.NoError(t, err)

	assert.Equal(t, instanceID, wi.ID)
	assert.Equal(t, sheetID, wi.SheetID)
	assert.Equal(t, "Test Widget", wi.Title)
	assert.Equal(t, "chart", wi.WidgetType)
	assert.Equal(t, DataMappingVersion1, wi.DataMappings.Version)
	assert.Len(t, wi.DataMappings.Mappings, 1)
}

func TestWidgetInstance_FromDB_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		dbModel *dbmodels.WidgetInstance
		wantErr bool
	}{
		{
			name: "invalid JSON in data mappings",
			dbModel: &dbmodels.WidgetInstance{
				ID:           uuid.New(),
				SheetID:      uuid.New(),
				Title:        "Test Widget",
				DataMappings: []byte(`{"invalid json`),
				WidgetType:   "chart",
			},
			wantErr: true,
		},
		{
			name: "with display config",
			dbModel: func() *dbmodels.WidgetInstance {
				displayConfig := json.RawMessage(`{"color":"blue"}`)
				dataMappings, _ := json.Marshal(DataMappings{Version: DataMappingVersion1})
				return &dbmodels.WidgetInstance{
					ID:            uuid.New(),
					SheetID:       uuid.New(),
					Title:         "Test Widget",
					DataMappings:  dataMappings,
					WidgetType:    "chart",
					DisplayConfig: &displayConfig,
				}
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wi := &WidgetInstance{}
			err := wi.FromDB(tt.dbModel)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			if tt.dbModel.DisplayConfig != nil {
				assert.NotNil(t, wi.DisplayConfig)
				assert.Equal(t, string(*tt.dbModel.DisplayConfig), string(*wi.DisplayConfig))
			}
		})
	}
}

func TestWidgetInstance_ToDB(t *testing.T) {
	sheetID := uuid.New()
	instanceID := uuid.New()

	wi := &WidgetInstance{
		ID:         instanceID,
		WidgetType: "bar_chart",
		SheetID:    sheetID,
		Title:      "Test Widget",
		DataMappings: DataMappings{
			Version: DataMappingVersion1,
			Mappings: []DataMappingFields{
				{
					DatasetID: "dataset1",
					Fields: map[string][]Field{
						"metrics": {{
							Column:      "revenue",
							Aggregation: "sum",
							Type:        "number",
							FieldType:   "metric",
						}},
					},
				},
			},
		},
	}

	dbModel, err := wi.ToDB()
	require.NoError(t, err)

	assert.Equal(t, instanceID, dbModel.ID)
	assert.Equal(t, sheetID, dbModel.SheetID)
	assert.Equal(t, "Test Widget", dbModel.Title)
	assert.Equal(t, "bar_chart", dbModel.WidgetType)
	var dataMappings DataMappings
	err = json.Unmarshal(dbModel.DataMappings, &dataMappings)
	require.NoError(t, err)

	assert.Equal(t, DataMappingVersion1, dataMappings.Version)
	assert.Len(t, dataMappings.Mappings, 1)
}

func TestWidgetInstance_ToDB_ErrorCase(t *testing.T) {
	// Test with display config
	wiWithConfig := &WidgetInstance{
		ID:         uuid.New(),
		WidgetType: "bar_chart",
		SheetID:    uuid.New(),
		Title:      "Test Widget",
		DataMappings: DataMappings{
			Version: DataMappingVersion1,
		},
		DisplayConfig: func() *json.RawMessage {
			raw := json.RawMessage(`{"color":"blue"}`)
			return &raw
		}(),
	}
	
	dbModel, err := wiWithConfig.ToDB()
	require.NoError(t, err)
	assert.NotNil(t, dbModel.DisplayConfig)
	assert.Equal(t, string(*wiWithConfig.DisplayConfig), string(*dbModel.DisplayConfig))
}

func TestCreateWidgetInstancePayload_ToModel(t *testing.T) {
	sheetID := uuid.New()

	dataMappings := map[string]interface{}{
		"version": "1",
		"mappings": []map[string]interface{}{
			{
				"dataset_id": "dataset1",
				"fields": map[string]interface{}{
					"metrics": []map[string]interface{}{
						{
							"column":      "revenue",
							"aggregation": "sum",
							"type":        "number",
							"field_type":  "metric",
						},
					},
				},
			},
		},
	}

	dataMappingsJSON, err := json.Marshal(dataMappings)
	require.NoError(t, err)

	require.NoError(t, err)

	payload := &CreateWidgetInstancePayload{
		WidgetType:   "bar_chart",
		SheetID:      sheetID.String(),
		Title:        "Test Widget",
		DataMappings: string(dataMappingsJSON),
	}

	model, err := payload.ToModel()
	require.NoError(t, err)

	assert.Equal(t, sheetID, model.SheetID)
	assert.Equal(t, "Test Widget", model.Title)
	assert.Equal(t, DataMappingVersion1, model.DataMappings.Version)
	assert.Len(t, model.DataMappings.Mappings, 1)
}

func TestCreateWidgetInstancePayload_ToModel_ErrorCase(t *testing.T) {
	payload := &CreateWidgetInstancePayload{
		WidgetType:   "bar_chart",
		SheetID:      uuid.New().String(),
		Title:        "Test Widget",
		DataMappings: `{"invalid json`,
	}

	_, err := payload.ToModel()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data mappings JSON")
}

func TestUpdateWidgetInstancePayload_ToModel(t *testing.T) {
	instanceID := uuid.New()
	sheetID := uuid.New()
	widgetType := "bar_chart"
	title := "Updated Widget"
	
	tests := []struct {
		name    string
		payload UpdateWidgetInstancePayload
		want    *WidgetInstance
		wantErr bool
	}{
		{
			name: "full update with all fields",
			payload: UpdateWidgetInstancePayload{
				WidgetInstanceID: instanceID.String(),
				WidgetType:       &widgetType,
				SheetID:          func() *string { s := sheetID.String(); return &s }(),
				Title:            &title,
				DataMappings:     func() *string { s := `{"version":"1","mappings":[{"dataset_id":"dataset1","ref":"ref1","fields":{"metrics":[{"column":"revenue","aggregation":"sum","type":"number","field_type":"metric"}]}}]}`; return &s }(),
				DisplayConfig:    func() *string { s := `{"color":"blue"}`; return &s }(),
			},
			want: &WidgetInstance{
				ID:         instanceID,
				WidgetType: "bar_chart",
				SheetID:    sheetID,
				Title:      "Updated Widget",
				DataMappings: DataMappings{
					Version: DataMappingVersion1,
					Mappings: []DataMappingFields{
						{
							DatasetID: "dataset1",
							Ref:       "ref1",
							Fields: map[string][]Field{
								"metrics": {
									{
										Column:      "revenue",
										Aggregation: "sum",
										Type:        "number",
										FieldType:   "metric",
									},
								},
							},
						},
					},
				},
				DisplayConfig: func() *json.RawMessage {
					raw := json.RawMessage(`{"color":"blue"}`)
					return &raw
				}(),
			},
			wantErr: false,
		},
		{
			name: "partial update with only title",
			payload: UpdateWidgetInstancePayload{
				WidgetInstanceID: instanceID.String(),
				Title:            &title,
			},
			want: &WidgetInstance{
				ID:           instanceID,
				Title:        "Updated Widget",
				DataMappings: DataMappings{},
			},
			wantErr: false,
		},
		{
			name: "invalid JSON in data mappings",
			payload: UpdateWidgetInstancePayload{
				WidgetInstanceID: instanceID.String(),
				DataMappings:     func() *string { s := `{"invalid json`; return &s }(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty fields should be handled correctly",
			payload: UpdateWidgetInstancePayload{
				WidgetInstanceID: instanceID.String(),
				WidgetType:       func() *string { s := ""; return &s }(),
				Title:            func() *string { s := ""; return &s }(),
				SheetID:          func() *string { s := ""; return &s }(),
				DataMappings:     func() *string { s := ""; return &s }(),
				DisplayConfig:    func() *string { s := ""; return &s }(),
			},
			want: &WidgetInstance{
				ID:           instanceID,
				DataMappings: DataMappings{},
			},
			wantErr: false,
		},
		{
			name: "invalid UUID in sheet ID",
			payload: UpdateWidgetInstancePayload{
				WidgetInstanceID: instanceID.String(),
				SheetID:          func() *string { s := "invalid-uuid"; return &s }(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.payload.ToModel()
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			
			if tt.want.DisplayConfig != nil {
				assert.NotNil(t, got.DisplayConfig)
				// Compare the string representation of the raw messages
				gotStr := string(*got.DisplayConfig)
				wantStr := string(*tt.want.DisplayConfig)
				assert.Equal(t, wantStr, gotStr)
				
				// Remove DisplayConfig for the rest of the comparison
				got.DisplayConfig = nil
				tt.want.DisplayConfig = nil
			}
			
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSortBy_GetColumn(t *testing.T) {
	tests := []struct {
		name   string
		sortBy SortBy
		want   string
	}{
		{
			name:   "returns alias when set",
			sortBy: SortBy{Column: "col1", Alias: "alias1"},
			want:   "alias1",
		},
		{
			name:   "returns column when alias not set",
			sortBy: SortBy{Column: "col1"},
			want:   "col1",
		},
		{
			name:   "handles empty column and alias",
			sortBy: SortBy{},
			want:   "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sortBy.GetColumn()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDataMappingFields_GetRef(t *testing.T) {
	tests := []struct {
		name    string
		mapping DataMappingFields
		want    string
	}{
		{
			name:    "returns ref when set",
			mapping: DataMappingFields{Ref: "test-ref"},
			want:    "test-ref",
		},
		{
			name:    "returns empty string when ref not set",
			mapping: DataMappingFields{},
			want:    "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mapping.GetRef()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  *string
	}{
		{
			name:  "nil input returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty string returns nil",
			input: func() *string { s := ""; return &s }(),
			want:  nil,
		},
		{
			name:  "non-empty string returns pointer to string",
			input: func() *string { s := "test"; return &s }(),
			want:  func() *string { s := "test"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringPtr(tt.input)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}
