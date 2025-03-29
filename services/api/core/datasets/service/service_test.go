package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	dataplatformactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	dataplatformactionmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	dataplatformConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	dataplatformDataTypesConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	servicemodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	datasetConstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	datasetErrors "github.com/Zampfi/application-platform/services/api/core/datasets/errors"
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	rulemodels "github.com/Zampfi/application-platform/services/api/core/rules/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	cloudservicemodels "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"
	"github.com/Zampfi/application-platform/services/api/pkg/s3"
	temporalmodels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"go.uber.org/zap"

	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	mockDataplatform "github.com/Zampfi/application-platform/services/api/mocks/core/dataplatform"
	mockDatasetService "github.com/Zampfi/application-platform/services/api/mocks/core/datasets/service"
	mockfileimports "github.com/Zampfi/application-platform/services/api/mocks/core/fileimports"
	mockruleservice "github.com/Zampfi/application-platform/services/api/mocks/core/rules/service"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mock_cache "github.com/Zampfi/application-platform/services/api/mocks/pkg/cache"
	mock_cloudservice "github.com/Zampfi/application-platform/services/api/mocks/pkg/cloudservices/service"
	mock_querybuilder "github.com/Zampfi/application-platform/services/api/mocks/pkg/querybuilder/service"
	mock_s3 "github.com/Zampfi/application-platform/services/api/mocks/pkg/s3"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	querybuildermodels "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
	mock_temporal "github.com/Zampfi/workflow-sdk-go/mocks/workflowmanagers/temporal"

	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetFilterConfigByDatasetId(t *testing.T) {
	tests := []struct {
		name           string
		merchantId     uuid.UUID
		datasetId      string
		mockSetup      func(*mockDataplatform.MockDataPlatformService, *mock_querybuilder.MockQueryBuilder, *mockDatasetService.MockDatasetServiceStore, *mockruleservice.MockRuleService, serverconfig.DatasetConfig, *mock_cache.MockCacheClient)
		expectedError  bool
		expectedConfig []models.FilterConfig
	}{
		{
			name:       "Success case with cache miss",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, serverConfig serverconfig.DatasetConfig, cache *mock_cache.MockCacheClient) {
				cache.EXPECT().FormatKey("dataset_filter_config", "dataset1").Return("dataset_filter_config:dataset1", nil)
				cache.EXPECT().Get(mock.Anything, "dataset_filter_config:dataset1", mock.Anything).Return(errors.New("error"))
				ds.EXPECT().GetDatasetById(mock.Anything, "dataset1").Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}, "display_config": [{"column": "column1", "is_hidden": false, "is_editable": false}, {"column": "column2", "is_hidden": false, "is_editable": false}, {"column": "column3", "is_hidden": false, "is_editable": false}, {"column": "column4", "is_hidden": false, "is_editable": true, "type": "tags"}]}`),
				}, nil)
				m.EXPECT().GetDatasetMetadata(mock.Anything, "123e4567-e89b-12d3-a456-426614174000", "dataset1").Return(
					dataplatformDataModels.DatasetMetadata{
						Schema: map[string]dataplatformDataModels.ColumnMetadata{
							"column1": {Type: "string"},
							"column2": {Type: "integer"},
							"column3": {Type: "date"},
							"column4": {Type: "boolean"},
						},
						Stats: dataplatformDataModels.DatasetStats{
							ColumnStats: map[string]dataplatformDataModels.ColumnStats{
								"column1": {DistinctCount: 10},
								"column4": {DistinctCount: 2},
							},
						},
					}, nil)
				m.EXPECT().GetDags(mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(map[string]*servicemodels.DAGNode{
					"dataset1": {
						NodeId:   "dataset1",
						NodeType: servicemodels.NodeTypeDataset,
						Parents: []*servicemodels.DAGNode{
							{
								NodeId:   "folder1",
								NodeType: servicemodels.NodeTypeFolder,
							},
						},
					},
				}, nil)
				m.EXPECT().QueryRealTime(
					mock.Anything,
					"123e4567-e89b-12d3-a456-426614174000",
					"SELECT DISTINCT column1 FROM {{.zamp_dataset1}} where _zamp_is_deleted = False LIMIT 20",
					map[string]string{"zamp_dataset1": "dataset1"},
				).Return(dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"column1": "1"},
						{"column1": "2"},
						{"column1": "3"},
						{"column1": "4"},
						{"column1": "5"},
						{"column1": "6"},
						{"column1": "7"},
						{"column1": "8"},
						{"column1": "9"},
						{"column1": "11"},
					},
				}, nil)
				m.EXPECT().QueryRealTime(
					mock.Anything,
					"123e4567-e89b-12d3-a456-426614174000",
					"SELECT DISTINCT column4 FROM {{.zamp_dataset1}} where _zamp_is_deleted = False LIMIT 20",
					map[string]string{"zamp_dataset1": "dataset1"},
				).Return(dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"column4": true},
						{"column4": false},
					},
				}, nil)
				cache.EXPECT().Set(mock.Anything, "dataset_filter_config:dataset1", mock.Anything, time.Minute*10).Return(errors.New("error"))
			},
			expectedError: false,
			expectedConfig: []models.FilterConfig{
				{
					Column: "column1",
					Type:   "multi-select",
					DataType: func() *dataplatformDataTypesConstants.Datatype {
						dt := dataplatformDataTypesConstants.StringDataType
						return &dt
					}(),
					Options:  []interface{}{"1", "2", "3", "4", "5", "6", "7", "8", "9", "11"},
					Metadata: map[string]interface{}{},
				},
				{
					Column: "column2",
					Type:   "amount-range",
					DataType: func() *dataplatformDataTypesConstants.Datatype {
						dt := dataplatformDataTypesConstants.IntegerDataType
						return &dt
					}(),
					Options:  []interface{}{},
					Metadata: map[string]interface{}{},
				},
				{
					Column: "column3",
					Type:   "date-range",
					DataType: func() *dataplatformDataTypesConstants.Datatype {
						dt := dataplatformDataTypesConstants.DateDataType
						return &dt
					}(),
					Options:  []interface{}{},
					Metadata: map[string]interface{}{},
				},
				{
					Column: "column4",
					Type:   "multi-select",
					DataType: func() *dataplatformDataTypesConstants.Datatype {
						dt := dataplatformDataTypesConstants.BooleanDataType
						return &dt
					}(),
					Options: []interface{}{true, false},
					Metadata: map[string]interface{}{
						"custom_type": func() *string {
							s := "tags"
							return &s
						}(),
						"is_editable": true,
					},
				},
			},
		},
		{
			name:       "Success case with cache hit",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, serverConfig serverconfig.DatasetConfig, cache *mock_cache.MockCacheClient) {
				cache.EXPECT().FormatKey("dataset_filter_config", "dataset1").Return("dataset_filter_config:dataset1", nil)
				cache.EXPECT().Get(mock.Anything, "dataset_filter_config:dataset1", mock.Anything).Return(nil)
			},
			expectedError:  false,
			expectedConfig: []models.FilterConfig{},
		},
		{
			name:       "Error getting dataset meta info",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, serverConfig serverconfig.DatasetConfig, cache *mock_cache.MockCacheClient) {
				cache.EXPECT().FormatKey("dataset_filter_config", "dataset1").Return("dataset_filter_config:dataset1", nil)
				cache.EXPECT().Get(mock.Anything, "dataset_filter_config:dataset1", mock.Anything).Return(errors.New("error"))
				ds.EXPECT().GetDatasetById(mock.Anything, "dataset1").Return(nil, fmt.Errorf("failed to get dataset meta info"))
			},
			expectedError:  true,
			expectedConfig: nil,
		},
		{
			name:       "Error getting dataset metadata",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, serverConfig serverconfig.DatasetConfig, cache *mock_cache.MockCacheClient) {
				cache.EXPECT().FormatKey("dataset_filter_config", "dataset1").Return("dataset_filter_config:dataset1", nil)
				cache.EXPECT().Get(mock.Anything, "dataset_filter_config:dataset1", mock.Anything).Return(errors.New("error"))
				ds.EXPECT().GetDatasetById(mock.Anything, "dataset1").Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}, "display_config": [{"column": "column1", "is_hidden": false, "is_editable": false}, {"column": "column2", "is_hidden": false, "is_editable": false}, {"column": "column3", "is_hidden": false, "is_editable": false}, {"column": "column4", "is_hidden": false, "is_editable": true}]}`),
				}, nil)
				m.EXPECT().GetDatasetMetadata(mock.Anything, "123e4567-e89b-12d3-a456-426614174000", "dataset1").
					Return(dataplatformDataModels.DatasetMetadata{}, fmt.Errorf("failed to get metadata"))
			},
			expectedError:  true,
			expectedConfig: nil,
		},
		{
			name:       "Error getting multi-select options",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, serverConfig serverconfig.DatasetConfig, cache *mock_cache.MockCacheClient) {
				cache.EXPECT().FormatKey("dataset_filter_config", "dataset1").Return("dataset_filter_config:dataset1", nil)
				cache.EXPECT().Get(mock.Anything, "dataset_filter_config:dataset1", mock.Anything).Return(errors.New("error"))
				ds.EXPECT().GetDatasetById(mock.Anything, "dataset1").Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}, "display_config": [{"column": "column1", "is_hidden": false, "is_editable": false}, {"column": "column2", "is_hidden": false, "is_editable": false}, {"column": "column3", "is_hidden": false, "is_editable": false}, {"column": "column4", "is_hidden": false, "is_editable": true}]}`),
				}, nil)
				m.EXPECT().GetDatasetMetadata(mock.Anything, "123e4567-e89b-12d3-a456-426614174000", "dataset1").Return(
					dataplatformDataModels.DatasetMetadata{
						Schema: map[string]dataplatformDataModels.ColumnMetadata{
							"column1": {Type: "string"},
						},
						Stats: dataplatformDataModels.DatasetStats{
							ColumnStats: map[string]dataplatformDataModels.ColumnStats{
								"column1": {DistinctCount: 2},
							},
						},
					}, nil)
				m.EXPECT().GetDags(mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(map[string]*servicemodels.DAGNode{
					"dataset1": {
						NodeId:   "dataset1",
						NodeType: servicemodels.NodeTypeDataset,
						Parents: []*servicemodels.DAGNode{
							{
								NodeId:   "folder1",
								NodeType: servicemodels.NodeTypeFolder,
							},
						},
					},
				}, nil)
				m.EXPECT().QueryRealTime(
					mock.Anything,
					"123e4567-e89b-12d3-a456-426614174000",
					mock.AnythingOfType("string"),
					mock.AnythingOfType("map[string]string"),
				).Return(dataplatformmodels.QueryResult{}, fmt.Errorf("query failed"))
			},
			expectedError:  true,
			expectedConfig: nil,
		},
		{
			name:       "Empty schema case",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, serverConfig serverconfig.DatasetConfig, cache *mock_cache.MockCacheClient) {
				cache.EXPECT().FormatKey("dataset_filter_config", "dataset1").Return("dataset_filter_config:dataset1", nil)
				cache.EXPECT().Get(mock.Anything, "dataset_filter_config:dataset1", mock.Anything).Return(errors.New("error"))
				ds.EXPECT().GetDatasetById(mock.Anything, "dataset1").Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}}`),
				}, nil)
				m.EXPECT().GetDags(mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(map[string]*servicemodels.DAGNode{
					"dataset1": {
						NodeId:   "dataset1",
						NodeType: servicemodels.NodeTypeDataset,
						Parents: []*servicemodels.DAGNode{
							{
								NodeId:   "folder1",
								NodeType: servicemodels.NodeTypeFolder,
							},
						},
					},
				}, nil)
				m.EXPECT().GetDatasetMetadata(mock.Anything, "123e4567-e89b-12d3-a456-426614174000", "dataset1").Return(
					dataplatformDataModels.DatasetMetadata{
						Schema: map[string]dataplatformDataModels.ColumnMetadata{},
						Stats:  dataplatformDataModels.DatasetStats{},
					}, nil)
				cache.EXPECT().Set(mock.Anything, "dataset_filter_config:dataset1", mock.Anything, time.Minute*10).Return(errors.New("error"))
			},
			expectedError:  false,
			expectedConfig: []models.FilterConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := mockfileimports.NewMockFileImportService(t)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}
			tt.mockSetup(mockDPS, mockQueryBuilder, mockDS, mockRuleService, serverConfig, mockCacheClient)

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			result, _, err := svc.GetFilterConfigByDatasetId(context.Background(), tt.merchantId, tt.datasetId)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				expectedMap := make(map[string]models.FilterConfig)
				resultMap := make(map[string]models.FilterConfig)

				for _, cfg := range tt.expectedConfig {
					expectedMap[cfg.Column] = cfg
				}
				for _, cfg := range result {
					resultMap[cfg.Column] = cfg
				}

				assert.Equal(t, expectedMap, resultMap, "FilterConfig elements should match regardless of order")
			}

			mockDPS.AssertExpectations(t)
			mockDS.AssertExpectations(t)
		})
	}
}

