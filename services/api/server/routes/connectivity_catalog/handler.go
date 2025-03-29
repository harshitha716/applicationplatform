package connectivity_catalog

import (
	"github.com/Zampfi/application-platform/services/api/core/connectivity_catalog/service"
	"github.com/Zampfi/application-platform/services/api/server/routes/connectivity_catalog/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleGetConnectors(c *gin.Context, connectorService service.ConnectorService) {
	connectors, err := connectorService.ListConnectors(c)
	connectorsDto := []dto.Connector{}
	for _, connector := range connectors {
		connectorDto := dto.Connector{}
		connectorDto.FromModel(connector)
		connectorsDto = append(connectorsDto, connectorDto)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, connectorsDto)
}

func handleGetConnectorByID(c *gin.Context, connectorService service.ConnectorService) {
	connectorIDStr := c.Param("connectorId")
	connectorID, err := uuid.Parse(connectorIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid connector id"})
		return
	}

	connector, err := connectorService.GetConnectorByID(c, connectorID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	connectorDetailsDto := dto.ConnectorDetails{}

	connectorDetailsDto.FromModel(*connector)

	c.JSON(200, connectorDetailsDto)
}
