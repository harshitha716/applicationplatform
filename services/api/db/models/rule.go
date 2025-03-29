package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rule struct {
	ID             uuid.UUID       `gorm:"column:rule_id"`
	OrganizationId uuid.UUID       `gorm:"column:organization_id"`
	DatasetId      uuid.UUID       `gorm:"column:dataset_id"`
	Column         string          `gorm:"column:column"`
	Value          string          `gorm:"column:value"`
	FilterConfig   json.RawMessage `gorm:"column:filter_config"`
	Title          string          `gorm:"column:title"`
	Description    string          `gorm:"column:description"`
	Priority       int             `gorm:"column:priority;unique"`
	CreatedAt      time.Time       `gorm:"column:created_at"`
	CreatedBy      uuid.UUID       `gorm:"column:created_by"`
	UpdatedAt      time.Time       `gorm:"column:updated_at"`
	UpdatedBy      uuid.UUID       `gorm:"column:updated_by"`
	DeletedAt      *time.Time      `gorm:"column:deleted_at"`
	DeletedBy      *uuid.UUID      `gorm:"column:deleted_by"`
}

type CreateRuleParams struct {
	Id             uuid.UUID
	Title          string
	Description    string
	OrganizationId uuid.UUID
	DatasetId      uuid.UUID
	Column         string
	Value          string
	FilterConfig   interface{}
	CreatedBy      uuid.UUID
}

type UpdateRuleParams struct {
	Title        string
	Description  string
	Value        string
	FilterConfig interface{}
	UpdatedBy    uuid.UUID
}

type FilterRuleParams struct {
	OrganizationId uuid.UUID
	DatasetColumns []DatasetColumn
}

type DatasetColumn struct {
	DatasetId uuid.UUID
	Columns   []string
}

type UpdateRulePriorityParams struct {
	DatasetId    uuid.UUID      `json:"dataset_id"`
	RulePriority []RulePriority `json:"rule_priority"`
	UpdatedBy    uuid.UUID      `json:"updated_by"`
}

type RulePriority struct {
	RuleId   uuid.UUID `json:"rule_id"`
	Priority int       `json:"priority"`
}

type DeleteRuleParams struct {
	RuleId    uuid.UUID
	DeletedBy uuid.UUID
}

func (Rule) TableName() string {
	return "rules"
}

func (r *Rule) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'dataset'
			AND frap.resource_id = rules.dataset_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (r *Rule) BeforeCreate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", "dataset", r.DatasetId, userId, "admin").Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("dataset access forbidden")
	}

	return nil

}

func (r *Rule) BeforeUpdate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", "dataset", r.DatasetId, userId, "admin").Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("dataset access forbidden")
	}

	return nil
}

func (r *Rule) BeforeDelete(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", "dataset", r.DatasetId, userId, "admin").Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("dataset access forbidden")
	}

	return nil

}
