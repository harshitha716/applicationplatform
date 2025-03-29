package organizations

import (
	"context"
	"fmt"
	"slices"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/mailer"
	"github.com/Zampfi/application-platform/services/api/core/organizations/teams"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/helper"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/errorreporting"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrganizationServiceStore interface {
	store.OrganizationStore
	store.UserStore
}

type OrganizationService interface {
	GetOrganizations(ctx context.Context) ([]models.Organization, error)
	GetOrganizationAudiences(ctx context.Context, organizationId uuid.UUID) ([]models.ResourceAudiencePolicy, error)
	UpdateMemberRole(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)
	BulkInviteMembers(ctx context.Context, organizationId uuid.UUID, payload BulkInvitationPayload) ([]models.OrganizationInvitation, BulkInvitationError)
	GetAllOrganizationInvitations(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationInvitation, error)
	RemoveOrganizationMember(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) error
	CreateOrganization(ctx context.Context, name string, description *string, ownerId uuid.UUID) (*models.Organization, error)
	GetOrganizationMembershipRequestsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationMembershipRequest, error)
	GetOrganizationMembershipRequestsAll(ctx context.Context) ([]models.OrganizationMembershipRequest, error)
	ApprovePendingOrganizationMembershipRequest(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.OrganizationMembershipRequest, error)
	ValidateAudienceInOrganization(ctx context.Context, organizationId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error
	TeamService() teams.TeamService
}

type organizationService struct {
	store        OrganizationServiceStore
	mailerClient mailer.MailerService
	teamService  teams.TeamService
}

func NewOrganizationService(serverConfig *serverconfig.ServerConfig) OrganizationService {
	mailer := mailer.NewMailerService(serverConfig.SparkpostClient, serverConfig.Env.ZampEmailUpdatesFrom, serverConfig.Env.EmailTemplatesPath)
	teamService := teams.NewTeamService(serverConfig.Store)
	return &organizationService{store: serverConfig.Store, mailerClient: mailer, teamService: teamService}
}

func (s *organizationService) GetOrganizations(ctx context.Context) ([]models.Organization, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	var orgs []models.Organization

	orgs, err := s.store.GetOrganizationsAll(ctx)
	if err != nil {
		ctxLogger.Error("failed to get organizations", zap.Error(err))
		return nil, err
	}

	return orgs, nil
}

func (s *organizationService) GetOrganizationAudiences(ctx context.Context, organizationId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	policies, err := s.store.GetOrganizationPolicies(ctx, organizationId)
	if err != nil {
		ctxLogger.Error("failed to get organization policies", zap.Error(err))
		return nil, err
	}

	return policies, nil
}

func (s *organizationService) UpdateMemberRole(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// check user's access to the organization
	_, currentUserId, orgIds := apicontext.GetAuthFromContext(ctx)
	if !slices.Contains(orgIds, organizationId) {
		ctxLogger.Error("user does not have access to the organization", zap.String("organizationId", organizationId.String()))
		return nil, fmt.Errorf("forbidden")
	}

	// restrit user from changing their own role
	if currentUserId == &userId {
		ctxLogger.Error("user cannot change their own role", zap.String("userId", userId.String()))
		return nil, fmt.Errorf("forbidden")

	}

	// validate if the privilege is an org privilege
	if !isOrganizationPrivilege(privilege) {
		ctxLogger.Error("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	// get existing policy of the given audiences
	policy, err := s.store.GetOrganizationPolicyByUser(ctx, organizationId, userId)

	// if the policy does not exist -- throw an error
	if err != nil {
		ctxLogger.Error("failed to get organization policy", zap.Error(err))
		return nil, err
	}

	// if the policy exists, update the privilege if it is different
	if policy != nil && policy.Privilege == privilege {
		return policy, nil
	}

	// update the privilege
	policy, err = s.store.UpdateOrganizationPolicy(ctx, organizationId, userId, privilege)
	if err != nil {
		ctxLogger.Error("failed to update member role", zap.Error(err))
		return nil, err
	}

	return policy, nil
}

func (s *organizationService) getOranizationAccessStateByEmail(ctx context.Context, organizationId uuid.UUID, userEmail string) (userMembershipState, error) {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	organization, err := s.store.GetOrganizationInvitationsAndMembershipRequests(ctx, organizationId)
	if err != nil {
		ctxLogger.Error("failed to get organization invitations and membership requests", zap.Error(err))
		return userMembershipStateNone, err
	}

	userMembershipState := userMembershipStateNone

	for _, invitation := range organization.Invitations {
		if helper.AreEmailsEqual(invitation.TargetEmail, userEmail) && len(invitation.InvitationStatuses) == 0 {
			userMembershipState = userMembershipStateInvited
		}
	}

	for _, membershipRequest := range organization.MembershipRequests {
		if helper.AreEmailsEqual(membershipRequest.User.Email, userEmail) && membershipRequest.Status == models.OrgMembershipStatusPending {
			userMembershipState = userMembershipStateUnderReview
		}
	}

	return userMembershipState, nil
}

func (s *organizationService) approvePendingMembershipRequestFromInvitee(ctx context.Context, organizationId uuid.UUID, userEmail string) error {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	organization, err := s.store.GetOrganizationInvitationsAndMembershipRequests(ctx, organizationId)
	if err != nil {
		ctxLogger.Error("failed to get organization invitations and membership requests", zap.Error(err))
		return err
	}

	inviteeUserId := uuid.Nil
	for _, request := range organization.MembershipRequests {
		if helper.AreEmailsEqual(request.User.Email, userEmail) && request.Status == models.OrgMembershipStatusPending {
			inviteeUserId = request.User.ID
		}
	}

	if inviteeUserId == uuid.Nil {
		return fmt.Errorf("no pending request found")
	}

	_, err = s.ApprovePendingOrganizationMembershipRequest(ctx, organizationId, inviteeUserId)
	if err != nil {
		ctxLogger.Error("failed to approve pending membership request", zap.Error(err))
		return err
	}

	return nil

}

func (s *organizationService) inviteMember(ctx context.Context, organizationId uuid.UUID, userEmail string, privilege models.ResourcePrivilege) (*models.OrganizationInvitation, error) {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// check user's access to the organization
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)
	if !slices.Contains(orgIds, organizationId) {
		ctxLogger.Error("user does not have access to the organization", zap.String("organizationId", organizationId.String()))
		return nil, fmt.Errorf("forbidden")
	}

	email := helper.SanitizeEmail(userEmail)
	// validate email
	valid := helper.IsValidEmail(email)
	if !valid {
		ctxLogger.Error("invalid email", zap.String("email", email))
		return nil, fmt.Errorf("invalid email")
	}

	// check if the given privilege is a valid organization privilege
	if !isOrganizationPrivilege(privilege) {
		ctxLogger.Error("invalid privilege", zap.String("privilege", string(privilege)))
		return nil, fmt.Errorf("invalid privilege")
	}

	// check if the user is not inviting an existing member
	var invitation *models.OrganizationInvitation
	var err error

	err = s.store.WithOrganizationTransaction(ctx, func(txStore store.OrganizationStore) error {

		existingInvitations, err := txStore.GetOrganizationInvitationsByOrganizationId(ctx, organizationId)
		if err != nil {
			ctxLogger.Error("failed to get organization invitations", zap.Error(err))
			return err
		}

		for _, invitation := range existingInvitations {
			if helper.AreEmailsEqual(invitation.TargetEmail, email) {
				ctxLogger.Error("user is already invited to the organization", zap.String("email", email))
				return fmt.Errorf("User is already invited to the organization.")
			}
		}

		existingMembers, errr := txStore.GetOrganizationPoliciesByEmail(ctx, organizationId, email)
		if errr != nil {
			ctxLogger.Error("failed to get organization policy by email", zap.Error(err))
			errorreporting.CaptureException(fmt.Errorf("failed to get organization policy by email: %s", err), ctx)
			return fmt.Errorf("Something went wrong. Please try again later.")
		}

		if len(existingMembers) > 0 {
			ctxLogger.Error("user is already a member of the organization", zap.String("email", email))
			err = fmt.Errorf("user is already a member of the organization")
			return fmt.Errorf("User is already a member of the organization")
		}

		invitation, err = txStore.CreateOrganizationInvitation(ctx, organizationId, email, privilege)
		if err != nil {
			ctxLogger.Error("failed to create organization invitation", zap.Error(err))
			errorreporting.CaptureException(fmt.Errorf("failed to create organization invitation: %s", err), ctx)
			return fmt.Errorf("Something went wrong. Please try again later.")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// if the invited user is already a member, approve the request
	err = s.approvePendingMembershipRequestFromInvitee(ctx, organizationId, userEmail)
	if err != nil {
		ctxLogger.Error("could not approve pending membership request", zap.Error(err))
		err = s.sendInvitationEmail(ctx, invitation)
		if err != nil {
			ctxLogger.Error("failed to send invitation email", zap.Error(err))
			errorreporting.CaptureException(fmt.Errorf("failed to send invitation email: %s", err), ctx)
		} else {
			ctxLogger.Info("invitation email sent", zap.String("email", email))
		}
	}
	return invitation, nil
}

func (s *organizationService) sendInvitationEmail(ctx context.Context, invitation *models.OrganizationInvitation) error {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	organization, err := s.store.GetOrganizationById(ctx, invitation.OrganizationID.String())
	if err != nil {
		ctxLogger.Error("failed to get organization by id", zap.Error(err))
		return fmt.Errorf("failed to get organization by id")
	}

	inviter, err := s.store.GetUserById(ctx, invitation.InvitedBy.String())
	if err != nil {
		ctxLogger.Error("failed to get user by id", zap.Error(err))
		return fmt.Errorf("failed to get user by id")
	}

	inviterName := inviter.Name
	if inviterName == "" {
		inviterName = helper.GetNameFromEmail(inviter.Email)
	}

	err = s.mailerClient.SendInvitationEmail(ctx, mailer.InvitationEmailData{
		OrganizationName:   organization.Name,
		RecipientEmail:     invitation.TargetEmail,
		InvitedByFirstName: inviterName,
		InvitationLink:     "https://app.zamp.ai", // TODO: construct from root domain
	})

	return err

}

func (s *organizationService) BulkInviteMembers(ctx context.Context, organizationId uuid.UUID, payload BulkInvitationPayload) ([]models.OrganizationInvitation, BulkInvitationError) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// check user's access to the organization
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)
	if !slices.Contains(orgIds, organizationId) {
		ctxLogger.Error("user does not have access to the organization", zap.String("organizationId", organizationId.String()))
		return nil, BulkInvitationError{
			Error: fmt.Errorf("forbidden"),
		}
	}

	// allow inviting only

	inviatationErrors := BulkInvitationError{}
	invitations := []models.OrganizationInvitation{}

	for _, p := range payload.Invitations {
		invitation, err := s.inviteMember(ctx, organizationId, p.Email, models.ResourcePrivilege(p.Privilege))
		if err != nil {
			ctxLogger.Error(fmt.Sprintf("failed to invite member: %s", p.Email), zap.Error(err))
			inviatationErrors.Invitations = append(inviatationErrors.Invitations, InvitationError{
				Email:        p.Email,
				ErrorMessage: err.Error(),
			})
			continue
		}

		if invitation == nil {
			ctxLogger.Error(fmt.Sprintf("invitation is nil: %s", p.Email))
			errorreporting.CaptureException(fmt.Errorf("invitation is nil: %s", p.Email), ctx)
			inviatationErrors.Invitations = append(inviatationErrors.Invitations, InvitationError{
				Email:        p.Email,
				ErrorMessage: "unexpected",
			})
			continue
		}

		invitations = append(invitations, *invitation)
	}

	// TODO send email

	return invitations, inviatationErrors
}

func (s *organizationService) GetAllOrganizationInvitations(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationInvitation, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// check user's access to the organization
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)
	if !slices.Contains(orgIds, organizationId) {
		ctxLogger.Error("user does not have access to the organization", zap.String("organizationId", organizationId.String()))
		return nil, fmt.Errorf("forbidden")
	}

	invitations, err := s.store.GetOrganizationInvitationsByOrganizationId(ctx, organizationId)
	if err != nil {
		ctxLogger.Error("failed to get organization invitations", zap.Error(err))
		return nil, err
	}

	return invitations, nil
}

func (s *organizationService) RemoveOrganizationMember(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) error {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// check user's access to the organization
	_, currentUserId, orgIds := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		ctxLogger.Error("no user id found in context")
		return fmt.Errorf("no user id found in context")
	}
	if !slices.Contains(orgIds, organizationId) {
		ctxLogger.Error("user does not have access to the organization", zap.String("organizationId", organizationId.String()))
		return fmt.Errorf("forbidden")
	}

	// ensure if the user is not trying to remove themselves
	if userId == *currentUserId {
		ctxLogger.Error("user cannot remove themselves", zap.String("userId", userId.String()))
		return fmt.Errorf("forbidden")
	}

	// ensure if current user is system admin
	policy, err := s.store.GetOrganizationPolicyByUser(ctx, organizationId, *currentUserId)
	if err != nil {
		ctxLogger.Error("failed to get organization policy", zap.Error(err))
		return err
	}

	if policy == nil {
		ctxLogger.Error("current user is not a member of the organization", zap.String("userId", userId.String()))
		return fmt.Errorf("user is not a member of the organization")
	}

	if policy.Privilege != models.PrivilegeOrganizationSystemAdmin {
		ctxLogger.Error("user does not have permission to remove members", zap.String("userId", userId.String()))
		return fmt.Errorf("forbidden")
	}

	err = s.store.DeleteOrganizationPolicy(ctx, organizationId, userId)
	if err != nil {
		ctxLogger.Error("failed to remove organization member", zap.Error(err))
		return err
	}

	return nil
}
func (s *organizationService) CreateOrganization(ctx context.Context, name string, description *string, ownerId uuid.UUID) (*models.Organization, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// check admin access
	role, _, _ := apicontext.GetAuthFromContext(ctx)
	if role != "admin" {
		ctxLogger.Error("user does not have admin access", zap.String("role", role))
		return nil, fmt.Errorf("forbidden")
	}

	var createdOrganization *models.Organization
	err := s.store.WithOrganizationTransaction(ctx, func(store.OrganizationStore) error {

		// get org policies by owner id
		organization, err := s.store.CreateOrganization(ctx, name, description, ownerId)
		if err != nil {
			ctxLogger.Error("failed to create organization", zap.Error(err))
			return err
		}

		// create org policy for the owner
		_, err = s.store.CreateOrganizationPolicy(ctx, organization.ID, models.AudienceTypeUser, ownerId, models.PrivilegeOrganizationSystemAdmin)
		if err != nil {
			ctxLogger.Error("failed to create organization policy", zap.Error(err))
			return err
		}

		createdOrganization = organization

		return nil
	})

	if err != nil {
		ctxLogger.Error("failed to create organization", zap.Error(err))
		return nil, err
	}

	return createdOrganization, nil
}

func (s *organizationService) GetOrganizationMembershipRequestsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationMembershipRequest, error) {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	// Get membership requests for organization
	requests, err := s.store.GetOrganizationMembershipRequestsByOrganizationId(ctx, organizationId)
	if err != nil {
		ctxLogger.Error("failed to get organization membership requests", zap.Error(err))
		return nil, err
	}

	return requests, nil
}

