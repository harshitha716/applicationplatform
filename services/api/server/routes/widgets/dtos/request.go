package dtos

import (
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
)

type WidgetQueryParams struct {
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

func (w *WidgetQueryParams) ToModels() widgetmodels.GetWidgetInstanceDataQueryParams {
	widgetFilters := make([]widgetmodels.WidgetFilters, 0)
	for _, filter := range w.Filters {
		widgetFilter := widgetmodels.WidgetFilters{
			DatasetID: filter.DatasetID,
			Filters:   filter.Filters,
		}

		if filter.Pagination != nil {
			widgetFilter.Pagination = &widgetmodels.PaginationParams{
				Page:     filter.Pagination.Page,
				PageSize: filter.Pagination.PageSize,
			}
		}

		widgetFilters = append(widgetFilters, widgetFilter)
	}

	timeColumns := make([]widgetmodels.ColumnMapping, 0)
	for _, column := range w.TimeColumns {
		timeColumns = append(timeColumns, widgetmodels.ColumnMapping{
			DatasetID: column.DatasetID,
			Column:    column.Column,
		})
	}

	return widgetmodels.GetWidgetInstanceDataQueryParams{
		Filters:     widgetFilters,
		TimeColumns: timeColumns,
		Periodicity: w.Periodicity,
		Currency:    w.Currency,
	}
}
