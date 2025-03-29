package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	dataplatformactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	dataplatformactionmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	dataplatformconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	dataplatformConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	dataplatformcoremodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	"github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	datasetConstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	"github.com/Zampfi/application-platform/services/api/core/datasets/errors"
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	rulemodels "github.com/Zampfi/application-platform/services/api/core/rules/models"
	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	dataplatformpkgmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	querybuilderconstants "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
	querybuildermodels "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *datasetService) convertToFilterConfig(datasetInfo dataplatformDataModels.DatasetMetadata, datasetMetaInfo models.DatasetMetadataConfig) []models.FilterConfig {
	var filterConfigs []models.FilterConfig
	mapDisplayConfig := make(map[string]models.DisplayConfig)

	for _, displayConfig := range datasetMetaInfo.DisplayConfig {
		mapDisplayConfig[displayConfig.Column] = displayConfig
	}

	for columnName, schemaInfo := range datasetInfo.Schema {
		if s.isHiddenColumn(columnName) {
			continue
		}

		columType := dataplatformdataconstants.Datatype(schemaInfo.Type)
		metadata := make(map[string]interface{})
		var alias *string

		if slices.Contains(datasetConstants.WhiteListedZampColumns, columnName) || strings.HasPrefix(columnName, datasetConstants.ZampUpdateColumnSourcePrefix) {
			metadata[datasetConstants.MetadataConfigIsHidden] = true
			metadata[datasetConstants.MetadataConfigIsEditable] = false
		}

		filterType := s.determineFilterType(columType, datasetInfo, columnName)

		if _, ok := datasetMetaInfo.Columns[columnName]; ok {
			metadata[datasetConstants.MetadataConfigCustomType] = datasetMetaInfo.Columns[columnName].CustomType
			metadata[datasetConstants.MetadataConfig] = datasetMetaInfo.Columns[columnName].CustomTypeConfig

			if datasetMetaInfo.Columns[columnName].CustomType == dataplatformconstants.DatabricksColumnCustomTypeTags {
				metadata[datasetConstants.MetadataConfigIsEditable] = true
			}
		}

		if _, ok := mapDisplayConfig[columnName]; ok {
			alias = mapDisplayConfig[columnName].Alias
			metadata = s.buildColumnMetadata(mapDisplayConfig[columnName])
		}

		filterConfigs = append(filterConfigs, models.FilterConfig{
			Column:   columnName,
			Alias:    alias,
			Type:     filterType,
			DataType: &columType,
			Options:  []interface{}{},
			Metadata: metadata,
		})
	}
	return filterConfigs
}

func (s *datasetService) buildColumnMetadata(columnInfo models.DisplayConfig) map[string]interface{} {
	metadata := make(map[string]interface{})

	if columnInfo.IsHidden {
		metadata[datasetConstants.MetadataConfigIsHidden] = true
	}

	if columnInfo.IsEditable {
		metadata[datasetConstants.MetadataConfigIsEditable] = true
	}

	if columnInfo.Type != nil {
		metadata[datasetConstants.MetadataConfigCustomType] = columnInfo.Type
	}

	if columnInfo.Config != nil {
		metadata[datasetConstants.MetadataConfig] = columnInfo.Config
	}

	return metadata
}

func (s *datasetService) determineFilterType(columType dataplatformdataconstants.Datatype, datasetInfo dataplatformDataModels.DatasetMetadata, columnName string) string {
	switch columType {
	case dataplatformdataconstants.DateDataType, dataplatformdataconstants.TimestampDataType, dataplatformdataconstants.TimestampNtzDataType:
		return datasetConstants.FilterTypeDateRange
	case dataplatformdataconstants.DecimalDataType, dataplatformdataconstants.DoubleDataType, dataplatformdataconstants.FloatDataType,
		dataplatformdataconstants.IntegerDataType, dataplatformdataconstants.SmallIntDataType, dataplatformdataconstants.TinyIntDataType,
		dataplatformdataconstants.BigIntDataType:
		return datasetConstants.FilterTypeAmountRange
	case dataplatformdataconstants.StringDataType:
		distinctCount := s.getDistinctValueCount(datasetInfo, columnName)
		return s.getStringColumnFilterType(distinctCount)
	case dataplatformdataconstants.ArrayOfStringDataType:
		return datasetConstants.FilterTypeArraySearch
	case dataplatformdataconstants.BooleanDataType:
		return datasetConstants.FilterTypeMultiSearch
	default:
		return datasetConstants.FilterTypeSearch
	}
}

func (s *datasetService) getStringColumnFilterType(distinctCount int) string {
	if distinctCount <= datasetConstants.MultiSelectThreshold {
		return datasetConstants.FilterTypeMultiSearch
	}
	return datasetConstants.FilterTypeSearch
}

func (s *datasetService) getDistinctValueCount(datasetInfo dataplatformDataModels.DatasetMetadata, colName string) int {
	if columnStats, exists := datasetInfo.Stats.ColumnStats[colName]; exists {
		return columnStats.DistinctCount
	}
	return datasetConstants.MultiSelectThreshold + 1
}

