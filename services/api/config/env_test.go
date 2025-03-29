package serverconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEnv(t *testing.T) {
	env := &ConfigVariables{
		Port: "notint",
	}

	err := env.validate()

	assert.NotNil(t, err)
	assert.Equal(t, "PORT must be a valid number", err.Error())

	env.Port = "8080"
	err = nil
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "AUTH_HOST is required", err.Error())

	env.AuthBaseUrl = "http://localhost:8080"

	err = nil
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "PG_DATABASE_URL_APP is required", err.Error())

	env.PgDatabaseUrl = "postgres://localhost:5432"
	err = nil

	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "CONTROL_PLANE_ADMIN_SECRETS is required", err.Error())

	env.ControlPlaneAdminSecrets = []string{"secret1"}

	err = nil

	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "ADMIN_HTML_TEMPLATES_PATH is required", err.Error())

	env.AdminHTMLTemplatesPath = "/app/templates"

	err = nil

	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "TEMPORAL_HOST is required", err.Error())

	env.TemporalConfig.Host = "localhost"

	err = nil
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "TEMPORAL_NAMESPACE is required", err.Error())

	env.TemporalConfig.Namespace = "default"
	err = nil
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "GCP_ACCESS_ID is required", err.Error())

	env.GcpAccessID = "GCP_ACCESS_ID"
	err = nil
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "GCP_BUCKET_NAME is required", err.Error())

	env.GcpBucketName = "GCP_BUCKET_NAME"
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "AWS_ACCESS_KEY_ID is required", err.Error())

	env.AWSAccessKeyID = "AWS_ACCESS_KEY_ID"
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "AWS_SECRET_ACCESS_KEY is required", err.Error())

	env.AWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "AWS_DEFAULT_BUCKET_NAME is required", err.Error())

	env.AWSDefaultBucketName = "AWS_DEFAULT_BUCKET_NAME"
	err = env.validate()
	assert.NotNil(t, err)
	assert.Equal(t, "DATAPLATFORM_PROVIDER is required", err.Error())

	env.DataplatformProvider = "pinot"
	err = env.validate()
	assert.Nil(t, err)

	err = nil

}

func TestParseAdminSecrets(t *testing.T) {
	adminSecrets := parseCommaSeparatedList("")
	assert.Empty(t, adminSecrets)

	adminSecrets = parseCommaSeparatedList("secret1")

	assert.Len(t, adminSecrets, 1)
	assert.Equal(t, "secret1", adminSecrets[0])

	adminSecrets = parseCommaSeparatedList("secret1,secret2")
	assert.Len(t, adminSecrets, 2)
	assert.Equal(t, "secret1", adminSecrets[0])
	assert.Equal(t, "secret2", adminSecrets[1])

	adminSecrets = parseCommaSeparatedList("secret1, secret2 ")
	assert.Len(t, adminSecrets, 2)
	assert.Equal(t, "secret1", adminSecrets[0])
	assert.Equal(t, "secret2", adminSecrets[1])

}

func TestGetIntEnvVariableWithDefault(t *testing.T) {
	value := getIntEnvVariableWithDefault("TEST_INT_ENV_VAR", 1)
	assert.Equal(t, 1, value)

	t.Setenv("TEST_INT_ENV_VAR", "1")

	value = getIntEnvVariableWithDefault("TEST_INT_ENV_VAR", 2)
	assert.Equal(t, 1, value)

	t.Setenv("TEST_INT_ENV_VAR", "abcd")

	value = getIntEnvVariableWithDefault("TEST_INT_ENV_VAR", 2)
	assert.Equal(t, 2, value)
}
