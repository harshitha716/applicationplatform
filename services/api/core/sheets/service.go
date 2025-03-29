package sheets

import (
	"context"
	"fmt"
	"time"

	datasetsService "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	sheetmodels "github.com/Zampfi/application-platform/services/api/core/sheets/models"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/cache"
	querybuilderconstants "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type SheetsServiceStore interface {
	store.SheetStore
}

type SheetsService interface {
	GetSheetsByPageId(ctx context.Context, pageId uuid.UUID) ([]models.Sheet, error)
	GetSheetById(ctx context.Context, pageId uuid.UUID) (*models.Sheet, error)
	GetSheetFilterConfigById(ctx context.Context, orgId uuid.UUID, sheetId uuid.UUID) (*sheetmodels.FilterOptionsConfig, error)
	CreateSheet(ctx context.Context, sheet sheetmodels.Sheet) (*models.Sheet, error)
	UpdateSheet(ctx context.Context, sheet *sheetmodels.Sheet) (*models.Sheet, error)
}

type sheetsService struct {
	store          SheetsServiceStore
	datasetService datasetsService.DatasetService
	cacheClient    cache.CacheClient
}

func NewSheetsService(appStore store.Store, datasetService datasetsService.DatasetService, cacheService cache.CacheClient) *sheetsService {
	return &sheetsService{store: appStore, datasetService: datasetService, cacheClient: cacheService}
}

func (s *sheetsService) GetSheetsByPageId(ctx context.Context, pageId uuid.UUID) ([]models.Sheet, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	sheets, err := s.store.GetSheetsAll(ctx, models.SheetFilters{PageIds: []uuid.UUID{pageId}, SortParams: []models.SheetSortParams{{Column: "created_at", Desc: false}}})
	if err != nil {
		ctxLogger.Error("failed to get pages", zap.Error(err))
		return nil, err
	}

	return sheets, nil
}

func (s *sheetsService) GetSheetById(ctx context.Context, pageId uuid.UUID) (*models.Sheet, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	sheet, err := s.store.GetSheetById(ctx, pageId)
	if err != nil {
		ctxLogger.Error("failed to get page", zap.Error(err))
		return nil, err
	}

	return sheet, nil
}

