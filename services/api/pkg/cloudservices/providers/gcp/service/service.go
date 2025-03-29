package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	apihelper "github.com/Zampfi/application-platform/services/api/helper"
	cloudservicemodels "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/helpers"
)

type GcpService interface {
	GetSignedUrlToDownload(ctx context.Context, objectName string, optionalConfigs []cloudservicemodels.GetDownloadsignedUrlConfigs) (*string, error)
	UploadFileToCloud(ctx context.Context, fileName string, fileData []byte) (cloudservicemodels.SignedUrlToUpload, error)
}

type gcpService struct {
	accessID      string
	defaultBucket string
}

func NewGcpService(gcpCloudConfig serverconfig.ConfigVariables) *gcpService {
	return &gcpService{
		accessID:      gcpCloudConfig.GcpAccessID,
		defaultBucket: gcpCloudConfig.GcpBucketName,
	}
}

func (gs *gcpService) GetSignedUrlToDownload(ctx context.Context, objectName string, optionalConfigs []cloudservicemodels.GetDownloadsignedUrlConfigs) (*string, error) {
	client, err := gs.getClient(ctx)
	if err != nil {
		return nil, err
	}

	options := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         constants.GetRequest,
		GoogleAccessID: gs.accessID,
		Expires:        time.Now().Add(time.Duration(constants.SignedUrlExpiryTimeInMinutes) * time.Minute),
	}

	if len(optionalConfigs) > 0 {
		configs, err := gs.getSignedUrlToDownloadConfigs(optionalConfigs)
		if err != nil {
			return nil, err
		}

		if configs.CustomDownloadName != nil {
			queryParams := url.Values{}
			encodedFileName := url.QueryEscape(*configs.CustomDownloadName)
			queryParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", encodedFileName))
			options.QueryParameters = queryParams
		}

		if configs.UrlTtlInMinutes != nil {
			options.Expires = time.Now().Add(*configs.UrlTtlInMinutes)
		}
	}

	bucket, filePath, err := helpers.ExtractBucketAndPath(objectName)
	if err == errors.ErrDefaultGcsBucket {
		bucket = gs.defaultBucket
		filePath = objectName
	} else if err != nil {
		return nil, err
	}

	url, err := client.Bucket(bucket).SignedURL(filePath, options)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (gs *gcpService) UploadFileToCloud(ctx context.Context, fileName string, fileData []byte) (cloudservicemodels.SignedUrlToUpload, error) {
	uploadUrl, err := gs.getSignedUrlToUploadForBytes(ctx, fileName)
	if err != nil {
		return cloudservicemodels.SignedUrlToUpload{}, err

	}

	headers := map[string]string{}
	httpClient := http.Client{}
	_, err = apihelper.HttpPut(&httpClient, uploadUrl.Url, fileData, headers)
	if err != nil {
		return cloudservicemodels.SignedUrlToUpload{}, err
	}
	return uploadUrl, nil
}
