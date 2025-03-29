package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/rosetta"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	serviceconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	servicemodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	models "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	provider "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/databricks"
	dataplatformservice "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/service"
	"go.uber.org/zap"
)

type DataService interface {
	QueryRealTime(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error)
	Query(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error)
	GetDatasetMetadata(ctx context.Context, merchantId string, datasetId string) (servicemodels.DatasetMetadata, error)
	GetDatasetParents(ctx context.Context, merchantId string, datasetId string) (servicemodels.DatasetParents, error)
	GetDatabricksWarehouseId(ctx context.Context, providerId string) (string, error)
	GetDatabricksServiceForMerchant(ctx context.Context, merchantId string) (databricks.DatabricksService, error)
	GetDatabricksServiceForProvider(ctx context.Context, providerId string) (databricks.DatabricksService, error)
	ProcessParamsForQuery(ctx context.Context, merchantId string, params map[string]string, providerType constants.ProviderType) (servicemodels.QueryMetadata, error)
	GetDataProviderIdForMerchant(merchantId string, providerType constants.ProviderType) (string, error)
	GetDataPlatformConfig() *serverconfig.DataPlatformConfig
	GetDatasetConfig(ctx context.Context, merchantId string, datasetId string) (servicemodels.DatasetConfig, error)
	TranslateQuery(ctx context.Context, query string, providerType constants.ProviderType) (string, error)
	GetDatasetEdgesByMerchant(ctx context.Context, merchantId string) ([]servicemodels.JobDatasetMapping, error)
}

type dataService struct {
	providerService    dataplatformservice.ProviderService
	dataPlatformConfig *serverconfig.DataPlatformConfig
	rosettaService     rosetta.RosettaService
}

// TODO: ADD PROPER FALLBACK FLOWS FOR REINIT AND ETC ..
func InitDataService(dataPlatformConfig *serverconfig.DataPlatformConfig) (DataService, error) {
	providerConfigs := []models.ProviderConfig{}
	for dataProviderId, dataProviderConfig := range dataPlatformConfig.DatabricksConfig.DataProviderConfigs {
		providerConfigs = append(providerConfigs, models.ProviderConfig{
			DataProviderId: dataProviderId,
			Provider:       constants.ProviderTypeDatabricks,
			Config:         dataProviderConfig,
		})
	}
	for dataProviderId, dataProviderConfig := range dataPlatformConfig.PinotConfig.DataProviderConfigs {
		providerConfigs = append(providerConfigs, models.ProviderConfig{
			DataProviderId: dataProviderId,
			Provider:       constants.ProviderTypePinot,
			Config:         dataProviderConfig,
		})
	}
	providerService, err := dataplatformservice.InitProviders(providerConfigs)
	if err != nil {
		return nil, err
	}
	rosettaService := rosetta.InitRosettaService(dataPlatformConfig.RosettaBaseUrl)
	return &dataService{
		providerService:    providerService,
		dataPlatformConfig: dataPlatformConfig,
		rosettaService:     rosettaService,
	}, nil
}

func (s *dataService) GetDataPlatformConfig() *serverconfig.DataPlatformConfig {
	return s.dataPlatformConfig
}

func (s *dataService) GetDataProviderIdForMerchant(merchantId string, providerType constants.ProviderType) (string, error) {
	if providerType == constants.ProviderTypeDatabricks {
		dataProviderId, ok := s.dataPlatformConfig.DatabricksConfig.MerchantDataProviderIdMapping[merchantId]
		if !ok {
			return s.dataPlatformConfig.DatabricksConfig.DefaultDataProviderId, nil
		}
		return dataProviderId, nil
	} else if providerType == constants.ProviderTypePinot {
		dataProviderId, ok := s.dataPlatformConfig.PinotConfig.MerchantDataProviderIdMapping[merchantId]
		if !ok {
			return s.dataPlatformConfig.PinotConfig.DefaultDataProviderId, nil
		}
		return dataProviderId, nil
	}
	return "", errors.ErrUnsupportedProviderType
}

func (s *dataService) getProviderService(ctx context.Context, merchantId string, providerType constants.ProviderType) (provider.ProviderService, error) {
	dataProviderId, err := s.GetDataProviderIdForMerchant(merchantId, providerType)
	if err != nil {
		return nil, err
	}
	return s.providerService.GetService(ctx, providerType, dataProviderId)
}

