package models

type CopyDatasetActionPayload struct {
	OriginalDatasetId string `json:"original_dataset_id"`
	NewDatasetId      string `json:"new_dataset_id"`
	MerchantId        string `json:"merchant_id"`
}
