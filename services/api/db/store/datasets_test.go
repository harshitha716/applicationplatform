package store

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetDatasetById(t *testing.T) {
	t.Parallel()

	datasetID := uuid.New()
	orgID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		datasetID string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:      "success",
			datasetID: datasetID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id",
					"title",
					"description",
					"organization_id",
					"created_by",
					"metadata",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					datasetID,
					"Test Dataset",
					"Test Description",
					orgID,
					userID,
					[]byte(`{"key": "value"}`),
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT (.+) FROM "datasets"`).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:      "not found",
			datasetID: uuid.New().String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "datasets"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:      "invalid id",
			datasetID: "invalid-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "datasets"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			dataset, err := store.GetDatasetById(context.Background(), tt.datasetID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dataset)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, dataset)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetDatasetsAll(t *testing.T) {
	t.Parallel()

	// Setup common test data
	orgID1 := uuid.New()
	orgID2 := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	datasetID1 := uuid.New()
	datasetID2 := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		filters   models.DatasetFilters
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantErr   bool
	}{
		{
			name:    "no filters - returns all datasets",
			filters: models.DatasetFilters{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id", "title", "organization_id", "created_by", "created_at",
				}).AddRow(
					datasetID1, "Dataset 1", orgID1, userID1, now,
				).AddRow(
					datasetID2, "Dataset 2", orgID2, userID2, now,
				)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE deleted_at IS NULL`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "pagination - first page",
			filters: models.DatasetFilters{
				Page:  1,
				Limit: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id", "title",
				}).AddRow(
					datasetID1, "Dataset 1",
				).AddRow(
					datasetID2, "Dataset 2",
				)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE deleted_at IS NULL LIMIT \$1`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "pagination - second page",
			filters: models.DatasetFilters{
				Page:  2,
				Limit: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id", "title",
				}).AddRow(
					datasetID1, "Dataset 1",
				)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE deleted_at IS NULL LIMIT \$1 OFFSET \$2`).
					WithArgs(2, 2).
					WillReturnRows(rows)
			},
			wantCount: 1,
		},
		{
			name: "sorting - single column ascending",
			filters: models.DatasetFilters{
				SortParams: []models.DatasetSortParam{{
					Column: "created_at",
					Desc:   false,
				}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id", "title", "created_at",
				}).AddRow(
					datasetID1, "Dataset 1", now,
				).AddRow(
					datasetID2, "Dataset 2", now.Add(time.Hour),
				)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE deleted_at IS NULL ORDER BY "created_at"`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "sorting - multiple columns mixed order",
			filters: models.DatasetFilters{
				SortParams: []models.DatasetSortParam{
					{Column: "created_at", Desc: true},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id", "title", "created_at",
				}).AddRow(
					datasetID1, "Dataset 1", now,
				).AddRow(
					datasetID2, "Dataset 2", now.Add(time.Hour),
				)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE deleted_at IS NULL ORDER BY "created_at" DESC`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "combine all filters",
			filters: models.DatasetFilters{
				OrganizationIds: []uuid.UUID{orgID1},
				DatasetIds:      []uuid.UUID{datasetID1},
				CreatedBy:       []uuid.UUID{userID1},
				Page:            1,
				Limit:           10,
				SortParams: []models.DatasetSortParam{
					{Column: "created_at", Desc: true},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"dataset_id", "title", "organization_id", "created_by",
				}).AddRow(
					datasetID1, "Dataset 1", orgID1, userID1,
				)
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE organization_id IN \(\$1\) AND dataset_id IN \(\$2\) AND created_by IN \(\$3\) AND deleted_at IS NULL ORDER BY "created_at" DESC LIMIT \$4`).
					WithArgs(orgID1, datasetID1, userID1, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
		},
		{
			name: "empty result set",
			filters: models.DatasetFilters{
				OrganizationIds: []uuid.UUID{uuid.New()}, // non-existent org
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE organization_id IN \(\$1\) AND deleted_at IS NULL`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"dataset_id"}))
			},
			wantCount: 0,
		},
		{
			name: "invalid sort column",
			filters: models.DatasetFilters{
				SortParams: []models.DatasetSortParam{{
					Column: "invalid_column",
					Desc:   false,
				}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "datasets" WHERE deleted_at IS NULL ORDER BY "invalid_column"`).
					WillReturnError(gorm.ErrInvalidField)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// Execute
			datasets, err := store.GetDatasetsAll(context.Background(), tt.filters)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, datasets)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, datasets, tt.wantCount)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
