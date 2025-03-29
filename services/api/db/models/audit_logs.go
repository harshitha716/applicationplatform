package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogKind string

const (
	AuditLogKindInfo     AuditLogKind = "info"
	AuditLogKindInternal AuditLogKind = "internal"
)

type AuditLog struct {
	ID             uuid.UUID       `json:"audit_log_id" gorm:"column:audit_log_id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Kind           AuditLogKind    `json:"kind"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	IPAddress      string          `json:"user_ip_address"`
	UserEmail      string          `json:"user_email"`
	UserAgent      string          `json:"user_agent"`
	ResourceType   ResourceType    `json:"resource_type"`
	ResourceID     uuid.UUID       `json:"resource_id"`
	EventName      string          `json:"event_name"`
	Payload        json.RawMessage `json:"payload"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (a *AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate hook to validate the audit log
func (a *AuditLog) BeforeCreate(db *gorm.DB) error {
	_, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
	if userId == nil {
		return errors.New("no user id found in context")
	}

	fraps := []FlattenedResourceAudiencePolicy{}
	db.Where("resource_type = ? AND resource_id = ? AND user_id = ?", ResourceTypeOrganization, a.OrganizationID, userId).Find(&fraps)
	if len(fraps) == 0 {
		return errors.New("user does not have access to the resource")
	}

	return nil
}

// BeforeUpdate hook to prevent updates to audit logs
func (a *AuditLog) BeforeUpdate(tx *gorm.DB) error {
	return fmt.Errorf("forbidden: audit logs cannot be updated")
}

// BeforeDelete hook to prevent deletion of audit logs
func (a *AuditLog) BeforeDelete(tx *gorm.DB) error {
	return errors.New("forbidden: audit logs cannot be deleted")
}

// Add query filters for access control
func (a *AuditLog) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`
		EXISTS (
			SELECT 1
			FROM app.flattened_resource_audience_policies frap
			WHERE
				frap.resource_type = audit_logs.resource_type AND
				frap.resource_id = audit_logs.resource_id AND
				frap.user_id = ? AND
				frap.deleted_at IS NULL
		)
	`, userId)
}
