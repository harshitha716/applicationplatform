package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Page struct {
	ID              uuid.UUID  `json:"page_id" gorm:"column:page_id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	Sheets          []Sheet    `json:"sheets,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	FractionalIndex float64    `json:"fractional_index"`
	OrganizationId  uuid.UUID  `json:"organization_id"`
}

type PageFilters struct {
	OrganizationIds []uuid.UUID
	PageIds         []uuid.UUID
	IncludeSheets   bool

	// Pagination
	Page  int
	Limit int

	// Sorting
	SortParams []PageSortParams
}

type PageSortParams struct {
	Column string
	Desc   bool
}

func (o *Page) TableName() string {
	return "pages"
}

func (o *Page) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'page'
			AND frap.resource_id = pages.page_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (p *Page) BeforeCreate(db *gorm.DB) error {

	// ensure if it is an authenticated user
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	// check in flattened_resource_audience_policies if user has access to organization
	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", "organization", p.OrganizationId, userId).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	// if user has no access to organization
	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil

}