func TestGetStringColumnFilterType(t *testing.T) {
	s := &datasetService{}

	tests := []struct {
		name          string
		distinctCount int
		expected      string
	}{
		{
			name:          "below threshold",
			distinctCount: datasetConstants.MultiSelectThreshold - 1,
			expected:      datasetConstants.FilterTypeMultiSearch,
		},
		{
			name:          "at threshold",
			distinctCount: datasetConstants.MultiSelectThreshold,
			expected:      datasetConstants.FilterTypeMultiSearch,
		},
		{
			name:          "above threshold",
			distinctCount: datasetConstants.MultiSelectThreshold + 1,
			expected:      datasetConstants.FilterTypeSearch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.getStringColumnFilterType(tt.distinctCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDistinctValueCount(t *testing.T) {
	s := &datasetService{}

	tests := []struct {
		name        string
		datasetInfo dataplatformDataModels.DatasetMetadata
		colName     string
		expected    int
	}{
		{
			name: "column exists",
			datasetInfo: dataplatformDataModels.DatasetMetadata{
				Stats: dataplatformDataModels.DatasetStats{
					ColumnStats: map[string]dataplatformDataModels.ColumnStats{
						"col1": {DistinctCount: 5},
					},
				},
			},
			colName:  "col1",
			expected: 5,
		},
		{
			name: "column does not exist",
			datasetInfo: dataplatformDataModels.DatasetMetadata{
				Stats: dataplatformDataModels.DatasetStats{
					ColumnStats: map[string]dataplatformDataModels.ColumnStats{},
				},
			},
			colName:  "nonexistent",
			expected: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.getDistinctValueCount(tt.datasetInfo, tt.colName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertToFilterConfig(t *testing.T) {
	s := &datasetService{}

	tests := []struct {
		name            string
		datasetInfo     dataplatformDataModels.DatasetMetadata
		datasetMetaInfo models.DatasetMetadataConfig
		expected        []models.FilterConfig
	}{
		{
			name: "various data types",
			datasetInfo: dataplatformDataModels.DatasetMetadata{
				Schema: map[string]dataplatformDataModels.ColumnMetadata{
					"date_col":   {Type: string(dataplatformDataTypesConstants.DateDataType)},
					"int_col":    {Type: string(dataplatformDataTypesConstants.IntegerDataType)},
					"string_col": {Type: string(dataplatformDataTypesConstants.StringDataType)},
					"bool_col":   {Type: string(dataplatformDataTypesConstants.BooleanDataType)},
				},
				Stats: dataplatformDataModels.DatasetStats{
					ColumnStats: map[string]dataplatformDataModels.ColumnStats{
						"string_col": {DistinctCount: datasetConstants.MultiSelectThreshold + 1},
					},
				},
			},
			datasetMetaInfo: models.DatasetMetadataConfig{
				DatasetConfig: dataplatformDataModels.DatasetConfig{
					Columns: map[string]dataplatformDataModels.DatasetColumnConfig{
						"string_col": {CustomType: dataplatformConstants.DatabricksColumnCustomTypeTags},
					},
				},
				DisplayConfig: []models.DisplayConfig{
					{Column: "string_col", IsHidden: false, IsEditable: true, Type: func() *string {
						s := "tags"
						return &s
					}()},
					{Column: "int_col", IsHidden: false, IsEditable: false},
					{Column: "date_col", IsHidden: false, IsEditable: false},
					{Column: "bool_col", IsHidden: false, IsEditable: true},
				},
			},
			expected: []models.FilterConfig{
				{Column: "date_col", Type: datasetConstants.FilterTypeDateRange, DataType: func() *dataplatformDataTypesConstants.Datatype {
					dt := dataplatformDataTypesConstants.DateDataType
					return &dt
				}(), Options: []interface{}{}, Metadata: map[string]interface{}{}},
				{Column: "int_col", Type: datasetConstants.FilterTypeAmountRange, DataType: func() *dataplatformDataTypesConstants.Datatype {
					dt := dataplatformDataTypesConstants.IntegerDataType
					return &dt
				}(), Options: []interface{}{}, Metadata: map[string]interface{}{}},
				{Column: "string_col", Type: datasetConstants.FilterTypeSearch, DataType: func() *dataplatformDataTypesConstants.Datatype {
					dt := dataplatformDataTypesConstants.StringDataType
					return &dt
				}(), Options: []interface{}{}, Metadata: map[string]interface{}{"custom_type": func() *string {
					s := "tags"
					return &s
				}(), "is_editable": true}},
				{Column: "bool_col", Type: datasetConstants.FilterTypeMultiSearch, DataType: func() *dataplatformDataTypesConstants.Datatype {
					dt := dataplatformDataTypesConstants.BooleanDataType
					return &dt
				}(), Options: []interface{}{}, Metadata: map[string]interface{}{"is_editable": true}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.convertToFilterConfig(tt.datasetInfo, tt.datasetMetaInfo)
			assert.Equal(t, len(tt.expected), len(result), "FilterConfig slices should have same length")

			expectedMap := make(map[string]models.FilterConfig)
			resultMap := make(map[string]models.FilterConfig)

			for _, cfg := range tt.expected {
				expectedMap[cfg.Column] = cfg
			}
			for _, cfg := range result {
				resultMap[cfg.Column] = cfg
			}

			assert.Equal(t, expectedMap, resultMap, "FilterConfig elements should match regardless of order")
		})
	}
}

func TestGetOptionsForColumn(t *testing.T) {
	tests := []struct {
		name          string
		merchantId    uuid.UUID
		datasetId     string
		column        string
		filterType    string
		mockSetup     func(*mockDataplatform.MockDataPlatformService)
		expectedError bool
		expected      []interface{}
	}{
		{
			name:       "Success case - multi-search filter",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			column:     "status",
			filterType: datasetConstants.FilterTypeMultiSearch,
			mockSetup: func(m *mockDataplatform.MockDataPlatformService) {
				m.EXPECT().QueryRealTime(
					mock.Anything,
					"123e4567-e89b-12d3-a456-426614174000",
					"SELECT DISTINCT status FROM {{.zamp_dataset1}} where _zamp_is_deleted = False LIMIT 20",
					map[string]string{"zamp_dataset1": "dataset1"},
				).Return(dataplatformmodels.QueryResult{
					Rows: []map[string]interface{}{
						{"status": "active"},
						{"status": "inactive"},
						{"status": "pending"},
					},
				}, nil)
			},
			expectedError: false,
			expected:      []interface{}{"active", "inactive", "pending"},
		},
		{
			name:       "Error case - query fails",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			column:     "status",
			filterType: datasetConstants.FilterTypeMultiSearch,
			mockSetup: func(m *mockDataplatform.MockDataPlatformService) {
				m.EXPECT().QueryRealTime(
					mock.Anything,
					"123e4567-e89b-12d3-a456-426614174000",
					"SELECT DISTINCT status FROM {{.zamp_dataset1}} where _zamp_is_deleted = False LIMIT 20",
					map[string]string{"zamp_dataset1": "dataset1"},
				).Return(dataplatformmodels.QueryResult{}, fmt.Errorf("database error"))
			},
			expectedError: true,
			expected:      nil,
		},
		{
			name:       "Default case - non-multi-search filter",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "dataset1",
			column:     "status",
			filterType: "unknown-filter-type",
			mockSetup:  func(m *mockDataplatform.MockDataPlatformService) {},
			expected:   []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDPS)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			result, err := svc.GetOptionsForColumn(
				context.Background(),
				tt.merchantId,
				tt.datasetId,
				tt.column,
				tt.filterType,
				true,
			)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockDPS.AssertExpectations(t)
		})
	}
}

func TestDatasetAudiences(t *testing.T) {
	t.Parallel()

	datasetId := uuid.New()
	audience1 := uuid.New()
	audience2 := uuid.New()
	testAudiences := []storemodels.ResourceAudiencePolicy{{
		ID:                   uuid.New(),
		ResourceID:           datasetId,
		ResourceType:         "dataset",
		ResourceAudienceType: "user",
		ResourceAudienceID:   audience1,
		User:                 &storemodels.User{ID: audience1, Email: "audience1@user.com", Name: "audience1"},
	}, {
		ID:                   uuid.New(),
		ResourceID:           datasetId,
		ResourceType:         "dataset",
		ResourceAudienceType: "user",
		ResourceAudienceID:   audience2,
		User:                 &storemodels.User{ID: audience2, Email: "audience2@user.com", Name: "audience2"},
	}}

	tests := []struct {
		name      string
		want      []storemodels.ResourceAudiencePolicy
		mockSetup func(*mockDatasetService.MockDatasetServiceStore)
		wantErr   bool
	}{
		{
			name: "success",
			want: testAudiences,
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {
				m.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return(testAudiences, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			want: nil,
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {
				m.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mockDatasetService.NewMockDatasetServiceStore(t)
			ctx := context.Background()
			tt.mockSetup(mockStore)

			// Execute
			service := &datasetService{datasetStore: mockStore}
			got, err := service.GetDatasetAudiences(ctx, datasetId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegisterDataset(t *testing.T) {
	tests := []struct {
		name                string
		merchantId          uuid.UUID
		userId              uuid.UUID
		datasetCreationInfo models.DatasetCreationInfo
		mockSetup           func(*mockDataplatform.MockDataPlatformService, *mockDatasetService.MockDatasetServiceStore)
		expectedError       bool
		expectedActionId    string
		expectedDatasetId   uuid.UUID
	}{
		{
			name:       "Success case",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			userId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),

			datasetCreationInfo: models.DatasetCreationInfo{
				DatasetTitle:       "Test Dataset",
				DatasetDescription: func() *string { s := "Test Description"; return &s }(),
				DatasetType:        storemodels.DatasetTypeSource,
				DatasetConfig: dataplatformDataModels.DatasetConfig{
					Columns: map[string]dataplatformDataModels.DatasetColumnConfig{
						"column1": {CustomType: "type1"},
					},
				},
				Provider: dataplatformdataconstants.ProviderPinot,
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.AnythingOfType("func(store.DatasetStore) error")).Return(nil).
					Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
						ds.EXPECT().CreateDataset(mock.Anything, mock.AnythingOfType("models.Dataset")).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), nil)
						ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, storemodels.AudienceTypeOrganization, mock.Anything, storemodels.PrivilegeDatasetAdmin).Return(&storemodels.ResourceAudiencePolicy{}, nil)
						m.EXPECT().RegisterDataset(mock.Anything, mock.MatchedBy(func(payload servicemodels.RegisterDatasetPayload) bool {
							return payload.MerchantID == "123e4567-e89b-12d3-a456-426614174000"
						})).Return(dataplatformactionmodels.CreateActionResponse{
							ActionID: "action123",
						}, nil)
						fn(ds)
					})
			},
			expectedError:     false,
			expectedActionId:  "action123",
			expectedDatasetId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
		},
		{
			name:       "Failed to create dataset",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			userId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			datasetCreationInfo: models.DatasetCreationInfo{
				DatasetConfig: dataplatformDataModels.DatasetConfig{},
				DatasetType:   storemodels.DatasetTypeSource,
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.AnythingOfType("func(store.DatasetStore) error")).Return(errors.New("failed to create dataset")).
					Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
						ds.EXPECT().CreateDataset(mock.Anything, mock.AnythingOfType("models.Dataset")).Return(uuid.Nil, fmt.Errorf("failed to create dataset"))
						fn(ds)
					})
			},
			expectedError:     true,
			expectedActionId:  "",
			expectedDatasetId: uuid.Nil,
		},
		{
			name:       "Success case - MV Dataset",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			userId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			datasetCreationInfo: models.DatasetCreationInfo{
				DatasetTitle:       "Test MV Dataset",
				DatasetDescription: func() *string { s := "Test MV Description"; return &s }(),
				DatasetType:        storemodels.DatasetTypeMV,
				MVConfig: &models.MVConfig{
					Query:            "SELECT * FROM source_dataset",
					QueryParams:      map[string]string{"param1": "value1"},
					ParentDatasetIds: []string{"parent1", "parent2"},
				},
				DatabricksConfig: dataplatformDataModels.DatabricksConfig{
					DedupColumns:  []string{"column1", "column2"},
					OrderByColumn: "timestamp",
				},
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.AnythingOfType("func(store.DatasetStore) error")).Return(nil).
					Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
						ds.EXPECT().CreateDataset(mock.Anything, mock.AnythingOfType("models.Dataset")).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), nil)
						ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, storemodels.AudienceTypeOrganization, mock.Anything, storemodels.PrivilegeDatasetAdmin).Return(&storemodels.ResourceAudiencePolicy{}, nil)

						m.EXPECT().CreateMV(mock.Anything, mock.MatchedBy(func(payload servicemodels.CreateMVPayload) bool {
							return payload.MerchantID == "123e4567-e89b-12d3-a456-426614174000" &&
								payload.ActorId == "123e4567-e89b-12d3-a456-426614174001" &&
								payload.ActionMetadataPayload.Query == "SELECT * FROM source_dataset" &&
								payload.ActionMetadataPayload.QueryParams["param1"] == "value1" &&
								len(payload.ActionMetadataPayload.ParentDatasetIds) == 2
						})).Return(dataplatformactionmodels.CreateActionResponse{
							ActionID: "action123",
						}, nil)

						fn(ds)
					})
			},
			expectedError:     false,
			expectedActionId:  "action123",
			expectedDatasetId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}
			tt.mockSetup(mockDPS, mockDS)

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			actionId, datasetId, err := svc.RegisterDataset(context.Background(), tt.merchantId, tt.userId, tt.datasetCreationInfo)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedActionId, actionId)
				assert.NotEqual(t, uuid.Nil, datasetId)
			}

			mockDPS.AssertExpectations(t)
			mockDS.AssertExpectations(t)
		})
	}
}