func (s *dataService) TranslateQuery(ctx context.Context, query string, providerType constants.ProviderType) (string, error) {
	return s.rosettaService.TranslateQuery(ctx, query, providerType)
}

func (s *dataService) query(ctx context.Context, providerType constants.ProviderType, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	traceId := apicontext.GetTraceIdFromContext(ctx)

	providerService, err := s.getProviderService(ctx, merchantId, providerType)
	if err != nil {
		return models.QueryResult{}, err
	}

	queryMetadata, err := s.ProcessParamsForQuery(ctx, merchantId, params, providerType)
	if err != nil {
		logger.Error(errors.ProcessingParamsForQueryFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, errors.ErrProcessingParamsForQueryFailed
	}

	filledQuery, err := helpers.FillQueryTemplate(ctx, query, queryMetadata.Params)
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, errors.ErrTemplateParsingFailed
	}

	filledQuery, err = s.rosettaService.TranslateQuery(ctx, filledQuery, providerType)
	if err != nil {
		logger.Error(errors.QueryTranslationFailedErrMessage, zap.Error(err))
	}

	filledQuery = helpers.AddCommentsToQuery(filledQuery, map[string]string{
		serviceconstants.MerchantIdQueryMetadataKey: merchantId,
		serviceconstants.TraceIdQueryMetadataKey:    traceId,
	})

	// THIS CONDITION IS ALWAYS TRUE IN ZAMP'S CASE AS WE NEED TO HAVE ONE TABLE NAME IN THE QUERY
	if len(params) > 0 && len(queryMetadata.TableNames) == 0 {
		logger.Error(errors.NoTableNamesFoundErrMessage)
		return models.QueryResult{}, errors.ErrNoTableNamesFound
	}

	var tableName string
	if len(queryMetadata.TableNames) > 0 {
		tableName = queryMetadata.TableNames[0]
	}

	logger.Info("filledQuery: ", zap.String("filledQuery", filledQuery))

	return providerService.Query(ctx, tableName, filledQuery, args...)
}

func (s *dataService) QueryRealTime(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error) {
	startTime := time.Now()
	logger := apicontext.GetLoggerFromCtx(ctx)

	pinotResult, err := s.query(ctx, constants.ProviderTypePinot, merchantId, query, params, args...)
	if err == nil {
		logger.Info("SUCCESSFULLY_EXECUTED_REAL_TIME_QUERY_VIA_PINOT", zap.Any("REAL_TIME_QUERY_PINOT_TIME_MS", time.Since(startTime).Milliseconds()))
		return pinotResult, nil
	}

	logger.Error("QUERYING_DATABRICKS_AS_PINOT_FAILED", zap.Error(err))

	// Fallback to databricks if pinot query fails
	databricksResult, err := s.query(ctx, constants.ProviderTypeDatabricks, merchantId, query, params, args...)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, err
	}

	logger.Info("SUCCESSFULLY_EXECUTED_REAL_TIME_QUERY_VIA_DATABRICKS", zap.Any("REAL_TIME_QUERY_DATABRICKS_TIME_MS", time.Since(startTime).Milliseconds()))
	return databricksResult, nil
}

func (s *dataService) Query(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error) {
	startTime := time.Now()
	logger := apicontext.GetLoggerFromCtx(ctx)

	databricksResult, err := s.query(ctx, constants.ProviderTypeDatabricks, merchantId, query, params, args...)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, err
	}

	logger.Info("SUCCESSFULLY_EXECUTED_QUERY", zap.Any("QUERY_TIME_MS", time.Since(startTime).Milliseconds()))
	return databricksResult, nil
}

func parseDatabricksFQTableName(databricksFQTableName string) string {
	parts := strings.Split(databricksFQTableName, ".")
	for i, part := range parts {
		parts[i] = fmt.Sprintf("\"%s\"", part)
	}

	return strings.Join(parts, ".")
}

func parsePinotTableName(pinotTableName string) string {
	return fmt.Sprintf("\"%s\"", pinotTableName)
}

