package organizations

type userMembershipState string

const (
	userMembershipStateInvited     userMembershipState = "invited"
	userMembershipStateUnderReview userMembershipState = "under_review"
	userMembershipStateNone        userMembershipState = "none"
)

type InvitationPayload struct {
	Privilege string `json:"role"`
	Email     string `json:"email"`
}

type BulkInvitationPayload struct {
	Invitations []InvitationPayload `json:"invitations"`
}

type InvitationError struct {
	Email        string `json:"email"`
	ErrorMessage string `json:"error_message"`
}

type BulkInvitationError struct {
	Error       error             `json:"error"`
	Invitations []InvitationError `json:"invitations"`
}
