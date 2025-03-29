package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFileUpload_GetQueryFilters(t *testing.T) {
	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New(), uuid.New()}

	tests := []struct {
		name       string
		model      *FileUpload
		userId     uuid.UUID
		orgIds     []uuid.UUID
		setupMock  func(mock sqlmock.Sqlmock)
		verifyFunc func(t *testing.T, db *gorm.DB)
	}{
		{
			name:   "should add uploaded_by_user_id filter to query",
			model:  &FileUpload{},
			userId: userId,
			orgIds: orgIds,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "file_uploads" WHERE uploaded_by_user_id = $1 AND deleted_at IS NULL`)).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"uploaded_by_user_id"}).AddRow(userId))
			},
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var result []FileUpload
				err := db.Find(&result).Error
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)

			tt.setupMock(mock)

			filteredDB := tt.model.GetQueryFilters(db, tt.userId, tt.orgIds)
			tt.verifyFunc(t, filteredDB)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFileUpload_BeforeCreate(t *testing.T) {
	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name          string
		model         *FileUpload
		currentUserId *uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock)
		wantErr       bool
		errMsg        string
	}{
		{
			name: "should succeed with valid organization access",
			model: &FileUpload{
				OrganizationID: orgId,
			},
			currentUserId: &userId,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id"}).
						AddRow("organization", orgId, userId))
			},
			wantErr: false,
		},
		{
			name: "should fail without user in context",
			model: &FileUpload{
				OrganizationID: orgId,
			},
			setupMock: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
			errMsg:    "no user id found in context",
		},
		{
			name: "should fail without organization access",
			model: &FileUpload{
				OrganizationID: orgId,
			},
			currentUserId: &userId,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_type", "resource_id", "user_id"}))
			},
			wantErr: true,
			errMsg:  "organization access forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)

			ctx := context.Background()
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.currentUserId, []uuid.UUID{})
			}
			db = db.WithContext(ctx)

			tt.setupMock(mock)

			err := tt.model.BeforeCreate(db)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userId, tt.model.UploadedByUserID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFileUpload_BeforeUpdate(t *testing.T) {
	userId := uuid.New()
	differentUserId := uuid.New()

	tests := []struct {
		name     string
		model    *FileUpload
		setupCtx func() context.Context
		wantErr  bool
		errMsg   string
	}{
		{
			name: "should succeed when user owns the upload",
			model: &FileUpload{
				UploadedByUserID: userId,
			},
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{})
			},
			wantErr: false,
		},
		{
			name: "should fail without user in context",
			model: &FileUpload{
				UploadedByUserID: userId,
			},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantErr: true,
			errMsg:  "no user id found in context",
		},
		{
			name: "should fail when user doesn't own the upload",
			model: &FileUpload{
				UploadedByUserID: userId,
			},
			setupCtx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", differentUserId, []uuid.UUID{})
			},
			wantErr: true,
			errMsg:  "user does not have permission to update this file upload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)

			ctx := tt.setupCtx()
			db = db.WithContext(ctx)

			err := tt.model.BeforeUpdate(db)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
