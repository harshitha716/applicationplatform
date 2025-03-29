package databricks

import (
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	provider "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers"

	dbsdk "github.com/databricks/databricks-sdk-go"
	sqlx "github.com/jmoiron/sqlx"
)

type DatabricksService interface {
	provider.ProviderService
	DatabricksSDKProxy
}

type databricksService struct {
	db *sqlx.DB
	ws *dbsdk.WorkspaceClient
}

func InitDatabricksService(configs models.DatabricksConfig) (DatabricksService, error) {

	db, err := InitDatabricksSQLService(configs)
	if err != nil {
		return nil, err
	}

	ws, err := InitDatabricksSDKProxy(configs)
	if err != nil {
		return nil, err
	}

	return &databricksService{
		db: db,
		ws: ws,
	}, nil
}
