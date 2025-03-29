package errors

import "errors"

const (
	PinotServiceAlreadyInitializedErrMessage                      = "ERR_PINOT_SERVICE_ALREADY_INITIALIZED"
	InvalidConfigurationForPinotErrMessage                        = "ERR_INVALID_CONFIGURATION_FOR_PINOT"
	DatabricksServiceAlreadyInitializedErrMessage                 = "ERR_DATABRICKS_SERVICE_ALREADY_INITIALIZED"
	InvalidConfigurationForDatabricksErrMessage                   = "ERR_INVALID_CONFIGURATION_FOR_DATABRICKS"
	UnsupportedProviderConfigurationErrMessage                    = "ERR_UNSUPPORTED_PROVIDER_CONFIGURATION"
	DatabricksServiceNotInitializedErrMessage                     = "ERR_DATABRICKS_SERVICE_NOT_INITIALIZED"
	PinotServiceNotInitializedErrMessage                          = "ERR_PINOT_SERVICE_NOT_INITIALIZED"
	PostgresServiceNotInitializedErrMessage                       = "ERR_POSTGRES_SERVICE_NOT_INITIALIZED"
	QueryingPinotFailedErrMessage                                 = "ERR_QUERYING_PINOT"
	BuildingQueryFailedErrMessage                                 = "ERR_BUILDING_QUERY_FAILED"
	QueryingDatabricksFailedErrMessage                            = "ERR_QUERYING_DATABRICKS_FAILED"
	UnsupportedProviderTypeErrMessage                             = "ERR_UNSUPPORTED_PROVIDER_TYPE"
	QueryingPostgresFailedErrMessage                              = "ERR_QUERYING_POSTGRES_FAILED"
	PostgresServiceInitializationFailedErrMessage                 = "ERR_POSTGRES_SERVICE_INITIALIZATION_FAILED"
	InvalidConfigurationForPostgresErrMessage                     = "ERR_INVALID_CONFIGURATION_FOR_POSTGRES"
	PostgresServiceAlreadyInitializedErrMessage                   = "ERR_POSTGRES_SERVICE_ALREADY_INITIALIZED"
	ProviderConfigNotInitializedErrMessage                        = "ERR_PROVIDER_CONFIG_NOT_INITIALIZED"
	ProviderServiceNotFoundErrMessage                             = "ERR_PROVIDER_SERVICE_NOT_FOUND"
	DatabricksServiceInitializationFailedErrMessage               = "ERR_DATABRICKS_SERVICE_INITIALIZATION_FAILED"
	PinotServiceInitializationFailedErrMessage                    = "ERR_PINOT_SERVICE_INITIALIZATION_FAILED"
	PinotQueryResultConversionFailedErrMessage                    = "ERR_PINOT_QUERY_RESULT_CONVERSION_FAILED"
	InvalidJsonNumberValueErrMessage                              = "ERR_INVALID_JSON_NUMBER_VALUE"
	PinotQueryResultTableNotFoundErrMessage                       = "ERR_PINOT_QUERY_RESULT_TABLE_NOT_FOUND"
	PinotQueryResultTableColumnNamesNotFoundErrMessage            = "ERR_PINOT_QUERY_RESULT_TABLE_COLUMN_NAMES_NOT_FOUND"
	PinotQueryResultTableColumnDataTypeConversionFailedErrMessage = "ERR_PINOT_QUERY_RESULT_TABLE_COLUMN_DATA_TYPE_CONVERSION_FAILED"
	PinotQueryExceptionsErrMessage                                = "ERR_PINOT_QUERY_EXCEPTIONS"
)

var (
	ErrPinotServiceAlreadyInitialized                      = errors.New(PinotServiceAlreadyInitializedErrMessage)
	ErrInvalidConfigurationForPinot                        = errors.New(InvalidConfigurationForPinotErrMessage)
	ErrDatabricksServiceAlreadyInitialized                 = errors.New(DatabricksServiceAlreadyInitializedErrMessage)
	ErrInvalidConfigurationForDatabricks                   = errors.New(InvalidConfigurationForDatabricksErrMessage)
	ErrUnsupportedProviderConfiguration                    = errors.New(UnsupportedProviderConfigurationErrMessage)
	ErrDatabricksServiceNotInitialized                     = errors.New(DatabricksServiceNotInitializedErrMessage)
	ErrPinotServiceNotInitialized                          = errors.New(PinotServiceNotInitializedErrMessage)
	ErrQueryingDatabricks                                  = errors.New(QueryingDatabricksFailedErrMessage)
	ErrBuildingQueryResult                                 = errors.New(BuildingQueryFailedErrMessage)
	ErrQueryingPinot                                       = errors.New(QueryingPinotFailedErrMessage)
	ErrUnsupportedProviderType                             = errors.New(UnsupportedProviderTypeErrMessage)
	ErrQueryingPostgres                                    = errors.New(QueryingPostgresFailedErrMessage)
	ErrPostgresServiceNotInitialized                       = errors.New(PostgresServiceNotInitializedErrMessage)
	ErrPostgresServiceInitializationFailed                 = errors.New(PostgresServiceInitializationFailedErrMessage)
	ErrPostgresServiceAlreadyInitialized                   = errors.New(PostgresServiceAlreadyInitializedErrMessage)
	ErrInvalidConfigurationForPostgres                     = errors.New(InvalidConfigurationForPostgresErrMessage)
	ErrProviderConfigNotInitialized                        = errors.New(ProviderConfigNotInitializedErrMessage)
	ErrProviderServiceNotFound                             = errors.New(ProviderServiceNotFoundErrMessage)
	ErrDatabricksServiceInitializationFailed               = errors.New(DatabricksServiceInitializationFailedErrMessage)
	ErrPinotServiceInitializationFailed                    = errors.New(PinotServiceInitializationFailedErrMessage)
	ErrPinotQueryResultConversionFailed                    = errors.New(PinotQueryResultConversionFailedErrMessage)
	ErrInvalidJsonNumberValue                              = errors.New(InvalidJsonNumberValueErrMessage)
	ErrPinotQueryResultTableNotFound                       = errors.New(PinotQueryResultTableNotFoundErrMessage)
	ErrPinotQueryResultTableColumnNamesNotFound            = errors.New(PinotQueryResultTableColumnNamesNotFoundErrMessage)
	ErrPinotQueryResultTableColumnDataTypeConversionFailed = errors.New(PinotQueryResultTableColumnDataTypeConversionFailedErrMessage)
	ErrPinotQueryExceptions                                = errors.New(PinotQueryExceptionsErrMessage)
)