func (s *dataService) ProcessParamsForQuery(ctx context.Context, merchantId string, params map[string]string, providerType constants.ProviderType) (servicemodels.QueryMetadata, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	queryMetadata := servicemodels.QueryMetadata{
		Params:     map[string]string{},
		TableNames: []string{},
	}
	for key, arg := range params {
		if strings.HasPrefix(strings.ToLower(key), serviceconstants.ZampTableName) {

			datasetInfo, err := s.getDataset(ctx, merchantId, arg)
			if err != nil {
				logger.Error(errors.GettingDatasetInfoFailedErrMessage, zap.Error(err))
				return servicemodels.QueryMetadata{}, errors.ErrGettingDatasetInfoFailed
			}

			switch providerType {
			case constants.ProviderTypeDatabricks:
				if datasetInfo.DatabricksFQTableName == "" {
					return servicemodels.QueryMetadata{}, errors.ErrDatasetNotFoundInProvider
				}
				databricksFQTableName := parseDatabricksFQTableName(datasetInfo.DatabricksFQTableName)
				queryMetadata.TableNames = append(queryMetadata.TableNames, databricksFQTableName)
				queryMetadata.Params[key] = databricksFQTableName
			case constants.ProviderTypePinot:
				if datasetInfo.PinotTableName == "" {
					return servicemodels.QueryMetadata{}, errors.ErrDatasetNotFoundInProvider
				}
				pinotTableName := parsePinotTableName(datasetInfo.PinotTableName)
				queryMetadata.TableNames = append(queryMetadata.TableNames, pinotTableName)
				queryMetadata.Params[key] = pinotTableName
			}
		} else {
			queryMetadata.Params[key] = arg
		}
	}
	return queryMetadata, nil
}

func (s *dataService) GetDatasetConfig(ctx context.Context, merchantId string, datasetId string) (servicemodels.DatasetConfig, error) {
	dataset, err := s.getDataset(ctx, merchantId, datasetId)
	if err != nil {
		return servicemodels.DatasetConfig{}, err
	}

	datasetConfig := servicemodels.DatasetConfig{}
	if dataset.DatasetConfig != "" {
		err = json.Unmarshal([]byte(dataset.DatasetConfig), &datasetConfig)
		if err != nil {
			return servicemodels.DatasetConfig{}, errors.ErrJSONUnmarshallingFailed
		}
	}
	return datasetConfig, nil
}

func (s *dataService) getDataset(ctx context.Context, merchantId string, datasetId string) (servicemodels.Dataset, error) {
	// TODO: ADD REDIS LAYER HERE
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetsTableName := helpers.BuildDatabricksTableName(s.dataPlatformConfig.DatabricksConfig.ZampDatabricksCatalog, s.dataPlatformConfig.DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.DatasetTableName)
	query, err := helpers.FillQueryTemplate(ctx, serviceconstants.QueryGetDatasetById, map[string]string{
		serviceconstants.DatasetTableNameQueryParam:  datasetsTableName,
		serviceconstants.DatasetMerchantIdColumnName: merchantId,
		serviceconstants.DatasetIdColumnName:         datasetId,
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return servicemodels.Dataset{}, errors.ErrTemplateParsingFailed
	}

	providerService, err := s.getProviderService(ctx, merchantId, constants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return servicemodels.Dataset{}, errors.ErrProviderServiceNotFound
	}

	queryResult, err := providerService.Query(ctx, datasetsTableName, query)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return servicemodels.Dataset{}, errors.ErrQueryingDatabricksFailed
	}

	if len(queryResult.Rows) == 0 {
		return servicemodels.Dataset{}, errors.ErrDatasetNotFound
	}

	dataset := servicemodels.Dataset{}
	jsonBytes, err := json.Marshal(queryResult.Rows[0])
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return servicemodels.Dataset{}, errors.ErrJSONUnmarshallingFailed
	}

	err = json.Unmarshal(jsonBytes, &dataset)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return servicemodels.Dataset{}, errors.ErrJSONUnmarshallingFailed
	}
	return dataset, nil
}

func (s *dataService) getDatabricksDatasetMetadata(ctx context.Context, datasetInfo servicemodels.Dataset) (servicemodels.InternalDatasetMetadata, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	databricksStats := servicemodels.DatasetStats{}
	databricksSchema := servicemodels.DatasetSchemaDetails{}

	if datasetInfo.DatabricksStats != "" {
		err := json.Unmarshal([]byte(datasetInfo.DatabricksStats), &databricksStats)
		if err != nil {
			logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
			return servicemodels.InternalDatasetMetadata{}, errors.ErrJSONUnmarshallingFailed
		}
	}
	if datasetInfo.DatabricksSchema != "" {
		err := json.Unmarshal([]byte(datasetInfo.DatabricksSchema), &databricksSchema)
		if err != nil {
			logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
			return servicemodels.InternalDatasetMetadata{}, errors.ErrJSONUnmarshallingFailed
		}
	}

	return servicemodels.InternalDatasetMetadata{
		Schema: databricksSchema,
		Stats:  databricksStats,
	}, nil
}

