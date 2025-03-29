package widgets

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	datasetconstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetconstants "github.com/Zampfi/application-platform/services/api/core/widgets/constants"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	querybuilderconstants "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
)

var widgetStrategies = map[string]DatasetParamsBuilder{
	"bar_chart":   BasicChartStrategy{BaseStrategy: *NewBaseStrategy()},
	"line_chart":  BasicChartStrategy{BaseStrategy: *NewBaseStrategy()},
	"pie_chart":   PieChartStrategy{BaseStrategy: *NewBaseStrategy()},
	"donut_chart": PieChartStrategy{BaseStrategy: *NewBaseStrategy()},
	"pivot_table": PivotTableStrategy{BaseStrategy: *NewBaseStrategy()},
	"kpi":         KPIStrategy{BaseStrategy: *NewBaseStrategy()},
}

func NewDatasetParamsBuilder(widgetType string) (DatasetParamsBuilder, error) {
	strategy, ok := widgetStrategies[widgetType]
	if !ok {
		return nil, fmt.Errorf("unsupported widget type: %s", widgetType)
	}
	return strategy, nil
}

type DatasetParamsBuilder interface {
	ToDatasetParams(instance *widgetmodels.WidgetInstance, datasetbuilderparams widgetmodels.DatasetBuilderParams) (map[string]widgetmodels.GetDataByDatasetIDParams, error)
}

// ProcessFieldsFunc is a function type for processing specific fields in a dataset params
type ProcessFieldsFunc func(*datasetmodels.DatasetParams, *widgetmodels.DataMappingFields, *datasetmodels.FilterModel, *widgetmodels.DatasetBuilderParams) error

// ParameterProcessor defines the interface for parameter processors
type ParameterProcessor interface {
	CanProcess(paramName string) bool
	Process(match string, parts []string, populationValues []string) string
}

// DateParameterProcessor processes date-related parameters
type DateParameterProcessor struct {
	BaseStrategy *BaseStrategy
}

// CanProcess checks if this processor can handle the given parameter
func (p *DateParameterProcessor) CanProcess(paramName string) bool {
	return paramName == widgetconstants.ParameterMethodToday ||
		paramName == widgetconstants.ParameterMethodEndDay ||
		paramName == widgetconstants.ParameterMethodStartDay
}

// Process processes a date parameter match
func (p *DateParameterProcessor) Process(match string, parts []string, populationValues []string) string {
	paramName := parts[1]
	baseTime := p.BaseStrategy.GetBaseTime(paramName, populationValues)

	// Simple {{.$today}} case
	if len(parts) < 3 || parts[2] == "" {
		return baseTime.Format(time.DateTime)
	}

	// Method case
	method := parts[2]
	args := ""
	if len(parts) >= 4 {
		args = parts[3]
	}
	computedTime := p.BaseStrategy.ApplyMethod(baseTime, method, args)
	return computedTime.Format(time.DateTime)
}

// BaseStrategy provides common functionality for all dataset parameter builder strategies
type BaseStrategy struct {
	// Pre-compiled regex patterns for better performance
	parameterizedValueRegex *regexp.Regexp
	defaultParametersRegex  *regexp.Regexp
}

// NewBaseStrategy creates a new BaseStrategy with pre-compiled regex patterns
func NewBaseStrategy() *BaseStrategy {
	return &BaseStrategy{
		parameterizedValueRegex: regexp.MustCompile(widgetconstants.ParametrizedValueFieldRegex),
		defaultParametersRegex:  regexp.MustCompile(widgetconstants.DefaultParametersRegex),
	}
}

// MergeFilters combines default filters with sheet filters
func (b *BaseStrategy) MergeFilters(defaultFilters, sheetFilters *datasetmodels.FilterModel) datasetmodels.FilterModel {
	if (defaultFilters == nil || len(defaultFilters.Conditions) == 0) &&
		(sheetFilters == nil || len(sheetFilters.Conditions) == 0) {
		return datasetmodels.FilterModel{}
	}

	if sheetFilters == nil || len(sheetFilters.Conditions) == 0 {
		return *defaultFilters
	}

	if defaultFilters == nil || len(defaultFilters.Conditions) == 0 {
		return *sheetFilters
	}

	if defaultFilters.LogicalOperator != sheetFilters.LogicalOperator {
		// FIXME: we dont' support OR . But handle this if needed
		return *sheetFilters
	}
	mergedConditions := defaultFilters.Conditions
	mergedConditions = append(mergedConditions, sheetFilters.Conditions...)

	return datasetmodels.FilterModel{
		LogicalOperator: datasetmodels.LogicalOperator(querybuilderconstants.LogicalOperatorAnd),
		Conditions:      mergedConditions,
	}
}

