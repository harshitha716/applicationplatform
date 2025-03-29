package service

import (
	"context"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/logger"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	provider "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/databricks"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/pinot"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/postgres"

	"go.uber.org/zap"
)

type Providers string

const (
	Databricks Providers = "databricks"
	Pinot      Providers = "pinot"
)

type ProviderRegistry struct {
	databricksServices map[string]databricks.DatabricksService
	pinotServices      map[string]pinot.PinotService
	postgresServices   map[string]postgres.PostgresService
	providerConfigs    map[string]models.ProviderConfig
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		databricksServices: make(map[string]databricks.DatabricksService),
		pinotServices:      make(map[string]pinot.PinotService),
		postgresServices:   make(map[string]postgres.PostgresService),
		providerConfigs:    make(map[string]models.ProviderConfig),
	}
}

type ProviderService interface {
	Query(ctx context.Context, providerType constants.ProviderType, dataProviderId string, table string, query string, args ...interface{}) (models.QueryResult, error)
	GetDatabricksService(ctx context.Context, dataProviderId string) (databricks.DatabricksService, error)
	GetPinotService(ctx context.Context, dataProviderId string) (pinot.PinotService, error)
	GetPostgresService(ctx context.Context, dataProviderId string) (postgres.PostgresService, error)
	GetService(ctx context.Context, providerType constants.ProviderType, dataProviderId string) (provider.ProviderService, error)
}

type providerService struct {
	providerRegistry *ProviderRegistry
}

func getProviderConfigKey(providerType constants.ProviderType, dataProviderId string) string {
	return fmt.Sprintf("%s_%s", providerType, dataProviderId)
}

func InitProviders(providerConfigs []models.ProviderConfig) (ProviderService, error) {
	providerRegistry := NewProviderRegistry()
	for _, config := range providerConfigs {
		providerRegistry.providerConfigs[getProviderConfigKey(config.Provider, config.DataProviderId)] = config
	}
	for _, config := range providerConfigs {
		switch config.Provider {
		case constants.ProviderTypeDatabricks:
			if _, exists := providerRegistry.databricksServices[config.DataProviderId]; exists {
				return nil, errors.ErrDatabricksServiceAlreadyInitialized
			}
			databricksConfig, ok := config.Config.(models.DatabricksConfig)
			if !ok {
				return nil, errors.ErrInvalidConfigurationForDatabricks
			}
			service, err := databricks.InitDatabricksService(databricksConfig)
			if err == nil {
				providerRegistry.databricksServices[config.DataProviderId] = service
			}
		case constants.ProviderTypePinot:
			if _, exists := providerRegistry.pinotServices[config.DataProviderId]; exists {
				return nil, errors.ErrPinotServiceAlreadyInitialized
			}
			pinotConfig, ok := config.Config.(models.PinotConfig)
			if !ok {
				return nil, errors.ErrInvalidConfigurationForPinot
			}
			service, err := pinot.InitPinotService(pinotConfig)
			if err == nil {
				providerRegistry.pinotServices[config.DataProviderId] = service
			}
		case constants.ProviderTypePostgres:
			if _, exists := providerRegistry.postgresServices[config.DataProviderId]; exists {
				return nil, errors.ErrPostgresServiceAlreadyInitialized
			}
			postgresConfig, ok := config.Config.(models.PostgresConfig)
			if !ok {
				return nil, errors.ErrInvalidConfigurationForPostgres
			}
			service, err := postgres.InitPostgresService(postgresConfig)
			if err == nil {
				providerRegistry.postgresServices[config.DataProviderId] = service
			}
		default:
			return nil, errors.ErrUnsupportedProviderConfiguration
		}

	}
	return &providerService{
		providerRegistry: providerRegistry,
	}, nil
}

