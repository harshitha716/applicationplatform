package store

import (
	"context"
	"encoding/json"
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

func getMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create gorm db: %v", err)
	}
	return gormdb, mock
}

func TestGetWidgetInstanceByID(t *testing.T) {
	t.Parallel()

	instanceID := uuid.New()
	sheetID := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name           string
		instanceID     uuid.UUID
		mockSetup      func(sqlmock.Sqlmock)
		wantErr        bool
		expectedFields models.WidgetInstance
	}{
		{
			name:       "success",
			instanceID: instanceID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"widget_instance_id",
					"widget_type",
					"sheet_id",
					"title",
					"data_mappings",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					instanceID,
					"bar_chart",
					sheetID,
					"Test Widget Instance",
					[]byte(`{"key": "value"}`),
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "widget_instances"`).
					WithArgs(instanceID, 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedFields: models.WidgetInstance{
				ID:           instanceID,
				WidgetType:   "bar_chart",
				SheetID:      sheetID,
				Title:        "Test Widget Instance",
				DataMappings: json.RawMessage(`{"key": "value"}`),
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name:       "not found",
			instanceID: uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "widget_instances"`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
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

			instance, err := store.GetWidgetInstanceByID(context.Background(), tt.instanceID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, models.WidgetInstance{}, instance)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFields, instance)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetWidgetTemplate(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	tests := []struct {
		name           string
		widgetType     string
		mockSetup      func(sqlmock.Sqlmock)
		wantErr        bool
		expectedFields models.Widget
	}{
		{
			name:       "success",
			widgetType: "chart",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"name",
					"type",
					"template_schema",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					"Test Widget",
					"chart",
					[]byte(`{"type": "bar"}`),
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "widgets" WHERE type = \$1 ORDER BY "widgets"."type" LIMIT \$2`).
					WithArgs("chart", 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedFields: models.Widget{
				Name: "Test Widget",
				Type: "chart",
			},
		},
		{
			name:       "not found",
			widgetType: "chart",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "widgets" WHERE type = \$1 ORDER BY "widgets"."type" LIMIT \$2`).
					WithArgs("chart", 1).
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
			widget, err := store.GetWidgetTemplate(context.Background(), tt.widgetType)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, models.Widget{}, widget)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFields.Name, widget.Name)
			assert.Equal(t, tt.expectedFields.Type, widget.Type)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_CreateWidget(t *testing.T) {
	t.Parallel()

	sheetID := uuid.New()
	userID := uuid.New()
	instanceID := uuid.New()
	rawMsg := json.RawMessage(`{}`)
	ptr := &rawMsg

	tests := []struct {
		name      string
		widget    models.WidgetInstance
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful widget creation",
			widget: models.WidgetInstance{
				ID:            instanceID,
				SheetID:       sheetID,
				Title:         "Test Widget",
				WidgetType:    "bar_chart",
				DataMappings:  json.RawMessage(`{"key": "value"}`),
				DisplayConfig: ptr,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE EXISTS \(.*\) LIMIT \$3`).
					WithArgs(sheetID, userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"sheet_id"}).AddRow(sheetID))

				mock.ExpectExec(`INSERT INTO "widget_instances"`).
					WithArgs(
						sqlmock.AnyArg(),           // widget_instance_id
						"bar_chart",                // widget_type
						sheetID,                    // sheet_id
						"Test Widget",              // title
						[]byte(`{"key": "value"}`), // data_mappings
						sqlmock.AnyArg(),           // created_at
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						[]byte(`{}`),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "permission denied",
			widget: models.WidgetInstance{
				SheetID:      sheetID,
				Title:        "Test Widget",
				WidgetType:   "bar_chart",
				DataMappings: json.RawMessage(`{"key": "value"}`),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE EXISTS \(.*\) LIMIT \$3`).
					WithArgs(sheetID, userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "database error",
			widget: models.WidgetInstance{
				SheetID:      sheetID,
				WidgetType:   "bar_chart",
				Title:        "Test Widget",
				DataMappings: json.RawMessage(`{"key": "value"}`),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE EXISTS \(.*\) LIMIT \$3`).
					WithArgs(sheetID, userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"sheet_id"}).AddRow(sheetID))

				mock.ExpectExec(`INSERT INTO "widget_instances"`).
					WithArgs(
						sqlmock.AnyArg(),
						"bar_chart",
						sheetID,
						"Test Widget",
						[]byte(`{"key": "value"}`),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
					).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
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

			ctx := apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{})

			result, err := store.CreateWidgetInstance(ctx, &tt.widget)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, result.ID)
			assert.Equal(t, tt.widget.SheetID, result.SheetID)
			assert.Equal(t, tt.widget.WidgetType, result.WidgetType)
			assert.Equal(t, tt.widget.Title, result.Title)
			assert.Equal(t, tt.widget.DataMappings, result.DataMappings)
			assert.NotZero(t, result.CreatedAt)
			assert.NotZero(t, result.UpdatedAt)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_UpdateWidgetInstance(t *testing.T) {
	t.Parallel()

	instanceID := uuid.New()
	sheetID := uuid.New()
	now := time.Now().UTC()

	rawMsg := json.RawMessage(`{}`)
	ptr := &rawMsg

	tests := []struct {
		name      string
		widget    models.WidgetInstance
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful widget update",
			widget: models.WidgetInstance{
				ID:            instanceID,
				SheetID:       sheetID,
				Title:         "Updated Widget",
				WidgetType:    "bar_chart",
				DataMappings:  json.RawMessage(`{"key": "updated_value"}`),
				DisplayConfig: ptr,
				UpdatedAt:     now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "widget_instances" SET`).
					WithArgs(
						"bar_chart",
						sheetID,
						"Updated Widget",
						[]byte(`{"key": "updated_value"}`),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
						[]byte(`{}`),
						instanceID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			widget: models.WidgetInstance{
				ID:            instanceID,
				SheetID:       sheetID,
				Title:         "Updated Widget",
				WidgetType:    "bar_chart",
				DataMappings:  json.RawMessage(`{"key": "updated_value"}`),
				DisplayConfig: ptr,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "widget_instances" SET`).
					WithArgs(
						"bar_chart",
						sheetID,
						"Updated Widget",
						[]byte(`{"key": "updated_value"}`),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
						[]byte(`{}`),
						instanceID,
					).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
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

			result, err := store.UpdateWidgetInstance(context.Background(), &tt.widget)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.widget.ID, result.ID)
			assert.Equal(t, tt.widget.SheetID, result.SheetID)
			assert.Equal(t, tt.widget.WidgetType, result.WidgetType)
			assert.Equal(t, tt.widget.Title, result.Title)
			assert.Equal(t, tt.widget.DataMappings, result.DataMappings)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
