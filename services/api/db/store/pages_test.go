package store

import (
	"context"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetPageById(t *testing.T) {
	t.Parallel()

	pageID := uuid.New()
	orgID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		pageId    uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:   "success",
			pageId: pageID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id",
					"name",
					"description",
					"organization_id",
					"created_at",
					"updated_at",
					"deleted_at",
					"fractional_index",
				}).AddRow(
					pageID,
					"Test Page",
					"Test Description",
					orgID,
					now,
					now,
					nil,
					1.5,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pages" WHERE "pages"."page_id" = $1 ORDER BY "pages"."page_id" LIMIT $2`)).
					WithArgs(pageID, 1).
					WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sheets" WHERE "sheets"."page_id" = $1`)).
					WithArgs(pageID).
					WillReturnRows(sqlmock.NewRows([]string{"sheet_id"}))
			},
			wantErr: false,
		},
		{
			name:   "not found",
			pageId: uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "pages"`).
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
			page, err := store.GetPageById(context.Background(), tt.pageId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, page)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, page)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPagesAll(t *testing.T) {
	t.Parallel()

	// Setup common test data
	orgID1 := uuid.New()
	orgID2 := uuid.New()
	pageId1 := uuid.New()
	pageId2 := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		filters   models.PageFilters
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantErr   bool
	}{
		{
			name:    "no filters - returns all pages",
			filters: models.PageFilters{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id", "name", "organization_id", "created_at",
				}).AddRow(
					pageId1, "Page 1", orgID1, now,
				).AddRow(
					pageId2, "Page 2", orgID2, now,
				)
				mock.ExpectQuery(`SELECT \* FROM "pages"`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "pagination - first page",
			filters: models.PageFilters{
				Page:  1,
				Limit: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id", "name",
				}).AddRow(
					pageId1, "Page 1",
				).AddRow(
					pageId2, "Page 2",
				)
				mock.ExpectQuery(`SELECT \* FROM "pages" LIMIT \$1`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "pagination - second page",
			filters: models.PageFilters{
				Page:  2,
				Limit: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id", "name",
				}).AddRow(
					pageId1, "Page 1",
				)
				mock.ExpectQuery(`SELECT \* FROM "pages" LIMIT \$1 OFFSET \$2`).
					WithArgs(2, 2).
					WillReturnRows(rows)
			},
			wantCount: 1,
		},
		{
			name: "sorting - single column ascending",
			filters: models.PageFilters{
				SortParams: []models.PageSortParams{{
					Column: "created_at",
					Desc:   false,
				}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id", "name", "created_at",
				}).AddRow(
					pageId1, "Page 1", now,
				).AddRow(
					pageId2, "Page 2", now.Add(time.Hour),
				)
				mock.ExpectQuery(`SELECT \* FROM "pages" ORDER BY "created_at"`).
					WillReturnRows(rows)
			},
			wantCount: 2,
		},
		{
			name: "combine all filters",
			filters: models.PageFilters{
				OrganizationIds: []uuid.UUID{orgID1},
				PageIds:         []uuid.UUID{pageId1},
				Page:            1,
				Limit:           10,
				SortParams: []models.PageSortParams{
					{Column: "created_at", Desc: true},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id", "name", "organization_id",
				}).AddRow(
					pageId1, "Page 1", orgID1,
				)
				mock.ExpectQuery(`SELECT \* FROM "pages" WHERE organization_id IN \(\$1\) AND page_id IN \(\$2\) ORDER BY "created_at" DESC LIMIT \$3`).
					WithArgs(orgID1, pageId1, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
		},
		{
			name: "empty result set",
			filters: models.PageFilters{
				OrganizationIds: []uuid.UUID{uuid.New()}, // non-existent org
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "pages" WHERE organization_id IN \(\$1\)`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"page_id"}))
			},
			wantCount: 0,
		},
		{
			name: "invalid sort column",
			filters: models.PageFilters{
				SortParams: []models.PageSortParams{{
					Column: "invalid_column",
					Desc:   false,
				}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "pages" ORDER BY "invalid_column"`).
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
			pages, err := store.GetPagesAll(context.Background(), tt.filters)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pages)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, pages, tt.wantCount)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_CreatePage(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name        string
		setupCtx    func(orgId uuid.UUID) context.Context
		pageName    string
		description string
		setupMock   func(mock sqlmock.Sqlmock)
		wantErr     bool
		errMsg      string
	}{
		{
			name: "successful page creation",
			setupCtx: func(orgId uuid.UUID) context.Context {

				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{orgId})
				return ctx
			},
			pageName:    "Test Page",
			description: "Test Description",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
						"created_at",
						"updated_at",
					}).AddRow(
						uuid.New(),
						"organization",
						orgId,
						userId,
						time.Now(),
						time.Now(),
					))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "pages" ("page_id","name","description","created_at","updated_at","deleted_at","fractional_index","organization_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`),
				).WithArgs(sqlmock.AnyArg(), "Test Page", "Test Description", sqlmock.AnyArg(), sqlmock.AnyArg(), nil, float64(0), orgId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "failure - no organization ID in context",
			setupCtx: func(orgId uuid.UUID) context.Context {
				return context.Background()
			},
			pageName:    "Test Page",
			description: "Test Description",
			setupMock: func(mock sqlmock.Sqlmock) {
				// No mock expectations needed as it should fail before DB query
			},
			wantErr: true,
			errMsg:  "organization access forbidden",
		},
		{
			name: "failure - database error",
			setupCtx: func(orgId uuid.UUID) context.Context {
				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{orgId})
				return ctx
			},
			pageName:    "Test Page",
			description: "Test Description",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId, userId, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"resource_type",
						"resource_id",
						"user_id",
						"created_at",
						"updated_at",
					}).AddRow(
						uuid.New(),
						"organization",
						orgId,
						userId,
						time.Now(),
						time.Now(),
					))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "pages" ("page_id","name","description","created_at","updated_at","deleted_at","fractional_index","organization_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`),
				).WithArgs(sqlmock.AnyArg(), "Test Page", "Test Description", sqlmock.AnyArg(), sqlmock.AnyArg(), nil, float64(0), orgId).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}

			// Setup test context
			ctx := tt.setupCtx(orgId)

			// Setup mock expectations
			tt.setupMock(mock)

			// Execute test
			page, err := store.CreatePage(ctx, tt.pageName, tt.description)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, page)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, page)
				assert.Equal(t, tt.pageName, page.Name)
				assert.Equal(t, tt.description, *page.Description)
			}

			// Verify all mock expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPagesByOrganizationId(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	pageId1 := uuid.New()
	pageId2 := uuid.New()
	sheetId1 := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		orgId     uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		wantLen   int
	}{
		{
			name:  "success - multiple pages",
			orgId: orgId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				pageRows := sqlmock.NewRows([]string{
					"page_id",
					"name",
					"description",
					"organization_id",
					"created_at",
					"updated_at",
					"deleted_at",
					"fractional_index",
				}).
					AddRow(pageId1, "Page 1", "Desc 1", orgId, now, now, nil, 1.0).
					AddRow(pageId2, "Page 2", "Desc 2", orgId, now, now, nil, 2.0)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pages" WHERE organization_id = $1`)).
					WithArgs(orgId).
					WillReturnRows(pageRows)

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
				}).AddRow(sheetId1, "Sheet 1", "Desc 1", now, now, nil, 1.0, pageId1, []byte(`{"config": "value"}`))

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sheets" WHERE "sheets"."page_id" IN ($1,$2)`)).
					WithArgs(pageId1, pageId2).
					WillReturnRows(sheetRows)

				widgetRows := sqlmock.NewRows([]string{
					"widget_instance_id",
					"name",
					"description",
					"created_at",
					"updated_at",
					"deleted_at",
					"fractional_index",
					"sheet_id",
					"widget_config",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "widget_instances" WHERE "widget_instances"."sheet_id" = $1`)).
					WithArgs(sheetId1).
					WillReturnRows(widgetRows)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:  "success - no pages",
			orgId: orgId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"page_id",
					"name",
					"description",
					"organization_id",
					"created_at",
					"updated_at",
					"deleted_at",
					"fractional_index",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pages" WHERE organization_id = $1`)).
					WithArgs(orgId).
					WillReturnRows(rows)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:  "error querying database",
			orgId: orgId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pages" WHERE organization_id = $1`)).
					WithArgs(orgId).
					WillReturnError(fmt.Errorf("database error"))
			},
			wantErr: true,
			wantLen: 0,
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
			pages, err := store.GetPagesByOrganizationId(context.Background(), tt.orgId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pages)
			} else {
				assert.NoError(t, err)
				assert.Len(t, pages, tt.wantLen)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
