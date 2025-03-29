package service

import (
	cloudstorageconstants "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/constants"
	cloudservicemodels "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/models"
	"github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/constants"
	gcpmodels "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/providers/gcp/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetSignedUrlToDownloadConfigs(t *testing.T) {
	svc := &gcpService{
		accessID:      "test-access-id",
		defaultBucket: "test-bucket",
	}

	tests := []struct {
		name           string
		configs        []cloudservicemodels.GetDownloadsignedUrlConfigs
		wantErr        bool
		expectedParams gcpmodels.DownloadParams
	}{
		{
			name: "success with custom name",
			configs: []cloudservicemodels.GetDownloadsignedUrlConfigs{
				{
					Key:   cloudstorageconstants.SignedDownloadUrlConfigCustomName,
					Value: "custom.txt",
				},
			},
			wantErr: false,
			expectedParams: gcpmodels.DownloadParams{
				CustomDownloadName: stringPtr("custom.txt"),
			},
		},
		{
			name: "success with custom TTL",
			configs: []cloudservicemodels.GetDownloadsignedUrlConfigs{
				{
					Key:   cloudstorageconstants.SignedDownloadUrlConfigTtl,
					Value: 30 * time.Minute,
				},
			},
			wantErr: false,
			expectedParams: gcpmodels.DownloadParams{
				UrlTtlInMinutes: durationPtr(30 * time.Minute),
			},
		},
		{
			name: "error with invalid config key",
			configs: []cloudservicemodels.GetDownloadsignedUrlConfigs{
				{
					Key:   "invalid-key",
					Value: "value",
				},
			},
			wantErr: true,
		},
		{
			name: "fallback to default TTL for invalid type",
			configs: []cloudservicemodels.GetDownloadsignedUrlConfigs{
				{
					Key:   cloudstorageconstants.SignedDownloadUrlConfigTtl,
					Value: "30m", // String instead of time.Duration
				},
			},
			wantErr: false,
			expectedParams: gcpmodels.DownloadParams{
				UrlTtlInMinutes: durationPtr(time.Duration(constants.SignedUrlExpiryTimeInMinutes) * time.Minute),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := svc.getSignedUrlToDownloadConfigs(tt.configs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedParams, params)
		})
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}