func TestRegisterDatasetJob(t *testing.T) {
	tests := []struct {
		name           string
		merchantId     uuid.UUID
		jobInfo        dataplatformactionmodels.RegisterJobActionPayload
		mockSetup      func(*mockDataplatform.MockDataPlatformService, *mock_querybuilder.MockQueryBuilder, *mockDatasetService.MockDatasetServiceStore, *mockruleservice.MockRuleService)
		expectedError  bool
		expectedAction string
	}{
		{
			name:       "Success case",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			jobInfo: dataplatformactionmodels.RegisterJobActionPayload{
				JobType:          "test_type",
				SourceType:       "test_source",
				SourceValue:      "test_value",
				DestinationType:  "test_dest",
				DestinationValue: "test_dest_value",
				TemplateId:       "template123",
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService) {
				m.EXPECT().RegisterJob(mock.Anything, mock.MatchedBy(func(payload servicemodels.RegisterJobPayload) bool {
					return payload.MerchantID == "123e4567-e89b-12d3-a456-426614174000"
				})).Return(dataplatformactionmodels.CreateActionResponse{
					ActionID: "action123",
				}, nil)
			},
			expectedError:  false,
			expectedAction: "action123",
		},
		{
			name:       "Failed to register job",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			jobInfo:    dataplatformactionmodels.RegisterJobActionPayload{},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService) {
				m.EXPECT().RegisterJob(mock.Anything, mock.Anything).Return(dataplatformactionmodels.CreateActionResponse{}, fmt.Errorf("failed to register job"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			tt.mockSetup(mockDPS, mockQueryBuilder, mockDS, mockRuleService)

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			actionId, err := svc.RegisterDatasetJob(context.Background(), tt.merchantId, tt.jobInfo)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAction, actionId)
			}

			mockDPS.AssertExpectations(t)
		})
	}
}

func TestUpdateDataset(t *testing.T) {
	tests := []struct {
		name           string
		merchantId     uuid.UUID
		datasetId      string
		params         models.UpdateDatasetParams
		mockSetup      func(*mockDataplatform.MockDataPlatformService, *mockDatasetService.MockDatasetServiceStore)
		expectedError  bool
		expectedAction string
	}{
		{
			name:       "Success case - update dataset config",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "123e4567-e89b-12d3-a456-426614174002",
			params: models.UpdateDatasetParams{
				DatasetConfig: &dataplatformDataModels.DatasetConfig{
					Columns: map[string]dataplatformDataModels.DatasetColumnConfig{
						"column1": {CustomType: "updated_type"},
					},
				},
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetById(mock.Anything, "123e4567-e89b-12d3-a456-426614174002").Return(&storemodels.Dataset{
					Metadata: []byte(`{"dataset_config":{},"display_config":[]}`),
				}, nil)

				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.AnythingOfType("func(store.DatasetStore) error")).Return(nil).
					Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
						ds.EXPECT().UpdateDataset(mock.Anything, mock.AnythingOfType("models.Dataset")).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), nil)

						m.EXPECT().UpdateDataset(mock.Anything, mock.MatchedBy(func(payload servicemodels.UpdateDatasetPayload) bool {
							return payload.MerchantID == "123e4567-e89b-12d3-a456-426614174000" &&
								payload.ActionMetadataPayload.EventType == dataplatformactionconstants.UpdateDatasetEventTypeUpdateCustomColumn
						})).Return(dataplatformactionmodels.CreateActionResponse{
							ActionID: "action123",
						}, nil)

						fn(ds)
					})
			},
			expectedError:  false,
			expectedAction: "action123",
		},
		{
			name:       "Success case - update display config only",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "123e4567-e89b-12d3-a456-426614174002",
			params: models.UpdateDatasetParams{
				DisplayConfig: &[]models.DisplayConfig{{
					Column:     "col1",
					IsHidden:   false,
					IsEditable: true,
				}},
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetById(mock.Anything, "123e4567-e89b-12d3-a456-426614174002").Return(&storemodels.Dataset{
					Metadata: []byte(`{"dataset_config":{},"display_config":[]}`),
				}, nil)

				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.AnythingOfType("func(store.DatasetStore) error")).Return(nil).
					Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
						ds.EXPECT().UpdateDataset(mock.Anything, mock.AnythingOfType("models.Dataset")).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), nil)
						fn(ds)
					})
			},
			expectedError:  false,
			expectedAction: "",
		},
		{
			name:       "Success case - update title, desc, type",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "123e4567-e89b-12d3-a456-426614174002",
			params: models.UpdateDatasetParams{
				Title:       func() *string { s := "updated title"; return &s }(),
				Description: func() *string { s := "updated desc"; return &s }(),
				Type:        func() *string { s := "staged"; return &s }(),
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetById(mock.Anything, "123e4567-e89b-12d3-a456-426614174002").Return(&storemodels.Dataset{
					Metadata: []byte(`{"dataset_config":{},"display_config":[]}`),
				}, nil)

				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.AnythingOfType("func(store.DatasetStore) error")).Return(nil).
					Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
						ds.EXPECT().UpdateDataset(mock.Anything, mock.AnythingOfType("models.Dataset")).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), nil)
						fn(ds)
					})
			},
			expectedError:  false,
			expectedAction: "",
		},
		{
			name:       "FAIL case - type ERR_INVALID_DATASET_TYPE",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "123e4567-e89b-12d3-a456-426614174002",
			params: models.UpdateDatasetParams{
				Type: func() *string { s := "unknown"; return &s }(),
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetById(mock.Anything, "123e4567-e89b-12d3-a456-426614174002").Return(&storemodels.Dataset{
					Metadata: []byte(`{"dataset_config":{},"display_config":[]}`),
				}, nil)
			},
			expectedError:  true,
			expectedAction: "",
		},
		{
			name:       "Failed to get dataset",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "123e4567-e89b-12d3-a456-426614174002",
			params:     models.UpdateDatasetParams{},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to get dataset"))
			},
			expectedError:  true,
			expectedAction: "",
		},
		{
			name:       "Failed to unmarshal metadata",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  "123e4567-e89b-12d3-a456-426614174002",
			params:     models.UpdateDatasetParams{},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					Metadata: []byte(`invalid json`),
				}, nil)
			},
			expectedError:  true,
			expectedAction: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}
			tt.mockSetup(mockDPS, mockDS)

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			actionId, err := svc.UpdateDataset(context.Background(), tt.merchantId, tt.datasetId, tt.params)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAction, actionId)
			}

			mockDPS.AssertExpectations(t)
			mockDS.AssertExpectations(t)
		})
	}
}

func TestUpsertTemplate(t *testing.T) {
	tests := []struct {
		name           string
		merchantId     uuid.UUID
		templateConfig dataplatformactionmodels.UpsertTemplateActionPayload
		mockSetup      func(*mockDataplatform.MockDataPlatformService, *mock_querybuilder.MockQueryBuilder, *mockDatasetService.MockDatasetServiceStore, *mockruleservice.MockRuleService, *mockfileimports.MockFileImportService)
		expectedError  bool
		expectedAction string
	}{
		{
			name:       "Success case",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			templateConfig: dataplatformactionmodels.UpsertTemplateActionPayload{
				Id:            "template123",
				Name:          "test_template",
				Configuration: "{\"key\": \"value\"}",
				TemplateType:  "join",
			},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, fis *mockfileimports.MockFileImportService) {
				m.EXPECT().UpsertTemplate(mock.Anything, mock.MatchedBy(func(payload servicemodels.UpsertTemplatePayload) bool {
					return payload.MerchantID == "123e4567-e89b-12d3-a456-426614174000"
				})).Return(dataplatformactionmodels.CreateActionResponse{
					ActionID: "action123",
				}, nil)
			},
			expectedError:  false,
			expectedAction: "action123",
		},
		{
			name:           "Failed to upsert template",
			merchantId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			templateConfig: dataplatformactionmodels.UpsertTemplateActionPayload{},
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, qb *mock_querybuilder.MockQueryBuilder, ds *mockDatasetService.MockDatasetServiceStore, rs *mockruleservice.MockRuleService, fis *mockfileimports.MockFileImportService) {
				m.EXPECT().UpsertTemplate(mock.Anything, mock.Anything).Return(dataplatformactionmodels.CreateActionResponse{}, fmt.Errorf("failed to upsert template"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			tt.mockSetup(mockDPS, mockQueryBuilder, mockDS, mockRuleService, mockFileUploadsService)

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			actionId, err := svc.UpsertTemplate(context.Background(), tt.merchantId, tt.templateConfig)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAction, actionId)
			}

			mockDPS.AssertExpectations(t)
		})
	}
}

