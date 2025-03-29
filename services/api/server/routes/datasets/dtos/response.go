package dtos

import (
	"time"

	dataplatformDataTypesConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"

	"github.com/google/uuid"
)

type GetDataResponse struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Data        DatasetData `json:"data"`
}

type DatasetData struct {
	Rows          Rows             `json:"rows"`
	Columns       []ColumnMetadata `json:"columns"`
	TotalCount    *int64           `json:"total_count,omitempty"`
	DatasetConfig DatasetConfig    `json:"config"`
}

type DatasetConfig struct {
	IsDrilldownEnabled bool `json:"is_drilldown_enabled"`
}

type Rows []map[string]interface{}

type ColumnMetadata struct {
	Name         string `json:"name"`
	DatabaseType string `json:"data_type"`
}

func (qr *DatasetData) FromModel(model datasetmodels.DatasetData) {
	qr.Rows = Rows(model.Rows)
	qr.Columns = []ColumnMetadata{}
	for _, column := range model.Columns {
		qr.Columns = append(qr.Columns, ColumnMetadata{
			Name:         column.Name,
			DatabaseType: column.DatabaseType,
		})
	}
	if model.TotalCount != nil {
		qr.TotalCount = model.TotalCount
	}
	qr.DatasetConfig = DatasetConfig{
		IsDrilldownEnabled: model.DatasetConfig.IsDrilldownEnabled,
	}
}

type ParentDatasetInfo struct {
	ParentDatasets []DatasetInfo `json:"tabs"`
}

type DatasetInfo struct {
	DatasetId          string                    `json:"dataset_id"`
	DatasetTitle       string                    `json:"dataset_title"`
	DatasetDescription string                    `json:"dataset_description"`
	DatasetType        string                    `json:"dataset_type"`
	Filters            datasetmodels.FilterModel `json:"filters"`
}

func (dr *ParentDatasetInfo) FromModel(model datasetmodels.ParentDatasetInfo) {
	dr.ParentDatasets = []DatasetInfo{}
	for _, tab := range model.ParentDatasets {
		datasetInfo := DatasetInfo{
			DatasetId:          tab.DatasetId,
			DatasetTitle:       tab.DatasetTitle,
			DatasetDescription: tab.DatasetDescription,
			DatasetType:        string(tab.DatasetType),
			Filters:            tab.Filters,
		}

		dr.ParentDatasets = append(dr.ParentDatasets, datasetInfo)
	}
}

type FilterConfigResponse struct {
	Data   []FilterConfig         `json:"data"`
	Config map[string]interface{} `json:"config"`
}

type FilterConfig struct {
	Column   string                                   `json:"column"`
	Alias    *string                                  `json:"alias,omitempty"`
	Type     string                                   `json:"type"`
	DataType *dataplatformDataTypesConstants.Datatype `json:"datatype"`
	Options  []interface{}                            `json:"options"`
	Metadata map[string]interface{}                   `json:"metadata"`
}

func (f *FilterConfig) FromModel(model datasetmodels.FilterConfig) {
	f.Column = model.Column
	f.Alias = model.Alias
	f.Type = string(model.Type)
	f.DataType = model.DataType
	f.Options = model.Options
	f.Metadata = model.Metadata
}

type GetDatasetListingResponse struct {
	Datasets   []DatasetListing `json:"datasets"`
	TotalCount int64            `json:"total_count"`
}

type DatasetListing struct {
	Id             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	CreatedBy      uuid.UUID   `json:"created_by"`
	OrganizationId uuid.UUID   `json:"organization_id"`
	Metadata       interface{} `json:"metadata"`
}

func (d *DatasetListing) FromModel(model datasetmodels.Dataset) {
	d.Id = model.ID
	d.Title = model.Title
	d.Description = model.Description
	d.CreatedAt = model.CreatedAt
	d.UpdatedAt = model.UpdatedAt
	d.CreatedBy = model.CreatedBy
	d.OrganizationId = model.OrganizationId
	d.Metadata = model.Metadata
}

type DatasetAction struct {
	ActionId    string    `json:"action_id"`
	ActionType  string    `json:"action_type"`
	DatasetId   uuid.UUID `json:"dataset_id"`
	Status      string    `json:"status"`
	ActionBy    uuid.UUID `json:"action_by"`
	IsCompleted bool      `json:"is_completed"`
}

func (u *DatasetAction) FromModel(model datasetmodels.DatasetAction) {
	u.ActionId = model.ActionId
	u.ActionType = string(model.ActionType)
	u.DatasetId = model.DatasetId
	u.Status = string(model.Status)
	u.ActionBy = model.ActionBy
	u.IsCompleted = model.IsCompleted
}

type GetDatasetDisplayConfigResponse struct {
	DisplayConfig []datasetmodels.DisplayConfig `json:"display_config"`
}
