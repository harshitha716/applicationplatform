package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/helper"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
)

type QueryBuilder interface {
	ToSQL(ctx context.Context, queryConfig models.QueryConfig) (string, map[string]interface{}, error)
	ToFilterSQL(ctx context.Context, filterConfig models.FilterModel) (string, map[string]interface{}, error)
}

type queryBuilder struct{}

func NewQueryBuilder() QueryBuilder {
	return &queryBuilder{}
}

func (qb *queryBuilder) ToSQL(ctx context.Context, queryConfig models.QueryConfig) (string, map[string]interface{}, error) {
	return qb.buildSelectQuery(ctx, queryConfig)
}

func (qb *queryBuilder) ToFilterSQL(ctx context.Context, filterConfig models.FilterModel) (string, map[string]interface{}, error) {
	var queryBuilder strings.Builder
	params := make(map[string]interface{})

	for i, filter := range filterConfig.Conditions {
		if i > 0 {
			if !helper.Contains([]models.LogicalOperator{constants.LogicalOperatorAnd, constants.LogicalOperatorOr}, filterConfig.LogicalOperator) {
				return "", nil, errors.ErrInvalidDataType
			}
			queryBuilder.WriteString(fmt.Sprintf(" %s ", filterConfig.LogicalOperator))
		}
		conditionString, err := qb.buildCondition(ctx, filter)
		if err != nil {
			return "", nil, err
		}
		queryBuilder.WriteString(conditionString)
	}

	return queryBuilder.String(), params, nil
}

