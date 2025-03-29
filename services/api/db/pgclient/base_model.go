package pgclient

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel interface {
	GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB
	// GetUpdateFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB
	// GetDeleteFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB
	BeforeCreate(db *gorm.DB) error
}
