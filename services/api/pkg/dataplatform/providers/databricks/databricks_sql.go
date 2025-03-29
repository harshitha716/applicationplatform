package databricks

import (
	"context"
	"database/sql"
	"time"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/logger"

	dbsql "github.com/databricks/databricks-sql-go"
	sqlx "github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DatabricksSQLService interface {
	Query(ctx context.Context, query string, args ...interface{}) (models.QueryResult, error)
}

func InitDatabricksSQLService(configs models.DatabricksConfig) (*sqlx.DB, error) {
	connector, err := dbsql.NewConnector(
		dbsql.WithServerHostname(configs.ServerHostname),
		dbsql.WithHTTPPath(configs.HttpPath),
		dbsql.WithPort(configs.Port),
		dbsql.WithAccessToken(configs.AccessToken),
		dbsql.WithTimeout(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	db := sqlx.NewDb(sql.OpenDB(connector), constants.DATABRICKS_DRIVER_NAME)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *databricksService) Query(ctx context.Context, table string, query string, args ...interface{}) (models.QueryResult, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	startTime := time.Now()
	logger.Info("QUERYING_DATABRICKS", zap.String("QUERY", query))
	rows, err := db.db.Queryx(query, args...)
	if err != nil {
		logger.Error(errors.QueryingDatabricksFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, errors.ErrQueryingDatabricks
	}

	defer rows.Close()

	queryResult := models.QueryResult{}
	err = queryResult.FromSqlRows(rows)
	if err != nil {
		logger.Error(errors.BuildingQueryFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, errors.ErrBuildingQueryResult
	}
	logger.Info("SUCCESSFULLY_QUERIED_DATABRICKS", zap.Any("DATABRICKS_QUERY_TIME_MS", time.Since(startTime).Milliseconds()))
	return queryResult, nil
}
