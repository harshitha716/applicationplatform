package service

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConnectionServiceStore interface {
	store.ConnectionStore
	store.ConnectionPoliciesStore
}

type ConnectionService interface {
	CreateConnection(ctx context.Context, params models.CreateConnectionParams, tx ConnectionServiceStore) (uuid.UUID, error)
	GetConnections(ctx context.Context, tx ConnectionServiceStore) ([]models.Connection, error)
	GetConnectionByID(ctx context.Context, connectionId uuid.UUID) (models.Connection, error)
}

type connectionService struct {
	store ConnectionServiceStore
}

func NewConnectionService(appStore store.Store) *connectionService {
	return &connectionService{
		store: appStore,
	}
}

func (s *connectionService) CreateConnection(ctx context.Context, params models.CreateConnectionParams, tx ConnectionServiceStore) (uuid.UUID, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)
	var connectionId uuid.UUID

	if tx == nil {
		tx = s.store
	}

	_, userID, _ := apicontext.GetAuthFromContext(ctx)

	var err error
	connectionId, err = tx.CreateConnection(ctx, &params)
	if err != nil {
		logger.Error("Failed to create connection", zap.Error(err))
		return uuid.Nil, err
	}

	_, err = tx.CreateConnectionPolicy(
		ctx,
		connectionId,
		models.AudienceTypeUser,
		*userID,
		models.PrivilegeConnectionAdmin,
	)

	if err != nil {
		logger.Error("Failed to create connection policy", zap.Error(err))
		return uuid.Nil, err
	}

	return connectionId, nil
}

func (s *connectionService) GetConnections(ctx context.Context, tx ConnectionServiceStore) ([]models.Connection, error) {
	connections, err := tx.GetConnections(ctx)

	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (s *connectionService) GetConnectionByID(ctx context.Context, connectionId uuid.UUID) (models.Connection, error) {
	connection, err := s.store.GetConnectionByID(ctx, connectionId)
	if err != nil {
		return models.Connection{}, err
	}
	return *connection, nil
}
