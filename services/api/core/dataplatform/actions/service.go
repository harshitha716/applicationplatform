package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	serviceconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	templates "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/templates"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	data "github.com/Zampfi/application-platform/services/api/core/dataplatform/data"
	dataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datamodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	helper "github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	dataplatformconstants "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/databricks"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"go.uber.org/zap"
)

type ActionService interface {
	CreateAction(ctx context.Context, payload models.CreateActionPayload) (models.CreateActionResponse, error)
	UpdateAction(ctx context.Context, jobStatusUpdate dataplatformmodels.DatabricksJobStatusUpdatePayload) (models.Action, error)
	GetActionById(ctx context.Context, merchantId string, actionId string) (models.Action, error)
}

type actionService struct {
	dataService data.DataService
}

func InitActionService(dataService data.DataService) ActionService {
	return &actionService{
		dataService: dataService,
	}
}

func (s *actionService) getCreateMVJobTemplate(ctx context.Context, providerId string, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	jobPayload := map[string]string{}
	warehouseId, err := s.dataService.GetDatabricksWarehouseId(ctx, providerId)
	if err != nil {
		logger.Error(errors.GettingDatabricksWarehouseIdFailedErrMessage, zap.Error(err))
		return nil, errors.ErrGettingDatabricksWarehouseIdFailed
	}

	createMvPayload, ok := payload.ActionMetadataPayload.(models.CreateMVActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of typeCreateMVActionPayload")))
		return nil, errors.ErrInvalidActionMetadataPayload
	}

	queryMetadata, err := s.dataService.ProcessParamsForQuery(ctx, payload.MerchantID, createMvPayload.QueryParams, dataplatformconstants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrTemplateParsingFailed
	}

	mvQuery, err := helper.FillQueryTemplate(ctx, createMvPayload.Query, queryMetadata.Params)
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrTemplateParsingFailed
	}

	mvQuery, err = s.dataService.TranslateQuery(ctx, mvQuery, dataplatformconstants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.QueryTranslationFailedErrMessage, zap.Error(err))
		return nil, errors.ErrQueryTranslationFailed
	}

	jobTemplate := templates.GetMVJobTemplate(s.dataService.GetDataPlatformConfig(), payload.MerchantID, warehouseId)

	createMVPayload := models.CreateMVPayload{
		ParentDatasetIds: createMvPayload.ParentDatasetIds,
		Query:            mvQuery,
		MerchantId:       payload.MerchantID,
		DatasetId:        createMvPayload.MVDatasetId,
		DedupColumns:     createMvPayload.DedupColumns,
		OrderByColumn:    createMvPayload.OrderByColumn,
	}

	createMVPayloadStr, err := helper.ConvertToJSONString(createMVPayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	jobPayload[serviceconstants.CreateMVParams] = createMVPayloadStr
	jobPayload[serviceconstants.DatasetIdParam] = createMVPayload.DatasetId

	return s.addJobPayloadToTemplate(jobTemplate, jobPayload), nil
}

