package store

import (
	"context"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageStore interface {
	GetPageById(ctx context.Context, pageId uuid.UUID) (*models.Page, error)
	GetPagesAll(ctx context.Context, filters models.PageFilters) ([]models.Page, error)
	GetPagesByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.Page, error)
	CreatePage(ctx context.Context, name string, description string) (*models.Page, error)
	pagePoliciesStore
	WithPageTransaction(ctx context.Context, fn func(PageStore) error) error
}

func (s *appStore) GetPageById(ctx context.Context, pageId uuid.UUID) (*models.Page, error) {
	var page models.Page
	page.ID = pageId
	err := s.client.WithContext(ctx).Preload("Sheets").Model(page).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (s *appStore) GetPagesByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.Page, error) {
	var pages []models.Page
	err := s.client.WithContext(ctx).Preload("Sheets").Preload("Sheets.WidgetInstances").Where("organization_id = ?", organizationId).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (s *appStore) GetPagesAll(ctx context.Context, filters models.PageFilters) ([]models.Page, error) {
	db := s.client.WithContext(ctx)

	pages := []models.Page{}
	if len(filters.OrganizationIds) > 0 {
		db = db.Where("organization_id IN (?)", filters.OrganizationIds)
	}

	if len(filters.PageIds) > 0 {
		db = db.Where("page_id IN (?)", filters.PageIds)
	}

	if filters.IncludeSheets {
		db = db.Preload("Sheets")
	}

	if filters.Page > 0 {
		db = db.Offset((filters.Page - 1) * filters.Limit)
	}

	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}

	if len(filters.SortParams) > 0 {
		for _, sort := range filters.SortParams {
			db = db.Order(clause.OrderByColumn{
				Column: clause.Column{Name: sort.Column},
				Desc:   sort.Desc,
			})
		}
	}

	result := db.Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}

	return pages, nil

}

func (a *appStore) CreatePage(ctx context.Context, name string, description string) (*models.Page, error) {

	// inject the organization ID from context
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)
	if len(orgIds) != 1 {
		return nil, fmt.Errorf("organization access forbidden")
	}

	page := models.Page{
		ID:             uuid.New(),
		Name:           name,
		Description:    &description,
		OrganizationId: orgIds[0],
	}

	err := a.client.WithContext(ctx).Model(page).Create(&page).Error
	if err != nil {
		return nil, err
	}

	return &page, nil

}

func (s *appStore) WithPageTransaction(ctx context.Context, fn func(PageStore) error) error {
	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txClient := pgclient.PostgresClient{DB: tx}
		return fn(&appStore{client: &txClient})
	})
}
