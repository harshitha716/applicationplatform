package postgres

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/logger"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgresSqlService interface {
	Query(ctx context.Context, table string, query string, args ...interface{}) (models.QueryResult, error)
}

func InitPostgresSqlService(configs models.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(constants.POSTGRES_DRIVER_NAME, configs.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.ErrPostgresServiceInitializationFailed
	}
	return db, nil
}

func (db *postgresService) Query(ctx context.Context, table string, query string, args ...interface{}) (models.QueryResult, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	rows, err := db.postgresClient.Queryx(query, args...)
	if err != nil {
		logger.Error(errors.QueryingPostgresFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, errors.ErrQueryingPostgres
	}

	defer rows.Close()

	queryResult := models.QueryResult{}
	err = queryResult.FromSqlRows(rows)
	if err != nil {
		logger.Error(errors.BuildingQueryFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, errors.ErrBuildingQueryResult
	}
	return queryResult, nil
}
