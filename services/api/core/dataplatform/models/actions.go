package models

import (
	actionmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
)

type CreateMVPayload struct {
	MerchantID            string                             `json:"merchantId"`
	ActorId               string                             `json:"actorId"`
	ActionMetadataPayload actionmodels.CreateMVActionPayload `json:"actionMetadataPayload"`
}

type UpdateDatasetDataPayload struct {
	MerchantID            string                                      `json:"merchantId"`
	ActorId               string                                      `json:"actorId"`
	ActionMetadataPayload actionmodels.UpdateDatasetDataActionPayload `json:"actionMetadataPayload"`
}

type RegisterDatasetPayload struct {
	MerchantID            string                                    `json:"merchantId"`
	ActorId               string                                    `json:"actorId"`
	ActionMetadataPayload actionmodels.RegisterDatasetActionPayload `json:"actionMetadataPayload"`
}

type RegisterJobPayload struct {
	MerchantID            string                                `json:"merchantId"`
	ActorId               string                                `json:"actorId"`
	ActionMetadataPayload actionmodels.RegisterJobActionPayload `json:"actionMetadataPayload"`
}

type UpsertTemplatePayload struct {
	MerchantID            string                                   `json:"merchantId"`
	ActorId               string                                   `json:"actorId"`
	ActionMetadataPayload actionmodels.UpsertTemplateActionPayload `json:"actionMetadataPayload"`
}

type UpdateDatasetPayload struct {
	MerchantID            string                          `json:"merchantId"`
	ActorId               string                          `json:"actorId"`
	ActionMetadataPayload actionmodels.UpdateDatasetEvent `json:"actionMetadataPayload"`
}

type CopyDatasetPayload struct {
	MerchantID            string                                `json:"merchantId"`
	ActorId               string                                `json:"actorId"`
	ActionMetadataPayload actionmodels.CopyDatasetActionPayload `json:"actionMetadataPayload"`
}
