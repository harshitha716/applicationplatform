package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMetadata struct {
	ColorHexCode string `json:"color_hex_code"`
}

type Team struct {
	TeamID          uuid.UUID        `json:"team_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID  uuid.UUID        `json:"organization_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	CreatedAt       time.Time        `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt       *gorm.DeletedAt  `json:"deleted_at"`
	Metadata        json.RawMessage  `json:"metadata" gorm:"type:jsonb;default:null"`
	TeamMemberships []TeamMembership `json:"team_memberships" gorm:"foreignKey:TeamID"`
	CreatedBy       uuid.UUID        `json:"created_by"`
}

func (t *Team) TableName() string {
	return "teams"
}

func (t *Team) GetQueryFilters(db *gorm.DB, currentUserId uuid.UUID, organizationIds []uuid.UUID) *gorm.DB {

	// TODO: Terrible query, need to optimize
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = ?
			AND frap.resource_id = teams.organization_id
			AND frap.user_id = ?
		)`,
		ResourceTypeOrganization,
		currentUserId,
	)
}

func (t *Team) validateContext(ctx context.Context) (*uuid.UUID, []uuid.UUID, error) {

	_, userId, orgIds := apicontext.GetAuthFromContext(ctx)

	if userId == nil {
		return nil, nil, fmt.Errorf("no user id found in context")
	}

	if len(orgIds) == 0 {
		return nil, nil, fmt.Errorf("user is not part of the organization")
	}

	return userId, orgIds, nil
}

func (t *Team) BeforeCreate(db *gorm.DB) error {

	userId, orgIds, err := t.validateContext(db.Statement.Context)

	if err != nil {
		return err
	}

	if orgIds[0] != t.OrganizationID {
		return fmt.Errorf("user is not part of the organization")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	db.
		Where("resource_type = ? AND resource_id = ?  AND user_id = ? AND privilege = ?", ResourceTypeOrganization, t.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).
		Find(&fraps)

	if len(fraps) == 0 {
		return fmt.Errorf("user does not have privileges to create a team")
	}

	return nil
}

func (t *Team) BeforeUpdate(db *gorm.DB) error {

	userId, orgIds, err := t.validateContext(db.Statement.Context)

	if err != nil {
		return err
	}

	if orgIds[0] != t.OrganizationID {
		return fmt.Errorf("user is not part of the organization")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	db.Where(`
		resource_type = ? AND resource_id = ?
		AND user_id = ? AND privilege = ?`, ResourceTypeOrganization, t.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).Find(&fraps)

	if len(fraps) == 0 {
		return fmt.Errorf("user does not have privileges to update a team")
	}

	return nil
}

func (t *Team) BeforeDelete(db *gorm.DB) error {

	userId, orgIds, err := t.validateContext(db.Statement.Context)

	if err != nil {
		return err
	}

	if orgIds[0] != t.OrganizationID {
		return fmt.Errorf("user is not part of the organization")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	db.Where(`
		resource_type = ? AND resource_id = ?
		AND user_id = ? AND privilege = ?`, ResourceTypeOrganization, t.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).Find(&fraps)

	if len(fraps) == 0 {
		return fmt.Errorf("user does not have privileges to delete a team")
	}

	return nil
}
