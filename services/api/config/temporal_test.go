package serverconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemporalConfig(t *testing.T) {
	config := GetTemporalConfig()
	assert.NotNil(t, config)
	//assert.Equal(t, config.Host, "localhost:7233")
	assert.Equal(t, config.Namespace, "default")
}

func TestGetTemporalConfigWithEnv(t *testing.T) {
	os.Setenv("TEMPORAL_HOST", "test-host:7233")
	os.Setenv("TEMPORAL_NAMESPACE", "test-namespace")
	config := GetTemporalConfig()
	assert.NotNil(t, config)
	assert.Equal(t, config.Host, "test-host:7233")
	assert.Equal(t, config.Namespace, "test-namespace")

}
