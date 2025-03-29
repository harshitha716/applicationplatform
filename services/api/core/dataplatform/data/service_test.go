package data

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	dataplatformconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	serviceconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	servicemodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	servicerrors "github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	mockrosetta "github.com/Zampfi/application-platform/services/api/mocks/core/dataplatform/rosetta"
	mockproviderregistry "github.com/Zampfi/application-platform/services/api/mocks/pkg/dataplatform/providers"
	mockprovider "github.com/Zampfi/application-platform/services/api/mocks/pkg/dataplatform/service"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	models "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func getDataPlatformMockConfig() *serverconfig.DataPlatformConfig {
	return &serverconfig.DataPlatformConfig{
		DatabricksConfig: serverconfig.DatabricksSetupConfig{
			MerchantDataProviderIdMapping: map[string]string{"merchant1": "workspace1"},
			DefaultDataProviderId:         "defaultWorkspace",
			DataProviderConfigs: map[string]models.DatabricksConfig{
				"workspace1": {
					WarehouseId: "warehouse1",
				},
			},
			ZampDatabricksCatalog:        "zamp",
			ZampDatabricksPlatformSchema: "platform",
		},
		PinotConfig: serverconfig.PinotSetupConfig{
			MerchantDataProviderIdMapping: map[string]string{"merchant2": "workspace2"},
			DefaultDataProviderId:         "defaultPinotWorkspace",
		},
		ActionsConfig: serverconfig.ActionsConfig{
			CreateMVJobTemplateConfig: serverconfig.CreateMVJobTemplateConfig{
				CreateMVNotebookPath:   "/path/to/create/mv/notebook",
				SideEffectNotebookPath: "/path/to/sideeffect/notebook",
			},
		},
	}
}

type DataServiceTestSuite struct {
	suite.Suite
	mockProviderService  *mockprovider.MockProviderService
	mockProviderRegistry *mockproviderregistry.MockProviderService
	service              *dataService
	mockRosettaService   *mockrosetta.MockRosettaService
}

func TestDataPlatformServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DataServiceTestSuite))
}

func (s *DataServiceTestSuite) SetupTest() {
	s.mockProviderService = new(mockprovider.MockProviderService)
	s.mockProviderRegistry = new(mockproviderregistry.MockProviderService)
	s.mockRosettaService = new(mockrosetta.MockRosettaService)
	s.service = &dataService{
		providerService:    s.mockProviderService,
		dataPlatformConfig: getDataPlatformMockConfig(),
		rosettaService:     s.mockRosettaService,
	}
}

func (s *DataServiceTestSuite) TestInitDataplatformService() {
	service, err := InitDataService(getDataPlatformMockConfig())
	s.NotNil(service)
	s.NoError(err)
}

