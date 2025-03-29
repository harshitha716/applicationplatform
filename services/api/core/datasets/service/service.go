package service

import (
	"context"
	"encoding/json"
	"time"

	"fmt"
	"slices"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	"github.com/Zampfi/application-platform/services/api/helper"
	"golang.org/x/sync/errgroup"

	dataplatformactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	dataplatformactionmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	dataplatformConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	datasetactionconstants "github.com/Zampfi/application-platform/services/api/core/datasets/actions/constants"
	datasetConstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	"github.com/Zampfi/application-platform/services/api/core/datasets/errors"
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	rulemodels "github.com/Zampfi/application-platform/services/api/core/rules/models"
	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/pkg/cache"
	dataplatformpkgmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	workersconstants "github.com/Zampfi/application-platform/services/api/workers/defaultworker/constants"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	fileimportsservice "github.com/Zampfi/application-platform/services/api/core/fileimports"
	ruleservice "github.com/Zampfi/application-platform/services/api/core/rules/service"
	datasetFileUploadsModels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	cloudservicemodels "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"
	cloudservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/service"
	querybuildermodels "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
	querybuilderservice "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	s3 "github.com/Zampfi/application-platform/services/api/pkg/s3"
	"github.com/google/uuid"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"

	dataplatformservice "github.com/Zampfi/application-platform/services/api/core/dataplatform"
	datasetactionservice "github.com/Zampfi/application-platform/services/api/core/datasets/actions/service"

	temporalsdk "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal"
	temporalmodels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"

	"go.uber.org/zap"
)

type DatasetService interface {
	GetFilterConfigByDatasetId(ctx context.Context, merchantId uuid.UUID, datasetId string) ([]models.FilterConfig, map[string]interface{}, error)
	GetDataByDatasetId(ctx context.Context, merchantId uuid.UUID, datasetId string, params models.DatasetParams) (models.DatasetData, error)
	ExecuteRawQuery(ctx context.Context, merchantId uuid.UUID, datasetId string, query string, queryParams map[string]interface{}) (models.DatasetData, error)
	GetRowDetailsByUUID(ctx context.Context, merchantId uuid.UUID, datasetId string, rowUUID string) (models.ParentDatasetInfo, error)
	GetDatasetListing(ctx context.Context, merchantId uuid.UUID, params models.DatsetListingParams) ([]models.Dataset, error)
	UpdateDatasetData(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID, params models.UpdateDatasetDataParams) (models.DatasetAction, error)
	GetDatasetCount(ctx context.Context, merchantId uuid.UUID, params models.DatsetListingParams) (int64, error)
	RegisterDataset(ctx context.Context, merchantId uuid.UUID, userId uuid.UUID, datasetCreationInfo models.DatasetCreationInfo) (string, uuid.UUID, error)
	CopyDataset(ctx context.Context, merchantId uuid.UUID, userId uuid.UUID, params models.CopyDatasetParams) (string, uuid.UUID, error)
	UpdateDataset(ctx context.Context, merchantId uuid.UUID, datasetId string, params models.UpdateDatasetParams) (string, error)
	RegisterDatasetJob(ctx context.Context, merchantId uuid.UUID, jobInfo dataplatformactionmodels.RegisterJobActionPayload) (string, error)
	UpsertTemplate(ctx context.Context, merchantId uuid.UUID, templateConfig dataplatformactionmodels.UpsertTemplateActionPayload) (string, error)
	GetOptionsForColumn(ctx context.Context, merchantId uuid.UUID, datasetId string, column string, filterType string, respectThreshold bool) ([]interface{}, error)
	GetDatasetAudiences(ctx context.Context, datasetId uuid.UUID) ([]storemodels.ResourceAudiencePolicy, error)
	GetDatasetActions(ctx context.Context, merchantId uuid.UUID, filters storemodels.DatasetActionFilters) ([]models.DatasetAction, error)
	UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error
	UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error
	AddAudienceToDataset(ctx context.Context, datasetId uuid.UUID, audienceType storemodels.AudienceType, audienceId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error)
	BulkAddAudienceToDataset(ctx context.Context, datasetId uuid.UUID, payload models.BulkAddDatasetAudiencePayload) ([]*storemodels.ResourceAudiencePolicy, models.BulkAddDatasetAudienceErrors)
	RemoveAudienceFromDataset(ctx context.Context, datasetId uuid.UUID, audienceId uuid.UUID) error
	UpdateDatasetAudiencePrivilege(ctx context.Context, datasetId uuid.UUID, audienceId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error)
	GetRulesByDatasetColumns(ctx context.Context, orgId uuid.UUID, datasetColumns []storemodels.DatasetColumn) (map[string]map[string][]rulemodels.Rule, error)
	CreateDatasetExportAction(ctx context.Context, merchantId uuid.UUID, datasetId string, queryConfig models.DatasetParams, userId uuid.UUID) (string, error)
	DatasetExportTemporalActivity(ctx context.Context, params models.DatasetExportParams, datasetId uuid.UUID, userId uuid.UUID, orgIds []uuid.UUID, workflowId string) (string, error)
	GetDownloadableDataExportUrl(ctx context.Context, workflowId string) (string, error)
	GetRulesByIds(ctx context.Context, ruleIds []string) ([]rulemodels.Rule, error)
	UpdateRulePriority(ctx context.Context, orgId uuid.UUID, userId uuid.UUID, params models.UpdateRulePriorityParams) (models.DatasetAction, error)
	InitiateFilePreparationForDatasetImport(ctx context.Context, datasetId uuid.UUID, fileId uuid.UUID) (datasetActionId *uuid.UUID, err error)
	CreateDatasetFileUpload(ctx context.Context, datasetId uuid.UUID, fileId uuid.UUID, metadata json.RawMessage) error
	GetDatasetFileUploads(ctx context.Context, datasetId uuid.UUID) ([]models.DatasetFileUpload, error)
	UpdateDatasetFileUploadStatus(ctx context.Context, datasetFileUploadId uuid.UUID, params models.UpdateDatasetFileUploadParams) error
	ImportDataFromFile(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID, fileUploadId uuid.UUID) (err error)
	GetFileUploadPreview(ctx context.Context, fileUploadId uuid.UUID) (datasetFileUploadsModels.DatasetPreview, error)
	GetDatasetImportPath(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID) (*models.FileImportConfig, error)
	DeleteDataset(ctx context.Context, merchantId uuid.UUID, datasetId string) (string, error)
	GetDatasetDisplayConfig(ctx context.Context, merchantId uuid.UUID, datasetId string) ([]models.DisplayConfig, error)
}

type DatasetServiceStore interface {
	store.DatasetStore
	store.DatasetActionStore
	store.DatasetFileUploadStore
	store.TransactionStore
	store.FlattenedResourceAudiencePoliciesStore
}

type datasetService struct {
	datasetStore         DatasetServiceStore
	queryBuilderService  querybuilderservice.QueryBuilder
	dataplatformService  dataplatformservice.DataPlatformService
	ruleService          ruleservice.RuleService
	datasetActionService datasetactionservice.DatasetActionService
	fileImportService    fileimportsservice.FileImportService
	temporalService      temporalsdk.TemporalService
	cloudService         cloudservice.CloudService
	s3Client             s3.S3Client
	serverDatasetConfig  serverconfig.DatasetConfig
	cacheClient          cache.CacheClient
}

func NewDatasetService(
	appStore DatasetServiceStore,
	queryBuilderService querybuilderservice.QueryBuilder,
	dataplatformService dataplatformservice.DataPlatformService,
	ruleService ruleservice.RuleService,
	fileimportsservice fileimportsservice.FileImportService,
	temporalService temporalsdk.TemporalService,
	cloudService cloudservice.CloudService,
	s3Client s3.S3Client,
	serverDatasetConfig serverconfig.DatasetConfig,
	cacheClient cache.CacheClient,
) DatasetService {
	return &datasetService{
		datasetStore:         appStore,
		queryBuilderService:  queryBuilderService,
		dataplatformService:  dataplatformService,
		ruleService:          ruleService,
		datasetActionService: datasetactionservice.NewDatasetActionService(appStore),
		fileImportService:    fileimportsservice,
		temporalService:      temporalService,
		cloudService:         cloudService,
		s3Client:             s3Client,
		serverDatasetConfig:  serverDatasetConfig,
		cacheClient:          cacheClient,
	}
}