func (s *actionService) getRegisterDatasetJobTemplate(ctx context.Context, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	jobPayload := map[string]string{}

	registerDatasetPayload, ok := payload.ActionMetadataPayload.(models.RegisterDatasetActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of type RegisterDatasetActionPayload")))
		return nil, errors.ErrInvalidActionMetadataPayload
	}

	jobTemplate := templates.GetRegisterDatasetJobTemplate(s.dataService.GetDataPlatformConfig(), payload.MerchantID, registerDatasetPayload.DatasetId)

	registerDatasetParamsStr, err := helper.ConvertToJSONString(registerDatasetPayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	jobPayload[serviceconstants.RegisterDatasetParams] = registerDatasetParamsStr
	return s.addJobPayloadToTemplate(jobTemplate, jobPayload), nil
}

func (s *actionService) getRegisterJobJobTemplate(ctx context.Context, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	jobPayload := map[string]string{}

	registerJobPayload, ok := payload.ActionMetadataPayload.(models.RegisterJobActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of type RegisterJobActionPayload")))
		return nil, errors.ErrInvalidActionMetadataPayload
	}

	jobTemplate := templates.GetRegisterJobJobTemplate(s.dataService.GetDataPlatformConfig(), payload.MerchantID, registerJobPayload.DestinationValue)

	registerJobParamsStr, err := helper.ConvertToJSONString(registerJobPayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	jobPayload[serviceconstants.RegisterJobParams] = registerJobParamsStr
	return s.addJobPayloadToTemplate(jobTemplate, jobPayload), nil
}

func (s *actionService) getUpsertTemplateJobTemplate(ctx context.Context, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	jobPayload := map[string]string{}

	upsertTemplatePayload, ok := payload.ActionMetadataPayload.(models.UpsertTemplateActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of type UpsertTemplateActionPayload")))
		return nil, errors.ErrInvalidActionMetadataPayload
	}

	jobTemplate := templates.GetUpsertTemplateJobTemplate(s.dataService.GetDataPlatformConfig(), payload.MerchantID, upsertTemplatePayload.Id)

	upsertTemplateParamsStr, err := helper.ConvertToJSONString(upsertTemplatePayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	jobPayload[serviceconstants.UpsertTemplateParams] = upsertTemplateParamsStr
	return s.addJobPayloadToTemplate(jobTemplate, jobPayload), nil
}

func (s *actionService) getCopyDatasetJobTemplate(ctx context.Context, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	copyDatasetPayload, ok := payload.ActionMetadataPayload.(models.CopyDatasetActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of type CopyDatasetActionPayload")))
		return nil, errors.ErrInvalidActionMetadataPayload
	}

	jobTemplate := templates.GetCopyDatasetJobTemplate(
		s.dataService.GetDataPlatformConfig(),
		copyDatasetPayload.NewDatasetId,
	)

	copyDatasetParamsStr, err := helper.ConvertToJSONString(copyDatasetPayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	jobPayload := map[string]string{
		serviceconstants.CopyDatasetParams: copyDatasetParamsStr,
	}

	return s.addJobPayloadToTemplate(jobTemplate, jobPayload), nil
}

func (s *actionService) getUpdateDatasetJobTemplate(ctx context.Context, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	jobPayload := map[string]string{}

	updateDatasetEventPayload, ok := payload.ActionMetadataPayload.(models.UpdateDatasetEvent)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of type UpdateDatasetActionPayload")))
		return nil, errors.ErrInvalidActionMetadataPayload
	}

	jobTemplate := templates.GetUpdateDatasetJobTemplate(s.dataService.GetDataPlatformConfig(), payload.MerchantID, updateDatasetEventPayload.EventData.DatasetId)

	updateDatasetParamsStr, err := helper.ConvertToJSONString(updateDatasetEventPayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return nil, errors.ErrJSONUnmarshallingFailed
	}

	jobPayload[serviceconstants.UpdateDatasetEventParams] = updateDatasetParamsStr
	jobPayload[serviceconstants.DatasetIdParam] = updateDatasetEventPayload.EventData.DatasetId

	return s.addJobPayloadToTemplate(jobTemplate, jobPayload), nil
}

func (s *actionService) addJobPayloadToTemplate(jobTemplate *jobs.SubmitRun, jobPayload map[string]string) *jobs.SubmitRun {
	// add default modules src to the job payload
	jobPayload[serviceconstants.DataPlatformModulesSrc] = s.dataService.GetDataPlatformConfig().ActionsConfig.DataPlatformModulesSrc
	for _, task := range jobTemplate.Tasks {
		if task.NotebookTask != nil {
			task.NotebookTask.BaseParameters = jobPayload

		}
		if task.SqlTask != nil {
			task.SqlTask.Parameters = jobPayload
		}
	}

	return jobTemplate
}

func (s *actionService) getJobTemplate(ctx context.Context, providerId string, payload models.CreateActionPayload) (*jobs.SubmitRun, error) {
	switch payload.ActionType {

	case serviceconstants.ActionTypeCreateMV:
		return s.getCreateMVJobTemplate(ctx, providerId, payload)
	case serviceconstants.ActionTypeRegisterDataset:
		return s.getRegisterDatasetJobTemplate(ctx, payload)
	case serviceconstants.ActionTypeRegisterJob:
		return s.getRegisterJobJobTemplate(ctx, payload)
	case serviceconstants.ActionTypeUpsertTemplate:
		return s.getUpsertTemplateJobTemplate(ctx, payload)
	case serviceconstants.ActionTypeUpdateDataset:
		return s.getUpdateDatasetJobTemplate(ctx, payload)
	case serviceconstants.ActionTypeCopyDataset:
		return s.getCopyDatasetJobTemplate(ctx, payload)
	}

	return nil, errors.ErrInvalidActionType
}

func (s *actionService) saveAction(ctx context.Context, databricksService databricks.DatabricksService, actionId string, providerId string, payload models.CreateActionPayload) error {
	logger := apicontext.GetLoggerFromCtx(ctx)
	actionMetadataJSONString, err := helper.ConvertToJSONStringWithReplacements(payload.ActionMetadataPayload, map[string]string{
		"'": "\"",
	})
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return errors.ErrJSONUnmarshallingFailed
	}

	actionsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.ActionsTableName)
	query, err := helper.FillQueryTemplate(ctx, serviceconstants.QueryCreateAction, map[string]string{
		serviceconstants.ActionsTableNameQueryParam:  actionsTableName,
		serviceconstants.ActionIdColumnName:          actionId,
		serviceconstants.ActionWorkspaceIdColumnName: providerId,
		serviceconstants.ActionTypeColumnName:        string(payload.ActionType),
		serviceconstants.ActionMetadataColumnName:    actionMetadataJSONString,
		serviceconstants.ActionStatusColumnName:      string(serviceconstants.ActionStatusInitiated),
		serviceconstants.ActionActorIdColumnName:     payload.ActorId,
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return errors.ErrTemplateParsingFailed
	}

	_, err = databricksService.Query(ctx, actionsTableName, query)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return errors.ErrQueryingDatabricksFailed
	}

	return nil
}

func (s *actionService) handleOneTimeJobAction(ctx context.Context, databricksService databricks.DatabricksService, providerId string, payload models.CreateActionPayload) (models.SubmitActionResponse, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	jobTemplate, err := s.getJobTemplate(ctx, providerId, payload)
	if err != nil {
		logger.Error(errors.GettingJobTemplateFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, err
	}

	submitResponse, err := databricksService.SubmitOneTimeJob(ctx, jobTemplate)
	if err != nil {
		logger.Error(errors.SubmittingJobFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, err
	}
	return models.SubmitActionResponse{RunId: submitResponse.RunId}, nil
}

func (s *actionService) handleUpdateDatasetDataAction(ctx context.Context, databricksService databricks.DatabricksService, payload models.CreateActionPayload) (models.SubmitActionResponse, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	updatePayload, ok := payload.ActionMetadataPayload.(models.UpdateDatasetDataActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of typeUpdateActionPayload")))
		return models.SubmitActionResponse{}, errors.ErrInvalidActionMetadataPayload
	}

	// Get The Update Job Id For The Dataset
	jobMappingsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, dataconstants.JobMappingsTableName)
	jobsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, dataconstants.JobsTableName)
	getUpdateJobIdQuery, err := helper.FillQueryTemplate(ctx, serviceconstants.QueryGetJobIdForDatasetForJobType, map[string]string{
		dataconstants.JobMappingsTableNameQueryParam:  jobMappingsTableName,
		dataconstants.JobsTableNameQueryParam:         jobsTableName,
		dataconstants.JobMappingSourceValueColumnName: updatePayload.DatasetId,
		dataconstants.JobMappingSourceTypeColumnName:  string(dataconstants.JobMappingTypeDataset),
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, err
	}

	queryResponse, err := databricksService.Query(ctx, jobMappingsTableName, getUpdateJobIdQuery)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, err
	}

	if len(queryResponse.Rows) == 0 {
		logger.Error(errors.GettingJobIdForDatasetForJobTypeFailedErrMessage, zap.Error(errors.ErrGettingJobIdForDatasetForJobTypeFailed))
		return models.SubmitActionResponse{}, errors.ErrGettingJobIdForDatasetForJobTypeFailed
	}

	jobIdModel := datamodels.JobIdModel{}
	jobIdJSONString, err := json.Marshal(queryResponse.Rows[0])
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, errors.ErrJSONUnmarshallingFailed
	}
	err = json.Unmarshal(jobIdJSONString, &jobIdModel)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, errors.ErrJSONUnmarshallingFailed
	}

	// Trigger the Update Job
	updatedValuesJSONString, err := helper.ConvertToJSONString(updatePayload)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, errors.ErrJSONUnmarshallingFailed
	}

	runResponse, err := databricksService.RunNow(ctx, jobs.RunNow{
		JobId: jobIdModel.JobId,
		JobParameters: map[string]string{
			serviceconstants.UpdateDatasetDataParams: updatedValuesJSONString,
		},
	})
	if err != nil {
		logger.Error(errors.SubmittingJobFailedErrMessage, zap.Error(err))
		return models.SubmitActionResponse{}, errors.ErrSubmittingJobFailed
	}

	return models.SubmitActionResponse{RunId: runResponse.RunId}, nil
}

