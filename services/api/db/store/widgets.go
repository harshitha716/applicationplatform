package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type WidgetStore interface {
	GetWidgetInstanceByID(ctx context.Context, widgetInstanceID uuid.UUID) (models.WidgetInstance, error)
	GetWidgetTemplate(ctx context.Context, widgetType string) (models.Widget, error)
	CreateWidgetInstance(ctx context.Context, widgetInstance *models.WidgetInstance) (*models.WidgetInstance, error)
	UpdateWidgetInstance(ctx context.Context, widgetInstance *models.WidgetInstance) (*models.WidgetInstance, error)
}

func (s *appStore) GetWidgetInstanceByID(ctx context.Context, widgetInstanceID uuid.UUID) (models.WidgetInstance, error) {
	db := s.client.WithContext(ctx)

	widgetInstance := models.WidgetInstance{}
	db = db.Where("widget_instance_id = ?", widgetInstanceID)

	return widgetInstance, db.First(&widgetInstance).Error
}

func (s *appStore) GetWidgetTemplate(ctx context.Context, widgetType string) (models.Widget, error) {
	db := s.client.WithContext(ctx)

	widget := models.Widget{}
	db = db.Where("type = ?", widgetType)

	return widget, db.First(&widget).Error
}

func (s *appStore) CreateWidgetInstance(ctx context.Context, widgetInstance *models.WidgetInstance) (*models.WidgetInstance, error) {
	db := s.client.WithContext(ctx)
	widgetInstance.ID = uuid.New()

	return widgetInstance, db.Create(&widgetInstance).Error
}

func (s *appStore) UpdateWidgetInstance(ctx context.Context, widgetInstance *models.WidgetInstance) (*models.WidgetInstance, error) {
	db := s.client.WithContext(ctx)

	return widgetInstance, db.Save(widgetInstance).Error
}
