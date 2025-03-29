package service

import (
	"context"
	"fmt"
	"strings"

	dataplatformConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/helper"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
)

func (qb *queryBuilder) prepareSql(ctx context.Context, column models.ColumnConfig, operator string, value interface{}) (string, error) {
	if column.Datatype == nil {
		return "", errors.ErrInvalidDataType
	}

	columnName := column.Column
	if column.Alias != nil {
		columnName = fmt.Sprintf("\"%s\"", *column.Alias)
	}

	switch *column.Datatype {
	case dataplatformConstants.ArrayOfStringDataType:
		switch operator {
		case constants.ArrayInOperator:
			return qb.buildArrayInClause(columnName, value)
		case constants.ArrayContainsOperator:
			return qb.buildArrayContainsClause(columnName, value)
		default:
			return "", errors.ErrInvalidOperator
		}
	default:
		switch operator {
		case constants.InOperator:
			return qb.buildInClause(columnName, value)
		case constants.NotInOperator:
			return qb.buildNotInClause(columnName, value)
		case constants.ContainsOperator:
			return qb.buildContainsClause(columnName, value)
		case constants.NotContainsOperator:
			return qb.buildNotContainsClause(columnName, value)
		case constants.StartsWithOperator:
			return qb.buildStartsWithClause(columnName, value)
		case constants.StartsWithCaseSensitiveOperator:
			return qb.buildStartsWithCaseSensitiveClause(columnName, value)
		case constants.EndsWithOperator:
			return qb.buildEndsWithClause(columnName, value)
		case constants.InBetweenOperator:
			return qb.buildInBetweenClause(columnName, value)
		case constants.IsNullOperator:
			return qb.buildIsNullClause(columnName)
		case constants.EqualOperator:
			return qb.buildEqualClause(columnName, value)
		default:
			sqlOperator, exists := constants.SqlOperatorMap[operator]
			if !exists {
				return "", errors.ErrInvalidOperator
			}
			valueString, ok := value.(string)
			if !ok {
				return "", errors.ErrInvalidDataType
			}
			return fmt.Sprintf("( %s %s %s )", columnName, sqlOperator, valueString), nil
		}
	}
}

func (qb *queryBuilder) buildInClause(column string, value interface{}) (string, error) {
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( %s %s (%s) )", column, constants.SqlOperatorMap[constants.InOperator], strings.Join(valueArr, ", ")), nil
}

func (qb *queryBuilder) buildNotInClause(column string, value interface{}) (string, error) {
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( %s %s (%s) )", column, constants.SqlOperatorMap[constants.NotInOperator], strings.Join(valueArr, ", ")), nil
}

func (qb *queryBuilder) buildContainsClause(column string, value interface{}) (string, error) {
	var constituents []string
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil {
		return "", errors.ErrInvalidDataType
	}

	for _, value := range valueArr {
		constituents = append(constituents, fmt.Sprintf("LOWER(%s) %s '%%%s%%'", column, constants.SqlOperatorMap[constants.ContainsOperator], strings.ToLower(helper.RemoveFirstAndLastQuote(value))))
	}
	return fmt.Sprintf("( %s )", strings.Join(constituents, " OR ")), nil
}

func (qb *queryBuilder) buildNotContainsClause(column string, value interface{}) (string, error) {
	var constituents []string
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil {
		return "", errors.ErrInvalidDataType
	}

	for _, value := range valueArr {
		constituents = append(constituents, fmt.Sprintf("LOWER(%s) %s '%%%s%%'", column, constants.SqlOperatorMap[constants.NotContainsOperator], strings.ToLower(helper.RemoveFirstAndLastQuote(value))))
	}
	return fmt.Sprintf("( %s )", strings.Join(constituents, " AND ")), nil
}

func (qb *queryBuilder) buildStartsWithClause(column string, value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( LOWER(%s) %s '%s%%' )", column, constants.SqlOperatorMap[constants.StartsWithOperator], strings.ToLower(helper.RemoveFirstAndLastQuote(valueString))), nil
}

