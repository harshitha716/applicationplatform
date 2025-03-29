package service

import (
	"context"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	mock_service "github.com/Zampfi/application-platform/services/api/mocks/pkg/cloudservices/providers/gcp/service"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCloudService(t *testing.T) {
	config := serverconfig.ConfigVariables{
		GcpAccessID:   "test-access-id",
		GcpBucketName: "test-bucket",
	}

	t.Run("invalid provider name", func(t *testing.T) {
		service, err := NewCloudService("invalid-provider", config)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidProviderName, err)
		assert.Nil(t, service)
	})

	t.Run("GCS provider success", func(t *testing.T) {
		service, err := NewCloudService(constants.GcpProviderName, config)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func setupTest(t *testing.T) (*cloudService, context.Context) {
	mockGcpService := mock_service.NewMockGcpService(t)
	service := &cloudService{
		serviceProvider: mockGcpService,
	}
	return service, context.Background()
}

func TestCloudService_GetSignedUrlToDownload_Success(t *testing.T) {
	service, ctx := setupTest(t)
	mockGcpService := service.serviceProvider.(*mock_service.MockGcpService)

	expectedURL := "https://storage.googleapis.com/test-bucket/test.txt"
	mockGcpService.EXPECT().
		GetSignedUrlToDownload(ctx, "test.txt", []models.GetDownloadsignedUrlConfigs{}).
		Return(&expectedURL, nil)

	url, err := service.GetSignedUrlToDownload(ctx, "test.txt", []models.GetDownloadsignedUrlConfigs{})

	assert.NoError(t, err)
	assert.Equal(t, &expectedURL, url)
}

func TestCloudService_GetSignedUrlToDownload_Failure(t *testing.T) {
	service, ctx := setupTest(t)
	mockGcpService := service.serviceProvider.(*mock_service.MockGcpService)

	mockGcpService.EXPECT().
		GetSignedUrlToDownload(ctx, "test.txt", []models.GetDownloadsignedUrlConfigs{}).
		Return(nil, assert.AnError)

	url, err := service.GetSignedUrlToDownload(ctx, "test.txt", []models.GetDownloadsignedUrlConfigs{})

	assert.Error(t, err)
	assert.Nil(t, url)
}

func TestCloudService_UploadFileToCloud_Success(t *testing.T) {
	service, ctx := setupTest(t)
	mockGcpService := service.serviceProvider.(*mock_service.MockGcpService)

	expectedResponse := models.SignedUrlToUpload{
		Url:        "https://storage.googleapis.com/test-bucket/test.txt",
		Identifier: "test.txt",
	}
	fileData := []byte("test content")

	mockGcpService.EXPECT().
		UploadFileToCloud(ctx, "test.txt", fileData).
		Return(expectedResponse, nil)

	response, err := service.UploadFileToCloud(ctx, "test.txt", fileData)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestCloudService_UploadFileToCloud_Failure(t *testing.T) {
	service, ctx := setupTest(t)
	mockGcpService := service.serviceProvider.(*mock_service.MockGcpService)

	fileData := []byte("test content")
	mockGcpService.EXPECT().
		UploadFileToCloud(ctx, "test.txt", fileData).
		Return(models.SignedUrlToUpload{}, assert.AnError)

	response, err := service.UploadFileToCloud(ctx, "test.txt", fileData)

	assert.Error(t, err)
	assert.Empty(t, response)
}
