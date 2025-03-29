package fileimports

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mock_s3 "github.com/Zampfi/application-platform/services/api/mocks/pkg/s3"
	"github.com/Zampfi/application-platform/services/api/pkg/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitiateFileImport(t *testing.T) {

	userId := uuid.New()
	orgId := uuid.New()
	uploadURL := "https://test-bucket.s3.amazonaws.com/test-file"
	fileUpload := &models.FileUpload{
		ID:               uuid.New(),
		OrganizationID:   orgId,
		UploadedByUserID: userId,
		Name:             "file-imports.csv",
		PresignedURL:     uploadURL,
		Expiry:           time.Now().Add(1 * time.Hour),
		FileType:         models.FileTypeCSV,
	}

	tests := []struct {
		name      string
		setupMock func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore)
		orgId     uuid.UUID
		fileName  string
		fileType  models.FileType
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockS3Client.EXPECT().GenerateUploadURL(mock.Anything, mock.Anything, mock.Anything, mock.Anything, CONTENT_TYPE_OCTET_STREAM).Return(uploadURL, nil)
				mockStore.EXPECT().CreateFileUpload(mock.Anything, mock.Anything).Return(fileUpload, nil)
			},
			orgId:    orgId,
			fileName: fileUpload.Name,
			fileType: models.FileTypeCSV,
			wantErr:  false,
		},
		{
			name: "error - no user in context",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				// No mocks needed
			},
			orgId:    orgId,
			fileName: "test.csv",
			fileType: models.FileTypeCSV,
			wantErr:  true,
		},
		{
			name: "error - user not in org",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				// No mocks needed
			},
			orgId:    uuid.New(), // Different org ID
			fileName: "test.csv",
			fileType: models.FileTypeCSV,
			wantErr:  true,
		},
		{
			name: "error - s3 client error",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockS3Client.EXPECT().GenerateUploadURL(mock.Anything, mock.Anything, mock.Anything, mock.Anything, CONTENT_TYPE_OCTET_STREAM).Return("", fmt.Errorf("s3 error"))
			},
			orgId:    orgId,
			fileName: "test.csv",
			fileType: models.FileTypeCSV,
			wantErr:  true,
		},
		{
			name: "error - store error",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockS3Client.EXPECT().GenerateUploadURL(mock.Anything, mock.Anything, mock.Anything, mock.Anything, CONTENT_TYPE_OCTET_STREAM).Return(uploadURL, nil)
				mockStore.EXPECT().CreateFileUpload(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("store error"))
			},
			orgId:    orgId,
			fileName: "test.csv",
			fileType: models.FileTypeCSV,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockStore := mock_store.NewMockFileUploadStore(t)
			service := NewFileImportService(mockS3Client, mockStore, "zamp-dev-us-application-platform")
			ctx := context.Background()

			if tt.name != "error - no user in context" {
				ctx = apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{orgId})
			}

			tt.setupMock(mockS3Client, mockStore)

			result, err := service.InitiateFileImport(ctx, tt.orgId, tt.fileName, tt.fileType)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, uploadURL, result.UploadURL)
				assert.Equal(t, fileUpload.ID, result.FileUploadID)
				assert.Equal(t, tt.fileType, result.FileType)
			}

			mockS3Client.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAcknowledgeFileImportCompletion(t *testing.T) {
	userId := uuid.New()
	fileUploadId := uuid.New()
	storagePath := "test/path/file.csv"

	fileUpload := &models.FileUpload{
		ID:               fileUploadId,
		UploadedByUserID: userId,
		Status:           models.FileUploadStatusPending,
		Expiry:           time.Now().Add(1 * time.Hour),
		StorageFilePath:  storagePath,
	}

	updatedFileUpload := &models.FileUpload{
		ID:               fileUploadId,
		UploadedByUserID: userId,
		Status:           models.FileUploadStatusCompleted,
		StorageFilePath:  storagePath,
	}

	fileDetails := &s3.FileInfo{
		Key:  storagePath,
		Size: 1000,
	}

	tests := []struct {
		name      string
		setupMock func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore)
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{*fileUpload}, nil)
				mockS3Client.On("GetFileDetails", mock.Anything, mock.Anything, fileUpload.StorageFilePath).Return(fileDetails, nil)
				mockStore.On("UpdateFileUploadStatus", mock.Anything, fileUploadId, models.FileUploadStatusCompleted).Return(updatedFileUpload, nil)
			},
			wantErr: false,
		},
		{
			name: "error - no user in context",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				// No mocks needed
			},
			wantErr: true,
		},
		{
			name: "error - file upload not found",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{}, nil)
			},
			wantErr: true,
		},
		{
			name: "error - file upload already processed",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				processedFileUpload := &models.FileUpload{
					ID:               fileUploadId,
					UploadedByUserID: userId,
					Status:           models.FileUploadStatusCompleted,
					Expiry:           time.Now().Add(1 * time.Hour),
				}
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{*processedFileUpload}, nil)
			},
			wantErr: true,
		},
		{
			name: "error - file upload expired",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				expiredFileUpload := &models.FileUpload{
					ID:               fileUploadId,
					UploadedByUserID: userId,
					Status:           models.FileUploadStatusPending,
					Expiry:           time.Now().Add(-1 * time.Hour),
				}
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{*expiredFileUpload}, nil)
			},
			wantErr: true,
		},
		{
			name: "error - unauthorized user",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				unauthorizedFileUpload := &models.FileUpload{
					ID:               fileUploadId,
					UploadedByUserID: uuid.New(), // Different user
					Status:           models.FileUploadStatusPending,
					Expiry:           time.Now().Add(1 * time.Hour),
				}
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{*unauthorizedFileUpload}, nil)
			},
			wantErr: true,
		},
		{
			name: "error - s3 file not found",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{*fileUpload}, nil)
				mockS3Client.On("GetFileDetails", mock.Anything, mock.Anything, fileUpload.StorageFilePath).Return(nil, fmt.Errorf("file not found"))
			},
			wantErr: true,
		},
		{
			name: "error - empty file",
			setupMock: func(mockS3Client *mock_s3.MockS3Client, mockStore *mock_store.MockFileUploadStore) {
				mockStore.On("GetFileUploadByIds", mock.Anything, []uuid.UUID{fileUploadId}).Return([]models.FileUpload{*fileUpload}, nil)
				emptyFileDetails := &s3.FileInfo{
					Key:  storagePath,
					Size: 0,
				}
				mockS3Client.On("GetFileDetails", mock.Anything, mock.Anything, fileUpload.StorageFilePath).Return(emptyFileDetails, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockS3Client := mock_s3.NewMockS3Client(t)
			mockStore := mock_store.NewMockFileUploadStore(t)
			service := NewFileImportService(mockS3Client, mockStore, "zamp-dev-us-application-platform")
			ctx := context.Background()

			if tt.name != "error - no user in context" {
				ctx = apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{})
			}

			tt.setupMock(mockS3Client, mockStore)

			result, err := service.AcknowledgeFileImportCompletion(ctx, fileUploadId)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, updatedFileUpload, result)
			}

			mockS3Client.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}