func (qb *queryBuilder) buildStartsWithCaseSensitiveClause(column string, value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( %s %s '%s%%' )", column, constants.SqlOperatorMap[constants.StartsWithCaseSensitiveOperator], helper.RemoveFirstAndLastQuote(valueString)), nil
}

func (qb *queryBuilder) buildEndsWithClause(column string, value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( LOWER(%s) %s '%%%s' )", column, constants.SqlOperatorMap[constants.EndsWithOperator], strings.ToLower(helper.RemoveFirstAndLastQuote(valueString))), nil
}

func (qb *queryBuilder) buildInBetweenClause(column string, value interface{}) (string, error) {
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil || len(valueArr) != 2 {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( %s %s %s AND %s )", column, constants.SqlOperatorMap[constants.InBetweenOperator], valueArr[0], valueArr[1]), nil
}

func (qb *queryBuilder) buildIsNullClause(column string) (string, error) {
	return fmt.Sprintf("( %s %s )", column, constants.SqlOperatorMap[constants.IsNullOperator]), nil
}

func (qb *queryBuilder) buildEqualClause(column string, value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( %s %s %s )", column, constants.SqlOperatorMap[constants.EqualOperator], valueString), nil
}

func (qb *queryBuilder) buildNotEqualClause(column string, value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", errors.ErrInvalidDataType
	}
	return fmt.Sprintf("( %s %s %s )", column, constants.SqlOperatorMap[constants.NotEqualOperator], valueString), nil
}

func (qb *queryBuilder) buildArrayInClause(column string, value interface{}) (string, error) {
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil {
		return "", errors.ErrInvalidDataType
	}
	var constituents []string
	for _, value := range valueArr {
		constituents = append(constituents, fmt.Sprintf("( LOWER(ARRAY_TO_STRING(%s, ',')) = %s )", column, strings.ToLower(value)))
	}
	return fmt.Sprintf("( %s )", strings.Join(constituents, " OR ")), nil
}

func (qb *queryBuilder) buildArrayContainsClause(column string, value interface{}) (string, error) {
	valueArr, err := helper.ConvertInterfaceSliceToStrings(value)
	if err != nil {
		return "", errors.ErrInvalidDataType
	}
	var constituents []string
	for _, value := range valueArr {
		constituents = append(constituents, fmt.Sprintf("( LOWER(ARRAY_TO_STRING(%s, ',')) LIKE '%%%s%%' )", column, strings.ToLower(helper.RemoveFirstAndLastQuote(value))))
	}
	return fmt.Sprintf("( %s )", strings.Join(constituents, " OR ")), nil
}

func (qb *queryBuilder) buildWindowFunction(window models.WindowConfig) (string, error) {
	var builder strings.Builder

	builder.WriteString(string(window.Function))
	builder.WriteString(constants.WindowStatement)

	if len(window.PartitionBy) > 0 {
		builder.WriteString(constants.PartitionByStatement)
		partitionCols := make([]string, len(window.PartitionBy))
		for i, col := range window.PartitionBy {
			colStr, err := col.GetSelectColumn()
			if err != nil {
				return "", err
			}
			partitionCols[i] = colStr
		}
		builder.WriteString(strings.Join(partitionCols, ", "))
	}

	if len(window.OrderBy) > 0 {
		builder.WriteString(constants.OrderByStatement)
		orderCols := make([]string, len(window.OrderBy))
		for i, order := range window.OrderBy {
			colStr, err := order.Column.GetOrderByColumn()
			if err != nil {
				return "", err
			}
			orderCols[i] = fmt.Sprintf("%s %s", colStr, order.Order)
		}
		builder.WriteString(strings.Join(orderCols, ", "))
	}

	builder.WriteString(" )")

	if window.Alias != "" {
		builder.WriteString(fmt.Sprintf(" AS \"%s\"", window.Alias))
	}

	return builder.String(), nil
}
