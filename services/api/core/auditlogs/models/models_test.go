package models

import (
	"encoding/json"
	"testing"
	"time"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuditLogFromSchema(t *testing.T) {
	// Create a test DB model
	id := uuid.New()
	orgId := uuid.New()
	now := time.Now().UTC()
	payload := map[string]interface{}{"test": "value"}
	payloadBytes, _ := json.Marshal(payload)

	dbAuditLog := dbmodels.AuditLog{
		ID:             id,
		Kind:           dbmodels.AuditLogKindInfo,
		OrganizationID: orgId,
		IPAddress:      "192.168.1.1",
		UserEmail:      "test@example.com",
		UserAgent:      "Mozilla/5.0",
		ResourceType:   dbmodels.ResourceTypeOrganization,
		ResourceID:     orgId,
		EventName:      "test_event",
		Payload:        payloadBytes,
		CreatedAt:      now,
	}

	// Convert to API model
	var auditLog AuditLog
	auditLog.FromSchema(dbAuditLog)

	// Verify fields
	assert.Equal(t, id, auditLog.ID)
	assert.Equal(t, string(dbmodels.AuditLogKindInfo), auditLog.Kind)
	assert.Equal(t, orgId, auditLog.OrganizationID)
	assert.Equal(t, "192.168.1.1", auditLog.IPAddress)
	assert.Equal(t, "test@example.com", auditLog.UserEmail)
	assert.Equal(t, "Mozilla/5.0", auditLog.UserAgent)
	assert.Equal(t, dbmodels.ResourceTypeOrganization, auditLog.ResourceType)
	assert.Equal(t, orgId, auditLog.ResourceID)
	assert.Equal(t, "test_event", auditLog.EventName)
	assert.Equal(t, now, auditLog.CreatedAt)

	// Check payload
	payloadMap, ok := auditLog.Payload.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", payloadMap["test"])
}

func TestAuditLogToSchema(t *testing.T) {
	// Create a test API model
	id := uuid.New()
	orgId := uuid.New()
	now := time.Now().UTC()

	auditLog := AuditLog{
		ID:             id,
		Kind:           string(dbmodels.AuditLogKindInfo),
		OrganizationID: orgId,
		IPAddress:      "192.168.1.1",
		UserEmail:      "test@example.com",
		UserAgent:      "Mozilla/5.0",
		ResourceType:   dbmodels.ResourceTypeOrganization,
		ResourceID:     orgId,
		EventName:      "test_event",
		Payload:        map[string]interface{}{"test": "value"},
		CreatedAt:      now,
	}

	// Convert to DB model
	dbAuditLog := auditLog.ToSchema()

	// Verify fields
	assert.Equal(t, id, dbAuditLog.ID)
	assert.Equal(t, dbmodels.AuditLogKindInfo, dbAuditLog.Kind)
	assert.Equal(t, orgId, dbAuditLog.OrganizationID)
	assert.Equal(t, "192.168.1.1", dbAuditLog.IPAddress)
	assert.Equal(t, "test@example.com", dbAuditLog.UserEmail)
	assert.Equal(t, "Mozilla/5.0", dbAuditLog.UserAgent)
	assert.Equal(t, dbmodels.ResourceTypeOrganization, dbAuditLog.ResourceType)
	assert.Equal(t, orgId, dbAuditLog.ResourceID)
	assert.Equal(t, "test_event", dbAuditLog.EventName)
	assert.Equal(t, now, dbAuditLog.CreatedAt)

	// Check payload
	var payload map[string]interface{}
	err := json.Unmarshal(dbAuditLog.Payload, &payload)
	assert.NoError(t, err)
	assert.Equal(t, "value", payload["test"])
}

func TestAuditLogEmitParams(t *testing.T) {
	// Test that the struct exists and can be populated
	orgId := uuid.New()

	params := AuditLogEmitParams{
		Kind:           dbmodels.AuditLogKindInfo,
		ResourceType:   dbmodels.ResourceTypeOrganization,
		ResourceID:     orgId,
		EventName:      "test_event",
		Payload:        map[string]interface{}{"test": "value"},
		IPAddress:      "192.168.1.1",
		UserEmail:      "test@example.com",
		UserAgent:      "Mozilla/5.0",
		OrganizationID: orgId,
	}

	// Verify fields
	assert.Equal(t, dbmodels.AuditLogKindInfo, params.Kind)
	assert.Equal(t, dbmodels.ResourceTypeOrganization, params.ResourceType)
	assert.Equal(t, orgId, params.ResourceID)
	assert.Equal(t, "test_event", params.EventName)
	assert.Equal(t, "192.168.1.1", params.IPAddress)
	assert.Equal(t, "test@example.com", params.UserEmail)
	assert.Equal(t, "Mozilla/5.0", params.UserAgent)
	assert.Equal(t, orgId, params.OrganizationID)
	assert.Equal(t, "value", params.Payload["test"])
}
