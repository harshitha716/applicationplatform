package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Widget struct {
	Type           string          `json:"type" gorm:"primaryKey;column:type"`
	Name           string          `json:"name"`
	TemplateSchema json.RawMessage `json:"template_schema"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty"`
}

func (Widget) TableName() string {
	return "widgets"
}

func (w *Widget) GetQueryFilters(db *gorm.DB, currentUserId uuid.UUID, organizationIds []uuid.UUID) *gorm.DB {
	return db
}

func (w *Widget) BeforeCreate(db *gorm.DB) error {

	return fmt.Errorf("insert forbidden")

}
