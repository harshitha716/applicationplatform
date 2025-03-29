package teams

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/helper"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TeamService interface {
	GetTeamsByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.Team, error)
	GetTeamById(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) (*models.Team, error)
	CreateTeam(ctx context.Context, orgId uuid.UUID, payload CreateTeamPayload) (*models.Team, error)
	DeleteTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID) error
	RenameTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID, payload RenameTeamPayload) error
	AddUserToTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID, userId uuid.UUID) (*models.TeamMembership, error)
	RemoveUserFromTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID, teamMembershipId uuid.UUID) error
}

type teamServiceStore interface {
	store.TeamStore
	GetOrganizationPolicyByUser(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.ResourceAudiencePolicy, error)
}

type teamService struct {
	store teamServiceStore
}

func NewTeamService(store teamServiceStore) TeamService {
	return &teamService{store: store}
}

func (s *teamService) validateContext(ctx context.Context) (uuid.UUID, uuid.UUID, error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, currentUserId, organizationIds := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		ctxlogger.Error("no user ID found in context")
		return uuid.Nil, uuid.Nil, fmt.Errorf("no user ID found in context")
	}

	if len(organizationIds) == 0 {
		ctxlogger.Error("no organization IDs found in context")
		return uuid.Nil, uuid.Nil, fmt.Errorf("no organization IDs found in context")
	}

	return *currentUserId, organizationIds[0], nil
}

func (s *teamService) validateContextWithOrganizationId(ctx context.Context, organizationId uuid.UUID) (uuid.UUID, uuid.UUID, error) {

	userId, orgId, err := s.validateContext(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	if orgId != organizationId {
		return uuid.Nil, uuid.Nil, fmt.Errorf("user is not a member of the organization")
	}

	return userId, orgId, nil
}

func (s *teamService) GetTeamsByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.Team, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	teams, err := s.store.GetTeams(ctx, organizationID)
	if err != nil {
		ctxlogger.Error("error getting teams by organization ID", zap.Error(err))
		return nil, fmt.Errorf("error getting teams by organization ID; please try again later")
	}

	return teams, nil
}

func (s *teamService) GetTeamById(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) (*models.Team, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	team, err := s.store.GetTeam(ctx, organizationId, teamId)
	if err != nil {
		ctxlogger.Error("error getting team by ID", zap.Error(err))
		return nil, fmt.Errorf("error getting team by ID; please try again later")
	}

	return team, nil
}

func (s *teamService) CreateTeam(ctx context.Context, orgId uuid.UUID, payload CreateTeamPayload) (*models.Team, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	userId, orgId, err := s.validateContextWithOrganizationId(ctx, orgId)
	if err != nil {
		ctxlogger.Error("error validating request", zap.Error(err))
		return nil, err
	}

	if payload.Name == "" {
		ctxlogger.Error("name is required")
		return nil, fmt.Errorf("name is required")
	} else if !helper.IsValidShortInput(payload.Name) {
		ctxlogger.Error("name is too long")
		return nil, fmt.Errorf("name is too long; must be less than 24 characters")
	}

	if payload.Description != "" {
		if !helper.IsValidMediumInput(payload.Description) {
			ctxlogger.Error("description is too long")
			return nil, fmt.Errorf("description is too long; must be less than 64 characters")
		}
	}

	if payload.ColorHexCode == "" {
		ctxlogger.Error("color hex code is required")
		return nil, fmt.Errorf("color hex code is required")
	} else if !helper.IsValidHexCode(payload.ColorHexCode) {
		ctxlogger.Error("invalid color hex code")
		return nil, fmt.Errorf("invalid color hex code")
	}

	metadata := models.TeamMetadata{
		ColorHexCode: payload.ColorHexCode,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		ctxlogger.Error("error marshalling metadata", zap.Error(err))
		return nil, fmt.Errorf("error creating team; please try again later")
	}

	var createdTeam *models.Team

	err = s.store.WithTeamTransaction(ctx, func(txnStore store.TeamStore) error {

		existingTeams, err := txnStore.GetTeamByName(ctx, orgId, payload.Name)
		if err != nil {
			ctxlogger.Error("error getting team by name", zap.Error(err))
			return fmt.Errorf("error getting team by name; please try again later")
		}

		if len(existingTeams) > 0 {
			ctxlogger.Error("team name already exists", zap.String("team_name", payload.Name))
			return fmt.Errorf("Team with name %s already exists", payload.Name)
		}

		team := models.Team{
			OrganizationID: orgId,
			Name:           payload.Name,
			Description:    payload.Description,
			CreatedBy:      userId,
			Metadata:       metadataJSON,
		}

		createdTeam, err = txnStore.CreateOrganizationTeam(ctx, orgId, team)
		if err != nil {
			ctxlogger.Error("error creating team", zap.Error(err))
			return fmt.Errorf("error creating team; please try again later")
		}

		return nil
	})

	if err != nil {
		ctxlogger.Error("error creating team", zap.Error(err))
		return nil, fmt.Errorf("error creating team; please try again later")
	}

	return createdTeam, nil
}

