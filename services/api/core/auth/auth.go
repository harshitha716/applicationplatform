package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/helper"
	"github.com/Zampfi/application-platform/services/api/helper/constants"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/errorreporting"
	"github.com/Zampfi/application-platform/services/api/pkg/kratosclient"

	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

type AuthServiceStore interface {
	store.OrganizationStore
	store.UserStore
}

type AuthService interface {
	ResolveSessionCookie(ctx context.Context, cookie string) (*kratos.Session, *http.Response, *kratosclient.Error)
	ResolveAdminInfo(ctx context.Context, header http.Header) (role string, emulatedUserId uuid.UUID, emulatedOrganizationIdss []uuid.UUID)
	GetUserOrganizations(ctx context.Context, userId uuid.UUID) ([]models.Organization, error)
	GetKratosProxy() (*httputil.ReverseProxy, url.URL)
	GetAuthFlowForUser(ctx context.Context, email string) (*kratos.LoginFlow, *http.Response, error)
	HandleNewUserCreated(ctx context.Context, adminSecret string, userId uuid.UUID) error
	SignupUserAsAdmin(ctx context.Context, email string, password string) (*models.User, error)
	IsUserExposedKratosPath(path string) bool
	GetCurrentUserInfo(ctx context.Context) (*models.User, error)
}

type authService struct {
	// The server context for the service
	adminSecrets     []string
	kratosClient     kratosclient.KratosClient
	authServiceStore AuthServiceStore
	environment      string
}

// NewAuthService creates a new AuthService
func NewAuthService(adminSecrets []string, authUrl string, appStore store.Store, environment string) (AuthService, error) {
	kratosClient, err := kratosclient.NewClient(authUrl)
	if err != nil {
		return nil, err
	}

	return &authService{
		adminSecrets:     adminSecrets,
		kratosClient:     kratosClient,
		authServiceStore: appStore,
		environment:      environment,
	}, nil
}

func (a *authService) ResolveSessionCookie(ctx context.Context, cookie string) (*kratos.Session, *http.Response, *kratosclient.Error) {

	logger := apicontext.GetLoggerFromCtx(ctx)

	session, httpResp, err := a.kratosClient.GetSessionInfo(ctx, logger, cookie)
	if err != nil {
		logger.Error("failed to get session", zap.String("error", err.Message))
		return nil, nil, err
	}

	return session, httpResp, nil
}

func (a *authService) ResolveAdminInfo(ctx context.Context, header http.Header) (role string, emulatedUserId uuid.UUID, emulatedOrganizationIds []uuid.UUID) {

	logger := apicontext.GetLoggerFromCtx(ctx)

	secret := helper.GetAdminSecretFromHeader(header)
	validAdminSecret := false
	for _, s := range a.adminSecrets {
		if s == secret {
			validAdminSecret = true
		}
	}
	if !validAdminSecret {
		return "anonymous", uuid.Nil, []uuid.UUID{}
	}

	userId := header.Get(helper.PROXY_USER_ID_HEADER)
	if userId == "" {
		logger.Info("no user id found in header")
		return "anonymous", uuid.Nil, []uuid.UUID{}
	}

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		logger.Info("invalid user id found in header, not uuid")
		return "anonymous", uuid.Nil, []uuid.UUID{}
	}

	// First check for the single organization ID header
	orgId := header.Get(helper.PROXY_WORKSPACE_ID_HEADER)
	if orgId != "" {
		orgIdUUID, err := uuid.Parse(orgId)
		if err == nil {
			return "admin", userIdUUID, []uuid.UUID{orgIdUUID}
		}
	}

	// Fall back to the multiple organization IDs header
	organizationIdsStrings := header.Values(helper.PROXY_WORKSPACE_IDS_HEADER)
	if len(organizationIdsStrings) == 0 {
		logger.Info("no organization ids found in header")
		return "anonymous", uuid.Nil, []uuid.UUID{}
	}

	var organizationIds []uuid.UUID
	for _, id := range organizationIdsStrings {
		idUUID, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		organizationIds = append(organizationIds, idUUID)
	}

	return "admin", userIdUUID, organizationIds

}

func (a *authService) GetUserOrganizations(ctx context.Context, userId uuid.UUID) ([]models.Organization, error) {
	return a.authServiceStore.GetOrganizationsByMemberId(ctx, userId)
}

func (a *authService) GetKratosProxy() (*httputil.ReverseProxy, url.URL) {
	return a.kratosClient.GetAuthProxy()
}

func (a *authService) GetAuthFlowForUser(ctx context.Context, emailRaw string) (*kratos.LoginFlow, *http.Response, error) {

	logger := apicontext.GetLoggerFromCtx(ctx)

	if !helper.IsValidEmail(emailRaw) {
		logger.Error("invalid email", zap.String("email", emailRaw))
		return nil, nil, fmt.Errorf("invalid email")
	}

	email := helper.SanitizeEmail(emailRaw)

	// extract the domain from the email
	domain := helper.GetDomainFromEmail(email)

	// fetch sso config by domain
	ssoConfig, err := a.authServiceStore.GetSSOConfigByDomain(ctx, domain)
	if err != nil || ssoConfig == nil {
		logger.Error("no sso config found for domain", zap.String("domain", domain), zap.Error(err))
		return nil, nil, fmt.Errorf("The given email %s is not registered with any organization", email)
	}

	providerId := ssoConfig.SSOProviderID

	// create a login flow on kratos
	loginFlow, httpResp, kerr := a.kratosClient.CreateLoginFlow(ctx, logger, email)
	if kerr != nil {
		logger.Error("failed to create login flow", zap.String("error", kerr.Message))
		return nil, nil, fmt.Errorf("Something went wrong. Please try again later.")
	}

	if a.environment == constants.ENVLOCAL {
		return loginFlow, httpResp, nil
	}

	// filter the flow's login methods to only include the ones that are configured in the sso config
	uiNodes := []kratos.UiNode{}
	for _, node := range loginFlow.Ui.Nodes {
		if node.Group == "oidc" {
			if node.Attributes.UiNodeInputAttributes.Name == "provider" && node.Attributes.UiNodeInputAttributes.Value == providerId {
				uiNodes = append(uiNodes, node)
			}
		}
	}

	if len(uiNodes) == 0 {
		return nil, nil, fmt.Errorf("The given email %s is not registered with any organization", email)
	}

	loginFlow.Ui.Nodes = uiNodes

	return loginFlow, httpResp, nil

}

