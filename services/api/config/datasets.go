package serverconfig

type DatasetConfig struct {
	DataplatformProvider string `json:"dataplatformProvider"`
}

func getDatasetConfig(configVariables *ConfigVariables) DatasetConfig {
	return DatasetConfig{
		DataplatformProvider: configVariables.DataplatformProvider,
	}
}
