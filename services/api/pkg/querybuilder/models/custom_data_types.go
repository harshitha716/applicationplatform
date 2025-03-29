package models

import (
	"fmt"

	dataplaformconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	datasetsconstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/errors"
)

type CustomDataTypeConfig struct {
	Type   dataplaformconstants.DatabricksColumnCustomType
	Config CustomDataTypeInterface
}

type CustomDataTypeInterface interface {
	Validate() error
	GetSelectColumn() string
	GetGroupByColumn() string
	GetFilterColumn() string
	GetAggregationColumn() string
}

func (c *CustomDataTypeConfig) ValidateConfig() error {
	switch c.Type {
	case dataplaformconstants.DatabricksColumnCustomTypeAmount:
		config, ok := c.Config.(*AmountCustomTypeConfig)
		if !ok {
			return errors.ErrInvalidCustomDataTypeConfig
		}
		return config.Validate()
	default:
		return errors.ErrInvalidCustomDataType
	}
}

type AmountCustomTypeConfig struct {
	CurrencyColumn string
	FxCurrency     string
	AmountColumn   string
}

func (c *AmountCustomTypeConfig) Validate() error {
	if c.CurrencyColumn == "" || c.AmountColumn == "" || c.FxCurrency == "" {
		return errors.ErrInvalidCustomDataType
	}
	return nil
}

func (c *AmountCustomTypeConfig) GetSelectColumn() string {
	return fmt.Sprintf("(%s%s->>'%s')::double AS \"%s\", '%s' AS \"%s\"", datasetsconstants.ZampFxColumnPrefix, c.AmountColumn, c.FxCurrency, c.AmountColumn, c.FxCurrency, c.CurrencyColumn)
}

func (c *AmountCustomTypeConfig) GetFilterColumn() string {
	return fmt.Sprintf("(%s%s->>'%s')::double", datasetsconstants.ZampFxColumnPrefix, c.AmountColumn, c.FxCurrency)
}

func (c *AmountCustomTypeConfig) GetGroupByColumn() string {
	return fmt.Sprintf("(%s%s->>'%s')::double", datasetsconstants.ZampFxColumnPrefix, c.AmountColumn, c.FxCurrency)
}

func (c *AmountCustomTypeConfig) GetAggregationColumn() string {
	return fmt.Sprintf("(%s%s->>'%s')::double", datasetsconstants.ZampFxColumnPrefix, c.AmountColumn, c.FxCurrency)
}
