package auth

import (
	"net/http"
	"testing"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/helper"
	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_kratosIdentityToZampUser(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name     string
		identity *kratos.Identity
		want     *models.User
		wantErr  string
	}{
		{
			name: "valid identity with all fields",
			identity: &kratos.Identity{
				Id: validUUID.String(),
				Traits: map[string]interface{}{
					"email": "test@example.com",
					"name":  "Test User",
				},
			},
			want: &models.User{
				ID:    validUUID,
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: "",
		},
		{
			name: "valid identity without name",
			identity: &kratos.Identity{
				Id: validUUID.String(),
				Traits: map[string]interface{}{
					"email": "test@example.com",
				},
			},
			want: &models.User{
				ID:    validUUID,
				Email: "test@example.com",
				Name:  "",
			},
			wantErr: "",
		},
		{
			name: "invalid UUID",
			identity: &kratos.Identity{
				Id: "invalid-uuid",
				Traits: map[string]interface{}{
					"email": "test@example.com",
				},
			},
			want:    nil,
			wantErr: "invalid UUID length: 12",
		},
		{
			name: "nil traits",
			identity: &kratos.Identity{
				Id:     validUUID.String(),
				Traits: nil,
			},
			want:    nil,
			wantErr: "user traits not found",
		},
		{
			name: "invalid traits type",
			identity: &kratos.Identity{
				Id:     validUUID.String(),
				Traits: "invalid",
			},
			want:    nil,
			wantErr: "invalid user traits",
		},
		{
			name: "missing email",
			identity: &kratos.Identity{
				Id:     validUUID.String(),
				Traits: map[string]interface{}{},
			},
			want:    nil,
			wantErr: "email not found in user traits",
		},
		{
			name: "invalid email type",
			identity: &kratos.Identity{
				Id: validUUID.String(),
				Traits: map[string]interface{}{
					"email": 123,
				},
			},
			want:    nil,
			wantErr: "invalid email",
		},
		{
			name: "invalid name type",
			identity: &kratos.Identity{
				Id: validUUID.String(),
				Traits: map[string]interface{}{
					"email": "test@example.com",
					"name":  123,
				},
			},
			want: &models.User{
				ID:    validUUID,
				Email: "test@example.com",
				Name:  "",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kratosIdentityToZampUser(tt.identity)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectOrganization(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create test UUIDs
	orgId1, _ := uuid.Parse("be166699-eeea-4c8a-a3ec-107764dc3e91")
	orgId2, _ := uuid.Parse("ce166699-eeea-4c8a-a3ec-107764dc3e92")
	invalidOrgId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")

	tests := []struct {
		name            string
		header          http.Header
		organizationIds []uuid.UUID
		expected        []uuid.UUID
	}{
		{
			name:            "Header with valid organization ID",
			header:          http.Header{helper.PROXY_WORKSPACE_ID_HEADER: []string{orgId1.String()}},
			organizationIds: []uuid.UUID{orgId1, orgId2},
			expected:        []uuid.UUID{orgId1},
		},
		{
			name:            "Header with invalid organization ID",
			header:          http.Header{helper.PROXY_WORKSPACE_ID_HEADER: []string{invalidOrgId.String()}},
			organizationIds: []uuid.UUID{orgId1, orgId2},
			expected:        []uuid.UUID{orgId1},
		},
		{
			name:            "No header with multiple organizations",
			header:          http.Header{},
			organizationIds: []uuid.UUID{orgId1, orgId2},
			expected:        []uuid.UUID{orgId1},
		},
		{
			name:            "No header with no organizations",
			header:          http.Header{},
			organizationIds: []uuid.UUID{},
			expected:        []uuid.UUID{},
		},
		{
			name:            "Header with invalid UUID format",
			header:          http.Header{helper.PROXY_WORKSPACE_ID_HEADER: []string{"invalid-uuid"}},
			organizationIds: []uuid.UUID{orgId1, orgId2},
			expected:        []uuid.UUID{orgId1},
		},
		{
			name:            "Header with valid organization ID but user has no organizations",
			header:          http.Header{helper.PROXY_WORKSPACE_ID_HEADER: []string{orgId1.String()}},
			organizationIds: []uuid.UUID{},
			expected:        []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectOrganization(tt.header, tt.organizationIds, logger)
			assert.Equal(t, tt.expected, result)
		})
	}
}
