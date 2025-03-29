package models

import (
	"time"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
)

type CreateActionPayload struct {
	// WorkflowID            string               `json:"workflowId"` // TODO: INTRODUCE THIS ONCE WE HAVE WORKFLOWS
	MerchantID            string               `json:"merchantId"`
	ActionType            constants.ActionType `json:"actionType"`
	ActionMetadataPayload interface{}          `json:"actionMetadataPayload"`
	ActorId               string               `json:"actorId"`
}

type CreateActionResponse struct {
	ActionID string `json:"actionId"`
}

type CreateMVActionPayload struct {
	Query            string            `json:"query"`
	QueryParams      map[string]string `json:"queryParams"`
	ParentDatasetIds []string          `json:"parentDatasetIds"`
	MVDatasetId      string            `json:"mvDatasetId"`
	DedupColumns     []string          `json:"dedupColumns"`
	OrderByColumn    string            `json:"orderByColumn"`
}

type UpdateDatasetDataActionPayload struct {
	DatasetId    string         `json:"dataset_id"`
	SqlCondition string         `json:"sql_condition"`
	UpdateValues map[string]any `json:"update_values"`
}

type Action struct {
	ID             string                 `json:"id"`
	WorkspaceId    string                 `json:"workspace_id"`
	RunId          int64                  `json:"run_id"`
	ActionType     constants.ActionType   `json:"action_type"`
	ActionStatus   constants.ActionStatus `json:"status"`
	ActionMetadata interface{}            `json:"action_metadata"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ActorId        string                 `json:"actor_id"`
}

type SubmitActionResponse struct {
	RunId int64 `json:"runId"`
}

type SourceColumnUpdateValue struct {
	SourceType      string `json:"source_type"`
	SourceId        string `json:"source_id"`
	SourceUpdatedAt string `json:"source_updated_at"`
}
