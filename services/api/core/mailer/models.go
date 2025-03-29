package mailer

type InvitationEmailData struct {
	RecipientEmail     string
	InvitedByFirstName string
	OrganizationName   string
	InvitationLink     string
}
