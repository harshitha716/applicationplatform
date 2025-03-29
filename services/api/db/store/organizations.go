package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrganizationStore interface {
	OrganizationReadStore
	OrganizationWriteStore
}

type OrganizationReadStore interface {
	GetOrganizationById(ctx context.Context, organizationId string) (*models.Organization, error)
	GetOrganizationsAll(ctx context.Context) ([]models.Organization, error)
	GetOrganizationsByMemberId(ctx context.Context, memberId uuid.UUID) ([]models.Organization, error)
	GetOrganizationInvitationsAll(ctx context.Context) ([]models.OrganizationInvitation, error)
	GetOrganizationInvitationsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationInvitation, error)
	GetOrganizationInvitationById(ctx context.Context, invitationId uuid.UUID) (*models.OrganizationInvitation, error)
	GetOrganizationSSOConfigsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationSSOConfig, error)
	GetSSOConfigByDomain(ctx context.Context, domain string) (*models.OrganizationSSOConfig, error)
	GetPrimarySSOConfigByDomain(ctx context.Context, domain string) (*models.OrganizationSSOConfig, error)
	GetOrganizationMembershipRequestsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationMembershipRequest, error)
	GetOrganizationMembershipRequestsAll(ctx context.Context) ([]models.OrganizationMembershipRequest, error)
	GetOrganizationInvitationsAndMembershipRequests(ctx context.Context, organizationId uuid.UUID) (*models.Organization, error)
	organizationPoliciesReadStore
}

type OrganizationWriteStore interface {
	CreateOrganizationInvitation(ctx context.Context, organizationId uuid.UUID, targetEmail string, privilege models.ResourcePrivilege) (*models.OrganizationInvitation, error)
	CreateSSOConfig(ctx context.Context, organizationId uuid.UUID, ssoProviderID string, ssoProviderName string, ssoConfig json.RawMessage, emailDomain string) (*models.OrganizationSSOConfig, error)
	CreateOrganizationInvitationStatus(ctx context.Context, invitationId uuid.UUID, status models.InvitationStatus) (*models.OrganizationInvitationStatus, error)
	UpdateOrganizationInvitationStatus(ctx context.Context, invitationId uuid.UUID, status models.InvitationStatus) (*models.OrganizationInvitationStatus, error)
	CreateOrganizationMembershipRequest(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, status models.OrgMembershipStatus) (*models.OrganizationMembershipRequest, error)
	UpdatePendingOrganizationMembershipRequest(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, status models.OrgMembershipStatus) (*models.OrganizationMembershipRequest, error)
	CreateOrganization(ctx context.Context, name string, description *string, ownerId uuid.UUID) (*models.Organization, error)
	WithOrganizationTransaction(ctx context.Context, fn func(OrganizationStore) error) error
	organizationPoliciesWriteStore
}

