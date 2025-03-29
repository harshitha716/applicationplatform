package dto

type CreateConnectionRequest struct {
	DisplayName      string                 `json:"display_name" binding:"required"`
	ConnectorName    string                 `json:"connector_name" binding:"required"`
	ConnectorID      string                 `json:"connector_id" binding:"required"`
	ConnectionConfig map[string]interface{} `json:"connection_config" binding:"required"`
}

type GetSchedulesRequest struct {
	ConnectionID string `json:"connection_id" binding:"required"`
}
