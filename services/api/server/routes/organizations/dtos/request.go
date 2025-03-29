package dtos

import "github.com/google/uuid"

type UpdateMemberRoleRequest struct {
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type InviteMemberPayload struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}

type BulkInvitationPayload struct {
	Invitations []InviteMemberPayload `json:"invitations"`
}

type DeleteMemberRequest struct {
	UserId uuid.UUID `json:"user_id"`
}

type ApproveOrganizationMembershipRequestRequest struct {
	UserId uuid.UUID `json:"user_id"`
}
