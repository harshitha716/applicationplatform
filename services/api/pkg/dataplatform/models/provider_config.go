package models

import "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"

type ProviderConfig struct {
	DataProviderId string                 `json:"dataProviderId"`
	Provider       constants.ProviderType `json:"provider"`
	Config         interface{}            `json:"config"`
}
