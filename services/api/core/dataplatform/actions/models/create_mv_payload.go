package models

type CreateMVPayload struct {
	ParentDatasetIds []string `json:"parent_dataset_ids"`
	Query            string   `json:"query"`
	MerchantId       string   `json:"merchant_id"`
	DatasetId        string   `json:"dataset_id"`
	DedupColumns     []string `json:"dedup_columns"`
	OrderByColumn    string   `json:"order_by_column"`
}