// ParametrizeDefaultFilters processes default filters with parametrized values
func (b *BaseStrategy) ParametrizeDefaultFilters(defaultFilters *datasetmodels.FilterModel, datasetID string, datasetBuilderParams widgetmodels.DatasetBuilderParams) *datasetmodels.FilterModel {
	if defaultFilters == nil {
		return nil
	}

	sheetFilters := datasetBuilderParams.Filters[datasetID].Filters
	sheetFilterConditionMap := make(map[string]datasetmodels.Filter)
	for _, condition := range sheetFilters.Conditions {
		sheetFilterConditionMap[condition.Column] = condition
	}

	for i := range defaultFilters.Conditions {
		if defaultFilters.Conditions[i].Value == nil {
			continue
		}
		sheetCondition, ok := sheetFilterConditionMap[defaultFilters.Conditions[i].Column]
		sheetConditionValues := []string{}
		if ok {
			sheetConditionValues = b.ParseAndGetValuesFromFilters(sheetCondition.Value)
		}
		values := b.ParseAndGetValuesFromFilters(defaultFilters.Conditions[i].Value)
		if values == nil {
			return defaultFilters
		}

		processedValues := b.ProcessParametrizedFilter(values, sheetConditionValues)

		// Maintain backward compatibility with existing code
		// If the original value was a string and we have a single processed value,
		// set the value as a string rather than a slice
		if _, ok := defaultFilters.Conditions[i].Value.(string); ok && len(processedValues) == 1 {
			defaultFilters.Conditions[i].Value = processedValues[0]
		} else if _, ok := defaultFilters.Conditions[i].Value.([]string); ok || len(processedValues) > 1 {
			// If the original value was a slice or we have multiple processed values,
			// set the value as a slice
			defaultFilters.Conditions[i].Value = processedValues
		} else {
			// For other cases, keep the original behavior
			defaultFilters.Conditions[i].Value = processedValues
		}
	}
	return defaultFilters
}

// ParseAndGetValuesFromFilters extracts string values from filter values
func (b *BaseStrategy) ParseAndGetValuesFromFilters(value interface{}) []string {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []string:
		return v

	case []interface{}:
		strValues := make([]string, len(v))
		for j, val := range v {
			if str, ok := val.(string); ok {
				strValues[j] = str
			}
		}
		return strValues

	default:
		// Try to convert to string using fmt.Sprint
		return nil
	}
}

// ProcessParametrizedFilter processes parametrized filter values
func (b *BaseStrategy) ProcessParametrizedFilter(values []string, sheetConditionValues []string) []string {
	if values == nil {
		return nil
	}

	parametrizedValues := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			// Skip empty values
			continue
		}

		if b.IsParametrizedString(value) {
			parametrizedValue := b.PopulateParams(value, sheetConditionValues)
			parametrizedValues = append(parametrizedValues, parametrizedValue)
		} else {
			parametrizedValues = append(parametrizedValues, value)
		}
	}

	return parametrizedValues
}

// IsParametrizedString checks if a string is parametrized
func (b *BaseStrategy) IsParametrizedString(value string) bool { // Supports format of {{.$<variable_name>.<method>(<args>)}}
	if b.parameterizedValueRegex == nil {
		b.parameterizedValueRegex = regexp.MustCompile(widgetconstants.ParametrizedValueFieldRegex)
	}
	return b.parameterizedValueRegex.MatchString(value)
}

// PopulateParams populates parameters in a string using registered processors
func (b *BaseStrategy) PopulateParams(value string, sheetConditionValues []string) string {
	processors := []ParameterProcessor{
		&DateParameterProcessor{BaseStrategy: b},
		// Add more processors here as needed
	}

	return b.processWithProcessors(value, sheetConditionValues, processors)
}

// processWithProcessors processes a string with the given processors
func (b *BaseStrategy) processWithProcessors(value string, populationValues []string, processors []ParameterProcessor) string {
	re := regexp.MustCompile(widgetconstants.DefaultParametersRegex)

	populatedValue := re.ReplaceAllStringFunc(value, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match // Return original if no match
		}

		paramName := parts[1]

		// Find a processor that can handle this parameter
		for _, processor := range processors {
			if processor.CanProcess(paramName) {
				return processor.Process(match, parts, populationValues)
			}
		}

		// No processor found, return original
		return match
	})

	return populatedValue
}

