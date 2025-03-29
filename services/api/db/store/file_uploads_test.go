package store

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateFileUpload(t *testing.T) {
	fileUploadId := uuid.New()
	userId := uuid.New()
	orgId := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		fileUpload *models.FileUpload
		mockSetup  func(mock sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success",
			fileUpload: &models.FileUpload{
				ID:               fileUploadId,
				OrganizationID:   orgId,
				UploadedByUserID: userId,
				Name:             "test.txt",
				PresignedURL:     "url",
				Expiry:           now.Add(time.Hour * 24),
				StorageProvider:  "provider",
				StorageBucket:    "bucket",
				StorageFilePath:  "path",
				Status:           models.FileUploadStatusPending,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id"}).
						AddRow("organization", orgId, userId))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "file_uploads" ("name","file_type","organization_id","uploaded_by_user_id","presigned_url","expiry","storage_provider","storage_bucket","storage_file_path","status","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id","created_at","updated_at"`)).
					WithArgs("test.txt", "", orgId, userId, "url", now.Add(time.Hour*24), "provider", "bucket", "path", models.FileUploadStatusPending, nil, fileUploadId).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fileUploadId))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			fileUpload: &models.FileUpload{
				ID:               fileUploadId,
				OrganizationID:   orgId,
				UploadedByUserID: userId,
				Name:             "test.txt",
				PresignedURL:     "url",
				Expiry:           now.Add(time.Hour * 24),
				StorageProvider:  "provider",
				StorageBucket:    "bucket",
				StorageFilePath:  "path",
				Status:           models.FileUploadStatusPending,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id"}).
						AddRow("organization", orgId, userId))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "file_uploads" ("name","file_type","organization_id","uploaded_by_user_id","presigned_url","expiry","storage_provider","storage_bucket","storage_file_path","status","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id","created_at","updated_at"`)).
					WithArgs("test.txt", "", orgId, userId, "url", now.Add(time.Hour*24), "provider", "bucket", "path", models.FileUploadStatusPending, nil, fileUploadId).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := apicontext.AddAuthToContext(context.Background(), "admin", userId, []uuid.UUID{orgId})
			tt.mockSetup(mock)

			_, err := store.CreateFileUpload(ctx, tt.fileUpload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetFileUploadByIDs(t *testing.T) {
	fileUploadId := uuid.New()
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name      string
		ids       []uuid.UUID
		mockSetup func(mock sqlmock.Sqlmock)
		want      []models.FileUpload
		wantErr   bool
	}{
		{
			name: "success",
			ids:  []uuid.UUID{fileUploadId},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "file_uploads" WHERE id IN ($1)`)).
					WithArgs(fileUploadId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "organization_id", "uploaded_by_user_id", "name"}).
						AddRow(fileUploadId, orgId, userId, "test.txt"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE "users_with_traits"."user_id" = $1`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(userId, "user@example.com", "Test User"))
			},
			want: []models.FileUpload{
				{
					ID:               fileUploadId,
					OrganizationID:   orgId,
					UploadedByUserID: userId,
					Name:             "test.txt",
					UploadedByUser: &models.User{
						ID:    userId,
						Email: "user@example.com",
						Name:  "Test User",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			ids:  []uuid.UUID{fileUploadId},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "file_uploads" WHERE id IN ($1)`)).
					WithArgs(fileUploadId).
					WillReturnError(assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := context.Background()
			tt.mockSetup(mock)

			got, err := store.GetFileUploadByIds(ctx, tt.ids)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllFileUploads(t *testing.T) {
	fileUploadId := uuid.New()
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		want      []*models.FileUpload
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "file_uploads"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "organization_id", "uploaded_by_user_id", "name"}).
						AddRow(fileUploadId, orgId, userId, "test.txt"))
			},
			want: []*models.FileUpload{
				{
					ID:               fileUploadId,
					OrganizationID:   orgId,
					UploadedByUserID: userId,
					Name:             "test.txt",
				},
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "file_uploads"`)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := context.Background()
			tt.mockSetup(mock)

			got, err := store.GetAllFileUploads(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			if len(got) > 0 {
				assert.Equal(t, tt.want[0].ID, got[0].ID)
				assert.Equal(t, tt.want[0].OrganizationID, got[0].OrganizationID)
				assert.Equal(t, tt.want[0].UploadedByUserID, got[0].UploadedByUserID)
				assert.Equal(t, tt.want[0].Name, got[0].Name)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateFileUploadStatus(t *testing.T) {
	fileUploadId := uuid.New()
	userId := uuid.New()
	status := models.FileUploadStatusCompleted

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "file_uploads" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`)).
					WithArgs(status, sqlmock.AnyArg(), fileUploadId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "file_uploads" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`)).
					WithArgs(status, sqlmock.AnyArg(), fileUploadId).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "error - no user in context",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No DB calls expected
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			ctx := context.Background()
			if tt.name != "error - no user in context" {
				ctx = apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{})
			}
			tt.mockSetup(mock)

			_, err := store.UpdateFileUploadStatus(ctx, fileUploadId, status)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
