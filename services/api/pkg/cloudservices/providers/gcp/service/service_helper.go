package service

import (
	"context"
	"time"

	cloudservicemodels "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"

	"cloud.google.com/go/storage"
	cloudstorageconstants "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/models"
)

func (gs *gcpService) getClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (gs *gcpService) getSignedUrlToDownloadConfigs(optionalConfigs []cloudservicemodels.GetDownloadsignedUrlConfigs) (models.DownloadParams, error) {
	configs := models.DownloadParams{}

	for _, config := range optionalConfigs {
		switch config.Key {
		case cloudstorageconstants.SignedDownloadUrlConfigCustomName:
			if customFileName, ok := config.Value.(string); ok {
				configs.CustomDownloadName = &customFileName
			}
		case cloudstorageconstants.SignedDownloadUrlConfigTtl:
			if urlTtl, ok := config.Value.(time.Duration); ok {
				configs.UrlTtlInMinutes = &urlTtl
			} else {
				urlTtl := time.Duration(constants.SignedUrlExpiryTimeInMinutes) * time.Minute
				configs.UrlTtlInMinutes = &urlTtl
			}

		default:
			return models.DownloadParams{}, errors.ErrInvalidDownloadParams
		}
	}

	return configs, nil
}

func (gs *gcpService) getSignedUrlToUploadForBytes(ctx context.Context, objectName string) (cloudservicemodels.SignedUrlToUpload, error) {
	client, err := gs.getClient(ctx)
	if err != nil {
		return cloudservicemodels.SignedUrlToUpload{}, err
	}

	options := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         constants.PutRequest,
		GoogleAccessID: gs.accessID,
		Expires:        time.Now().Add(time.Duration(constants.SignedUrlExpiryTimeInMinutes) * time.Minute),
	}

	url, err := client.Bucket(gs.defaultBucket).SignedURL(objectName, options)
	if err != nil {
		return cloudservicemodels.SignedUrlToUpload{}, err
	}

	return cloudservicemodels.SignedUrlToUpload{
		Url:        url,
		Identifier: objectName,
	}, nil
}
