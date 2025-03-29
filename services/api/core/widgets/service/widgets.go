package widgets

import (
	"context"

	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	widgetconstants "github.com/Zampfi/application-platform/services/api/core/widgets/constants"
	"github.com/Zampfi/application-platform/services/api/core/widgets/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type WidgetsServiceStore interface {
	store.WidgetStore
}

type WidgetsService interface {
	GetWidgetInstance(ctx context.Context, widgetInstanceID uuid.UUID) (models.WidgetInstance, error)
	GetWidgetInstanceData(ctx context.Context, orgId uuid.UUID, widgetInstanceID uuid.UUID, params models.GetWidgetInstanceDataQueryParams) ([]datasetmodels.DatasetData, error)
	CreateWidgetInstance(ctx context.Context, widgetInstance models.WidgetInstance) (*models.WidgetInstance, error)
	UpdateWidgetInstance(ctx context.Context, widgetInstance models.WidgetInstance) (*models.WidgetInstance, error)
}

type widgetsService struct {
	store          WidgetsServiceStore
	datasetService datasetservice.DatasetService
}

func NewWidgetsService(appStore store.Store, datasetService datasetservice.DatasetService) *widgetsService {
	return &widgetsService{store: appStore, datasetService: datasetService}
}

func (s *widgetsService) GetWidgetInstance(ctx context.Context, widgetInstanceID uuid.UUID) (models.WidgetInstance, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)
	widgetInstance, err := s.store.GetWidgetInstanceByID(ctx, widgetInstanceID)
	if err != nil {
		ctxLogger.Error("failed to get widget instance", zap.String("error", err.Error()))
		return models.WidgetInstance{}, err
	}

	widgetInstanceModel := models.WidgetInstance{}
	if err := widgetInstanceModel.FromDB(&widgetInstance); err != nil {
		ctxLogger.Error("failed to parse widget instance", zap.String("error", err.Error()))
		return models.WidgetInstance{}, err
	}

	return widgetInstanceModel, nil
}

func (s *widgetsService) GetWidgetInstanceData(ctx context.Context, orgId uuid.UUID, widgetInstanceID uuid.UUID, params models.GetWidgetInstanceDataQueryParams) ([]datasetmodels.DatasetData, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	datasetFilters := make(map[string]models.WidgetFilters)
	for _, filter := range params.Filters {
		datasetFilters[filter.DatasetID] = filter
	}

	widgetInstance, err := s.store.GetWidgetInstanceByID(ctx, widgetInstanceID)
	if err != nil {
		ctxLogger.Error("failed to get widget instance", zap.String("error", err.Error()))
		return []datasetmodels.DatasetData{}, err
	}

	widgetInstanceModel := models.WidgetInstance{}
	if err := widgetInstanceModel.FromDB(&widgetInstance); err != nil {
		ctxLogger.Error("failed to convert widget instance to model", zap.String("error", err.Error()))
		return []datasetmodels.DatasetData{}, err
	}

	datasetParamsBuilder, err := NewDatasetParamsBuilder(widgetInstanceModel.WidgetType)
	if err != nil {
		ctxLogger.Error("failed to get dataset params builder", zap.String("error", err.Error()))
		return []datasetmodels.DatasetData{}, err
	}

	timeColumnMap := make(map[string]string)
	for _, timeColumn := range params.TimeColumns {
		timeColumnMap[timeColumn.DatasetID] = timeColumn.Column
	}

	datasetParams, err := datasetParamsBuilder.ToDatasetParams(&widgetInstanceModel, models.DatasetBuilderParams{
		Filters:     datasetFilters,
		TimeColumns: timeColumnMap,
		Periodicity: params.Periodicity,
		Currency:    params.Currency,
	})
	if err != nil {
		ctxLogger.Error("failed to get dataset params", zap.String("error", err.Error()))
		return []datasetmodels.DatasetData{}, err
	}

	ctxLogger.Info("dataset params", zap.Any("dataset_params", datasetParams))

	errGroup := errgroup.Group{}
	type datasetResult struct {
		data datasetmodels.DatasetData
		ref  string
	}

	dataResultsChan := make(chan datasetResult, len(datasetParams))
	for ref, params := range datasetParams {
		params := params
		ref := ref
		//params.Params.FxCurrency = nil
		params.Params.Pagination = &datasetmodels.Pagination{
			Page:     1,
			PageSize: widgetconstants.MAX_PAGE_SIZE,
		}
		errGroup.Go(func() error {
			data, err := s.datasetService.GetDataByDatasetId(ctx, orgId, params.DatasetID, params.Params)
			if err != nil {
				ctxLogger.Error("failed to get data by dataset id", zap.String("error", err.Error()), zap.String("dataset_id", params.DatasetID))
				return err
			}

			tagColumn, ok := data.GetTagsColumns()
			if ok {
				data = s.flattenTags(&data, tagColumn)
			}

			data = s.addRefToDataResults(&data, &ref, len(datasetParams) > 1)

			dataResultsChan <- datasetResult{data: data, ref: ref}
			return nil
		})
	}

	err = errGroup.Wait()
	if err != nil {
		ctxLogger.Error("failed to get data", zap.String("error", err.Error()))
		return []datasetmodels.DatasetData{}, err
	}

	close(dataResultsChan)

	resultMap := make(map[string]datasetmodels.DatasetData)
	for result := range dataResultsChan {
		resultMap[result.ref] = result.data
	}

	dataResults := make([]datasetmodels.DatasetData, 0, len(datasetParams))
	for _, mapping := range widgetInstanceModel.DataMappings.Mappings {
		if data, exists := resultMap[mapping.Ref]; exists {
			dataResults = append(dataResults, data)
		}
	}
	return dataResults, nil
}