// PopulateDateParams populates date parameters in a string (legacy method for backward compatibility)
func (b *BaseStrategy) PopulateDateParams(value string, populationValues []string) string {
	processor := &DateParameterProcessor{BaseStrategy: b}
	return b.processWithProcessors(value, populationValues, []ParameterProcessor{processor})
}

// GetBaseTime gets the base time for date parameters
func (b *BaseStrategy) GetBaseTime(paramName string, populationValues []string) time.Time {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch paramName {
	case widgetconstants.ParameterMethodToday:
		return startOfDay

	case widgetconstants.ParameterMethodEndDay:
		if len(populationValues) == 0 {
			// Log error and return default
			fmt.Printf("Error: No population values for end date parameter\n")
			return startOfDay
		}

		endDateStr := populationValues[len(populationValues)-1]
		endDate, err := time.Parse(time.DateTime, endDateStr)
		if err != nil {
			// Log error and return default
			fmt.Printf("Error parsing end date '%s': %v\n", endDateStr, err)
			return startOfDay
		}

		return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, now.Location())

	case widgetconstants.ParameterMethodStartDay:
		if len(populationValues) == 0 {
			// Log error and return default
			fmt.Printf("Error: No population values for start date parameter\n")
			return startOfDay
		}

		startDateStr := populationValues[0]
		startDate, err := time.Parse(time.DateTime, startDateStr)
		if err != nil {
			// Log error and return default
			fmt.Printf("Error parsing start date '%s': %v\n", startDateStr, err)
			return startOfDay
		}

		return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, now.Location())

	default:
		// Log error and return default
		fmt.Printf("Error: Unknown parameter name '%s'\n", paramName)
		return startOfDay
	}
}

// ApplyMethod applies a method to a base time
func (b *BaseStrategy) ApplyMethod(baseTime time.Time, method string, args string) time.Time {
	switch method {
	case widgetconstants.ParameterMethodAddDays:
		days, err := strconv.Atoi(strings.TrimSpace(args))
		if err != nil {
			// Log error and return default
			fmt.Printf("Error parsing days argument '%s': %v\n", args, err)
			return baseTime
		}
		return baseTime.AddDate(0, 0, days)

	case widgetconstants.ParameterMethodAddSeconds:
		seconds, err := strconv.Atoi(strings.TrimSpace(args))
		if err != nil {
			// Log error and return default
			fmt.Printf("Error parsing seconds argument '%s': %v\n", args, err)
			return baseTime
		}
		return baseTime.Add(time.Duration(seconds) * time.Second)

	default:
		// Log error and return default
		fmt.Printf("Error: Unknown method '%s'\n", method)
		return baseTime
	}
}

// AddTimeColumn adds time column to dataset params
func (b *BaseStrategy) AddTimeColumn(params *datasetmodels.DatasetParams, timeColumnMap map[string]string, periodicity *string, datasetID string) error {
	if len(timeColumnMap) == 0 || periodicity == nil {
		return nil
	}

	if _, ok := widgetconstants.Periodicities[*periodicity]; !ok {
		return fmt.Errorf("invalid periodicity: %s", *periodicity)
	}

	for i := range params.Columns {
		if params.Columns[i].Column == timeColumnMap[datasetID] {
			params.Columns[i].Column = fmt.Sprintf("date_trunc('%s', %s)", *periodicity, params.Columns[i].Column)
		}
	}

	for i := range params.GroupBy {
		if params.GroupBy[i].Column == timeColumnMap[datasetID] {
			params.GroupBy[i].Column = fmt.Sprintf("date_trunc('%s', %s)", *periodicity, params.GroupBy[i].Column)
		}
	}

	for i := range params.OrderBy {
		if params.OrderBy[i].Column == timeColumnMap[datasetID] {
			params.OrderBy[i].Column = fmt.Sprintf("date_trunc('%s', %s)", *periodicity, params.OrderBy[i].Column)
		}
	}

	return nil
}