func (qb *queryBuilder) buildSelectQuery(ctx context.Context, queryConfig models.QueryConfig) (string, map[string]interface{}, error) {
	var queryBuilder strings.Builder
	params := make(map[string]interface{})

	// Start building the SQL query
	queryBuilder.WriteString(constants.SelectStatement)

	// Prepare list of columns to select, including group by columns
	selectedColumns := []string{}
	for _, column := range queryConfig.TableConfig.Columns {
		selectedColumn, err := column.GetSelectColumn()
		if err != nil {
			return "", nil, err
		}
		selectedColumns = append(selectedColumns, selectedColumn)
	}

	// Add window functions if present in main query
	for _, window := range queryConfig.Windows {
		windowStr, err := qb.buildWindowFunction(window)
		if err != nil {
			return "", nil, err
		}
		selectedColumns = append(selectedColumns, windowStr)
	}

	groupByColumns := make([]string, len(queryConfig.GroupBy))
	if len(queryConfig.GroupBy) > 0 {
		selectedColumns = []string{}
		for i, group := range queryConfig.GroupBy {
			groupByColumn := group.Column
			groupByColumnString, err := groupByColumn.GetGroupByColumn()
			if err != nil {
				return "", nil, err
			}
			groupByColumns[i] = groupByColumnString
			groupBySelectColumn, err := groupByColumn.GetGroupBySelectColumn()
			if err != nil {
				return "", nil, err
			}
			selectedColumns = append(selectedColumns, groupBySelectColumn)
		}
	}

	// Adding aggregations
	if len(queryConfig.Aggregations) > 0 {
		if len(queryConfig.GroupBy) == 0 {
			selectedColumns = []string{}
		}

		for _, aggregation := range queryConfig.Aggregations {
			aggregationColumn, err := aggregation.Column.GetAggregationColumn()
			if err != nil {
				return "", nil, err
			}
			selectedColumns = append(selectedColumns, fmt.Sprintf("%s(%s)%s\"%s\"", aggregation.Function, aggregationColumn, constants.AsStatement, aggregation.Alias))
		}
	}

	if len(selectedColumns) == 0 {
		queryBuilder.WriteString(constants.AllColumns)
	} else {
		queryBuilder.WriteString(strings.Join(selectedColumns, ", "))
	}

	// Handle FROM clause with potential subquery
	if queryConfig.Subquery != nil {
		queryBuilder.WriteString(constants.FromStatement)
		queryBuilder.WriteString("( ")
		subqueryStr, subParams, err := qb.buildSelectQuery(ctx, *queryConfig.Subquery)
		if err != nil {
			return "", nil, err
		}
		queryBuilder.WriteString(subqueryStr)
		queryBuilder.WriteString(" ) subquery")

		// Merge subquery params with main params
		for k, v := range subParams {
			params[k] = v
		}
	} else {
		// Regular FROM clause for table
		queryBuilder.WriteString(fmt.Sprintf("%s{{.%s%s}}",
			constants.FromStatement,
			constants.ZampDataset,
			queryConfig.TableConfig.DatasetId))
		params[fmt.Sprintf("%s%s", constants.ZampDataset, queryConfig.TableConfig.DatasetId)] = queryConfig.TableConfig.DatasetId
	}

	// Adding filters
	if len(queryConfig.Filters.Conditions) > 0 {
		queryBuilder.WriteString(constants.WhereStatement)
		for i, filter := range queryConfig.Filters.Conditions {
			if i > 0 {
				if !helper.Contains([]models.LogicalOperator{constants.LogicalOperatorAnd, constants.LogicalOperatorOr}, queryConfig.Filters.LogicalOperator) {
					return "", nil, errors.ErrInvalidDataType
				}
				queryBuilder.WriteString(fmt.Sprintf(" %s ", queryConfig.Filters.LogicalOperator))
			}
			conditionString, err := qb.buildCondition(ctx, filter)
			if err != nil {
				return "", nil, err
			}
			queryBuilder.WriteString(conditionString)
		}
	}

	// Adding GROUP BY clauses
	if len(queryConfig.GroupBy) > 0 {
		queryBuilder.WriteString(fmt.Sprintf("%s%s", constants.GroupByStatement, strings.Join(groupByColumns, ", ")))
	}

	// Adding ORDER BY clauses
	if len(queryConfig.OrderBy) > 0 {
		orderColumns := make([]string, len(queryConfig.OrderBy))
		for i, order := range queryConfig.OrderBy {
			orderColumnString, err := order.Column.GetOrderByColumn()
			if err != nil {
				return "", nil, err
			}
			orderColumns[i] = fmt.Sprintf("%s %s", orderColumnString, order.Order)
		}
		queryBuilder.WriteString(fmt.Sprintf("%s%s", constants.OrderByStatement, strings.Join(orderColumns, ", ")))
	}

	// Adding pagination if present
	if queryConfig.Pagination != nil {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", queryConfig.Pagination.PageSize, queryConfig.Pagination.PageSize*(queryConfig.Pagination.Page-1)))
	}

	// Return the final query string along with the params map
	return queryBuilder.String(), params, nil

}

// Helper function to build SQL condition strings from filters
func (qb *queryBuilder) buildCondition(ctx context.Context, filter models.Filter) (string, error) {
	subConditions := make([]string, len(filter.Conditions))
	conditionString := ""
	if len(filter.Conditions) > 0 && filter.LogicalOperator != nil {
		for i, subFilter := range filter.Conditions {
			condition, err := qb.buildCondition(ctx, subFilter)
			if err != nil {
				return "", err
			}
			subConditions[i] = condition
		}
		if !helper.Contains([]models.LogicalOperator{constants.LogicalOperatorAnd, constants.LogicalOperatorOr}, *filter.LogicalOperator) {
			return "", errors.ErrInvalidDataType
		}
		conditionString = strings.Join(subConditions, fmt.Sprintf(" %s ", *filter.LogicalOperator))
	}

	if filter.Operator != "" {
		value, err := helper.ToSqlValue(filter.Value)
		if err != nil {
			return "", err
		}

		sqlCondition, err := qb.prepareSql(ctx, filter.Column, filter.Operator, value)
		if err != nil {
			return "", err
		}

		if len(filter.Conditions) > 0 && filter.LogicalOperator != nil {
			return fmt.Sprintf("%s %s (%s)", sqlCondition, *filter.LogicalOperator, conditionString), nil
		}

		return fmt.Sprintf("%s", sqlCondition), nil
	}

	return fmt.Sprintf("(%s)", conditionString), nil
}
