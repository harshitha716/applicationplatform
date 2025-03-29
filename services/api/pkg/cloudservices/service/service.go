package service

import (
	"context"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"
	gcpservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/service"
)

type CloudService interface {
	GetSignedUrlToDownload(ctx context.Context, objectName string, optionalConfigs []models.GetDownloadsignedUrlConfigs) (*string, error)
	UploadFileToCloud(ctx context.Context, fileName string, fileData []byte) (models.SignedUrlToUpload, error)
}

type cloudService struct {
	serviceProvider CloudService
}

func NewCloudService(providerName string, serverConfig serverconfig.ConfigVariables) (CloudService, error) {
	switch providerName {
	case constants.GcpProviderName:
		return &cloudService{
			serviceProvider: gcpservice.NewGcpService(serverConfig),
		}, nil
	default:
		return nil, errors.ErrInvalidProviderName
	}
}

func (c *cloudService) GetSignedUrlToDownload(ctx context.Context, objectName string, optionalConfigs []models.GetDownloadsignedUrlConfigs) (*string, error) {
	return c.serviceProvider.GetSignedUrlToDownload(ctx, objectName, optionalConfigs)
}

func (c *cloudService) UploadFileToCloud(ctx context.Context, fileName string, fileData []byte) (models.SignedUrlToUpload, error) {
	return c.serviceProvider.UploadFileToCloud(ctx, fileName, fileData)
}
