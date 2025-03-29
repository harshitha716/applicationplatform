package models

import (
	"encoding/json"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationSSOConfig struct {
	OrganizationSSOConfigID uuid.UUID       `json:"organization_sso_config_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID          uuid.UUID       `json:"organization_id"`
	SSOProviderID           string          `json:"sso_provider_id"`
	SSOProviderName         string          `json:"sso_provider_name"`
	SSOConfig               json.RawMessage `json:"sso_config"`
	EmailDomain             string          `json:"email_domain"`
	IsPrimary               bool            `json:"is_primary"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

func (o *OrganizationSSOConfig) TableName() string {
	return "organization_sso_configs"
}

func (o *OrganizationSSOConfig) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {

	ctx := db.Statement.Context

	role, _, _ := apicontext.GetAuthFromContext(ctx)

	if role == "admin" || role == "anonymous" {
		return db
	}

	return db.Where(`
		EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = ?
			AND frap.resource_id in ?
			AND frap.user_id = ?
			AND frap.deleted_at IS NULL
		)
	`, ResourceTypeOrganization, orgIds, userId,
	)

}

func (o *OrganizationSSOConfig) BeforeCreate(db *gorm.DB) error {

	role, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)

	if userId == nil {
		return fmt.Errorf("no user ID in context")
	}

	if role == "admin" {
		return nil
	}

	frap := []FlattenedResourceAudiencePolicy{}

	err := db.Model(&FlattenedResourceAudiencePolicy{}).Where("resource_type = ? AND resource_id = ? AND user_id = ?", ResourceTypeOrganization, o.OrganizationID, userId).Find(&frap).Error
	if err != nil {
		return fmt.Errorf("failed to check if user has access to organization: %w", err)
	}

	if len(frap) == 0 {
		return fmt.Errorf("user does not have access to organization")
	}

	return nil
}

func (o *OrganizationSSOConfig) BeforeUpdate(db *gorm.DB) error {

	role, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)

	if role != "admin" {
		return fmt.Errorf("only admin can update organization sso configs")
	}

	frap := []FlattenedResourceAudiencePolicy{}

	err := db.Model(&FlattenedResourceAudiencePolicy{}).Where("resource_type = ? AND resource_id = ? AND user_id = ?", ResourceTypeOrganization, o.OrganizationID, userId).Find(&frap).Error
	if err != nil {
		return fmt.Errorf("failed to check if user has access to organization: %w", err)
	}

	if len(frap) == 0 {
		return fmt.Errorf("user does not have access to organization")
	}

	return nil
}