func (s *DataServiceTestSuite) TestGetWorkspaceIdForMerchant() {
	tests := []struct {
		merchantId   string
		providerType constants.ProviderType
		expectedId   string
		expectError  bool
	}{
		{"merchant1", constants.ProviderTypeDatabricks, "workspace1", false},
		{"unknownMerchant", constants.ProviderTypeDatabricks, "defaultWorkspace", false},
		{"merchant2", constants.ProviderTypePinot, "workspace2", false},
		{"unknownMerchant", constants.ProviderTypePinot, "defaultPinotWorkspace", false},
		{"merchant1", constants.ProviderType("unknown"), "", true},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			id, err := s.service.GetDataProviderIdForMerchant(tt.merchantId, tt.providerType)
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedId, id)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestInternalQuery() {
	tests := []struct {
		merchantId          string
		datasetId           string
		query               string
		params              map[string]string
		providerType        constants.ProviderType
		rosettaMockResponse string
		mockDatasetResponse models.QueryResult
		expectedErr         bool
		err                 error
	}{
		{"merchant1", "dataset1", "SELECT * FROM table WHERE id = ?", map[string]string{"zamp_table_name_1": "dataset1"}, constants.ProviderTypeDatabricks, "SELECT * FROM \"dataset1\" WHERE id = ?", models.QueryResult{Rows: []map[string]interface{}{{"id": "1", "databricks_fq_table_name": "dataset1"}}}, false, nil},
		{"merchant2", "dataset2", "SELECT * FROM table WHERE id = ?", map[string]string{"zamp_table_name_2": "dataset2"}, constants.ProviderTypePinot, "SELECT * FROM \"dataset2\" WHERE id = ?", models.QueryResult{Rows: []map[string]interface{}{{"id": "2", "pinot_table_name": "dataset2"}}}, false, nil},
		{"unknownMerchant", "dataset1", "SELECT * FROM table WHERE id = ?", map[string]string{"id": "1"}, constants.ProviderTypeDatabricks, "", models.QueryResult{}, true, errors.New("error")},
		{"merchant3", "unknownDataset", "SELECT * FROM table WHERE id = ?", map[string]string{"id": "1"}, constants.ProviderTypeDatabricks, "", models.QueryResult{}, true, errors.New("error")},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, tt.providerType, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockRosettaService.On("TranslateQuery", ctx, tt.query, tt.providerType).Return(tt.rosettaMockResponse, nil).Once()
			s.mockProviderRegistry.On("Query", ctx, mock.Anything, fmt.Sprintf("SELECT %s FROM `zamp`.`platform`.`datasets` WHERE id = '%s' AND merchant_id = '%s' AND is_deleted = false", serviceconstants.SelectDatasetColumnNames, tt.datasetId, tt.merchantId)).Return(tt.mockDatasetResponse, tt.err).Once()
			if tt.providerType == constants.ProviderTypeDatabricks {
				s.mockProviderRegistry.On("Query", ctx, "\""+tt.datasetId+"\"", mock.Anything).Return(tt.mockDatasetResponse, tt.err).Once()
			} else {
				s.mockProviderRegistry.On("Query", ctx, "\""+tt.datasetId+"\"", mock.Anything).Return(tt.mockDatasetResponse, tt.err).Once()
			}
			result, err := s.service.query(ctx, tt.providerType, tt.merchantId, tt.query, tt.params)
			if tt.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.mockDatasetResponse, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetDataset() {
	tests := []struct {
		merchantId          string
		datasetId           string
		mockResponse        models.QueryResult
		expectedErr         bool
		err                 error
		expectedFunctionErr error
	}{
		{"merchant1", "dataset1", models.QueryResult{Rows: []map[string]interface{}{{"id": "1"}}}, false, nil, nil},
		{"merchant2", "dataset2", models.QueryResult{Rows: []map[string]interface{}{{"id": "2"}}}, false, nil, nil},
		{"merchant2", "dataset3", models.QueryResult{Rows: []map[string]interface{}{}}, true, nil, servicerrors.ErrDatasetNotFound},
		{"unknownMerchant", "dataset1", models.QueryResult{}, true, errors.New("error"), servicerrors.ErrQueryingDatabricksFailed},
		{"merchant1", "unknownDataset", models.QueryResult{}, true, errors.New("error"), servicerrors.ErrQueryingDatabricksFailed},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()

			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.err).Once()

			result, err := s.service.getDataset(ctx, tt.merchantId, tt.datasetId)

			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.expectedFunctionErr, err)
			} else {
				s.NoError(err)
				s.NotNil(result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetDatasetConfig() {
	tests := []struct {
		merchantId          string
		datasetId           string
		mockResponse        models.QueryResult
		expectedErr         bool
		expectedFunctionErr error
		expectedConfig      servicemodels.DatasetConfig
	}{
		{"merchant1", "dataset1", models.QueryResult{Rows: []map[string]interface{}{{"id": "1", "dataset_config": `{"columns": {"tags": {"custom_type": "tags"}}}`}}}, false, nil, servicemodels.DatasetConfig{Columns: map[string]servicemodels.DatasetColumnConfig{"tags": {CustomType: dataplatformconstants.DatabricksColumnCustomTypeTags}}}},
		{"merchant2", "dataset2", models.QueryResult{Rows: []map[string]interface{}{{"id": "1"}}}, false, nil, servicemodels.DatasetConfig{}},
		{"merchant1", "dataset3", models.QueryResult{Rows: []map[string]interface{}{{"id": "1", "dataset_config": `wrong_json}`}}}, true, servicerrors.ErrJSONUnmarshallingFailed, servicemodels.DatasetConfig{}},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()
			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, nil).Once()

			result, err := s.service.GetDatasetConfig(ctx, tt.merchantId, tt.datasetId)
			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.expectedFunctionErr, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedConfig, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetDatabricksDatasetMetadata() {
	tests := []struct {
		datasetInfo         servicemodels.Dataset
		mockResponse        models.QueryResult
		expectedErr         bool
		expectedMetadata    servicemodels.InternalDatasetMetadata
		expectedFunctionErr error
	}{
		{
			datasetInfo: servicemodels.Dataset{
				DatabricksStats:  `{"columns": {"id": {"distinct_count": 100}}}`,
				DatabricksSchema: `{"columns": {"id": {"type": "string"}}}`,
			},
			mockResponse: models.QueryResult{Rows: []map[string]interface{}{{"id": "1"}}},
			expectedErr:  false,
			expectedMetadata: servicemodels.InternalDatasetMetadata{
				Schema: servicemodels.DatasetSchemaDetails{
					Columns: map[string]servicemodels.ColumnMetadata{
						"id": {Type: "string"},
					},
				},
				Stats: servicemodels.DatasetStats{
					ColumnStats: map[string]servicemodels.ColumnStats{
						"id": {DistinctCount: 100},
					},
				},
			},
			expectedFunctionErr: nil,
		},
		{
			datasetInfo: servicemodels.Dataset{
				DatabricksStats:  `{"countx": 100}`,
				DatabricksSchema: `{"columns": {"id": {"type": "string"}}}`,
			},
			mockResponse: models.QueryResult{},
			expectedErr:  false,
			expectedMetadata: servicemodels.InternalDatasetMetadata{
				Schema: servicemodels.DatasetSchemaDetails{
					Columns: map[string]servicemodels.ColumnMetadata{
						"id": {Type: "string"},
					},
				},
				Stats: servicemodels.DatasetStats{},
			},
			expectedFunctionErr: servicerrors.ErrJSONUnmarshallingFailed,
		},
	}

	for _, tt := range tests {
		s.Run(tt.datasetInfo.DatabricksStats, func() {
			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()

			ctx := context.Background()
			result, err := s.service.getDatabricksDatasetMetadata(ctx, tt.datasetInfo)
			if tt.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedMetadata, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetPinotDatasetMetadata() {
	tests := []struct {
		datasetInfo         servicemodels.Dataset
		mockResponse        models.QueryResult
		expectedErr         bool
		expectedMetadata    servicemodels.InternalDatasetMetadata
		expectedFunctionErr error
	}{
		{
			datasetInfo: servicemodels.Dataset{
				PinotStats:  `{"columns": {"id": {"distinct_count": 100}}}`,
				PinotSchema: `{"columns": {"id": {"type": "string"}}}`,
			},
			mockResponse: models.QueryResult{Rows: []map[string]interface{}{{"id": "1"}}},
			expectedErr:  false,
			expectedMetadata: servicemodels.InternalDatasetMetadata{
				Schema: servicemodels.DatasetSchemaDetails{
					Columns: map[string]servicemodels.ColumnMetadata{
						"id": {Type: "string"},
					},
				},
				Stats: servicemodels.DatasetStats{
					ColumnStats: map[string]servicemodels.ColumnStats{
						"id": {DistinctCount: 100},
					},
				},
			},
			expectedFunctionErr: nil,
		},
		{
			datasetInfo: servicemodels.Dataset{
				PinotStats:  `{"countx": 100}`,
				PinotSchema: `{"columns": {"id": {"type": "string"}}}`,
			},
			mockResponse: models.QueryResult{},
			expectedErr:  false,
			expectedMetadata: servicemodels.InternalDatasetMetadata{
				Schema: servicemodels.DatasetSchemaDetails{
					Columns: map[string]servicemodels.ColumnMetadata{
						"id": {Type: "string"},
					},
				},
				Stats: servicemodels.DatasetStats{},
			},
			expectedFunctionErr: servicerrors.ErrJSONUnmarshallingFailed,
		},
	}

	for _, tt := range tests {
		s.Run(tt.datasetInfo.PinotStats, func() {
			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()

			ctx := context.Background()
			result, err := s.service.getPinotDatasetMetadata(ctx, tt.datasetInfo)
			if tt.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedMetadata, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetDatasetMetadata() {
	tests := []struct {
		merchantId          string
		datasetId           string
		expectedErr         bool
		expectedMetadata    servicemodels.DatasetMetadata
		queryResult         models.QueryResult
		expectedFunctionErr error
	}{
		{
			merchantId:  "merchant1",
			datasetId:   "dataset1",
			expectedErr: false,
			expectedMetadata: servicemodels.DatasetMetadata{
				Schema: map[string]servicemodels.ColumnMetadata{
					"id": {Type: "string"},
				},
				Stats: servicemodels.DatasetStats{
					ColumnStats: map[string]servicemodels.ColumnStats{
						"id": {DistinctCount: 100},
					},
				},
			},
			expectedFunctionErr: nil,
			queryResult: models.QueryResult{Rows: []map[string]interface{}{{
				"id":                    "dataset1",
				"merchant_id":           "merchant1",
				"databricks_table_name": "dataset1",
				"databricks_schema":     `{"columns": {"id": {"type": "string"}}}`,
				"databricks_config":     `{}`,
				"databricks_stats":      `{"columns": {"id": {"distinct_count": 100}}}`,
				"pinot_table_name":      "dataset1",
				"pinot_schema":          `{"columns": {"id": {"type": "string32"}}}`,
				"pinot_config":          `{}`,
				"pinot_stats":           `{"columns": {"id": {"distinct_count": 100}}}`,
				"created_at":            time.Time{},
				"updated_at":            time.Time{},
				"is_deleted":            false,
				"deleted_at":            time.Time{},
			}}},
		},
		{
			merchantId:  "merchant2",
			datasetId:   "dataset2",
			expectedErr: false,
			expectedMetadata: servicemodels.DatasetMetadata{
				Schema: nil,
				Stats:  servicemodels.DatasetStats{},
			},
			expectedFunctionErr: servicerrors.ErrJSONUnmarshallingFailed,
			queryResult: models.QueryResult{Rows: []map[string]interface{}{{
				"id":                    "dataset1",
				"merchant_id":           "merchant1",
				"databricks_table_name": "dataset1",
				"databricks_schema":     `{}`,
				"databricks_config":     `{}`,
				"databricks_stats":      `{}`,
				"pinot_table_name":      "dataset1",
				"pinot_schema":          `{}`,
				"pinot_config":          `{}`,
				"pinot_stats":           `{}`,
				"created_at":            time.Time{},
				"updated_at":            time.Time{},
				"is_deleted":            false,
				"deleted_at":            time.Time{},
			}}},
		},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()
			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.queryResult, nil).Once()

			result, err := s.service.GetDatasetMetadata(ctx, tt.merchantId, tt.datasetId)
			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.expectedFunctionErr, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedMetadata, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetDatasetParents() {
	tests := []struct {
		merchantId          string
		datasetId           string
		mockQueryResult     models.QueryResult
		expectedErr         bool
		expectedParents     servicemodels.DatasetParents
		expectedFunctionErr error
	}{
		{
			merchantId: "merchant1",
			datasetId:  "dataset1",
			mockQueryResult: models.QueryResult{
				Rows: []map[string]interface{}{
					{"id": "parent1", "source_type": "type1", "source_value": "value1"},
					{"id": "parent2", "source_type": "type2", "source_value": "value2"},
				},
			},
			expectedErr: false,
			expectedParents: servicemodels.DatasetParents{
				Parents: []servicemodels.DAGNode{
					{Id: "value1", Type: "type1"},
					{Id: "value2", Type: "type2"},
				},
			},
			expectedFunctionErr: nil,
		},
		{
			merchantId: "merchant2",
			datasetId:  "dataset2",
			mockQueryResult: models.QueryResult{
				Rows: []map[string]interface{}{},
			},
			expectedErr:         true,
			expectedParents:     servicemodels.DatasetParents{},
			expectedFunctionErr: servicerrors.ErrDatasetNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()

			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.mockQueryResult, nil).Once()
			result, err := s.service.GetDatasetParents(ctx, tt.merchantId, tt.datasetId)
			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.expectedFunctionErr, err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedParents, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestQueryRealTime() {
	tests := []struct {
		merchantId          string
		datasetId           string
		query               string
		params              map[string]string
		mockDatasetResponse models.QueryResult
		rosettaInputQuery   string
		rosettaMockResponse string
		mockResponse        models.QueryResult
		expectedErr         bool
		err                 error
	}{
		{
			merchantId: "merchant1",
			datasetId:  "dataset1",
			query:      "SELECT * FROM {{.zamp_table_name_1}} WHERE id = ?",
			params:     map[string]string{"zamp_table_name_1": "dataset1"},
			mockDatasetResponse: models.QueryResult{
				Rows: []map[string]interface{}{{"id": "1", "pinot_table_name": "dataset1"}},
			},
			rosettaInputQuery:   "SELECT * FROM \"dataset1\" WHERE id = ?",
			rosettaMockResponse: "SELECT * FROM \"dataset1\" WHERE id = ?",
			mockResponse:        models.QueryResult{Rows: []map[string]interface{}{{"id": "1"}}},
			expectedErr:         false,
			err:                 nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypePinot, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockRosettaService.On("TranslateQuery", ctx, tt.rosettaInputQuery, constants.ProviderTypePinot).Return(tt.rosettaMockResponse, nil).Once()
			s.mockProviderRegistry.On("Query", ctx, mock.Anything, fmt.Sprintf("SELECT %s FROM `zamp`.`platform`.`datasets` WHERE id = '%s' AND merchant_id = '%s' AND is_deleted = false", serviceconstants.SelectDatasetColumnNames, tt.datasetId, tt.merchantId)).Return(tt.mockDatasetResponse, tt.err).Once()
			s.mockProviderRegistry.On("Query", ctx, "\""+tt.datasetId+"\"", mock.Anything).Return(tt.mockDatasetResponse, tt.err).Once()
			result, err := s.service.QueryRealTime(ctx, tt.merchantId, tt.query, tt.params)

			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.err, err)
			} else {
				s.NoError(err)
				s.Equal(tt.mockDatasetResponse, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestQueryRealTimeFallback() {
	tests := []struct {
		name                          string
		merchantId                    string
		datasetId                     string
		query                         string
		params                        map[string]string
		mockDatasetResponse           models.QueryResult
		rosettaInputQueryPinot        string
		rosettaMockResponsePinot      string
		rosettaInputQueryDatabricks   string
		rosettaMockResponseDatabricks string
		mockResponse                  models.QueryResult
		pinotQueryErr                 error
		expectedErr                   bool
		err                           error
	}{
		{
			name:       "Fallback to Databricks as Pinot table name not found",
			merchantId: "merchant1",
			datasetId:  "dataset1",
			query:      "SELECT * FROM {{.zamp_table_name_1}} WHERE id = ?",
			params:     map[string]string{"zamp_table_name_1": "dataset1"},
			mockDatasetResponse: models.QueryResult{
				Rows: []map[string]interface{}{{"id": "1", "databricks_fq_table_name": "dataset1"}},
			},
			rosettaInputQueryPinot:        "SELECT * FROM \"dataset1\" WHERE id = ?",
			rosettaMockResponsePinot:      "SELECT * FROM \"dataset1\" WHERE id = ?",
			rosettaInputQueryDatabricks:   "SELECT * FROM \"dataset1\" WHERE id = ?",
			rosettaMockResponseDatabricks: "SELECT * FROM `dataset1` WHERE id = ?",
			pinotQueryErr:                 servicerrors.ErrQueryingPinotFailed,
			expectedErr:                   false,
			err:                           nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypePinot, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockRosettaService.On("TranslateQuery", ctx, tt.rosettaInputQueryPinot, constants.ProviderTypePinot).Return(tt.rosettaMockResponsePinot, nil).Once()
			s.mockRosettaService.On("TranslateQuery", ctx, tt.rosettaInputQueryDatabricks, constants.ProviderTypeDatabricks).Return(tt.rosettaMockResponseDatabricks, nil).Once()
			s.mockProviderRegistry.On("Query", ctx, mock.Anything, fmt.Sprintf("SELECT %s FROM `zamp`.`platform`.`datasets` WHERE id = '%s' AND merchant_id = '%s' AND is_deleted = false", serviceconstants.SelectDatasetColumnNames, tt.datasetId, tt.merchantId)).Return(tt.mockDatasetResponse, tt.err)
			s.mockProviderRegistry.On("Query", ctx, tt.datasetId, mock.Anything).Return(tt.mockDatasetResponse, tt.pinotQueryErr).Once()
			s.mockProviderRegistry.On("Query", ctx, "\""+tt.datasetId+"\"", mock.Anything).Return(tt.mockDatasetResponse, tt.err).Once()
			result, err := s.service.QueryRealTime(ctx, tt.merchantId, tt.query, tt.params)

			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.err, err)
			} else {
				s.NoError(err)
				s.Equal(tt.mockDatasetResponse, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestQueryRealTimeError() {
	tests := []struct {
		name                string
		merchantId          string
		datasetId           string
		query               string
		params              map[string]string
		mockDatasetResponse models.QueryResult
		rosettaInputQuery   string
		rosettaMockResponse string
		expectedErr         bool
		err                 error
	}{
		{
			name:       "Fallback also fails as table names not found",
			merchantId: "merchant1",
			datasetId:  "dataset1",
			query:      "SELECT * FROM {{.zamp_table_name_1}} WHERE id = ?",
			params:     map[string]string{"zamp_table_name_1": "dataset1"},
			mockDatasetResponse: models.QueryResult{
				Rows: []map[string]interface{}{},
			},
			expectedErr: true,
			err:         servicerrors.ErrProcessingParamsForQueryFailed,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypePinot, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockProviderRegistry.On("Query", ctx, mock.Anything, fmt.Sprintf("SELECT %s FROM `zamp`.`platform`.`datasets` WHERE id = '%s' AND merchant_id = '%s' AND is_deleted = false", serviceconstants.SelectDatasetColumnNames, tt.datasetId, tt.merchantId)).Return(tt.mockDatasetResponse, tt.err)
			s.mockProviderRegistry.On("Query", ctx, tt.datasetId, tt.rosettaMockResponse).Return(tt.mockDatasetResponse, tt.err).Once()
			result, err := s.service.QueryRealTime(ctx, tt.merchantId, tt.query, tt.params)

			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.err, err)
			} else {
				s.NoError(err)
				s.Equal(tt.mockDatasetResponse, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestQuery() {
	tests := []struct {
		merchantId          string
		datasetId           string
		query               string
		params              map[string]string
		mockDatasetResponse models.QueryResult
		rosettaInputQuery   string
		rosettaMockResponse string
		expectedErr         bool
		err                 error
	}{
		{
			merchantId: "merchant1",
			datasetId:  "dataset1",
			query:      "SELECT * FROM {{.zamp_table_name_1}} WHERE id = ?",
			params:     map[string]string{"zamp_table_name_1": "dataset1"},
			mockDatasetResponse: models.QueryResult{
				Rows: []map[string]interface{}{{"id": "1", "databricks_fq_table_name": "dataset1"}},
			},
			rosettaInputQuery:   "SELECT * FROM \"dataset1\" WHERE id = ?",
			rosettaMockResponse: "SELECT * FROM `dataset1` WHERE id = ?",
			expectedErr:         false,
			err:                 nil,
		},
		{
			merchantId: "merchant2",
			datasetId:  "dataset2",
			query:      "SELECT * FROM {{.zamp_table_name_1}} WHERE id = ?",
			params:     map[string]string{"zamp_table_name_1": "dataset2"},
			mockDatasetResponse: models.QueryResult{
				Rows: []map[string]interface{}{},
			},
			rosettaInputQuery:   "SELECT * FROM \"dataset2\" WHERE id = ?",
			rosettaMockResponse: "SELECT * FROM `dataset2` WHERE id = ?",
			expectedErr:         true,
			err:                 servicerrors.ErrProcessingParamsForQueryFailed,
		},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil)
			s.mockRosettaService.On("TranslateQuery", ctx, tt.rosettaInputQuery, constants.ProviderTypeDatabricks).Return(tt.rosettaMockResponse, nil).Once()

			s.mockProviderRegistry.On("Query", ctx, mock.Anything, fmt.Sprintf("SELECT %s FROM `zamp`.`platform`.`datasets` WHERE id = '%s' AND merchant_id = '%s' AND is_deleted = false", serviceconstants.SelectDatasetColumnNames, tt.datasetId, tt.merchantId)).Return(tt.mockDatasetResponse, tt.err).Once()

			s.mockProviderRegistry.On("Query", ctx, "\""+tt.datasetId+"\"", mock.Anything).Return(tt.mockDatasetResponse, tt.err).Once()

			result, err := s.service.Query(ctx, tt.merchantId, tt.query, tt.params)

			if tt.expectedErr {
				s.Error(err)
				s.Equal(tt.err, err)
			} else {
				s.NoError(err)
				s.Equal(tt.mockDatasetResponse, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestProcessParamsForQuery() {
	tests := []struct {
		merchantId         string
		params             map[string]string
		providerType       constants.ProviderType
		expectedTableNames []string
		expectedParams     map[string]string
		expectError        bool
	}{
		{
			merchantId: "merchant1",
			params: map[string]string{
				"zamp_table_name_1": "dataset1",
			},
			providerType:       constants.ProviderTypeDatabricks,
			expectedTableNames: []string{"\"databricks_table_name\""},
			expectedParams:     map[string]string{"zamp_table_name_1": "\"databricks_table_name\""},
			expectError:        false,
		},
		{
			merchantId: "merchant2",
			params: map[string]string{
				"zamp_table_name_2": "dataset2",
			},
			providerType:       constants.ProviderTypePinot,
			expectedTableNames: []string{"\"pinot_table_name\""},
			expectedParams:     map[string]string{"zamp_table_name_2": "\"pinot_table_name\""},
			expectError:        false,
		},
		{
			merchantId: "unknownMerchant",
			params: map[string]string{
				"zamp_table_name_3": "unknownDataset",
			},
			providerType:       constants.ProviderTypeDatabricks,
			expectedTableNames: []string{},
			expectedParams:     map[string]string{},
			expectError:        true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()
			if !tt.expectError {
				s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.QueryResult{Rows: []map[string]interface{}{{"databricks_fq_table_name": "databricks_table_name", "pinot_table_name": "pinot_table_name"}}}, nil).Once()
			} else {
				s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.QueryResult{Rows: []map[string]interface{}{}}, servicerrors.ErrDatasetNotFound).Once()
			}

			queryMetadata, err := s.service.ProcessParamsForQuery(ctx, tt.merchantId, tt.params, tt.providerType)

			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedTableNames, queryMetadata.TableNames)
				s.Equal(tt.expectedParams, queryMetadata.Params)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetProviderLevelDatasetMetadata() {
	tests := []struct {
		merchantId          string
		datasetId           string
		mockResponse        models.QueryResult
		expectedErr         bool
		expectedMetadata    servicemodels.ProviderLevelDatasetMetadata
		expectedFunctionErr error
	}{
		{
			merchantId:   "merchant1",
			datasetId:    "dataset1",
			mockResponse: models.QueryResult{Rows: []map[string]interface{}{{"id": "1", "pinot_table_name": "dataset1", "databricks_table_name": "dataset1", "pinot_schema": `{"columns": {"id": {"type": "string"}}}`, "databricks_schema": `{"columns": {"id": {"type": "string"}}}`, "pinot_stats": `{"columns": {"id": {"distinct_count": 100}}}`, "databricks_stats": `{"columns": {"id": {"distinct_count": 100}}}`}}},
			expectedErr:  false,
			expectedMetadata: servicemodels.ProviderLevelDatasetMetadata{
				Pinot: servicemodels.InternalDatasetMetadata{
					Schema: servicemodels.DatasetSchemaDetails{
						Columns: map[string]servicemodels.ColumnMetadata{
							"id": {Type: "string"},
						},
					},
					Stats: servicemodels.DatasetStats{
						ColumnStats: map[string]servicemodels.ColumnStats{
							"id": {DistinctCount: 100},
						},
					},
				},
				Databricks: servicemodels.InternalDatasetMetadata{
					Schema: servicemodels.DatasetSchemaDetails{
						Columns: map[string]servicemodels.ColumnMetadata{
							"id": {Type: "string"},
						},
					},
					Stats: servicemodels.DatasetStats{
						ColumnStats: map[string]servicemodels.ColumnStats{
							"id": {DistinctCount: 100},
						},
					},
				},
			},
			expectedFunctionErr: nil,
		},
		{
			merchantId: "merchant2",
			datasetId:  "dataset2",
			mockResponse: models.QueryResult{
				Rows: []map[string]interface{}{},
			},
			expectedErr: true,
			expectedMetadata: servicemodels.ProviderLevelDatasetMetadata{
				Pinot: servicemodels.InternalDatasetMetadata{
					Schema: servicemodels.DatasetSchemaDetails{
						Columns: map[string]servicemodels.ColumnMetadata{
							"id": {Type: "string"},
						},
					},
					Stats: servicemodels.DatasetStats{},
				},
				Databricks: servicemodels.InternalDatasetMetadata{},
			},
			expectedFunctionErr: servicerrors.ErrJSONUnmarshallingFailed,
		},
	}

	for _, tt := range tests {
		s.Run(tt.merchantId, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()

			s.mockProviderRegistry.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, nil).Once()
			result, err := s.service.getProviderLevelDatasetMetadata(ctx, tt.merchantId, tt.datasetId)
			if tt.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedMetadata, result)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestGetDatabricksWarehouseId() {
	tests := []struct {
		name           string
		dataProviderId string
		mockConfig     *serverconfig.DataPlatformConfig
		expectedId     string
		expectError    bool
	}{
		{
			name:           "Valid provider ID",
			dataProviderId: "workspace1",
			mockConfig: &serverconfig.DataPlatformConfig{
				DatabricksConfig: serverconfig.DatabricksSetupConfig{
					ZampDatabricksCatalog:        "zamp",
					ZampDatabricksPlatformSchema: "platform",
					DataProviderConfigs: map[string]models.DatabricksConfig{
						"workspace1": {
							WarehouseId: "warehouse1",
						},
					},
				},
			},
			expectedId:  "warehouse1",
			expectError: false,
		},
		{
			name:           "Invalid provider ID",
			dataProviderId: "unknownWorkspace",
			mockConfig: &serverconfig.DataPlatformConfig{
				DatabricksConfig: serverconfig.DatabricksSetupConfig{
					ZampDatabricksCatalog:        "zamp",
					ZampDatabricksPlatformSchema: "platform",
					DataProviderConfigs: map[string]models.DatabricksConfig{
						"workspace1": {
							WarehouseId: "warehouse1",
						},
					},
				},
			},
			expectedId:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.service.dataPlatformConfig = tt.mockConfig
			ctx := context.Background()
			warehouseId, err := s.service.GetDatabricksWarehouseId(ctx, tt.dataProviderId)

			if tt.expectError {
				s.Error(err)
				s.Equal("", warehouseId)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedId, warehouseId)
			}
		})
	}
}

func (s *DataServiceTestSuite) TestParseDatabricksFQTableName() {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"database.table", "\"database\".\"table\""},
		{"database-with-dash.table", "\"database-with-dash\".\"table\""},
		{"database.table-with-dash", "\"database\".\"table-with-dash\""},
		{"database.with.dots.table", "\"database\".\"with\".\"dots\".\"table\""},
		{"table", "\"table\""},
	}

	for _, tt := range tests {
		s.Run(tt.input, func() {
			output := parseDatabricksFQTableName(tt.input)
			s.Equal(tt.expectedOutput, output)
		})
	}
}

func (s *DataServiceTestSuite) TestParsePinotTableName() {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"table", "\"table\""},
	}

	for _, tt := range tests {
		s.Run(tt.input, func() {
			output := parsePinotTableName(tt.input)
			s.Equal(tt.expectedOutput, output)
		})
	}
}

func (s *DataServiceTestSuite) TestGetDatasetEdgesByMerchant() {
	tests := []struct {
		name              string
		merchantId        string
		mockQueryResult   models.QueryResult
		expectedErr       bool
		expectedMappings  []servicemodels.JobDatasetMapping
		mockProviderError error
	}{
		{
			name:       "Success - Returns dataset edges",
			merchantId: "merchant1",
			mockQueryResult: models.QueryResult{
				Rows: []map[string]interface{}{
					{
						"id":                "mapping1",
						"merchant_id":       "merchant1",
						"source_type":       "dataset",
						"source_value":      "dataset1",
						"destination_type":  "dataset",
						"destination_value": "dataset2",
					},
					{
						"id":                "mapping2",
						"merchant_id":       "merchant1",
						"source_type":       "dataset",
						"source_value":      "dataset2",
						"destination_type":  "dataset",
						"destination_value": "dataset3",
					},
				},
			},
			expectedErr: false,
			expectedMappings: []servicemodels.JobDatasetMapping{
				{
					Id:               "mapping1",
					MerchantId:       "merchant1",
					SourceType:       "dataset",
					SourceValue:      "dataset1",
					DestinationType:  "dataset",
					DestinationValue: "dataset2",
				},
				{
					Id:               "mapping2",
					MerchantId:       "merchant1",
					SourceType:       "dataset",
					SourceValue:      "dataset2",
					DestinationType:  "dataset",
					DestinationValue: "dataset3",
				},
			},
			mockProviderError: nil,
		},
		{
			name:       "Error - Provider service query fails",
			merchantId: "merchant1",
			mockQueryResult: models.QueryResult{
				Rows: nil,
			},
			expectedErr:       true,
			expectedMappings:  nil,
			mockProviderError: errors.New("provider service error"),
		},
		{
			name:       "Success - Empty result",
			merchantId: "merchant1",
			mockQueryResult: models.QueryResult{
				Rows: []map[string]interface{}{},
			},
			expectedErr:       false,
			expectedMappings:  []servicemodels.JobDatasetMapping{},
			mockProviderError: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockProviderService.On("GetService", ctx, constants.ProviderTypeDatabricks, mock.Anything).Return(s.mockProviderRegistry, nil).Once()
			s.mockProviderRegistry.On("Query", ctx, mock.Anything, mock.Anything).Return(tt.mockQueryResult, tt.mockProviderError).Once()

			result, err := s.service.GetDatasetEdgesByMerchant(ctx, tt.merchantId)

			if tt.expectedErr {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedMappings, result)
			}
		})
	}
}
