package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sheet struct {
	ID              uuid.UUID        `json:"sheet_id" gorm:"column:sheet_id"`
	Name            string           `json:"name"`
	Description     *string          `json:"description"`
	WidgetInstances []WidgetInstance `json:"widget_instances,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at,omitempty"`
	FractionalIndex float64          `json:"fractional_index"`
	PageId          uuid.UUID        `json:"page_id"`
	SheetConfig     json.RawMessage  `json:"sheet_config"`
}

type SheetFilters struct {
	PageIds  []uuid.UUID
	SheetIds []uuid.UUID

	IncludeWidgetInstances bool

	// Pagination
	Page  int
	Limit int

	// Sorting
	SortParams []SheetSortParams
}

type SheetSortParams struct {
	Column string
	Desc   bool
}

func (s *Sheet) TableName() string {
	return "sheets"
}

func (s *Sheet) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'page'
			AND frap.resource_id = sheets.page_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (s *Sheet) BeforeCreate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", "page", s.PageId, userId, "admin").Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("page access forbidden")
	}

	return nil

}
