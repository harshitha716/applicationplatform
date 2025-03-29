package serverconfig

import (
	"encoding/json"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
)

type DatabricksSetupConfig struct {
	DefaultDataProviderId         string                             `json:"defaultDataProviderId"`
	DataProviderConfigs           map[string]models.DatabricksConfig `json:"dataProviderConfigs"`
	MerchantDataProviderIdMapping map[string]string                  `json:"merchantDataProviderIdMapping"`
	ZampDatabricksCatalog         string                             `json:"zampDatabricksCatalog"`
	ZampDatabricksPlatformSchema  string                             `json:"zampDatabricksPlatformSchema"`
}

type PinotSetupConfig struct {
	DefaultDataProviderId         string                        `json:"defaultDataProviderId"`
	DataProviderConfigs           map[string]models.PinotConfig `json:"dataProviderConfigs"`
	MerchantDataProviderIdMapping map[string]string             `json:"merchantDataProviderIdMapping"`
}

type CreateMVJobTemplateConfig struct {
	CreateMVNotebookPath   string `json:"createMVNotebookPath"`
	SideEffectNotebookPath string `json:"sideEffectNotebookPath"`
}

type RegisterDatasetJobTemplateConfig struct {
	RegisterDatasetNotebookPath string `json:"registerDatasetNotebookPath"`
}

type RegisterJobJobTemplateConfig struct {
	RegisterJobNotebookPath string `json:"registerJobNotebookPath"`
}

type UpdateDatasetJobTemplateConfig struct {
	UpdateDatasetNotebookPath string `json:"updateDatasetNotebookPath"`
}

type UpsertTemplateJobTemplateConfig struct {
	UpsertTemplateNotebookPath string `json:"upsertTemplateNotebookPath"`
}

type CopyDatasetJobTemplateConfig struct {
	CopyDatasetNotebookPath string `json:"copyDatasetNotebookPath"`
}

type ActionsConfig struct {
	CreateMVJobTemplateConfig        CreateMVJobTemplateConfig        `json:"createMVJobTemplateConfig"`
	RegisterDatasetJobTemplateConfig RegisterDatasetJobTemplateConfig `json:"registerDatasetJobTemplateConfig"`
	RegisterJobJobTemplateConfig     RegisterJobJobTemplateConfig     `json:"registerJobJobTemplateConfig"`
	UpdateDatasetJobTemplateConfig   UpdateDatasetJobTemplateConfig   `json:"updateDatasetJobTemplateConfig"`
	UpsertTemplateJobTemplateConfig  UpsertTemplateJobTemplateConfig  `json:"upsertTemplateJobTemplateConfig"`
	CopyDatasetJobTemplateConfig     CopyDatasetJobTemplateConfig     `json:"copyDatasetJobTemplateConfig"`
	WebhookConfig                    WebhookConfig                    `json:"webhookConfig"`
	RunAsUserName                    string                           `json:"runAsUserName"`
	PinotIngestionNotebookPath       string                           `json:"pinotIngestionNotebookPath"`
	PinotClusterId                   string                           `json:"pinotClusterId"`
	DataPlatformModulesSrc           string                           `json:"dataPlatformModulesSrc"`
}

type WebhookConfig struct {
	WebhookId string `json:"webhookId"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
}

type DataPlatformConfig struct {
	DatabricksConfig DatabricksSetupConfig `json:"databricks"`
	PinotConfig      PinotSetupConfig      `json:"pinot"`
	ActionsConfig    ActionsConfig         `json:"actionsConfig"`
	RosettaBaseUrl   string                `json:"rosettaBaseUrl"`
}

// TODO Add a validation function for the data platform config

func getDataPlatformConfig(configVariables *ConfigVariables) (DataPlatformConfig, error) {
	rawDataPlatformConfig := configVariables.DataPlatformConfig
	var dataPlatformConfig DataPlatformConfig
	err := json.Unmarshal([]byte(rawDataPlatformConfig), &dataPlatformConfig)
	if err != nil {
		return DataPlatformConfig{}, err
	}
	return dataPlatformConfig, nil
}
