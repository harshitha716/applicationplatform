package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrgMembershipStatus string

const (
	OrgMembershipStatusPending  OrgMembershipStatus = "pending"
	OrgMembershipStatusApproved OrgMembershipStatus = "approved"
	OrgMembershipStatusRejected OrgMembershipStatus = "rejected"
)

type OrganizationMembershipRequest struct {
	ID             uuid.UUID           `json:"id" gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID           `json:"organization_id" gorm:"column:organization_id"`
	UserID         uuid.UUID           `json:"user_id" gorm:"column:user_id"`
	User           User                `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time           `json:"created_at" gorm:"column:created_at;default:now()"`
	Status         OrgMembershipStatus `json:"status" gorm:"column:status"`
	UpdatedAt      time.Time           `json:"updated_at" gorm:"column:updated_at;default:now()"`
	DeletedAt      time.Time           `json:"deleted_at" gorm:"column:deleted_at"`
}

func (o *OrganizationMembershipRequest) TableName() string {
	return "organization_membership_requests"
}

func (o *OrganizationMembershipRequest) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`
		EXISTS (
			SELECT 1
			FROM app.flattened_resource_audience_policies frap
			WHERE
				frap.resource_type = ?
				AND frap.resource_id = organization_membership_requests.organization_id
				AND frap.user_id = ?
				AND frap.privilege = ?	
		) OR organization_membership_requests.user_id = ?
	`, ResourceTypeOrganization, userId, PrivilegeOrganizationSystemAdmin, userId)
}

func (o *OrganizationMembershipRequest) BeforeCreate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	if o.UserID != *userId {
		return fmt.Errorf("user cannot create request for another user")
	}

	return nil
}

func (o *OrganizationMembershipRequest) BeforeUpdate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)

	frap := []FlattenedResourceAudiencePolicy{}
	err := db.Model(&frap).Where("resource_type = ? AND resource_id = ? AND user_id = ? AND privilege = ?", ResourceTypeOrganization, o.OrganizationID, userId, PrivilegeOrganizationSystemAdmin).Limit(1).Find(&frap).Error
	if err != nil {
		return fmt.Errorf("failed to get flattened resource audience policy")
	}

	if len(frap) == 0 {
		return fmt.Errorf("organization access forbidden")
	}

	return nil
}
