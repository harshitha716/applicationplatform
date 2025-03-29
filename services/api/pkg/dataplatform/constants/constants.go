package constants

const DATABRICKS_DRIVER_NAME = "databricks"

const DatabaseTypeArray = "ARRAY"

const PINOT_DRIVER_NAME = "pinot"

const POSTGRES_DRIVER_NAME = "postgres"

type ProviderType string

const (
	ProviderTypeDatabricks ProviderType = DATABRICKS_DRIVER_NAME
	ProviderTypePinot      ProviderType = PINOT_DRIVER_NAME
	ProviderTypePostgres   ProviderType = POSTGRES_DRIVER_NAME
)
