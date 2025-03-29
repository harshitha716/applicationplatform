package models

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWidget_TableName(t *testing.T) {
	t.Parallel()
	widget := Widget{}
	assert.Equal(t, "widgets", widget.TableName())
}

func TestWidget_GetQueryFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	org1ID := uuid.New()
	org2ID := uuid.New()
	orgIDs := []uuid.UUID{org1ID, org2ID}

	tests := []struct {
		name            string
		userId          uuid.UUID
		organizationIDs []uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock)
		wantErr         bool
	}{
		{
			name:            "returns unmodified db instance",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"type"}).AddRow("chart")
				mock.ExpectQuery(`SELECT \* FROM "widgets"`).
					WillReturnRows(rows)
			},
		},
		{
			name:            "handles empty organization list",
			userId:          userId,
			organizationIDs: []uuid.UUID{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "widgets"`).
					WillReturnRows(sqlmock.NewRows([]string{"widget_id"}))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupTestDB(t)
			widget := &Widget{Type: "chart"}
			tt.setupMock(mock)

			query := widget.GetQueryFilters(db.Model(&Widget{}), tt.userId, tt.organizationIDs)

			// Verify query executes without error
			var results []Widget
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

func TestStructImplementsBaseModel_Widget(t *testing.T) {
	var _ pgclient.BaseModel = &Widget{}
}

func TestBeforeCreate_Widgets(t *testing.T) {

	db, _ := setupTestDB(t)

	user := &Widget{}

	err := user.BeforeCreate(db)

	assert.NotNil(t, err)
	assert.Equal(t, "insert forbidden", err.Error())

}