// HandleAggregation handles aggregation for a field
func (b *BaseStrategy) HandleAggregation(params *datasetmodels.DatasetParams, field widgetmodels.Field, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
	alias := field.GetAlias()
	if alias == nil {
		alias = &field.Column
	}

	switch field.Aggregation {
	case widgetconstants.WindowFunctionFirst, widgetconstants.WindowFunctionLast:
		var sortBy []widgetmodels.SortBy
		if len(field.SortBy) == 0 && len(mapping.SortBy) == 0 {
			return fmt.Errorf("sort by is required for window functions")
		}

		if len(field.SortBy) == 0 {
			sortBy = []widgetmodels.SortBy{mapping.SortBy[0]}
		} else {
			sortBy = field.SortBy
		}

		windowParams, err := b.BuildWindowBasedParams(field, sortBy, params.GroupBy, filters)
		if err != nil {
			return err
		}

		if params.Subquery == nil {
			params.Subquery = windowParams
			if len(params.Subquery.Windows) > 0 {
				for i := range params.Subquery.Windows[0].PartitionBy {
					if len(datasetBuilderParams.TimeColumns) > 0 && datasetBuilderParams.Periodicity != nil && params.Subquery.Windows[0].PartitionBy[i].Column == datasetBuilderParams.TimeColumns[mapping.DatasetID] {
						params.Subquery.Windows[0].PartitionBy[i].Column = fmt.Sprintf("date_trunc('%s', %s)", *datasetBuilderParams.Periodicity, params.Subquery.Windows[0].PartitionBy[i].Column)
					}
				}
			}
		}

		if params.Filters.LogicalOperator == "" {
			params.Filters.LogicalOperator = datasetmodels.LogicalOperator(querybuilderconstants.LogicalOperatorAnd)
		}

		params.Filters.Conditions = []datasetmodels.Filter{
			{
				Column:   "rn",
				Operator: querybuilderconstants.EqualOperator,
				Value:    1,
			},
		}

		params.Aggregations = append(params.Aggregations, datasetmodels.Aggregation{
			Column:   field.Column,
			Function: datasetmodels.AggregationFunction("sum"), // Dummy aggregation function
			Alias:    *alias,
		})

	default:
		params.Aggregations = append(params.Aggregations, datasetmodels.Aggregation{
			Column:   field.Column,
			Function: datasetmodels.AggregationFunction(field.Aggregation),
			Alias:    *alias,
		})
	}
	return nil
}

// BuildWindowBasedParams builds window based params for window functions
func (b *BaseStrategy) BuildWindowBasedParams(field widgetmodels.Field, sortBy []widgetmodels.SortBy, groupBy []datasetmodels.GroupBy, filters *datasetmodels.FilterModel) (*datasetmodels.DatasetParams, error) {
	if len(sortBy) == 0 {
		return nil, fmt.Errorf("sort by is required for window functions")
	}

	newFilters := datasetmodels.FilterModel{}
	if filters != nil {
		newFilters = *filters
	}

	newFilters.Conditions = append(newFilters.Conditions, datasetmodels.Filter{
		Column:   datasetconstants.ZampIsDeletedColumn,
		Operator: querybuilderconstants.EqualOperator,
		Value:    false,
	})
	// Build columns list including the value column and all group by columns
	columns := []datasetmodels.ColumnConfig{{Column: field.Column}}
	partitionBy := make([]datasetmodels.ColumnConfig, len(groupBy))
	for i, group := range groupBy {
		columns = append(columns, datasetmodels.ColumnConfig{Column: group.Column})
		partitionBy[i] = datasetmodels.ColumnConfig{Column: group.Column}
	}

	columns = append(columns, datasetmodels.ColumnConfig{Column: datasetconstants.ZampIsDeletedColumn})

	orderBy := make([]datasetmodels.OrderBy, len(sortBy))
	for i, sort := range sortBy {
		orderBy[i] = datasetmodels.OrderBy{
			Column: sort.Column,
			Order:  datasetmodels.OrderType(sort.Order),
		}
	}
	subqueryParams := &datasetmodels.DatasetParams{
		Columns: columns,
		Windows: []datasetmodels.WindowConfig{
			{
				Function:    string(querybuilderconstants.WindowFunctionRowNumber),
				PartitionBy: partitionBy,
				OrderBy:     orderBy,
				Alias:       "rn",
			},
		},
		Filters: newFilters,
	}

	return subqueryParams, nil
}

