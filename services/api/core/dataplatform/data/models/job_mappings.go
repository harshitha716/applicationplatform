package models

import "time"

type JobDatasetMapping struct {
	Id               string    `json:"id"`
	SourceType       string    `json:"source_type"`
	SourceValue      string    `json:"source_value"`
	DestinationType  string    `json:"destination_type"`
	DestinationValue string    `json:"destination_value"`
	JobId            int       `json:"job_id"`
	JobParams        string    `json:"job_params"`
	MerchantId       string    `json:"merchant_id"`
	IsDeleted        bool      `json:"is_deleted"`
	DeletedAt        time.Time `json:"deleted_at"`
	DeletedBy        string    `json:"deleted_by"`
}

type DAGNode struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type DatasetParents struct {
	Parents []DAGNode `json:"parents"`
}

type JobIdModel struct {
	JobId int64 `json:"job_id"`
}
