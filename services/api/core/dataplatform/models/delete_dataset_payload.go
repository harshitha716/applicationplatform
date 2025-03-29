package models

type DeleteDatasetPayload struct {
	MerchantID string `json:"merchant_id"`
	DatasetID  string `json:"dataset_id"`
	ActorId    string `json:"actor_id"`
}