func (s *datasetService) GetFilterConfigByDatasetId(ctx context.Context, merchantId uuid.UUID, datasetId string) ([]models.FilterConfig, map[string]interface{}, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	filterConfigCacheKey, err := s.cacheClient.FormatKey(datasetConstants.DatasetFilterConfigCacheKey, datasetId)
	if err != nil {
		logger.Error("failed to format cache key", zap.String("error", err.Error()))
		return nil, nil, fmt.Errorf("failed to format cache key")
	}

	cacheFilterConfig := &models.CacheFilterConfig{}
	if err := s.cacheClient.Get(ctx, filterConfigCacheKey, cacheFilterConfig); err != nil {
		logger.Warn("failed to fetch filter config from cache", zap.String("dataset_id", datasetId), zap.String("error", err.Error()))
	} else {
		return cacheFilterConfig.FilterConfig, cacheFilterConfig.DatsetConfig, nil
	}

	datasetMetaInfo, err := s.datasetStore.GetDatasetById(ctx, datasetId)
	if err != nil {
		logger.Error("failed to get dataset meta info", zap.String("error", err.Error()))
		return nil, nil, errors.ErrFailedToGetDatasetById
	}

	var datasetMetaData models.DatasetMetadataConfig
	if err := json.Unmarshal([]byte(datasetMetaInfo.Metadata), &datasetMetaData); err != nil {
		logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
		return nil, nil, errors.ErrFailedToUnmarshalMetadata
	}

	datasetInfo, err := s.dataplatformService.GetDatasetMetadata(ctx, merchantId.String(), datasetId)
	if err != nil {
		logger.Error("failed to get dataset metadata", zap.String("error", err.Error()))
		return nil, nil, errors.ErrFailedToGetDatasetMetadata
	}

	datasetConfig := make(map[string]interface{})
	datasetConfig[datasetConstants.DatasetConfigIsFxEnabled] = s.isFxEnabled(datasetInfo.Schema)
	datasetConfig[datasetConstants.DatasetConfigIsFileImportEnabled] = s.isFileImportEnabled(ctx, merchantId, datasetId)

	filterConfigs := s.convertToFilterConfig(datasetInfo, datasetMetaData)

	err = s.populateFilterOptions(ctx, merchantId, datasetId, filterConfigs)
	if err != nil {
		logger.Error("failed to populate filter options", zap.String("error", err.Error()))
		return nil, nil, fmt.Errorf("failed to populate filter options")
	}

	cacheFilterConfig.FilterConfig = filterConfigs
	cacheFilterConfig.DatsetConfig = datasetConfig

	if err := s.cacheClient.Set(ctx, filterConfigCacheKey, cacheFilterConfig, datasetConstants.DatasetFilterConfigCacheExpiry); err != nil {
		logger.Warn("failed to set filter config in cache", zap.String("dataset_id", datasetId), zap.String("error", err.Error()))
	}

	return filterConfigs, datasetConfig, nil
}

func (s *datasetService) GetDataByDatasetId(ctx context.Context, merchantId uuid.UUID, datasetId string, params models.DatasetParams) (models.DatasetData, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetMetaInfo, err := s.datasetStore.GetDatasetById(ctx, datasetId)
	if err != nil {
		logger.Error("failed to get dataset meta info", zap.String("error", err.Error()))
		return models.DatasetData{}, errors.ErrFailedToGetDatasetById
	}

	var datasetMetaData models.DatasetMetadataConfig
	if err := json.Unmarshal([]byte(datasetMetaInfo.Metadata), &datasetMetaData); err != nil {
		logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
		return models.DatasetData{}, errors.ErrFailedToUnmarshalMetadata
	}

	datasetInfo, err := s.dataplatformService.GetDatasetMetadata(ctx, merchantId.String(), datasetId)
	if err != nil {
		logger.Error("failed to get dataset metadata", zap.String("error", err.Error()))
		return models.DatasetData{}, errors.ErrFailedToGetDatasetMetadata
	}

	columnDatatypes, err := s.getColumnDatatypes(datasetInfo)
	if err != nil {
		logger.Error("failed to get column datatypes", zap.String("error", err.Error()))
		return models.DatasetData{}, err
	}

	queryConfigMapped := s.mapToQueryConfig(datasetId, params, datasetInfo, columnDatatypes, datasetMetaData)

	query, _, err := s.queryBuilderService.ToSQL(ctx, queryConfigMapped)
	if err != nil {
		logger.Error("failed to build query", zap.String("error", err.Error()))
		return models.DatasetData{}, errors.ErrFailedToBuildQuery
	}

	logger.Info("QUERY BEFORE ROSETTA", zap.String("QUERY", query), zap.Any("DATASETPARAMS", params))

	errgrp := errgroup.Group{}
	var result dataplatformpkgmodels.QueryResult
	var totalCount *int64

	errgrp.Go(func() error {
		if s.serverDatasetConfig.DataplatformProvider == datasetConstants.DataplatformProviderDatabricks || params.GetDatafromLake {
			result, err = s.dataplatformService.Query(ctx, merchantId.String(), query, map[string]string{
				datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
			})
		} else if s.serverDatasetConfig.DataplatformProvider == datasetConstants.DataplatformProviderPinot {
			result, err = s.dataplatformService.QueryRealTime(ctx, merchantId.String(), query, map[string]string{
				datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
			})
		} else {
			return errors.ErrInvalidDataplatformProvider
		}

		if err != nil {
			logger.Error("failed to get the data", zap.String("error", err.Error()))
			return errors.ErrFailedToGetData
		}

		return nil
	})

	errgrp.Go(func() error {
		if queryConfigMapped.CountAll && queryConfigMapped.Pagination.Page == 1 {
			count, err := s.getTotalCount(ctx, merchantId, datasetId, queryConfigMapped)
			if err != nil {
				logger.Error("failed to get total count", zap.String("error", err.Error()))
				return err
			}

			totalCount = &count
		}

		return nil
	})

	err = errgrp.Wait()
	if err != nil {
		return models.DatasetData{}, err
	}

	queryResultWithConfig := models.DatasetData{
		QueryResult: result,
		Title:       datasetMetaInfo.Title,
		Description: datasetMetaInfo.Description,
		TotalCount:  totalCount,
	}

	queryResultWithConfig.DatasetConfig.IsDrilldownEnabled = s.isDrilldownEnabled(datasetInfo)

	queryResultWithConfig.Metadata = datasetMetaData

	return queryResultWithConfig, nil
}

func (s *datasetService) ExecuteRawQuery(ctx context.Context, merchantId uuid.UUID, datasetId string, query string, queryParams map[string]interface{}) (models.DatasetData, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	queryParamsString := make(map[string]string)
	for key, value := range queryParams {
		queryParamsString[key] = fmt.Sprintf("%v", value)
	}
	queryParamsString[datasetConstants.ZampDatasetPrefix+datasetId] = datasetId

	var result dataplatformpkgmodels.QueryResult
	var err error

	switch s.serverDatasetConfig.DataplatformProvider {
	case datasetConstants.DataplatformProviderDatabricks:
		result, err = s.dataplatformService.Query(ctx, merchantId.String(), query, queryParamsString)
	case datasetConstants.DataplatformProviderPinot:
		result, err = s.dataplatformService.QueryRealTime(ctx, merchantId.String(), query, queryParamsString)
	default:
		return models.DatasetData{}, errors.ErrInvalidDataplatformProvider
	}

	if err != nil {
		logger.Error("failed to get the data", zap.String("error", err.Error()))
		return models.DatasetData{}, err
	}

	datasetMetaInfo, err := s.datasetStore.GetDatasetById(ctx, datasetId)
	if err != nil {
		logger.Error("failed to get dataset meta info", zap.String("error", err.Error()))
		return models.DatasetData{}, errors.ErrFailedToGetDatasetById
	}

	var datasetMetaData models.DatasetMetadataConfig
	if err := json.Unmarshal([]byte(datasetMetaInfo.Metadata), &datasetMetaData); err != nil {
		logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
		return models.DatasetData{}, errors.ErrFailedToUnmarshalMetadata
	}

	return models.DatasetData{
		QueryResult: result,
		Metadata:    datasetMetaData,
	}, nil
}

func (s *datasetService) GetRowDetailsByUUID(ctx context.Context, merchantId uuid.UUID, datasetId string, rowUUID string) (models.ParentDatasetInfo, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	rowDetailsResponse, parentDatasetResponse, err := s.fetchParentDatasetAndRowDetails(ctx, logger, merchantId.String(), datasetId, rowUUID)
	if err != nil {
		logger.Error("failed to fetch parents and row details", zap.Error(err))
		return models.ParentDatasetInfo{}, errors.ErrFailedToGetData
	}

	parentDatasets := s.mergeDatasetsInfoWithParents(rowDetailsResponse.ParentDatasets, parentDatasetResponse)

	updatedParentDatasets, err := s.fetchDatasetDetails(ctx, logger, parentDatasets)
	if err != nil {
		return models.ParentDatasetInfo{}, err
	}

	return models.ParentDatasetInfo{
		ParentDatasets: updatedParentDatasets,
	}, nil
}

