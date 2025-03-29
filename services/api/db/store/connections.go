package store

import (
	"context"
	"fmt"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

type ConnectionStore interface {
	CreateConnection(ctx context.Context, connection *models.CreateConnectionParams) (uuid.UUID, error)
	GetConnectionByID(ctx context.Context, id uuid.UUID) (*models.Connection, error)
	GetConnections(ctx context.Context) ([]models.Connection, error)
}

type connectionStore struct {
	db *pgclient.PostgresClient
}

func NewConnectionStore(db *pgclient.PostgresClient) *connectionStore {
	return &connectionStore{db: db}
}

func (s *appStore) CreateConnection(ctx context.Context, connection *models.CreateConnectionParams) (uuid.UUID, error) {
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)

	if len(orgIds) != 1 {
		return uuid.Nil, fmt.Errorf("organization access forbidden")
	}

	connectionId := uuid.New()

	cnx := models.Connection{
		ID:             connectionId,
		ConnectorID:    connection.ConnectorID,
		OrganizationID: orgIds[0],
		Name:           connection.Name,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result := s.client.WithContext(ctx).Create(&cnx)

	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return connectionId, nil
}

func (s *appStore) GetConnectionByID(ctx context.Context, id uuid.UUID) (*models.Connection, error) {
	db := s.client.WithContext(ctx)

	connection := models.Connection{}

	result := db.First(&connection, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &connection, nil

}

func (s *appStore) GetConnections(ctx context.Context) ([]models.Connection, error) {
	_, _, orgIds := apicontext.GetAuthFromContext(ctx)

	if len(orgIds) != 1 {
		return nil, fmt.Errorf("organization access forbidden")
	}

	db := s.client.WithContext(ctx)

	connections := []models.Connection{}

	err := db.Preload("Connector").Preload("Schedules").Order("created_at ASC").Find(&connections, "organization_id = ?", orgIds[0])
	if err.Error != nil {
		return nil, err.Error
	}

	return connections, nil

}
