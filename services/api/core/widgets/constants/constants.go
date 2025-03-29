package constants

const (
	MeasureType   = "measure"
	DimensionType = "dimension"
)

const ValuesField = "values"

// Basic chart fields
const (
	XAxisField   = "x_axis"
	YAxisField   = "y_axis"
	GroupByField = "group_by"
)

// Pie chart fields
const (
	SlicesField = "slices"
)

// Pivot table fields
const (
	RowsField    = "rows"
	ColumnsField = "columns"
)

// KPI fields
const (
	PrimaryValueField    = "primary_value"
	ComparisonValueField = "comparison_value"
	TimeComparisonField  = "time_comparison"
)

const (
	REF_PREFIX       = "__REF"
	HEIRARCHY_SUFFIX = "LEVEL"
)

var Periodicities map[string]bool = map[string]bool{
	"day":     true,
	"week":    true,
	"month":   true,
	"quarter": true,
	"year":    true,
}

const DEFAULT_CURRENCY = "USD"

type UserFacingDatatype string

const (
	UserFacingDatatypeString    UserFacingDatatype = "string"
	UserFacingDatatypeNumber    UserFacingDatatype = "number"
	UserFacingDatatypeBoolean   UserFacingDatatype = "boolean"
	UserFacingDatatypeTimestamp UserFacingDatatype = "timestamp"
	UserFacingDatatypeCountry   UserFacingDatatype = "country"
	UserFacingDatatypeAmount    UserFacingDatatype = "amount"
	UserFacingDatatypeTag       UserFacingDatatype = "tag"
	UserFacingDatatypeBank      UserFacingDatatype = "bank"
	UserFacingDatatypeStatus    UserFacingDatatype = "status"
)

const (
	WindowFunctionFirst string = "first"
	WindowFunctionLast  string = "last"
)

const (
	MAX_PAGE_SIZE = 50000
)

const (
	ParameterMethodAddDays      = "addDays"
	ParameterMethodAddSeconds   = "addSeconds"
	DefaultVariableSymbol       = "$"
	ParameterMethodToday        = DefaultVariableSymbol + "today"
	ParameterMethodEndDay       = DefaultVariableSymbol + "end_date"
	ParameterMethodStartDay     = DefaultVariableSymbol + "start_date"
	DefaultParametersRegex      = `{{\.(\$[a-zA-Z0-9_]+)(?:\.([a-zA-Z0-9_]+)(?:\(([^)]*)\))?)?}}`
	ParametrizedValueFieldRegex = `^{{\.(\$[^}]+)}}$`
)