func (s *datasetService) GetDatasetListing(ctx context.Context, merchantId uuid.UUID, params models.DatsetListingParams) ([]models.Dataset, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	var sortParams []storemodels.DatasetSortParam
	for _, sortParam := range params.SortParams {
		if !slices.Contains(storemodels.DatasetListingSortingColumns, sortParam.Column) {
			return nil, errors.ErrInvalidDatalistinSortColumn
		}
		sortParams = append(sortParams, storemodels.DatasetSortParam{Column: sortParam.Column, Desc: sortParam.Desc})
	}

	datasetsSchema, err := s.datasetStore.GetDatasetsAll(ctx, storemodels.DatasetFilters{
		OrganizationIds: []uuid.UUID{merchantId},
		CreatedBy:       params.CreatedBy,
		Page:            params.Pagination.Page,
		Limit:           params.Pagination.PageSize,
		SortParams:      sortParams,
		Type:            storemodels.UserVisibleDatasetTypes,
	})
	if err != nil {
		logger.Error("failed to get datasets for merchant", zap.String("merchant_id", merchantId.String()), zap.String("error", err.Error()))
		return nil, err
	}

	var datasets []models.Dataset
	for _, datasetSchema := range datasetsSchema {
		datasetModel := models.Dataset{}
		datasetModel.FromSchema(datasetSchema)
		datasets = append(datasets, datasetModel)
	}

	return datasets, nil
}

func (s *datasetService) UpdateDatasetData(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID, params models.UpdateDatasetDataParams) (models.DatasetAction, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetInfo, err := s.dataplatformService.GetDatasetMetadata(ctx, merchantId.String(), datasetId.String())
	if err != nil {
		logger.Error("failed to get dataset metadata", zap.String("error", err.Error()))
		return models.DatasetAction{}, errors.ErrFailedToGetData
	}

	columnDatatypes := make(map[string]dataplatformConstants.Datatype)
	for columnName, columnMetadata := range datasetInfo.Schema {
		columnDatatypes[columnName] = dataplatformConstants.Datatype(columnMetadata.Type)
	}

	datasetMetaInfo, err := s.datasetStore.GetDatasetById(ctx, datasetId.String())
	if err != nil {
		logger.Error("failed to get dataset meta info", zap.String("error", err.Error()))
		return models.DatasetAction{}, errors.ErrFailedToGetDatasetById
	}

	var datasetMetaData models.DatasetMetadataConfig
	if err := json.Unmarshal([]byte(datasetMetaInfo.Metadata), &datasetMetaData); err != nil {
		logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
		return models.DatasetAction{}, errors.ErrFailedToUnmarshalMetadata
	}

	customColumnConfig := make(map[string]querybuildermodels.CustomDataTypeConfig)

	queryConfig, err := s.mapUpdateDatasetDataParamsToQueryConfig(datasetId, params, columnDatatypes, customColumnConfig)
	if err != nil {
		logger.Error("failed to map update dataset data params to query config", zap.String("error", err.Error()))
		return models.DatasetAction{}, err
	}

	var dataplatformAction dataplatformactionmodels.CreateActionResponse

	switch params.SourceType {
	case datasetConstants.UpdateColumnSourceTypeUser:
		query, queryParams, err := s.queryBuilderService.ToFilterSQL(ctx, queryConfig.Filters)
		if err != nil {
			return models.DatasetAction{}, err
		}

		logger.Info("UPDATE DATASET DATA", zap.String("filter query", query), zap.Any("queryParams", queryParams))

		dataplatformAction, err = s.dataplatformService.UpdateDatasetData(ctx, dataplatformmodels.UpdateDatasetDataPayload{
			MerchantID: merchantId.String(),
			ActorId:    params.UserId.String(),
			ActionMetadataPayload: dataplatformactionmodels.UpdateDatasetDataActionPayload{
				DatasetId:    datasetId.String(),
				SqlCondition: query,
				UpdateValues: map[string]any{
					params.Update.Column: params.Update.Value,
				},
			},
		})
		if err != nil {
			return models.DatasetAction{}, err
		}
	case datasetConstants.UpdateColumnSourceTypeRule:
		query, queryParams, err := s.queryBuilderService.ToFilterSQL(ctx, queryConfig.Filters)
		if err != nil {
			return models.DatasetAction{}, err
		}

		logger.Info("UPDATE DATASET DATA FILTER SQL", zap.String("query", query), zap.Any("queryParams", queryParams))
		ruleParams, err := s.buildCreateRuleParams(params.SourceId, merchantId, datasetId, params, query, queryParams, columnDatatypes, customColumnConfig)
		if err != nil {
			logger.Error("failed to build create rule params", zap.String("error", err.Error()))
			return models.DatasetAction{}, err
		}

		err = s.ruleService.CreateRule(ctx, ruleParams)
		if err != nil {
			logger.Error("failed to create rule", zap.String("user_id", params.UserId.String()), zap.String("dataset_id", datasetId.String()), zap.String("error", err.Error()))
			return models.DatasetAction{}, err
		}

		dataplatformAction, err = s.handleRuleBasedDatasetUpdate(ctx, merchantId, datasetId, params)
		if err != nil {
			logger.Error("failed to handle rule based dataset update", zap.String("error", err.Error()))
			return models.DatasetAction{}, err
		}

	}

	action, err := s.dataplatformService.GetActionById(ctx, merchantId.String(), dataplatformAction.ActionID)
	if err != nil {
		return models.DatasetAction{}, err
	}

	isCompleted := false
	if slices.Contains(dataplatformactionconstants.ActionTerminationStatuses, action.ActionStatus) {
		isCompleted = true
	}

	err = s.datasetActionService.CreateDatasetAction(ctx, merchantId, storemodels.CreateDatasetActionParams{
		ActionId:    action.ID,
		ActionType:  string(action.ActionType),
		DatasetId:   datasetId,
		Status:      string(action.ActionStatus),
		Config:      action.ActionMetadata,
		ActionBy:    params.UserId,
		IsCompleted: isCompleted,
	})
	if err != nil {
		logger.Error("failed to create dataset action", zap.String("error", err.Error()))
		return models.DatasetAction{}, err
	}

	actionBy, err := uuid.Parse(action.ActorId)
	if err != nil {
		logger.Warn("failed to parse actor id", zap.String("error", err.Error()))
	}

	return models.DatasetAction{
		ActionId:    action.ID,
		ActionType:  action.ActionType,
		DatasetId:   datasetId,
		Status:      action.ActionStatus,
		Config:      action.ActionMetadata,
		ActionBy:    actionBy,
		IsCompleted: isCompleted,
	}, nil
}

func (s *datasetService) GetDatasetCount(ctx context.Context, merchantId uuid.UUID, params models.DatsetListingParams) (int64, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	totalCount, err := s.datasetStore.GetDatasetCount(ctx, storemodels.DatasetFilters{
		OrganizationIds: []uuid.UUID{merchantId},
		CreatedBy:       params.CreatedBy,
		Type:            storemodels.UserVisibleDatasetTypes,
	})

	if err != nil {
		logger.Error("failed to get dataset count", zap.String("error", err.Error()))
		return 0, err
	}

	return totalCount, nil
}

func (s *datasetService) RegisterDataset(ctx context.Context, merchantId uuid.UUID, userId uuid.UUID, datasetCreationInfo models.DatasetCreationInfo) (string, uuid.UUID, error) {
	datasetId := uuid.New()
	logger := apicontext.GetLoggerFromCtx(ctx)

	var actionResponse dataplatformactionmodels.CreateActionResponse
	err := s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {
		// Create dataset within transaction using the transaction store
		_, err := s.createDataset(ctx, ds, merchantId, userId, datasetId, datasetCreationInfo)
		if err != nil {
			logger.Error("failed to create dataset", zap.String("error", err.Error()))
			return errors.ErrFailedToRegisterDataset
		}

		// Create audience policy within transaction
		_, err = ds.CreateDatasetPolicy(ctx, datasetId, storemodels.AudienceTypeOrganization, merchantId, storemodels.PrivilegeDatasetAdmin)
		if err != nil {
			logger.Error("failed to create dataset policy", zap.String("error", err.Error()))
			return errors.ErrFailedToRegisterDataset
		}

		// Register dataset with platform service
		var regErr error
		switch datasetCreationInfo.DatasetType {
		case storemodels.DatasetTypeSource, storemodels.DatasetTypeBronze, storemodels.DatasetTypeStaged:
			actionResponse, regErr = s.dataplatformService.RegisterDataset(ctx, dataplatformmodels.RegisterDatasetPayload{
				MerchantID: merchantId.String(),
				ActionMetadataPayload: dataplatformactionmodels.RegisterDatasetActionPayload{
					MerchantId:       merchantId.String(),
					DatasetId:        datasetId.String(),
					DatasetConfig:    datasetCreationInfo.DatasetConfig,
					DatabricksConfig: datasetCreationInfo.DatabricksConfig,
					Provider:         datasetCreationInfo.Provider,
				},
			})
		case storemodels.DatasetTypeMV:
			actionResponse, regErr = s.dataplatformService.CreateMV(ctx, dataplatformmodels.CreateMVPayload{
				MerchantID: merchantId.String(),
				ActorId:    userId.String(),
				ActionMetadataPayload: dataplatformactionmodels.CreateMVActionPayload{
					Query:            datasetCreationInfo.MVConfig.Query,
					QueryParams:      datasetCreationInfo.MVConfig.QueryParams,
					ParentDatasetIds: datasetCreationInfo.MVConfig.ParentDatasetIds,
					MVDatasetId:      datasetId.String(),
					DedupColumns:     datasetCreationInfo.DatabricksConfig.DedupColumns,
					OrderByColumn:    datasetCreationInfo.DatabricksConfig.OrderByColumn,
				},
			})
		}

		return regErr
	})

	if err != nil {
		return "", uuid.Nil, err
	}

	return actionResponse.ActionID, datasetId, nil
}

