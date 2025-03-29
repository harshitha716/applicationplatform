package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationInvitation struct {
	OrganizationInvitationID uuid.UUID                      `gorm:"column:organization_invitation_id;type:uuid;primary_key;default:gen_random_uuid()" json:"organization_invitation_id"`
	OrganizationID           uuid.UUID                      `gorm:"column:organization_id;type:uuid;not null" json:"organization_id"`
	TargetEmail              string                         `gorm:"column:email;type:text;not null" json:"email"`
	Privilege                ResourcePrivilege              `gorm:"column:privilege;type:text;not null" json:"privilege"`
	CreatedAt                time.Time                      `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
	InvitedBy                uuid.UUID                      `gorm:"column:invited_by;type:uuid;not null" json:"invited_by"`
	InvitationStatuses       []OrganizationInvitationStatus `gorm:"foreignKey:OrganizationInvitationID;references:OrganizationInvitationID" json:"organization_invitation_statuses,omitempty"`
	Inviter                  *User                          `gorm:"foreignKey:InvitedBy;references:ID" json:"inviter,omitempty"`
	EmailRetryCount          int                            `gorm:"column:email_retry_count;type:int;not null;default:0" json:"email_retry_count"`
	EmailSentAt              *time.Time                     `gorm:"column:email_sent_at;type:timestamptz" json:"email_sent_at"`
	UpdatedAt                time.Time                      `gorm:"column:updated_at;type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt                *time.Time                     `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (o *OrganizationInvitation) TableName() string {
	return "organization_invitations"
}

func (o *OrganizationInvitation) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`
		(
			EXISTS (
				SELECT 1
				FROM app.flattened_resource_audience_policies frap
				WHERE
					frap.resource_type = 'organization'
					AND frap.resource_id = organization_invitations.organization_id
					AND frap.user_id = ?
					AND frap.privilege = ?
					AND frap.deleted_at IS NULL
			) OR 
			EXISTS (
				SELECT 1
				FROM app.users_with_traits uwt
				WHERE
					uwt.user_id = ?
					AND uwt.email = organization_invitations.email
			)
		) AND (
		 	organization_invitation_id NOT IN (
				SELECT organization_invitation_id
				FROM app.organization_invitation_statuses ois
				WHERE ois.organization_invitation_id = organization_invitations.organization_invitation_id
			)
		)
	`, userId, PrivilegeOrganizationSystemAdmin, userId)
}

func (o *OrganizationInvitation) BeforeCreate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	// check in flattened_resource_audience_policies if user has access to organization
	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", ResourceTypeOrganization, o.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	// if user has no access to organization
	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil

}

func (o *OrganizationInvitation) BeforeUpdate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	// check in flattened_resource_audience_policies if user has access to organization
	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&fraps).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ? AND deleted_at IS NULL", ResourceTypeOrganization, o.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).Limit(1).Find(&fraps).Error
	if err != nil {
		return err
	}

	// if user has no access to organization
	if len(fraps) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil
}

func (o *OrganizationInvitation) BeforeDelete(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	// check in flattened_resource_audience_policies if user has access to organization
	invitations := []OrganizationInvitation{}
	err := db.Model(&invitations).Joins(`
		JOIN flattened_resource_audience_policies frap ON
			frap.resource_type = ?
			frap.resource_id = organization_invitations.organization_id
			AND organization_invitations.organization_invitation_id = ?
			AND frap.user_id = ?
			AND frap.privilege = ?
			AND frap.deleted_at IS NULL
	`, ResourceTypeOrganization, o.OrganizationInvitationID, userId, PrivilegeOrganizationSystemAdmin).Limit(1).Find(&invitations).Error

	if err != nil {
		return fmt.Errorf("error checking user access to invitation: %w", err)
	}

	// if user has no access to organization
	if len(invitations) == 0 {
		return fmt.Errorf("organization invitation access forbidden")
	}

	return nil
}
