package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrganization_GetQueryFilters(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgID := uuid.New()
	orgIDs := []uuid.UUID{orgID}

	tests := []struct {
		name            string
		userId          uuid.UUID
		organizationIDs []uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock)
		wantErr         bool
	}{
		{
			name:            "active membership",
			userId:          userId,
			organizationIDs: orgIDs,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"organization_id"}).AddRow(orgID)
				mock.ExpectQuery(`SELECT \* FROM "organizations" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnRows(rows)
			},
		},
		{
			name:            "no active memberships",
			userId:          userId,
			organizationIDs: []uuid.UUID{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "organizations" WHERE EXISTS \(.*\)`).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"organization_id"}))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupTestDB(t)
			org := &Organization{ID: orgID}
			tt.setupMock(mock)

			query := org.GetQueryFilters(db.Model(&Organization{}), tt.userId, tt.organizationIDs)

			var results []Organization
			err := query.Find(&results).Error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStructImplementsBaseModel_Org(t *testing.T) {
	var _ pgclient.BaseModel = &Organization{}
}

func TestBeforeCreate_Org(t *testing.T) {
	userId := uuid.New()
	tests := []struct {
		name    string
		setup   func(context.Context) context.Context
		org     *Organization
		wantErr string
	}{
		{
			name: "no user id in context",
			setup: func(ctx context.Context) context.Context {
				return ctx
			},
			org:     &Organization{},
			wantErr: "no user id found in context",
		},
		{
			name: "non-admin user - forbidden",
			setup: func(ctx context.Context) context.Context {
				return apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{})
			},
			org:     &Organization{},
			wantErr: "insert forbidden",
		},
		{
			name: "admin user - different owner - allowed",
			setup: func(ctx context.Context) context.Context {
				return apicontext.AddAuthToContext(ctx, "admin", userId, []uuid.UUID{})
			},
			org: &Organization{
				OwnerId: uuid.New(),
			},
			wantErr: "",
		},
		{
			name: "admin user - same owner - allowed",
			setup: func(ctx context.Context) context.Context {
				return apicontext.AddAuthToContext(ctx, "admin", userId, []uuid.UUID{})
			},
			org: &Organization{
				OwnerId: userId,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := setupTestDB(t)
			db = db.WithContext(tt.setup(context.Background()))

			err := tt.org.BeforeCreate(db)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