func (a *authService) IsUserExposedKratosPath(path string) bool {

	startsWithPaths := []string{"/sessions", "/self-service/methods/oidc"}

	exactPaths := []string{"/self-service/login", "/self-service/logout/browser", "/self-service/logout"}

	for _, p := range startsWithPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	for _, p := range exactPaths {
		if path == p {
			return true
		}
	}

	return false
}

// assigns the new user to the relevant organization
func (a *authService) HandleNewUserCreated(ctx context.Context, adminSecret string, userId uuid.UUID) error {

	logger := apicontext.GetLoggerFromCtx(ctx).With(zap.String("user_id", userId.String()), zap.String("function", "HandleNewUserCreated"))

	if !slices.Contains(a.adminSecrets, adminSecret) {
		logger.Error("invalid admin secret")
		return fmt.Errorf("unauthorized")
	}

	ctx = apicontext.AddAuthToContext(ctx, "admin", userId, []uuid.UUID{})

	// TODO: make the store accept uuid.UUID
	user, err := a.authServiceStore.GetUserById(ctx, userId.String())
	if err != nil {
		logger.Error("failed to get user", zap.String("error", err.Error()))
		return err
	}

	email := helper.SanitizeEmail(user.Email)

	err = a.authServiceStore.WithOrganizationTransaction(ctx, func(orgStore store.OrganizationStore) error {
		// get invitations
		invitations, err := orgStore.GetOrganizationInvitationsAll(ctx)
		if err != nil {
			logger.Error("failed to get invitations", zap.String("error", err.Error()))
			return err
		}

		var currentUserInvitation *models.OrganizationInvitation
		for _, invitation := range invitations {
			if helper.AreEmailsEqual(invitation.TargetEmail, email) {
				currentUserInvitation = &invitation
			}
		}

		// if the user has no invitation, queue a request to the organization to create an invitation
		if currentUserInvitation == nil {
			// create membership request
			// getting a primary SSO config for the given email domain

			ssoConfig, err := a.authServiceStore.GetPrimarySSOConfigByDomain(ctx, helper.GetDomainFromEmail(email))
			if err != nil {
				logger.Error("failed to get sso config", zap.String("error", err.Error()))
				errorreporting.CaptureMessage("user tried to signup with an email that is not registered with any organization; no sso config found", ctx)
				return nil
			}

			logger.Info("creating membership request")
			membershipRequest, err := orgStore.CreateOrganizationMembershipRequest(ctx, ssoConfig.OrganizationID, userId, models.OrgMembershipStatusPending)
			if err != nil {
				logger.Error("failed to create membership request", zap.String("error", err.Error()))
				apicontext.AddCtxVariableToCtx(ctx, "organization_id", ssoConfig.OrganizationID)
				errorreporting.CaptureException(err, ctx)
			} else {
				logger.Info("membership request created", zap.String("membership_request_id", membershipRequest.ID.String()))
			}

		} else {
			// if the user already has an invitation, update the status to accepted
			_, err := orgStore.CreateOrganizationInvitationStatus(ctx, currentUserInvitation.OrganizationInvitationID, models.InvitationStatusAccepted)
			if err != nil {
				logger.Error("failed to create invitation status", zap.String("error", err.Error()))
				return err
			}

			// add resource audience policy
			_, err = orgStore.CreateOrganizationPolicy(ctx, currentUserInvitation.OrganizationID, models.AudienceTypeUser, userId, currentUserInvitation.Privilege)
			if err != nil {
				logger.Error("failed to create resource audience policy", zap.String("error", err.Error()))
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("failed to handle new user created", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *authService) SignupUserAsAdmin(ctx context.Context, email string, password string) (*models.User, error) {

	logger := apicontext.GetLoggerFromCtx(ctx)

	role, _, _ := apicontext.GetAuthFromContext(ctx)
	if role != "admin" {
		logger.Error("unauthorized")
		return nil, fmt.Errorf("unauthorized")
	}

	identity, httpResp, kerr := a.kratosClient.SignupUserEmailPassword(ctx, logger, email, password)

	if kerr != nil {
		statusCode := http.StatusInternalServerError
		if httpResp != nil {
			statusCode = httpResp.StatusCode
		}
		logger.Error("failed to signup user as admin", zap.Int("status_code", statusCode), zap.String("error", kerr.Message))
		return nil, fmt.Errorf("failed to signup user as admin, status code: %d", statusCode)
	}

	logger.Info("signup user as admin success", zap.Any("identity", identity))

	user, err := kratosIdentityToZampUser(identity)
	if err != nil {
		logger.Error("failed to convert kratos identity to zamp user", zap.String("error", err.Error()))
		return nil, err
	}

	return user, nil

}

func (a *authService) GetCurrentUserInfo(ctx context.Context) (*models.User, error) {
	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	user, err := a.authServiceStore.GetUserById(ctx, userId.String())
	if err != nil {
		return nil, err
	}
	return user, nil
}