func (s *datasetService) CopyDataset(ctx context.Context, merchantId uuid.UUID, userId uuid.UUID, params models.CopyDatasetParams) (string, uuid.UUID, error) {
	copyDatasetId := uuid.New()
	logger := apicontext.GetLoggerFromCtx(ctx)

	var actionResponse dataplatformactionmodels.CreateActionResponse
	err := s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {
		// Create dataset within transaction using the transaction store
		dataset, err := s.datasetStore.GetDatasetById(ctx, params.DatasetId)
		if err != nil {
			logger.Error("failed to get dataset", zap.String("error", err.Error()))
			return errors.ErrFailedToGetDataset
		}

		var datasetMetaData models.DatasetMetadataConfig
		if err := json.Unmarshal([]byte(dataset.Metadata), &datasetMetaData); err != nil {
			logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
			return errors.ErrFailedToUnmarshalMetadata
		}

		_, err = s.createDataset(ctx, ds, merchantId, userId, copyDatasetId, models.DatasetCreationInfo{
			DatasetTitle:       params.DatasetTitle,
			DatasetDescription: params.DatasetDescription,
			DatasetType:        dataset.Type,
			DatasetConfig:      datasetMetaData.DatasetConfig,
			DatabricksConfig:   datasetMetaData.DatabricksConfig,
			DisplayConfig:      datasetMetaData.DisplayConfig,
		})
		if err != nil {
			logger.Error("failed to create dataset", zap.String("error", err.Error()))
			return errors.ErrFailedToRegisterDataset
		}

		// Create audience policy within transaction
		_, err = ds.CreateDatasetPolicy(ctx, copyDatasetId, storemodels.AudienceTypeOrganization, merchantId, storemodels.PrivilegeDatasetAdmin)
		if err != nil {
			logger.Error("failed to create dataset policy", zap.String("error", err.Error()))
			return errors.ErrFailedToRegisterDataset
		}

		actionResponse, err = s.dataplatformService.CopyDataset(ctx, dataplatformmodels.CopyDatasetPayload{
			MerchantID: merchantId.String(),
			ActorId:    userId.String(),
			ActionMetadataPayload: dataplatformactionmodels.CopyDatasetActionPayload{
				OriginalDatasetId: params.DatasetId,
				NewDatasetId:      copyDatasetId.String(),
				MerchantId:        merchantId.String(),
			},
		})

		return err
	})
	if err != nil {
		return "", uuid.Nil, err
	}

	return actionResponse.ActionID, copyDatasetId, nil
}

func (s *datasetService) UpdateDataset(ctx context.Context, merchantId uuid.UUID, datasetId string, params models.UpdateDatasetParams) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	dataset, err := s.datasetStore.GetDatasetById(ctx, datasetId)
	if err != nil {
		logger.Error("failed to get dataset", zap.String("error", err.Error()))
		return "", err
	}

	metadata := models.DatasetMetadataConfig{}
	err = json.Unmarshal(dataset.Metadata, &metadata)
	if err != nil {
		logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
		return "", err
	}

	if params.Title != nil {
		dataset.Title = *params.Title
	}

	if params.Description != nil {
		dataset.Description = params.Description
	}

	if params.Type != nil {
		if !slices.Contains(storemodels.ValidDatasetTypes, storemodels.DatasetType(*params.Type)) {
			return "", errors.ErrInvalidDatasetType
		}
		dataset.Type = storemodels.DatasetType(*params.Type)
	}

	if params.DisplayConfig != nil {
		metadata.DisplayConfig = *params.DisplayConfig
	}

	if params.DatasetConfig != nil {
		metadata.DatasetConfig = *params.DatasetConfig
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		logger.Error("failed to marshal dataset metadata", zap.String("error", err.Error()))
		return "", errors.ErrFailedToMarshalMetadata
	}

	dataset.Metadata = metadataJSON
	dataset.UpdatedAt = time.Now().UTC()

	var actionResponse dataplatformactionmodels.CreateActionResponse

	err = s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {
		_, err = ds.UpdateDataset(ctx, *dataset)
		if err != nil {
			logger.Error("failed to update dataset", zap.String("error", err.Error()))
			return errors.ErrFailedToUpdateDataset
		}

		if params.DatasetConfig != nil {
			actionResponse, err = s.dataplatformService.UpdateDataset(ctx, dataplatformmodels.UpdateDatasetPayload{
				MerchantID: merchantId.String(),
				ActionMetadataPayload: dataplatformactionmodels.UpdateDatasetEvent{
					EventType: dataplatformactionconstants.UpdateDatasetEventTypeUpdateCustomColumn,
					EventData: dataplatformactionmodels.UpdateDatasetActionPayload{
						DatasetId:     datasetId,
						DatasetConfig: *params.DatasetConfig,
					},
				},
			})
			if err != nil {
				logger.Error("failed to update dataset action", zap.String("error", err.Error()))
				return errors.ErrFailedToUpdateDatasetAction
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return actionResponse.ActionID, nil
}

func (s *datasetService) RegisterDatasetJob(ctx context.Context, merchantId uuid.UUID, jobInfo dataplatformactionmodels.RegisterJobActionPayload) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actionResponse, err := s.dataplatformService.RegisterJob(ctx, dataplatformmodels.RegisterJobPayload{
		MerchantID: merchantId.String(),
		ActionMetadataPayload: dataplatformactionmodels.RegisterJobActionPayload{
			MerchantId:           merchantId.String(),
			JobType:              constants.DatabricksJobType(jobInfo.JobType),
			SourceType:           constants.DatabricksJobSourceType(jobInfo.SourceType),
			SourceValue:          jobInfo.SourceValue,
			DestinationType:      constants.DatabricksJobDestinationType(jobInfo.DestinationType),
			DestinationValue:     jobInfo.DestinationValue,
			TemplateId:           jobInfo.TemplateId,
			QuartzCronExpression: jobInfo.QuartzCronExpression,
		},
	})
	if err != nil {
		logger.Error("failed to register dataset job", zap.String("error", err.Error()))
		return "", errors.ErrFailedToRegisterDatasetJob
	}

	return actionResponse.ActionID, nil
}

func (s *datasetService) UpsertTemplate(ctx context.Context, merchantId uuid.UUID, templateConfig dataplatformactionmodels.UpsertTemplateActionPayload) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actionResponse, err := s.dataplatformService.UpsertTemplate(ctx, dataplatformmodels.UpsertTemplatePayload{
		MerchantID:            merchantId.String(),
		ActionMetadataPayload: templateConfig,
	})
	if err != nil {
		logger.Error("failed to upsert template", zap.String("error", err.Error()))
		return "", errors.ErrFailedToUpsertTemplate
	}

	return actionResponse.ActionID, nil
}

func (s *datasetService) GetOptionsForColumn(
	ctx context.Context,
	merchantId uuid.UUID,
	datasetId string,
	column string,
	filterType string,
	respectThreshold bool,
) ([]interface{}, error) {
	switch filterType {
	case datasetConstants.FilterTypeMultiSearch, datasetConstants.FilterTypeSelect:

		query := ""
		if respectThreshold {
			query = fmt.Sprintf(datasetConstants.GetDistinctValuesQuery, column, datasetId, datasetConstants.ZampIsDeletedColumn, datasetConstants.MultiSelectThreshold)
		} else {
			query = fmt.Sprintf(datasetConstants.GetDistinctValuesQueryWithoutLimit, column, datasetId, datasetConstants.ZampIsDeletedColumn)
		}

		var result dataplatformpkgmodels.QueryResult
		var err error

		switch s.serverDatasetConfig.DataplatformProvider {
		case datasetConstants.DataplatformProviderDatabricks:
			result, err = s.dataplatformService.Query(ctx, merchantId.String(), query, map[string]string{
				datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
			})
		case datasetConstants.DataplatformProviderPinot:
			result, err = s.dataplatformService.QueryRealTime(ctx, merchantId.String(), query, map[string]string{
				datasetConstants.ZampDatasetPrefix + datasetId: datasetId,
			})
		default:
			return nil, errors.ErrInvalidDataplatformProvider
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get distinct values for %s: %w", column, err)
		}

		var options []interface{}
		for _, row := range result.Rows {
			options = append(options, row[column])
		}
		return options, nil

	default:
		return []interface{}{}, nil
	}
}

