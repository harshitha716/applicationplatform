package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestGetConnectors(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	Id := uuid.New()
	Name := "Test Connector"
	Description := "Test Description"
	DisplayName := "Test Display Name"
	Documentation := "Test Documentation"
	LogoURL := "Test Logo URL"
	ConfigTemplate := json.RawMessage(`{"key": "value"}`)
	Category := "Test Category"
	CreatedAt := time.Now()
	UpdatedAt := time.Now()
	ActiveConnectionsCount := int64(2)

	tests := []struct {
		name      string
		orgIds    []uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      []models.ConnectorWithActiveConnectionsCount
		wantErr   bool
	}{
		{
			name:   "success",
			orgIds: []uuid.UUID{orgId},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "description", "display_name", "documentation", "logo_url", "config_template", "category", "created_at", "updated_at", "deleted_at", "is_deleted", "status", "active_connections_count",
				}).AddRow(
					Id, Name, Description, DisplayName, Documentation, LogoURL, ConfigTemplate, Category, CreatedAt, UpdatedAt, nil, false, "active", ActiveConnectionsCount,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT connectors.*, COUNT(connections.id) AS active_connections_count FROM "connectors" LEFT JOIN connections ON connectors.id = connections.connector_id AND connections.status = 'active' AND connections.organization_id = $1 WHERE "connectors"."deleted_at" IS NULL GROUP BY "connectors"."id" ORDER BY connectors.name ASC`)).
					WithArgs(orgId).
					WillReturnRows(rows)
			},
			want: []models.ConnectorWithActiveConnectionsCount{
				{
					Connector: models.Connector{
						ID:             Id,
						Name:           Name,
						Description:    Description,
						DisplayName:    DisplayName,
						Documentation:  Documentation,
						LogoURL:        LogoURL,
						ConfigTemplate: ConfigTemplate,
						Category:       Category,
						CreatedAt:      CreatedAt,
						UpdatedAt:      UpdatedAt,
						DeletedAt:      gorm.DeletedAt{},
						IsDeleted:      false,
						Status:         "active",
					},
					ActiveConnectionsCount: int(ActiveConnectionsCount),
				},
			},
		},
		{
			name:      "error - no org id",
			orgIds:    []uuid.UUID{},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "error - multiple org ids",
			orgIds:    []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)

			tt.mockSetup(mock)

			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "user", uuid.New(), tt.orgIds)
			store := &appStore{client: &pgclient.PostgresClient{DB: gormDB}}

			got, err := store.GetAllConnectors(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("appStore.GetAllConnectors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appStore.GetAllConnectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConnectorById(t *testing.T) {
	t.Parallel()

	Id := uuid.New()
	Name := "Test Connector"
	Description := "Test Description"
	DisplayName := "Test Display Name"
	Documentation := "Test Documentation"
	LogoURL := "Test Logo URL"
	ConfigTemplate := json.RawMessage(`{"key": "value"}`)
	Category := "Test Category"
	CreatedAt := time.Now()
	UpdatedAt := time.Now()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      *models.Connector
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "description", "display_name", "documentation", "logo_url", "config_template", "category", "created_at", "updated_at", "deleted_at", "is_deleted", "status",
				}).AddRow(
					Id, Name, Description, DisplayName, Documentation, LogoURL, ConfigTemplate, Category, CreatedAt, UpdatedAt, nil, false, "active",
				)

				mock.ExpectQuery(`SELECT \* FROM "connectors" WHERE id = \$1`).
					WithArgs(Id).
					WillReturnRows(rows)
			},
			want: &models.Connector{
				ID:             Id,
				Name:           Name,
				Description:    Description,
				DisplayName:    DisplayName,
				Documentation:  Documentation,
				LogoURL:        LogoURL,
				ConfigTemplate: ConfigTemplate,
				Category:       Category,
				CreatedAt:      CreatedAt,
				UpdatedAt:      UpdatedAt,
				DeletedAt:      gorm.DeletedAt{},
				IsDeleted:      false,
				Status:         "active",
			},
		},
		{
			name: "not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "connectors" WHERE id = \$1`).
					WithArgs(Id).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})
	}
}
