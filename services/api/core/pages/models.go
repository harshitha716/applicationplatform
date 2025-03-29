package pages

import (
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type AddPageAudiencePayload struct {
	AudienceType models.AudienceType
	AudienceId   uuid.UUID
	Privilege    models.ResourcePrivilege
}

type BulkAddPageAudiencePayload struct {
	Audiences []AddPageAudiencePayload
}

type AddPageAudienceError struct {
	AudienceId   uuid.UUID `json:"audience_id"`
	ErrorMessage string    `json:"error_message"`
}

type BulkAddPageAudienceErrors struct {
	Error     error
	Audiences []AddPageAudienceError `json:"audiences"`
}

type CreatePagePayload struct {
	PageName        string `form:"label=Page Name"`
	PageDescription string `form:"label=Page Description"`
}
