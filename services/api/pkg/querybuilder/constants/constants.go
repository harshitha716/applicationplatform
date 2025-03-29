package constants

import (
	models "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
)

const (
	AggregationFunctionSum models.AggregationFunction = "SUM"
	AggregationFunctionAvg models.AggregationFunction = "AVG"
	AggregationFunctionMin models.AggregationFunction = "MIN"
	AggregationFunctionMax models.AggregationFunction = "MAX"
)

const (
	OrderAsc  models.OrderType = "ASC"
	OrderDesc models.OrderType = "DESC"
)

var (
	LogicalOperatorAnd models.LogicalOperator = "AND"
	LogicalOperatorOr  models.LogicalOperator = "OR"
)

const (
	EqualOperator                   string = "eq"
	NotEqualOperator                string = "neq"
	GreaterThanOperator             string = "gt"
	GreaterThanOrEqualOperator      string = "gte"
	LessThanOperator                string = "lt"
	LessThanOrEqualOperator         string = "lte"
	ContainsOperator                string = "contains"
	NotContainsOperator             string = "ncontains"
	InOperator                      string = "in"
	NotInOperator                   string = "nin"
	StartsWithOperator              string = "startswith"
	EndsWithOperator                string = "endswith"
	InBetweenOperator               string = "inbetween"
	StartsWithCaseSensitiveOperator string = "startswith_cs"
	IsNullOperator                  string = "is_null"
	AndOperator                     string = "and"
	OrOperator                      string = "or"
	ArrayContainsOperator           string = "array_contains"
	ArrayInOperator                 string = "array_in"
)

var (
	OperatorEqual                   models.Operator = "="
	OperatorNotEqual                models.Operator = "!="
	OperatorGreaterThan             models.Operator = ">"
	OperatorLessThan                models.Operator = "<"
	OperatorGreaterThanOrEqual      models.Operator = ">="
	OperatorLessThanOrEqual         models.Operator = "<="
	OperatorIn                      models.Operator = "IN"
	OperatorNotIn                   models.Operator = "NOT IN"
	OperatorContains                models.Operator = "LIKE"
	OperatorNotContains             models.Operator = "NOT LIKE"
	OperatorStartsWith              models.Operator = "LIKE"
	OperatorEndsWith                models.Operator = "LIKE"
	OperatorInBetween               models.Operator = "BETWEEN"
	OperatorAnd                     models.Operator = "AND"
	OperatorOr                      models.Operator = "OR"
	OperatorStartsWithCaseSensitive models.Operator = "LIKE"
	OperatorIsNull                  models.Operator = "IS NULL"
)

var SqlOperatorMap = map[string]models.Operator{
	EqualOperator:                   OperatorEqual,
	NotEqualOperator:                OperatorNotEqual,
	GreaterThanOperator:             OperatorGreaterThan,
	GreaterThanOrEqualOperator:      OperatorGreaterThanOrEqual,
	LessThanOperator:                OperatorLessThan,
	LessThanOrEqualOperator:         OperatorLessThanOrEqual,
	InOperator:                      OperatorIn,
	NotInOperator:                   OperatorNotIn,
	ContainsOperator:                OperatorContains,
	NotContainsOperator:             OperatorNotContains,
	StartsWithOperator:              OperatorStartsWith,
	EndsWithOperator:                OperatorEndsWith,
	InBetweenOperator:               OperatorInBetween,
	AndOperator:                     OperatorAnd,
	OrOperator:                      OperatorOr,
	StartsWithCaseSensitiveOperator: OperatorStartsWithCaseSensitive,
	IsNullOperator:                  OperatorIsNull,
}

const (
	ZampDataset = "zamp_"
	AllColumns  = " * "
)

const (
	SelectStatement  = "SELECT "
	UpdateStatement  = "UPDATE "
	FromStatement    = " FROM "
	AsStatement      = " AS "
	WhereStatement   = " WHERE "
	GroupByStatement = " GROUP BY "
	OrderByStatement = " ORDER BY "
)

const (
	WindowStatement      = " OVER ( "
	PartitionByStatement = " PARTITION BY "
)

type WindowFunction string

const (
	WindowFunctionRowNumber WindowFunction = "ROW_NUMBER()"
)
