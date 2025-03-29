package models

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFlattenedResourceAudiencePolicy_TableName(t *testing.T) {
	tests := []struct {
		name     string
		model    *FlattenedResourceAudiencePolicy
		expected string
	}{
		{
			name:     "should return correct table name",
			model:    &FlattenedResourceAudiencePolicy{},
			expected: "flattened_resource_audience_policies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.TableName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlattenedResourceAudiencePolicy_GetQueryFilters(t *testing.T) {

	userId := uuid.New()
	orgIds := []uuid.UUID{uuid.New(), uuid.New()}

	tests := []struct {
		name       string
		model      *FlattenedResourceAudiencePolicy
		userId     uuid.UUID
		orgIds     []uuid.UUID
		setupMock  func(mock sqlmock.Sqlmock)
		verifyFunc func(t *testing.T, db *gorm.DB)
	}{
		{
			name:   "should add user_id filter to query",
			model:  &FlattenedResourceAudiencePolicy{},
			userId: userId,
			orgIds: orgIds,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE EXISTS`)).
					WithArgs(userId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userId))
			},
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var result FlattenedResourceAudiencePolicy
				err := db.Find(&result).Error
				assert.NoError(t, err)
			},
		},
		{
			name:   "should handle empty org IDs",
			model:  &FlattenedResourceAudiencePolicy{},
			userId: userId,
			orgIds: []uuid.UUID{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE EXISTS`)).
					WithArgs(userId, userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userId))
			},
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var result FlattenedResourceAudiencePolicy
				err := db.Find(&result).Error
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test database and mock
			db, mock := setupTestDB(t)

			// Setup mock expectations
			tt.setupMock(mock)

			// Execute GetQueryFilters
			filteredDB := tt.model.GetQueryFilters(db, tt.userId, tt.orgIds)

			// Verify the results
			tt.verifyFunc(t, filteredDB)

			// Verify all expectations were met
			err := mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestStructImplementsBaseModel_Frap(t *testing.T) {
	var _ pgclient.BaseModel = &FlattenedResourceAudiencePolicy{}
}

func TestBeforeCreate_Frap(t *testing.T) {

	db, _ := setupTestDB(t)

	user := &FlattenedResourceAudiencePolicy{}

	err := user.BeforeCreate(db)

	assert.NotNil(t, err)
	assert.Equal(t, "insert forbidden", err.Error())

}
