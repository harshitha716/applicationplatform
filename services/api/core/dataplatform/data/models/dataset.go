package models

import (
	"time"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"
)

type Dataset struct {
	Id                    string    `json:"id"`
	MerchantId            string    `json:"merchant_id"`
	DatabricksTableName   string    `json:"databricks_table_name"`
	DatabricksFQTableName string    `json:"databricks_fq_table_name"`
	DatasetConfig         string    `json:"dataset_config"`
	DatabricksSchema      string    `json:"databricks_schema"`
	DatabricksConfig      string    `json:"databricks_config"`
	DatabricksStats       string    `json:"databricks_stats"`
	PinotTableName        string    `json:"pinot_table_name"`
	PinotSchema           string    `json:"pinot_schema"`
	PinotConfig           string    `json:"pinot_config"`
	PinotStats            string    `json:"pinot_stats"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	IsDeleted             bool      `json:"is_deleted"`
	DeletedAt             time.Time `json:"deleted_at"`
}

type QueryMetadata struct {
	Params     map[string]string `json:"params"`
	TableNames []string          `json:"tableNames"`
}

type CustomColumnGroup struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type ComputedColumn struct {
	OutputColumn   string `json:"output_column"`
	Expression     string `json:"expression"`
	ForceRecompute bool   `json:"force_recompute"`
}

type Rule struct {
	Id           string                 `json:"id"`
	Priority     int                    `json:"priority"`
	ValueToApply string                 `json:"value_to_apply"`
	SqlCondition string                 `json:"sql_condition"`
	SqlArgs      map[string]interface{} `json:"sql_args"`
}

type DatasetConfig struct {
	Columns            map[string]DatasetColumnConfig `json:"columns"`
	CustomColumnGroups []CustomColumnGroup            `json:"custom_column_groups"`
	Rules              map[string][]Rule              `json:"rules"`
	ComputedColumns    []ComputedColumn               `json:"computed_columns"`
}

type DatasetColumnConfig struct {
	CustomType       constants.DatabricksColumnCustomType `json:"custom_type"`
	CustomTypeConfig map[string]string                    `json:"custom_type_config"`
}

type DatabricksConfig struct {
	DedupColumns     []string `json:"dedup_columns"`
	OrderByColumn    string   `json:"order_by_column"`
	PartitionColumns []string `json:"partition_columns"`
	ClusterColumns   []string `json:"cluster_columns"`
}
