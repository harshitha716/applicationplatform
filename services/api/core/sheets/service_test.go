package sheets

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"

	datasetsconstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	sheetmodels "github.com/Zampfi/application-platform/services/api/core/sheets/models"
	"github.com/Zampfi/application-platform/services/api/db/models"
	mock_datasetsService "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mockcache "github.com/Zampfi/application-platform/services/api/mocks/pkg/cache"
	querybuilderconstants "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(t *testing.T) (SheetsService, *mock_store.MockStore, context.Context) {
	mockStore := mock_store.NewMockStore(t)
	mockDatasetService := mock_datasetsService.NewMockDatasetService(t)
	mockCacheService := mockcache.NewMockCacheClient(t)
	service := NewSheetsService(mockStore, mockDatasetService, mockCacheService)

	// Create a context with logger
	ctx := context.Background()

	return service, mockStore, ctx
}

func TestNewSheetsService(t *testing.T) {
	mockStore := mock_store.NewMockStore(t)
	mockDatasetService := mock_datasetsService.NewMockDatasetService(t)
	mockCacheService := mockcache.NewMockCacheClient(t)
	service := NewSheetsService(mockStore, mockDatasetService, mockCacheService)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
}

func TestStructImplementsInterface(t *testing.T) {
	var _ SheetsService = &sheetsService{}
}

