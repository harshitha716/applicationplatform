package serverconfig

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/Zampfi/application-platform/services/api/helper/constants"
	temporalSdk "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal"
	temporalConfigModels "github.com/Zampfi/workflow-sdk-go/workflowmanagers/temporal/models"
	"go.temporal.io/sdk/client"
)

type TemporalConfig struct {
	Host        string
	Namespace   string
	Certificate tls.Certificate
}

func GetTemporalConfig() TemporalConfig {
	TemporalCertificateFilePath := os.Getenv("TEMPORAL_CERTIFICATE_FILE")
	TemporalKeyFilePath := os.Getenv("TEMPORAL_KEY_FILE")

	temporalCert, _ := tls.LoadX509KeyPair(TemporalCertificateFilePath, TemporalKeyFilePath)

	host := getEnvVariableWithDefault("TEMPORAL_HOST", "localhost:7233")
	namespace := getEnvVariableWithDefault("TEMPORAL_NAMESPACE", "default")

	return TemporalConfig{
		Host:        host,
		Namespace:   namespace,
		Certificate: temporalCert,
	}
}

func ConnectToTemporalBasedOnEnv(temporalService temporalSdk.TemporalService, env *ConfigVariables) error {
	ctx := context.Background()
	if env.Environment == constants.ENVLOCAL {
		return temporalService.Connect(
			ctx,
			temporalConfigModels.ConnectClientParams{
				Options: temporalConfigModels.ConnectClientOptions{
					HostPort:  env.TemporalConfig.Host,
					Namespace: env.TemporalConfig.Namespace,
				},
			},
		)
	}

	return temporalService.Connect(
		ctx,
		temporalConfigModels.ConnectClientParams{
			Options: temporalConfigModels.ConnectClientOptions{
				HostPort:    env.TemporalConfig.Host,
				Namespace:   env.TemporalConfig.Namespace,
				Credentials: client.NewMTLSCredentials(env.TemporalConfig.Certificate),
			},
		},
	)

}
