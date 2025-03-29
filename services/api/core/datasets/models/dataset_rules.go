package models

import (
	storemodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type UpdateRulePriorityParams struct {
	DatasetId      uuid.UUID                            `json:"dataset_id"`
	Column         string                               `json:"column"`
	RulePriorities storemodels.UpdateRulePriorityParams `json:"rule_priorities"`
}