func (s *appStore) GetOrganizationById(ctx context.Context, organizationId string) (*models.Organization, error) {

	var organization models.Organization
	err := s.client.WithContext(ctx).Model(&organization).Where("organization_id = ?", organizationId).First(&organization).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (s *appStore) GetOrganizationsAll(ctx context.Context) ([]models.Organization, error) {

	var organizations []models.Organization
	err := s.client.WithContext(ctx).Model(organizations).Order("created_at DESC").Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (s *appStore) GetOrganizationsByMemberId(ctx context.Context, memberId uuid.UUID) ([]models.Organization, error) {
	var organizations []models.Organization
	err := s.client.WithContext(ctx).Model(organizations).Preload("ResourceAudiencePolicies", "resource_audience_policies.resource_audience_id = ? AND resource_audience_policies.resource_audience_type = ? AND resource_audience_policies.resource_type = ?", memberId, models.AudienceTypeUser, models.ResourceTypeOrganization).Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (s *appStore) GetOrganizationInvitationsAll(ctx context.Context) ([]models.OrganizationInvitation, error) {
	var invitations []models.OrganizationInvitation
	err := s.client.WithContext(ctx).Model(invitations).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (s *appStore) GetOrganizationInvitationsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationInvitation, error) {
	var invitations []models.OrganizationInvitation
	err := s.client.WithContext(ctx).Model(invitations).Where("organization_id = ?", organizationId).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (s *appStore) GetOrganizationInvitationById(ctx context.Context, invitationId uuid.UUID) (*models.OrganizationInvitation, error) {
	var invitation models.OrganizationInvitation
	err := s.client.WithContext(ctx).Model(&invitation).Where("organization_invitation_id = ?", invitationId).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (s *appStore) CreateOrganizationInvitation(ctx context.Context, organizationId uuid.UUID, targetEmail string, privilege models.ResourcePrivilege) (*models.OrganizationInvitation, error) {

	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return nil, fmt.Errorf("user not found in context")
	}

	invitation := &models.OrganizationInvitation{
		OrganizationID: organizationId,
		TargetEmail:    targetEmail,
		InvitedBy:      *userId,
		Privilege:      privilege,
	}

	err := s.client.WithContext(ctx).Create(invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

func (s *appStore) CreateOrganizationInvitationStatus(ctx context.Context, invitationId uuid.UUID, status models.InvitationStatus) (*models.OrganizationInvitationStatus, error) {
	invitationStatus := &models.OrganizationInvitationStatus{
		OrganizationInvitationID: invitationId,
		Status:                   status,
	}

	err := s.client.WithContext(ctx).Create(invitationStatus).Error
	if err != nil {
		return nil, err
	}
	return invitationStatus, nil
}

func (s *appStore) UpdateOrganizationInvitationStatus(ctx context.Context, invitationId uuid.UUID, status models.InvitationStatus) (*models.OrganizationInvitationStatus, error) {
	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return nil, fmt.Errorf("user not found in context")
	}

	invitationStatus := &models.OrganizationInvitationStatus{
		OrganizationInvitationID: invitationId,
		Status:                   status,
	}

	err := s.client.WithContext(ctx).Create(&invitationStatus).Error
	if err != nil {
		return nil, err
	}

	return invitationStatus, nil
}

func (s *appStore) DeleteOrganizationInvitation(ctx context.Context, invitationId uuid.UUID) error {
	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return fmt.Errorf("user not found in context")
	}

	invitation := models.OrganizationInvitation{
		OrganizationInvitationID: invitationId,
	}

	err := s.client.WithContext(ctx).Delete(&invitation).Error
	if err != nil {
		return err
	}

	org := models.Organization{
		ResourceAudiencePolicies: []models.ResourceAudiencePolicy{
			{
				ResourceID:         invitation.OrganizationID,
				ResourceType:       models.ResourceTypeOrganization,
				ResourceAudienceID: invitation.InvitedBy,
				Privilege:          invitation.Privilege,
			},
		},
	}

	s.client.DB.Create(&org)

	return nil
}

func (s *appStore) GetOrganizationSSOConfigsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationSSOConfig, error) {
	var configs []models.OrganizationSSOConfig
	err := s.client.WithContext(ctx).Model(configs).Where("organization_id = ?", organizationId).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *appStore) GetSSOConfigByDomain(ctx context.Context, domain string) (*models.OrganizationSSOConfig, error) {
	var ssoConfig models.OrganizationSSOConfig
	err := s.client.WithContext(ctx).Model(&ssoConfig).Where("email_domain = ?", domain).Order("is_primary DESC").First(&ssoConfig).Error
	if err != nil {
		return nil, err
	}
	return &ssoConfig, nil
}

func (s *appStore) GetPrimarySSOConfigByDomain(ctx context.Context, domain string) (*models.OrganizationSSOConfig, error) {
	var ssoConfig models.OrganizationSSOConfig
	err := s.client.WithContext(ctx).Model(&ssoConfig).Where("email_domain = ? AND is_primary = ?", domain, true).First(&ssoConfig).Error
	if err != nil {
		return nil, err
	}
	return &ssoConfig, nil
}

func (s *appStore) CreateSSOConfig(ctx context.Context, organizationId uuid.UUID, ssoProviderID string, ssoProviderName string, ssoConfig json.RawMessage, emailDomain string) (*models.OrganizationSSOConfig, error) {
	config := &models.OrganizationSSOConfig{
		OrganizationID:  organizationId,
		SSOProviderID:   ssoProviderID,
		SSOProviderName: ssoProviderName,
		SSOConfig:       ssoConfig,
		EmailDomain:     emailDomain,
	}
	err := s.client.WithContext(ctx).Create(config).Error
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *appStore) WithOrganizationTransaction(ctx context.Context, fn func(OrganizationStore) error) error {
	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txClient := pgclient.PostgresClient{DB: tx}
		return fn(&appStore{client: &txClient})
	})
}

func (s *appStore) GetOrganizationMembershipRequestsAll(ctx context.Context) ([]models.OrganizationMembershipRequest, error) {
	var requests []models.OrganizationMembershipRequest
	err := s.client.WithContext(ctx).Model(requests).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil

}

func (s *appStore) GetOrganizationMembershipRequestsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationMembershipRequest, error) {
	var requests []models.OrganizationMembershipRequest
	err := s.client.WithContext(ctx).Model(requests).Where("organization_id = ? AND status = ?", organizationId, models.OrgMembershipStatusPending).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *appStore) CreateOrganizationMembershipRequest(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, status models.OrgMembershipStatus) (*models.OrganizationMembershipRequest, error) {
	request := &models.OrganizationMembershipRequest{
		OrganizationID: organizationId,
		UserID:         userId,
		Status:         status,
	}

	err := s.client.WithContext(ctx).Create(request).Error
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (s *appStore) UpdatePendingOrganizationMembershipRequest(ctx context.Context, orgId uuid.UUID, memberId uuid.UUID, status models.OrgMembershipStatus) (*models.OrganizationMembershipRequest, error) {

	request := &models.OrganizationMembershipRequest{
		OrganizationID: orgId,
		UserID:         memberId,
		Status:         status,
	}

	err := s.client.WithContext(ctx).Model(request).Where("organization_id = ? AND user_id = ? AND status = ?", orgId, memberId, models.OrgMembershipStatusPending).Clauses(clause.Returning{}).Update("status", status).Error
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (s *appStore) CreateOrganization(ctx context.Context, name string, description *string, ownerId uuid.UUID) (*models.Organization, error) {
	organization := &models.Organization{
		Name:        name,
		Description: description,
		OwnerId:     ownerId,
	}

	err := s.client.WithContext(ctx).Create(organization).Error
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (s *appStore) GetOrganizationInvitationsAndMembershipRequests(ctx context.Context, organizationId uuid.UUID) (*models.Organization, error) {

	organization := &models.Organization{}
	err := s.client.WithContext(ctx).
		Model(organization).
		Where("organization_id = ?", organizationId).
		Preload("Invitations.InvitationStatuses").
		Preload("MembershipRequests.User").
		First(organization).Error
	if err != nil {
		return nil, err
	}
	return organization, nil
}
