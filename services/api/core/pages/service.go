package pages

import (
	"context"
	"fmt"
	"slices"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PagesServiceStore interface {
	store.PageStore
	GetPagesPolicies(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
}

type PagesService interface {
	GetPagesAll(ctx context.Context) ([]models.Page, error)
	GetPageByID(ctx context.Context, pageId uuid.UUID) (*models.Page, error)
	GetPagesByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.Page, error)
	GetPageAudiences(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
	AddAudienceToPage(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	BulkAddAudienceToPage(ctx context.Context, pageId uuid.UUID, payload BulkAddPageAudiencePayload) ([]*models.ResourceAudiencePolicy, BulkAddPageAudienceErrors)
	RemoveAudienceFromPage(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID) error
	UpdatePageAudiencePrivilege(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
}

type pagesService struct {
	store PagesServiceStore
}

func NewPagesService(appStore store.Store) *pagesService {
	return &pagesService{store: appStore}
}

func (s *pagesService) getPagePolicies(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	policies, err := s.store.GetPagesPolicies(ctx, pageId)
	if err != nil {
		return nil, fmt.Errorf("failed to get page policies: %w", err)
	}

	return policies, nil
}

func (s *pagesService) GetPagesAll(ctx context.Context) ([]models.Page, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	pages, err := s.store.GetPagesAll(ctx, models.PageFilters{SortParams: []models.PageSortParams{{Column: "created_at", Desc: false}}})
	if err != nil {
		ctxLogger.Error("failed to get pages", zap.Error(err))
		return nil, err
	}

	return pages, nil
}

func (s *pagesService) GetPageByID(ctx context.Context, pageId uuid.UUID) (*models.Page, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	page, err := s.store.GetPageById(ctx, pageId)
	if err != nil {
		ctxLogger.Error("failed to get page", zap.Error(err))
		return nil, err
	}

	return page, nil
}

func (s *pagesService) GetPageAudiences(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	policies, err := s.store.GetPagesPolicies(ctx, pageId)
	if err != nil {
		ctxLogger.Error("failed to get page policies", zap.Error(err))
		return nil, err
	}

	return policies, nil
}

func (s *pagesService) AddAudienceToPage(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)
	if !slices.Contains(models.PagePrivileges, privilege) {
		ctxlogger.Error("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	switch audienceType {
	case models.AudienceTypeUser:
		return s.addUserAudienceToPage(ctx, pageId, audienceId, privilege)
	case models.AudienceTypeOrganization:
		return s.addOrganizationAudienceToPage(ctx, pageId, audienceId, privilege)
	case models.AudienceTypeTeam:
		return s.addTeamAudienceToPage(ctx, pageId, audienceId, privilege)
	default:
		ctxlogger.Error("only user and organization audience is supported")
		return nil, fmt.Errorf("only user audience is supported")
	}
}

func (s *pagesService) addUserAudienceToPage(ctx context.Context, pageId uuid.UUID, userId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in the context")
	}

	var createdPolicy *models.ResourceAudiencePolicy

	err := s.store.WithPageTransaction(ctx, func(ps store.PageStore) error {
		// current user should be an admin on the page
		policies, err := s.getPagePolicies(ctx, pageId)
		if err != nil {
			ctxlogger.Info("failed to get page policies", zap.String("error", err.Error()))
			return err
		}

		// ensure that user is not already added on the page
		err = ensureAudienceNotAlreadyAdded(models.AudienceTypeUser, userId, policies)
		if err != nil {
			ctxlogger.Info("user already exists on the page")
			return err
		}

		// ensure that the current user is an admin on the page
		err = ensureCurrentUsersAdminAccess(ctx, policies)
		if err != nil {
			ctxlogger.Info("current user does not have access to change permissions on the page")
			return err
		}

		createdPolicy, err = ps.CreatePagePolicy(ctx, pageId, models.AudienceTypeUser, userId, privilege)
		if err != nil {
			ctxlogger.Info("failed to create page policy", zap.String("error", err.Error()))
			return err
		}

		ctxlogger.Info("created page policy", zap.Any("policy", createdPolicy))

		return nil
	})

	return createdPolicy, err

}

func (s *pagesService) addOrganizationAudienceToPage(ctx context.Context, pageId uuid.UUID, organizationId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)
	_, currentUserId, orgIds := apicontext.GetAuthFromContext(ctx)

	if currentUserId == nil {
		ctxlogger.Info("no user id found in context")
		return nil, fmt.Errorf("no user id found in context")
	}

	if !slices.Contains(orgIds, organizationId) {
		ctxlogger.Info("current user does not have access to add organizations on the page")
		return nil, fmt.Errorf("current user does not have access to add organizations on the page")
	}

	// ensure valid privilege
	if !slices.Contains(models.PagePrivileges, privilege) {
		ctxlogger.Info("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	var createdPolicy *models.ResourceAudiencePolicy

	err := s.store.WithPageTransaction(ctx, func(ds store.PageStore) error {
		// check if the page belongs to the organization
		page, err := ds.GetPageById(ctx, pageId)
		if err != nil {
			ctxlogger.Info("failed to get page", zap.String("error", err.Error()))
			return err
		}

		if page.OrganizationId != organizationId {
			ctxlogger.Info("page does not belong to the organization")
			return fmt.Errorf("page does not belong to the organization")
		}

		// Get existing policies
		policies, err := s.getPagePolicies(ctx, pageId)
		if err != nil {
			ctxlogger.Info("failed to get page policies", zap.String("error", err.Error()))
			return err
		}

		// ensure that user is not already added on the page
		err = ensureAudienceNotAlreadyAdded(models.AudienceTypeOrganization, organizationId, policies)
		if err != nil {
			ctxlogger.Info("organization already exists on the page")
			return err
		}

		// ensure that the current user is an admin on the page
		err = ensureCurrentUsersAdminAccess(ctx, policies)
		if err != nil {
			ctxlogger.Info("current user does not have access to change permissions on the page")
			return err
		}

		// Create new policy
		createdPolicy, err = ds.CreatePagePolicy(ctx, pageId, models.AudienceTypeOrganization, organizationId, privilege)
		if err != nil {
			ctxlogger.Info("failed to create page policy", zap.String("error", err.Error()))
			return err
		}

		ctxlogger.Info("created page policy", zap.Any("policy", createdPolicy))

		return nil
	})

	return createdPolicy, err
}

func (s *pagesService) addTeamAudienceToPage(ctx context.Context, pageId uuid.UUID, teamId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure valid privilege
	if !slices.Contains(models.PagePrivileges, privilege) {
		ctxlogger.Info("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	var createdPolicy *models.ResourceAudiencePolicy
	err := s.store.WithPageTransaction(ctx, func(ps store.PageStore) error {
		policies, err := s.getPagePolicies(ctx, pageId)
		if err != nil {
			ctxlogger.Info("failed to get page policies", zap.String("error", err.Error()))
			return err
		}

		// ensure that the team is not already added on the page
		err = ensureAudienceNotAlreadyAdded(models.AudienceTypeTeam, teamId, policies)
		if err != nil {
			ctxlogger.Info("team already exists on the page")
			return err
		}

		// ensure that the current user is an admin on the page
		err = ensureCurrentUsersAdminAccess(ctx, policies)
		if err != nil {
			ctxlogger.Info("current user does not have access to change permissions on the page")
			return err
		}

		p, err := ps.CreatePagePolicy(ctx, pageId, models.AudienceTypeTeam, teamId, privilege)
		if err != nil {
			ctxlogger.Info("failed to create page policy", zap.String("error", err.Error()))
			return err
		}

		createdPolicy = p

		ctxlogger.Info("created page policy", zap.Any("policy", createdPolicy))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdPolicy, nil
}

// Helper functions moved to service_helper.go

func (s *pagesService) BulkAddAudienceToPage(ctx context.Context, pageId uuid.UUID, payload BulkAddPageAudiencePayload) ([]*models.ResourceAudiencePolicy, BulkAddPageAudienceErrors) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)
	var createdPolicies []*models.ResourceAudiencePolicy
	var bulkErrors BulkAddPageAudienceErrors

	for _, audience := range payload.Audiences {
		policy, err := s.AddAudienceToPage(ctx, pageId, audience.AudienceType, audience.AudienceId, audience.Privilege)
		if err != nil {
			ctxlogger.Info("failed to add audience to page",
				zap.String("error", err.Error()),
				zap.String("audience_id", audience.AudienceId.String()))

			bulkErrors.Audiences = append(bulkErrors.Audiences, AddPageAudienceError{
				AudienceId:   audience.AudienceId,
				ErrorMessage: err.Error(),
			})
			continue
		}
		createdPolicies = append(createdPolicies, policy)
	}

	return createdPolicies, bulkErrors

}

func (s *pagesService) RemoveAudienceFromPage(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID) error {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return fmt.Errorf("no user ID found in the context")
	}

	// current user should be an admin on the page
	policies, err := s.getPagePolicies(ctx, pageId)
	if err != nil {
		ctxlogger.Info("failed to get page policies", zap.String("error", err.Error()))
		return err
	}

	// ensure that the current user is an admin on the page
	err = ensureCurrentUsersAdminAccess(ctx, policies)
	if err != nil {
		ctxlogger.Info("current user does not have access to change permissions on the page")
		return err
	}

	var policyToBeRemoved *models.ResourceAudiencePolicy
	for _, policy := range policies {
		if policy.ResourceAudienceID == audienceId {
			policyToBeRemoved = &policy
			break
		}
	}

	if policyToBeRemoved == nil {
		ctxlogger.Info("policy not found", zap.String("audience_id", audienceId.String()))
		return fmt.Errorf("invalid audience id")
	}

	err = ensureUserIsNotChangingTheirOwnAdminPolicy(*policyToBeRemoved, policies, *currentUserId)
	if err != nil {
		ctxlogger.Info("current user is changing their own admin policy", zap.String("error", err.Error()))
		return err
	}

	err = s.store.DeletePagePolicy(ctx, pageId, policyToBeRemoved.ResourceAudienceType, audienceId)
	if err != nil {
		ctxlogger.Error("failed to delete page policy", zap.Error(err))
		return err
	}

	return nil
}

func (s *pagesService) UpdatePageAudiencePrivilege(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// ensure authenticated
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in the context")
	}

	// ensure valid privilege
	if !slices.Contains(models.PagePrivileges, privilege) {
		ctxlogger.Info("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	// ensure that the current user is an admin on the page
	policies, err := s.getPagePolicies(ctx, pageId)
	if err != nil {
		ctxlogger.Info("failed to get page policies", zap.String("error", err.Error()))
		return nil, err
	}

	// ensure that the current user is an admin on the page
	err = ensureCurrentUsersAdminAccess(ctx, policies)
	if err != nil {
		ctxlogger.Info("current user does not have access to change permissions on the page")
		return nil, err
	}

	var policyToBeUpdated *models.ResourceAudiencePolicy
	for _, policy := range policies {
		if policy.ResourceAudienceID == audienceId {
			policyToBeUpdated = &policy
			break
		}
	}
	if policyToBeUpdated == nil {
		ctxlogger.Info("policy not found", zap.String("audience_id", audienceId.String()))
		return nil, fmt.Errorf("invalid audience id")
	}

	err = ensureUserIsNotChangingTheirOwnAdminPolicy(*policyToBeUpdated, policies, *currentUserId)
	if err != nil {
		ctxlogger.Info("current user is changing their own admin policy", zap.String("error", err.Error()))
		return nil, err
	}

	updatedPolicy, err := s.store.UpdatePagePolicy(ctx, pageId, audienceId, privilege)
	if err != nil {
		ctxlogger.Error("failed to update page policy", zap.Error(err))
		return nil, err
	}

	return updatedPolicy, nil
}

func (s *pagesService) CreatePage(ctx context.Context, payload CreatePagePayload) (*models.Page, error) {

	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in the context")
	}

	var createdPage *models.Page

	err := s.store.WithPageTransaction(ctx, func(ps store.PageStore) error {
		page, err := s.store.CreatePage(ctx, payload.PageName, payload.PageDescription)
		if err != nil {
			return err
		}

		_, err = s.store.CreatePagePolicy(ctx, page.ID, models.AudienceTypeUser, *currentUserId, models.PrivilegePageAdmin)
		if err != nil {
			return err
		}

		page = createdPage

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdPage, nil
}

func (s *pagesService) GetPagesByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.Page, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	pages, err := s.store.GetPagesByOrganizationId(ctx, organizationId)
	if err != nil {
		ctxlogger.Error("failed to get pages by organization id", zap.Error(err))
		return nil, err
	}

	return pages, nil
}
