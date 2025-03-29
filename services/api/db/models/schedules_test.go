package models

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestScheduleBeforeCreate(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	userId := uuid.New()

	tests := []struct {
		name      string
		schedule  *Schedule
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - user has organization access",
			schedule: &Schedule{
				OrganizationID: orgId,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId.String(), userId.String(), 1).
					WillReturnRows(sqlmock.NewRows([]string{"resource_audience_type", "resource_audience_id", "privilege", "resource_type", "resource_id", "created_at", "updated_at", "deleted_at"}).
						AddRow("user", userId.String(), "read", "organization", orgId.String(), time.Now(), time.Now(), nil))
			},
			wantErr: false,
		},
		{
			name: "error - no user in context",
			schedule: &Schedule{
				OrganizationID: orgId,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
			errMsg:    "no user id found in context",
		},
		{
			name: "error - no organization access",
			schedule: &Schedule{
				OrganizationID: orgId,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND deleted_at IS NULL LIMIT $4`)).
					WithArgs("organization", orgId.String(), userId.String(), 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			wantErr: true,
			errMsg:  "organization access forbidden",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock db
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			dialector := postgres.New(postgres.Config{
				Conn:       mockDB,
				DriverName: "postgres",
			})

			db, err := gorm.Open(dialector, &gorm.Config{})
			assert.NoError(t, err)

			// Setup mock expectations
			tt.mockSetup(mock)

			// Setup context
			ctx := context.Background()
			if tt.name != "error - no user in context" {
				ctx = apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{orgId})
			}
			db = db.WithContext(ctx)

			// Execute test
			err = tt.schedule.BeforeCreate(db)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