func (s *datasetService) GetDatasetAudiences(ctx context.Context, datasetId uuid.UUID) ([]storemodels.ResourceAudiencePolicy, error) {
	return s.datasetStore.GetDatasetPolicies(ctx, datasetId)
}

func (s *datasetService) GetDatasetActions(ctx context.Context, merchantId uuid.UUID, filters storemodels.DatasetActionFilters) ([]models.DatasetAction, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actions, err := s.datasetActionService.GetDatasetActions(ctx, merchantId, filters)
	if err != nil {
		logger.Error("failed to get dataset actions", zap.String("merchant_id", merchantId.String()), zap.Any("filters", filters), zap.String("error", err.Error()))
		return nil, err
	}

	var datasetActions []models.DatasetAction
	for _, action := range actions {

		isCompleted := false
		if slices.Contains(dataplatformactionconstants.ActionTerminationStatuses, dataplatformactionconstants.ActionStatus(action.Status)) {
			isCompleted = true
		}

		datasetActions = append(datasetActions, models.DatasetAction{
			ActionId:    action.ActionId,
			ActionType:  dataplatformactionconstants.ActionType(action.ActionType),
			DatasetId:   action.DatasetId,
			Status:      dataplatformactionconstants.ActionStatus(action.Status),
			Config:      action.Config,
			ActionBy:    action.ActionBy,
			IsCompleted: isCompleted,
		})
	}
	return datasetActions, nil
}

func (s *datasetService) AddAudienceToDataset(ctx context.Context, datasetId uuid.UUID, audienceType storemodels.AudienceType, audienceId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	if !slices.Contains(storemodels.DatasetPrivileges, privilege) {
		ctxlogger.Error("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	switch audienceType {
	case storemodels.AudienceTypeUser:
		return s.addUserAudienceToDataset(ctx, datasetId, audienceId, privilege)
	case storemodels.AudienceTypeOrganization:
		return s.addOrganizationAudienceToDataset(ctx, datasetId, audienceId, privilege)
	case storemodels.AudienceTypeTeam:
		return s.addTeamAudienceToDataset(ctx, datasetId, audienceId, privilege)
	default:
		ctxlogger.Error("only user and organization audience is supported")
		return nil, fmt.Errorf("only user audience is supported")
	}
}

func (s *datasetService) addUserAudienceToDataset(ctx context.Context, datasetId uuid.UUID, userId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in the context")
	}

	var createdPolicy *storemodels.ResourceAudiencePolicy

	err := s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {

		// current user should be an admin on the dataset
		policies, err := s.datasetStore.GetDatasetPolicies(ctx, datasetId)
		if err != nil {
			ctxlogger.Info("failed to get dataset policies", zap.String("error", err.Error()))
			return err
		}

		// ensure that user is not already added on the dataset
		err = ensureAudienceNotAlreadyAdded(storemodels.AudienceTypeUser, userId, policies)
		if err != nil {
			ctxlogger.Info("user already exists on the dataset")
			return err
		}

		// ensure that the current user is an admin on the dataset
		err = ensureCurrentUsersAdminAccess(ctx, policies)
		if err != nil {
			ctxlogger.Info("current user does not have access to change permissions on the dataset")
			return err
		}

		createdPolicy, err = ds.CreateDatasetPolicy(ctx, datasetId, storemodels.AudienceTypeUser, userId, privilege)
		if err != nil {
			ctxlogger.Info("failed to create dataset policy", zap.String("error", err.Error()))
			return err
		}

		ctxlogger.Info("created dataset policy", zap.Any("policy", createdPolicy))

		return nil
	})

	return createdPolicy, err

}

func (s *datasetService) addOrganizationAudienceToDataset(ctx context.Context, datasetId uuid.UUID, organizationId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)
	_, currentUserId, orgIds := apicontext.GetAuthFromContext(ctx)

	if currentUserId == nil {
		ctxlogger.Info("no user id found in context")
		return nil, fmt.Errorf("no user id found in context")
	}

	if !slices.Contains(orgIds, organizationId) {
		ctxlogger.Info("current user does not have access to add organizations on the dataset")
		return nil, fmt.Errorf("current user does not have access to add organizations on the dataset")
	}

	var createdPolicy *storemodels.ResourceAudiencePolicy

	err := s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {

		// check if the dataset belongs to the organization
		dataset, err := ds.GetDatasetById(ctx, datasetId.String())
		if err != nil {
			ctxlogger.Info("failed to get dataset", zap.String("error", err.Error()))
			return err
		}

		if dataset.OrganizationId != organizationId {
			ctxlogger.Info("dataset does not belong to the organization")
			return fmt.Errorf("dataset does not belong to the organization")
		}

		// Get existing policies
		policies, err := s.datasetStore.GetDatasetPolicies(ctx, datasetId)
		if err != nil {
			ctxlogger.Info("failed to get dataset policies", zap.String("error", err.Error()))
			return err
		}

		// ensure that user is not already added on the dataset
		err = ensureAudienceNotAlreadyAdded(storemodels.AudienceTypeOrganization, organizationId, policies)
		if err != nil {
			ctxlogger.Info("user already exists on the dataset")
			return err
		}

		// ensure that the current user is an admin on the dataset
		err = ensureCurrentUsersAdminAccess(ctx, policies)
		if err != nil {
			ctxlogger.Info("current user does not have access to change permissions on the dataset")
			return err
		}

		// Create new policy
		createdPolicy, err = ds.CreateDatasetPolicy(ctx, datasetId, storemodels.AudienceTypeOrganization, organizationId, privilege)
		if err != nil {
			ctxlogger.Info("failed to create dataset policy", zap.String("error", err.Error()))
			return err
		}

		ctxlogger.Info("created dataset policy", zap.Any("policy", createdPolicy))

		return nil
	})

	return createdPolicy, err
}

func (s *datasetService) addTeamAudienceToDataset(ctx context.Context, datasetId uuid.UUID, teamId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in the context")
	}

	var createdPolicy *storemodels.ResourceAudiencePolicy
	err := s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {

		// ensure that the current user is an admin on the dataset
		policies, err := s.datasetStore.GetDatasetPolicies(ctx, datasetId)
		if err != nil {
			ctxlogger.Info("failed to get dataset policies", zap.String("error", err.Error()))
			return fmt.Errorf("something went wrong; please try again later")
		}

		// ensure that the current user is an admin on the dataset
		err = ensureCurrentUsersAdminAccess(ctx, policies)
		if err != nil {
			ctxlogger.Info("current user does not have access to change permissions on the dataset")
			return fmt.Errorf("current user does not have access to change permissions on the dataset")
		}

		// ensure that the team is not already added on the dataset
		err = ensureAudienceNotAlreadyAdded(storemodels.AudienceTypeTeam, teamId, policies)
		if err != nil {
			ctxlogger.Info("team already exists on the dataset")
			return fmt.Errorf("team already exists on the dataset")
		}

		p, err := s.datasetStore.CreateDatasetPolicy(ctx, datasetId, storemodels.AudienceTypeTeam, teamId, privilege)
		if err != nil {
			ctxlogger.Info("failed to create dataset policy", zap.String("error", err.Error()))
			return fmt.Errorf("failed to create dataset policy")
		}
		createdPolicy = p

		return nil
	})

	if err != nil {
		ctxlogger.Info("failed to add team audience to dataset", zap.String("error", err.Error()))
		return nil, err
	}

	return createdPolicy, nil
}

func (s *datasetService) BulkAddAudienceToDataset(ctx context.Context, datasetId uuid.UUID, payload models.BulkAddDatasetAudiencePayload) ([]*storemodels.ResourceAudiencePolicy, models.BulkAddDatasetAudienceErrors) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)
	var createdPolicies []*storemodels.ResourceAudiencePolicy
	var bulkErrors models.BulkAddDatasetAudienceErrors

	for _, audience := range payload.Audiences {
		policy, err := s.AddAudienceToDataset(ctx, datasetId, storemodels.AudienceType(audience.AudienceType), audience.AudienceId, storemodels.ResourcePrivilege(audience.Privilege))
		if err != nil {
			ctxlogger.Info("failed to add audience to dataset",
				zap.String("error", err.Error()),
				zap.String("audience_id", audience.AudienceId.String()))

			bulkErrors.Audiences = append(bulkErrors.Audiences, models.AddDatasetAudienceError{
				AudienceId:   audience.AudienceId,
				ErrorMessage: err.Error(),
			})
			continue
		}
		createdPolicies = append(createdPolicies, policy)
	}

	return createdPolicies, bulkErrors

}

