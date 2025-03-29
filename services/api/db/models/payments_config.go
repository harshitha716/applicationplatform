package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentsConfigStatus string

const (
	PaymentsConfigStatusConnected PaymentsConfigStatus = "connected"
)

type PaymentsConfig struct {
	ID                uuid.UUID            `json:"id" gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	OrganizationID    uuid.UUID            `json:"organization_id" gorm:"column:organization_id"`
	AccountsDatasetID uuid.UUID            `json:"accounts_dataset_id" gorm:"column:accounts_dataset_id"`
	MappingConfig     json.RawMessage      `json:"mapping_config" gorm:"column:mapping_config"`
	Status            PaymentsConfigStatus `json:"status" gorm:"column:status"`
	DeletedAt         gorm.DeletedAt       `json:"deleted_at" gorm:"column:deleted_at"`
	CreatedAt         time.Time            `json:"created_at" gorm:"column:created_at;default:now()"`
	UpdatedAt         time.Time            `json:"updated_at" gorm:"column:updated_at;default:now()"`
}

func (p *PaymentsConfig) TableName() string {
	return "payments_config"
}

func (p *PaymentsConfig) BeforeCreate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", ResourceTypeOrganization, p.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil
}

func (p *PaymentsConfig) BeforeUpdate(tx *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(tx.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := tx.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", ResourceTypePayments, p.ID, userId, PrivilegePaymentsAdmin).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("payments access forbidden")
	}

	return nil
}

func (p *PaymentsConfig) BeforeDelete(tx *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(tx.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := tx.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", ResourceTypePayments, p.ID, userId, PrivilegePaymentsAdmin).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("payments access forbidden")
	}
	return nil
}

func (o *PaymentsConfig) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'payments'
			AND frap.resource_id = payments_config.id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}
