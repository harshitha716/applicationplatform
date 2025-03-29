package models

import (
	"time"

	dataplatformconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	"github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	querybuildermodels "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
	"github.com/google/uuid"
)

type (
	OrderType           string
	Operator            string
	LogicalOperator     string
	AggregationFunction string
)

type DatasetParams struct {
	Columns         []ColumnConfig
	Subquery        *DatasetParams
	Windows         []WindowConfig
	Filters         FilterModel
	Aggregations    []Aggregation
	GroupBy         []GroupBy
	OrderBy         []OrderBy
	CountAll        bool
	Pagination      *Pagination
	FxCurrency      *string
	GetDatafromLake bool
}

type ColumnConfig struct {
	Column string
	Alias  *string
}

type UpdateColumn struct {
	Column string
	Value  interface{}
}

type FilterModel struct {
	LogicalOperator LogicalOperator `json:"logical_operator"`
	Conditions      []Filter        `json:"conditions"`
}

type Filter struct {
	LogicalOperator *LogicalOperator `json:"logical_operator"`
	Column          string           `json:"column"`
	Operator        string           `json:"operator"`
	Value           interface{}      `json:"value"`
	Conditions      []Filter         `json:"conditions"`
}

type Aggregation struct {
	Column   string              `json:"column"`
	Alias    string              `json:"alias"`
	Function AggregationFunction `json:"function"`
}

type GroupBy struct {
	Column string  `json:"column"`
	Alias  *string `json:"alias"`
}

type OrderBy struct {
	Column string    `json:"column"`
	Alias  *string   `json:"alias"`
	Order  OrderType `json:"order"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type WindowConfig struct {
	Function    string         `json:"function"`
	PartitionBy []ColumnConfig `json:"partition_by"`
	OrderBy     []OrderBy      `json:"order_by"`
	Alias       string         `json:"alias"`
}

type UpdateDatasetDataParams struct {
	Filters         FilterModel
	Update          UpdateColumn
	SourceType      constants.UpdateColumnSourceType
	SourceId        uuid.UUID
	UserId          uuid.UUID
	RuleTitle       string
	RuleDescription string
}

type DatasetData struct {
	dataplatformmodels.QueryResult
	TotalCount    *int64                `json:"total_count"`
	Title         string                `json:"title"`
	Description   *string               `json:"description"`
	DatasetConfig DatasetConfig         `json:"dataset_config"`
	Metadata      DatasetMetadataConfig `json:"metadata"`
}

type DatasetConfig struct {
	IsDrilldownEnabled bool `json:"is_drilldown_enabled"`
}

type DatasetInfo struct {
	DatasetId          string
	DatasetTitle       string
	DatasetDescription string
	DatasetType        dbmodels.DatasetType
	Filters            FilterModel
}

type ParentDatasetInfo struct {
	ParentDatasets []DatasetInfo
}

type DatsetListingParams struct {
	CreatedBy  []uuid.UUID
	Pagination *querybuildermodels.Pagination
	SortParams []DatasetListingSortParams
}

type DatasetListingSortParams struct {
	Column string `json:"column"`
	Desc   bool   `json:"desc"`
}

type DatasetCreationInfo struct {
	DatasetTitle       string
	DatasetDescription *string
	DatasetType        dbmodels.DatasetType
	DatasetConfig      dataplatformDataModels.DatasetConfig
	DatabricksConfig   dataplatformDataModels.DatabricksConfig
	DisplayConfig      []DisplayConfig
	MVConfig           *MVConfig
	Provider           dataplatformdataconstants.Provider
}

type CopyDatasetParams struct {
	DatasetId          string  `json:"dataset_id"`
	DatasetTitle       string  `json:"dataset_title"`
	DatasetDescription *string `json:"dataset_description"`
}

type UpdateColumnSource struct {
	SourceType constants.UpdateColumnSourceType `json:"source_type"`
	SourceId   uuid.UUID                        `json:"source_id"`
	UpdatedAt  time.Time                        `json:"updated_at"`
}

func (d *DatasetData) GetTagsColumns() (string, bool) {
	for columnName, columnType := range d.Metadata.Columns {
		if columnType.CustomType == (dataplatformconstants.DatabricksColumnCustomTypeTags) {
			for _, column := range d.Columns {
				if column.Name == columnName {
					return column.Name, true
				}
			}
		}
	}
	return "", false
}

type DatasetExportParams struct {
	QueryConfig DatasetParams `json:"query_config"`
	ExportPath  string        `json:"export_path"`
}

type ExportMetadata struct {
	FilePath string `json:"file_path"`
}
