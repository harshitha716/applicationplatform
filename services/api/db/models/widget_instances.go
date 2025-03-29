package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WidgetInstance struct {
	ID            uuid.UUID        `json:"widget_instance_id" gorm:"column:widget_instance_id"`
	WidgetType    string           `json:"widget_type" gorm:"foreignKey:WidgetType"`
	SheetID       uuid.UUID        `json:"sheet_id"`
	Title         string           `json:"title"`
	DataMappings  json.RawMessage  `json:"data_mappings"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     *time.Time       `json:"deleted_at,omitempty"`
	DisplayConfig *json.RawMessage `json:"display_config,omitempty"`
}

func (WidgetInstance) TableName() string {
	return "widget_instances"
}

func (WidgetInstance) GetQueryFilters(db *gorm.DB, userId uuid.UUID, organizationIDs []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap, "app".sheets
			WHERE frap.resource_type = 'page'
			AND frap.resource_id = sheets.page_id
			AND sheets.sheet_id = widget_instances.sheet_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId.String(),
	)

}

func (w *WidgetInstance) BeforeCreate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	sheets := []Sheet{}
	err := db.Where(
		`
		EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap, "app".sheets
			WHERE frap.resource_type = 'page'
			AND frap.resource_id = sheets.page_id
			AND frap.privilege = 'admin'
			AND sheets.sheet_id = ?
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)
	`, w.SheetID, userId).Limit(1).Find(&sheets).Error

	if err != nil {
		return err
	}

	if len(sheets) == 0 {
		return fmt.Errorf("page access forbidden")
	}

	return nil

}

type CreateWidgetInstanceParams struct {
	WidgetID     uuid.UUID
	SheetID      uuid.UUID
	Title        string
	DataMappings json.RawMessage
}
