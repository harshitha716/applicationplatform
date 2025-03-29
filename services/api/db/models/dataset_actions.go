package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatasetAction struct {
	ID             uuid.UUID       `json:"id" gorm:"column:id"`
	ActionId       string          `json:"action_id" gorm:"column:action_id"`
	ActionType     string          `json:"action_type" gorm:"column:action_type"`
	DatasetId      uuid.UUID       `json:"dataset_id" gorm:"column:dataset_id"`
	OrganizationId uuid.UUID       `json:"organization_id" gorm:"column:organization_id"`
	Status         string          `json:"status" gorm:"column:status"`
	Config         json.RawMessage `json:"config" gorm:"column:config"`
	ActionBy       uuid.UUID       `json:"action_by" gorm:"column:action_by"`
	StartedAt      time.Time       `json:"started_at" gorm:"column:started_at"`
	CompletedAt    *time.Time      `json:"completed_at" gorm:"column:completed_at"`
}

func (DatasetAction) TableName() string {
	return "dataset_actions"
}

type CreateDatasetActionParams struct {
	ActionId    string
	ActionType  string
	DatasetId   uuid.UUID
	Status      string
	Config      interface{}
	ActionBy    uuid.UUID
	IsCompleted bool
}

type DatasetActionFilters struct {
	DatasetIds []uuid.UUID
	ActionIds  []string
	ActionType []string
	ActionBy   []uuid.UUID
	Status     []string
}

func (d *DatasetAction) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'dataset'
			AND frap.resource_id = dataset_actions.dataset_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (d *DatasetAction) BeforeCreate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", ResourceTypeDataset, d.DatasetId, userId, PrivilegeDatasetAdmin).Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("dataset write access forbidden")
	}

	return nil

}

func (d *DatasetAction) BeforeUpdate(db *gorm.DB) error {

	// 	return d.BeforeCreate(db)
	return nil
}