func (s *datasetService) RemoveAudienceFromDataset(ctx context.Context, datasetId uuid.UUID, audienceId uuid.UUID) error {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return fmt.Errorf("no user ID found in the context")
	}

	// current user should be an admin on the dataset
	policies, err := s.datasetStore.GetDatasetPolicies(ctx, datasetId)
	if err != nil {
		ctxlogger.Info("failed to get dataset policies", zap.String("error", err.Error()))
		return err
	}

	// ensure that the current user is an admin on the dataset
	err = ensureCurrentUsersAdminAccess(ctx, policies)
	if err != nil {
		ctxlogger.Info("current user does not have access to change permissions on the dataset")
		return err
	}

	var policyToBeRemoved *storemodels.ResourceAudiencePolicy
	for _, policy := range policies {
		if policy.ResourceAudienceID == audienceId {
			policyToBeRemoved = &policy
			break
		}
	}

	if policyToBeRemoved == nil {
		ctxlogger.Info("policy not found", zap.String("audience_id", audienceId.String()))
		return fmt.Errorf("invalid audience id")
	}

	err = ensureUserIsNotChangingTheirOwnAdminPolicy(*policyToBeRemoved, policies, *currentUserId)
	if err != nil {
		ctxlogger.Info("current user is changing their own admin policy", zap.String("error", err.Error()))
		return err
	}

	err = s.datasetStore.DeleteDatasetPolicy(ctx, datasetId, storemodels.AudienceType(policyToBeRemoved.ResourceAudienceType), audienceId)
	if err != nil {
		ctxlogger.Error("failed to delete dataset policy", zap.Error(err))
		return err
	}

	return nil
}

func (s *datasetService) UpdateDatasetAudiencePrivilege(ctx context.Context, datasetId uuid.UUID, audienceId uuid.UUID, privilege storemodels.ResourcePrivilege) (*storemodels.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in the context")
	}

	// ensure valid privilege
	if !slices.Contains(storemodels.DatasetPrivileges, privilege) {
		ctxlogger.Info("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	// ensure that the current user is an admin on the dataset
	policies, err := s.datasetStore.GetDatasetPolicies(ctx, datasetId)
	if err != nil {
		ctxlogger.Info("failed to get dataset policies", zap.String("error", err.Error()))
		return nil, err
	}

	// ensure that the current user is an admin on the dataset
	err = ensureCurrentUsersAdminAccess(ctx, policies)
	if err != nil {
		ctxlogger.Info("current user does not have access to change permissions on the dataset")
		return nil, err
	}

	var policyToBeUpdated *storemodels.ResourceAudiencePolicy
	for _, policy := range policies {
		if policy.ResourceAudienceID == audienceId {
			policyToBeUpdated = &policy
			break
		}
	}
	if policyToBeUpdated == nil {
		ctxlogger.Info("policy not found", zap.String("audience_id", audienceId.String()))
		return nil, fmt.Errorf("invalid audience id")
	}

	err = ensureUserIsNotChangingTheirOwnAdminPolicy(*policyToBeUpdated, policies, *currentUserId)
	if err != nil {
		ctxlogger.Info("current user is changing their own admin policy", zap.String("error", err.Error()))
		return nil, err
	}

	updatedPolicy, err := s.datasetStore.UpdateDatasetPolicy(ctx, datasetId, audienceId, privilege)
	if err != nil {
		ctxlogger.Error("failed to update dataset policy", zap.Error(err))
		return nil, err
	}

	return updatedPolicy, nil
}

func (s *datasetService) UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error {
	return s.datasetActionService.UpdateDatasetActionStatus(ctx, actionId, status)
}

func (s *datasetService) GetRulesByDatasetColumns(ctx context.Context, organizationId uuid.UUID, datasetColumns []storemodels.DatasetColumn) (map[string]map[string][]rulemodels.Rule, error) {
	return s.ruleService.GetRules(ctx, storemodels.FilterRuleParams{
		OrganizationId: organizationId,
		DatasetColumns: datasetColumns,
	})
}

func (s *datasetService) CreateDatasetExportAction(ctx context.Context, merchantId uuid.UUID, datasetId string, params models.DatasetParams, userId uuid.UUID) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	workflowId := uuid.New().String()

	timestamp := time.Now().UTC().Format(datasetConstants.DatasetExportTimestampFormat)
	fileName := fmt.Sprintf(datasetConstants.DatasetExportFileNameFormat, timestamp, uuid.New().String())

	filePath := fmt.Sprintf(datasetConstants.DatasetExportFilePathFormat, datasetId, workflowId, fileName)

	metadata := models.ExportMetadata{
		FilePath: filePath,
	}

	datasetIdUUID, err := uuid.Parse(datasetId)
	if err != nil {
		logger.Error("failed to parse dataset id", zap.Error(err))
		return "", fmt.Errorf("failed to parse dataset id: %w", err)
	}

	orgIds := []uuid.UUID{merchantId}

	err = s.datasetActionService.CreateDatasetAction(ctx, merchantId, storemodels.CreateDatasetActionParams{
		ActionId:    workflowId,
		ActionType:  string(datasetactionconstants.ActionTypeDatasetExport),
		DatasetId:   datasetIdUUID,
		Status:      string(dataplatformactionconstants.ActionStatusInitiated),
		Config:      metadata,
		ActionBy:    userId,
		IsCompleted: false,
	})
	if err != nil {
		logger.Error("failed to create dataset action", zap.Error(err))
		return "", fmt.Errorf("failed to create dataset action: %w", err)
	}

	_, err = s.temporalService.ExecuteAsyncWorkflow(ctx, temporalmodels.ExecuteWorkflowParams{
		Options: temporalmodels.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: workersconstants.DefaultTaskQueueName,
		},
		Workflow: workersconstants.DatasetExportWorkflowName,
		Args: []interface{}{
			models.DatasetExportParams{
				QueryConfig: params,
				ExportPath:  metadata.FilePath,
			},
			datasetIdUUID,
			userId,
			orgIds,
			workflowId,
		},
	})

	if err != nil {
		logger.Error("failed to start export workflow", zap.Error(err))
		return "", fmt.Errorf("failed to start export: %w", err)
	}

	return workflowId, nil
}

func (s *datasetService) GetDownloadableDataExportUrl(ctx context.Context, workflowId string) (string, error) {

	datasetAction, err := s.datasetActionService.GetDatasetActionFromActionId(ctx, workflowId)
	if err != nil {
		return "", fmt.Errorf("failed to get dataset action: %w", err)
	}

	if datasetAction.Status != string(dataplatformactionconstants.ActionStatusSuccessful) {
		return "", fmt.Errorf("dataset action is not successful")
	}

	configMap, ok := datasetAction.Config.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid config type: %T", datasetAction.Config)
	}

	configBytes, err := json.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config map: %w", err)
	}

	var metadata models.ExportMetadata
	if err := json.Unmarshal(configBytes, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal config to ExportMetadata: %w", err)
	}

	signedURL, err := s.cloudService.GetSignedUrlToDownload(ctx, metadata.FilePath, []cloudservicemodels.GetDownloadsignedUrlConfigs{})

	if err != nil {
		return "", fmt.Errorf("failed to get signed url: %w", err)
	}

	return *signedURL, nil
}

func (s *datasetService) GetRulesByIds(ctx context.Context, ruleIds []string) ([]rulemodels.Rule, error) {
	ruleIdsUUID := make([]uuid.UUID, len(ruleIds))
	for i, ruleId := range ruleIds {
		ruleIdUUID, err := uuid.Parse(ruleId)
		if err != nil {
			return nil, fmt.Errorf("invalid rule id: %s", ruleId)
		}
		ruleIdsUUID[i] = ruleIdUUID
	}
	return s.ruleService.GetRuleByIds(ctx, ruleIdsUUID)
}