func (s *dataService) getPinotDatasetMetadata(ctx context.Context, datasetInfo servicemodels.Dataset) (servicemodels.InternalDatasetMetadata, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	pinotStats := servicemodels.DatasetStats{}
	pinotSchema := servicemodels.DatasetSchemaDetails{}

	if datasetInfo.PinotStats != "" {
		err := json.Unmarshal([]byte(datasetInfo.PinotStats), &pinotStats)
		if err != nil {
			logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
			return servicemodels.InternalDatasetMetadata{}, errors.ErrJSONUnmarshallingFailed
		}
	}

	if datasetInfo.PinotSchema != "" {
		err := json.Unmarshal([]byte(datasetInfo.PinotSchema), &pinotSchema)
		if err != nil {
			logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
			return servicemodels.InternalDatasetMetadata{}, errors.ErrJSONUnmarshallingFailed
		}
	}

	return servicemodels.InternalDatasetMetadata{
		Schema: pinotSchema,
		Stats:  pinotStats,
	}, nil
}

func (s *dataService) getProviderLevelDatasetMetadata(ctx context.Context, merchantId string, datasetId string) (servicemodels.ProviderLevelDatasetMetadata, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	datasetInfo, err := s.getDataset(ctx, merchantId, datasetId)
	if err != nil {
		logger.Error(errors.GettingDatasetInfoFailedErrMessage, zap.Error(err))
		return servicemodels.ProviderLevelDatasetMetadata{}, errors.ErrGettingDatasetInfoFailed
	}

	databricksMetadata, err := s.getDatabricksDatasetMetadata(ctx, datasetInfo)
	if err != nil {
		logger.Error(errors.GettingDatabricksDatasetMetadataFailedErrMessage, zap.Error(err))
		return servicemodels.ProviderLevelDatasetMetadata{}, errors.ErrGettingDatabricksDatasetMetadataFailed
	}

	pinotMetadata, err := s.getPinotDatasetMetadata(ctx, datasetInfo)
	if err != nil {
		logger.Error(errors.GettingPinotDatasetMetadataFailedErrMessage, zap.Error(err))
		return servicemodels.ProviderLevelDatasetMetadata{}, errors.ErrGettingPinotDatasetMetadataFailed
	}

	return servicemodels.ProviderLevelDatasetMetadata{
		Databricks: databricksMetadata,
		Pinot:      pinotMetadata,
	}, nil
}

func (s *dataService) GetDatasetMetadata(ctx context.Context, merchantId string, datasetId string) (servicemodels.DatasetMetadata, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	providerLevelDatasetMetadata, err := s.getProviderLevelDatasetMetadata(ctx, merchantId, datasetId)
	if err != nil {
		logger.Error(errors.GettingProviderLevelDatasetMetadataFailedErrMessage, zap.Error(err))
		return servicemodels.DatasetMetadata{}, errors.ErrGettingProviderLevelDatasetMetadataFailed
	}
	schema := providerLevelDatasetMetadata.Databricks.Schema.Columns
	tableMetadata := servicemodels.DatasetMetadata{
		Schema: schema,
		Stats:  providerLevelDatasetMetadata.Databricks.Stats,
	}
	return tableMetadata, nil
}

