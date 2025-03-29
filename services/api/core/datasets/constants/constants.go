package constants

import "time"

const (
	FilterTypeAmountRange string = "amount-range"
	FilterTypeDateRange   string = "date-range"
	FilterTypeSearch      string = "search"
	FilterTypeArraySearch string = "array-search"
	FilterTypeMultiSearch string = "multi-select"
	FilterTypeSelect      string = "select"
)

const MultiSelectThreshold = 20
const ZampDatasetPrefix = "zamp_"

const ZampColumnPrefix = "_zamp"
const UnderscorePrefix = "_"

const ZampUpdateColumnSourcePrefix = "_zamp_source_json_"
const ZampFxColumnPrefix = "_zamp_fx_json_"

const GetDistinctValuesQuery = "SELECT DISTINCT %s FROM {{.zamp_%s}} where %s = False LIMIT %d"
const GetDistinctValuesQueryWithoutLimit = "SELECT DISTINCT %s FROM {{.zamp_%s}} where %s = False"
const GetRowDetailsQuery = "SELECT * FROM {{.zamp_%s}} WHERE _zamp_id = '%s'"
const GetRowCountQuery = "SELECT COUNT(*) FROM (%s)"

const (
	ZampIDColumn                = "_zamp_id"
	ZampDrilldownMetadataColumn = "_zamp_drilldown_metadata_json"
	ZampIsDeletedColumn         = "_zamp_is_deleted"
)

var WhiteListedZampColumns = []string{
	ZampIDColumn,
}

const (
	MetadataConfigFormat               = "format"
	MetadataConfigFormatDefault        = "dd/MM/yyyy"
	MetadataConfigCurrencyColumnPrefix = "currency_column_prefix"
	MetadataConfig                     = "config"
	MetadataConfigCustomType           = "custom_type"
	MetadataConfigIsHidden             = "is_hidden"
	MetadataConfigIsEditable           = "is_editable"
)

const (
	AmountColumn   = "amount_column"
	CurrencyColumn = "currency_column"
)

var CustomColumnGroupConfig = map[string]string{
	AmountColumn:   "amount_column",
	CurrencyColumn: "currency_column",
}

const (
	DefaultPaginationPage        = 1
	DefaultMaxPaginationPageSize = 1000
	DefaultPaginationPageSize    = 100
)

type UpdateColumnSourceType string

const (
	UpdateColumnSourceTypeRule UpdateColumnSourceType = "rule"
	UpdateColumnSourceTypeUser UpdateColumnSourceType = "user"
)

const (
	CustomColumnGroupTypeCurrencyAmount = "currency_amount"
	CustomColumnGroupTypeDateTime       = "date_time"
	CustomColumnGroupTypeTags           = "tags"
	CustomColumnGroupTypeStatus         = "status"
	CustomColumnGroupTypePG             = "pg"
)

const (
	DatasetExportPage          = 1
	DatasetExportPageSizeLimit = 100 * 1000
)

const (
	DatasetExportTimestampFormat = "2006-01-02_15-04-05"
	DatasetExportFilePathFormat  = "export-data/dataset/%s/workflow/%s/%s"
	DatasetExportFileNameFormat  = "%s_%s.csv"
)

const (
	DatasetConfigIsFxEnabled         = "is_fx_enabled"
	DatasetConfigIsFileImportEnabled = "is_file_import_enabled"
)

const (
	DataplatformProviderDatabricks = "databricks"
	DataplatformProviderPinot      = "pinot"
)

const (
	DatasetFilterConfigCacheKey = "dataset_filter_config"
)

const (
	DatasetFilterConfigCacheExpiry = time.Minute * 10
)
