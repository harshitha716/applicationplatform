package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SheetStore interface {
	GetSheetById(ctx context.Context, sheetId uuid.UUID) (*models.Sheet, error)
	GetSheetsAll(ctx context.Context, filters models.SheetFilters) ([]models.Sheet, error)
	CreateSheet(ctx context.Context, sheet models.Sheet) (*models.Sheet, error)
	UpdateSheet(ctx context.Context, sheet *models.Sheet) (*models.Sheet, error)
}

func (s *appStore) GetSheetById(ctx context.Context, sheetId uuid.UUID) (*models.Sheet, error) {

	var sheet models.Sheet
	err := s.client.WithContext(ctx).Preload("WidgetInstances").Model(sheet).Where("sheet_id = ?", sheetId).First(&sheet).Error
	if err != nil {
		return nil, err
	}
	return &sheet, nil
}

func (s *appStore) GetSheetsAll(ctx context.Context, filters models.SheetFilters) ([]models.Sheet, error) {
	db := s.client.WithContext(ctx)

	sheets := []models.Sheet{}
	if len(filters.PageIds) > 0 {
		db = db.Where("page_id IN (?)", filters.PageIds)
	}

	if filters.IncludeWidgetInstances {
		db = db.Preload("WidgetInstances", func(db *gorm.DB) *gorm.DB {
			return db.Order("widget_instances.created_at ASC")
		})
	}

	if len(filters.SheetIds) > 0 {
		db = db.Where("sheet_id IN (?)", filters.SheetIds)
	}

	if filters.Page > 0 {
		db = db.Offset((filters.Page - 1) * filters.Limit)
	}

	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}

	if len(filters.SortParams) > 0 {
		for _, sort := range filters.SortParams {
			db = db.Order(clause.OrderByColumn{
				Column: clause.Column{Name: sort.Column},
				Desc:   sort.Desc,
			})
		}
	}

	result := db.Find(&sheets)
	if result.Error != nil {
		return nil, result.Error
	}

	return sheets, nil

}

func (s *appStore) CreateSheet(ctx context.Context, sheet models.Sheet) (*models.Sheet, error) {
	sheet.ID = uuid.New()
	sheet.FractionalIndex = 0
	err := s.client.WithContext(ctx).Create(&sheet).Error
	if err != nil {
		return nil, err
	}
	return &sheet, nil
}

func (s *appStore) UpdateSheet(ctx context.Context, sheet *models.Sheet) (*models.Sheet, error) {
	err := s.client.WithContext(ctx).Save(sheet).Error
	if err != nil {
		return nil, err
	}
	return sheet, nil
}
