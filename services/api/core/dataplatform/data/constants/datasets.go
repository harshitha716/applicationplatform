package constants

import "fmt"

const DatasetTableNameQueryParam = "datasets_table_name"
const DatasetTableName = "datasets"
const ZampTableName = "zamp"

const (
	DatasetIdColumnName                    = "id"
	DatasetMerchantIdColumnName            = "merchant_id"
	DatasetDatabricksTableNameColumnName   = "databricks_table_name"
	DatasetDatabricksFQTableNameColumnName = "databricks_fq_table_name"
	DatasetDatabricksSchemaColumnName      = "databricks_schema"
	DatasetDatabricksConfigColumnName      = "databricks_config"
	DatasetDatabricksStatsColumnName       = "databricks_stats"
	DatasetPinotTableNameColumnName        = "pinot_table_name"
	DatasetPinotSchemaColumnName           = "pinot_schema"
	DatasetPinotConfigColumnName           = "pinot_config"
	DatasetPinotStatsColumnName            = "pinot_stats"
	DatasetCreatedAtColumnName             = "created_at"
	DatasetUpdatedAtColumnName             = "updated_at"
	DatasetIsDeletedColumnName             = "is_deleted"
	DatasetDeletedAtColumnName             = "deleted_at"
	DatasetDatasetConfigColumnName         = "dataset_config"
)

var SelectDatasetColumnNames string = fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s", DatasetIdColumnName, DatasetMerchantIdColumnName, DatasetDatabricksFQTableNameColumnName, DatasetDatabricksTableNameColumnName, DatasetDatabricksSchemaColumnName, DatasetDatabricksConfigColumnName, DatasetDatabricksStatsColumnName, DatasetPinotTableNameColumnName, DatasetPinotSchemaColumnName, DatasetPinotConfigColumnName, DatasetPinotStatsColumnName, DatasetCreatedAtColumnName, DatasetUpdatedAtColumnName, DatasetIsDeletedColumnName, DatasetDeletedAtColumnName, DatasetDatasetConfigColumnName)

var QueryGetDatasetById string = fmt.Sprintf("SELECT %s FROM {{.%s}} WHERE %s = '{{.%s}}' AND %s = '{{.%s}}' AND %s = false", SelectDatasetColumnNames, DatasetTableNameQueryParam, DatasetIdColumnName, DatasetIdColumnName, DatasetMerchantIdColumnName, DatasetMerchantIdColumnName, DatasetIsDeletedColumnName)

type DAGDestination string

const (
	DAGDatasetDestination DAGDestination = "dataset"
)

type Datatype string

const (
	BooleanDataType Datatype = "boolean"

	DecimalDataType  Datatype = "decimal"
	DoubleDataType   Datatype = "double"
	FloatDataType    Datatype = "float"
	IntegerDataType  Datatype = "integer"
	SmallIntDataType Datatype = "smallint"
	TinyIntDataType  Datatype = "tinyint"
	BigIntDataType   Datatype = "bigint"

	StringDataType Datatype = "string"

	TimestampDataType    Datatype = "timestamp"
	TimestampNtzDataType Datatype = "timestamp_ntz"
	DateDataType         Datatype = "date"

	ArrayOfStringDataType Datatype = "array<string>"
)

type JSDatatype string

const (
	JSBoolean JSDatatype = "boolean"
	JSNumber  JSDatatype = "number"
	JSString  JSDatatype = "string"
	JSDate    JSDatatype = "Date"
	JSArray   JSDatatype = "Array"
)

// DataTypeToJSType maps backend Datatype to JavaScript/TypeScript types
var DataTypeToJSType = map[Datatype]JSDatatype{
	BooleanDataType: JSBoolean,

	DecimalDataType:  JSNumber,
	DoubleDataType:   JSNumber,
	FloatDataType:    JSNumber,
	IntegerDataType:  JSNumber,
	SmallIntDataType: JSNumber,
	TinyIntDataType:  JSNumber,
	BigIntDataType:   JSNumber,

	StringDataType: JSString,

	TimestampDataType:    JSDate,
	TimestampNtzDataType: JSDate,
	DateDataType:         JSDate,

	ArrayOfStringDataType: JSArray,
}

type Provider string

const (
	ProviderDatabricks Provider = "databricks"
	ProviderPinot      Provider = "pinot"
)
