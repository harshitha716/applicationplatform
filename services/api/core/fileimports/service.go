package fileimports

import (
	"context"
	"fmt"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileImportService interface {
	InitiateFileImport(ctx context.Context, orgId uuid.UUID, fileName string, fileType models.FileType) (*FileImportSignedURL, error)
	AcknowledgeFileImportCompletion(ctx context.Context, fileUploadId uuid.UUID) (*models.FileUpload, error)
	GetFileUploadByIds(ctx context.Context, ids []uuid.UUID) ([]models.FileUpload, error)
}

type fileImportService struct {
	s3Client   s3.S3Client
	store      store.FileUploadStore
	bucketName string
}

func NewFileImportService(s3Client s3.S3Client, store store.FileUploadStore, bucketName string) FileImportService {
	return &fileImportService{
		s3Client:   s3Client,
		store:      store,
		bucketName: bucketName,
	}
}

func (s *fileImportService) InitiateFileImport(ctx context.Context, orgId uuid.UUID, fileName string, fileType models.FileType) (*FileImportSignedURL, error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	_, currentUserId, orgIds := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		ctxlogger.Error("no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}

	if len(orgIds) != 1 || orgIds[0] != orgId {
		ctxlogger.Error("user does not have access to org", zap.String("org_id", orgId.String()))
		return nil, fmt.Errorf("user does not have access to org %s", orgId)
	}

	// generate file import path
	fileUploadId := uuid.New()
	fileUploadPath := getFileImportPath(orgId, fileUploadId, fileName)

	expiry := time.Now().Add(SIGNED_URL_EXPIRY)
	uploadURL, err := s.s3Client.GenerateUploadURL(ctx, s.bucketName, fileUploadPath, expiry, CONTENT_TYPE_OCTET_STREAM)
	if err != nil {
		ctxlogger.Error("error generating upload url", zap.Error(err))
		return nil, err
	}

	fileUpload, err := s.store.CreateFileUpload(ctx, &models.FileUpload{
		ID:               fileUploadId,
		OrganizationID:   orgId,
		UploadedByUserID: *currentUserId,
		Name:             fileName,
		PresignedURL:     uploadURL,
		Expiry:           expiry,
		Status:           models.FileUploadStatusPending,
		StorageFilePath:  fileUploadPath,
		FileType:         fileType,
		StorageBucket:    s.bucketName,
		StorageProvider:  models.StorageTypeS3,
	})

	if err != nil {
		ctxlogger.Error("error creating file upload", zap.Error(err))
		return nil, err
	}

	ctxlogger.Info("file upload created", zap.String("file_upload_id", fileUpload.ID.String()))
	return &FileImportSignedURL{
		UploadURL:    uploadURL,
		Key:          getFileImportPath(orgId, fileUploadId, fileName),
		FileName:     fileName,
		FileType:     fileType,
		FileUploadID: fileUpload.ID,
	}, nil

}

func (s *fileImportService) AcknowledgeFileImportCompletion(ctx context.Context, fileUploadId uuid.UUID) (*models.FileUpload, error) {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	// Get current user from context
	_, currentUserId, _ := apicontext.GetAuthFromContext(ctx)
	if currentUserId == nil {
		ctxlogger.Error("no user found in context")
		return nil, fmt.Errorf("unauthorized")
	}

	// TODO: Check if user has access to file upload
	fileUploads, err := s.store.GetFileUploadByIds(ctx, []uuid.UUID{fileUploadId})
	if err != nil {
		ctxlogger.Error("error getting file upload", zap.Error(err))
		return nil, fmt.Errorf("File upload not found")
	}

	if len(fileUploads) == 0 {
		ctxlogger.Error("file upload not found", zap.String("file_upload_id", fileUploadId.String()))
		return nil, fmt.Errorf("File upload not found")
	}

	fileUpload := fileUploads[0]

	// Check if file upload is pending
	if fileUpload.Status != models.FileUploadStatusPending {
		ctxlogger.Error("file upload is not pending", zap.String("file_upload_id", fileUpload.ID.String()))
		return nil, fmt.Errorf("This file upload is already processed")
	}

	// Check if file upload has expired
	if fileUpload.Expiry.Before(time.Now()) {
		ctxlogger.Error("file upload has expired", zap.String("file_upload_id", fileUpload.ID.String()))
		return nil, fmt.Errorf("This file upload has expired. Please re-initiate the file import.")
	}

	// Check if user is the owner of the file upload
	if fileUpload.UploadedByUserID != *currentUserId {
		ctxlogger.Error("user does not have access to file upload", zap.String("file_upload_id", fileUpload.ID.String()))
		return nil, fmt.Errorf("user does not have access to file upload")
	}

	// check in s3 if the file exists
	fileDetails, err := s.s3Client.GetFileDetails(ctx, s.bucketName, fileUpload.StorageFilePath)
	if err != nil {
		ctxlogger.Error("error getting file details", zap.Error(err))
		return nil, fmt.Errorf("error getting file details")
	}

	// check if the file is empty
	if fileDetails.Size == 0 {
		ctxlogger.Error("file is empty", zap.String("file_upload_id", fileUpload.ID.String()))
		return nil, fmt.Errorf("file is empty")
	}

	// update the file upload status
	updatedFileUpload, err := s.store.UpdateFileUploadStatus(ctx, fileUpload.ID, models.FileUploadStatusCompleted)
	if err != nil {
		ctxlogger.Error("error updating file upload status", zap.Error(err))
		return nil, fmt.Errorf("error updating file upload status")
	}

	return updatedFileUpload, nil

}

func (s *fileImportService) GetFileUploadByIds(ctx context.Context, ids []uuid.UUID) ([]models.FileUpload, error) {
	fileUploads, err := s.store.GetFileUploadByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return fileUploads, nil
}
