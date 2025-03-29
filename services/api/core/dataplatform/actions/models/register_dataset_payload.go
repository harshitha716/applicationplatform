package models

import (
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
)

type RegisterDatasetActionPayload struct {
	MerchantId       string                         `json:"merchant_id"`
	DatasetId        string                         `json:"dataset_id"`
	DatasetConfig    datasetmodels.DatasetConfig    `json:"dataset_config"`
	DatabricksConfig datasetmodels.DatabricksConfig `json:"databricks_config"`
	Provider         constants.Provider             `json:"provider"`
}
