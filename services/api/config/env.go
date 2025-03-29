package serverconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ConfigVariables struct {
	// The port the server will listen on
	Port string
	// The authentication base url
	AuthBaseUrl string
	// Postgres database URL
	PgDatabaseUrl string
	// Maximum number of idle connections in the pool
	PgMaxIdleConnections int
	// Maximum number of open connections in the pool
	PgMaxOpenConnections int
	// Allowed CORS origins
	AllowedCORSOrigins string
	// Pantheon URL
	PantheonURL string
	// Runtime environment
	Environment string
	// Admin secrets
	ControlPlaneAdminSecrets []string
	// Sentry DSN
	SentryDsn string
	// Dataplatform config
	DataPlatformConfig string
	// HTML templates path
	AdminHTMLTemplatesPath string
	// Temporal config
	TemporalConfig TemporalConfig
	// Google Access id
	GcpAccessID string
	// GcpBucketName
	GcpBucketName string
	// AWS Access Key ID
	AWSAccessKeyID string
	// AWS Secret Access Key
	AWSSecretAccessKey string
	// AWS Default Bucket Name
	AWSDefaultBucketName string
	// Dataplatform provider
	DataplatformProvider string
	// Sparkpost API Key
	SparkpostAPIKey string
	// Sparkpost API URL
	SparkpostAPIURL string
	// Zamp email updates from
	ZampEmailUpdatesFrom string
	// Email templates path
	EmailTemplatesPath string
	// Redis config
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
	RedisUseTLS   string
}

func (e *ConfigVariables) validate() error {

	if e.Port == "" {
		return fmt.Errorf("PORT is required")
	} else {
		// check if port is a valid number
		_, err := strconv.Atoi(e.Port)
		if err != nil {
			return fmt.Errorf("PORT must be a valid number")
		}
	}

	if e.AuthBaseUrl == "" {
		return fmt.Errorf("AUTH_HOST is required")
	}

	if e.PgDatabaseUrl == "" {
		return fmt.Errorf("PG_DATABASE_URL_APP is required")
	}

	if e.ControlPlaneAdminSecrets == nil || len(e.ControlPlaneAdminSecrets) == 0 {
		return fmt.Errorf("CONTROL_PLANE_ADMIN_SECRETS is required")
	}

	if e.AdminHTMLTemplatesPath == "" {
		return fmt.Errorf("ADMIN_HTML_TEMPLATES_PATH is required")
	}

	if e.TemporalConfig.Host == "" {
		return fmt.Errorf("TEMPORAL_HOST is required")
	}

	if e.TemporalConfig.Namespace == "" {
		return fmt.Errorf("TEMPORAL_NAMESPACE is required")
	}

	if e.GcpAccessID == "" {
		return fmt.Errorf("GCP_ACCESS_ID is required")
	}

	if e.GcpBucketName == "" {
		return fmt.Errorf("GCP_BUCKET_NAME is required")
	}

	if e.AWSAccessKeyID == "" {
		return fmt.Errorf("AWS_ACCESS_KEY_ID is required")
	}

	if e.AWSSecretAccessKey == "" {
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY is required")
	}

	if e.AWSDefaultBucketName == "" {
		return fmt.Errorf("AWS_DEFAULT_BUCKET_NAME is required")
	}
	if e.DataplatformProvider == "" {
		return fmt.Errorf("DATAPLATFORM_PROVIDER is required")
	}

	return nil
}

func getIntEnvVariableWithDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvVariableWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseCommaSeparatedList(commaSeparatedSecrets string) []string {
	if commaSeparatedSecrets == "" {
		return nil
	}

	adminSecrets := strings.Split(commaSeparatedSecrets, ",")
	for i, secret := range adminSecrets {
		adminSecrets[i] = strings.TrimSpace(secret)
	}

	return adminSecrets
}

// NewEnv creates a new Env struct with the default values
func getServerConfigVariables() (*ConfigVariables, error) {

	env := &ConfigVariables{
		Port:                     getEnvVariableWithDefault("PORT", "8080"),
		AuthBaseUrl:              getEnvVariableWithDefault("AUTH_HOST", ""),
		PgDatabaseUrl:            getEnvVariableWithDefault("PG_DATABASE_URL_APP", ""),
		PgMaxIdleConnections:     getIntEnvVariableWithDefault("PG_MAX_IDLE_CONNECTIONS", 2),
		PgMaxOpenConnections:     getIntEnvVariableWithDefault("PG_MAX_OPEN_CONNECTIONS", 2),
		AllowedCORSOrigins:       getEnvVariableWithDefault("CORS_ALLOWED_ORIGINS", ""),
		PantheonURL:              getEnvVariableWithDefault("PANTHEON_URL", ""),
		Environment:              getEnvVariableWithDefault("ENVIRONMENT", "production"),
		ControlPlaneAdminSecrets: parseCommaSeparatedList(getEnvVariableWithDefault("CONTROL_PLANE_ADMIN_SECRETS", "")),
		SentryDsn:                getEnvVariableWithDefault("SENTRY_DSN", ""),
		DataPlatformConfig:       getEnvVariableWithDefault("DATAPLATFORM_CONFIG", ""),
		AdminHTMLTemplatesPath:   getEnvVariableWithDefault("ADMIN_HTML_TEMPLATES_PATH", "/app/templates"),
		TemporalConfig:           GetTemporalConfig(),
		GcpAccessID:              getEnvVariableWithDefault("GCP_ACCESS_ID", ""),
		GcpBucketName:            getEnvVariableWithDefault("GCP_BUCKET_NAME", ""),
		AWSAccessKeyID:           getEnvVariableWithDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:       getEnvVariableWithDefault("AWS_SECRET_ACCESS_KEY", ""),
		AWSDefaultBucketName:     getEnvVariableWithDefault("AWS_DEFAULT_BUCKET_NAME", ""),
		DataplatformProvider:     getEnvVariableWithDefault("DATAPLATFORM_PROVIDER", "databricks"),
		SparkpostAPIKey:          getEnvVariableWithDefault("SPARKPOST_API_KEY", ""),
		SparkpostAPIURL:          getEnvVariableWithDefault("SPARKPOST_API_URL", ""),
		ZampEmailUpdatesFrom:     getEnvVariableWithDefault("ZAMP_EMAIL_UPDATES_FROM", "noreply@zamp.ai"),
		EmailTemplatesPath:       getEnvVariableWithDefault("EMAIL_TEMPLATES_PATH", "/app/email_templates"),
		RedisHost:                getEnvVariableWithDefault("REDIS_HOST", ""),
		RedisPort:                getEnvVariableWithDefault("REDIS_PORT", ""),
		RedisPassword:            getEnvVariableWithDefault("REDIS_PASSWORD", ""),
		RedisDB:                  getEnvVariableWithDefault("REDIS_DB", ""),
	}

	err := env.validate()
	if err != nil {
		return nil, err
	}
	return env, nil
}