// AddSortBy adds sort by to dataset params
func (b *BaseStrategy) AddSortBy(params *datasetmodels.DatasetParams, sortBy []widgetmodels.SortBy) {
	// First add explicit sort columns
	for _, sort := range sortBy {
		column := sort.GetColumn()
		params.OrderBy = append(params.OrderBy, datasetmodels.OrderBy{
			Column: column,
			Order:  datasetmodels.OrderType(sort.Order),
			Alias:  &column,
		})
	}

	// Track which columns are already in the sort order
	sortedColumns := make(map[string]bool)
	for _, sort := range params.OrderBy {
		sortedColumns[sort.Column] = true
	}

	for _, group := range params.GroupBy {
		if !sortedColumns[*group.Alias] {
			params.OrderBy = append(params.OrderBy, datasetmodels.OrderBy{
				Column: *group.Alias,
				Order:  datasetmodels.OrderType(querybuilderconstants.OrderAsc),
				Alias:  group.Alias,
			})
			sortedColumns[*group.Alias] = true
		}
	}
}

// AddCurrency adds currency to dataset params
func (b *BaseStrategy) AddCurrency(params *datasetmodels.DatasetParams, currency *string) {
	if params.Subquery != nil {
		params.Subquery.FxCurrency = currency
		return
	}
	params.FxCurrency = currency
}

// InitializeDatasetParams initializes dataset params with empty collections
func (b *BaseStrategy) InitializeDatasetParams(filters datasetmodels.FilterModel) datasetmodels.DatasetParams {
	return datasetmodels.DatasetParams{
		Columns:      []datasetmodels.ColumnConfig{},
		Aggregations: []datasetmodels.Aggregation{},
		GroupBy:      []datasetmodels.GroupBy{},
		Filters:      filters,
	}
}

// ProcessDatasetParams is a template method that handles the common flow of operations for processing dataset parameters
func (b *BaseStrategy) ProcessDatasetParams(mapping *widgetmodels.DataMappingFields, datasetbuilderparams widgetmodels.DatasetBuilderParams, processFields ProcessFieldsFunc) (widgetmodels.GetDataByDatasetIDParams, error) {
	var sheetFilters datasetmodels.FilterModel
	if filterSet, exists := datasetbuilderparams.Filters[mapping.DatasetID]; exists {
		sheetFilters = filterSet.Filters
	}

	populatedDefaultFilters := b.ParametrizeDefaultFilters(mapping.DefaultFilters, mapping.DatasetID, datasetbuilderparams)
	combinedFilters := b.MergeFilters(populatedDefaultFilters, &sheetFilters)
	params := b.InitializeDatasetParams(combinedFilters)

	// Process specific fields based on widget type (provided by the strategy)
	if err := processFields(&params, mapping, &combinedFilters, &datasetbuilderparams); err != nil {
		return widgetmodels.GetDataByDatasetIDParams{}, err
	}

	// Add common elements
	if err := b.AddTimeColumn(&params, datasetbuilderparams.TimeColumns, datasetbuilderparams.Periodicity, mapping.DatasetID); err != nil {
		return widgetmodels.GetDataByDatasetIDParams{}, err
	}

	b.AddSortBy(&params, mapping.SortBy)
	b.AddCurrency(&params, datasetbuilderparams.Currency)

	return widgetmodels.GetDataByDatasetIDParams{
		DatasetID: mapping.DatasetID,
		Params:    params,
	}, nil
}

type BasicChartStrategy struct {
	BaseStrategy
}

