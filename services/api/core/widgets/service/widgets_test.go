package widgets

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	dataplatformdataConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	models "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	dbModels "github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mockDatasetService "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mockWidgets "github.com/Zampfi/application-platform/services/api/mocks/core/widgets/service"
	mockStore "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewWidgetsService(t *testing.T) {
	t.Parallel()

	ms := mockStore.NewMockStore(t)
	mockDatasetSvc := mockDatasetService.NewMockDatasetService(t)
	service := NewWidgetsService(ms, mockDatasetSvc)

	assert.NotNil(t, service)
	assert.Equal(t, ms, service.store)
	assert.Equal(t, mockDatasetSvc, service.datasetService)
}

func TestGetWidgetInstance(t *testing.T) {
	t.Parallel()

	widgetID := uuid.New()

	// Initialize default data mappings
	defaultDataMappings := models.DataMappings{
		Version:  models.DataMappingVersion1,
		Mappings: []models.DataMappingFields{},
	}

	// Initialize default filters

	// Marshal both to JSON
	dataMappingsJSON, err := json.Marshal(defaultDataMappings)
	require.NoError(t, err)
	require.NoError(t, err)

	testWidget := dbModels.WidgetInstance{
		ID:           widgetID,
		Title:        "Test Widget",
		DataMappings: dataMappingsJSON,
	}

	tests := []struct {
		name      string
		widgetID  uuid.UUID
		want      models.WidgetInstance
		mockSetup func(*mockWidgets.MockWidgetsServiceStore)
		wantErr   bool
	}{
		{
			name:     "success",
			widgetID: widgetID,
			want: models.WidgetInstance{
				ID:           widgetID,
				Title:        "Test Widget",
				DataMappings: defaultDataMappings,
			},
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(testWidget, nil)
			},
			wantErr: false,
		},
		{
			name:     "store error",
			widgetID: widgetID,
			want:     models.WidgetInstance{},
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(dbModels.WidgetInstance{}, errors.New("test error"))
			},
			wantErr: true,
		},
		{
			name:     "invalid JSON in data mappings",
			widgetID: widgetID,
			want:     models.WidgetInstance{},
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				invalidWidget := testWidget
				invalidWidget.DataMappings = []byte(`{"invalid json`)
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(invalidWidget, nil)
			},
			wantErr: true,
		},
		{
			name:     "empty widget ID",
			widgetID: uuid.Nil,
			want:     models.WidgetInstance{},
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, uuid.Nil).Return(dbModels.WidgetInstance{}, errors.New("widget not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mockWidgets.NewMockWidgetsServiceStore(t)
			tt.mockSetup(mockStore)

			service := &widgetsService{store: mockStore}
			got, err := service.GetWidgetInstance(context.Background(), tt.widgetID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetWidgetInstanceData(t *testing.T) {

	// Only keep truly shared constants
	datasetID := "sales_dataset"
	orgID := uuid.New()
	userID := uuid.New()

	type testSetup struct {
		widgetID uuid.UUID
		widget   dbModels.WidgetInstance
	}

	tests := []struct {
		name      string
		setup     testSetup
		want      []datasetmodels.DatasetData
		mockSetup func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup)
		wantErr   bool
	}{
		{
			name: "success - single mapping",
			setup: func() testSetup {
				widgetID := uuid.New()
				dataMappings := models.DataMappings{
					Version: models.DataMappingVersion1,
					Mappings: []models.DataMappingFields{
						{
							DatasetID: datasetID,
							Ref:       "ref1",
							Fields: map[string][]models.Field{
								"sales_measure": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
								"category_dimension": {
									{
										FieldType: "dimension",
										Column:    "category",
									},
								},
							},
						},
					},
				}
				dataMappingsJSON, _ := json.Marshal(dataMappings)

				return testSetup{
					widgetID: widgetID,
					widget: dbModels.WidgetInstance{
						ID:           widgetID,
						WidgetType:   "bar_chart",
						DataMappings: dataMappingsJSON,
					},
				}
			}(),
			want: []datasetmodels.DatasetData{
				{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "category", DatabaseType: "STRING"},
							{Name: "sales", DatabaseType: "NUMBER"},
							{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						},
						Rows: []map[string]interface{}{
							{"category": "A", "sales": 100, "__REF": "ref1"},
							{"category": "B", "sales": 200, "__REF": "ref1"},
						},
					},
				},
			},
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, nil)
				mds.EXPECT().GetDataByDatasetId(
					mock.Anything,
					orgID,
					datasetID,
					mock.MatchedBy(func(params datasetmodels.DatasetParams) bool {
						return true
					}),
				).Return(datasetmodels.DatasetData{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "category", DatabaseType: "STRING"},
							{Name: "sales", DatabaseType: "NUMBER"},
						},
						Rows: []map[string]interface{}{
							{"category": "A", "sales": 100},
							{"category": "B", "sales": 200},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - widget instance not found",
			setup: func() testSetup {
				widgetID := uuid.New()
				return testSetup{
					widgetID: widgetID,
					widget:   dbModels.WidgetInstance{},
				}
			}(),
			want: nil,
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "error - invalid widget instance data mappings",
			setup: func() testSetup {
				widgetID := uuid.New()
				return testSetup{
					widgetID: widgetID,
					widget: dbModels.WidgetInstance{
						ID:           widgetID,
						DataMappings: []byte(`{"invalid json`),
					},
				}
			}(),
			want: nil,
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, nil)
			},
			wantErr: true,
		},
		{
			name: "error - invalid widget type",
			setup: func() testSetup {
				widgetID := uuid.New()
				dataMappings := models.DataMappings{
					Version:  models.DataMappingVersion1,
					Mappings: []models.DataMappingFields{},
				}
				dataMappingsJSON, _ := json.Marshal(dataMappings)
				return testSetup{
					widgetID: widgetID,
					widget: dbModels.WidgetInstance{
						ID:           widgetID,
						WidgetType:   "invalid_widget_type",
						DataMappings: dataMappingsJSON,
					},
				}
			}(),
			want: nil,
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, nil)
			},
			wantErr: true,
		},
		{
			name: "error - dataset service error",
			setup: func() testSetup {
				widgetID := uuid.New()
				dataMappings := models.DataMappings{
					Version: models.DataMappingVersion1,
					Mappings: []models.DataMappingFields{
						{
							DatasetID: datasetID,
							Ref:       "ref1",
							Fields: map[string][]models.Field{
								"sales_measure": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
							},
						},
					},
				}
				dataMappingsJSON, _ := json.Marshal(dataMappings)
				return testSetup{
					widgetID: widgetID,
					widget: dbModels.WidgetInstance{
						ID:           widgetID,
						WidgetType:   "bar_chart",
						DataMappings: dataMappingsJSON,
					},
				}
			}(),
			want: nil,
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, nil)
				mds.EXPECT().GetDataByDatasetId(
					mock.Anything,
					orgID,
					datasetID,
					mock.MatchedBy(func(params datasetmodels.DatasetParams) bool {
						return true
					}),
				).Return(datasetmodels.DatasetData{}, errors.New("dataset service error"))
			},
			wantErr: true,
		},
		{
			name: "success - with tag column",
			setup: func() testSetup {
				widgetID := uuid.New()
				dataMappings := models.DataMappings{
					Version: models.DataMappingVersion1,
					Mappings: []models.DataMappingFields{
						{
							DatasetID: datasetID,
							Ref:       "ref1",
							Fields: map[string][]models.Field{
								"sales_measure": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
							},
						},
					},
				}
				dataMappingsJSON, _ := json.Marshal(dataMappings)
				return testSetup{
					widgetID: widgetID,
					widget: dbModels.WidgetInstance{
						ID:           widgetID,
						WidgetType:   "bar_chart",
						DataMappings: dataMappingsJSON,
					},
				}
			}(),
			want: []datasetmodels.DatasetData{
				{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "sales", DatabaseType: "NUMBER"},
							{Name: "tags", DatabaseType: "STRING"},
							{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						},
						Rows: []map[string]interface{}{
							{"sales": 100, "tags": "electronics", "__REF": "ref1"},
							{"sales": 200, "tags": "clothing", "__REF": "ref1"},
						},
					},
				},
			},
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, nil)
				mds.EXPECT().GetDataByDatasetId(
					mock.Anything,
					orgID,
					datasetID,
					mock.MatchedBy(func(params datasetmodels.DatasetParams) bool {
						return true
					}),
				).Return(datasetmodels.DatasetData{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "sales", DatabaseType: "NUMBER"},
							{Name: "tags", DatabaseType: "STRING"},
						},
						Rows: []map[string]interface{}{
							{"sales": 100, "tags": "electronics"},
							{"sales": 200, "tags": "clothing"},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success - multiple mappings",
			setup: func() testSetup {
				widgetID := uuid.New()
				dataMappings := models.DataMappings{
					Version: models.DataMappingVersion1,
					Mappings: []models.DataMappingFields{
						{
							DatasetID: datasetID,
							Ref:       "ref1",
							Fields: map[string][]models.Field{
								"rows": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
								"columns": {
									{
										FieldType: "dimension",
										Column:    "category",
									},
								},
								"values": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
							},
						},
						{
							DatasetID: datasetID,
							Ref:       "ref2",
							Fields: map[string][]models.Field{
								"rows": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
								"columns": {
									{
										FieldType: "dimension",
										Column:    "category",
									},
								},
								"values": {
									{
										FieldType:   "measure",
										Column:      "sales",
										Aggregation: "sum",
									},
								},
							},
						},
					},
				}
				dataMappingsJSON, _ := json.Marshal(dataMappings)

				return testSetup{
					widgetID: widgetID,
					widget: dbModels.WidgetInstance{
						ID:           widgetID,
						WidgetType:   "pivot_table",
						DataMappings: dataMappingsJSON,
					},
				}
			}(),
			want: []datasetmodels.DatasetData{
				{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "category", DatabaseType: "STRING"},
							{Name: "sales", DatabaseType: "NUMBER"},
							{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						},
						Rows: []map[string]interface{}{
							{"category": "A", "sales": 100, "__REF": "ref1"},
							{"category": "B", "sales": 200, "__REF": "ref1"},
						},
					},
				},
				{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "category", DatabaseType: "STRING"},
							{Name: "sales", DatabaseType: "NUMBER"},
							{Name: "__REF", DatabaseType: string(dataplatformdataConstants.StringDataType)},
						},
						Rows: []map[string]interface{}{
							{"category": "A", "sales": 100, "__REF": "ref2"},
							{"category": "B", "sales": 200, "__REF": "ref2"},
						},
					},
				},
			},
			mockSetup: func(t *testing.T, ms *mockWidgets.MockWidgetsServiceStore, mds *mockDatasetService.MockDatasetService, setup testSetup) {
				ms.EXPECT().GetWidgetInstanceByID(mock.Anything, setup.widgetID).Return(setup.widget, nil)

				// First call returns data set A
				mds.EXPECT().GetDataByDatasetId(
					mock.Anything,
					orgID,
					datasetID,
					mock.MatchedBy(func(params datasetmodels.DatasetParams) bool {
						return true
					}),
				).Return(datasetmodels.DatasetData{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "category", DatabaseType: "STRING"},
							{Name: "sales", DatabaseType: "NUMBER"},
						},
						Rows: []map[string]interface{}{
							{"category": "A", "sales": 100},
							{"category": "B", "sales": 200},
						},
					},
				}, nil).Once()

				// Second call returns data set B
				mds.EXPECT().GetDataByDatasetId(
					mock.Anything,
					orgID,
					datasetID,
					mock.MatchedBy(func(params datasetmodels.DatasetParams) bool {
						return true
					}),
				).Return(datasetmodels.DatasetData{
					QueryResult: dataplatformmodels.QueryResult{
						Columns: []dataplatformmodels.ColumnMetadata{
							{Name: "category", DatabaseType: "STRING"},
							{Name: "sales", DatabaseType: "NUMBER"},
						},
						Rows: []map[string]interface{}{
							{"category": "A", "sales": 100},
							{"category": "B", "sales": 200},
						},
					},
				}, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create new mocks for each test case
			mockStore := mockWidgets.NewMockWidgetsServiceStore(t)
			mockDatasetSvc := mockDatasetService.NewMockDatasetService(t)

			// Pass the testing.T to mockSetup
			tt.mockSetup(t, mockStore, mockDatasetSvc, tt.setup)

			service := &widgetsService{
				store:          mockStore,
				datasetService: mockDatasetSvc,
			}

			ctx := apicontext.AddAuthToContext(
				context.Background(),
				"user",
				userID,
				[]uuid.UUID{orgID},
			)

			got, err := service.GetWidgetInstanceData(ctx, orgID, tt.setup.widgetID, models.GetWidgetInstanceDataQueryParams{
				Filters: []models.WidgetFilters{
					{
						DatasetID: datasetID,
						Filters:   datasetmodels.FilterModel{},
					},
				},
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateWidgetInstance(t *testing.T) {
	//t.Parallel()

	widgetID := uuid.New()

	// Create model for input/output
	testWidget := models.WidgetInstance{
		ID:    widgetID,
		Title: "Test Widget",
		DataMappings: models.DataMappings{
			Version: models.DataMappingVersion1,
			Mappings: []models.DataMappingFields{
				{
					DatasetID: "sales_dataset",
					Fields: map[string][]models.Field{
						"sales_measure": {
							{FieldType: "measure", Column: "sales", Aggregation: "sum"},
						},
					},
				},
			},
		},
		WidgetType: "bar_chart",
	}

	// Initialize data mappings JSON
	dataMappingsJSON, err := json.Marshal(testWidget.DataMappings)
	require.NoError(t, err)

	// Create DB model for mock return with all necessary fields
	testWidgetDB := &dbModels.WidgetInstance{
		ID:           widgetID,
		Title:        "Test Widget",
		DataMappings: dataMappingsJSON,
		WidgetType:   "bar_chart",
	}

	tests := []struct {
		name      string
		input     models.WidgetInstance
		want      *models.WidgetInstance
		mockSetup func(*mockWidgets.MockWidgetsServiceStore)
		wantErr   bool
	}{
		{
			name:  "success",
			input: testWidget,
			want:  &testWidget,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().CreateWidgetInstance(
					mock.Anything,
					mock.MatchedBy(func(wi *dbModels.WidgetInstance) bool {
						return wi.ID == widgetID && wi.Title == "Test Widget"
					}),
				).Return(testWidgetDB, nil)
			},
			wantErr: false,
		},
		{
			name:  "store error",
			input: testWidget,
			want:  nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().CreateWidgetInstance(
					mock.Anything,
					mock.MatchedBy(func(wi *dbModels.WidgetInstance) bool {
						return wi.ID == widgetID && wi.Title == "Test Widget"
					}),
				).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
		{
			name: "error converting to DB model",
			input: func() models.WidgetInstance {
				// Create a widget with a channel as a field, which can't be marshaled to JSON
				invalidWidget := models.WidgetInstance{
					ID:         widgetID,
					Title:      "Test Widget",
					WidgetType: "bar_chart",
					DataMappings: models.DataMappings{
						Version:  models.DataMappingVersion1,
						Mappings: []models.DataMappingFields{},
					},
				}
				// Set a field that can't be marshaled to JSON
				invalidWidget.DataMappings.Version = "invalid\xffversion"
				return invalidWidget
			}(),
			want: nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				// We need to set up the mock even though it shouldn't be called
				// because the test framework checks expectations before the function runs
				m.EXPECT().CreateWidgetInstance(mock.Anything, mock.Anything).Return(nil, errors.New("should not be called")).Maybe()
			},
			wantErr: true,
		},
		{
			name: "success with display config",
			input: func() models.WidgetInstance {
				wi := testWidget
				displayConfig := json.RawMessage(`{"color":"blue"}`)
				wi.DisplayConfig = &displayConfig
				return wi
			}(),
			want: func() *models.WidgetInstance {
				wi := testWidget
				displayConfig := json.RawMessage(`{"color":"blue"}`)
				wi.DisplayConfig = &displayConfig
				return &wi
			}(),
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().CreateWidgetInstance(
					mock.Anything,
					mock.MatchedBy(func(wi *dbModels.WidgetInstance) bool {
						return wi.ID == widgetID && wi.Title == "Test Widget" && wi.DisplayConfig != nil
					}),
				).Return(func() *dbModels.WidgetInstance {
					db := *testWidgetDB
					displayConfig := json.RawMessage(`{"color":"blue"}`)
					db.DisplayConfig = &displayConfig
					return &db
				}(), nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockStore := mockWidgets.NewMockWidgetsServiceStore(t)
			tt.mockSetup(mockStore)

			service := &widgetsService{store: mockStore}
			got, err := service.CreateWidgetInstance(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateWidgetInstance(t *testing.T) {
	t.Parallel()

	widgetID := uuid.New()
	sheetID := uuid.New()

	// Initialize default data mappings
	defaultDataMappings := models.DataMappings{
		Version: models.DataMappingVersion1,
		Mappings: []models.DataMappingFields{
			{
				DatasetID: "test_dataset",
				Fields: map[string][]models.Field{
					"test_field": {{
						FieldType:   "measure",
						Column:      "test",
						Aggregation: "sum",
					}},
				},
			},
		},
	}

	// Marshal both to JSON for DB model
	dataMappingsJSON, err := json.Marshal(defaultDataMappings)
	require.NoError(t, err)

	// Create existing widget instance
	existingWidget := dbModels.WidgetInstance{
		ID:           widgetID,
		Title:        "Test Widget",
		WidgetType:   "bar_chart",
		DataMappings: dataMappingsJSON,
		SheetID:      uuid.Nil,
	}

	tests := []struct {
		name      string
		input     models.WidgetInstance
		want      *models.WidgetInstance
		mockSetup func(*mockWidgets.MockWidgetsServiceStore)
		wantErr   bool
	}{
		{
			name: "success - update widget type",
			input: models.WidgetInstance{
				ID:         widgetID,
				WidgetType: "line_chart",
			},
			want: &models.WidgetInstance{
				ID:           widgetID,
				Title:        "Test Widget",
				WidgetType:   "line_chart",
				DataMappings: defaultDataMappings,
			},
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				// First call for GetWidgetInstanceByID
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(existingWidget, nil)

				// Expect UpdateWidgetInstance with updated widget type
				m.EXPECT().UpdateWidgetInstance(
					mock.Anything,
					mock.MatchedBy(func(wi *dbModels.WidgetInstance) bool {
						return wi.ID == widgetID && wi.WidgetType == "line_chart"
					}),
				).Return(&dbModels.WidgetInstance{
					ID:           widgetID,
					Title:        "Test Widget",
					WidgetType:   "line_chart",
					DataMappings: dataMappingsJSON,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - widget not found",
			input: models.WidgetInstance{
				ID:         widgetID,
				WidgetType: "line_chart",
			},
			want: nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(dbModels.WidgetInstance{}, errors.New("widget not found"))
			},
			wantErr: true,
		},
		{
			name: "error - converting to DB model",
			input: func() models.WidgetInstance {
				// Create a widget with invalid data that will cause ToDB to fail
				wi := models.WidgetInstance{
					ID:         widgetID,
					WidgetType: "line_chart",
					DataMappings: models.DataMappings{
						Version: models.DataMappingVersion1,
						Mappings: []models.DataMappingFields{
							{
								DatasetID: "invalid",
								Fields: map[string][]models.Field{
									"invalid": {
										{
											FieldType: "invalid",
											// Use a function that can't be marshaled to JSON
											Column: string([]byte{0xff, 0xfe}),
										},
									},
								},
							},
						},
					},
				}
				return wi
			}(),
			want: nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(existingWidget, nil)
				// We need to set up the mock even though it shouldn't be called
				// because the test framework checks expectations before the function runs
				m.EXPECT().UpdateWidgetInstance(mock.Anything, mock.Anything).Return(nil, errors.New("should not be called")).Maybe()
			},
			wantErr: true,
		},
		{
			name: "error - update fails",
			input: models.WidgetInstance{
				ID:         widgetID,
				WidgetType: "line_chart",
			},
			want: nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(existingWidget, nil)
				m.EXPECT().UpdateWidgetInstance(
					mock.Anything,
					mock.MatchedBy(func(wi *dbModels.WidgetInstance) bool {
						return wi.ID == widgetID && wi.WidgetType == "line_chart"
					}),
				).Return(nil, errors.New("update failed"))
			},
			wantErr: true,
		},
		{
			name: "success - update sheet ID",
			input: models.WidgetInstance{
				ID:      widgetID,
				SheetID: sheetID,
			},
			want: &models.WidgetInstance{
				ID:           widgetID,
				Title:        "Test Widget",
				WidgetType:   "bar_chart",
				DataMappings: defaultDataMappings,
				SheetID:      sheetID,
			},
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(existingWidget, nil)
				m.EXPECT().UpdateWidgetInstance(
					mock.Anything,
					mock.MatchedBy(func(wi *dbModels.WidgetInstance) bool {
						return wi.ID == widgetID && wi.SheetID == sheetID
					}),
				).Return(&dbModels.WidgetInstance{
					ID:           widgetID,
					Title:        "Test Widget",
					WidgetType:   "bar_chart",
					DataMappings: dataMappingsJSON,
					SheetID:      sheetID,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - widget instance not found",
			input: models.WidgetInstance{
				ID:         widgetID,
				WidgetType: "line_chart",
			},
			want: nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(dbModels.WidgetInstance{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "error - update failed",
			input: models.WidgetInstance{
				ID:         widgetID,
				WidgetType: "line_chart",
			},
			want: nil,
			mockSetup: func(m *mockWidgets.MockWidgetsServiceStore) {
				m.EXPECT().GetWidgetInstanceByID(mock.Anything, widgetID).Return(existingWidget, nil)
				m.EXPECT().UpdateWidgetInstance(
					mock.Anything,
					mock.Anything,
				).Return(nil, errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mockWidgets.NewMockWidgetsServiceStore(t)
			tt.mockSetup(mockStore)

			service := &widgetsService{store: mockStore}
			got, err := service.UpdateWidgetInstance(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