func TestGetSheetsAll(t *testing.T) {

	pageId := uuid.New()

	sheetId := uuid.New()
	tests := []struct {
		name           string
		setupMock      func(*mock_store.MockStore)
		expectedSheets []models.Sheet
		expectedErr    error
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetsAll(mock.Anything, models.SheetFilters{PageIds: []uuid.UUID{pageId}, SortParams: []models.SheetSortParams{{Column: "created_at", Desc: false}}}).
					Return([]models.Sheet{{ID: sheetId}}, nil)
			},
			expectedSheets: []models.Sheet{{ID: sheetId}},
			expectedErr:    nil,
		},
		{
			name: "store error",
			setupMock: func(m *mock_store.MockStore) {
				m.On("GetSheetsAll", mock.Anything, models.SheetFilters{PageIds: []uuid.UUID{pageId}, SortParams: []models.SheetSortParams{{Column: "created_at", Desc: false}}}).
					Return(nil, errors.New("store error"))
			},
			expectedSheets: nil,
			expectedErr:    errors.New("store error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore, ctx := setupTest(t)

			// Setup mock expectations
			tt.setupMock(mockStore)

			// Execute
			sheets, err := service.GetSheetsByPageId(ctx, pageId)

			// Verify
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, sheets, len(tt.expectedSheets))
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestGetSheetByID(t *testing.T) {
	sheetID := uuid.New()

	widgetInstanceId := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(*mock_store.MockStore)
		expectedSheet *models.Sheet
		expectedErr   error
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetById(mock.Anything, sheetID).
					Return(&models.Sheet{ID: sheetID, WidgetInstances: []models.WidgetInstance{{ID: widgetInstanceId}}}, nil)
			},
			expectedSheet: &models.Sheet{ID: sheetID, WidgetInstances: []models.WidgetInstance{{ID: widgetInstanceId}}},
			expectedErr:   nil,
		},
		{
			name: "store error",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetById(mock.Anything, sheetID).
					Return(nil, errors.New("store error"))
			},
			expectedSheet: nil,
			expectedErr:   errors.New("store error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore, ctx := setupTest(t)

			// Setup mock expectations
			tt.setupMock(mockStore)

			// Execute
			sheet, err := service.GetSheetById(ctx, sheetID)

			// Verify
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSheet.ID, sheet.ID)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestIsRangeOperator(t *testing.T) {
	tests := []struct {
		name           string
		operator       string
		expectedResult bool
	}{
		{
			name:           "greater than operator",
			operator:       querybuilderconstants.GreaterThanOperator,
			expectedResult: true,
		},
		{
			name:           "equal operator",
			operator:       querybuilderconstants.EqualOperator,
			expectedResult: true,
		},
		{
			name:           "invalid operator",
			operator:       "invalid",
			expectedResult: false,
		},
		{
			name:           "in operator",
			operator:       querybuilderconstants.InOperator,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRangeOperator(tt.operator)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetSheetFilterConfigByIdFromDB(t *testing.T) {
	orgId := uuid.MustParse("f4149aae-7c15-450c-a5a9-da358955a22a")
	sheetId := uuid.New()
	datasetId1 := uuid.New()
	datasetId2 := uuid.New()

	sheetModelSingleTarget := sheetmodels.Sheet{
		ID: sheetId,
		SheetConfig: sheetmodels.SheetConfig{
			Version: "1",
			NativeFilterConfig: []sheetmodels.NativeFilterConfig{
				{
					Name:           "Test Filter",
					Id:             "test_filter",
					FilterType:     "multi-select",
					DataType:       "string",
					WidgetsInScope: []string{},
					Targets: []sheetmodels.FilterTarget{
						{DatasetId: datasetId1, Column: "test_column"},
					},
				},
			},
		},
	}

	sheetModelMultipleTargets := sheetmodels.Sheet{
		ID: sheetId,
		SheetConfig: sheetmodels.SheetConfig{
			Version: "1",
			NativeFilterConfig: []sheetmodels.NativeFilterConfig{
				{
					Name:           "Test Filter",
					Id:             "test_filter",
					FilterType:     "multi-select",
					WidgetsInScope: []string{},
					Targets: []sheetmodels.FilterTarget{
						{DatasetId: datasetId1, Column: "test_column1"},
						{DatasetId: datasetId2, Column: "test_column2"},
					},
				},
			},
		},
	}

	tests := []struct {
		name           string
		sheetModel     sheetmodels.Sheet
		setupMocks     func(*mock_store.MockStore, *mock_datasetsService.MockDatasetService)
		expectedConfig *sheetmodels.FilterOptionsConfig
		expectedErr    error
	}{
		{
			name:       "successful retrieval",
			sheetModel: sheetModelSingleTarget,
			setupMocks: func(ms *mock_store.MockStore, mds *mock_datasetsService.MockDatasetService) {

				mds.EXPECT().GetOptionsForColumn(
					mock.Anything,
					orgId,
					datasetId1.String(),
					"test_column",
					datasetsconstants.FilterTypeMultiSearch,
					false,
				).Return([]interface{}{"option1", "option2"}, nil)
			},
			expectedConfig: &sheetmodels.FilterOptionsConfig{
				NativeFilterConfig: []sheetmodels.FilterOptionsModel{
					{
						Name:           "Test Filter",
						FilterType:     datasetsconstants.FilterTypeMultiSearch,
						DataType:       "",
						WidgetsInScope: nil,
						Targets: []sheetmodels.FilterTarget{
							{
								DatasetId: datasetId1,
								Column:    "test_column",
							},
						},
						Options:      []interface{}{"option1", "option2"},
						DefaultValue: nil,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:       "deduplicates options from multiple targets",
			sheetModel: sheetModelMultipleTargets,
			setupMocks: func(ms *mock_store.MockStore, mds *mock_datasetsService.MockDatasetService) {

				// First target returns some options
				mds.EXPECT().GetOptionsForColumn(
					mock.Anything,
					orgId,
					datasetId1.String(),
					"test_column1",
					datasetsconstants.FilterTypeMultiSearch,
					false,
				).Return([]interface{}{"option1", "option2", "option3"}, nil)

				// Second target returns overlapping options
				mds.EXPECT().GetOptionsForColumn(
					mock.Anything,
					orgId,
					datasetId2.String(),
					"test_column2",
					datasetsconstants.FilterTypeMultiSearch,
					false,
				).Return([]interface{}{"option2", "option3", "option4"}, nil)
			},
			expectedConfig: &sheetmodels.FilterOptionsConfig{
				NativeFilterConfig: []sheetmodels.FilterOptionsModel{
					{
						Name:           "Test Filter",
						FilterType:     datasetsconstants.FilterTypeMultiSearch,
						DataType:       "",
						WidgetsInScope: nil,
						Targets: []sheetmodels.FilterTarget{
							{
								DatasetId: datasetId1,
								Column:    "test_column1",
							},
							{
								DatasetId: datasetId2,
								Column:    "test_column2",
							},
						},
						// The order may vary as we're using a map, so we need to check this differently in the test
						Options:      []interface{}{"option1", "option2", "option3", "option4"},
						DefaultValue: nil,
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mock_store.NewMockStore(t)
			mockDatasetService := mock_datasetsService.NewMockDatasetService(t)
			mockCacheService := mockcache.NewMockCacheClient(t)
			service := NewSheetsService(mockStore, mockDatasetService, mockCacheService)
			ctx := context.Background()

			// Setup mock expectations
			tt.setupMocks(mockStore, mockDatasetService)

			// Execute
			config, err := service.getSheetFilterConfigFromDB(ctx, orgId, tt.sheetModel)

			// Verify
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				// Verify everything except options ordering
				assert.Equal(t, len(tt.expectedConfig.NativeFilterConfig), len(config.NativeFilterConfig))
				assert.Equal(t, tt.expectedConfig.NativeFilterConfig[0].Name, config.NativeFilterConfig[0].Name)
				assert.Equal(t, tt.expectedConfig.NativeFilterConfig[0].FilterType, config.NativeFilterConfig[0].FilterType)
				assert.Equal(t, tt.expectedConfig.NativeFilterConfig[0].Targets, config.NativeFilterConfig[0].Targets)

				// Check options regardless of order
				expectedOptions := map[interface{}]bool{}
				for _, opt := range tt.expectedConfig.NativeFilterConfig[0].Options {
					expectedOptions[opt] = true
				}

				for _, opt := range config.NativeFilterConfig[0].Options {
					assert.True(t, expectedOptions[opt], "Option %v should be in the result", opt)
				}

				assert.Equal(t, len(expectedOptions), len(config.NativeFilterConfig[0].Options))
			}

			mock.AssertExpectationsForObjects(t, mockStore, mockDatasetService)
		})
	}
}

func TestCreateSheet(t *testing.T) {
	t.Parallel()

	sheetID := uuid.New()
	pageId := uuid.New()

	testSheetDB := models.Sheet{
		ID:     sheetID,
		PageId: pageId,
		Name:   "Test Sheet",
	}

	tests := []struct {
		name      string
		input     sheetmodels.Sheet
		want      *models.Sheet
		mockSetup func(*mock_store.MockStore)
		wantErr   bool
	}{
		{
			name: "success",
			input: sheetmodels.Sheet{
				ID:     sheetID,
				PageId: pageId,
				Name:   "Test Sheet",
			},
			want: &models.Sheet{
				ID:     sheetID,
				PageId: pageId,
				Name:   "Test Sheet",
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().CreateSheet(mock.Anything, mock.Anything).Return(&testSheetDB, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			input: sheetmodels.Sheet{
				ID:     sheetID,
				PageId: pageId,
				Name:   "Test Sheet",
			},
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().CreateSheet(mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &sheetsService{store: mockStore}
			got, err := service.CreateSheet(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateSheet(t *testing.T) {
	t.Parallel()

	sheetID := uuid.New()
	pageId := uuid.New()
	newPageId := uuid.New()

	existingSheet := &models.Sheet{
		ID:          sheetID,
		PageId:      pageId,
		Name:        "Original Sheet",
		Description: ptr("Original description"),
		SheetConfig: json.RawMessage(`{"version": "1.0"}`),
	}

	tests := []struct {
		name      string
		input     *sheetmodels.Sheet
		want      *models.Sheet
		mockSetup func(*mock_store.MockStore)
		wantErr   bool
	}{
		{
			name: "successful update all fields",
			input: &sheetmodels.Sheet{
				ID:          sheetID,
				PageId:      newPageId,
				Name:        "Updated Sheet",
				Description: ptr("Updated description"),
				SheetConfig: sheetmodels.SheetConfig{Version: "2.0"},
			},
			want: &models.Sheet{
				ID:          sheetID,
				PageId:      newPageId,
				Name:        "Updated Sheet",
				Description: ptr("Updated description"),
				SheetConfig: json.RawMessage(`{"version":"2.0"}`),
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetById(mock.Anything, sheetID).Return(existingSheet, nil)
				m.EXPECT().UpdateSheet(mock.Anything, mock.Anything).Return(&models.Sheet{
					ID:          sheetID,
					PageId:      newPageId,
					Name:        "Updated Sheet",
					Description: ptr("Updated description"),
					SheetConfig: json.RawMessage(`{"version":"2.0"}`),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "partial update",
			input: &sheetmodels.Sheet{
				ID:   sheetID,
				Name: "Updated Sheet",
			},
			want: &models.Sheet{
				ID:          sheetID,
				PageId:      pageId,
				Name:        "Updated Sheet",
				Description: ptr("Original description"),
				SheetConfig: json.RawMessage(`{"version": "1.0"}`),
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetById(mock.Anything, sheetID).Return(existingSheet, nil)
				m.EXPECT().UpdateSheet(mock.Anything, mock.Anything).Return(&models.Sheet{
					ID:          sheetID,
					PageId:      pageId,
					Name:        "Updated Sheet",
					Description: ptr("Original description"),
					SheetConfig: json.RawMessage(`{"version": "1.0"}`),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "sheet not found",
			input: &sheetmodels.Sheet{
				ID:   sheetID,
				Name: "Updated Sheet",
			},
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetById(mock.Anything, sheetID).Return(nil, errors.New("sheet not found"))
			},
			wantErr: true,
		},
		{
			name: "update store error",
			input: &sheetmodels.Sheet{
				ID:   sheetID,
				Name: "Updated Sheet",
			},
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetSheetById(mock.Anything, sheetID).Return(existingSheet, nil)
				m.EXPECT().UpdateSheet(mock.Anything, mock.Anything).Return(nil, errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &sheetsService{store: mockStore}
			got, err := service.UpdateSheet(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetSheetFilterConfigById(t *testing.T) {
	t.Parallel()

	// Fixed test data
	orgId := uuid.MustParse("f4149aae-7c15-450c-a5a9-da358955a22a")
	sheetId := uuid.MustParse("b5259aae-8c15-450c-a5a9-da358955a33b")
	datasetId := uuid.MustParse("c6369aae-9c15-450c-a5a9-da358955a44c")

	// Pre-sort the options to ensure consistent comparison
	sortedOptions := []interface{}{"option1", "option2"}
	sort.Slice(sortedOptions, func(i, j int) bool {
		return sortedOptions[i].(string) < sortedOptions[j].(string)
	})

	validSheetConfig := &sheetmodels.FilterOptionsConfig{
		NativeFilterConfig: []sheetmodels.FilterOptionsModel{
			{
				Name:       "Test Filter",
				FilterType: "multi-select",
				Targets: []sheetmodels.FilterTarget{
					{
						DatasetId: datasetId,
						Column:    "test_column",
					},
				},
				Options: sortedOptions,
			},
		},
	}

	sheetConfigJSON := `{
		"native_filter_config": [{
			"name": "Test Filter",
			"filter_type": "multi-select",
			"targets": [{
				"dataset_id": "` + datasetId.String() + `",
				"column": "test_column"
			}]
		}]
	}`

	cacheKey := "sheet_filter_config:" + sheetId.String()

	tests := []struct {
		name       string
		setupMocks func(t *testing.T, ms *mock_store.MockStore, mds *mock_datasetsService.MockDatasetService, mc *mockcache.MockCacheClient)
		want       *sheetmodels.FilterOptionsConfig
		wantErr    bool
	}{
		{
			name: "successful cache hit",
			setupMocks: func(t *testing.T, ms *mock_store.MockStore, mds *mock_datasetsService.MockDatasetService, mc *mockcache.MockCacheClient) {
				ms.EXPECT().
					GetSheetById(mock.Anything, sheetId).
					Return(&models.Sheet{
						ID:          sheetId,
						SheetConfig: json.RawMessage(`{}`),
					}, nil).
					Once()

				mc.EXPECT().
					FormatKey("sheet_filter_config", sheetId.String()).
					Return(cacheKey, nil).
					Once()

				mc.EXPECT().
					Get(mock.Anything, cacheKey, mock.MatchedBy(func(v interface{}) bool {
						config, ok := v.(*sheetmodels.FilterOptionsConfig)
						if !ok {
							return false
						}
						*config = *validSheetConfig
						return true
					})).
					Return(nil).
					Once()
			},
			want:    validSheetConfig,
			wantErr: false,
		},
		{
			name: "cache miss, successful DB retrieval",
			setupMocks: func(t *testing.T, ms *mock_store.MockStore, mds *mock_datasetsService.MockDatasetService, mc *mockcache.MockCacheClient) {
				ms.EXPECT().
					GetSheetById(mock.Anything, sheetId).
					Return(&models.Sheet{
						ID:          sheetId,
						SheetConfig: json.RawMessage(sheetConfigJSON),
					}, nil).
					Once()

				mc.EXPECT().
					FormatKey("sheet_filter_config", sheetId.String()).
					Return(cacheKey, nil).
					Once()

				mc.EXPECT().
					Get(mock.Anything, cacheKey, mock.Anything).
					Return(errors.New("cache miss")).
					Once()

				mds.EXPECT().
					GetOptionsForColumn(
						mock.Anything,
						orgId,
						datasetId.String(),
						"test_column",
						"multi-select",
						false,
					).
					Return(sortedOptions, nil).
					Once()

				mc.EXPECT().
					Set(
						mock.Anything,
						cacheKey,
						mock.MatchedBy(func(v interface{}) bool {
							config, ok := v.(*sheetmodels.FilterOptionsConfig)
							if !ok {
								return false
							}
							// Sort the options before comparison
							for i := range config.NativeFilterConfig {
								opts := config.NativeFilterConfig[i].Options
								sort.Slice(opts, func(i, j int) bool {
									return opts[i].(string) < opts[j].(string)
								})
							}
							return reflect.DeepEqual(config, validSheetConfig)
						}),
						time.Minute*30,
					).
					Return(nil).
					Once()
			},
			want:    validSheetConfig,
			wantErr: false,
		},
		{
			name: "initial sheet retrieval error",
			setupMocks: func(t *testing.T, ms *mock_store.MockStore, mds *mock_datasetsService.MockDatasetService, mc *mockcache.MockCacheClient) {
				ms.EXPECT().
					GetSheetById(mock.Anything, sheetId).
					Return(nil, errors.New("db error")).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			mockDatasetService := mock_datasetsService.NewMockDatasetService(t)
			mockCacheService := mockcache.NewMockCacheClient(t)

			tt.setupMocks(t, mockStore, mockDatasetService, mockCacheService)

			service := NewSheetsService(mockStore, mockDatasetService, mockCacheService)

			got, err := service.GetSheetFilterConfigById(context.Background(), orgId, sheetId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Sort the options in the result before comparison
			if got != nil {
				for i := range got.NativeFilterConfig {
					opts := got.NativeFilterConfig[i].Options
					sort.Slice(opts, func(i, j int) bool {
						return opts[i].(string) < opts[j].(string)
					})
				}
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper function to create string pointer
func ptr(s string) *string {
	return &s
}