func TestAddAudienceToDataset(t *testing.T) {

	currentUserId := uuid.New()
	audienceId := uuid.New()
	datasetId := uuid.New()
	organizationId := uuid.New()
	teamId := uuid.New()

	tests := []struct {
		name          string
		datasetId     uuid.UUID
		audienceType  storemodels.AudienceType
		privilege     storemodels.ResourcePrivilege
		audienceId    uuid.UUID
		mockSetup     func(*mockDatasetService.MockDatasetServiceStore)
		wantErr       bool
		expectedError string
	}{
		{
			name:         "Success case - add user audience",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeUser,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			audienceId:   audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeUser),
									ResourceAudienceId:   currentUserId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
					}, nil)
					ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&storemodels.ResourceAudiencePolicy{}, nil)
					fn(ds)
				})
			},
			wantErr: false,
		},
		{
			name:         "Success case - add team audience",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeTeam,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			audienceId:   teamId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeTeam),
									ResourceAudienceId:   teamId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
					}, nil)
					ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&storemodels.ResourceAudiencePolicy{}, nil)
					fn(ds)
				})
			},
			wantErr: false,
		},
		{
			name:         "Success - add organization audience",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeOrganization,
			audienceId:   organizationId,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
						OrganizationId: organizationId,
					}, nil)
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeUser),
									ResourceAudienceId:   currentUserId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
					}, nil)
					ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&storemodels.ResourceAudiencePolicy{}, nil)
					fn(ds)
				})
			},
			wantErr: false,
		},
		{
			name:         "Success - user has admin access through organization",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeUser,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			audienceId:   audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeOrganization,
							ResourceAudienceID:   organizationId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeOrganization),
									ResourceAudienceId:   organizationId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
					}, nil)
					ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&storemodels.ResourceAudiencePolicy{}, nil)
					fn(ds)
				})
			},
			wantErr: false,
		},
		{
			name:         "Success - user has admin access through team",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeUser,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			audienceId:   audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeTeam,
							ResourceAudienceID:   teamId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeTeam),
									ResourceAudienceId:   teamId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
					}, nil)
					ds.EXPECT().CreateDatasetPolicy(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&storemodels.ResourceAudiencePolicy{}, nil)
					fn(ds)
				})
			},
			wantErr: false,
		},
		{
			name:          "invalid privilege",
			datasetId:     datasetId,
			audienceType:  storemodels.AudienceTypeUser,
			privilege:     storemodels.PrivilegeOrganizationSystemAdmin,
			audienceId:    audienceId,
			mockSetup:     func(ds *mockDatasetService.MockDatasetServiceStore) {},
			wantErr:       true,
			expectedError: "invalid privilege",
		},
		{
			name:          "Error - unsupported audience type",
			datasetId:     datasetId,
			audienceType:  storemodels.AudienceType("invalid_type"),
			privilege:     storemodels.PrivilegeDatasetAdmin,
			audienceId:    teamId,
			mockSetup:     func(ds *mockDatasetService.MockDatasetServiceStore) {},
			wantErr:       true,
			expectedError: "only user audience is supported",
		},
		{
			name:         "Error - user already exists",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeUser,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			audienceId:   audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("user already exists on the dataset")).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
						},
						{
							ResourceAudienceType: storemodels.AudienceTypeUser,
							ResourceAudienceID:   audienceId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
						},
					}, nil)
					fn(ds)
				})
			},
			wantErr:       true,
			expectedError: "user already exists on the dataset",
		},
		{
			name:         "Error - user does not have admin access",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeUser,
			audienceId:   audienceId,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("current user does not have access to change permissions on the dataset")).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{}, nil)
					fn(ds)
				})
			},
			wantErr:       true,
			expectedError: "current user does not have access to change permissions on the dataset",
		},
		{
			name:         "Error - organization already exists",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeOrganization,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			audienceId:   organizationId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("organization already exists on the dataset")).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
						OrganizationId: organizationId,
					}, nil)
					ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
						{
							ResourceAudienceType: storemodels.AudienceTypeOrganization,
							ResourceAudienceID:   organizationId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeOrganization),
									ResourceAudienceId:   organizationId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
						{
							ResourceAudienceType: storemodels.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            storemodels.PrivilegeDatasetAdmin,
							UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
								{
									ResourceAudienceType: string(storemodels.AudienceTypeUser),
									ResourceAudienceId:   currentUserId,
									UserId:               currentUserId,
									Privilege:            storemodels.PrivilegeDatasetAdmin,
								},
							},
						},
					}, nil)
					fn(ds)
				})
			},
			wantErr:       true,
			expectedError: "organization already exists on the dataset",
		},
		{
			name:         "Error - dataset does not belong to organization",
			datasetId:    datasetId,
			audienceType: storemodels.AudienceTypeOrganization,
			audienceId:   organizationId,
			privilege:    storemodels.PrivilegeDatasetAdmin,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().WithDatasetTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("dataset does not belong to the organization")).Run(func(ctx context.Context, fn func(store.DatasetStore) error) {
					ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
						OrganizationId: uuid.New(),
					}, nil)
					fn(ds)
				})
			},
			wantErr:       true,
			expectedError: "dataset does not belong to the organization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			ctx = apicontext.AddAuthToContext(ctx, "user", currentUserId, []uuid.UUID{organizationId})

			policy, err := svc.AddAudienceToDataset(ctx, tt.datasetId, tt.audienceType, tt.audienceId, tt.privilege)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				assert.Nil(t, policy)
			} else {
				assert.NoError(t, err)
			}

			mockDS.AssertExpectations(t)
		})
	}
}

func TestRemoveAudienceFromDataset(t *testing.T) {
	t.Parallel()

	currentUserId := uuid.New()
	audienceId := uuid.New()
	datasetId := uuid.New()
	organizationId := uuid.New()
	teamId := uuid.New()
	policyId1 := uuid.New()
	policyId2 := uuid.New()
	tests := []struct {
		name          string
		datasetId     uuid.UUID
		audienceId    uuid.UUID
		mockSetup     func(*mockDatasetService.MockDatasetServiceStore)
		wantErr       bool
		expectedError string
	}{
		{
			name:       "Success case - remove user audience",
			datasetId:  datasetId,
			audienceId: audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
					{
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						ID:                   policyId1,
						Privilege:            storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       currentUserId,
								ResourceAudiencePolicyId: policyId1,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
					{
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   audienceId,
						ID:                   policyId2,
						Privilege:            storemodels.PrivilegeDatasetViewer,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       audienceId,
								ResourceAudiencePolicyId: policyId2,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetViewer,
							},
						},
					},
				}, nil)
				ds.EXPECT().DeleteDatasetPolicy(mock.Anything, datasetId, storemodels.AudienceTypeUser, audienceId).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "Success case - remove team audience",
			datasetId:  datasetId,
			audienceId: teamId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
					{
						ID:                   policyId1,
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       currentUserId,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
					{
						ID:                   policyId2,
						ResourceAudienceType: storemodels.AudienceTypeTeam,
						ResourceAudienceID:   teamId,
						Privilege:            storemodels.PrivilegeDatasetViewer,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId2,
								ResourceAudienceType:     string(storemodels.AudienceTypeTeam),
								ResourceAudienceId:       teamId,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetViewer,
							},
						},
					},
				}, nil)
				ds.EXPECT().DeleteDatasetPolicy(mock.Anything, datasetId, storemodels.AudienceTypeTeam, teamId).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "Error - user trying to remove their own admin policy",
			datasetId:  datasetId,
			audienceId: currentUserId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
					{
						ID:                   policyId1,
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       currentUserId,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
				}, nil)
			},
			wantErr:       true,
			expectedError: "you cannot change own permissions",
		},
		{
			name:       "Error - audience not found",
			datasetId:  datasetId,
			audienceId: audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
					{
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudienceType: string(storemodels.AudienceTypeUser),
								ResourceAudienceId:   currentUserId,
								UserId:               currentUserId,
								Privilege:            storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
				}, nil)
			},
			wantErr:       true,
			expectedError: "invalid audience id",
		},
		{
			name:       "Error - failed to get dataset policies",
			datasetId:  datasetId,
			audienceId: audienceId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return(nil, fmt.Errorf("failed to get policies"))
			},
			wantErr:       true,
			expectedError: "failed to get policies",
		},
		{
			name:       "Success removing team with no users",
			datasetId:  datasetId,
			audienceId: teamId,
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore) {
				ds.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return([]storemodels.ResourceAudiencePolicy{
					{
						ID:                   policyId1,
						ResourceAudienceType: storemodels.AudienceTypeTeam,
						ResourceAudienceID:   teamId,
						Privilege:            storemodels.PrivilegeDatasetViewer,
						UserPolicies:         []storemodels.FlattenedResourceAudiencePolicy{},
					},
					{
						ID:                   policyId2,
						ResourceAudienceType: storemodels.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId2,
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       currentUserId,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
				}, nil)
				ds.EXPECT().DeleteDatasetPolicy(mock.Anything, datasetId, storemodels.AudienceTypeTeam, teamId).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}
			mockCacheClient := mock_cache.NewMockCacheClient(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			ctx = apicontext.AddAuthToContext(ctx, "user", currentUserId, []uuid.UUID{organizationId})

			err := svc.RemoveAudienceFromDataset(ctx, tt.datasetId, tt.audienceId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockDS.AssertExpectations(t)
		})
	}
}

func TestUpdateDatasetAudiencePrivilege(t *testing.T) {
	t.Parallel()

	currentUserId := uuid.New()
	datasetId := uuid.New()
	organizationId := uuid.New()
	policyId1 := uuid.New()
	policyId2 := uuid.New()
	audienceId := uuid.New()

	tests := []struct {
		name          string
		datasetId     uuid.UUID
		audienceId    uuid.UUID
		privilege     storemodels.ResourcePrivilege
		mockSetup     func(*mockDatasetService.MockDatasetServiceStore)
		wantErr       bool
		expectedError string
	}{
		{
			name:       "Success - update dataset audience privilege",
			datasetId:  datasetId,
			audienceId: audienceId,
			privilege:  storemodels.PrivilegeDatasetAdmin,
			mockSetup: func(mockDS *mockDatasetService.MockDatasetServiceStore) {
				policies := []storemodels.ResourceAudiencePolicy{
					{
						ID:                 policyId1,
						ResourceAudienceID: currentUserId,
						Privilege:          storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       currentUserId,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
					{
						ID:                 policyId2,
						ResourceAudienceID: audienceId,
						Privilege:          storemodels.PrivilegeDatasetViewer,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId2,
								ResourceAudienceType:     string(storemodels.AudienceTypeUser),
								ResourceAudienceId:       audienceId,
								UserId:                   currentUserId,
								Privilege:                storemodels.PrivilegeDatasetViewer,
							},
						},
					},
				}
				mockDS.EXPECT().GetDatasetPolicies(mock.Anything, datasetId).Return(policies, nil)
				mockDS.EXPECT().UpdateDatasetPolicy(mock.Anything, datasetId, audienceId, storemodels.PrivilegeDatasetAdmin).Return(&storemodels.ResourceAudiencePolicy{}, nil)
			},
			wantErr: false,
		},
		{
			name:          "Error - invalid privilege",
			datasetId:     datasetId,
			audienceId:    uuid.New(),
			privilege:     "invalid_privilege",
			mockSetup:     func(mockDS *mockDatasetService.MockDatasetServiceStore) {},
			wantErr:       true,
			expectedError: "invalid privilege",
		},
		{
			name:       "Error - failed to get dataset policies",
			datasetId:  datasetId,
			audienceId: uuid.New(),
			privilege:  storemodels.PrivilegeDatasetViewer,
			mockSetup: func(mockDS *mockDatasetService.MockDatasetServiceStore) {
				mockDS.On("GetDatasetPolicies", mock.Anything, datasetId).Return(nil, fmt.Errorf("db error"))
			},
			wantErr: true,
		},
		{
			name:       "Error - user not admin",
			datasetId:  datasetId,
			audienceId: uuid.New(),
			privilege:  storemodels.PrivilegeDatasetViewer,
			mockSetup: func(mockDS *mockDatasetService.MockDatasetServiceStore) {
				policies := []storemodels.ResourceAudiencePolicy{
					{
						ResourceAudienceID: currentUserId,
						Privilege:          storemodels.PrivilegeDatasetViewer,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								UserId:    currentUserId,
								Privilege: storemodels.PrivilegeDatasetViewer,
							},
						},
					},
				}
				mockDS.On("GetDatasetPolicies", mock.Anything, datasetId).Return(policies, nil)
			},
			wantErr:       true,
			expectedError: "current user does not have access to change permissions on the dataset",
		},
		{
			name:       "Error - audience not found",
			datasetId:  datasetId,
			audienceId: uuid.New(),
			privilege:  storemodels.PrivilegeDatasetViewer,
			mockSetup: func(mockDS *mockDatasetService.MockDatasetServiceStore) {
				policies := []storemodels.ResourceAudiencePolicy{
					{
						ResourceAudienceID: currentUserId,
						Privilege:          storemodels.PrivilegeDatasetAdmin,
						UserPolicies: []storemodels.FlattenedResourceAudiencePolicy{
							{
								UserId:    currentUserId,
								Privilege: storemodels.PrivilegeDatasetAdmin,
							},
						},
					},
				}
				mockDS.On("GetDatasetPolicies", mock.Anything, datasetId).Return(policies, nil)
			},
			wantErr:       true,
			expectedError: "invalid audience id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			ctx = apicontext.AddAuthToContext(ctx, "user", currentUserId, []uuid.UUID{organizationId})

			updatedPolicy, err := svc.UpdateDatasetAudiencePrivilege(ctx, tt.datasetId, tt.audienceId, tt.privilege)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				assert.Nil(t, updatedPolicy)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedPolicy)
			}

			mockDS.AssertExpectations(t)
		})
	}
}

