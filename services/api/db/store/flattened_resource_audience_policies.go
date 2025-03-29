package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
)

type FlattenedResourceAudiencePoliciesStore interface {
	GetFlattenedResourceAudiencePolicies(ctx context.Context, filters models.FlattenedResourceAudiencePoliciesFilters) ([]models.FlattenedResourceAudiencePolicy, error)
}

func (store *appStore) GetFlattenedResourceAudiencePolicies(ctx context.Context, filters models.FlattenedResourceAudiencePoliciesFilters) ([]models.FlattenedResourceAudiencePolicy, error) {
	var flattenedResourceAudiencePolicies []models.FlattenedResourceAudiencePolicy
	query := store.client.WithContext(ctx).Model(&flattenedResourceAudiencePolicies)
	if len(filters.ResourceIds) > 0 {
		query = query.Where("resource_id IN ?", filters.ResourceIds)
	}
	if len(filters.UserIds) > 0 {
		query = query.Where("user_id IN ?", filters.UserIds)
	}
	if len(filters.ResourceTypes) > 0 {
		query = query.Where("resource_type IN ?", filters.ResourceTypes)
	}
	if len(filters.Privileges) > 0 {
		query = query.Where("privilege IN ?", filters.Privileges)
	}
	query = query.Where("deleted_at IS NULL")

	err := query.Find(&flattenedResourceAudiencePolicies).Error
	if err != nil {
		return nil, err
	}
	return flattenedResourceAudiencePolicies, nil
}