func (s *datasetService) UpdateRulePriority(ctx context.Context, orgId uuid.UUID, userId uuid.UUID, params models.UpdateRulePriorityParams) (models.DatasetAction, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	rules, err := s.ruleService.GetRules(ctx, storemodels.FilterRuleParams{
		OrganizationId: orgId,
		DatasetColumns: []storemodels.DatasetColumn{
			{
				DatasetId: params.DatasetId,
				Columns:   []string{params.Column},
			},
		},
	})
	if err != nil {
		logger.Error("failed to get rules", zap.Error(err))
		return models.DatasetAction{}, errors.ErrFailedToGetRule
	}

	var datasetRules []rulemodels.Rule
	if _, ok := rules[params.DatasetId.String()]; ok {
		if _, ok := rules[params.DatasetId.String()][params.Column]; ok {
			datasetRules = rules[params.DatasetId.String()][params.Column]
		}
	}

	err = s.validateAllDatasetColumnRulesPresentInPriorityUpdate(datasetRules, params)
	if err != nil {
		logger.Error("failed to validate all dataset column rules present in priority update", zap.Error(err))
		return models.DatasetAction{}, err
	}

	err = s.ruleService.UpdateRulePriority(ctx, storemodels.UpdateRulePriorityParams{
		DatasetId:    params.DatasetId,
		RulePriority: params.RulePriorities.RulePriority,
		UpdatedBy:    userId,
	})

	if err != nil {
		logger.Error("failed to update rule priority", zap.Error(err))
		return models.DatasetAction{}, err
	}

	createDatasetAction, err := s.handleRulePriorityUpdate(ctx, orgId, userId, params.DatasetId, params)
	if err != nil {
		logger.Error("failed to handle rule priority update", zap.Error(err))
		return models.DatasetAction{}, err
	}

	action, err := s.dataplatformService.GetActionById(ctx, orgId.String(), createDatasetAction.ActionID)
	if err != nil {
		return models.DatasetAction{}, err
	}

	isCompleted := false
	if slices.Contains(dataplatformactionconstants.ActionTerminationStatuses, dataplatformactionconstants.ActionStatus(action.ActionStatus)) {
		isCompleted = true
	}

	err = s.datasetActionService.CreateDatasetAction(ctx, orgId, storemodels.CreateDatasetActionParams{
		ActionId:    action.ID,
		ActionType:  string(action.ActionType),
		DatasetId:   params.DatasetId,
		Status:      string(action.ActionStatus),
		Config:      action.ActionMetadata,
		ActionBy:    userId,
		IsCompleted: isCompleted,
	})
	if err != nil {
		logger.Error("failed to create dataset action", zap.String("error", err.Error()))
		return models.DatasetAction{}, err
	}

	actionBy, err := uuid.Parse(action.ActorId)
	if err != nil {
		logger.Warn("failed to parse actor id", zap.String("error", err.Error()))
	}

	return models.DatasetAction{
		ActionId:    action.ID,
		ActionType:  action.ActionType,
		DatasetId:   params.DatasetId,
		Status:      action.ActionStatus,
		Config:      action.ActionMetadata,
		ActionBy:    actionBy,
		IsCompleted: isCompleted,
	}, nil
}

func (s *datasetService) getDatsetDags(ctx context.Context, merchantId string, datasetId string) (*dataplatformmodels.DAGNode, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	dags, err := s.dataplatformService.GetDags(ctx, merchantId)
	if err != nil {
		logger.Error("failed to get dags", zap.Error(err), zap.String("merchant_id", merchantId))
		return nil, fmt.Errorf("failed to get dags: %w", err)
	}

	if _, ok := dags[datasetId]; !ok {
		return nil, errors.ErrFailedToGetDatasetDags
	}

	return dags[datasetId], nil
}

func (s *datasetService) InitiateFilePreparationForDatasetImport(ctx context.Context, datasetId uuid.UUID, fileId uuid.UUID) (datasetActionId *uuid.UUID, err error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, currentUserId, orgIds := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		ctxlogger.Error("no user ID in context")
		return nil, fmt.Errorf("unauthorized")
	}

	if len(orgIds) == 0 {
		ctxlogger.Error("no organization ID in context")
		return nil, fmt.Errorf("unauthorized")
	}

	workflowId := uuid.New()
	actionId := uuid.New()

	err = s.datasetStore.WithTx(ctx, func(dsStore store.Store) error {

		datasetFileUpload, errr := dsStore.CreateDatasetFileUpload(ctx, &storemodels.DatasetFileUpload{
			ID:                   uuid.New(),
			DatasetID:            datasetId,
			FileUploadID:         fileId,
			FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusPending,
			Metadata:             json.RawMessage(`{}`),
		})
		if errr != nil {
			ctxlogger.Error("failed to create dataset file upload", zap.Error(err))
			return errr
		}

		fileImportActionConfig := models.FileImportDatasetActionConfig{
			Version:    1,
			FileId:     fileId,
			WorkflowId: workflowId,
			WorkflowInitPayload: models.FileImportWorkflowInitPayload{
				DatasetId:           datasetId,
				UserId:              *currentUserId,
				OrganizationId:      orgIds[0],
				DatasetActionId:     actionId,
				FileUploadId:        fileId,
				DatasetFileUploadId: datasetFileUpload.ID,
			},
		}

		errr = dsStore.CreateDatasetAction(ctx, orgIds[0], storemodels.CreateDatasetActionParams{
			ActionId:   actionId.String(),
			ActionType: string(datasetactionconstants.ActionTypeDatasetFileImport),
			ActionBy:   *currentUserId,
			DatasetId:  datasetId,
			Status:     string(dataplatformactionconstants.ActionStatusInitiated),
			Config:     fileImportActionConfig,
		})

		if errr != nil {
			ctxlogger.Error("failed to create dataset action", zap.Error(err))
			return errr
		}

		_, errr = s.temporalService.ExecuteAsyncWorkflow(ctx, temporalmodels.ExecuteWorkflowParams{
			Options: temporalmodels.StartWorkflowOptions{
				ID:        workflowId.String(),
				TaskQueue: workersconstants.DefaultTaskQueueName,
			},
			Workflow: workersconstants.DatasetFileImportWorkflowName,
			Args: []interface{}{
				fileImportActionConfig.WorkflowInitPayload,
			},
		})

		if errr != nil {
			ctxlogger.Error("failed to execute workflow", zap.Error(err))
			return errr
		}

		return nil
	})

	if err != nil {
		ctxlogger.Error("failed to import data from file", zap.Error(err))
		return nil, err
	}

	return &actionId, nil
}

func (s *datasetService) UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error {
	return s.datasetActionService.UpdateDatasetActionConfig(ctx, actionId, config)
}

func (s *datasetService) CreateDatasetFileUpload(ctx context.Context, datasetId uuid.UUID, fileId uuid.UUID, metadata json.RawMessage) error {
	datasetFileUpload := storemodels.DatasetFileUpload{
		ID:                   uuid.New(),
		DatasetID:            datasetId,
		FileUploadID:         fileId,
		FileAllignmentStatus: storemodels.DatasetFileAllignmentStatusPending,
		Metadata:             metadata,
	}

	logger := apicontext.GetLoggerFromCtx(ctx)
	logger.Info("creating dataset file upload", zap.Any("datasetFileUpload", datasetFileUpload))

	_, err := s.datasetStore.CreateDatasetFileUpload(ctx, &datasetFileUpload)
	if err != nil {
		return fmt.Errorf("failed to create dataset file upload: %w", err)
	}

	return nil
}

func (s *datasetService) GetDatasetFileUploads(ctx context.Context, datasetId uuid.UUID) ([]models.DatasetFileUpload, error) {
	datasetFileUploads, err := s.datasetStore.GetDatasetFileUploadByDatasetId(ctx, datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset file uploads: %w", err)
	}

	mapDatasetFileUploads := make(map[uuid.UUID]storemodels.DatasetFileUpload)
	var fileIds []uuid.UUID

	for _, datasetFileUpload := range datasetFileUploads {
		fileIds = append(fileIds, datasetFileUpload.FileUploadID)
		mapDatasetFileUploads[datasetFileUpload.FileUploadID] = datasetFileUpload
	}

	fileUploads, err := s.fileImportService.GetFileUploadByIds(ctx, fileIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get file uploads: %w", err)
	}

	mapFileUploads := make(map[uuid.UUID]storemodels.FileUpload)

	for _, fileUpload := range fileUploads {
		mapFileUploads[fileUpload.ID] = fileUpload
	}

	return s.buildDatasetFileUploads(ctx, mapDatasetFileUploads, mapFileUploads)
}

func (s *datasetService) UpdateDatasetFileUploadStatus(ctx context.Context, datasetFileUploadId uuid.UUID, params models.UpdateDatasetFileUploadParams) error {
	_, err := s.datasetStore.UpdateDatasetFileUploadStatus(ctx, datasetFileUploadId, params.FileAllignmentStatus, params.Metadata)
	if err != nil {
		return fmt.Errorf("failed to update dataset file upload status: %w", err)
	}

	return nil
}

