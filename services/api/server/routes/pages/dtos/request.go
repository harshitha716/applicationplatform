package dtos

import "github.com/google/uuid"

type UpdateAudienceRoleRequest struct {
	AudiencId uuid.UUID `json:"audience_id"`
	Role      string    `json:"role"`
}

type DeleteAudienceRoleRequest struct {
	AudiencId uuid.UUID `json:"audience_id"`
}

type AddAudienceRequest struct {
	AudienceType string    `json:"audience_type"`
	AudienceId   uuid.UUID `json:"audience_id"`
	Role         string    `json:"role"`
}

type BulkAddAudienceRequest struct {
	Audiences []AddAudienceRequest `json:"audiences"`
}
