package models

import (
	"encoding/json"
	"time"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type AuditLogEmitParams struct {
	Kind           dbmodels.AuditLogKind
	ResourceType   dbmodels.ResourceType
	ResourceID     uuid.UUID
	EventName      string
	Payload        map[string]interface{}
	IPAddress      string
	UserEmail      string
	UserAgent      string
	OrganizationID uuid.UUID
}

type AuditLog struct {
	ID             uuid.UUID
	Kind           string
	OrganizationID uuid.UUID
	IPAddress      string
	UserEmail      string
	UserAgent      string
	ResourceType   dbmodels.ResourceType
	ResourceID     uuid.UUID
	EventName      string
	Payload        interface{}
	CreatedAt      time.Time
}

func (a *AuditLog) FromSchema(schema dbmodels.AuditLog) {
	a.ID = schema.ID
	a.Kind = string(schema.Kind)
	a.OrganizationID = schema.OrganizationID
	a.IPAddress = schema.IPAddress
	a.UserEmail = schema.UserEmail
	a.UserAgent = schema.UserAgent
	a.ResourceType = schema.ResourceType
	a.ResourceID = schema.ResourceID
	a.EventName = schema.EventName
	var payload interface{}
	_ = json.Unmarshal(schema.Payload, &payload)
	a.Payload = payload
	a.CreatedAt = schema.CreatedAt
}

func (a *AuditLog) ToSchema() dbmodels.AuditLog {
	payloadBytes, _ := json.Marshal(a.Payload)
	return dbmodels.AuditLog{
		ID:             a.ID,
		Kind:           dbmodels.AuditLogKind(a.Kind),
		OrganizationID: a.OrganizationID,
		IPAddress:      a.IPAddress,
		UserEmail:      a.UserEmail,
		UserAgent:      a.UserAgent,
		ResourceType:   a.ResourceType,
		ResourceID:     a.ResourceID,
		EventName:      a.EventName,
		Payload:        payloadBytes,
		CreatedAt:      a.CreatedAt,
	}
}
