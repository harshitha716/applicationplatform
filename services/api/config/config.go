package serverconfig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/pkg/cache"
	"github.com/Zampfi/application-platform/services/api/pkg/errorreporting"
	"github.com/Zampfi/application-platform/services/api/pkg/kratosclient"
	"github.com/Zampfi/application-platform/services/api/pkg/s3"
	"github.com/Zampfi/application-platform/services/api/pkg/sparkpost"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	temporalsdk "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal"
)

type ServerConfig struct {
	AuthClient         *kratosclient.Client
	Env                *ConfigVariables
	DataPlatformConfig *DataPlatformConfig
	DatasetConfig      *DatasetConfig
	Store              store.Store
	TemporalSdk        temporalsdk.TemporalService
	SparkpostClient    sparkpost.SparkPostClient
	DefaultS3Client    s3.S3Client
	CacheClient        cache.CacheClient
}

func Createserverconfig(logger *zap.Logger) (*ServerConfig, func(), error) {
	configVars, err := getServerConfigVariables()
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing server config variables: %v", err)
	}
	if configVars == nil {
		return nil, nil, fmt.Errorf("config variables are nil")
	}
	dataPlatformConfig, err := getDataPlatformConfig(configVars)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing data platform config: %v", err)
	}
	datasetConfig := getDatasetConfig(configVars)
	pgClient, err := pgclient.NewPostgresClient(configVars.PgDatabaseUrl, configVars.PgMaxIdleConnections, configVars.PgMaxOpenConnections)
	if err != nil {
		return nil, nil, err
	}
	if pgClient == nil {
		return nil, nil, fmt.Errorf("pg client is nil")
	}

	kratosClient, err := kratosclient.NewClient(configVars.AuthBaseUrl)
	if kratosClient == nil || err != nil {
		return nil, nil, fmt.Errorf("error initializing auth client")
	}

	temporalService := temporalsdk.NewTemporalService()
	err = ConnectToTemporalBasedOnEnv(
		temporalService,
		configVars,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing temporal service: %v", err)
	}

	sparkpostClient, err := sparkpost.NewClient(sparkpost.Config{
		APIKey: configVars.SparkpostAPIKey,
		APIUrl: configVars.SparkpostAPIURL,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing sparkpost client: %v", err)
	}

	store, cleanup := store.NewStore(pgClient)

	defaultS3Client, err := s3.NewDefaultS3Client(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing default s3 client: %v", err)
	}

	// initialize error reporting
	err = errorreporting.InitializeErrorReporting(configVars.SentryDsn, configVars.Environment)
	if err != nil {
		logger.Warn("Error initializing error reporting", zap.Error(err))
	}

	redisDB, err := strconv.Atoi(configVars.RedisDB)
	if err != nil {
		logger.Error("error converting redis db to int", zap.Error(err))
	}

	cacheClient := cache.NewRedisCache(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configVars.RedisHost, configVars.RedisPort),
		Password: configVars.RedisPassword,
		DB:       redisDB,
	}))

	return &ServerConfig{
		Env:                configVars,
		AuthClient:         kratosClient,
		Store:              store,
		DataPlatformConfig: &dataPlatformConfig,
		DatasetConfig:      &datasetConfig,
		TemporalSdk:        temporalService,
		SparkpostClient:    sparkpostClient,
		DefaultS3Client:    defaultS3Client,
		CacheClient:        cacheClient,
	}, cleanup, nil
}

func GetEmptyServerConfig() *ServerConfig {
	return &ServerConfig{
		Env: &ConfigVariables{
			Port: "8080",
		},
		AuthClient: &kratosclient.Client{},
	}
}