func (s *teamService) DeleteTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID) error {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, orgId, err := s.validateContextWithOrganizationId(ctx, orgId)
	if err != nil {
		ctxlogger.Error("error validating request", zap.Error(err))
		return err
	}

	if teamId == uuid.Nil {
		ctxlogger.Error("team ID is required")
		return fmt.Errorf("team ID is required")
	}

	err = s.store.DeleteTeam(ctx, orgId, teamId)
	if err != nil {
		ctxlogger.Error("error deleting team", zap.Error(err))
		return fmt.Errorf("error deleting team; please try again later")
	}

	return nil
}

func (s *teamService) RenameTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID, payload RenameTeamPayload) error {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, orgId, err := s.validateContextWithOrganizationId(ctx, orgId)
	if err != nil {
		ctxlogger.Error("error validating request", zap.Error(err))
		return err
	}

	if teamId == uuid.Nil {
		ctxlogger.Error("team ID is required")
		return fmt.Errorf("team ID is required")
	}

	if payload.Name == "" {
		ctxlogger.Error("name is required")
		return fmt.Errorf("name is required")
	} else if !helper.IsValidShortInput(payload.Name) {
		ctxlogger.Error("name is too long")
		return fmt.Errorf("name is too long; must be less than 24 characters")
	}

	team := models.Team{
		TeamID:         teamId,
		OrganizationID: orgId,
		Name:           payload.Name,
	}

	_, err = s.store.UpdateTeam(ctx, orgId, team)
	if err != nil {
		ctxlogger.Error("error renaming team", zap.Error(err))
		return fmt.Errorf("error renaming team; please try again later")
	}

	return nil
}

func (s *teamService) AddUserToTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID, userId uuid.UUID) (*models.TeamMembership, error) {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, orgId, err := s.validateContextWithOrganizationId(ctx, orgId)
	if err != nil {
		ctxlogger.Error("error validating request", zap.Error(err))
		return nil, err
	}

	if teamId == uuid.Nil {
		ctxlogger.Error("team ID is required")
		return nil, fmt.Errorf("team ID is required")
	}

	var teamMembership *models.TeamMembership

	// check if user is a member of the organization
	// TODO: move this inside the transaction
	policy, err := s.store.GetOrganizationPolicyByUser(ctx, orgId, userId)
	if err != nil || policy == nil {
		ctxlogger.Error("error getting organization policy by user", zap.Error(err))
		return nil, fmt.Errorf("user is not a member of this organization")
	}

	err = s.store.WithTeamTransaction(ctx, func(txnStore store.TeamStore) error {
		team, err := txnStore.GetTeam(ctx, orgId, teamId)
		if err != nil {
			ctxlogger.Error("error getting team", zap.Error(err))
			return fmt.Errorf("team not found")
		}

		for _, membership := range team.TeamMemberships {
			if membership.UserID == userId {
				return fmt.Errorf("user is already a member of this team")
			}
		}

		teamMembership, err = txnStore.CreateTeamMembership(ctx, teamId, userId)
		if err != nil {
			ctxlogger.Error("error creating team membership", zap.Error(err))
			return fmt.Errorf("error creating team membership; please try again later")
		}

		return nil
	})

	if err != nil {
		ctxlogger.Error("error adding user to team", zap.Error(err))
		return nil, fmt.Errorf("error adding user to team; please try again later")
	}

	return teamMembership, nil

}

func (s *teamService) RemoveUserFromTeam(ctx context.Context, orgId uuid.UUID, teamId uuid.UUID, teamMembershipId uuid.UUID) error {
	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, orgId, err := s.validateContextWithOrganizationId(ctx, orgId)
	if err != nil {
		ctxlogger.Error("error validating request", zap.Error(err))
		return err
	}

	if teamId == uuid.Nil {
		ctxlogger.Error("team ID is required")
		return fmt.Errorf("team ID is required")
	}

	if teamMembershipId == uuid.Nil {
		ctxlogger.Error("team membership ID is required")
		return fmt.Errorf("team membership ID is required")
	}

	err = s.store.WithTeamTransaction(ctx, func(txnStore store.TeamStore) error {
		team, err := txnStore.GetTeam(ctx, orgId, teamId)
		if err != nil {
			ctxlogger.Error("error getting team", zap.Error(err))
			return fmt.Errorf("team not found")
		}

		for _, membership := range team.TeamMemberships {
			if membership.TeamMembershipID == teamMembershipId {
				err = txnStore.DeleteTeamMembership(ctx, teamId, teamMembershipId)
				if err != nil {
					ctxlogger.Error("error deleting team membership", zap.Error(err))
					return fmt.Errorf("error deleting team membership; please try again later")
				}
				return nil
			}
		}

		return fmt.Errorf("team membership not found")
	})

	if err != nil {
		ctxlogger.Error("error removing user from team", zap.Error(err))
		return fmt.Errorf("error removing user from team")
	}

	return nil
}