func (s *actionService) handleJobAction(ctx context.Context, databricksService databricks.DatabricksService, payload models.CreateActionPayload) (models.SubmitActionResponse, error) {

	switch payload.ActionType {
	case serviceconstants.ActionTypeUpdateDatasetData:
		return s.handleUpdateDatasetDataAction(ctx, databricksService, payload)
	}

	return models.SubmitActionResponse{}, errors.ErrInvalidActionType
}

func (s *actionService) submitAction(ctx context.Context, databricksService databricks.DatabricksService, providerId string, payload models.CreateActionPayload) (models.SubmitActionResponse, error) {

	if slices.Contains(serviceconstants.SubmitOneTimeJobActions, payload.ActionType) {
		return s.handleOneTimeJobAction(ctx, databricksService, providerId, payload)
	}

	if slices.Contains(serviceconstants.SubmitJobActions, payload.ActionType) {
		return s.handleJobAction(ctx, databricksService, payload)
	}

	return models.SubmitActionResponse{}, errors.ErrInvalidActionType
}

func (s *actionService) updateActionRunId(ctx context.Context, databricksService databricks.DatabricksService, actionId string, runId int64) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actionsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.ActionsTableName)
	query, err := helper.FillQueryTemplate(ctx, serviceconstants.QueryUpdateActionRunId, map[string]string{
		serviceconstants.ActionsTableNameQueryParam: actionsTableName,
		serviceconstants.ActionIdColumnName:         actionId,
		serviceconstants.ActionRunIdColumnName:      strconv.FormatInt(runId, 10),
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return err
	}

	_, err = databricksService.Query(ctx, actionsTableName, query)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return err
	}
	return nil
}

