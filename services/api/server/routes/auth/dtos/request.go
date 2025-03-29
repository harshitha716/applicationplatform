package dtos

type GetLoginFlowRequest struct {
	Email string `json:"email"`
}

type AfterRegistrationWebhookRequest struct {
	UserId string `json:"user_id"`
}