func (r *providerService) reinitializeProvider(ctx context.Context, providerType constants.ProviderType, dataProviderId string) (provider.ProviderService, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	logger.Info("REINITIALIZING_PROVIDER", zap.String("providerType", string(providerType)), zap.String("dataProviderId", dataProviderId))
	config, exists := r.providerRegistry.providerConfigs[getProviderConfigKey(providerType, dataProviderId)]
	if !exists {
		logger.Error(errors.ProviderConfigNotInitializedErrMessage)
		return nil, errors.ErrProviderConfigNotInitialized
	}

	switch providerType {
	case constants.ProviderTypeDatabricks:
		databricksConfig, ok := config.Config.(models.DatabricksConfig)
		if !ok {
			logger.Error(errors.InvalidConfigurationForDatabricksErrMessage)
			return nil, errors.ErrInvalidConfigurationForDatabricks
		}

		service, err := databricks.InitDatabricksService(databricksConfig)
		if err != nil {
			logger.Error(errors.DatabricksServiceInitializationFailedErrMessage, zap.Error(err))
			return nil, err
		}

		r.providerRegistry.databricksServices[dataProviderId] = service
		return service, nil

	case constants.ProviderTypePinot:
		pinotConfig, ok := config.Config.(models.PinotConfig)
		if !ok {
			logger.Error(errors.InvalidConfigurationForPinotErrMessage)
			return nil, errors.ErrInvalidConfigurationForPinot
		}

		service, err := pinot.InitPinotService(pinotConfig)
		if err != nil {
			logger.Error(errors.PinotServiceInitializationFailedErrMessage, zap.Error(err))
			return nil, err
		}

		r.providerRegistry.pinotServices[dataProviderId] = service
		return service, nil

	case constants.ProviderTypePostgres:
		postgresConfig, ok := config.Config.(models.PostgresConfig)
		if !ok {
			logger.Error(errors.InvalidConfigurationForPostgresErrMessage)
			return nil, errors.ErrInvalidConfigurationForPostgres
		}

		service, err := postgres.InitPostgresService(postgresConfig)
		if err != nil {
			logger.Error(errors.PostgresServiceInitializationFailedErrMessage, zap.Error(err))
			return nil, err
		}

		r.providerRegistry.postgresServices[dataProviderId] = service
		return service, nil
	}

	return nil, errors.ErrUnsupportedProviderType
}

func (r *providerService) GetDatabricksService(ctx context.Context, dataProviderId string) (databricks.DatabricksService, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	service, exists := r.providerRegistry.databricksServices[dataProviderId]
	if !exists {
		service, err := r.reinitializeProvider(ctx, constants.ProviderTypeDatabricks, dataProviderId)
		if err != nil {
			logger.Error(errors.DatabricksServiceInitializationFailedErrMessage, zap.Error(err))
			return nil, errors.ErrDatabricksServiceInitializationFailed
		}

		return service.(databricks.DatabricksService), nil
	}
	return service, nil
}

func (r *providerService) GetPinotService(ctx context.Context, dataProviderId string) (pinot.PinotService, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	service, exists := r.providerRegistry.pinotServices[dataProviderId]
	if !exists {
		service, err := r.reinitializeProvider(ctx, constants.ProviderTypePinot, dataProviderId)
		if err != nil {
			logger.Error(errors.PinotServiceInitializationFailedErrMessage, zap.Error(err))
			return nil, errors.ErrPinotServiceInitializationFailed
		}

		return service.(pinot.PinotService), nil
	}
	return service, nil
}

func (r *providerService) GetPostgresService(ctx context.Context, dataProviderId string) (postgres.PostgresService, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	service, exists := r.providerRegistry.postgresServices[dataProviderId]
	if !exists {
		service, err := r.reinitializeProvider(ctx, constants.ProviderTypePostgres, dataProviderId)
		if err != nil {
			logger.Error(errors.PostgresServiceInitializationFailedErrMessage, zap.Error(err))
			return nil, errors.ErrPostgresServiceInitializationFailed
		}

		return service.(postgres.PostgresService), nil
	}
	return service, nil
}

func (r *providerService) GetService(ctx context.Context, providerType constants.ProviderType, dataProviderId string) (provider.ProviderService, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	var service provider.ProviderService
	var err error

	switch providerType {

	case constants.ProviderTypeDatabricks:
		service, err = r.GetDatabricksService(ctx, dataProviderId)

	case constants.ProviderTypePinot:
		service, err = r.GetPinotService(ctx, dataProviderId)

	case constants.ProviderTypePostgres:
		service, err = r.GetPostgresService(ctx, dataProviderId)

	default:
		return nil, errors.ErrUnsupportedProviderType
	}

	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return nil, errors.ErrProviderServiceNotFound
	}

	return service, nil
}

func (r *providerService) Query(ctx context.Context, providerType constants.ProviderType, dataProviderId string, table string, query string, args ...interface{}) (models.QueryResult, error) {
	logger := logger.GetLoggerFromCtx(ctx)
	service, err := r.GetService(ctx, providerType, dataProviderId)
	if err != nil {
		logger.Error(errors.ProviderServiceNotFoundErrMessage, zap.Error(err))
		return models.QueryResult{}, err
	}

	return service.Query(ctx, table, query, args...)
}
