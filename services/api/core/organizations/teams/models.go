package teams

import (
	"github.com/google/uuid"
)

type CreateTeamPayload struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ColorHexCode string `json:"color_hex_code"`
}

type AddUserToTeamPayload struct {
	UserID uuid.UUID `json:"user_id"`
}

type RemoveUserFromTeamPayload struct {
	TeamMembershipID uuid.UUID `json:"team_membership_id"`
}

type RenameTeamPayload struct {
	Name string `json:"name"`
}
