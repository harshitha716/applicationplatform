package store

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetSheetById(t *testing.T) {
	t.Parallel()

	sheetId := uuid.New()
	pageId := uuid.New()
	name := "Test Sheet"
	description := "Test Description"
	now := time.Now().UTC()
	sheetConfig := json.RawMessage(`{"key": "value"}`)

	tests := []struct {
		name      string
		sheetId   uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      *models.Sheet
		wantErr   bool
	}{
		{
			name:    "success",
			sheetId: sheetId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Sheet query
				sheetRows := sqlmock.NewRows([]string{
					"sheet_id",
					"name",
					"description",
					"created_at",
					"updated_at",
					"deleted_at",
					"fractional_index",
					"page_id",
					"sheet_config",
				}).AddRow(
					sheetId,
					name,
					description,
					now,
					now,
					nil,
					1.0,
					pageId,
					sheetConfig,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sheets" WHERE sheet_id = $1 ORDER BY "sheets"."sheet_id" LIMIT $2`)).
					WithArgs(sheetId, 1).
					WillReturnRows(sheetRows)

				// Widget instances query
				widgetRows := sqlmock.NewRows([]string{
					"widget_instance_id",
					"widget_type",
					"sheet_id",
					"title",
					"data_mappings",
					"default_filters",
					"created_at",
					"updated_at",
					"deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "widget_instances" WHERE "widget_instances"."sheet_id" = $1`)).
					WithArgs(sheetId).
					WillReturnRows(widgetRows)
			},
			want: &models.Sheet{
				ID:              sheetId,
				Name:            name,
				Description:     &description,
				CreatedAt:       now,
				UpdatedAt:       now,
				FractionalIndex: 1.0,
				PageId:          pageId,
				SheetConfig:     sheetConfig,
				WidgetInstances: []models.WidgetInstance{},
			},
			wantErr: false,
		},
		{
			name:    "not found",
			sheetId: sheetId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sheets" WHERE sheet_id = $1 ORDER BY "sheets"."sheet_id" LIMIT $2`)).
					WithArgs(sheetId, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "database error",
			sheetId: sheetId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sheets" WHERE sheet_id = $1 ORDER BY "sheets"."sheet_id" LIMIT $2`)).
					WithArgs(sheetId, 1).
					WillReturnError(gorm.ErrInvalidDB)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			got, err := store.GetSheetById(context.Background(), tt.sheetId)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.CreatedAt.UTC(), got.CreatedAt.UTC())
			assert.Equal(t, tt.want.UpdatedAt.UTC(), got.UpdatedAt.UTC())
			assert.Equal(t, tt.want.DeletedAt, got.DeletedAt)
			assert.Equal(t, tt.want.FractionalIndex, got.FractionalIndex)
			assert.Equal(t, tt.want.PageId, got.PageId)
			assert.Equal(t, string(tt.want.SheetConfig), string(got.SheetConfig))
			assert.Equal(t, len(tt.want.WidgetInstances), len(got.WidgetInstances))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetSheetsAll(t *testing.T) {
	t.Parallel()

	// Setup common test data
	pageId1 := uuid.New()
	pageId2 := uuid.New()
	sheetID1 := uuid.New()
	sheetID2 := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		filters   models.SheetFilters
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantErr   bool
	}{
		{
			name:    "no filters - returns all sheets",
			filters: models.SheetFilters{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id", "name", "page_id", "created_at",
				}).AddRow(
					sheetID1, "Sheet 1", pageId1, now,
				).AddRow(
					sheetID2, "Sheet 2", pageId2, now,
				)
				mock.ExpectQuery(`SELECT \* FROM "sheets"`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "include widget instances",
			filters: models.SheetFilters{
				IncludeWidgetInstances: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				sheetRows := sqlmock.NewRows([]string{
					"sheet_id", "name", "page_id", "created_at",
				}).AddRow(
					sheetID1, "Sheet 1", pageId1, now,
				)

				widgetRows := sqlmock.NewRows([]string{
					"id", "sheet_id", "created_at",
				}).AddRow(
					uuid.New(), sheetID1, now,
				).AddRow(
					uuid.New(), sheetID1, now.Add(time.Hour),
				)

				mock.ExpectQuery(`SELECT \* FROM "sheets"`).
					WillReturnRows(sheetRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "widget_instances" WHERE "widget_instances"."sheet_id" = $1 ORDER BY widget_instances.created_at ASC`)).
					WithArgs(sheetID1).
					WillReturnRows(widgetRows)
			},
			wantCount: 1,
		},
		{
			name: "pagination - first page",
			filters: models.SheetFilters{
				Page:  1,
				Limit: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"sheet_id", "name",
				}).AddRow(
					sheetID1, "Sheet 1",
				).AddRow(
					sheetID2, "Sheet 2",
				)
				mock.ExpectQuery(`SELECT \* FROM "sheets" LIMIT \$1`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "pagination - second page",
			filters: models.SheetFilters{
				Page:  2,
				Limit: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"sheet_id", "name",
				}).AddRow(
					sheetID1, "Sheet 1",
				)
				mock.ExpectQuery(`SELECT \* FROM "sheets" LIMIT \$1 OFFSET \$2`).
					WithArgs(2, 2).
					WillReturnRows(rows)
			},
			wantCount: 1,
		},
		{
			name: "sorting - single column ascending",
			filters: models.SheetFilters{
				SortParams: []models.SheetSortParams{{
					Column: "created_at",
					Desc:   false,
				}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"sheet_id", "name", "created_at",
				}).AddRow(
					sheetID1, "Sheet 1", now,
				).AddRow(
					sheetID2, "Sheet 2", now.Add(time.Hour),
				)
				mock.ExpectQuery(`SELECT \* FROM "sheets" ORDER BY "created_at"`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "combine all filters",
			filters: models.SheetFilters{
				PageIds:  []uuid.UUID{pageId1},
				SheetIds: []uuid.UUID{sheetID1},
				Page:     1,
				Limit:    10,
				SortParams: []models.SheetSortParams{
					{Column: "created_at", Desc: true},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"sheet_id", "name", "page_id",
				}).AddRow(
					sheetID1, "Page 1", pageId1,
				)
				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE page_id IN \(\$1\) AND sheet_id IN \(\$2\) ORDER BY "created_at" DESC LIMIT \$3`).
					WithArgs(pageId1, sheetID1, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
		},
		{
			name: "empty result set",
			filters: models.SheetFilters{
				PageIds: []uuid.UUID{uuid.New()}, // non-existent org
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE page_id IN \(\$1\)`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"sheet_id"}))
			},
			wantCount: 0,
		},
		{
			name: "invalid sort column",
			filters: models.SheetFilters{
				SortParams: []models.SheetSortParams{{
					Column: "invalid_column",
					Desc:   false,
				}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "sheets" ORDER BY "invalid_column"`).
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
			sheets, err := store.GetSheetsAll(context.Background(), tt.filters)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, sheets)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, sheets, tt.wantCount)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_CreateSheet(t *testing.T) {
	t.Parallel()

	pageId := uuid.New()
	userId := uuid.New()
	name := "Test Sheet"
	description := "Test Description"

	tests := []struct {
		name      string
		sheet     models.Sheet
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful sheet creation",
			sheet: models.Sheet{
				PageId:      pageId,
				Name:        name,
				Description: &description,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("page", pageId, userId, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
					}).AddRow(uuid.New(), "page", pageId, userId))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "sheets" ("sheet_id","name","description","created_at","updated_at","deleted_at","fractional_index","page_id","sheet_config") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,(NULL))`)).
					WithArgs(sqlmock.AnyArg(), name, description, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, float64(0), pageId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			sheet: models.Sheet{
				PageId:      pageId,
				Name:        name,
				Description: &description,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				// Fix authorization check query and arguments
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("page", pageId, userId, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
					}).AddRow(uuid.New(), "page", pageId, userId))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "sheets" ("sheet_id","name","description","created_at","updated_at","deleted_at","page_id") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
					WithArgs(sqlmock.AnyArg(), name, description, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, pageId).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
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

			// Create context with user authentication
			ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})

			// Execute
			sheet, err := store.CreateSheet(ctx, tt.sheet)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, sheet)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, sheet)
			assert.Equal(t, tt.sheet.Name, sheet.Name)
			assert.Equal(t, tt.sheet.Description, sheet.Description)
			assert.Equal(t, tt.sheet.PageId, sheet.PageId)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_UpdateSheet(t *testing.T) {
	t.Parallel()

	sheetId := uuid.New()
	pageId := uuid.New()
	name := "Updated Sheet"
	description := "Updated Description"
	now := time.Now().UTC()
	sheetConfig := json.RawMessage(`{"key": "updated_value"}`)

	tests := []struct {
		name      string
		sheet     *models.Sheet
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful sheet update",
			sheet: &models.Sheet{
				ID:              sheetId,
				PageId:          pageId,
				Name:            name,
				Description:     &description,
				CreatedAt:       now,
				UpdatedAt:       now,
				FractionalIndex: 1.0,
				SheetConfig:     sheetConfig,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sheets" SET "name"=$1,"description"=$2,"created_at"=$3,"updated_at"=$4,"deleted_at"=$5,"fractional_index"=$6,"page_id"=$7,"sheet_config"=$8 WHERE "sheet_id" = $9`)).
					WithArgs(
						name,             // name
						description,      // description
						now,              // created_at
						sqlmock.AnyArg(), // updated_at
						nil,              // deleted_at
						1.0,              // fractional_index
						pageId,           // page_id
						sheetConfig,      // sheet_config
						sheetId,          // WHERE sheet_id = ?
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			sheet: &models.Sheet{
				ID:          sheetId,
				PageId:      pageId,
				Name:        name,
				Description: &description,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sheets" SET "name"=$1,"description"=$2,"created_at"=$3,"updated_at"=$4,"deleted_at"=$5,"fractional_index"=$6,"page_id"=$7,"sheet_config"=$8 WHERE "sheet_id" = $9`)).
					WithArgs(
						name,             // name
						description,      // description
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
						nil,              // deleted_at
						0.0,              // fractional_index
						pageId,           // page_id
						nil,              // sheet_config
						sheetId,          // WHERE sheet_id = ?
					).WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
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
			sheet, err := store.UpdateSheet(context.Background(), tt.sheet)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, sheet)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, sheet)
			assert.Equal(t, tt.sheet.ID, sheet.ID)
			assert.Equal(t, tt.sheet.Name, sheet.Name)
			assert.Equal(t, tt.sheet.Description, sheet.Description)
			assert.Equal(t, tt.sheet.PageId, sheet.PageId)
			assert.Equal(t, string(tt.sheet.SheetConfig), string(sheet.SheetConfig))

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