func TestGetRulesByDatasetColumns(t *testing.T) {
	t.Parallel()

	organizationId := uuid.New()
	dataset1Id := uuid.New()
	dataset2Id := uuid.New()
	datasetColumns := []storemodels.DatasetColumn{
		{
			DatasetId: dataset1Id,
			Columns:   []string{"column1"},
		},
		{
			DatasetId: dataset2Id,
			Columns:   []string{"column2"},
		},
	}

	tests := []struct {
		name          string
		mockSetup     func(*mockruleservice.MockRuleService)
		wantErr       bool
		expectedRules map[string]map[string][]rulemodels.Rule
	}{
		{
			name: "Success - returns rules for dataset columns",
			mockSetup: func(rs *mockruleservice.MockRuleService) {
				ruleId1 := uuid.New()
				ruleId2 := uuid.New()
				createdBy := uuid.New()
				updatedBy := uuid.New()
				now := time.Now()

				expectedRules := map[string]map[string][]rulemodels.Rule{
					dataset1Id.String(): {
						"column1": {
							{
								ID:             ruleId1,
								OrganizationId: organizationId,
								DatasetId:      dataset1Id,
								Column:         "column1",
								Value:          "test_value_1",
								FilterConfig: rulemodels.FilterConfig{
									QueryConfig: querybuildermodels.QueryConfig{},
									Sql:         "test_sql_1",
									Args:        map[string]interface{}{"arg1": "value1"},
								},
								Title:       "Test Rule 1",
								Description: "Test Description 1",
								Priority:    1,
								CreatedAt:   now,
								CreatedBy:   createdBy,
								UpdatedAt:   now,
								UpdatedBy:   updatedBy,
							},
						},
					},
					dataset2Id.String(): {
						"column2": {
							{
								ID:             ruleId2,
								OrganizationId: organizationId,
								DatasetId:      dataset2Id,
								Column:         "column2",
								Value:          "test_value_2",
								FilterConfig: rulemodels.FilterConfig{
									QueryConfig: querybuildermodels.QueryConfig{},
									Sql:         "test_sql_2",
									Args:        map[string]interface{}{"arg2": "value2"},
								},
								Title:       "Test Rule 2",
								Description: "Test Description 2",
								Priority:    2,
								CreatedAt:   now,
								CreatedBy:   createdBy,
								UpdatedAt:   now,
								UpdatedBy:   updatedBy,
							},
						},
					},
				}

				rs.EXPECT().GetRules(mock.Anything, storemodels.FilterRuleParams{
					OrganizationId: organizationId,
					DatasetColumns: datasetColumns,
				}).Return(expectedRules, nil)
			},
			wantErr: false,
			expectedRules: map[string]map[string][]rulemodels.Rule{
				dataset1Id.String(): {
					"column1": {
						{
							ID:             uuid.New(),
							OrganizationId: organizationId,
							DatasetId:      dataset1Id,
							Column:         "column1",
							Value:          "test_value_1",
							FilterConfig: rulemodels.FilterConfig{
								QueryConfig: querybuildermodels.QueryConfig{},
								Sql:         "test_sql_1",
								Args:        map[string]interface{}{"arg1": "value1"},
							},
							Title:       "Test Rule 1",
							Description: "Test Description 1",
							Priority:    1,
							CreatedAt:   time.Now(),
							CreatedBy:   uuid.New(),
							UpdatedAt:   time.Now(),
							UpdatedBy:   uuid.New(),
						},
					},
				},
				dataset2Id.String(): {
					"column2": {
						{
							ID:             uuid.New(),
							OrganizationId: organizationId,
							DatasetId:      dataset2Id,
							Column:         "column2",
							Value:          "test_value_2",
							FilterConfig: rulemodels.FilterConfig{
								QueryConfig: querybuildermodels.QueryConfig{},
								Sql:         "test_sql_2",
								Args:        map[string]interface{}{"arg2": "value2"},
							},
							Title:       "Test Rule 2",
							Description: "Test Description 2",
							Priority:    2,
							CreatedAt:   time.Now(),
							CreatedBy:   uuid.New(),
							UpdatedAt:   time.Now(),
							UpdatedBy:   uuid.New(),
						},
					},
				},
			},
		},
		{
			name: "Error - rule service returns error",
			mockSetup: func(rs *mockruleservice.MockRuleService) {
				rs.EXPECT().GetRules(mock.Anything, storemodels.FilterRuleParams{
					OrganizationId: organizationId,
					DatasetColumns: datasetColumns,
				}).Return(nil, fmt.Errorf("failed to get rules"))
			},
			wantErr:       true,
			expectedRules: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockRuleService)
			}

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())

			service := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			// Execute
			rules, err := service.GetRulesByDatasetColumns(ctx, organizationId, datasetColumns)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rules)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedRules), len(rules))
				for datasetID, columnRules := range rules {
					expectedColumnRules := tt.expectedRules[datasetID]
					assert.Equal(t, len(expectedColumnRules), len(columnRules))
					for column, rulesList := range columnRules {
						expectedRules := expectedColumnRules[column]
						assert.Equal(t, len(expectedRules), len(rulesList))
						for i, rule := range rulesList {
							expectedRule := expectedRules[i]
							assert.Equal(t, expectedRule.Column, rule.Column)
							assert.Equal(t, expectedRule.Value, rule.Value)
							assert.Equal(t, expectedRule.Title, rule.Title)
							assert.Equal(t, expectedRule.Description, rule.Description)
							assert.Equal(t, expectedRule.Priority, rule.Priority)
							assert.Equal(t, expectedRule.FilterConfig.Sql, rule.FilterConfig.Sql)
							assert.Equal(t, expectedRule.FilterConfig.Args, rule.FilterConfig.Args)
						}
					}
				}
			}

			mockRuleService.AssertExpectations(t)
		})
	}
}

func TestGetRulesByIds(t *testing.T) {
	t.Parallel()

	ruleId1 := uuid.New()
	ruleId2 := uuid.New()
	organizationId := uuid.New()
	datasetId := uuid.New()
	createdBy := uuid.New()
	updatedBy := uuid.New()
	now := time.Now()

	tests := []struct {
		name          string
		ruleIds       []string
		mockSetup     func(*mockruleservice.MockRuleService)
		wantErr       bool
		expectedRules []rulemodels.Rule
	}{
		{
			name:    "Success - returns rules for given ids",
			ruleIds: []string{ruleId1.String(), ruleId2.String()},
			mockSetup: func(rs *mockruleservice.MockRuleService) {
				expectedRules := []rulemodels.Rule{
					{
						ID:             ruleId1,
						OrganizationId: organizationId,
						DatasetId:      datasetId,
						Column:         "column1",
						Value:          "test_value_1",
						FilterConfig: rulemodels.FilterConfig{
							QueryConfig: querybuildermodels.QueryConfig{},
							Sql:         "test_sql_1",
							Args:        map[string]interface{}{"arg1": "value1"},
						},
						Title:       "Test Rule 1",
						Description: "Test Description 1",
						Priority:    1,
						CreatedAt:   now,
						CreatedBy:   createdBy,
						UpdatedAt:   now,
						UpdatedBy:   updatedBy,
					},
					{
						ID:             ruleId2,
						OrganizationId: organizationId,
						DatasetId:      datasetId,
						Column:         "column2",
						Value:          "test_value_2",
						FilterConfig: rulemodels.FilterConfig{
							QueryConfig: querybuildermodels.QueryConfig{},
							Sql:         "test_sql_2",
							Args:        map[string]interface{}{"arg2": "value2"},
						},
						Title:       "Test Rule 2",
						Description: "Test Description 2",
						Priority:    2,
						CreatedAt:   now,
						CreatedBy:   createdBy,
						UpdatedAt:   now,
						UpdatedBy:   updatedBy,
					},
				}

				rs.EXPECT().GetRuleByIds(mock.Anything, []uuid.UUID{ruleId1, ruleId2}).Return(expectedRules, nil)
			},
			wantErr: false,
			expectedRules: []rulemodels.Rule{
				{
					ID:             ruleId1,
					OrganizationId: organizationId,
					DatasetId:      datasetId,
					Column:         "column1",
					Value:          "test_value_1",
					FilterConfig: rulemodels.FilterConfig{
						QueryConfig: querybuildermodels.QueryConfig{},
						Sql:         "test_sql_1",
						Args:        map[string]interface{}{"arg1": "value1"},
					},
					Title:       "Test Rule 1",
					Description: "Test Description 1",
					Priority:    1,
					CreatedAt:   now,
					CreatedBy:   createdBy,
					UpdatedAt:   now,
					UpdatedBy:   updatedBy,
				},
				{
					ID:             ruleId2,
					OrganizationId: organizationId,
					DatasetId:      datasetId,
					Column:         "column2",
					Value:          "test_value_2",
					FilterConfig: rulemodels.FilterConfig{
						QueryConfig: querybuildermodels.QueryConfig{},
						Sql:         "test_sql_2",
						Args:        map[string]interface{}{"arg2": "value2"},
					},
					Title:       "Test Rule 2",
					Description: "Test Description 2",
					Priority:    2,
					CreatedAt:   now,
					CreatedBy:   createdBy,
					UpdatedAt:   now,
					UpdatedBy:   updatedBy,
				},
			},
		},
		{
			name:    "Error - invalid rule id",
			ruleIds: []string{"invalid-uuid"},
			mockSetup: func(rs *mockruleservice.MockRuleService) {
				// No mock expectations since the function should return early
			},
			wantErr:       true,
			expectedRules: nil,
		},
		{
			name:    "Error - rule service returns error",
			ruleIds: []string{ruleId1.String(), ruleId2.String()},
			mockSetup: func(rs *mockruleservice.MockRuleService) {
				rs.EXPECT().GetRuleByIds(mock.Anything, []uuid.UUID{ruleId1, ruleId2}).Return(nil, fmt.Errorf("failed to get rules"))
			},
			wantErr:       true,
			expectedRules: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockRuleService)
			}

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())

			service := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			// Execute
			rules, err := service.GetRulesByIds(ctx, tt.ruleIds)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rules)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedRules), len(rules))
				for i, rule := range rules {
					expectedRule := tt.expectedRules[i]
					assert.Equal(t, expectedRule.ID, rule.ID)
					assert.Equal(t, expectedRule.OrganizationId, rule.OrganizationId)
					assert.Equal(t, expectedRule.DatasetId, rule.DatasetId)
					assert.Equal(t, expectedRule.Column, rule.Column)
					assert.Equal(t, expectedRule.Value, rule.Value)
				}
			}

			mockRuleService.AssertExpectations(t)
		})
	}
}

func TestInitiateDatasetExport(t *testing.T) {
	tests := []struct {
		name          string
		merchantId    uuid.UUID
		datasetId     uuid.UUID
		queryConfig   models.DatasetParams
		userId        uuid.UUID
		mockSetup     func(*mockDatasetService.MockDatasetServiceStore, *mock_temporal.MockTemporalService)
		wantErr       bool
		expectedError string
	}{
		{
			name:       "Success case",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			userId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			queryConfig: models.DatasetParams{
				Filters: models.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []models.Filter{},
				},
			},
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().CreateDatasetAction(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				ts.EXPECT().ExecuteAsyncWorkflow(mock.Anything, mock.Anything).Return(temporalmodels.WorkflowResponse{
					WorkflowID: "test-workflow",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:       "Error - dataset not found",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			userId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			queryConfig: models.DatasetParams{
				Filters: models.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []models.Filter{},
				},
			},
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().CreateDatasetAction(mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("action creation failed"))
			},
			wantErr:       true,
			expectedError: "action creation failed",
		},
		{
			name:       "Error - workflow execution failed",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			datasetId:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			userId:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			queryConfig: models.DatasetParams{
				Filters: models.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []models.Filter{},
				},
			},
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().CreateDatasetAction(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				ts.EXPECT().ExecuteAsyncWorkflow(mock.Anything, mock.Anything).Return(temporalmodels.WorkflowResponse{}, fmt.Errorf("workflow execution failed"))
			},
			wantErr:       true,
			expectedError: "workflow execution failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS, mockTemporalService)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			ctx = apicontext.AddAuthToContext(ctx, "user", tt.userId, []uuid.UUID{tt.merchantId})

			workflowId, err := svc.CreateDatasetExportAction(ctx, tt.merchantId, tt.datasetId.String(), tt.queryConfig, tt.userId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
				assert.Empty(t, workflowId)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, workflowId)
			}

			mockDS.AssertExpectations(t)
			mockTemporalService.AssertExpectations(t)
		})
	}
}

