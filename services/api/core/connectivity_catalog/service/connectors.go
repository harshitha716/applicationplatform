package service

import (
	"context"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/google/uuid"
)

type ConnectorServiceStore interface {
	store.ConnectorStore
}

type ConnectorService interface {
	ListConnectors(ctx context.Context) ([]dbmodels.ConnectorWithActiveConnectionsCount, error)
	GetConnectorByID(ctx context.Context, id uuid.UUID) (*dbmodels.Connector, error)
}

type connectorService struct {
	store ConnectorServiceStore
}

func NewConnectorService(appStore store.Store) *connectorService {
	return &connectorService{
		store: appStore,
	}
}

func (s *connectorService) ListConnectors(ctx context.Context) ([]dbmodels.ConnectorWithActiveConnectionsCount, error) {
	connectors, error := s.store.GetAllConnectors(ctx)

	if error != nil {
		return nil, error
	}

	return connectors, nil
}

func (s *connectorService) GetConnectorByID(ctx context.Context, id uuid.UUID) (*dbmodels.Connector, error) {
	connector, err := s.store.GetConnectorById(ctx, id)

	if err != nil {
		return nil, err
	}

	return connector, nil
}

func (s *connectorService) CreateConnector(ctx context.Context, connector *dbmodels.Connector) error {
	return s.store.CreateConnector(ctx, connector)
}
