package models

import (
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
)

type UpdateDatasetActionPayload struct {
	DatasetId     string                      `json:"dataset_id"`
	DatasetConfig datasetmodels.DatasetConfig `json:"dataset_config"`
}

type UpsertRuleEventMetadata struct {
	DeltaRuleId string                        `json:"delta_rule_id"`
	Column      string                        `json:"column"`
	Type        constants.UpsertRuleOperation `json:"type"`
}

type UpdateDatasetEvent struct {
	EventType     constants.UpdateDatasetEventType `json:"event_type"`
	EventData     UpdateDatasetActionPayload       `json:"event_data"`
	EventMetadata UpsertRuleEventMetadata          `json:"event_metadata"`
}
