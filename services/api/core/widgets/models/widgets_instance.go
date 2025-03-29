package models

import (
	"encoding/json"
	"fmt"

	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetconstants "github.com/Zampfi/application-platform/services/api/core/widgets/constants"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

/*
	Add source_datasets field to mappings with a structure like this
	{
	"source_datasets": {
		"datasets": [
		{ "id": "dataset_A", "alias": "A" },
		{ "id": "dataset_B", "alias": "B" },
		{ "id": "dataset_C", "alias": "C" }
		],
		"joins": [
		{
			"left_dataset_alias": "A",
			"right_dataset_alias": "B",
			"join_type": "left",
			"conditions": [
			{ "left_column": "A_id", "right_column": "B_id" }
			]
		},
		{
			"left_dataset_alias": "A",
			"right_dataset_alias": "C",
			"join_type": "inner",
			"conditions": [
			{ "left_column": "A_something", "right_column": "C_something" }
			]
		}
		]
	}
}

*/

type DataMappingVersion string

const (
	DataMappingVersion1 DataMappingVersion = "1"
)

type WidgetInstance struct {
	ID            uuid.UUID        `json:"widget_instance_id"`
	WidgetType    string           `json:"widget_type"`
	SheetID       uuid.UUID        `json:"sheet_id"`
	Title         string           `json:"title"`
	DataMappings  DataMappings     `json:"data_mappings"`
	DisplayConfig *json.RawMessage `json:"display_config,omitempty"`
}

type DataMappings struct {
	Version  DataMappingVersion  `json:"version"`
	Mappings []DataMappingFields `json:"mappings"`
}

type DataMappingFields struct {
	DatasetID      string                     `json:"dataset_id"`
	Ref            string                     `json:"ref"`
	SourceDatasets *SourceDatasets            `json:"source_datasets,omitempty"`
	Fields         map[string][]Field         `json:"fields"`
	DefaultFilters *datasetmodels.FilterModel `json:"default_filters,omitempty"`
	SortBy         []SortBy                   `json:"sort_by,omitempty"`
}

type SortBy struct {
	Column string `json:"column"`
	Order  string `json:"order"`
	Alias  string `json:"alias,omitempty"`
}

func (s *SortBy) GetColumn() string {
	if s.Alias != "" {
		return s.Alias
	}
	return s.Column
}

func (d *DataMappingFields) GetRef() string {
	return d.Ref
}

type SourceDatasets struct {
	Datasets []SourceDataset `json:"datasets"`
	Joins    []SourceJoin    `json:"joins"`
}

type SourceDataset struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

type SourceJoin struct {
	LeftDatasetAlias  string                `json:"left_dataset_alias"`
	RightDatasetAlias string                `json:"right_dataset_alias"`
	JoinType          string                `json:"join_type"`
	Conditions        []SourceJoinCondition `json:"conditions"`
}

type SourceJoinCondition struct {
	LeftColumn  string `json:"left_column"`
	RightColumn string `json:"right_column"`
}

type Field struct {
	Column                  string                             `json:"column"`
	Aggregation             string                             `json:"aggregation,omitempty"`
	Type                    widgetconstants.UserFacingDatatype `json:"type"`
	DrilldownFilterType     string                             `json:"drilldown_filter_type,omitempty"`
	DrilldownFilterOperator string                             `json:"drilldown_filter_operator,omitempty"`
	FieldType               string                             `json:"field_type"`
	Alias                   string                             `json:"alias,omitempty"`
	Expression              string                             `json:"expression,omitempty"`
	SortBy                  []SortBy                           `json:"sort_by,omitempty"`
}

func (f *Field) GetAlias() *string {
	if f.Alias != "" {
		return &f.Alias
	}
	return &f.Column
}

func (f *Field) GetExpression() string {
	if f.Expression != "" {
		return f.Expression
	}
	return f.Column
}

func (wi *WidgetInstance) FromDB(dbModel *dbmodels.WidgetInstance) error {
	wi.ID = dbModel.ID
	wi.WidgetType = dbModel.WidgetType
	wi.SheetID = dbModel.SheetID
	wi.Title = dbModel.Title

	if err := json.Unmarshal(dbModel.DataMappings, &wi.DataMappings); err != nil {
		return fmt.Errorf("unmarshal data mappings: %w", err)
	}

	if dbModel.DisplayConfig != nil {
		wi.DisplayConfig = dbModel.DisplayConfig
	}

	return nil
}

func (wi *WidgetInstance) ToDB() (*dbmodels.WidgetInstance, error) {

	dbModel := &dbmodels.WidgetInstance{
		ID:            wi.ID,
		SheetID:       wi.SheetID,
		Title:         wi.Title,
		WidgetType:    wi.WidgetType,
		DisplayConfig: wi.DisplayConfig,
	}

	mappings, err := json.Marshal(wi.DataMappings)
	if err != nil {
		return nil, fmt.Errorf("marshal data mappings: %w", err)
	}
	dbModel.DataMappings = mappings

	return dbModel, nil
}

type CreateWidgetInstancePayload struct {
	WidgetType   string `form:"label=Widget Type"`
	SheetID      string `form:"label=Sheet ID"`
	Title        string `form:"label=Title"`
	DataMappings string `form:"label=Data Mappings"`
}

func (c *CreateWidgetInstancePayload) ToModel() (*WidgetInstance, error) {
	var dataMappings DataMappings
	if err := json.Unmarshal([]byte(c.DataMappings), &dataMappings); err != nil {
		return nil, fmt.Errorf("invalid data mappings JSON: %w", err)
	}

	sheetID, err := uuid.Parse(c.SheetID)
	if err != nil {
		return nil, fmt.Errorf("invalid sheet ID: %w", err)
	}

	return &WidgetInstance{
		WidgetType:   c.WidgetType,
		SheetID:      sheetID,
		Title:        c.Title,
		DataMappings: dataMappings,
	}, nil
}

type UpdateWidgetInstancePayload struct {
	WidgetInstanceID string  `form:"label=Widget Instance ID"`
	WidgetType       *string `form:"label=Widget Type,omitempty"`
	SheetID          *string `form:"label=Sheet ID,omitempty"`
	Title            *string `form:"label=Title,omitempty"`
	DataMappings     *string `form:"label=Data Mappings,omitempty"`
	DisplayConfig    *string `form:"label=Display Config,omitempty"`
}

func (u *UpdateWidgetInstancePayload) ToModel() (*WidgetInstance, error) {
	var dataMappings DataMappings

	if mappings := stringPtr(u.DataMappings); mappings != nil {
		if err := json.Unmarshal([]byte(*mappings), &dataMappings); err != nil {
			return nil, fmt.Errorf("invalid data mappings JSON: %w", err)
		}
	}

	widgetType := ""
	if wType := stringPtr(u.WidgetType); wType != nil {
		widgetType = *wType
	}

	sheetID := uuid.Nil
	if sID := stringPtr(u.SheetID); sID != nil {
		var err error
		sheetID, err = uuid.Parse(*sID)
		if err != nil {
			return nil, fmt.Errorf("invalid sheet ID: %w", err)
		}
	}

	title := ""
	if t := stringPtr(u.Title); t != nil {
		title = *t
	}

	var displayConfig *json.RawMessage
	if config := stringPtr(u.DisplayConfig); config != nil && *config != "" {
		rawMsg := json.RawMessage(*config)
		displayConfig = &rawMsg
	}

	widgetInstanceID, err := uuid.Parse(u.WidgetInstanceID)
	if err != nil {
		return nil, fmt.Errorf("invalid widget instance ID: %w", err)
	}

	return &WidgetInstance{
		ID:            widgetInstanceID,
		WidgetType:    widgetType,
		SheetID:       sheetID,
		Title:         title,
		DataMappings:  dataMappings,
		DisplayConfig: displayConfig,
	}, nil
}

func stringPtr(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

type GetWidgetInstanceDataQueryParams struct {
	Filters     []WidgetFilters `json:"filters"`
	TimeColumns []ColumnMapping `json:"time_columns"`
	Periodicity *string         `json:"periodicity,omitempty"`
	Currency    *string         `json:"currency,omitempty"`
}

type ColumnMapping struct {
	DatasetID string `json:"dataset_id"`
	Column    string `json:"column"`
}

type WidgetFilters struct {
	DatasetID  string                    `json:"dataset_id"`
	Filters    datasetmodels.FilterModel `json:"filters"`
	Pagination *PaginationParams         `json:"pagination,omitempty"`
}

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type DatasetBuilderParams struct {
	Filters     map[string]WidgetFilters
	TimeColumns map[string]string
	Periodicity *string
	Currency    *string
}

type GetDataByDatasetIDParams struct {
	DatasetID string
	Params    datasetmodels.DatasetParams
}