func verifyCountryColumn(ctx context.Context, updatedValues interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	valueString, ok := updatedValues.(string)
	if !ok {
		logger.Error(errors.InvalidCountryValueErrMessage, zap.Error(errors.ErrInvalidCountryValue))
		return errors.ErrInvalidCountryValue
	}

	if !slices.Contains(serviceconstants.Countries, valueString) {
		logger.Error(errors.InvalidCountryValueErrMessage, zap.Error(errors.ErrInvalidCountryValue))
		return errors.ErrInvalidCountryValue
	}

	return nil

}

func verifyBankColumn(ctx context.Context, updatedValues interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	valueString, ok := updatedValues.(string)
	if !ok {
		logger.Error(errors.InvalidBankValueErrMessage, zap.Error(errors.ErrInvalidBankValue))
		return errors.ErrInvalidBankValue
	}

	if !slices.Contains(serviceconstants.Banks, valueString) {
		logger.Error(errors.InvalidBankValueErrMessage, zap.Error(errors.ErrInvalidBankValue))
		return errors.ErrInvalidBankValue
	}

	return nil
}

func verifyCurrencyColumn(ctx context.Context, updatedValues interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	valueString, ok := updatedValues.(string)
	if !ok {
		logger.Error(errors.InvalidCurrencyValueErrMessage, zap.Error(errors.ErrInvalidCurrencyValue))
		return errors.ErrInvalidCurrencyValue
	}

	if !slices.Contains(serviceconstants.Currencies, valueString) {
		logger.Error(errors.InvalidCurrencyValueErrMessage, zap.Error(errors.ErrInvalidCurrencyValue))
		return errors.ErrInvalidCurrencyValue
	}

	return nil
}

