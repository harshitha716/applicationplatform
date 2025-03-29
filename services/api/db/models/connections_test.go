package models

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConnection_TableName(t *testing.T) {
	t.Parallel()
	connection := Connection{}
	assert.Equal(t, "connections", connection.TableName())
}

func TestConnection_GetAccessControlFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	org1ID := uuid.New()
	orgIDs := []uuid.UUID{org1ID}

	db, mock := setupTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "connections" WHERE EXISTS ( SELECT 1 FROM "app"."flattened_resource_audience_policies" frap WHERE frap.resource_type = 'connection' AND frap.resource_id = connections.id AND frap.user_id = $1 AND frap.deleted_at IS NULL )`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

	connection := Connection{}
	db.Model(&Connection{}).Where(connection.GetQueryFilters(db, userId, orgIDs)).Scan(&connection)

	assert.NoError(t, mock.ExpectationsWereMet())
}
