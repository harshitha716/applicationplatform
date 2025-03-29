package dtos

import (
	"encoding/json"

	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	"github.com/google/uuid"
)

type WidgetInstanceDataResponse struct {
	Status      string   `json:"status"`
	Periodicity string   `json:"periodicity"`
	Currency    *string  `json:"currency"`
	Result      []Result `json:"result"`
}

type Result struct {
	RowCount int              `json:"rowcount"`
	Columns  []Column         `json:"columns"`
	Data     []map[string]any `json:"data"`
}

type Column struct {
	ColumnName string `json:"column_name"`
	ColumnType string `json:"column_type"`
}

// Remove the pointer receiver method and replace with constructor
func NewWidgetInstanceDataResponse(qrs []datasetmodels.DatasetData, periodicity string, currency *string) *WidgetInstanceDataResponse {
	widgetInstanceDataResponse := WidgetInstanceDataResponse{
		Status:      "success",
		Periodicity: periodicity,
		Currency:    currency,
		Result:      []Result{},
	}

	for _, qr := range qrs {
		columns := make([]Column, len(qr.Columns))
		for i, col := range qr.Columns {
			columns[i] = Column{
				ColumnName: col.Name,
				ColumnType: col.DatabaseType,
			}
		}

		widgetInstanceDataResponse.Result = append(widgetInstanceDataResponse.Result, Result{
			RowCount: len(qr.Rows),
			Columns:  columns,
			Data:     qr.Rows,
		})
	}

	return &widgetInstanceDataResponse
}

type WidgetInstanceResponse struct {
	ID            uuid.UUID       `json:"widget_instance_id"`
	WidgetType    string          `json:"widget_type"`
	SheetID       uuid.UUID       `json:"sheet_id"`
	Title         string          `json:"title"`
	DataMappings  json.RawMessage `json:"data_mappings"`
	DisplayConfig json.RawMessage `json:"display_config,omitempty"`
}

type WidgetInstanceMappingDTO struct {
	Version  string                `json:"version"`
	Mappings []DataMappingFieldDTO `json:"mappings"`
}

type DataMappingFieldDTO struct {
	DatasetID string                     `json:"dataset_id"`
	Fields    map[string]json.RawMessage `json:"fields"`
}

type DefaultFilterValue struct {
	Column   string        `json:"column,omitempty"`
	Operator string        `json:"operator,omitempty"`
	Value    []interface{} `json:"value,omitempty"`
	Type     *string       `json:"type,omitempty"`
}

func NewWidgetInstanceResponse(wi *widgetmodels.WidgetInstance) (*WidgetInstanceResponse, error) {
	mappings, err := json.Marshal(wi.DataMappings)
	if err != nil {
		return nil, err
	}

	response := &WidgetInstanceResponse{
		ID:           wi.ID,
		WidgetType:   wi.WidgetType,
		SheetID:      wi.SheetID,
		Title:        wi.Title,
		DataMappings: mappings,
	}

	if wi.DisplayConfig != nil {
		response.DisplayConfig = *wi.DisplayConfig
	}

	return response, nil
}

func transformDefaultFilter(filter datasetmodels.FilterModel) []DefaultFilterValue {
	// Assuming a flat hierarchy of filters
	// TODO: Evaluate if this is the best approach
	if len(filter.Conditions) == 0 {
		return nil
	}

	defaultFilters := make([]DefaultFilterValue, 0)
	for _, condition := range filter.Conditions {
		defaultFilters = append(defaultFilters, DefaultFilterValue{
			Operator: condition.Operator,
			Value:    []interface{}{condition.Value},
			Column:   condition.Column,
		})
	}
	return defaultFilters
}