func verifyTagsColumn(ctx context.Context, updatedValues interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	_, ok := updatedValues.(string)
	if !ok {
		logger.Error(errors.InvalidTagsValueErrMessage, zap.Error(errors.ErrInvalidTagsValue))
		return errors.ErrInvalidTagsValue
	}

	return nil
}

func verifyAmountColumn(ctx context.Context, updatedValues interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)
	switch updatedValues.(type) {
	case int:
		return nil
	case int32:
		return nil
	case float32:
		return nil
	case int64:
		return nil
	case float64:
		return nil
	default:
		logger.Error(errors.InvalidAmountValueErrMessage, zap.Error(errors.ErrInvalidAmountValue))
		return errors.ErrInvalidAmountValue
	}

}

func getSourceColumnUpdateValue(actorId string) models.SourceColumnUpdateValue {
	return models.SourceColumnUpdateValue{
		SourceType:      "user",
		SourceId:        actorId,
		SourceUpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func getSourceColumnName(columnName string) string {
	return "_zamp_source_json_" + columnName
}

func (s *actionService) handleValidationsAndSourceUpdates(ctx context.Context, payload models.CreateActionPayload) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	if payload.ActionType != serviceconstants.ActionTypeUpdateDatasetData {
		return nil
	}

	actionMetadataPayload, ok := payload.ActionMetadataPayload.(models.UpdateDatasetDataActionPayload)
	if !ok {
		logger.Error(errors.InvalidActionMetadataPayloadErrMessage, zap.Error(fmt.Errorf("actionMetadataPayload is not of type UpdateDatasetDataActionPayload")))
		return errors.ErrInvalidActionMetadataPayload
	}

	datasetConfig, err := s.dataService.GetDatasetConfig(ctx, payload.MerchantID, actionMetadataPayload.DatasetId)
	if err != nil {
		logger.Error(errors.GettingDatasetConfigFailedErrMessage, zap.Error(err))
		return err
	}

	// only adding source column for tags
	for columnName := range actionMetadataPayload.UpdateValues {
		if _, ok := datasetConfig.Columns[columnName]; !ok {
			continue
		}
		if datasetConfig.Columns[columnName].CustomType == constants.DatabricksColumnCustomTypeTags {
			sourceColumnUpdateValue, err := helper.ConvertToJSONString(getSourceColumnUpdateValue(payload.ActorId))
			if err != nil {
				logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
				continue
			}
			actionMetadataPayload.UpdateValues[getSourceColumnName(columnName)] = sourceColumnUpdateValue
		}
	}

	var validationError error
	// check if the updated column exists in the dataset column map if yes then check the type and if no then return error
	for columnName, value := range actionMetadataPayload.UpdateValues {
		if _, ok := datasetConfig.Columns[columnName]; !ok {
			continue
		}

		switch datasetConfig.Columns[columnName].CustomType {
		case constants.DatabricksColumnCustomTypeCurrency:
			validationError = verifyCurrencyColumn(ctx, value)
		case constants.DatabricksColumnCustomTypeCountry:
			validationError = verifyCountryColumn(ctx, value)
		case constants.DatabricksColumnCustomTypeAmount:
			validationError = verifyAmountColumn(ctx, value)
		case constants.DatabricksColumnCustomTypeTags:
			validationError = verifyTagsColumn(ctx, value)
		case constants.DatabricksColumnCustomTypeBank:
			validationError = verifyBankColumn(ctx, value)
		}

		if validationError != nil {
			return validationError
		}
	}
	return nil
}

func (s *actionService) CreateAction(ctx context.Context, payload models.CreateActionPayload) (models.CreateActionResponse, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := s.handleValidationsAndSourceUpdates(ctx, payload)
	if err != nil {
		logger.Error(errors.ActionValidationFailedErrMessage, zap.Error(err))
		return models.CreateActionResponse{}, err
	}

	providerId, err := s.dataService.GetDataProviderIdForMerchant(payload.MerchantID, dataplatformconstants.ProviderTypeDatabricks)
	if err != nil {
		logger.Error(errors.GettingDataProviderIdForMerchantFailedErrMessage, zap.Error(err))
		return models.CreateActionResponse{}, errors.ErrGettingDataProviderIdForMerchantFailed
	}

	databricksService, err := s.dataService.GetDatabricksServiceForMerchant(ctx, payload.MerchantID)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return models.CreateActionResponse{}, errors.ErrProviderServiceNotFound
	}

	actionId := helper.GenerateUUIDWithUnderscores()

	err = s.saveAction(ctx, databricksService, actionId, providerId, payload)
	if err != nil {
		logger.Error(errors.CreatingActionFailedErrMessage, zap.Error(err))
		return models.CreateActionResponse{}, errors.ErrCreatingActionFailed
	}

	submitResponse, err := s.submitAction(ctx, databricksService, providerId, payload)
	if err != nil {
		logger.Error(errors.SubmittingActionFailedErrMessage, zap.Error(err))
		return models.CreateActionResponse{}, errors.ErrSubmittingActionFailed
	}

	err = s.updateActionRunId(ctx, databricksService, actionId, submitResponse.RunId)
	if err != nil {
		logger.Error(errors.UpdatingActionRunIdFailedErrMessage, zap.Error(err))
		return models.CreateActionResponse{}, errors.ErrUpdatingActionRunIdFailed
	}

	return models.CreateActionResponse{
		ActionID: actionId,
	}, nil
}

