package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceType string
type AudienceType string
type ResourcePrivilege string
type ConnectionPrivilege string
type SchedulePrivilege string

const (
	PrivilegeOrganizationMember      ResourcePrivilege = "member"
	PrivilegeOrganizationSystemAdmin ResourcePrivilege = "system_admin"
	PrivilegeDatasetAdmin            ResourcePrivilege = "admin"
	PrivilegeDatasetViewer           ResourcePrivilege = "viewer"
	PrivilegePageAdmin               ResourcePrivilege = "admin"
	PrivilegePageRead                ResourcePrivilege = "viewer"
	PrivilegeConnectionAdmin         ResourcePrivilege = "admin"
	PrivilegeConnectionRead          ResourcePrivilege = "viewer"
	PrivilegeScheduleAdmin           ResourcePrivilege = "admin"
	PrivilegeScheduleRead            ResourcePrivilege = "viewer"
	PrivilegePaymentsAdmin           ResourcePrivilege = "admin"
	PrivilegePaymentsViewer          ResourcePrivilege = "viewer"
	PrivilegePaymentsInitiator       ResourcePrivilege = "initiator"
)

const (
	ResourceTypeOrganization ResourceType = "organization"
	ResourceTypeDataset      ResourceType = "dataset"
	ResourceTypePage         ResourceType = "page"
	ResourceTypeConnection   ResourceType = "connection"
	ResourceTypeSchedule     ResourceType = "schedule"
	ResourceTypePayments     ResourceType = "payments"
)

var OrganizationPrivileges = []ResourcePrivilege{PrivilegeOrganizationMember, PrivilegeOrganizationSystemAdmin}
var DatasetPrivileges = []ResourcePrivilege{PrivilegeDatasetAdmin, PrivilegeDatasetViewer}
var PagePrivileges = []ResourcePrivilege{PrivilegePageAdmin, PrivilegePageRead}
var ConnectionPrivileges = []ResourcePrivilege{PrivilegeConnectionAdmin, PrivilegeConnectionRead}
var SchedulePrivileges = []ResourcePrivilege{PrivilegeScheduleAdmin, PrivilegeScheduleRead}
var PaymentsPrivileges = []ResourcePrivilege{PrivilegePaymentsAdmin, PrivilegePaymentsViewer, PrivilegePaymentsInitiator}

const (
	AudienceTypeUser         AudienceType = "user"
	AudienceTypeOrganization AudienceType = "organization"
	AudienceTypeTeam         AudienceType = "team"
)

type ResourceAudiencePolicy struct {
	ID                   uuid.UUID                         `json:"resource_audience_policy_id" gorm:"column:resource_audience_policy_id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ResourceAudienceType AudienceType                      `json:"resource_audience_type"`
	ResourceAudienceID   uuid.UUID                         `json:"resource_audience_id"`
	Privilege            ResourcePrivilege                 `json:"privilege"`
	ResourceType         ResourceType                      `json:"resource_type"`
	User                 *User                             `json:"user,omitempty" gorm:"foreignKey:ResourceAudienceID;references:ID"`
	ResourceID           uuid.UUID                         `json:"resource_id"`
	CreatedAt            time.Time                         `json:"created_at"`
	UpdatedAt            time.Time                         `json:"updated_at"`
	DeletedAt            *time.Time                        `json:"deleted_at"`
	Metadata             *json.RawMessage                  `json:"metadata"`
	UserPolicies         []FlattenedResourceAudiencePolicy `json:"-" gorm:"foreignKey:ResourceAudiencePolicyId;references:ID"`
}

func (r *ResourceAudiencePolicy) TableName() string {
	return "resource_audience_policies"
}

func (r *ResourceAudiencePolicy) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`
		EXISTS (
			SELECT 1
			FROM app.flattened_resource_audience_policies frap
			WHERE
				frap.resource_type = resource_audience_policies.resource_type
				AND frap.resource_id = resource_audience_policies.resource_id
				AND frap.user_id = ?
				AND frap.deleted_at IS NULL
		)
	`, userId)
}

func (r *ResourceAudiencePolicy) BeforeCreate(db *gorm.DB) error {
	return nil
}

func (r *ResourceAudiencePolicy) BeforeUpdate(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where(`
		(resource_id = ? ANd resource_type = ? AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = ? AND deleted_at IS NULL)
	`, r.ResourceID, r.ResourceType, userId).Limit(1).Find(&fraps).Error

	if err != nil {
		return fmt.Errorf("error checking user's access to the resource")
	}

	if len(fraps) == 0 {
		return fmt.Errorf("user does not have permission to update the policy on resource")
	}

	return nil
}

func (r *ResourceAudiencePolicy) BeforeDelete(db *gorm.DB) error {

	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)

	fraps := []FlattenedResourceAudiencePolicy{}
	err := db.Where(`
		(resource_id = ? ANd resource_type = ? AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = ? AND deleted_at IS NULL)
	`, r.ResourceID, r.ResourceType, userId).Limit(1).Find(&fraps).Error

	if err != nil {
		return fmt.Errorf("error checking user's access to the resource")
	}

	if len(fraps) == 0 {
		return fmt.Errorf("user does not have permission to update the policy on resource")
	}

	return nil
}
