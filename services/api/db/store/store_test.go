package store

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*pgclient.PostgresClient, sqlmock.Sqlmock) {
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

	return &pgclient.PostgresClient{DB: db}, mock
}

func TestNewStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*testing.T) (*pgclient.PostgresClient, sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T) (*pgclient.PostgresClient, sqlmock.Sqlmock) {
				return setupMockDB(t)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, mock := tt.setup(t)
			store, cleanup := NewStore(client)

			assert.NotNil(t, store)
			assert.NotNil(t, cleanup)
			assert.Implements(t, (*Store)(nil), store)
			assert.Implements(t, (*UserStore)(nil), store)
			assert.Implements(t, (*OrganizationStore)(nil), store)
			assert.Implements(t, (*DatasetStore)(nil), store)
			assert.Implements(t, (*TransactionStore)(nil), store)

			mock.ExpectClose()
			cleanup()
		})
	}
}

func TestAppStore_WithTx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		txFunc    func(Store) error
		wantErr   bool
	}{
		{
			name: "successful transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			txFunc: func(s Store) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "transaction rollback on error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			txFunc: func(s Store) error {
				return errors.New("transaction error")
			},
			wantErr: true,
		},
		{
			name: "nested transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			txFunc: func(s Store) error {
				return s.WithTx(context.Background(), func(nestedStore Store) error {
					return nil
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			client, mock := setupMockDB(t)
			store, cleanup := NewStore(client)
			defer cleanup()

			tt.setupMock(mock)

			// Execute
			err := store.WithTx(context.Background(), tt.txFunc)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppStore_WithTx_Integration(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		txFunc    func(Store) error
		wantErr   bool
	}{
		{
			name: "successful complex transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Transaction begin
				mock.ExpectBegin()

				// User query
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits"`).
					WithArgs(userID.String(), 1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "name"}).
						AddRow(userID, "test@example.com", "Test User"))

				// Organization query
				mock.ExpectQuery(`SELECT \* FROM "organizations"`).
					WithArgs(orgID.String(), 1).
					WillReturnRows(sqlmock.NewRows([]string{"organization_id", "name"}).
						AddRow(orgID, "Test Org"))

				// Commit
				mock.ExpectCommit()
			},
			txFunc: func(s Store) error {
				_, err := s.GetUserById(context.Background(), userID.String())
				if err != nil {
					return err
				}

				_, err = s.GetOrganizationById(context.Background(), orgID.String())
				return err
			},
			wantErr: false,
		},
		{
			name: "rollback on error in transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits"`).
					WithArgs(userID.String(), 1).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			txFunc: func(s Store) error {
				_, err := s.GetUserById(context.Background(), userID.String())
				return err
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			client, mock := setupMockDB(t)
			store, cleanup := NewStore(client)
			defer cleanup()

			tt.setupMock(mock)

			// Execute
			err := store.WithTx(context.Background(), tt.txFunc)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