func (s *actionService) handleJobStatusUpdate(ctx context.Context, databricksService databricks.DatabricksService, runDetails *jobs.Run, runId int64, workspaceId string) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	finalStatus := runDetails.State.ResultState
	var actionStatus serviceconstants.ActionStatus
	switch finalStatus {
	case jobs.RunResultStateSuccess:
		actionStatus = serviceconstants.ActionStatusSuccessful
	case jobs.RunResultStateFailed:
		actionStatus = serviceconstants.ActionStatusFailed
	default:
		logger.Info("JOB CURRENT STATUS IS NOT SUCCESS OR FAILED")
		return errors.ErrJobStatusNotSuccessOrFailed
	}

	actionsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.ActionsTableName)
	query, err := helper.FillQueryTemplate(ctx, serviceconstants.QueryUpdateActionStatus, map[string]string{
		serviceconstants.ActionsTableNameQueryParam:  actionsTableName,
		serviceconstants.ActionRunIdColumnName:       strconv.FormatInt(runId, 10),
		serviceconstants.ActionWorkspaceIdColumnName: workspaceId,
		serviceconstants.ActionStatusColumnName:      string(actionStatus),
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return err
	}

	_, err = databricksService.Query(ctx, actionsTableName, query)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return err
	}
	return nil
}