func TestExportDatasetActivity(t *testing.T) {
	tests := []struct {
		name          string
		params        models.DatasetExportParams
		userId        uuid.UUID
		orgIds        []uuid.UUID
		datasetId     uuid.UUID
		workflowId    string
		exportPath    string
		queryConfig   models.DatasetParams
		mockSetup     func(*mockDatasetService.MockDatasetServiceStore, *mockDataplatform.MockDataPlatformService, *mock_cloudservice.MockCloudService, *mock_querybuilder.MockQueryBuilder)
		wantErr       bool
		expectedError string
		expectedURL   string
	}{
		{
			name:       "Success case",
			userId:     uuid.New(),
			orgIds:     []uuid.UUID{uuid.New()},
			datasetId:  uuid.New(),
			workflowId: "test-workflow",
			exportPath: "test/export/path.csv",
			queryConfig: models.DatasetParams{
				Filters: models.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []models.Filter{},
				},
			},
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore, dps *mockDataplatform.MockDataPlatformService, cs *mock_cloudservice.MockCloudService, qb *mock_querybuilder.MockQueryBuilder) {
				ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}, "display_config": [{"column": "column1", "is_hidden": false, "is_editable": false}, {"column": "column2", "is_hidden": false, "is_editable": false}, {"column": "column3", "is_hidden": false, "is_editable": false}, {"column": "column4", "is_hidden": false, "is_editable": true, "type": "tags"}]}`),
				}, nil)
				dps.EXPECT().GetDatasetMetadata(mock.Anything, mock.Anything, mock.Anything).Return(dataplatformDataModels.DatasetMetadata{
					Schema: map[string]dataplatformDataModels.ColumnMetadata{
						"column1": {Type: "string"},
						"column2": {Type: "string"},
						"column3": {Type: "string"},
						"column4": {Type: "tags"},
					},
				}, nil)

				qb.EXPECT().ToSQL(mock.Anything, mock.Anything).Return("SELECT * FROM dataset", map[string]interface{}{}, nil)

				dps.EXPECT().Query(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{{Name: "col1"}, {Name: "col2"}},
					Rows:    []map[string]interface{}{{"col1": "val1", "col2": "val2"}},
				}, nil)

				cs.EXPECT().UploadFileToCloud(mock.Anything, mock.Anything, mock.Anything).Return(cloudservicemodels.SignedUrlToUpload{
					Url: "https://exported-file.csv",
				}, nil)
				ds.EXPECT().UpdateDatasetActionStatus(mock.Anything, "test-workflow", "SUCCESSFUL").Return(nil)
			},
			expectedURL: "https://exported-file.csv",
		},
		{
			name:       "Error - query failure",
			userId:     uuid.New(),
			orgIds:     []uuid.UUID{uuid.New()},
			datasetId:  uuid.New(),
			workflowId: "test-workflow",
			exportPath: "test/export/path.csv",
			queryConfig: models.DatasetParams{
				Filters: models.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []models.Filter{},
				},
			},
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore, dps *mockDataplatform.MockDataPlatformService, cs *mock_cloudservice.MockCloudService, qb *mock_querybuilder.MockQueryBuilder) {
				ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}, "display_config": [{"column": "column1", "is_hidden": false, "is_editable": false}, {"column": "column2", "is_hidden": false, "is_editable": false}, {"column": "column3", "is_hidden": false, "is_editable": false}, {"column": "column4", "is_hidden": false, "is_editable": true, "type": "tags"}]}`),
				}, nil)
				dps.EXPECT().GetDatasetMetadata(mock.Anything, mock.Anything, mock.Anything).Return(dataplatformDataModels.DatasetMetadata{
					Schema: map[string]dataplatformDataModels.ColumnMetadata{
						"column1": {Type: "string"},
						"column2": {Type: "string"},
						"column3": {Type: "string"},
						"column4": {Type: "tags"},
					},
				}, nil)

				qb.EXPECT().ToSQL(mock.Anything, mock.Anything).Return("SELECT * FROM dataset", map[string]interface{}{}, nil)

				dps.EXPECT().Query(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{{Name: "col1"}, {Name: "col2"}},
					Rows:    []map[string]interface{}{{"col1": "val1", "col2": "val2"}},
				}, fmt.Errorf("query failed"))
			},
			wantErr:       true,
			expectedError: "failed to get dataset data",
		},
		{
			name:       "Error - upload failure",
			userId:     uuid.New(),
			orgIds:     []uuid.UUID{uuid.New()},
			datasetId:  uuid.New(),
			workflowId: "test-workflow",
			exportPath: "test/export/path.csv",
			queryConfig: models.DatasetParams{
				Filters: models.FilterModel{
					LogicalOperator: "AND",
					Conditions:      []models.Filter{},
				},
			},
			mockSetup: func(ds *mockDatasetService.MockDatasetServiceStore, dps *mockDataplatform.MockDataPlatformService, cs *mock_cloudservice.MockCloudService, qb *mock_querybuilder.MockQueryBuilder) {
				ds.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Metadata: json.RawMessage(`{"columns": {"column4": {"custom_type": "tags"}}, "display_config": [{"column": "column1", "is_hidden": false, "is_editable": false}, {"column": "column2", "is_hidden": false, "is_editable": false}, {"column": "column3", "is_hidden": false, "is_editable": false}, {"column": "column4", "is_hidden": false, "is_editable": true, "type": "tags"}]}`),
				}, nil)
				dps.EXPECT().GetDatasetMetadata(mock.Anything, mock.Anything, mock.Anything).Return(dataplatformDataModels.DatasetMetadata{
					Schema: map[string]dataplatformDataModels.ColumnMetadata{
						"column1": {Type: "string"},
						"column2": {Type: "string"},
						"column3": {Type: "string"},
						"column4": {Type: "tags"},
					},
				}, nil)

				qb.EXPECT().ToSQL(mock.Anything, mock.Anything).Return("SELECT * FROM dataset", map[string]interface{}{}, nil)

				dps.EXPECT().Query(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataplatformmodels.QueryResult{
					Columns: []dataplatformmodels.ColumnMetadata{{Name: "col1"}, {Name: "col2"}},
					Rows:    []map[string]interface{}{{"col1": "val1", "col2": "val2"}},
				}, nil)

				cs.EXPECT().UploadFileToCloud(mock.Anything, mock.Anything, mock.Anything).Return(cloudservicemodels.SignedUrlToUpload{}, fmt.Errorf("upload failed"))
			},
			wantErr:       true,
			expectedError: "upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := new(mockfileimports.MockFileImportService)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS, mockDPS, mockCloudService, mockQueryBuilder)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			url, err := svc.DatasetExportTemporalActivity(context.Background(), tt.params, tt.datasetId, tt.userId, tt.orgIds, tt.workflowId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}

			mockDS.AssertExpectations(t)
			mockDPS.AssertExpectations(t)
			mockCloudService.AssertExpectations(t)
		})
	}
}

func TestInitiateFilePreparationForDatasetImport(t *testing.T) {
	datasetId := uuid.New()
	fileId := uuid.New()
	userId := uuid.New()
	orgId := uuid.New()
	tests := []struct {
		name          string
		datasetId     uuid.UUID
		fileId        uuid.UUID
		mockSetup     func(*mock_store.MockStore, *mock_temporal.MockTemporalService)
		wantErr       bool
		userId        *uuid.UUID
		orgIds        []uuid.UUID
		expectedError string
	}{
		{
			name:      "Success case",
			datasetId: datasetId,
			fileId:    fileId,
			mockSetup: func(ds *mock_store.MockStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().WithTx(mock.Anything, mock.AnythingOfType("func(store.Store) error")).Return(nil).
					Run(func(ctx context.Context, fn func(store.Store) error) {
						ds.EXPECT().CreateDatasetFileUpload(mock.Anything, mock.Anything).Return(&storemodels.DatasetFileUpload{}, nil)
						ds.EXPECT().CreateDatasetAction(mock.Anything, mock.Anything, mock.Anything).Return(nil)
						ts.EXPECT().ExecuteAsyncWorkflow(mock.Anything, mock.Anything).Return(temporalmodels.WorkflowResponse{}, nil)
						fn(ds)
					})
			},
			wantErr:       false,
			expectedError: "",
			userId:        &userId,
			orgIds:        []uuid.UUID{orgId},
		},
		{
			name:      "Error - CreateDatasetFileUpload fails",
			datasetId: datasetId,
			fileId:    fileId,
			mockSetup: func(ds *mock_store.MockStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().WithTx(mock.Anything, mock.AnythingOfType("func(store.Store) error")).Return(fmt.Errorf("failed to create dataset file upload")).
					Run(func(ctx context.Context, fn func(store.Store) error) {
						ds.EXPECT().CreateDatasetFileUpload(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to create dataset file upload"))
						fn(ds)
					})
			},
			wantErr:       true,
			expectedError: "failed to create dataset file upload",
			userId:        &userId,
			orgIds:        []uuid.UUID{orgId},
		},
		{
			name:      "Error - CreateDatasetAction fails",
			datasetId: datasetId,
			fileId:    fileId,
			mockSetup: func(ds *mock_store.MockStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().WithTx(mock.Anything, mock.AnythingOfType("func(store.Store) error")).Return(fmt.Errorf("failed to create dataset action")).
					Run(func(ctx context.Context, fn func(store.Store) error) {
						ds.EXPECT().CreateDatasetFileUpload(mock.Anything, mock.Anything).Return(&storemodels.DatasetFileUpload{}, nil)
						ds.EXPECT().CreateDatasetAction(mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("failed to create dataset action"))
						fn(ds)
					})
			},
			wantErr:       true,
			expectedError: "failed to create dataset action",
			userId:        &userId,
			orgIds:        []uuid.UUID{orgId},
		},
		{
			name:      "Error - ExecuteAsyncWorkflow fails",
			datasetId: datasetId,
			fileId:    fileId,
			mockSetup: func(ds *mock_store.MockStore, ts *mock_temporal.MockTemporalService) {
				ds.EXPECT().WithTx(mock.Anything, mock.AnythingOfType("func(store.Store) error")).Return(fmt.Errorf("failed to execute workflow")).
					Run(func(ctx context.Context, fn func(store.Store) error) {
						ds.EXPECT().CreateDatasetFileUpload(mock.Anything, mock.Anything).Return(&storemodels.DatasetFileUpload{}, nil)
						ds.EXPECT().CreateDatasetAction(mock.Anything, mock.Anything, mock.Anything).Return(nil)
						ts.EXPECT().ExecuteAsyncWorkflow(mock.Anything, mock.Anything).Return(temporalmodels.WorkflowResponse{}, fmt.Errorf("failed to execute workflow"))
						fn(ds)
					})
			},
			wantErr:       true,
			expectedError: "failed to execute workflow",
			userId:        &userId,
			orgIds:        []uuid.UUID{orgId},
		},
		{
			name:          "Error - No user ID in context",
			datasetId:     datasetId,
			fileId:        fileId,
			mockSetup:     func(ds *mock_store.MockStore, ts *mock_temporal.MockTemporalService) {},
			wantErr:       true,
			expectedError: "unauthorized",
			userId:        nil,
			orgIds:        []uuid.UUID{orgId},
		},
		{
			name:          "Error - No org IDs in context",
			datasetId:     datasetId,
			fileId:        fileId,
			mockSetup:     func(ds *mock_store.MockStore, ts *mock_temporal.MockTemporalService) {},
			wantErr:       true,
			expectedError: "unauthorized",
			userId:        &userId,
			orgIds:        []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := mock_store.NewMockStore(t)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := mockfileimports.NewMockFileImportService(t)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS, mockTemporalService)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			ctx := context.Background()
			if tt.userId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.userId, tt.orgIds)
			}

			actionId, err := svc.InitiateFilePreparationForDatasetImport(ctx, tt.datasetId, tt.fileId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
				assert.Nil(t, actionId)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, actionId)
			}

			mockDS.AssertExpectations(t)
			mockTemporalService.AssertExpectations(t)
		})
	}
}

