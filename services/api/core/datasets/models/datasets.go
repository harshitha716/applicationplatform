package models

import (
	"time"

	dataplatformactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type Dataset struct {
	ID             uuid.UUID
	Title          string
	Description    string
	OrganizationId uuid.UUID
	CreatedBy      uuid.UUID
	Metadata       interface{}
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (d *Dataset) FromSchema(schema dbmodels.Dataset) {
	var description string
	if schema.Description != nil {
		description = *schema.Description
	}

	d.ID = schema.ID
	d.Title = schema.Title
	d.Description = description
	d.OrganizationId = schema.OrganizationId
	d.CreatedBy = schema.CreatedBy
	d.Metadata = schema.Metadata
	d.CreatedAt = schema.CreatedAt
	d.UpdatedAt = schema.UpdatedAt
}

type DatasetAction struct {
	ActionId    string                                   `json:"action_id"`
	ActionType  dataplatformactionconstants.ActionType   `json:"action_type"`
	DatasetId   uuid.UUID                                `json:"dataset_id"`
	Status      dataplatformactionconstants.ActionStatus `json:"status"`
	Config      interface{}                              `json:"config"`
	ActionBy    uuid.UUID                                `json:"action_by"`
	IsCompleted bool                                     `json:"is_completed"`
}

type AddDatasetAudiencePayload struct {
	AudienceType string
	AudienceId   uuid.UUID
	Privilege    string
}

type BulkAddDatasetAudiencePayload struct {
	Audiences []AddDatasetAudiencePayload
}

type AddDatasetAudienceError struct {
	AudienceId   uuid.UUID `json:"audience_id"`
	ErrorMessage string    `json:"error_message"`
}

type BulkAddDatasetAudienceErrors struct {
	Error     error
	Audiences []AddDatasetAudienceError `json:"audiences;omitempty"`
}

type FileImportWorkflowInitPayload struct {
	DatasetActionId     uuid.UUID `json:"dataset_action_id"`
	DatasetFileUploadId uuid.UUID `json:"dataset_file_upload_id"`
	DatasetId           uuid.UUID `json:"dataset_id"`
	FileUploadId        uuid.UUID `json:"file_upload_id"`
	UserId              uuid.UUID `json:"user_id"`
	OrganizationId      uuid.UUID `json:"organization_id"`
}

type FileImportWorkflowExitPayload struct {
	SourceFilePath      string                          `json:"file_path"`
	NormalizedFilePath  string                          `json:"base_normalized_file_path"`
	TransformedFilePath string                          `json:"transformed_file_path"`
	NormalizationResult AINormalizationWorkflowResponse `json:"normalization_result"`
	Error               string                          `json:"error"`
}

type AINormalizationWorkflowResponse struct {
	TransformedDataBucket string                  `json:"transformed_data_bucket"`
	TransformedDataPath   string                  `json:"transformed_data_path"`
	ColumnMapping         map[string]interface{}  `json:"column_mapping"`
	DataPreview           dbmodels.DatasetPreview `json:"data_preview"`
	ExtractedMetadata     map[string]interface{}  `json:"extracted_metadata"`
}

type FileImportDatasetActionConfig struct {
	Version             int                           `json:"version"`
	FileId              uuid.UUID                     `json:"file_id"`
	WorkflowId          uuid.UUID                     `json:"workflow_id"`
	WorkflowInitPayload FileImportWorkflowInitPayload `json:"workflow_init_payload"`
	WorkflowExitPayload FileImportWorkflowExitPayload `json:"workflow_exit_payload"`
}

type DataplatformOptions struct {
	GetDataFromLake bool
}
