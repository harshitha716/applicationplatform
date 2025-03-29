package store

import (
	"context"
	"errors"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

type ConnectorStore interface {
	GetAllConnectors(ctx context.Context) ([]models.ConnectorWithActiveConnectionsCount, error)
	GetConnectorById(ctx context.Context, id uuid.UUID) (*models.Connector, error)
	CreateConnector(ctx context.Context, connector *models.Connector) error
}

type connectorStore struct {
	db *pgclient.PostgresClient
}

func NewConnectorStore(db *pgclient.PostgresClient) *connectorStore {
	return &connectorStore{db: db}
}

func (s *appStore) GetAllConnectors(ctx context.Context) ([]models.ConnectorWithActiveConnectionsCount, error) {
	db := s.client.WithContext(ctx)

	_, _, orgIds := apicontext.GetAuthFromContext(ctx)
	var results []models.ConnectorWithActiveConnectionsCount

	if len(orgIds) != 1 {
		return nil, errors.New("organization access forbidden")
	}

	result := db.Model(&models.Connector{}).
		Select("connectors.*, COUNT(connections.id) AS active_connections_count").
		Joins("LEFT JOIN connections ON connectors.id = connections.connector_id AND connections.status = 'active' AND connections.organization_id = ?", orgIds[0]).
		Group("connectors.id").
		Order("connectors.name ASC").
		Scan(&results)

	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil

}

func (s *appStore) GetConnectorById(ctx context.Context, id uuid.UUID) (*models.Connector, error) {
	var connector models.Connector
	connector.ID = id

	err := s.client.WithContext(ctx).Model(connector).First(&connector).Error
	if err != nil {
		return nil, err
	}

	return &connector, nil

}

func (s *appStore) CreateConnector(ctx context.Context, connector *models.Connector) error {
	return s.client.WithContext(ctx).Create(connector).Error
}