func (s *dataService) GetDatasetParents(ctx context.Context, merchantId string, datasetId string) (servicemodels.DatasetParents, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	jobMappingsTableName := helpers.BuildDatabricksTableName(s.dataPlatformConfig.DatabricksConfig.ZampDatabricksCatalog, s.dataPlatformConfig.DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.JobMappingsTableName)
	query, err := helpers.FillQueryTemplate(ctx, serviceconstants.QueryGetDatasetParents, map[string]string{
		serviceconstants.JobMappingsTableNameQueryParam:       jobMappingsTableName,
		serviceconstants.JobMappingDestinationTypeColumnName:  string(serviceconstants.DAGDatasetDestination),
		serviceconstants.JobMappingDestinationValueColumnName: datasetId,
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return servicemodels.DatasetParents{}, errors.ErrTemplateParsingFailed
	}

	providerService, err := s.getProviderService(ctx, merchantId, constants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.UnsupportedProviderTypeErrMessage, zap.Error(err))
		return servicemodels.DatasetParents{}, errors.ErrUnsupportedProviderType
	}

	queryResult, err := providerService.Query(ctx, jobMappingsTableName, query)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return servicemodels.DatasetParents{}, errors.ErrQueryingDatabricksFailed
	}

	if len(queryResult.Rows) == 0 {
		logger.Error(errors.DatasetNotFoundErrMessage, zap.String("datasetId", datasetId))
		return servicemodels.DatasetParents{}, errors.ErrDatasetNotFound
	}

	jobDatasetMappings := []servicemodels.JobDatasetMapping{}
	jsonBytes, err := json.Marshal(queryResult.Rows)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return servicemodels.DatasetParents{}, errors.ErrJSONUnmarshallingFailed
	}

	err = json.Unmarshal(jsonBytes, &jobDatasetMappings)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return servicemodels.DatasetParents{}, errors.ErrJSONUnmarshallingFailed
	}

	dagNodes := []servicemodels.DAGNode{}
	for _, jobDatasetMapping := range jobDatasetMappings {
		dagNode := servicemodels.DAGNode{
			Id:   jobDatasetMapping.SourceValue,
			Type: jobDatasetMapping.SourceType,
		}
		dagNodes = append(dagNodes, dagNode)
	}

	return servicemodels.DatasetParents{Parents: dagNodes}, nil
}

func (s *dataService) GetDatabricksWarehouseId(ctx context.Context, dataProviderId string) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	databricksConfig, ok := s.dataPlatformConfig.DatabricksConfig.DataProviderConfigs[dataProviderId]
	if !ok {
		logger.Error(errors.ProviderServiceNotFoundErrMessage)
		return "", errors.ErrProviderServiceNotFound
	}
	return databricksConfig.WarehouseId, nil
}

func (s *dataService) GetDatabricksServiceForMerchant(ctx context.Context, merchantId string) (databricks.DatabricksService, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	dataProviderId, err := s.GetDataProviderIdForMerchant(merchantId, constants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.GettingDataProviderIdForMerchantFailedErrMessage, zap.Error(err))
		return nil, err
	}
	databricksService, err := s.providerService.GetDatabricksService(ctx, dataProviderId)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return nil, errors.ErrProviderServiceNotFound
	}
	return databricksService, nil
}

func (s *dataService) GetDatabricksServiceForProvider(ctx context.Context, providerId string) (databricks.DatabricksService, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	databricksService, err := s.providerService.GetDatabricksService(ctx, providerId)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return nil, err
	}
	return databricksService, nil
}

func (s *dataService) GetDatasetEdgesByMerchant(ctx context.Context, merchantId string) ([]servicemodels.JobDatasetMapping, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	providerService, err := s.getProviderService(ctx, merchantId, constants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.UnsupportedProviderTypeErrMessage, zap.Error(err))
		return nil, errors.ErrUnsupportedProviderType
	}

	jobMappingsTableName := helpers.BuildDatabricksTableName(s.dataPlatformConfig.DatabricksConfig.ZampDatabricksCatalog, s.dataPlatformConfig.DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.JobMappingsTableName)
	query, err := helpers.FillQueryTemplate(ctx, serviceconstants.QueryGetDatasetEdgesByMerchant, map[string]string{
		serviceconstants.JobMappingMerchantIdColumnName:        merchantId,
		serviceconstants.JobMappingsTableNameQueryParam:        jobMappingsTableName,
		serviceconstants.JobMappingSourceTypeFolderColumnName:  string(serviceconstants.JobMappingTypeFolder),
		serviceconstants.JobMappingSourceTypeDatasetColumnName: string(serviceconstants.JobMappingTypeDataset),
		serviceconstants.JobMappingDestinationTypeColumnName:   string(serviceconstants.JobMappingTypeDataset),
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrTemplateParsingFailed
	}

	queryResult, err := providerService.Query(ctx, jobMappingsTableName, query)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return nil, errors.ErrQueryingDatabricksFailed
	}

	jobDatasetMappings := []servicemodels.JobDatasetMapping{}
	jsonBytes, err := json.Marshal(queryResult.Rows)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	err = json.Unmarshal(jsonBytes, &jobDatasetMappings)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	return jobDatasetMappings, nil
}
