package postgres

import (
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	provider "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers"
	"github.com/jmoiron/sqlx"
)

type PostgresService interface {
	provider.ProviderService
}

type postgresService struct {
	postgresClient *sqlx.DB
}

func InitPostgresService(configs models.PostgresConfig) (PostgresService, error) {
	postgresClient, err := InitPostgresSqlService(configs)
	if err != nil {
		return nil, err
	}
	return &postgresService{postgresClient: postgresClient}, nil
}
