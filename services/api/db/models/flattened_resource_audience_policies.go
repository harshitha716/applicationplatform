package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlattenedResourceAudiencePolicy struct {
	ResourceAudiencePolicyId uuid.UUID         `json:"resource_audience_policy_id"`
	ResourceAudienceType     string            `json:"resource_audience_type"`
	UserId                   uuid.UUID         `json:"user_id"`
	ResourceId               uuid.UUID         `json:"resource_id"`
	ResourceAudienceId       uuid.UUID         `json:"resource_audience_id"`
	ResourceType             string            `json:"resource_type"`
	Privilege                ResourcePrivilege `json:"privilege"`
	CreatedAt                string            `json:"created_at"`
	UpdatedAt                string            `json:"updated_at"`
	DeletedAt                string            `json:"deleted_at"`
}

type FlattenedResourceAudiencePoliciesFilters struct {
	ResourceIds   []uuid.UUID
	UserIds       []uuid.UUID
	ResourceTypes []string
	Privileges    []ResourcePrivilege
}

func (f *FlattenedResourceAudiencePolicy) TableName() string {
	return "flattened_resource_audience_policies"
}

func (f *FlattenedResourceAudiencePolicy) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`
		EXISTS (
			SELECT 1
			FROM app.flattened_resource_audience_policies nfrap
			WHERE 
				nfrap.resource_id = flattened_resource_audience_policies.resource_id
				AND nfrap.resource_type = flattened_resource_audience_policies.resource_type
				AND nfrap.user_id = ?
		) OR flattened_resource_audience_policies.user_id = ?
	`, userId, userId)
}

func (f *FlattenedResourceAudiencePolicy) BeforeCreate(db *gorm.DB) error {

	return fmt.Errorf("insert forbidden")

}
