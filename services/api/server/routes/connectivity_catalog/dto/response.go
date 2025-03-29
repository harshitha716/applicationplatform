package dto

import (
	"encoding/json"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
)

type Connector struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	DisplayName            string `json:"display_name"`
	LogoURL                string `json:"logo_url"`
	Category               string `json:"category"`
	Status                 string `json:"status"`
	ActiveConnectionsCount int    `json:"active_connections_count"`
}

type ConnectorDetails struct {
	Connector
	Documentation  string                 `json:"documentation"`
	ConfigTemplate map[string]interface{} `json:"config_template"` // JSON schema for connector configuration
}

func (c *Connector) FromModel(connector dbmodels.ConnectorWithActiveConnectionsCount) {
	c.ID = connector.ID.String()
	c.Name = connector.Name
	c.Description = connector.Description
	c.DisplayName = connector.DisplayName
	c.LogoURL = connector.LogoURL
	c.Category = connector.Category
	c.Status = connector.Status

	c.ActiveConnectionsCount = connector.ActiveConnectionsCount
}

func (c *ConnectorDetails) FromModel(connector dbmodels.Connector) {
	config := map[string]interface{}{}

	json.Unmarshal(connector.ConfigTemplate, &config)

	c.Connector = Connector{
		ID:          connector.ID.String(),
		Name:        connector.Name,
		Description: connector.Description,
		DisplayName: connector.DisplayName,
		LogoURL:     connector.LogoURL,
		Category:    connector.Category,
		Status:      connector.Status,
	}
	c.Documentation = connector.Documentation
	c.ConfigTemplate = config
}
