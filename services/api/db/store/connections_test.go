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
)

func TestCreateConnection(t *testing.T) {
	t.Parallel()

	connectorId := uuid.New()
	orgId := uuid.New()
	userId := uuid.New()

	tests := []struct {
		name         string
		setupContext func() context.Context
		params       models.CreateConnectionParams
		mockSetup    func(sqlmock.Sqlmock)
		wantErr      bool
	}{
		{
			name: "success",
			setupContext: func() context.Context {
				ctx := context.Background()
				// Add both user ID and org ID to context
				return apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{orgId})
			},
			params: models.CreateConnectionParams{
				ConnectorID: connectorId,
				Name:        "test",
				Status:      "active",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId.String(), userId.String(), 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_type", "resource_audience_id", "privilege", "resource_type", "resource_id", "created_at", "updated_at", "deleted_at"}).AddRow("user", userId.String(), "viewer", "organization", orgId.String(), time.Now(), time.Now(), nil))
				mock.ExpectExec(`INSERT INTO "connections"`).
					WithArgs(sqlmock.AnyArg(), connectorId.String(), orgId.String(), "test", "active", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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

			ctx := context.Background()

			ctx = apicontext.AddAuthToContext(ctx, "user_id", userId, []uuid.UUID{orgId})

			connectionId, err := store.CreateConnection(ctx, &tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "no organization in context" {
					assert.Equal(t, "organization access forbidden", err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, connectionId)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetConnections(t *testing.T) {
	t.Parallel()

	connectorId := uuid.New()
	orgId := uuid.New()
	userId := uuid.New()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful retrieval",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "connector_id", "organization_id", "name", "status", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), connectorId, orgId, "test conn", "active", time.Now(), time.Now(), nil)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "connections" WHERE organization_id = $1 AND "connections"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(orgId.String()).
					WillReturnRows(rows)

				// For Preload("Connector")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "connectors" WHERE "connectors"."id" = $1 AND "connectors"."deleted_at" IS NULL`)).
					WithArgs(connectorId.String()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(connectorId))

				// For Preload("Schedules")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE "schedules"."connection_id" = $1`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
		{
			name: "no organization in context",
			mockSetup: func(mock sqlmock.Sqlmock) {
			},
			wantErr: true,
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "connections" WHERE organization_id = $1 AND "connections"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(orgId.String()).
					WillReturnError(fmt.Errorf("database error"))
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

			ctx := context.Background()
			if tt.name != "no organization in context" {
				ctx = apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{orgId})
			}

			connections, err := store.GetConnections(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "no organization in context" {
					assert.Equal(t, "organization access forbidden", err.Error())
				}
				assert.Nil(t, connections)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, connections)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
