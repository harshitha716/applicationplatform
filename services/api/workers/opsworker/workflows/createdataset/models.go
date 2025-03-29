package createdataset

import (
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	"github.com/google/uuid"
)

type CreateDatasetWorkflowInitPayload struct {
	RegisterDatasetPayload models.DatasetCreationInfo `json:"register_dataset_payload"`
	OrganizationId         uuid.UUID                  `json:"organization_id"`
	UserId                 uuid.UUID                  `json:"user_id"`
}

func (p *CreateDatasetWorkflowInitPayload) GetAccessControlParams() (*uuid.UUID, []uuid.UUID) {
	return &p.OrganizationId, []uuid.UUID{p.UserId}
}

type CreateDatasetWorkflowExitPayload struct {
	DatasetId uuid.UUID `json:"dataset_id"`
}

type RegisterDatasetResponse struct {
	ActionId  string    `json:"action_id"`
	DatasetId uuid.UUID `json:"dataset_id"`
}
