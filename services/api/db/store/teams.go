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

type TeamStore interface {
	GetTeams(ctx context.Context, organizationId uuid.UUID) ([]models.Team, error)
	GetTeam(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) (*models.Team, error)
	GetTeamByName(ctx context.Context, organizationId uuid.UUID, name string) ([]models.Team, error)
	CreateOrganizationTeam(ctx context.Context, organizationId uuid.UUID, team models.Team) (*models.Team, error)
	UpdateTeam(ctx context.Context, organizationId uuid.UUID, team models.Team) (*models.Team, error)
	DeleteTeam(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) error
	GetTeamMemberships(ctx context.Context, teamId uuid.UUID) ([]models.TeamMembership, error)
	GetTeamMembershipById(ctx context.Context, teamMembershipId uuid.UUID) (*models.TeamMembership, error)
	GetTeamMembershipByUserIdTeamId(ctx context.Context, userId uuid.UUID, teamId uuid.UUID) (*models.TeamMembership, error)
	CreateTeamMembership(ctx context.Context, teamId uuid.UUID, userId uuid.UUID) (*models.TeamMembership, error)
	DeleteTeamMembership(ctx context.Context, teamId uuid.UUID, teamMembershipId uuid.UUID) error
	WithTeamTransaction(ctx context.Context, fn func(TeamStore) error) error
}

func (s *appStore) GetTeams(ctx context.Context, organizationId uuid.UUID) ([]models.Team, error) {
	teams := []models.Team{}
	err := s.client.DB.WithContext(ctx).Preload("TeamMemberships").Where("organization_id = ?", organizationId).Order("name ASC").Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *appStore) GetTeam(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) (*models.Team, error) {
	team := models.Team{}
	err := s.client.DB.WithContext(ctx).Preload("TeamMemberships").Where("team_id = ? AND organization_id = ?", teamId, organizationId).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *appStore) GetTeamByName(ctx context.Context, organizationId uuid.UUID, name string) ([]models.Team, error) {
	teams := []models.Team{}
	err := s.client.DB.WithContext(ctx).Where("organization_id = ? AND name = ?", organizationId, name).Limit(1).Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *appStore) CreateOrganizationTeam(ctx context.Context, organizationId uuid.UUID, team models.Team) (*models.Team, error) {

	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return nil, fmt.Errorf("no user ID found in context")
	}

	team.OrganizationID = organizationId
	team.CreatedBy = *userId

	err := s.client.DB.WithContext(ctx).Create(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *appStore) UpdateTeam(ctx context.Context, organizationId uuid.UUID, team models.Team) (*models.Team, error) {

	if team.TeamID == uuid.Nil {
		return nil, fmt.Errorf("team id is required")
	}

	team.OrganizationID = organizationId

	err := s.client.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *appStore) DeleteTeam(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) error {
	team := models.Team{
		TeamID:         teamId,
		OrganizationID: organizationId,
	}
	err := s.client.DB.WithContext(ctx).Where("team_id = ?", teamId).Where("organization_id = ?", organizationId).Delete(&team).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *appStore) GetTeamMemberships(ctx context.Context, teamId uuid.UUID) ([]models.TeamMembership, error) {
	teamMemberships := []models.TeamMembership{}
	err := s.client.DB.WithContext(ctx).Preload("Team").Preload("User").Where("team_id = ?", teamId).Find(&teamMemberships).Error
	if err != nil {
		return nil, err
	}
	return teamMemberships, nil
}

func (s *appStore) GetTeamMembershipById(ctx context.Context, teamMembershipId uuid.UUID) (*models.TeamMembership, error) {
	teamMembership := models.TeamMembership{}
	err := s.client.DB.WithContext(ctx).Preload("Team").Preload("User").Where("team_membership_id = ?", teamMembershipId).First(&teamMembership).Error
	if err != nil {
		return nil, err
	}
	return &teamMembership, nil
}

func (s *appStore) GetTeamMembershipByUserIdTeamId(ctx context.Context, userId uuid.UUID, teamId uuid.UUID) (*models.TeamMembership, error) {
	teamMembership := models.TeamMembership{}
	err := s.client.DB.WithContext(ctx).Preload("Team").Preload("User").Where("user_id = ? AND team_id = ?", userId, teamId).First(&teamMembership).Error
	if err != nil {
		return nil, err
	}
	return &teamMembership, nil
}

func (s *appStore) CreateTeamMembership(ctx context.Context, teamId uuid.UUID, userId uuid.UUID) (*models.TeamMembership, error) {

	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		return nil, fmt.Errorf("no user ID found in context")
	}

	if teamId == uuid.Nil {
		return nil, fmt.Errorf("team id is required")
	}

	if userId == uuid.Nil {
		return nil, fmt.Errorf("user id is required")
	}

	teamMembership := models.TeamMembership{
		TeamID:    teamId,
		UserID:    userId,
		CreatedBy: *currentUserId,
	}
	err := s.client.DB.WithContext(ctx).Create(&teamMembership).Error
	if err != nil {
		return nil, err
	}
	return &teamMembership, nil
}

func (s *appStore) DeleteTeamMembership(ctx context.Context, teamId uuid.UUID, teamMembershipId uuid.UUID) error {
	err := s.client.DB.WithContext(ctx).Where("team_id = ?", teamId).Delete(&models.TeamMembership{TeamMembershipID: teamMembershipId, TeamID: teamId}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *appStore) WithTeamTransaction(ctx context.Context, fn func(TeamStore) error) error {
	return s.client.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txStore := &appStore{client: &pgclient.PostgresClient{DB: tx}}
		return fn(txStore)
	})
}