func TestImportDataFromFile(t *testing.T) {
	t.Parallel()

	datasetId := uuid.New()
	merchantId := uuid.New()
	fileId := uuid.New()

	tests := []struct {
		name          string
		merchantId    uuid.UUID
		datasetId     uuid.UUID
		fileId        uuid.UUID
		mockSetup     func(*mock_store.MockStore, *mock_s3.MockS3Client, *mockDataplatform.MockDataPlatformService)
		wantErr       bool
		expectedError string
	}{
		{
			name:       "success",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						ID:                   fileId,
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
						Metadata:             []byte(`{"transformed_data_bucket":"bucket","transformed_data_path":"/path"}`),
					},
				}, nil)

				mockS3.EXPECT().GetFileDetails(mock.Anything, "bucket", "/path").Return(&s3.FileInfo{
					Size: 100,
				}, nil)
				m.EXPECT().GetDags(mock.Anything, merchantId.String()).Return(map[string]*servicemodels.DAGNode{
					datasetId.String(): {
						NodeId:   datasetId.String(),
						NodeType: servicemodels.NodeTypeDataset,
						Parents: []*servicemodels.DAGNode{
							{
								NodeId:   "s3a://zamp-prd-file-imports/careem-prd/copilot-accounts",
								NodeType: servicemodels.NodeTypeFolder,
							},
						},
					},
				}, nil)
				mockS3.EXPECT().GetSampleFilePathFromFolder(mock.Anything, "zamp-prd-file-imports", "careem-prd/copilot-accounts/").Return("s3a://zamp-prd-file-imports/careem-prd/copilot-accounts/file1.parquet", nil)
				mockS3.EXPECT().CopyFile(mock.Anything, "bucket", "/path", "zamp-prd-file-imports", "s3a:/zamp-prd-file-imports/careem-prd/copilot-accounts/"+fileId.String()+".parquet").Return(nil)
			},
		},
		{
			name:       "error - failed to get dataset file uploads",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return(nil, errors.New("db error"))
			},
			wantErr:       true,
			expectedError: "Failed to get dataset file uploads",
		},
		{
			name:       "error - file not associated with dataset",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{}, nil)
			},
			wantErr:       true,
			expectedError: "The given file is not associated with the dataset",
		},
		{
			name:       "error - file not prepared",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusPending,
					},
				}, nil)
			},
			wantErr:       true,
			expectedError: "The given file is not prepared for dataset import",
		},
		{
			name:       "error - invalid metadata",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
						Metadata:             []byte(`invalid json`),
					},
				}, nil)
			},
			wantErr:       true,
			expectedError: "Invalid file upload metadata",
		},
		{
			name:       "error - failed to get file details",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
						Metadata:             []byte(`{"transformed_data_bucket":"bucket","transformed_data_path":"/path"}`),
					},
				}, nil)

				mockS3.EXPECT().GetFileDetails(mock.Anything, "bucket", "/path").Return(nil, errors.New("s3 error"))
			},
			wantErr:       true,
			expectedError: "Failed to retrieve the given file details",
		},
		{
			name:       "error - empty file",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
						Metadata:             []byte(`{"transformed_data_bucket":"bucket","transformed_data_path":"/path"}`),
					},
				}, nil)

				mockS3.EXPECT().GetFileDetails(mock.Anything, "bucket", "/path").Return(&s3.FileInfo{
					Size: 0,
				}, nil)
			},
			wantErr:       true,
			expectedError: "The given file is empty",
		},
		{
			name:       "error - failed to get dataset import path",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
						Metadata:             []byte(`{"transformed_data_bucket":"bucket","transformed_data_path":"/path"}`),
					},
				}, nil)

				mockS3.EXPECT().GetFileDetails(mock.Anything, "bucket", "/path").Return(&s3.FileInfo{
					Size: 100,
				}, nil)

				m.EXPECT().GetDags(mock.Anything, merchantId.String()).Return(nil, errors.New("dag error"))
			},
			wantErr:       true,
			expectedError: "Failed to get dataset import path",
		},
		{
			name:       "error - failed to copy file",
			merchantId: merchantId,
			datasetId:  datasetId,
			fileId:     fileId,
			mockSetup: func(mockDS *mock_store.MockStore, mockS3 *mock_s3.MockS3Client, m *mockDataplatform.MockDataPlatformService) {
				mockDS.EXPECT().GetDatasetFileUploadByDatasetId(mock.Anything, datasetId).Return([]storemodels.DatasetFileUpload{
					{
						ID:                   fileId,
						FileUploadID:         fileId,
						FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
						Metadata:             []byte(`{"transformed_data_bucket":"bucket","transformed_data_path":"/path"}`),
					},
				}, nil)

				mockS3.EXPECT().GetFileDetails(mock.Anything, "bucket", "/path").Return(&s3.FileInfo{
					Size: 100,
				}, nil)

				m.EXPECT().GetDags(mock.Anything, merchantId.String()).Return(map[string]*servicemodels.DAGNode{
					datasetId.String(): {
						NodeId:   datasetId.String(),
						NodeType: servicemodels.NodeTypeDataset,
						Parents: []*servicemodels.DAGNode{
							{
								NodeId:   "s3a://zamp-prd-file-imports/careem-prd/copilot-accounts",
								NodeType: servicemodels.NodeTypeFolder,
							},
						},
					},
				}, nil)

				mockS3.EXPECT().GetSampleFilePathFromFolder(mock.Anything, "zamp-prd-file-imports", "careem-prd/copilot-accounts/").Return("s3a://zamp-prd-file-imports/careem-prd/copilot-accounts/file1.parquet", nil)
				mockS3.EXPECT().CopyFile(mock.Anything, "bucket", "/path", "zamp-prd-file-imports", "s3a:/zamp-prd-file-imports/careem-prd/copilot-accounts/"+fileId.String()+".parquet").Return(errors.New("copy error"))
			},
			wantErr:       true,
			expectedError: "Failed to import the given file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDS := mock_store.NewMockStore(t)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)
			mockFileUploadsService := mockfileimports.NewMockFileImportService(t)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS, mockS3Client, mockDPS)
			}

			svc := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			err := svc.ImportDataFromFile(context.Background(), tt.merchantId, tt.datasetId, tt.fileId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}

			mockDS.AssertExpectations(t)
			mockS3Client.AssertExpectations(t)
		})
	}
}

func TestGetFileUploadPreview(t *testing.T) {
	t.Parallel()

	fileUploadId := uuid.New()

	tests := []struct {
		name            string
		fileUploadId    uuid.UUID
		mockSetup       func(*mock_store.MockStore)
		wantErr         bool
		expectedError   string
		expectedPreview storemodels.DatasetPreview
	}{
		{
			name:         "success",
			fileUploadId: fileUploadId,
			mockSetup: func(mockDS *mock_store.MockStore) {
				mockDS.EXPECT().GetDatasetFileUploadById(mock.Anything, fileUploadId).Return(storemodels.DatasetFileUpload{
					ID:                   fileUploadId,
					FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
					Metadata:             []byte(`{"data_preview":{"columns":["col1","col2"],"rows":[{"col1":"val1","col2":"val2"}]}}`),
				}, nil)
			},
			expectedPreview: storemodels.DatasetPreview{
				Columns: []string{"col1", "col2"},
				Rows:    []map[string]interface{}{{"col1": "val1", "col2": "val2"}},
			},
		},
		{
			name:         "error - failed to get dataset file upload",
			fileUploadId: fileUploadId,
			mockSetup: func(mockDS *mock_store.MockStore) {
				mockDS.EXPECT().GetDatasetFileUploadById(mock.Anything, fileUploadId).Return(storemodels.DatasetFileUpload{}, fmt.Errorf("failed to get dataset file uploads: %w", errors.New("db error")))
			},
			wantErr:       true,
			expectedError: "failed to get dataset file uploads",
		},
		{
			name:         "error - no dataset file upload found",
			fileUploadId: fileUploadId,
			mockSetup: func(mockDS *mock_store.MockStore) {
				mockDS.EXPECT().GetDatasetFileUploadById(mock.Anything, fileUploadId).Return(storemodels.DatasetFileUpload{
					ID: uuid.Nil,
				}, nil)
			},
			wantErr:       true,
			expectedError: "no dataset file upload found",
		},
		{
			name:         "error - file alignment not completed",
			fileUploadId: fileUploadId,
			mockSetup: func(mockDS *mock_store.MockStore) {
				mockDS.EXPECT().GetDatasetFileUploadById(mock.Anything, fileUploadId).Return(storemodels.DatasetFileUpload{
					ID:                   fileUploadId,
					FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusPending,
				}, nil)
			},
			wantErr:       true,
			expectedError: "Preview is not ready yet",
		},
		{
			name:         "error - invalid metadata",
			fileUploadId: fileUploadId,
			mockSetup: func(mockDS *mock_store.MockStore) {
				mockDS.EXPECT().GetDatasetFileUploadById(mock.Anything, fileUploadId).Return(storemodels.DatasetFileUpload{
					ID:                   fileUploadId,
					FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusCompleted,
					Metadata:             []byte(`invalid json`),
				}, nil)
			},
			wantErr:       true,
			expectedError: "failed to unmarshal dataset file upload metadata",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockDS := mock_store.NewMockStore(t)
			mockDPS := mockDataplatform.NewMockDataPlatformService(t)
			mockQueryBuilder := mock_querybuilder.NewMockQueryBuilder(t)
			mockRuleService := mockruleservice.NewMockRuleService(t)
			mockFileUploadsService := mockfileimports.NewMockFileImportService(t)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS)
			}

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			service := NewDatasetService(mockDS, mockQueryBuilder, mockDPS, mockRuleService, mockFileUploadsService, mockTemporalService, mockCloudService, mockS3Client, serverConfig, mockCacheClient)

			// Execute
			preview, err := service.GetFileUploadPreview(ctx, tt.fileUploadId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, preview)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPreview, preview)
			}

			mockDS.AssertExpectations(t)
		})
	}
}

