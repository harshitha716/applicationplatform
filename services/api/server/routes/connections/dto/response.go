package dto

type CreateConnectionResponse struct {
	ConnectionID string `json:"connection_id"`
}

type ConnectionDetails struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ConnectorIconURL string `json:"icon_url"`
	LastSyncedAt     string `json:"last_synced_at"`
	CreatedAt        string `json:"created_at"`
}

type ScheduleDetails struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	NextRunAt string `json:"next_run_at"`
	LastRunAt string `json:"last_run_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	LogoURL   string `json:"logo_url"`
}
