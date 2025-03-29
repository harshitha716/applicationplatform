package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUserById(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	email := "test@example.com"
	name := "Test User"

	tests := []struct {
		name      string
		userID    string
		mockSetup func(sqlmock.Sqlmock)
		want      *models.User
		wantErr   bool
	}{
		{
			name:   "success",
			userID: userID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"user_id",
					"email",
					"name",
				}).AddRow(
					userID,
					email,
					name,
				)

				mock.ExpectQuery(`SELECT \* FROM "users_with_traits" WHERE user_id = \$1 ORDER BY "users_with_traits"."user_id" LIMIT \$2`).
					WithArgs(userID.String(), 1).
					WillReturnRows(rows)
			},
			want: &models.User{
				ID:    userID,
				Email: email,
				Name:  name,
			},
		},
		{
			name:   "not found",
			userID: uuid.New().String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits" WHERE user_id = \$1 ORDER BY "users_with_traits"."user_id" LIMIT \$2`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:   "invalid uuid",
			userID: "invalid-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits" WHERE user_id = \$1 ORDER BY "users_with_traits"."user_id" LIMIT \$2`).
					WithArgs("invalid-uuid", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:   "database error",
			userID: userID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits" WHERE user_id = \$1 ORDER BY "users_with_traits"."user_id" LIMIT \$2`).
					WithArgs(userID.String(), 1).
					WillReturnError(gorm.ErrInvalidDB)
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
			got, err := store.GetUserById(context.Background(), tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUsersAll(t *testing.T) {
	t.Parallel()

	user1ID := uuid.New()
	user2ID := uuid.New()
	email1 := "test1@example.com"
	email2 := "test2@example.com"
	name1 := "Test User 1"
	name2 := "Test User 2"

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []models.User
		wantErr   bool
	}{
		{
			name: "success - multiple users",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"user_id",
					"email",
					"name",
				}).AddRow(
					user1ID,
					email1,
					name1,
				).AddRow(
					user2ID,
					email2,
					name2,
				)

				mock.ExpectQuery(`SELECT \* FROM "users_with_traits"`).
					WillReturnRows(rows)
			},
			want: []models.User{
				{
					ID:    user1ID,
					Email: email1,
					Name:  name1,
				},
				{
					ID:    user2ID,
					Email: email2,
					Name:  name2,
				},
			},
		},
		{
			name: "success - empty result",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits"`).
					WillReturnRows(sqlmock.NewRows([]string{
						"user_id",
						"email",
						"name",
					}))
			},
			want: []models.User{},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits"`).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantErr: true,
		},
		{
			name: "context canceled",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users_with_traits"`).
					WillReturnError(context.Canceled)
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
			got, err := store.GetUsersAll(context.Background())

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))

			for i := range got {
				assert.Equal(t, tt.want[i].ID, got[i].ID)
				assert.Equal(t, tt.want[i].Email, got[i].Email)
				assert.Equal(t, tt.want[i].Name, got[i].Name)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