func (s *actionService) getActionDetails(ctx context.Context, databricksService databricks.DatabricksService, runId int64, workspaceId string) (models.Action, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actionsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.ActionsTableName)
	query, err := helper.FillQueryTemplate(ctx, serviceconstants.QueryGetActionByRunId, map[string]string{
		serviceconstants.ActionsTableNameQueryParam:  actionsTableName,
		serviceconstants.ActionRunIdColumnName:       strconv.FormatInt(runId, 10),
		serviceconstants.ActionWorkspaceIdColumnName: workspaceId,
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	actionRawResponse, err := databricksService.Query(ctx, actionsTableName, query)
	if err != nil {
		logger.Error(errors.GettingActionByRunIdFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	if len(actionRawResponse.Rows) == 0 {
		logger.Error(errors.GettingActionByRunIdFailedErrMessage, zap.Error(errors.ErrGettingActionByRunIdFailed))
		return models.Action{}, errors.ErrGettingActionByRunIdFailed
	}

	action := models.Action{}
	actionJSONString, err := json.Marshal(actionRawResponse.Rows[0])
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}
	err = json.Unmarshal(actionJSONString, &action)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}
	return action, nil
}

func (s *actionService) UpdateAction(ctx context.Context, jobStatusUpdate dataplatformmodels.DatabricksJobStatusUpdatePayload) (models.Action, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	runId := jobStatusUpdate.Run.RunId
	workspaceId := strconv.FormatInt(jobStatusUpdate.WorkspaceId, 10)

	databricksService, err := s.dataService.GetDatabricksServiceForProvider(ctx, workspaceId)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	action, err := s.getActionDetails(ctx, databricksService, runId, workspaceId)
	if err != nil {
		logger.Error(errors.GettingActionByRunIdFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	// jobs can go from failed to successful, in case of retries
	if action.ActionStatus != serviceconstants.ActionStatusInitiated && action.ActionStatus != serviceconstants.ActionStatusFailed {
		logger.Info(errors.ActionStatusNotInitiatedOrFailedErrMessage)
		return action, nil
	}

	runDetails, err := databricksService.GetRunDetails(ctx, runId)
	if err != nil {
		logger.Error(errors.GettingRunDetailsFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	err = s.handleJobStatusUpdate(ctx, databricksService, runDetails, runId, workspaceId)
	if err != nil {
		logger.Error(errors.UpdatingActionStatusFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	action, err = s.getActionDetails(ctx, databricksService, runId, workspaceId)
	if err != nil {
		logger.Error(errors.GettingActionByRunIdFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	return action, nil
}

func (s *actionService) GetActionById(ctx context.Context, merchantId string, actionId string) (models.Action, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	actionsTableName := helpers.BuildDatabricksTableName(s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksCatalog, s.dataService.GetDataPlatformConfig().DatabricksConfig.ZampDatabricksPlatformSchema, serviceconstants.ActionsTableName)
	query, err := helper.FillQueryTemplate(ctx, serviceconstants.QueryGetActionById, map[string]string{
		serviceconstants.ActionsTableNameQueryParam: actionsTableName,
		serviceconstants.ActionIdColumnName:         actionId,
	})
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	databricksService, err := s.dataService.GetDatabricksServiceForMerchant(ctx, merchantId)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	actionRawResponse, err := databricksService.Query(ctx, actionsTableName, query)
	if err != nil {
		logger.Error(errors.GettingActionByRunIdFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}

	if len(actionRawResponse.Rows) == 0 {
		logger.Error(errors.GettingActionByRunIdFailedErrMessage, zap.Error(errors.ErrGettingActionByRunIdFailed))
		return models.Action{}, errors.ErrGettingActionByRunIdFailed
	}

	action := models.Action{}
	actionJSONString, err := json.Marshal(actionRawResponse.Rows[0])
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}
	err = json.Unmarshal(actionJSONString, &action)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return models.Action{}, err
	}
	return action, nil
}
