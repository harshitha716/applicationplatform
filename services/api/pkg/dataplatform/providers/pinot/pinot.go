package pinot

import (
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	provider "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers"

	pinot "github.com/startreedata/pinot-client-go/pinot"
)

type PinotService interface {
	provider.ProviderService
}

type pinotService struct {
	pinotClient *pinot.Connection
}

func InitPinotService(configs models.PinotConfig) (PinotService, error) {
	pinotClient, err := InitPinotSqlService(configs)
	if err != nil {
		return nil, err
	}
	return &pinotService{
		pinotClient: pinotClient,
	}, nil
}
