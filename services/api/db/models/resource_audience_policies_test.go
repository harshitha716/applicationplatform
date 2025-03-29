package models

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestResourceAudiencePolicy_TableName(t *testing.T) {
	t.Parallel()
	policy := ResourceAudiencePolicy{}
	assert.Equal(t, "resource_audience_policies", policy.TableName())
}

func TestResourceAudiencePolicy_GetAccessControlFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	org1ID := uuid.New()
	org2ID := uuid.New()
	orgIDs := []uuid.UUID{org1ID, org2ID}
	policyID := uuid.New()

	tests := []struct {
		name            string
		userId          uuid.UUID
		organizationIDs []uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock)
		wantErr         bool
	}{
		{
			name:            "successful query execution",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"resource_audience_policy_id"}).AddRow(policyID)
				mock.ExpectQuery(`SELECT \* FROM "resource_audience_policies" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnRows(rows)
			},
		},
		{
			name:            "empty result set",
			userId:          userId,
			organizationIDs: []uuid.UUID{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "resource_audience_policies" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_policy_id"}))
			},
		},
		{
			name:            "database error",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "resource_audience_policies" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			db, mock := setupTestDB(t)
			policy := &ResourceAudiencePolicy{ID: policyID}
			tt.setupMock(mock)

			// Execute query
			query := policy.GetQueryFilters(db.Model(&ResourceAudiencePolicy{}), tt.userId, tt.organizationIDs)

			// Verify query execution
			var results []ResourceAudiencePolicy
			err := query.Find(&results).Error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStructImplementsBaseModel_RAP(t *testing.T) {
	var _ pgclient.BaseModel = &ResourceAudiencePolicy{}
}

func TestBeforeCreate_RAP(t *testing.T) {

	db, _ := setupTestDB(t)

	user := &ResourceAudiencePolicy{}

	err := user.BeforeCreate(db)

	assert.Nil(t, err)

}

func TestBeforeUpdate_RAP(t *testing.T) {

	db, mock := setupTestDB(t)

	userId := uuid.New()

	ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{uuid.New()})

	rap := &ResourceAudiencePolicy{
		ResourceType: "organization",
		ResourceID:   uuid.New(),
	}

	// successful query execution
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
		WithArgs(rap.ResourceID, rap.ResourceType, userId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(rap.ResourceID))

	err := rap.BeforeUpdate(db.WithContext(ctx))

	assert.Nil(t, err)

	// no frap access
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
		WithArgs(rap.ResourceID, rap.ResourceType, rap.ResourceID).
		WillReturnRows(sqlmock.NewRows([]string{"resource_id"}))

	err = rap.BeforeUpdate(db.WithContext(ctx))

	assert.NotNil(t, err)

	// db error
	err = nil
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
		WillReturnError(gorm.ErrInvalidDB)
	err = rap.BeforeUpdate(db.WithContext(ctx))
	assert.NotNil(t, err)

}

func TestBeforeDelete_RAP(t *testing.T) {
	db, mock := setupTestDB(t)

	userId := uuid.New()

	ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{uuid.New()})

	rap := &ResourceAudiencePolicy{
		ResourceType: "organization",
		ResourceID:   uuid.New(),
	}

	// successful query execution
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
		WithArgs(rap.ResourceID, rap.ResourceType, userId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"resource_id"}).AddRow(rap.ResourceID))

	err := rap.BeforeDelete(db.WithContext(ctx))

	assert.Nil(t, err)

	// no frap access
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
		WithArgs(rap.ResourceID, rap.ResourceType, rap.ResourceID).
		WillReturnRows(sqlmock.NewRows([]string{"resource_id"}))

	err = rap.BeforeDelete(db.WithContext(ctx))

	assert.NotNil(t, err)

	// db error
	err = nil
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE (resource_id = $1 ANd resource_type = $2 AND ((resource_type = 'page' AND privilege = 'admin') OR (resource_type = 'dataset' AND privilege = 'admin') OR (resource_type = 'organization' AND privilege = 'system_admin') OR (resource_type = 'connection' AND privilege = 'admin')) AND user_id = $3 AND deleted_at IS NULL) LIMIT $4`)).
		WillReturnError(gorm.ErrInvalidDB)
	err = rap.BeforeDelete(db.WithContext(ctx))
	assert.NotNil(t, err)
}
