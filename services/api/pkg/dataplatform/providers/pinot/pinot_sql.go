package pinot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Zampfi/application-platform/services/api/pkg/errorreporting"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/logger"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"

	pinot "github.com/startreedata/pinot-client-go/pinot"
	"go.uber.org/zap"
)

type PinotSQLService interface {
	Query(ctx context.Context, query string, args ...interface{}) (models.QueryResult, error)
}

func InitPinotSqlService(configs models.PinotConfig) (*pinot.Connection, error) {
	pinotClient, err := pinot.NewWithConfig(&pinot.ClientConfig{
		BrokerList: configs.BrokerList,
		ExtraHTTPHeader: map[string]string{
			"Authorization": "Bearer " + configs.AccessToken,
		},
		HTTPTimeout: 60 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	pinotClient.UseMultistageEngine(true)
	return pinotClient, nil
}

func (p *pinotService) Query(ctx context.Context, table string, query string, args ...interface{}) (models.QueryResult, error) {
	startTime := time.Now()
	logger := logger.GetLoggerFromCtx(ctx)
	logger.Info("QUERYING_PINOT", zap.String("QUERY", query))

	sqlResponse, err := p.pinotClient.ExecuteSQL(table, query)
	if err != nil {
		logger.Error(errors.QueryingPinotFailedErrMessage, zap.Error(err))
		// TODO, FIXME: pkg should not be calling any IO directly (e.g. errorreporting)
		errorreporting.CaptureException(fmt.Errorf("error querying Pinot: %w", err), ctx)
		return models.QueryResult{}, errors.ErrQueryingPinot
	}

	queryResult, err := convertToQueryResult(ctx, sqlResponse)
	if err != nil {
		logger.Error(errors.PinotQueryResultConversionFailedErrMessage, zap.Error(err))
		return models.QueryResult{}, err
	}

	logger.Info("SUCCESSFULLY_QUERIED_PINOT", zap.Any("PINOT_QUERY_TIME_MS", time.Since(startTime).Milliseconds()))
	return queryResult, nil
}

func convertToQueryResult(ctx context.Context, sqlResponse *pinot.BrokerResponse) (models.QueryResult, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	queryResult := models.QueryResult{}
	queryResult.Rows = models.Rows{}
	queryResult.Columns = []models.ColumnMetadata{}

	if sqlResponse.Exceptions != nil && len(sqlResponse.Exceptions) > 0 {
		jsonBytes, err := json.Marshal(sqlResponse.Exceptions)
		if err != nil {
			logger.Error(errors.PinotQueryExceptionsErrMessage, zap.Error(err))
			// TODO, FIXME: pkg should not be calling any IO directly (e.g. errorreporting)
			errorreporting.CaptureException(fmt.Errorf("error querying Pinot: %w", err), ctx, errorreporting.Tag{Key: "category", Value: "data_platform"})
			return models.QueryResult{}, errors.ErrPinotQueryExceptions
		}
		logger.Error(errors.PinotQueryExceptionsErrMessage, zap.Any("PINOT_SQL_RESPONSE_EXCEPTIONS", string(jsonBytes)))
		errorreporting.CaptureException(fmt.Errorf("error querying Pinot: %s", jsonBytes), ctx, errorreporting.Tag{Key: "category", Value: "data_platform"})
		return models.QueryResult{}, errors.ErrPinotQueryExceptions
	}

	if sqlResponse.AggregationResults != nil {
		jsonBytes, err := json.Marshal(sqlResponse.AggregationResults)
		if err == nil {
			logger.Info("PINOT_SQL_RESPONSE_AGGREGATION_RESULTS", zap.Any("AGGREGATION_RESULTS", string(jsonBytes)))
		}
	}

	if sqlResponse.SelectionResults != nil {
		jsonBytes, err := json.Marshal(sqlResponse.SelectionResults)
		if err == nil {
			logger.Info("PINOT_SQL_RESPONSE_SELECTION_RESULTS", zap.Any("SELECTION_RESULTS", string(jsonBytes)))
		}
	}

	if sqlResponse.ResultTable == nil {
		logger.Error(errors.PinotQueryResultTableNotFoundErrMessage)
		return models.QueryResult{}, errors.ErrPinotQueryResultTableNotFound
	}

	logger.Info("PINOT_SQL_RESPONSE_COLUMN_NAMES", zap.Any("COLUMN_NAMES", sqlResponse.ResultTable.DataSchema.ColumnNames))
	logger.Info("PINOT_SQL_RESPONSE_COLUMN_DATA_TYPES", zap.Any("COLUMN_DATA_TYPES", sqlResponse.ResultTable.DataSchema.ColumnDataTypes))

	if len(sqlResponse.ResultTable.DataSchema.ColumnNames) == 0 {
		logger.Error(errors.PinotQueryResultTableColumnNamesNotFoundErrMessage)
		return models.QueryResult{}, errors.ErrPinotQueryResultTableColumnNamesNotFound
	}

	columnNames := sqlResponse.ResultTable.DataSchema.ColumnNames
	var err error

	for _, row := range sqlResponse.ResultTable.Rows {
		rowData := map[string]interface{}{}
		for colIndex, col := range row {
			rowData[columnNames[colIndex]], err = handleDataType(col)
			if err != nil {
				logger.Error(errors.PinotQueryResultTableColumnDataTypeConversionFailedErrMessage, zap.Error(err))
				return models.QueryResult{}, err
			}
		}
		queryResult.Rows = append(queryResult.Rows, rowData)
	}

	for colIndex, col := range columnNames {
		queryResult.Columns = append(queryResult.Columns, models.ColumnMetadata{
			Name:         col,
			DatabaseType: sqlResponse.ResultTable.GetColumnDataType(colIndex),
		})
	}

	return queryResult, nil
}

func handleDataType(value interface{}) (interface{}, error) {
	switch value.(type) {
	case json.Number:
		// Try parsing as an int
		if i, err := value.(json.Number).Int64(); err == nil {
			return i, nil
		}
		// Try parsing as a float
		if f, err := value.(json.Number).Float64(); err == nil {
			return f, nil
		}
		// If neither works, return the original value with an error
		return nil, errors.ErrInvalidJsonNumberValue
	}
	return value, nil
}
