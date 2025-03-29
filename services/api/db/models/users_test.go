package models

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return db, mock
}

func TestUser_TableName(t *testing.T) {
	t.Parallel()
	user := User{}
	assert.Equal(t, "users_with_traits", user.TableName())
}

func TestUser_GetQueryFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	org1ID := uuid.New()
	orgIDs := []uuid.UUID{org1ID}

	tests := []struct {
		name            string
		userId          uuid.UUID
		organizationIDs []uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock)
		wantErr         bool
	}{
		{
			name:            "multiple organization memberships",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'organization' AND frap.resource_id IN ($1) AND ( frap.user_id = users_with_traits.user_id OR ( EXISTS ( SELECT 1 FROM "app"."organization_membership_requests" omr WHERE omr.organization_id = frap.resource_id AND omr.user_id = users_with_traits.user_id ) ) ) AND frap.deleted_at IS NULL ) OR users_with_traits.user_id = $2`)).
					WithArgs(sqlmock.AnyArg(), userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(userId, "test@example.com", "Test User"))
			},
		},
		{
			name:            "no organization memberships",
			userId:          userId,
			organizationIDs: []uuid.UUID{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users_with_traits" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'organization' AND frap.resource_id IN ($1) AND ( frap.user_id = users_with_traits.user_id OR ( EXISTS ( SELECT 1 FROM "app"."organization_membership_requests" omr WHERE omr.organization_id = frap.resource_id AND omr.user_id = users_with_traits.user_id ) ) ) AND frap.deleted_at IS NULL ) OR users_with_traits.user_id = $2`)).
					WithArgs(sqlmock.AnyArg(), userId).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "name"}))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupUserTestDB(t)
			user := &User{ID: userId}
			tt.setupMock(mock)

			baseQuery := db.Model(&User{})
			query := user.GetQueryFilters(baseQuery, tt.userId, tt.organizationIDs)

			var results []User
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

func TestStructImplementsBaseModel_Users(t *testing.T) {
	var _ pgclient.BaseModel = &User{}
}

func TestBeforeCreat(t *testing.T) {

	db, _ := setupTestDB(t)

	user := &User{}

	err := user.BeforeCreate(db)

	assert.NotNil(t, err)
	assert.Equal(t, "insert forbidden", err.Error())

}
