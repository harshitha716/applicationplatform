package models

import (
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	ID                       uuid.UUID                       `json:"organization_id" gorm:"column:organization_id;default:gen_random_uuid();primaryKey"`
	Name                     string                          `json:"name"`
	Description              *string                         `json:"description"`
	CreatedAt                time.Time                       `json:"created_at" gorm:"default:now()"`
	UpdatedAt                time.Time                       `json:"updated_at" gorm:"default:now()"`
	DeletedAt                *time.Time                      `json:"deleted_at,omitempty"`
	OwnerId                  uuid.UUID                       `json:"owner_id"`
	ResourceAudiencePolicies []ResourceAudiencePolicy        `json:"resource_audience_policies" gorm:"foreignKey:ResourceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Invitations              []OrganizationInvitation        `json:"invitations" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MembershipRequests       []OrganizationMembershipRequest `json:"membership_requests" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SSOConfigs               []OrganizationSSOConfig         `json:"sso_configs" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (o *Organization) TableName() string {
	return "organizations"
}

func (o *Organization) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {

	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'organization'
			AND frap.resource_id = organizations.organization_id
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)`, userId,
	)
}

func (o *Organization) BeforeCreate(db *gorm.DB) error {

	ctx := db.Statement.Context

	role, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return fmt.Errorf("no user id found in context")
	}

	if role == "admin" {
		return nil
	}

	return fmt.Errorf("insert forbidden")

}
