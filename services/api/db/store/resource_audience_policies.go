package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type organizationPoliciesReadStore interface {
	GetOrganizationPolicies(ctx context.Context, orgId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
	GetOrganizationPolicyByUser(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.ResourceAudiencePolicy, error)
	GetOrganizationPoliciesByEmail(ctx context.Context, organizationId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error)
}

type organizationPoliciesWriteStore interface {
	CreateOrganizationPolicy(ctx context.Context, orgId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	UpdateOrganizationPolicy(ctx context.Context, orgId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	DeleteOrganizationPolicy(ctx context.Context, orgId uuid.UUID, audienceId uuid.UUID) error
}

type datasetPoliciesStore interface {
	datasetPoliciesReadStore
	datasetPoliciesWriteStore
}

type datasetPoliciesReadStore interface {
	GetDatasetPolicies(ctx context.Context, datasetId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
	GetDatasetPoliciesByEmail(ctx context.Context, datasetId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error)
}

type datasetPoliciesWriteStore interface {
	CreateDatasetPolicy(ctx context.Context, datasetId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	UpdateDatasetPolicy(ctx context.Context, datasetId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	DeleteDatasetPolicy(ctx context.Context, datasetId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error
}

type pagePoliciesStore interface {
	pagePoliciesReadStore
	pagePoliciesWriteStore
}

type pagePoliciesReadStore interface {
	GetPagesPolicies(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
	GetPagePoliciesByEmail(ctx context.Context, pageId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error)
}

type pagePoliciesWriteStore interface {
	CreatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	UpdatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	DeletePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error
}

type ConnectionPoliciesStore interface {
	connectionPoliciesReadStore
	connectionPoliciesWriteStore
}

type connectionPoliciesReadStore interface {
	GetConnectionPolicies(ctx context.Context, connectionId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
}

type connectionPoliciesWriteStore interface {
	CreateConnectionPolicy(ctx context.Context, connectionId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	UpdateConnectionPolicy(ctx context.Context, connectionId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
}

func (s *appStore) getResourcePolicies(ctx context.Context, resourceId uuid.UUID, resourceType models.ResourceType) ([]models.ResourceAudiencePolicy, error) {
	var policies []models.ResourceAudiencePolicy
	err := s.client.WithContext(ctx).Preload("UserPolicies").Model(policies).Preload("User").Where("resource_id = ? AND resource_type = ?", resourceId, resourceType).Order("created_at DESC").Find(&policies).Error
	if err != nil {
		return nil, err
	}
	return policies, nil
}

func (s *appStore) getResourcePoliciesByAudienceEmail(ctx context.Context, resourceId uuid.UUID, resourceType models.ResourceType, audienceEmail string) ([]models.ResourceAudiencePolicy, error) {
	var policies []models.ResourceAudiencePolicy
	err := s.client.WithContext(ctx).Preload("UserPolicies").Model(policies).Joins("JOIN users_with_traits ON resource_audience_id = users_with_traits.user_id AND users_with_traits.email ILIKE ?", audienceEmail).Where("resource_id = ? AND resource_type = ? AND resource_audience_type = ?", resourceId, resourceType, models.AudienceTypeUser).Find(&policies).Error
	if err != nil {
		return nil, err
	}
	return policies, nil
}

func (s *appStore) GetOrganizationPolicies(ctx context.Context, organizationId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePolicies(ctx, organizationId, models.ResourceTypeOrganization)
}

func (s *appStore) GetOrganizationPolicyByUser(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.ResourceAudiencePolicy, error) {
	var policy models.ResourceAudiencePolicy
	err := s.client.WithContext(ctx).Model(policy).Preload("User").Where("resource_id = ? AND resource_type = ? AND resource_audience_id = ?", organizationId, models.ResourceTypeOrganization, userId).First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *appStore) GetDatasetPolicies(ctx context.Context, datasetId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePolicies(ctx, datasetId, models.ResourceTypeDataset)
}

func (s *appStore) GetPagesPolicies(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePolicies(ctx, pageId, models.ResourceTypePage)
}

func (s *appStore) GetConnectionPolicies(ctx context.Context, connectionId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePolicies(ctx, connectionId, models.ResourceTypeConnection)
}

func (s *appStore) GetOrganizationPoliciesByEmail(ctx context.Context, organizationId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePoliciesByAudienceEmail(ctx, organizationId, models.ResourceTypeOrganization, email)
}

func (s *appStore) GetDatasetPoliciesByEmail(ctx context.Context, datasetId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePoliciesByAudienceEmail(ctx, datasetId, models.ResourceTypeDataset, email)
}

func (s *appStore) GetPagePoliciesByEmail(ctx context.Context, pageId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error) {
	return s.getResourcePoliciesByAudienceEmail(ctx, pageId, models.ResourceTypePage, email)
}

func (s *appStore) CreateOrganizationPolicy(ctx context.Context, orgId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.createResourcePolicy(ctx, models.ResourceTypeOrganization, orgId, audienceType, audienceId, privilege)
}

func (s *appStore) CreateDatasetPolicy(ctx context.Context, datasetId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.createResourcePolicy(ctx, models.ResourceTypeDataset, datasetId, audienceType, audienceId, privilege)
}

func (s *appStore) CreatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.createResourcePolicy(ctx, models.ResourceTypePage, pageId, audienceType, audienceId, privilege)
}

func (s *appStore) CreateConnectionPolicy(ctx context.Context, connectionId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.createResourcePolicy(ctx, models.ResourceTypeConnection, connectionId, audienceType, audienceId, privilege)
}

func (s *appStore) createResourcePolicy(ctx context.Context, resourceType models.ResourceType, resourceId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	policy := models.ResourceAudiencePolicy{
		ResourceAudienceType: audienceType,
		ResourceAudienceID:   audienceId,
		ResourceID:           resourceId,
		ResourceType:         resourceType,
		Privilege:            privilege,
	}
	err := s.client.WithContext(ctx).Create(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *appStore) UpdateOrganizationPolicy(ctx context.Context, orgId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.updateResourcePolicy(ctx, models.ResourceTypeOrganization, orgId, audienceId, privilege)
}

func (s *appStore) UpdateDatasetPolicy(ctx context.Context, datasetId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.updateResourcePolicy(ctx, models.ResourceTypeDataset, datasetId, audienceId, privilege)
}

func (s *appStore) UpdatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.updateResourcePolicy(ctx, models.ResourceTypePage, pageId, audienceId, privilege)
}

func (s *appStore) UpdateConnectionPolicy(ctx context.Context, connectionId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	return s.updateResourcePolicy(ctx, models.ResourceTypeConnection, connectionId, audienceId, privilege)
}

func (s *appStore) updateResourcePolicy(ctx context.Context, resourceType models.ResourceType, resourceId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	policy := models.ResourceAudiencePolicy{
		ResourceType:       resourceType,
		ResourceID:         resourceId,
		ResourceAudienceID: audienceId,
	}
	err := s.client.WithContext(ctx).
		Model(&policy).
		Where("resource_type = ? AND resource_id = ? AND resource_audience_id = ?", resourceType, resourceId, audienceId).Clauses(clause.Returning{}).
		Update("privilege", privilege).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *appStore) DeleteOrganizationPolicy(ctx context.Context, orgId uuid.UUID, audienceId uuid.UUID) error {
	return s.deleteResourcePolicy(ctx, models.ResourceTypeOrganization, orgId, models.AudienceTypeUser, audienceId)
}

func (s *appStore) DeleteDatasetPolicy(ctx context.Context, datasetId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {
	return s.deleteResourcePolicy(ctx, models.ResourceTypeDataset, datasetId, audienceType, audienceId)
}

func (s *appStore) DeletePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {
	return s.deleteResourcePolicy(ctx, models.ResourceTypePage, pageId, audienceType, audienceId)
}

func (s *appStore) deleteResourcePolicy(ctx context.Context, resourceType models.ResourceType, resourceId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {

	policy := models.ResourceAudiencePolicy{
		ResourceType:         resourceType,
		ResourceID:           resourceId,
		ResourceAudienceType: audienceType,
		ResourceAudienceID:   audienceId,
	}

	return s.client.WithContext(ctx).Where("resource_type = ? AND resource_id = ? AND resource_audience_type = ? AND resource_audience_id = ?", resourceType, resourceId, audienceType, audienceId).Delete(&policy).Error
}