func (s *organizationService) GetOrganizationMembershipRequestsAll(ctx context.Context) ([]models.OrganizationMembershipRequest, error) {
	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	request, err := s.store.GetOrganizationMembershipRequestsAll(ctx)
	if err != nil {
		ctxLogger.Error("failed to get organization membership request", zap.Error(err))
		return nil, err
	}

	return request, nil
}

func (s *organizationService) updatePendingOrganizationMembershipRequest(ctx context.Context, txnStore store.OrganizationStore, organizationId uuid.UUID, userId uuid.UUID, status models.OrgMembershipStatus) (*models.OrganizationMembershipRequest, error) {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	request, err := txnStore.UpdatePendingOrganizationMembershipRequest(ctx, organizationId, userId, status)
	if err != nil {
		ctxLogger.Error("failed to update organization membership request", zap.Error(err))
		return nil, err
	}

	return request, nil
}

func (s *organizationService) ApprovePendingOrganizationMembershipRequest(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.OrganizationMembershipRequest, error) {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		ctxLogger.Error("no user id found in context")
		return nil, fmt.Errorf("no user id found in context")
	}

	var updatedRequest *models.OrganizationMembershipRequest
	var err error
	err = s.store.WithOrganizationTransaction(ctx, func(txnStore store.OrganizationStore) error {

		// get organization policies
		policies, err := txnStore.GetOrganizationPolicies(ctx, organizationId)
		if err != nil {
			ctxLogger.Error("failed to get organization policies", zap.Error(err))
			return err
		}

		var currentUserPolicy *models.ResourceAudiencePolicy
		for _, policy := range policies {
			if policy.ResourceAudienceType == models.AudienceTypeUser && policy.ResourceAudienceID == *currentUserId {
				currentUserPolicy = &policy
			}
		}

		if currentUserPolicy == nil {
			ctxLogger.Error("current user is not a member of the organization", zap.String("userId", userId.String()))
			return fmt.Errorf("current user is not a member of the organization")
		}

		if currentUserPolicy.Privilege != models.PrivilegeOrganizationSystemAdmin {
			ctxLogger.Error("current user does not have permission to approve membership request", zap.String("userId", userId.String()))
			return fmt.Errorf("You do not have the necessary privileges to approve membership requests")
		}

		var requesterPolicy *models.ResourceAudiencePolicy
		for _, policy := range policies {
			if policy.ResourceAudienceType == models.AudienceTypeUser && policy.ResourceAudienceID == userId {
				requesterPolicy = &policy
			}
		}

		if requesterPolicy != nil {
			ctxLogger.Error("requester is already a member of the organization", zap.String("userId", userId.String()))
			return fmt.Errorf("requester is already a member of the organization")
		}

		// create org policy for the user
		_, err = txnStore.CreateOrganizationPolicy(ctx, organizationId, models.AudienceTypeUser, userId, models.PrivilegeOrganizationMember)
		if err != nil {
			ctxLogger.Error("failed to create organization policy", zap.Error(err))
			return err
		}

		updatedRequest, err = s.updatePendingOrganizationMembershipRequest(ctx, txnStore, organizationId, userId, models.OrgMembershipStatusApproved)
		if err != nil {
			ctxLogger.Error("failed to update organization membership request", zap.Error(err))
			return err
		}
		return nil
	})

	return updatedRequest, err
}

func (s *organizationService) TeamService() teams.TeamService {
	return s.teamService
}

func (s *organizationService) ValidateAudienceInOrganization(ctx context.Context, organizationId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {

	ctxLogger := apicontext.GetLoggerFromCtx(ctx)

	policies, err := s.store.GetOrganizationPolicies(ctx, organizationId)
	if err != nil {
		ctxLogger.Error("failed to get organization policies", zap.Error(err))
		return fmt.Errorf("failed to find organization information")
	}

	for _, policy := range policies {
		if policy.ResourceAudienceType == audienceType && policy.ResourceAudienceID == audienceId {
			return nil
		}
	}

	return fmt.Errorf("audience not found in organization")
}