func (s *datasetService) GetFileUploadPreview(ctx context.Context, fileUploadId uuid.UUID) (datasetFileUploadsModels.DatasetPreview, error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	datasetFileUpload, err := s.datasetStore.GetDatasetFileUploadById(ctx, fileUploadId)
	if err != nil {
		ctxlogger.Error("failed to get dataset file uploads", zap.Error(err))
		return datasetFileUploadsModels.DatasetPreview{}, fmt.Errorf("failed to get dataset file uploads: %w", err)
	}

	if datasetFileUpload.ID == uuid.Nil {
		ctxlogger.Error("no dataset file upload found")
		return datasetFileUploadsModels.DatasetPreview{}, fmt.Errorf("no dataset file upload found")
	}

	if datasetFileUpload.FileAllignmentStatus != storemodels.DatasetFileAllignmentStatusCompleted {
		ctxlogger.Error("dataset file upload is not completed")
		return datasetFileUploadsModels.DatasetPreview{}, fmt.Errorf("Preview is not ready yet.")
	}

	metadata := datasetFileUpload.Metadata

	var datasetFileUploadMetadata datasetFileUploadsModels.DatasetFileUploadMetadata
	err = json.Unmarshal(metadata, &datasetFileUploadMetadata)
	if err != nil {
		ctxlogger.Error("failed to unmarshal dataset file upload metadata", zap.Error(err))
		return datasetFileUploadsModels.DatasetPreview{}, fmt.Errorf("failed to unmarshal dataset file upload metadata: %w", err)
	}

	return datasetFileUploadMetadata.DataPreview, nil
}

func (s *datasetService) GetDatasetImportPath(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID) (*models.FileImportConfig, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	dag, err := s.getDatsetDags(ctx, merchantId.String(), datasetId.String())
	if err != nil {
		logger.Error("failed to get dags", zap.Error(err), zap.String("datasetId", datasetId.String()))
		return nil, err
	}

	importUrl, err := dag.GetImportFilePath()
	if err != nil {
		return nil, err
	}

	bucket, folderPrefix := helper.ExtractBucketNameAndFolderPrefix(importUrl)
	bronzeSourcePath, err := s.s3Client.GetSampleFilePathFromFolder(ctx, bucket, folderPrefix)
	if err != nil {
		return nil, err
	}

	return &models.FileImportConfig{
		IsFileImportEnabled: true,
		BronzeSourceBucket:  bucket,
		BronzeSourcePath:    bronzeSourcePath,
		BronzeSourceConfig:  map[string]interface{}{},
	}, nil
}

func (s *datasetService) ImportDataFromFile(ctx context.Context, merchantId uuid.UUID, datasetId uuid.UUID, fileUploadId uuid.UUID) (err error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx).With(zap.String("dataset_id", datasetId.String()), zap.String("file_upload_id", fileUploadId.String()))

	datasetFileUploads, err := s.datasetStore.GetDatasetFileUploadByDatasetId(ctx, datasetId)
	if err != nil {
		return fmt.Errorf("Failed to get dataset file uploads: %w", err)
	}

	var datasetFileUpload *storemodels.DatasetFileUpload
	for _, dfu := range datasetFileUploads {
		if dfu.FileUploadID == fileUploadId {
			datasetFileUpload = &dfu
		}
	}

	if datasetFileUpload == nil {
		return fmt.Errorf("The given file is not associated with the dataset")
	}

	if datasetFileUpload.FileAllignmentStatus != storemodels.DatasetFileAllignmentStatusCompleted {
		return fmt.Errorf("The given file is not prepared for dataset import")
	}

	var datasetFileUploadMetadata storemodels.DatasetFileUploadMetadata
	err = json.Unmarshal(datasetFileUpload.Metadata, &datasetFileUploadMetadata)
	if err != nil {
		ctxlogger.Error("failed to unmarshal dataset file upload metadata", zap.Error(err))
		return fmt.Errorf("Unexpected. Invalid file upload metadata.")
	}

	bucketName := datasetFileUploadMetadata.TransformedDataBucket
	filePath := datasetFileUploadMetadata.TransformedDataPath

	// validate if the file exists
	fileDetails, err := s.s3Client.GetFileDetails(ctx, bucketName, filePath)
	if err != nil {
		ctxlogger.Error("failed to retrieve the file details", zap.Error(err))
		return fmt.Errorf("Unexpected. Failed to retrieve the given file details.")
	}

	if fileDetails.Size == 0 {
		ctxlogger.Error("the given file is empty")
		return fmt.Errorf("The given file is empty")
	}

	var destinationBucketName string

	datasetImportPath, err := s.GetDatasetImportPath(ctx, merchantId, datasetId)
	if err != nil {
		ctxlogger.Error("failed to get dataset import path", zap.Error(err))
		return fmt.Errorf("Unexpected. Failed to get dataset import path. Please try again later.")
	}

	if datasetImportPath == nil {
		ctxlogger.Error("dataset import path is nil")
		return fmt.Errorf("Unexpected. Failed to get dataset import path. Please try again later.")
	}

	destinationBucketName = datasetImportPath.BronzeSourceBucket

	// extract destination folder from the destinationfile path; destination folder is the same folder as the file path but without the file name
	desitnationFilePath := helper.GetRenamedFilePath(datasetImportPath.BronzeSourcePath, datasetFileUpload.ID.String())

	ctxlogger.Info("copying file to destination bucket", zap.String("source_bucket_name", bucketName), zap.String("source_file_path", filePath), zap.String("destination_bucket_name", destinationBucketName), zap.String("destination_path", desitnationFilePath))

	// copy the file to destination bucket
	err = s.s3Client.CopyFile(ctx, bucketName, filePath, destinationBucketName, desitnationFilePath)
	if err != nil {
		ctxlogger.Error("failed to copy the file", zap.Error(err))
		return fmt.Errorf("Unexpected. Failed to import the given file. Please try again later.")
	}

	return nil
}

func (s *datasetService) DeleteDataset(ctx context.Context, merchantId uuid.UUID, datasetId string) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	dataset, err := s.datasetStore.GetDatasetById(ctx, datasetId)
	if err != nil {
		logger.Error("failed to get dataset", zap.Error(err))
		return "", errors.ErrFailedToGetDatasetById
	}

	now := time.Now()
	dataset.DeletedAt = &now

	err = s.datasetStore.WithDatasetTransaction(ctx, func(ds store.DatasetStore) error {
		err = s.datasetStore.DeleteDataset(ctx, *dataset)
		if err != nil {
			logger.Error("failed to delete dataset", zap.Error(err))
			return fmt.Errorf("failed to delete dataset: %w", err)
		}

		// TODO - Create delete action in dataplatform
		// actionId, err := s.dataplatformService.DeleteDataset(ctx, dataplatformmodels.DeleteDatasetPayload{
		// 	MerchantID: merchantId.String(),
		// 	DatasetID:  datasetId,
		// })
		// if err != nil {
		// 	logger.Error("failed to delete dataset in dataplatform", zap.Error(err))
		// 	return "", fmt.Errorf("failed to delete dataset in dataplatform: %w", err)
		// }

		return nil
	})

	if err != nil {
		return "", err
	}

	return "", nil
}

func (s *datasetService) GetDatasetDisplayConfig(ctx context.Context, merchantId uuid.UUID, datasetId string) ([]models.DisplayConfig, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	datasetMetaInfo, err := s.datasetStore.GetDatasetById(ctx, datasetId)
	if err != nil {
		logger.Error("failed to get dataset meta info", zap.String("error", err.Error()))
		return nil, errors.ErrFailedToGetDatasetById
	}

	var datasetMetaData models.DatasetMetadataConfig
	if err := json.Unmarshal([]byte(datasetMetaInfo.Metadata), &datasetMetaData); err != nil {
		logger.Error("failed to unmarshal dataset metadata", zap.String("error", err.Error()))
		return nil, errors.ErrFailedToUnmarshalMetadata
	}

	// If display config exists, return it
	if len(datasetMetaData.DisplayConfig) > 0 {
		return datasetMetaData.DisplayConfig, nil
	}

	// If display config doesn't exist, build it from dataplatform metadata
	datasetInfo, err := s.dataplatformService.GetDatasetMetadata(ctx, merchantId.String(), datasetId)
	if err != nil {
		logger.Error("failed to get dataset metadata", zap.String("error", err.Error()))
		return nil, errors.ErrFailedToGetDatasetMetadata
	}

	// Build display config from schema
	var displayConfig []models.DisplayConfig
	for columnName := range datasetInfo.Schema {
		displayConfig = append(displayConfig, models.DisplayConfig{
			Column:     columnName,
			IsHidden:   false,
			IsEditable: false,
		})
	}

	return displayConfig, nil
}

func (s *datasetService) getFlattenedDatasetPolicies(ctx context.Context, datasetId uuid.UUID) ([]storemodels.FlattenedResourceAudiencePolicy, error) {

	policies, err := s.datasetStore.GetFlattenedResourceAudiencePolicies(ctx, storemodels.FlattenedResourceAudiencePoliciesFilters{
		ResourceIds:   []uuid.UUID{datasetId},
		ResourceTypes: []string{string(storemodels.ResourceTypeDataset)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset policies: %w", err)
	}

	return policies, nil
}