func TestGetDatasetDisplayConfig(t *testing.T) {
	tests := []struct {
		name           string
		merchantId     uuid.UUID
		datasetId      string
		setupMocks     func(*mock_store.MockStore, *mockDataplatform.MockDataPlatformService)
		expectedResult []models.DisplayConfig
		expectedError  error
	}{
		{
			name:       "Success case - existing display config",
			merchantId: uuid.New(),
			datasetId:  uuid.New().String(),
			setupMocks: func(mockStore *mock_store.MockStore, mockDataplatformService *mockDataplatform.MockDataPlatformService) {
				// Setup mock for GetDatasetById
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					Metadata: json.RawMessage(`{"display_config":[{"column":"test_column","is_hidden":false,"is_editable":true}]}`),
				}, nil)
			},
			expectedResult: []models.DisplayConfig{
				{
					Column:     "test_column",
					IsHidden:   false,
					IsEditable: true,
				},
			},
			expectedError: nil,
		},
		{
			name:       "Success case - empty display config",
			merchantId: uuid.New(),
			datasetId:  uuid.New().String(),
			setupMocks: func(mockStore *mock_store.MockStore, mockDataplatformService *mockDataplatform.MockDataPlatformService) {
				// Setup mock for GetDatasetById
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					Metadata: json.RawMessage(`{}`),
				}, nil)

				// Setup mock for GetDatasetMetadata
				mockDataplatformService.EXPECT().GetDatasetMetadata(mock.Anything, mock.Anything, mock.Anything).Return(dataplatformDataModels.DatasetMetadata{
					Schema: map[string]dataplatformDataModels.ColumnMetadata{
						"column1": {},
						"column2": {},
					},
				}, nil)
			},
			expectedResult: []models.DisplayConfig{
				{
					Column:     "column1",
					IsHidden:   false,
					IsEditable: false,
				},
				{
					Column:     "column2",
					IsHidden:   false,
					IsEditable: false,
				},
			},
			expectedError: nil,
		},
		{
			name:       "Error case - failed to get dataset by id",
			merchantId: uuid.New(),
			datasetId:  uuid.New().String(),
			setupMocks: func(mockStore *mock_store.MockStore, mockDataplatformService *mockDataplatform.MockDataPlatformService) {
				// Setup mock for GetDatasetById
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(nil, errors.New("failed to get dataset"))
			},
			expectedResult: nil,
			expectedError:  datasetErrors.ErrFailedToGetDatasetById,
		},
		{
			name:       "Error case - failed to unmarshal metadata",
			merchantId: uuid.New(),
			datasetId:  uuid.New().String(),
			setupMocks: func(mockStore *mock_store.MockStore, mockDataplatformService *mockDataplatform.MockDataPlatformService) {
				// Setup mock for GetDatasetById with invalid JSON
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					Metadata: json.RawMessage(`{invalid json`),
				}, nil)
			},
			expectedResult: nil,
			expectedError:  datasetErrors.ErrFailedToUnmarshalMetadata,
		},
		{
			name:       "Error case - failed to get dataset metadata",
			merchantId: uuid.New(),
			datasetId:  uuid.New().String(),
			setupMocks: func(mockStore *mock_store.MockStore, mockDataplatformService *mockDataplatform.MockDataPlatformService) {
				// Setup mock for GetDatasetById
				mockStore.EXPECT().GetDatasetById(mock.Anything, mock.Anything).Return(&storemodels.Dataset{
					Metadata: json.RawMessage(`{}`),
				}, nil)

				// Setup mock for GetDatasetMetadata with error
				mockDataplatformService.EXPECT().GetDatasetMetadata(mock.Anything, mock.Anything, mock.Anything).Return(dataplatformDataModels.DatasetMetadata{}, errors.New("failed to get metadata"))
			},
			expectedResult: nil,
			expectedError:  datasetErrors.ErrFailedToGetDatasetMetadata,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockStore := mock_store.NewMockStore(t)
			mockDataplatformService := mockDataplatform.NewMockDataPlatformService(t)
			mockQueryBuilder := mock_querybuilder.NewMockQueryBuilder(t)
			mockRuleService := mockruleservice.NewMockRuleService(t)
			mockFileImportService := mockfileimports.NewMockFileImportService(t)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)
			mockS3Client := mock_s3.NewMockS3Client(t)

			tt.setupMocks(mockStore, mockDataplatformService)

			// Create service
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}
			service := NewDatasetService(
				mockStore,
				mockQueryBuilder,
				mockDataplatformService,
				mockRuleService,
				mockFileImportService,
				mockTemporalService,
				mockCloudService,
				mockS3Client,
				serverConfig,
				mockCacheClient,
			)

			// Create context with logger
			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())

			// Call the method
			result, err := service.GetDatasetDisplayConfig(ctx, tt.merchantId, tt.datasetId)

			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			// Check result
			if tt.expectedResult != nil {
				assert.Equal(t, len(tt.expectedResult), len(result))

				// Create a map of columns for easier comparison
				expectedColumns := make(map[string]models.DisplayConfig)
				for _, config := range tt.expectedResult {
					expectedColumns[config.Column] = config
				}

				actualColumns := make(map[string]models.DisplayConfig)
				for _, config := range result {
					actualColumns[config.Column] = config
				}

				// Compare maps
				assert.Equal(t, len(expectedColumns), len(actualColumns))
				for column, expectedConfig := range expectedColumns {
					actualConfig, exists := actualColumns[column]
					assert.True(t, exists)
					assert.Equal(t, expectedConfig.IsHidden, actualConfig.IsHidden)
					assert.Equal(t, expectedConfig.IsEditable, actualConfig.IsEditable)
				}
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestGetAllDatasetUserPolicies(t *testing.T) {
	t.Parallel()

	datasetId := uuid.New()
	currentUserId := uuid.New()

	tests := []struct {
		name             string
		datasetId        uuid.UUID
		mockSetup        func(m *mockDatasetService.MockDatasetServiceStore)
		expectedErr      string
		expectedPolicies []storemodels.FlattenedResourceAudiencePolicy
	}{
		{
			name:      "Success",
			datasetId: datasetId,
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {
				expectedPolicies := []storemodels.FlattenedResourceAudiencePolicy{
					{
						ResourceId:   datasetId,
						ResourceType: string(storemodels.ResourceTypeDataset),
						UserId:       currentUserId,
						Privilege:    storemodels.PrivilegeDatasetAdmin,
					},
				}
				m.EXPECT().GetFlattenedResourceAudiencePolicies(mock.Anything, storemodels.FlattenedResourceAudiencePoliciesFilters{
					ResourceIds:   []uuid.UUID{datasetId},
					ResourceTypes: []string{string(storemodels.ResourceTypeDataset)},
				}).Return(expectedPolicies, nil)
			},
			expectedPolicies: []storemodels.FlattenedResourceAudiencePolicy{
				{
					ResourceId:   datasetId,
					ResourceType: string(storemodels.ResourceTypeDataset),
					UserId:       currentUserId,
					Privilege:    storemodels.PrivilegeDatasetAdmin,
				},
			},
		},
		{
			name:      "Error - failed to get policies",
			datasetId: datasetId,
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {
				m.EXPECT().GetFlattenedResourceAudiencePolicies(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("database error"))
			},
			expectedErr: "failed to get dataset policies: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDS := new(mockDatasetService.MockDatasetServiceStore)
			mockDPS := new(mockDataplatform.MockDataPlatformService)
			mockQueryBuilder := new(mock_querybuilder.MockQueryBuilder)
			mockRuleService := new(mockruleservice.MockRuleService)

			if tt.mockSetup != nil {
				tt.mockSetup(mockDS)
			}

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())

			service := datasetService{
				datasetStore:        mockDS,
				queryBuilderService: mockQueryBuilder,
				dataplatformService: mockDPS,
				ruleService:         mockRuleService,
			}

			// Execute
			policies, err := service.getFlattenedDatasetPolicies(ctx, tt.datasetId)

			// Assert
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, policies)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPolicies, policies)
		})
	}
}

func TestGetDatasetListing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		merchantId uuid.UUID
		params     models.DatsetListingParams
		mockSetup  func(*mockDatasetService.MockDatasetServiceStore)
		wantErr    error
	}{
		{
			name:       "Success case",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			params: models.DatsetListingParams{
				CreatedBy: []uuid.UUID{uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")},
				Pagination: &querybuildermodels.Pagination{
					Page:     1,
					PageSize: 10,
				},
				SortParams: []models.DatasetListingSortParams{
					{
						Column: "updated_at",
						Desc:   true,
					},
				},
			},
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {
				m.EXPECT().GetDatasetsAll(mock.Anything, storemodels.DatasetFilters{
					OrganizationIds: []uuid.UUID{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")},
					CreatedBy:       []uuid.UUID{uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")},
					Page:            1,
					Limit:           10,
					SortParams: []storemodels.DatasetSortParam{
						{
							Column: "updated_at",
							Desc:   true,
						},
					},
					Type: storemodels.UserVisibleDatasetTypes,
				}).Return([]storemodels.Dataset{}, nil)
			},
			wantErr: nil,
		},
		{
			name:       "Error - Invalid sort column",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			params: models.DatsetListingParams{
				SortParams: []models.DatasetListingSortParams{
					{
						Column: "invalid_column",
						Desc:   true,
					},
				},
			},
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {},
			wantErr:   errors.New("ERR_INVALID_DATALISTIN_SORT_COLUMN"),
		},
		{
			name:       "Error - Database error",
			merchantId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			params: models.DatsetListingParams{
				Pagination: &querybuildermodels.Pagination{
					Page:     1,
					PageSize: 10,
				},
			},
			mockSetup: func(m *mockDatasetService.MockDatasetServiceStore) {
				m.EXPECT().GetDatasetsAll(mock.Anything, storemodels.DatasetFilters{
					OrganizationIds: []uuid.UUID{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")},
					Page:            1,
					Limit:           10,
					Type:            storemodels.UserVisibleDatasetTypes,
				}).Return(nil, errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockStore := mockDatasetService.NewMockDatasetServiceStore(t)
			tt.mockSetup(mockStore)

			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			service := &datasetService{datasetStore: mockStore}

			// Execute
			_, err := service.GetDatasetListing(ctx, tt.merchantId, tt.params)

			// Assert
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGetDatasetImportPath(t *testing.T) {
	merchantId := uuid.New()
	datasetId := uuid.New()
	tests := []struct {
		name          string
		mockSetup     func(*mockDataplatform.MockDataPlatformService, *mock_s3.MockS3Client)
		expected      models.FileImportConfig
		wantErr       bool
		expectedError string
	}{
		{
			name: "Success case",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, mockS3Client *mock_s3.MockS3Client) {
				m.EXPECT().GetDags(mock.Anything, mock.Anything).Return(map[string]*servicemodels.DAGNode{
					datasetId.String(): {
						NodeId:   datasetId.String(),
						NodeType: servicemodels.NodeTypeDataset,
						Parents: []*servicemodels.DAGNode{
							{
								NodeId:   "s3a://bucket/path",
								NodeType: servicemodels.NodeTypeFolder,
							},
						},
					},
				}, nil)
				mockS3Client.EXPECT().GetSampleFilePathFromFolder(
					mock.Anything,
					"bucket",
					"path/",
				).Return("path/file1.parquet", nil)
			},
			expected: models.FileImportConfig{
				BronzeSourceBucket:  "bucket",
				BronzeSourcePath:    "path/file1.parquet",
				BronzeSourceConfig:  map[string]interface{}{},
				IsFileImportEnabled: true,
			},
			wantErr: false,
		},
		{
			name: "Error - failed to get dags",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, mockS3Client *mock_s3.MockS3Client) {
				m.EXPECT().GetDags(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to get dags"))
			},
			wantErr:       true,
			expectedError: "failed to get dags",
		},
		{
			name: "Error - dataset not found in dags",
			mockSetup: func(m *mockDataplatform.MockDataPlatformService, mockS3Client *mock_s3.MockS3Client) {
				m.EXPECT().GetDags(mock.Anything, mock.Anything).Return(map[string]*servicemodels.DAGNode{}, nil)
			},
			wantErr:       true,
			expectedError: "failed to get dataset dags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockDS := mock_store.NewMockStore(t)
			mockDPS := mockDataplatform.NewMockDataPlatformService(t)
			mockQueryBuilder := mock_querybuilder.NewMockQueryBuilder(t)
			mockRuleService := mockruleservice.NewMockRuleService(t)
			mockFileImportService := mockfileimports.NewMockFileImportService(t)
			mockTemporalService := mock_temporal.NewMockTemporalService(t)
			mockCloudService := mock_cloudservice.NewMockCloudService(t)
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockCacheClient := mock_cache.NewMockCacheClient(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockDPS, mockS3Client)
			}

			// Create service
			serverConfig := serverconfig.DatasetConfig{
				DataplatformProvider: "pinot",
			}
			service := NewDatasetService(
				mockDS,
				mockQueryBuilder,
				mockDPS,
				mockRuleService,
				mockFileImportService,
				mockTemporalService,
				mockCloudService,
				mockS3Client,
				serverConfig,
				mockCacheClient,
			)

			// Create context with logger
			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())

			// Call the method
			path, err := service.GetDatasetImportPath(ctx, merchantId, datasetId)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.BronzeSourceBucket, path.BronzeSourceBucket)
				assert.Equal(t, tt.expected.BronzeSourcePath, path.BronzeSourcePath)
				assert.Equal(t, tt.expected.BronzeSourceConfig, path.BronzeSourceConfig)
				assert.Equal(t, tt.expected.IsFileImportEnabled, path.IsFileImportEnabled)
			}

			// Assert mock expectations
			mockDPS.AssertExpectations(t)
		})
	}
}
