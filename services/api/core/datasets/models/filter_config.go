package models

import (
	dataplatformDataTypesConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
)

type FilterConfig struct {
	Column   string                                   `json:"column"`
	Alias    *string                                  `json:"alias"`
	Type     string                                   `json:"type"`
	DataType *dataplatformDataTypesConstants.Datatype `json:"datatype"`
	Options  []interface{}                            `json:"options"`
	Metadata map[string]interface{}                   `json:"metadata"`
}

type DatasetMetadataConfig struct {
	dataplatformDataModels.DatasetConfig
	DatabricksConfig dataplatformDataModels.DatabricksConfig `json:"databricks_config"`
	DisplayConfig    []DisplayConfig                         `json:"display_config"`
}

type FileImportConfig struct {
	IsFileImportEnabled bool                   `json:"is_file_import_enabled"`
	BronzeSourceBucket  string                 `json:"bronze_source_bucket"`
	BronzeSourcePath    string                 `json:"bronze_source_path"`
	BronzeSourceConfig  map[string]interface{} `json:"bronze_source_config"`
}

type DisplayConfig struct {
	Column        string                 `json:"column"`
	Alias         *string                `json:"alias,omitempty"`
	IsHidden      bool                   `json:"is_hidden"`
	IsEditable    bool                   `json:"is_editable"`
	Type          *string                `json:"type,omitempty"`
	Config        map[string]interface{} `json:"config,omitempty"`
	SqlExpression *string                `json:"sql_expression,omitempty"`
}

type MVConfig struct {
	Query            string            `json:"query"`
	QueryParams      map[string]string `json:"query_params"`
	ParentDatasetIds []string          `json:"parent_dataset_ids"`
}

type UpdateDatasetParams struct {
	Title            *string                               `json:"title,omitempty"`
	Description      *string                               `json:"description,omitempty"`
	Type             *string                               `json:"type,omitempty"`
	DatasetConfig    *dataplatformDataModels.DatasetConfig `json:"dataset_config,omitempty"`
	DisplayConfig    *[]DisplayConfig                      `json:"display_config,omitempty"`
	FileImportConfig *FileImportConfig                     `json:"file_import_config,omitempty"`
}

type CacheFilterConfig struct {
	FilterConfig []FilterConfig
	DatsetConfig map[string]interface{}
}
