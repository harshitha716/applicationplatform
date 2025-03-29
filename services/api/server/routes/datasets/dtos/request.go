package dtos

import (
	"fmt"
	"strings"

	dataplatformdataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	datasetConstants "github.com/Zampfi/application-platform/services/api/core/datasets/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type GetDataRequest struct {
	Filters         datasetmodels.FilterModel   `json:"filters"`
	Aggregations    []datasetmodels.Aggregation `json:"aggregations"`
	GroupBy         []datasetmodels.GroupBy     `json:"group_by"`
	OrderBy         []datasetmodels.OrderBy     `json:"order_by"`
	GetTotalRecords bool                        `json:"get_total_records,omitempty"`
	Pagination      *datasetmodels.Pagination   `json:"pagination,omitempty"`
	FxCurrency      *string                     `json:"fx_currency,omitempty"`
}

func (g *GetDataRequest) ToModel() datasetmodels.DatasetParams {
	pagination := g.Pagination
	if pagination == nil {
		pagination = &datasetmodels.Pagination{}
	}
	if pagination.Page <= 0 {
		pagination.Page = datasetConstants.DefaultPaginationPage
	}

	if pagination.PageSize <= 0 {
		pagination.PageSize = datasetConstants.DefaultPaginationPageSize
	}

	if pagination.PageSize > datasetConstants.DefaultMaxPaginationPageSize {
		pagination.PageSize = datasetConstants.DefaultMaxPaginationPageSize
	}

	return datasetmodels.DatasetParams{
		Filters:      g.Filters,
		Aggregations: g.Aggregations,
		GroupBy:      g.GroupBy,
		OrderBy:      g.OrderBy,
		CountAll:     g.GetTotalRecords,
		Pagination:   pagination,
		FxCurrency:   g.FxCurrency,
	}
}

type RegisterDatasetRequest struct {
	DatasetTitle       string                                  `json:"dataset_title"`
	DatasetDescription *string                                 `json:"dataset_description"`
	DatasetType        storemodels.DatasetType                 `json:"dataset_type"`
	DatasetConfig      dataplatformDataModels.DatasetConfig    `json:"dataset_config"`
	DatabricksConfig   dataplatformDataModels.DatabricksConfig `json:"databricks_config"`
	DisplayConfig      []datasetmodels.DisplayConfig           `json:"display_config"`
	MVConfig           *datasetmodels.MVConfig                 `json:"mv_config,omitempty"`
	Provider           dataplatformdataconstants.Provider      `json:"provider"`
}

func (r *RegisterDatasetRequest) ToModel() datasetmodels.DatasetCreationInfo {
	return datasetmodels.DatasetCreationInfo{
		DatasetTitle:       r.DatasetTitle,
		DatasetDescription: r.DatasetDescription,
		DatasetType:        r.DatasetType,
		DatasetConfig:      r.DatasetConfig,
		DatabricksConfig:   r.DatabricksConfig,
		DisplayConfig:      r.DisplayConfig,
		MVConfig:           r.MVConfig,
		Provider:           r.Provider,
	}
}

type UpdateDatasetDataRequest struct {
	Filters         datasetmodels.FilterModel  `json:"filters"`
	Update          datasetmodels.UpdateColumn `json:"update"`
	SaveAsRule      bool                       `json:"save_as_rule"`
	RuleTitle       *string                    `json:"rule_title"`
	RuleDescription *string                    `json:"rule_description"`
}

func (u *UpdateDatasetDataRequest) ToModel(userId uuid.UUID) datasetmodels.UpdateDatasetDataParams {
	var ruleTitle string
	var ruleDescription string

	if u.RuleTitle != nil {
		ruleTitle = *u.RuleTitle
	}

	if u.RuleDescription != nil {
		ruleDescription = *u.RuleDescription
	}

	var sourceType datasetConstants.UpdateColumnSourceType
	var sourceId uuid.UUID
	if u.SaveAsRule {
		sourceType = datasetConstants.UpdateColumnSourceTypeRule
		sourceId = uuid.New()
	} else {
		sourceType = datasetConstants.UpdateColumnSourceTypeUser
		sourceId = userId
	}

	return datasetmodels.UpdateDatasetDataParams{
		Filters:         u.Filters,
		Update:          u.Update,
		SourceType:      sourceType,
		SourceId:        sourceId,
		UserId:          userId,
		RuleTitle:       ruleTitle,
		RuleDescription: ruleDescription,
	}
}

