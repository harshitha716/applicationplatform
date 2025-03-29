package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Connection struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id"`
	ConnectorID    uuid.UUID      `json:"connector_id" gorm:"column:connector_id"`
	OrganizationID uuid.UUID      `json:"organization_id" gorm:"column:organization_id"`
	Name           string         `json:"name" gorm:"column:name"`
	Status         string         `json:"status" gorm:"column:status"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`

	Connector Connector  `json:"connector" gorm:"foreignKey:ConnectorID;references:ID"`
	Schedules []Schedule `json:"schedules" gorm:"foreignKey:ConnectionID; references:ID"`
}

type CreateConnectionParams struct {
	ConnectorID uuid.UUID
	Name        string
	Status      string
	Config      map[string]interface{}
}

func (o *Connection) TableName() string {
	return "connections"
}

func (o *Connection) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'connection'
			AND frap.resource_id = connections.id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (o *Connection) BeforeCreate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", "organization", o.OrganizationID, userId).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil
}
