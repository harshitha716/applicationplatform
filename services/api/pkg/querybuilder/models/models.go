package models

import (
	"fmt"

	dataplaformconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/errors"
)

type (
	OrderType           string
	Operator            string
	LogicalOperator     string
	AggregationFunction string
)

type WindowFunction string

const (
	WindowFunctionRowNumber WindowFunction = "ROW_NUMBER()"
)

type QueryConfig struct {
	TableConfig  TableConfig    `json:"table_config"`
	Subquery     *QueryConfig   `json:"subquery"`
	Windows      []WindowConfig `json:"windows"`
	Filters      FilterModel    `json:"filters"`
	Aggregations []Aggregation  `json:"aggregations"`
	GroupBy      []GroupBy      `json:"group_by"`
	OrderBy      []OrderBy      `json:"order_by"`
	CountAll     bool           `json:"count_all"`
	Pagination   *Pagination    `json:"pagination"`
}

type TableConfig struct {
	DatasetId string         `json:"dataset_id"`
	Columns   []ColumnConfig `json:"columns"`
}

type ColumnConfig struct {
	Column           string                              `json:"column"`
	Datatype         *dataplatformdataconstants.Datatype `json:"datatype"`
	CustomDataConfig *CustomDataTypeConfig               `json:"custom_data_config"`
	Alias            *string                             `json:"alias"`
}

func (c *ColumnConfig) GetSelectColumn() (string, error) {
	if c.CustomDataConfig != nil {
		switch c.CustomDataConfig.Type {
		case dataplaformconstants.DatabricksColumnCustomTypeAmount:
			if c.CustomDataConfig.Config == nil {
				return "", errors.ErrInvalidCustomDataTypeConfig
			}
			return c.CustomDataConfig.Config.GetSelectColumn(), nil
		default:
			return "", errors.ErrInvalidCustomDataType
		}
	}
	if c.Alias != nil {
		return fmt.Sprintf("%s AS \"%s\"", c.Column, *c.Alias), nil
	}
	return c.Column, nil
}

func (c *ColumnConfig) GetGroupByColumn() (string, error) {
	if c.Datatype == nil {
		return "", errors.ErrInvalidDataType
	}

	switch *c.Datatype {
	case dataplatformdataconstants.ArrayOfStringDataType:
		if c.Alias != nil {
			return fmt.Sprintf("\"%s\"", *c.Alias), nil
		}
		return fmt.Sprintf("unnest(%s)", c.Column), nil
	default:
		if c.CustomDataConfig != nil {
			switch c.CustomDataConfig.Type {
			case dataplaformconstants.DatabricksColumnCustomTypeAmount:
				if c.CustomDataConfig.Config == nil {
					return "", errors.ErrInvalidCustomDataTypeConfig
				}
				return c.CustomDataConfig.Config.GetGroupByColumn(), nil
			default:
				return "", errors.ErrInvalidCustomDataType
			}
		}
		if c.Alias != nil {
			return fmt.Sprintf("\"%s\"", *c.Alias), nil
		}
		return c.Column, nil
	}
}

func (c *ColumnConfig) GetOrderByColumn() (string, error) {
	if c.Alias != nil {
		return fmt.Sprintf("\"%s\"", *c.Alias), nil
	}
	return c.Column, nil
}

func (c *ColumnConfig) GetAggregationColumn() (string, error) {
	if c.CustomDataConfig != nil {
		switch c.CustomDataConfig.Type {
		case dataplaformconstants.DatabricksColumnCustomTypeAmount:
			if c.CustomDataConfig.Config == nil {
				return "", errors.ErrInvalidCustomDataTypeConfig
			}
			return c.CustomDataConfig.Config.GetAggregationColumn(), nil
		default:
			return "", errors.ErrInvalidCustomDataType
		}
	}
	return c.Column, nil
}

func (c *ColumnConfig) GetGroupBySelectColumn() (string, error) {
	if c.Datatype == nil {
		return "", errors.ErrInvalidDataType
	}

	switch *c.Datatype {
	case dataplatformdataconstants.ArrayOfStringDataType:
		if c.Alias != nil {
			return fmt.Sprintf("unnest(%s) AS \"%s\"", c.Column, *c.Alias), nil
		}
		return fmt.Sprintf("unnest(%s)", c.Column), nil
	default:
		if c.CustomDataConfig != nil {
			switch c.CustomDataConfig.Type {
			case dataplaformconstants.DatabricksColumnCustomTypeAmount:
				if c.CustomDataConfig.Config == nil {
					return "", errors.ErrInvalidCustomDataTypeConfig
				}
				return c.CustomDataConfig.Config.GetSelectColumn(), nil
			default:
				return "", errors.ErrInvalidCustomDataType
			}
		}
		if c.Alias != nil {
			return fmt.Sprintf("%s AS \"%s\"", c.Column, *c.Alias), nil
		}
		return c.Column, nil
	}
}

type FilterModel struct {
	LogicalOperator LogicalOperator `json:"logical_operator"`
	Conditions      []Filter        `json:"conditions"`
}

type Filter struct {
	LogicalOperator *LogicalOperator `json:"logical_operator"`
	Column          ColumnConfig     `json:"column"`
	Operator        string           `json:"operator"`
	Value           interface{}      `json:"value"`
	Conditions      []Filter         `json:"conditions"`
}

type Aggregation struct {
	Column   ColumnConfig        `json:"column"`
	Alias    string              `json:"alias"`
	Function AggregationFunction `json:"function"`
}

type GroupBy struct {
	Column ColumnConfig `json:"column"`
}

type OrderBy struct {
	Column ColumnConfig `json:"column"`
	Order  OrderType    `json:"order"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type WindowConfig struct {
	Function    WindowFunction `json:"function"`
	PartitionBy []ColumnConfig `json:"partition_by"`
	OrderBy     []OrderBy      `json:"order_by"`
	Alias       string         `json:"alias"`
}