func (s *datasetService) getQueryBuilderFilterModel(conditions []models.Filter, columnDatatypes map[string]dataplatformdataconstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.Filter {
	filters := make([]querybuildermodels.Filter, len(conditions))
	for i, condition := range conditions {
		filters[i] = querybuildermodels.Filter{
			LogicalOperator: (*querybuildermodels.LogicalOperator)(condition.LogicalOperator),
			Column: querybuildermodels.ColumnConfig{
				Column: condition.Column,
				Datatype: func() *dataplatformdataconstants.Datatype {
					dt := columnDatatypes[condition.Column]
					return &dt
				}(),
				CustomDataConfig: func() *querybuildermodels.CustomDataTypeConfig {
					if customColumnConfig[condition.Column].Type != "" {
						return &querybuildermodels.CustomDataTypeConfig{
							Type:   customColumnConfig[condition.Column].Type,
							Config: customColumnConfig[condition.Column].Config,
						}
					}
					return nil
				}(),
			},
			Operator: condition.Operator,
			Value:    condition.Value,
		}

		if condition.Conditions != nil {
			filters[i].Conditions = s.getQueryBuilderFilterModel(condition.Conditions, columnDatatypes, customColumnConfig)
		}
	}
	return filters
}

func (s *datasetService) addDefaultZampIsDeletedFilterModel(filterParams models.FilterModel, columnDatatypes map[string]dataplatformdataconstants.Datatype,
	customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) querybuildermodels.FilterModel {
	var booleanDatatype dataplatformdataconstants.Datatype = "boolean"
	logicalOperator := querybuildermodels.LogicalOperator(filterParams.LogicalOperator)

	return querybuildermodels.FilterModel{
		LogicalOperator: querybuilderconstants.LogicalOperatorAnd,
		Conditions: []querybuildermodels.Filter{
			{
				LogicalOperator: &logicalOperator,
				Column: querybuildermodels.ColumnConfig{
					Column:   datasetConstants.ZampIsDeletedColumn,
					Datatype: &booleanDatatype,
				},
				Operator:   querybuilderconstants.EqualOperator,
				Value:      false,
				Conditions: s.getQueryBuilderFilterModel(filterParams.Conditions, columnDatatypes, customColumnConfig),
			},
		},
	}
}

func (s *datasetService) mapToQueryConfig(datasetId string, queryConfig models.DatasetParams, datasetInfo dataplatformDataModels.DatasetMetadata,
	columnDatatypes map[string]dataplatformdataconstants.Datatype, datasetMetaData models.DatasetMetadataConfig) querybuildermodels.QueryConfig {

	customColumnConfig := s.getCustomColumnConfig(datasetMetaData.DatasetConfig, queryConfig)

	filters := s.addDefaultZampIsDeletedFilterModel(queryConfig.Filters, columnDatatypes, customColumnConfig)

	filteredColumns := s.buildFilteredColumns(queryConfig, datasetInfo, columnDatatypes, customColumnConfig)

	aggregations := s.getQueryBuilderAggregationModel(queryConfig.Aggregations, columnDatatypes, customColumnConfig)
	groupBy := s.getQueryBuilderGroupByModel(queryConfig.GroupBy, columnDatatypes, customColumnConfig)
	orderBy := s.getQueryBuilderOrderByModel(queryConfig.OrderBy, columnDatatypes, customColumnConfig)

	var pagination *querybuildermodels.Pagination
	if queryConfig.Pagination != nil {
		pagination = &querybuildermodels.Pagination{
			Page:     queryConfig.Pagination.Page,
			PageSize: queryConfig.Pagination.PageSize,
		}
	}

	var subquery *querybuildermodels.QueryConfig
	if queryConfig.Subquery != nil {
		subquery = func() *querybuildermodels.QueryConfig {
			subqueryConfig := s.mapToQueryConfig(datasetId, *queryConfig.Subquery, datasetInfo, columnDatatypes, datasetMetaData)
			return &subqueryConfig
		}()
	}

	var windows []querybuildermodels.WindowConfig
	if queryConfig.Windows != nil {
		windows = func() []querybuildermodels.WindowConfig {
			windowsConfig := s.getQueryBuilderWindowModel(queryConfig.Windows, columnDatatypes, customColumnConfig)
			return windowsConfig
		}()
	}

	return querybuildermodels.QueryConfig{
		TableConfig: querybuildermodels.TableConfig{
			DatasetId: datasetId,
			Columns:   filteredColumns,
		},
		Subquery:     subquery,
		Windows:      windows,
		Filters:      filters,
		Aggregations: aggregations,
		GroupBy:      groupBy,
		OrderBy:      orderBy,
		CountAll:     queryConfig.CountAll,
		Pagination:   pagination,
	}
}

func (s *datasetService) getQueryBuilderWindowModel(windows []models.WindowConfig, columnDatatypes map[string]dataplatformConstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.WindowConfig {
	result := make([]querybuildermodels.WindowConfig, len(windows))
	for i, window := range windows {
		result[i] = querybuildermodels.WindowConfig{
			Function:    querybuildermodels.WindowFunction(window.Function),
			PartitionBy: s.getQueryBuilderColumnConfigModel(window.PartitionBy, columnDatatypes, customColumnConfig),
			OrderBy:     s.getQueryBuilderOrderByModel(window.OrderBy, columnDatatypes, customColumnConfig),
			Alias:       window.Alias,
		}
	}
	return result
}

func (s *datasetService) getQueryBuilderAggregationModel(aggregations []models.Aggregation, columnDatatypes map[string]dataplatformdataconstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.Aggregation {
	result := make([]querybuildermodels.Aggregation, len(aggregations))
	for i, agg := range aggregations {
		result[i] = querybuildermodels.Aggregation{
			Column: querybuildermodels.ColumnConfig{
				Column: agg.Column,
				Datatype: func() *dataplatformdataconstants.Datatype {
					dt := columnDatatypes[agg.Column]
					return &dt
				}(),
				CustomDataConfig: func() *querybuildermodels.CustomDataTypeConfig {
					if customColumnConfig[agg.Column].Type != "" {
						return &querybuildermodels.CustomDataTypeConfig{
							Type:   customColumnConfig[agg.Column].Type,
							Config: customColumnConfig[agg.Column].Config,
						}
					}
					return nil
				}(),
				Alias: &agg.Alias,
			},
			Alias:    agg.Alias,
			Function: querybuildermodels.AggregationFunction(agg.Function),
		}
	}
	return result
}

func (s *datasetService) getQueryBuilderGroupByModel(groupBy []models.GroupBy, columnDatatypes map[string]dataplatformdataconstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.GroupBy {
	result := make([]querybuildermodels.GroupBy, len(groupBy))
	for i, gb := range groupBy {
		result[i] = querybuildermodels.GroupBy{
			Column: querybuildermodels.ColumnConfig{
				Column: gb.Column,
				Datatype: func() *dataplatformdataconstants.Datatype {
					dt := columnDatatypes[gb.Column]
					return &dt
				}(),
				CustomDataConfig: func() *querybuildermodels.CustomDataTypeConfig {
					if customColumnConfig[gb.Column].Type != "" {
						return &querybuildermodels.CustomDataTypeConfig{
							Type:   customColumnConfig[gb.Column].Type,
							Config: customColumnConfig[gb.Column].Config,
						}
					}
					return nil
				}(),
				Alias: gb.Alias,
			},
		}
	}
	return result
}

func (s *datasetService) getCustomColumnConfig(datasetConfig dataplatformDataModels.DatasetConfig, params models.DatasetParams) map[string]querybuildermodels.CustomDataTypeConfig {
	customColumnConfig := make(map[string]querybuildermodels.CustomDataTypeConfig)

	FXColumnConfigMap := make(map[string]querybuildermodels.AmountCustomTypeConfig)
	for _, customColumnGroup := range datasetConfig.CustomColumnGroups {
		if customColumnGroup.Type == "FX" && params.FxCurrency != nil {
			FXColumnConfigMap[customColumnGroup.Config[datasetConstants.AmountColumn]] = querybuildermodels.AmountCustomTypeConfig{
				CurrencyColumn: customColumnGroup.Config[datasetConstants.CurrencyColumn],
				FxCurrency:     *params.FxCurrency,
				AmountColumn:   customColumnGroup.Config[datasetConstants.AmountColumn],
			}
		}
	}

	for columnName, columnConfig := range datasetConfig.Columns {
		var customTypeConfig querybuildermodels.CustomDataTypeInterface
		if columnConfig.CustomType == dataplatformconstants.DatabricksColumnCustomTypeAmount {
			if _, ok := FXColumnConfigMap[columnName]; ok {
				customTypeConfig = &querybuildermodels.AmountCustomTypeConfig{
					CurrencyColumn: FXColumnConfigMap[columnName].CurrencyColumn,
					FxCurrency:     FXColumnConfigMap[columnName].FxCurrency,
					AmountColumn:   FXColumnConfigMap[columnName].AmountColumn,
				}

				customColumnConfig[columnName] = querybuildermodels.CustomDataTypeConfig{
					Type:   dataplatformconstants.DatabricksColumnCustomType(columnConfig.CustomType),
					Config: customTypeConfig,
				}
			}
		}
	}
	return customColumnConfig
}

func (s *datasetService) getQueryBuilderOrderByModel(orderBy []models.OrderBy, columnDatatypes map[string]dataplatformdataconstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.OrderBy {
	result := make([]querybuildermodels.OrderBy, len(orderBy))
	for i, ob := range orderBy {
		result[i] = querybuildermodels.OrderBy{
			Column: querybuildermodels.ColumnConfig{
				Column: ob.Column,
				Datatype: func() *dataplatformdataconstants.Datatype {
					dt := columnDatatypes[ob.Column]
					return &dt
				}(),
				CustomDataConfig: func() *querybuildermodels.CustomDataTypeConfig {
					if customColumnConfig[ob.Column].Type != "" {
						return &querybuildermodels.CustomDataTypeConfig{
							Type:   customColumnConfig[ob.Column].Type,
							Config: customColumnConfig[ob.Column].Config,
						}
					}
					return nil
				}(),
				Alias: ob.Alias,
			},
			Order: querybuildermodels.OrderType(ob.Order),
		}
	}
	return result
}

func (s *datasetService) createCountQueryConfig(queryConfigMapped querybuildermodels.QueryConfig) querybuildermodels.QueryConfig {
	return querybuildermodels.QueryConfig{
		Filters:      queryConfigMapped.Filters,
		TableConfig:  queryConfigMapped.TableConfig,
		GroupBy:      queryConfigMapped.GroupBy,
		Aggregations: queryConfigMapped.Aggregations,
	}
}

func (s *datasetService) getTotalCount(ctx context.Context, merchantId uuid.UUID, datasetId string, queryConfigMapped querybuildermodels.QueryConfig) (int64, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	countQueryConfig := s.createCountQueryConfig(queryConfigMapped)

	query, _, err := s.queryBuilderService.ToSQL(ctx, countQueryConfig)
	if err != nil {
		logger.Error("failed to build query", zap.String("error", err.Error()))
		return 0, errors.ErrFailedToBuildQuery
	}

	countQuery := fmt.Sprintf(datasetConstants.GetRowCountQuery, query)

	var result dataplatformpkgmodels.QueryResult

	switch s.serverDatasetConfig.DataplatformProvider {
	case datasetConstants.DataplatformProviderDatabricks:
		result, err = s.dataplatformService.Query(ctx, merchantId.String(), countQuery, map[string]string{
			datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
		})
	case datasetConstants.DataplatformProviderPinot:
		result, err = s.dataplatformService.QueryRealTime(ctx, merchantId.String(), countQuery, map[string]string{
			datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
		})
	default:
		return 0, errors.ErrInvalidDataplatformProvider
	}

	if err != nil {
		logger.Error("failed to query dataset", zap.String("dataset_id", datasetId), zap.Error(err))
		return 0, errors.ErrFailedToGetData
	}

	totalCount, ok := result.Rows[0][result.Columns[0].Name].(int64)
	if !ok {
		logger.Error("failed to get total count", zap.String("dataset_id", datasetId), zap.Error(errors.ErrInvalidMetadataFormat))
		return 0, errors.ErrInvalidMetadataFormat
	}

	return totalCount, nil
}

func (s *datasetService) fetchRowDetails(ctx context.Context, merchantId string, rowDetails dataplatformpkgmodels.QueryResult) ([]models.DatasetInfo, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	if len(rowDetails.Rows) == 0 {
		logger.Info("no rows found for the specific rowUUID", zap.Any("row_details", rowDetails))
		return nil, errors.ErrNoRow
	}

	if len(rowDetails.Rows) > 1 {
		logger.Warn("more than one row found for the specific rowUUID", zap.Any("row_details", rowDetails), zap.Error(errors.ErrMoreThanOneRow))
	}

	drilldownMetadata, ok := rowDetails.Rows[0][datasetConstants.ZampDrilldownMetadataColumn].(string)
	if !ok {
		return nil, errors.ErrInvalidMetadataFormat
	}

	var metadata map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(drilldownMetadata), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	var parentDatasets []models.DatasetInfo
	for datasetId, rows := range metadata {
		var conditions []models.Filter
		for _, rowData := range rows {
			for columnName, value := range rowData {
				conditions = append(conditions, models.Filter{
					Column:   columnName,
					Operator: querybuilderconstants.EqualOperator,
					Value:    value,
				})
			}
		}

		parentDatasets = append(parentDatasets, models.DatasetInfo{
			DatasetId: datasetId,
			Filters: models.FilterModel{
				LogicalOperator: models.LogicalOperator(querybuilderconstants.LogicalOperatorAnd),
				Conditions:      conditions,
			},
		})
	}

	return parentDatasets, nil
}

func (s *datasetService) processRowDetails(ctx context.Context, logger *zap.Logger, merchantId, datasetId, rowUUID string) (models.ParentDatasetInfo, error) {
	var err error
	var rowDetails dataplatformpkgmodels.QueryResult

	switch s.serverDatasetConfig.DataplatformProvider {
	case datasetConstants.DataplatformProviderDatabricks:
		rowDetails, err = s.dataplatformService.Query(ctx, merchantId, fmt.Sprintf(datasetConstants.GetRowDetailsQuery, datasetId, rowUUID), map[string]string{
			datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
		})
	case datasetConstants.DataplatformProviderPinot:
		rowDetails, err = s.dataplatformService.QueryRealTime(ctx, merchantId, fmt.Sprintf(datasetConstants.GetRowDetailsQuery, datasetId, rowUUID), map[string]string{
			datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
		})
	default:
		return models.ParentDatasetInfo{}, errors.ErrInvalidDataplatformProvider
	}

	if err != nil {
		logger.Error("failed to get the row details", zap.Error(err))
		return models.ParentDatasetInfo{}, fmt.Errorf("fetch row details query: %w", err)
	}

	parentDatasets, err := s.fetchRowDetails(ctx, merchantId, rowDetails)
	if err != nil {
		logger.Error("failed to fetch row details", zap.Error(err))
		return models.ParentDatasetInfo{}, fmt.Errorf("process row details: %w", err)
	}

	return models.ParentDatasetInfo{ParentDatasets: parentDatasets}, nil
}

func (s *datasetService) processParentDataset(ctx context.Context, logger *zap.Logger, merchantId, datasetId string) (dataplatformDataModels.DatasetParents, error) {
	parents, err := s.dataplatformService.GetDatasetParents(ctx, merchantId, datasetId)
	if err != nil {
		logger.Error("failed to get the parent dataset", zap.Error(err))
		return dataplatformDataModels.DatasetParents{}, fmt.Errorf("fetch parent dataset: %w", err)
	}

	return parents, nil
}

func (s *datasetService) mergeDatasetsInfoWithParents(parentDatasets []models.DatasetInfo, parents dataplatformDataModels.DatasetParents) []models.DatasetInfo {
	existingDatasets := make(map[string]struct{}, len(parentDatasets))
	for _, datasetInfo := range parentDatasets {
		existingDatasets[datasetInfo.DatasetId] = struct{}{}
	}

	for _, parent := range parents.Parents {
		if _, exists := existingDatasets[parent.Id]; !exists {
			parentDatasets = append(parentDatasets, models.DatasetInfo{
				DatasetId: parent.Id,
			})
		}
	}

	return parentDatasets
}

func (s *datasetService) fetchParentDatasetAndRowDetails(ctx context.Context, logger *zap.Logger, merchantId, datasetId, rowUUID string) (models.ParentDatasetInfo, dataplatformDataModels.DatasetParents, error) {
	g, ctx := errgroup.WithContext(ctx)

	var rowDetailsRes models.ParentDatasetInfo
	var parentDatasetRes dataplatformDataModels.DatasetParents

	g.Go(func() error {
		res, err := s.processRowDetails(ctx, logger, merchantId, datasetId, rowUUID)
		if err != nil {
			logger.Error("failed to process row details", zap.Error(err))
			return err
		}
		rowDetailsRes = res
		return nil
	})

	g.Go(func() error {
		res, err := s.processParentDataset(ctx, logger, merchantId, datasetId)
		if err != nil {
			logger.Error("failed to process parent dataset", zap.Error(err))
			return err
		}
		parentDatasetRes = res
		return nil
	})

	if err := g.Wait(); err != nil {
		return models.ParentDatasetInfo{}, dataplatformDataModels.DatasetParents{}, err
	}

	return rowDetailsRes, parentDatasetRes, nil
}

func (s *datasetService) createDataset(ctx context.Context, ds store.DatasetStore, merchantId uuid.UUID, userId uuid.UUID, datasetId uuid.UUID, datasetCreationInfo models.DatasetCreationInfo) (uuid.UUID, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	metadata := models.DatasetMetadataConfig{
		DatasetConfig:    datasetCreationInfo.DatasetConfig,
		DisplayConfig:    datasetCreationInfo.DisplayConfig,
		DatabricksConfig: datasetCreationInfo.DatabricksConfig,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		logger.Error("failed to marshal dataset metadata", zap.String("error", err.Error()))
		return uuid.Nil, errors.ErrFailedToMarshalMetadata
	}

	_, err = ds.CreateDataset(ctx, storemodels.Dataset{
		ID:             datasetId,
		Title:          datasetCreationInfo.DatasetTitle,
		Description:    datasetCreationInfo.DatasetDescription,
		Type:           datasetCreationInfo.DatasetType,
		OrganizationId: merchantId,
		CreatedBy:      userId,
		Metadata:       metadataJSON,
	})
	if err != nil {
		logger.Error("failed to insert dataset into db", zap.String("error", err.Error()))
		return uuid.Nil, errors.ErrFailedToRegisterDataset
	}

	return datasetId, nil
}

func (s *datasetService) mapUpdateDatasetDataParamsToQueryConfig(datasetId uuid.UUID, queryConfig models.UpdateDatasetDataParams, columnDatatypes map[string]dataplatformdataconstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) (querybuildermodels.QueryConfig, error) {
	return querybuildermodels.QueryConfig{
		TableConfig: querybuildermodels.TableConfig{
			DatasetId: datasetId.String(),
		},
		Filters: s.addDefaultZampIsDeletedFilterModel(queryConfig.Filters, columnDatatypes, customColumnConfig),
	}, nil
}

func (s *datasetService) populateFilterOptions(ctx context.Context, merchantId uuid.UUID, datasetId string, filterConfigs []models.FilterConfig) error {
	errgrp := errgroup.Group{}
	resultCh := make(chan struct {
		Index   int
		Options []interface{}
	}, len(filterConfigs))

	for i, config := range filterConfigs {
		index, cfg := i, config
		errgrp.Go(func() error {
			options, err := s.GetOptionsForColumn(ctx, merchantId, datasetId, cfg.Column, cfg.Type, true)
			if err != nil {
				return fmt.Errorf("failed to get options for %s: %w", cfg.Column, err)
			}

			resultCh <- struct {
				Index   int
				Options []interface{}
			}{Index: index, Options: options}

			return nil
		})
	}

	err := errgrp.Wait()
	if err != nil {
		return err
	}

	close(resultCh)

	for res := range resultCh {
		filterConfigs[res.Index].Options = append(filterConfigs[res.Index].Options, res.Options...)
	}

	return nil
}

type CreateRuleParams struct {
	Id             uuid.UUID
	Title          string
	Description    string
	OrganizationId uuid.UUID
	DatasetId      uuid.UUID
	Column         string
	Value          string
	FilterConfig   interface{}
	CreatedBy      uuid.UUID
}

func (s *datasetService) buildCreateRuleParams(ruleId uuid.UUID, organizationId uuid.UUID, datasetId uuid.UUID, params models.UpdateDatasetDataParams, Sql string, args map[string]interface{}, columnDatatypes map[string]dataplatformdataconstants.Datatype, customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) (storemodels.CreateRuleParams, error) {
	queryConfig, err := s.mapUpdateDatasetDataParamsToQueryConfig(datasetId, params, columnDatatypes, customColumnConfig)
	if err != nil {
		return storemodels.CreateRuleParams{}, err
	}

	return storemodels.CreateRuleParams{
		Id:             ruleId,
		Title:          params.RuleTitle,
		Description:    params.RuleDescription,
		OrganizationId: organizationId,
		DatasetId:      datasetId,
		Column:         params.Update.Column,
		Value:          fmt.Sprintf("%v", params.Update.Value),
		FilterConfig: rulemodels.FilterConfig{
			QueryConfig: queryConfig,
			Sql:         Sql,
			Args:        args,
		},
		CreatedBy: params.UserId,
	}, nil
}

func (s *datasetService) getDatasetRulesForDataPlatfrom(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID, column string) ([]dataplatformDataModels.Rule, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetIdString := datasetId.String()

	rules, err := s.ruleService.GetRules(ctx, storemodels.FilterRuleParams{
		OrganizationId: merchantId,
		DatasetColumns: []storemodels.DatasetColumn{
			{
				DatasetId: datasetId,
				Columns:   []string{column},
			},
		},
	})
	if err != nil {
		logger.Error("failed to get rules for dataset", zap.String("dataset_id", datasetIdString), zap.String("error", err.Error()))
		return nil, err
	}

	var datasetRules []dataplatformDataModels.Rule
	if _, ok := rules[datasetIdString]; ok {
		if _, ok := rules[datasetIdString][column]; ok {
			for _, rule := range rules[datasetIdString][column] {
				datasetRules = append(datasetRules, dataplatformDataModels.Rule{
					Id:           rule.ID.String(),
					Priority:     rule.Priority,
					ValueToApply: rule.Value,
					SqlCondition: rule.FilterConfig.Sql,
					SqlArgs:      rule.FilterConfig.Args,
				})
			}
		}
	}

	return datasetRules, nil
}

func (s *datasetService) handleRuleBasedDatasetUpdate(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID, params models.UpdateDatasetDataParams) (dataplatformactionmodels.CreateActionResponse, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetRules, err := s.getDatasetRulesForDataPlatfrom(ctx, merchantId, datasetId, params.Update.Column)
	if err != nil {
		logger.Error("failed to get dataset rules for data platfrom", zap.String("dataset_id", datasetId.String()), zap.String("error", err.Error()))
		return dataplatformactionmodels.CreateActionResponse{}, err
	}

	dataplatformAction, err := s.dataplatformService.UpdateDataset(ctx, dataplatformcoremodels.UpdateDatasetPayload{
		MerchantID: merchantId.String(),
		ActorId:    params.UserId.String(),
		ActionMetadataPayload: dataplatformactionmodels.UpdateDatasetEvent{
			EventType: dataplatformactionconstants.UpdateDatasetEventTypeUpsertRules,
			EventData: dataplatformactionmodels.UpdateDatasetActionPayload{
				DatasetId: datasetId.String(),
				DatasetConfig: dataplatformDataModels.DatasetConfig{
					Rules: map[string][]dataplatformDataModels.Rule{
						params.Update.Column: datasetRules,
					},
				},
			},
			EventMetadata: dataplatformactionmodels.UpsertRuleEventMetadata{
				DeltaRuleId: params.SourceId.String(),
				Column:      params.Update.Column,
				Type:        dataplatformactionconstants.UpsertRuleOperationCreate,
			},
		},
	})
	if err != nil {
		logger.Error("failed to update dataset", zap.String("error", err.Error()))
		return dataplatformactionmodels.CreateActionResponse{}, errors.ErrFailedToUpdateDataset
	}

	return dataplatformAction, nil
}

func (s *datasetService) handleRulePriorityUpdate(ctx context.Context, merchantId uuid.UUID, userId uuid.UUID, datasetId uuid.UUID, params models.UpdateRulePriorityParams) (dataplatformactionmodels.CreateActionResponse, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetRules, err := s.getDatasetRulesForDataPlatfrom(ctx, merchantId, datasetId, params.Column)
	if err != nil {
		logger.Error("failed to get dataset rules for data platfrom", zap.String("dataset_id", datasetId.String()), zap.String("error", err.Error()))
		return dataplatformactionmodels.CreateActionResponse{}, err
	}

	dataplatformAction, err := s.dataplatformService.UpdateDataset(ctx, dataplatformcoremodels.UpdateDatasetPayload{
		MerchantID: merchantId.String(),
		ActorId:    userId.String(),
		ActionMetadataPayload: dataplatformactionmodels.UpdateDatasetEvent{
			EventType: dataplatformactionconstants.UpdateDatasetEventTypeUpsertRules,
			EventData: dataplatformactionmodels.UpdateDatasetActionPayload{
				DatasetId: datasetId.String(),
				DatasetConfig: dataplatformDataModels.DatasetConfig{
					Rules: map[string][]dataplatformDataModels.Rule{
						params.Column: datasetRules,
					},
				},
			},
			EventMetadata: dataplatformactionmodels.UpsertRuleEventMetadata{
				Column: params.Column,
				Type:   dataplatformactionconstants.UpsertRuleOperationReorder,
			},
		},
	})
	if err != nil {
		logger.Error("failed to update rule priority", zap.String("error", err.Error()))
		return dataplatformactionmodels.CreateActionResponse{}, errors.ErrFailedToUpdateRulePriority
	}

	return dataplatformAction, nil
}

func (s *datasetService) fetchDatasetDetails(ctx context.Context, logger *zap.Logger, parentDatasets []models.DatasetInfo) ([]models.DatasetInfo, error) {
	errgrp, ctx := errgroup.WithContext(ctx)
	resultCh := make(chan struct {
		Index   int
		Dataset models.DatasetInfo
	}, len(parentDatasets))

	for i, parentDataset := range parentDatasets {
		index, dataset := i, parentDataset
		errgrp.Go(func() error {
			datasetDetails, err := s.datasetStore.GetDatasetById(ctx, dataset.DatasetId)
			if err != nil {
				logger.Error("failed to get dataset info", zap.Error(err))
				return errors.ErrFailedToGetData
			}

			dataset.DatasetTitle = datasetDetails.Title
			dataset.DatasetDescription = *datasetDetails.Description
			dataset.DatasetType = storemodels.DatasetType(datasetDetails.Type)

			resultCh <- struct {
				Index   int
				Dataset models.DatasetInfo
			}{Index: index, Dataset: dataset}

			return nil
		})
	}

	if err := errgrp.Wait(); err != nil {
		return nil, err
	}
	close(resultCh)

	result := make([]models.DatasetInfo, len(parentDatasets))
	for res := range resultCh {
		result[res.Index] = res.Dataset
	}

	return result, nil
}

func (s *datasetService) getQueryBuilderColumnConfigModel(params []models.ColumnConfig, columnDatatypes map[string]dataplatformdataconstants.Datatype,
	customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.ColumnConfig {
	var filteredColumns []querybuildermodels.ColumnConfig
	for _, column := range params {
		filteredColumns = append(filteredColumns, querybuildermodels.ColumnConfig{
			Column: column.Column,
			Datatype: func() *dataplatformdataconstants.Datatype {
				dt := columnDatatypes[column.Column]
				return &dt
			}(),
			CustomDataConfig: func() *querybuildermodels.CustomDataTypeConfig {
				if customColumnConfig[column.Column].Type != "" {
					return &querybuildermodels.CustomDataTypeConfig{
						Type:   customColumnConfig[column.Column].Type,
						Config: customColumnConfig[column.Column].Config,
					}
				}
				return nil
			}(),
			Alias: column.Alias,
		})
	}
	return filteredColumns
}

func (s *datasetService) buildFilteredColumns(params models.DatasetParams, datasetInfo dataplatformDataModels.DatasetMetadata, columnDatatypes map[string]dataplatformdataconstants.Datatype,
	customColumnConfig map[string]querybuildermodels.CustomDataTypeConfig) []querybuildermodels.ColumnConfig {
	var filteredColumns []querybuildermodels.ColumnConfig

	if params.Columns != nil {
		filteredColumns = s.getQueryBuilderColumnConfigModel(params.Columns, columnDatatypes, customColumnConfig)
	} else {
		for columnName := range datasetInfo.Schema {
			if !s.isHiddenColumn(columnName) {
				filteredColumns = append(filteredColumns, querybuildermodels.ColumnConfig{
					Column: columnName,
					Datatype: func() *dataplatformdataconstants.Datatype {
						dt := columnDatatypes[columnName]
						return &dt
					}(),
					CustomDataConfig: func() *querybuildermodels.CustomDataTypeConfig {
						if customColumnConfig[columnName].Type != "" {
							return &querybuildermodels.CustomDataTypeConfig{
								Type:   customColumnConfig[columnName].Type,
								Config: customColumnConfig[columnName].Config,
							}
						}
						return nil
					}(),
				})
			}
		}
	}

	return filteredColumns
}

func (s *datasetService) isDrilldownEnabled(datasetInfo dataplatformDataModels.DatasetMetadata) bool {
	for columnName := range datasetInfo.Schema {
		if columnName == datasetConstants.ZampDrilldownMetadataColumn {
			return true
		}
	}
	return false
}

func (s *datasetService) createCSVFromQueryResult(data models.DatasetData) (*bytes.Buffer, error) {
	var csvData bytes.Buffer
	writer := csv.NewWriter(&csvData)

	headers := make([]string, len(data.QueryResult.Columns))
	for i, col := range data.QueryResult.Columns {
		headers[i] = col.Name
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, row := range data.QueryResult.Rows {
		rowData := make([]string, len(headers))
		for i, header := range headers {
			if val, ok := row[header]; ok && val != nil {
				rowData[i] = fmt.Sprintf("%v", val)
			} else {
				rowData[i] = ""
			}
		}
		if err := writer.Write(rowData); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return &csvData, nil
}

func (s *datasetService) isFxEnabled(datasetSchema map[string]dataplatformDataModels.ColumnMetadata) bool {
	for columnName := range datasetSchema {
		if strings.HasPrefix(columnName, constants.ZampFxColumnPrefix) {
			return true
		}
	}
	return false
}

func (s *datasetService) isFileImportEnabled(ctx context.Context, merchantId uuid.UUID, datasetId string) bool {
	logger := apicontext.GetLoggerFromCtx(ctx)

	dag, err := s.getDatsetDags(ctx, merchantId.String(), datasetId)
	if err != nil {
		logger.Error("failed to get dags", zap.Error(err), zap.String("datasetId", datasetId))
		return false
	}

	_, err = dag.GetImportFilePath()
	if err != nil {
		logger.Error("failed to get import file path", zap.Error(err), zap.String("datasetId", datasetId))
		return false
	}

	return true
}

func (s *datasetService) getColumnDatatypes(datasetInfo dataplatformDataModels.DatasetMetadata) (map[string]dataplatformdataconstants.Datatype, error) {
	result := make(map[string]dataplatformdataconstants.Datatype)
	for columnName, columnMetadata := range datasetInfo.Schema {
		if columnMetadata.Type != "" {
			result[columnName] = dataplatformdataconstants.Datatype(columnMetadata.Type)
		} else {
			return nil, errors.ErrColumnMetadataTypeIsEmpty
		}
	}
	return result, nil
}

func (s *datasetService) isHiddenColumn(columnName string) bool {
	// If column is whitelisted, it's not hidden
	if slices.Contains(datasetConstants.WhiteListedZampColumns, columnName) || strings.HasPrefix(columnName, datasetConstants.ZampUpdateColumnSourcePrefix) {
		return false
	}

	// Hidden if starts with underscore
	if strings.HasPrefix(columnName, datasetConstants.UnderscorePrefix) {
		return true
	}

	return false
}

func (s *datasetService) buildDatasetFileUploads(ctx context.Context, storeDatasetFileUploadsMap map[uuid.UUID]storemodels.DatasetFileUpload, storeFileUploadsMap map[uuid.UUID]storemodels.FileUpload) ([]models.DatasetFileUpload, error) {
	var datasetFileUploads []models.DatasetFileUpload
	logger := apicontext.GetLoggerFromCtx(ctx)

	for _, datasetFileUpload := range storeDatasetFileUploadsMap {
		fileUpload, ok := storeFileUploadsMap[datasetFileUpload.FileUploadID]
		if !ok {
			continue
		}

		metadata := storemodels.DatasetFileUploadMetadata{}
		err := json.Unmarshal(datasetFileUpload.Metadata, &metadata)
		if err != nil {
			logger.Error("failed to unmarshal metadata", zap.String("error", err.Error()))
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		datasetFileUploads = append(datasetFileUploads, models.DatasetFileUpload{
			ID:                   datasetFileUpload.ID,
			DatasetID:            datasetFileUpload.DatasetID,
			FileID:               datasetFileUpload.FileUploadID,
			FileName:             fileUpload.Name,
			UploadedByUserID:     fileUpload.UploadedByUserID,
			UploadedByUser:       fileUpload.UploadedByUser,
			FileUploadStatus:     fileUpload.Status,
			FileUploadCreatedAt:  fileUpload.CreatedAt,
			FileUploadDeletedAt:  fileUpload.DeletedAt,
			FileAllignmentStatus: datasetFileUpload.FileAllignmentStatus,
			Metadata:             metadata,
		})
	}

	return datasetFileUploads, nil
}
func (s *datasetService) validateAllDatasetColumnRulesPresentInPriorityUpdate(rules []rulemodels.Rule, params models.UpdateRulePriorityParams) error {
	if len(rules) == 0 {
		return errors.ErrNoRulesPresent
	}

	upadateRulesMap := make(map[uuid.UUID]int)
	for _, rulePriority := range params.RulePriorities.RulePriority {
		upadateRulesMap[rulePriority.RuleId] = rulePriority.Priority
	}

	for _, rule := range rules {
		if _, ok := upadateRulesMap[rule.ID]; !ok {
			return fmt.Errorf("rule not found: %s", rule.ID.String())
		}
	}

	return nil
}

func ensureCurrentUsersAdminAccess(ctx context.Context, policies []storemodels.ResourceAudiencePolicy) error {
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)

	if currentUserId == nil {
		return fmt.Errorf("no user ID found in the context")
	}

	for _, policy := range policies {

		for _, userPolicy := range policy.UserPolicies {
			// user has admin access directly
			if userPolicy.UserId == *currentUserId && userPolicy.Privilege == storemodels.PrivilegeDatasetAdmin {
				return nil
			}
		}
	}

	return fmt.Errorf("current user does not have access to change permissions on the dataset")
}

func ensureAudienceNotAlreadyAdded(audienceType storemodels.AudienceType, audienceId uuid.UUID, policies []storemodels.ResourceAudiencePolicy) error {
	for _, policy := range policies {
		if policy.ResourceAudienceType == audienceType && policy.ResourceAudienceID == audienceId {
			return fmt.Errorf("audience already exists on the dataset")
		}
	}

	return nil
}

func ensureUserIsNotChangingTheirOwnAdminPolicy(policyToBeUpdated storemodels.ResourceAudiencePolicy, allPolicies []storemodels.ResourceAudiencePolicy, currentUserId uuid.UUID) error {

	if policyToBeUpdated.ResourceAudienceType == storemodels.AudienceTypeUser && policyToBeUpdated.ResourceAudienceID == currentUserId && policyToBeUpdated.Privilege == storemodels.PrivilegeDatasetAdmin {
		return fmt.Errorf("you cannot change own permissions")
	}

	isCurrentUserSeparatelyAdded := false
	if storemodels.AudienceType(policyToBeUpdated.ResourceAudienceType) == storemodels.AudienceTypeOrganization || storemodels.AudienceType(policyToBeUpdated.ResourceAudienceType) == storemodels.AudienceTypeTeam {
		for _, policy := range allPolicies {
			for _, userPolicy := range policy.UserPolicies {
				if policyToBeUpdated.ID != userPolicy.ResourceAudiencePolicyId && userPolicy.UserId == currentUserId && userPolicy.Privilege == storemodels.PrivilegeDatasetAdmin {
					isCurrentUserSeparatelyAdded = true
					return nil
				}
			}
		}

		if !isCurrentUserSeparatelyAdded {
			return fmt.Errorf("you cannot change own permissions")
		}
	}

	return nil
}
