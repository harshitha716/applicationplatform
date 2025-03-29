package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvitationStatus string

const (
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
)

type OrganizationInvitationStatus struct {
	OrganizationInvitationID uuid.UUID        `gorm:"column:organization_invitation_id;type:uuid;not null;primaryKey" json:"organization_invitation_id"`
	Status                   InvitationStatus `gorm:"column:status;type:text;not null" json:"status"`
	CreatedAt                time.Time        `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt                time.Time        `gorm:"column:updated_at;type:timestamptz;not null;default:now()" json:"updated_at"`
}

func (o *OrganizationInvitationStatus) TableName() string {
	return "organization_invitation_statuses"
}

func (o *OrganizationInvitationStatus) GetQueryFilters(db *gorm.DB, userID uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`
		EXISTS (
			SELECT 1
			FROM app.flattened_resource_audience_policies frap, app.organization_invitations oi
			WHERE
				frap.resource_type = ?
				AND frap.resource_id = oi.organization_id
				AND oi.organization_invitation_id = organization_invitation_statuses.organization_invitation_id
				AND frap.user_id = ?
				AND frap.privilege = ?
				AND frap.deleted_at IS NULL
		) OR 
		EXISTS (
			SELECT 1
			FROM app.users_with_traits uwt, app.organization_invitations oi
			WHERE
				oi.organization_invitation_id = organization_invitation_statuses.organization_invitation_id
				AND uwt.email = oi.email
				AND uwt.user_id = ?
		)
	`, ResourceTypeOrganization, userID, PrivilegeOrganizationSystemAdmin, userID)
}

func (o *OrganizationInvitationStatus) BeforeCreate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	orgInvitations := []OrganizationInvitation{}

	err := db.Where(`
		EXISTS (
			SELECT 1
			FROM users_with_traits uwt
			WHERE
				organization_invitation_id = ?
				AND uwt.email = organization_invitations.email
				AND uwt.user_id = ?
		)
	`, o.OrganizationInvitationID, userId).Find(&orgInvitations).Error

	if err != nil {
		return fmt.Errorf("error checking user access to invitation: %w", err)
	}

	if len(orgInvitations) == 0 {
		return fmt.Errorf("user does not have access to the invitation")
	}

	return nil
}