func (b BasicChartStrategy) ToDatasetParams(instance *widgetmodels.WidgetInstance, datasetbuilderparams widgetmodels.DatasetBuilderParams) (map[string]widgetmodels.GetDataByDatasetIDParams, error) {
	if len(instance.DataMappings.Mappings) == 0 {
		return nil, fmt.Errorf("no mappings found for BasicChart widget")
	}

	result, err := b.ProcessDatasetParams(&instance.DataMappings.Mappings[0], datasetbuilderparams, func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
		// Handle x-axis (dimension)
		if xAxis, ok := mapping.Fields[widgetconstants.XAxisField]; ok {
			expression := xAxis[0].GetExpression()
			if expression == "" {
				expression = xAxis[0].Column
			}
			params.GroupBy = append(params.GroupBy, datasetmodels.GroupBy{Column: expression, Alias: xAxis[0].GetAlias()})
		}

		// Handle y-axis (measure)
		if yAxis, ok := mapping.Fields[widgetconstants.YAxisField]; ok {
			if err := b.HandleAggregation(params, yAxis[0], mapping, filters, datasetBuilderParams); err != nil {
				return fmt.Errorf("failed to handle aggregation for bar-chart y-axis: %w", err)
			}
		}

		// Handle group by fields
		groupBy, ok := mapping.Fields[widgetconstants.GroupByField]
		if ok {
			for _, field := range groupBy {
				params.GroupBy = append(params.GroupBy, datasetmodels.GroupBy{Column: field.Column, Alias: field.GetAlias()})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]widgetmodels.GetDataByDatasetIDParams{
		instance.DataMappings.Mappings[0].Ref: result,
	}, nil
}

type PieChartStrategy struct {
	BaseStrategy
}

func (p PieChartStrategy) ToDatasetParams(instance *widgetmodels.WidgetInstance, datasetbuilderparams widgetmodels.DatasetBuilderParams) (map[string]widgetmodels.GetDataByDatasetIDParams, error) {
	if len(instance.DataMappings.Mappings) == 0 {
		return nil, fmt.Errorf("no mappings found for PieChart widget")
	}

	result, err := p.ProcessDatasetParams(&instance.DataMappings.Mappings[0], datasetbuilderparams, func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
		// Handles slices (dimensions)
		if slices, ok := mapping.Fields[widgetconstants.SlicesField]; ok {
			expression := slices[0].GetExpression()
			if expression == "" {
				expression = slices[0].Column
			}
			params.GroupBy = append(params.GroupBy, datasetmodels.GroupBy{Column: expression, Alias: slices[0].GetAlias()})
		}

		// Handles value (measure)
		if value, ok := mapping.Fields[widgetconstants.ValuesField]; ok {
			if err := p.HandleAggregation(params, value[0], mapping, filters, datasetBuilderParams); err != nil {
				return fmt.Errorf("failed to handle aggregation for pie-chart value: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]widgetmodels.GetDataByDatasetIDParams{
		instance.DataMappings.Mappings[0].Ref: result,
	}, nil
}

type PivotTableStrategy struct {
	BaseStrategy
}

func (p PivotTableStrategy) ToDatasetParams(instance *widgetmodels.WidgetInstance, datasetbuilderparams widgetmodels.DatasetBuilderParams) (map[string]widgetmodels.GetDataByDatasetIDParams, error) {
	result := make(map[string]widgetmodels.GetDataByDatasetIDParams)

	for _, mapping := range instance.DataMappings.Mappings {
		datasetResult, err := p.ProcessDatasetParams(&mapping, datasetbuilderparams, func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
			// Handle row fields
			if rows, ok := mapping.Fields[widgetconstants.RowsField]; ok {
				for _, row := range rows {
					expression := row.GetExpression()
					if expression == "" {
						expression = row.Column
					}
					params.GroupBy = append(params.GroupBy, datasetmodels.GroupBy{Column: expression, Alias: row.GetAlias()})
				}
			}

			// Handle column fields
			if cols, ok := mapping.Fields[widgetconstants.ColumnsField]; ok {
				for _, col := range cols {
					expression := col.GetExpression()
					if expression == "" {
						expression = col.Column
					}
					params.GroupBy = append(params.GroupBy, datasetmodels.GroupBy{Column: expression, Alias: col.GetAlias()})
				}
			}

			// Handle value fields
			if values, ok := mapping.Fields[widgetconstants.ValuesField]; ok {
				for _, value := range values {
					if err := p.HandleAggregation(params, value, mapping, filters, datasetBuilderParams); err != nil {
						return fmt.Errorf("failed to handle aggregation for pivot-table value: %w", err)
					}
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		result[mapping.Ref] = datasetResult
	}

	return result, nil
}

type KPIStrategy struct {
	BaseStrategy
}

func (k KPIStrategy) ToDatasetParams(instance *widgetmodels.WidgetInstance, datasetbuilderparams widgetmodels.DatasetBuilderParams) (map[string]widgetmodels.GetDataByDatasetIDParams, error) {
	if len(instance.DataMappings.Mappings) == 0 {
		return nil, fmt.Errorf("no mappings found for KPI widget")
	}

	result, err := k.ProcessDatasetParams(&instance.DataMappings.Mappings[0], datasetbuilderparams, func(params *datasetmodels.DatasetParams, mapping *widgetmodels.DataMappingFields, filters *datasetmodels.FilterModel, datasetBuilderParams *widgetmodels.DatasetBuilderParams) error {
		if primaryValue, ok := mapping.Fields[widgetconstants.PrimaryValueField]; ok {
			if err := k.HandleAggregation(params, primaryValue[0], mapping, filters, datasetBuilderParams); err != nil {
				return fmt.Errorf("failed to handle aggregation for kpi value: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]widgetmodels.GetDataByDatasetIDParams{
		instance.DataMappings.Mappings[0].Ref: result,
	}, nil
}
