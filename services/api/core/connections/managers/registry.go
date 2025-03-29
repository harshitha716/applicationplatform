package manager

import (
	"fmt"

	"github.com/Zampfi/application-platform/services/api/core/connections/constants"
	"github.com/Zampfi/application-platform/services/api/core/connections/managers/gcs"
	"github.com/Zampfi/application-platform/services/api/core/connections/managers/snowflake"
)

type ConnectionManagerRegistry struct {
	managers map[string]ConnectionManager
}

func NewConnectionManagerRegistry() ConnectionManagerRegistry {
	registry := ConnectionManagerRegistry{
		managers: make(map[string]ConnectionManager),
	}

	// Register connection managers
	registry.managers[constants.GCS] = gcs.InitGCSConnectionManager()
	registry.managers[constants.Snowflake] = snowflake.InitSnowflakeConnectionManager()

	return registry
}

func (r *ConnectionManagerRegistry) GetManager(connectorType string) (ConnectionManager, error) {
	manager, exists := r.managers[connectorType]
	if !exists {
		return nil, fmt.Errorf("no connection manager found for connector type: %s", connectorType)
	}
	return manager, nil
}
