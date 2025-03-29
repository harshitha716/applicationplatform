package models

import (
	"context"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMembership struct {
	TeamMembershipID uuid.UUID       `json:"team_membership_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TeamID           uuid.UUID       `json:"team_id" gorm:"type:uuid;not null"`
	Team             *Team           `json:"team" gorm:"foreignKey:TeamID;references:TeamID"`
	User             *User           `json:"user" gorm:"foreignKey:UserID;references:ID"`
	UserID           uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	CreatedAt        time.Time       `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt        *gorm.DeletedAt `json:"deleted_at"`
	CreatedBy        uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"`
}

func (t *TeamMembership) TableName() string {
	return "team_memberships"
}

func (t *TeamMembership) GetQueryFilters(db *gorm.DB, currentUserId uuid.UUID, organizationIds []uuid.UUID) *gorm.DB {
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap, "app"."teams" teams
			WHERE frap.resource_type = ?
			AND frap.resource_id = teams.organization_id
			AND teams.team_id = team_memberships.team_id
			AND frap.user_id = ?
		)`,
		ResourceTypeOrganization,
		currentUserId,
	)
}

func (t *TeamMembership) validateContext(ctx context.Context) (*uuid.UUID, []uuid.UUID, error) {

	_, userId, orgIds := apicontext.GetAuthFromContext(ctx)

	if userId == nil {
		return nil, nil, fmt.Errorf("no user id found in context")
	}

	if len(orgIds) == 0 {
		return nil, nil, fmt.Errorf("user is not part of the organization")
	}

	return userId, orgIds, nil
}

func (t *TeamMembership) BeforeCreate(db *gorm.DB) error {

	userId, _, err := t.validateContext(db.Statement.Context)

	if err != nil {
		return err
	}

	teams := []Team{}

	err = db.Where(`
		EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE teams.team_id = ?
			AND frap.resource_type = ?
			AND frap.resource_id = teams.organization_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
			AND frap.privilege = ?
		)
	`, t.TeamID, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).Find(&teams).Error

	if err != nil {
		return err
	}

	if len(teams) == 0 {
		return fmt.Errorf("no access to modify team memberships")
	}

	return nil
}

func (t *TeamMembership) BeforeUpdate(db *gorm.DB) error {
	return fmt.Errorf("update forbidden")
}

func (t *TeamMembership) BeforeDelete(db *gorm.DB) error {

	userId, _, err := t.validateContext(db.Statement.Context)

	if err != nil {
		return err
	}

	teams := []Team{}

	err = db.Where(`
		EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE teams.team_id = ?
			AND frap.resource_type = ?
			AND frap.resource_id = teams.organization_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
			AND frap.privilege = ?
		)
	`, t.TeamID, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin).Find(&teams).Error

	if err != nil {
		return err
	}

	if len(teams) == 0 {
		return fmt.Errorf("no access to modify team memberships")
	}

	return nil
}