func (s *sheetsService) GetSheetFilterConfigById(ctx context.Context, orgId uuid.UUID, sheetId uuid.UUID) (*sheetmodels.FilterOptionsConfig, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	sheet, err := s.store.GetSheetById(ctx, sheetId)
	if err != nil {
		logger.Error("failed to get sheet", zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get sheet: %w", err)
	}

	sheetModel := sheetmodels.Sheet{}
	if err := sheetModel.FromDB(sheet); err != nil {
		return nil, fmt.Errorf("failed to convert sheet to model: %w", err)
	}

	cacheKey, err := s.cacheClient.FormatKey("sheet_filter_config", sheetModel.ID.String())
	if err != nil {
		logger.Warn("failed to format cache key", zap.Error(err))
	}

	sheetFilterConfig := &sheetmodels.FilterOptionsConfig{}
	if err := s.cacheClient.Get(ctx, cacheKey, sheetFilterConfig); err != nil {
		logger.Warn("failed to get sheet filter config from cache", zap.Error(err))
	} else {
		return sheetFilterConfig, nil
	}

	sheetFilterConfig, err = s.getSheetFilterConfigFromDB(ctx, orgId, sheetModel)
	if err != nil {
		logger.Error("failed to get sheet filter config from db", zap.Error(err))
		return nil, err
	}

	if err := s.cacheClient.Set(ctx, cacheKey, sheetFilterConfig, time.Minute*30); err != nil {
		logger.Warn("failed to set sheet filter config in cache", zap.Error(err))
	}

	return sheetFilterConfig, nil
}

func (s *sheetsService) CreateSheet(ctx context.Context, sheet sheetmodels.Sheet) (*models.Sheet, error) {
	sheetModel, err := sheet.ToDB()
	if err != nil {
		return nil, fmt.Errorf("failed to convert sheet to model: %w", err)
	}
	return s.store.CreateSheet(ctx, *sheetModel)
}

func isRangeOperator(op string) bool {
	switch op {
	case querybuilderconstants.GreaterThanOperator, querybuilderconstants.GreaterThanOrEqualOperator, querybuilderconstants.LessThanOperator, querybuilderconstants.LessThanOrEqualOperator, querybuilderconstants.EqualOperator, querybuilderconstants.NotEqualOperator:
		return true
	default:
		return false
	}
}

func (s *sheetsService) UpdateSheet(ctx context.Context, updatedSheet *sheetmodels.Sheet) (*models.Sheet, error) {
	// Get the sheet
	sheet, err := s.store.GetSheetById(ctx, updatedSheet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sheet: %w", err)
	}

	sheetModel := sheetmodels.Sheet{}
	if err := sheetModel.FromDB(sheet); err != nil {
		return nil, fmt.Errorf("failed to convert sheet to model: %w", err)
	}

	if updatedSheet.SheetConfig.Version != "" {
		sheetModel.SheetConfig = updatedSheet.SheetConfig
	}

	if updatedSheet.Name != "" {
		sheetModel.Name = updatedSheet.Name
	}

	if updatedSheet.Description != nil {
		sheetModel.Description = updatedSheet.Description
	}

	if updatedSheet.PageId != uuid.Nil {
		sheetModel.PageId = updatedSheet.PageId
	}

	sheetModelDB, err := sheetModel.ToDB()
	if err != nil {
		return nil, fmt.Errorf("failed to convert sheet to model: %w", err)
	}
	return s.store.UpdateSheet(ctx, sheetModelDB)
}

func (s *sheetsService) getSheetFilterConfigFromDB(ctx context.Context, orgId uuid.UUID, sheetModel sheetmodels.Sheet) (*sheetmodels.FilterOptionsConfig, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	sheetFilterConfig := sheetmodels.FilterOptionsConfig{
		NativeFilterConfig: make([]sheetmodels.FilterOptionsModel, len(sheetModel.SheetConfig.NativeFilterConfig)),
	}

	filterErrgrp := errgroup.Group{}

	for i, filterConfig := range sheetModel.SheetConfig.NativeFilterConfig {
		index := i
		cfg := filterConfig

		sheetFilterConfig.NativeFilterConfig[index] = sheetmodels.FilterOptionsModel{
			Id:             cfg.Id,
			Name:           cfg.Name,
			FilterType:     cfg.FilterType,
			WidgetsInScope: cfg.WidgetsInScope,
			Targets:        cfg.Targets,
			DataType:       cfg.DataType,
			Options:        []interface{}{},
			DefaultValue:   cfg.DefaultValue,
		}

		filterErrgrp.Go(func() error {
			targetErrgrp := errgroup.Group{}
			optionsCh := make(chan []interface{}, len(cfg.Targets))

			done := make(chan struct{})
			optionSet := make(map[interface{}]struct{})

			go func() {
				defer close(done)
				for options := range optionsCh {
					for _, opt := range options {
						optionSet[opt] = struct{}{}
					}
				}
			}()

			for _, target := range cfg.Targets {
				tgt := target
				targetErrgrp.Go(func() error {
					options, err := s.datasetService.GetOptionsForColumn(
						ctx,
						orgId,
						tgt.DatasetId.String(),
						tgt.Column,
						cfg.FilterType,
						false,
					)
					if err != nil {
						return fmt.Errorf("failed to get options for target %s: %w", tgt.Column, err)
					}
					optionsCh <- options
					return nil
				})
			}

			err := targetErrgrp.Wait()
			close(optionsCh)
			if err != nil {
				return err
			}

			<-done

			var uniqueOptions []interface{}
			for opt := range optionSet {
				uniqueOptions = append(uniqueOptions, opt)
			}

			sheetFilterConfig.NativeFilterConfig[index].Options = uniqueOptions
			return nil
		})
	}

	if err := filterErrgrp.Wait(); err != nil {
		logger.Error("failed to populate filter options", zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to populate filter options: %w", err)
	}

	return &sheetFilterConfig, nil
}
