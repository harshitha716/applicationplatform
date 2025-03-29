package serverconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataPlatformConfig_Error(t *testing.T) {
	configVariables := &ConfigVariables{
		DataPlatformConfig: "invalid_json",
	}
	dataPlatformConfig, err := getDataPlatformConfig(configVariables)
	assert.NotNil(t, err)
	assert.Equal(t, DataPlatformConfig{}, dataPlatformConfig)
}

func TestGetDataPlatformConfig_Success(t *testing.T) {
	configVariables := &ConfigVariables{
		DataPlatformConfig: "{}",
	}
	dataPlatformConfig, err := getDataPlatformConfig(configVariables)
	assert.Nil(t, err)
	assert.Equal(t, DataPlatformConfig{}, dataPlatformConfig)
}
