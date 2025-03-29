package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Connector struct {
	ID             uuid.UUID       `json:"id" gorm:"column:id"`
	Name           string          `json:"name" gorm:"column:name"`
	Description    string          `json:"description" gorm:"column:description"`
	DisplayName    string          `json:"display_name" gorm:"column:display_name"`
	Documentation  string          `json:"documentation" gorm:"column:documentation"`
	LogoURL        string          `json:"logo_url" gorm:"column:logo_url"`
	ConfigTemplate json.RawMessage `json:"config_template" gorm:"column:config_template"`
	Category       string          `json:"category" gorm:"column:category"`
	CreatedAt      time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt  `json:"deleted_at" gorm:"column:deleted_at"`
	IsDeleted      bool            `json:"is_deleted" gorm:"column:is_deleted"`
	Status         string          `json:"status" gorm:"column:status"`
}

type ConnectorWithActiveConnectionsCount struct {
	Connector
	ActiveConnectionsCount int `json:"active_connections_count"`
}

func (c *Connector) TableName() string {
	return "connectors"
}

func (c *Connector) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db
}

func (c *Connector) BeforeCreate(db *gorm.DB) error {
	return nil
}
