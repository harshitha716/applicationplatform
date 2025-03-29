package models

type ColumnMetadata struct {
	Type string `json:"type"`
}

type ColumnStats struct {
	DistinctCount int `json:"distinct_count"`
}

type DatasetStats struct {
	ColumnStats map[string]ColumnStats `json:"columns"`
}

type DatasetSchemaDetails struct {
	Columns map[string]ColumnMetadata `json:"columns"`
}

type DatasetMetadata struct {
	Schema map[string]ColumnMetadata `json:"schema"`
	Stats  DatasetStats              `json:"stats"`
}

type InternalDatasetMetadata struct {
	Schema DatasetSchemaDetails `json:"schema"`
	Stats  DatasetStats         `json:"stats"`
}

type ProviderLevelDatasetMetadata struct {
	Databricks InternalDatasetMetadata `json:"databricks"`
	Pinot      InternalDatasetMetadata `json:"pinot"`
}
