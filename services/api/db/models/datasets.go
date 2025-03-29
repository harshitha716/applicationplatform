package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatasetType string

const (
	DatasetTypeBronze DatasetType = "bronze"
	DatasetTypeSource DatasetType = "source"
	DatasetTypeMV     DatasetType = "materialised"
	DatasetTypeStaged DatasetType = "staged"
)

var UserVisibleDatasetTypes = []DatasetType{
	DatasetTypeSource,
	DatasetTypeMV,
}

var ValidDatasetTypes = []DatasetType{
	DatasetTypeBronze,
	DatasetTypeSource,
	DatasetTypeMV,
	DatasetTypeStaged,
}

type Dataset struct {
	ID             uuid.UUID       `json:"dataset_id" gorm:"column:dataset_id"`
	Title          string          `json:"title"`
	Description    *string         `json:"description"`
	OrganizationId uuid.UUID       `json:"organization_id" gorm:"column:organization_id"`
	Type           DatasetType     `json:"type"`
	CreatedBy      uuid.UUID       `json:"created_by" gorm:"column:created_by"`
	Metadata       json.RawMessage `json:"metadata"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty"`
}

var DatasetListingSortingColumns = []string{
	"title",
	"updated_at",
	"description",
}

type DatasetFilters struct {
	OrganizationIds []uuid.UUID
	DatasetIds      []uuid.UUID
	CreatedBy       []uuid.UUID
	Type            []DatasetType

	// Pagination
	Page  int
	Limit int

	// Sorting
	SortParams []DatasetSortParam
}

type DatasetSortParam struct {
	Column string
	Desc   bool
}

func (d *Dataset) TableName() string {
	return "datasets"
}

func (d *Dataset) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'dataset'
			AND frap.resource_id = datasets.dataset_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (a *Dataset) BeforeCreate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", "organization", a.OrganizationId, userId).Limit(1).Find(&fraps).Error

	if err != nil {
		return err
	}

	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil

}