func (s *widgetsService) CreateWidgetInstance(ctx context.Context, widgetInstance models.WidgetInstance) (*models.WidgetInstance, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	widgetInstanceDB, err := widgetInstance.ToDB()
	if err != nil {
		ctxLogger.Error("failed to convert widget instance to db model", zap.String("error", err.Error()))
		return nil, err
	}

	widgetInstanceDB, err = s.store.CreateWidgetInstance(ctx, widgetInstanceDB)
	if err != nil {
		ctxLogger.Error("failed to create widget instance", zap.String("error", err.Error()))
		return nil, err
	}

	widgetInstanceModel := models.WidgetInstance{}
	if err := widgetInstanceModel.FromDB(widgetInstanceDB); err != nil {
		ctxLogger.Error("failed to convert widget instance to model", zap.String("error", err.Error()))
		return nil, err
	}

	return &widgetInstanceModel, nil
}

func (s *widgetsService) UpdateWidgetInstance(ctx context.Context, updatedInstance models.WidgetInstance) (*models.WidgetInstance, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// Get the widget instance
	widgetInstance, err := s.GetWidgetInstance(ctx, updatedInstance.ID)
	if err != nil {
		ctxLogger.Error("failed to get widget instance", zap.String("error", err.Error()))
		return nil, err
	}

	if updatedInstance.WidgetType != "" {
		widgetInstance.WidgetType = updatedInstance.WidgetType
	}
	if updatedInstance.DataMappings.Version != "" {
		widgetInstance.DataMappings = updatedInstance.DataMappings
	}
	if updatedInstance.Title != "" {
		widgetInstance.Title = updatedInstance.Title
	}
	if updatedInstance.DisplayConfig != nil {
		widgetInstance.DisplayConfig = updatedInstance.DisplayConfig
	}

	if updatedInstance.SheetID != uuid.Nil {
		widgetInstance.SheetID = updatedInstance.SheetID
	}

	widgetInstanceDB, err := widgetInstance.ToDB()
	if err != nil {
		ctxLogger.Error("failed to convert widget instance to db model", zap.String("error", err.Error()))
		return nil, err
	}

	// Update in the database
	updatedWidgetInstanceDB, err := s.store.UpdateWidgetInstance(ctx, widgetInstanceDB)
	if err != nil {
		ctxLogger.Error("failed to update widget instance", zap.String("error", err.Error()))
		return nil, err
	}

	// Convert back to API model
	var result models.WidgetInstance
	if err := result.FromDB(updatedWidgetInstanceDB); err != nil {
		ctxLogger.Error("failed to convert updated widget instance to model", zap.String("error", err.Error()))
		return nil, err
	}

	return &result, nil
}