type DatasetActionQueryParams struct {
	ActionIds  []string    `json:"action_ids"`
	ActionType []string    `json:"action_type"`
	ActionBy   []uuid.UUID `json:"action_by"`
	Status     []string    `json:"status"`
}

func (d *DatasetActionQueryParams) ToModel(datasetId uuid.UUID) storemodels.DatasetActionFilters {
	return storemodels.DatasetActionFilters{
		DatasetIds: []uuid.UUID{datasetId},
		ActionIds:  d.ActionIds,
		ActionType: d.ActionType,
		ActionBy:   d.ActionBy,
		Status:     d.Status,
	}
}

type UpdateAudienceRoleRequest struct {
	AudiencId uuid.UUID `json:"audience_id"`
	Role      string    `json:"role"`
}

type DeleteAudienceRoleRequest struct {
	AudiencId uuid.UUID `json:"audience_id"`
}

type AddAudienceRequest struct {
	AudienceType string    `json:"audience_type"`
	AudienceId   uuid.UUID `json:"audience_id"`
	Role         string    `json:"role"`
}

type BulkAddAudienceRequest struct {
	Audiences []AddAudienceRequest `json:"audiences"`
}

type UpdateDatasetRequest struct {
	Title            *string                               `json:"title,omitempty"`
	Description      *string                               `json:"description,omitempty"`
	Type             *string                               `json:"type,omitempty"`
	DatasetConfig    *dataplatformDataModels.DatasetConfig `json:"dataset_config,omitempty"`
	DisplayConfig    *[]datasetmodels.DisplayConfig        `json:"display_config,omitempty"`
	FileImportConfig *datasetmodels.FileImportConfig       `json:"file_import_config,omitempty"`
}

func (u *UpdateDatasetRequest) ToModel() datasetmodels.UpdateDatasetParams {
	return datasetmodels.UpdateDatasetParams{
		Title:            u.Title,
		Description:      u.Description,
		Type:             u.Type,
		DatasetConfig:    u.DatasetConfig,
		DisplayConfig:    u.DisplayConfig,
		FileImportConfig: u.FileImportConfig,
	}
}

type DatasetColumnRequest struct {
	DatasetId string   `json:"dataset_id"`
	Columns   []string `json:"columns"`
}

type GetRulesByDatasetColumnsRequest []DatasetColumnRequest

func (g GetRulesByDatasetColumnsRequest) ToModel() []storemodels.DatasetColumn {
	var datasetColumns []storemodels.DatasetColumn
	for _, dc := range g {
		datasetId, _ := uuid.Parse(dc.DatasetId)
		datasetColumns = append(datasetColumns, storemodels.DatasetColumn{
			DatasetId: datasetId,
			Columns:   dc.Columns,
		})
	}
	return datasetColumns
}

type InitiateDatasetFileImportRequest struct {
	FileName string               `json:"file_name"`
	FileType storemodels.FileType `json:"file_type"`
}

func (r *InitiateDatasetFileImportRequest) Validate() error {
	if strings.TrimSpace(r.FileName) == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	validFileTypes := map[storemodels.FileType]bool{
		storemodels.FileTypeCSV:     true,
		storemodels.FileTypeXLSX:    true,
		storemodels.FileTypeXLS:     true,
		storemodels.FileTypeParquet: true,
	}
	if !validFileTypes[r.FileType] {
		return fmt.Errorf("invalid file type. Supported types: csv, xlsx, xls, parquet")
	}

	return nil
}

type AckDatasetFileImportRequest struct {
	DatasetId uuid.UUID `json:"dataset_id"`
}

type ConfirmDatasetImportRequest struct {
	DatasetId uuid.UUID `json:"dataset_id"`
}

type SetDatasetDisplayConfigRequest struct {
	DisplayConfig []datasetmodels.DisplayConfig `json:"display_config"`
}
