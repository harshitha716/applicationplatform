package store

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetFlattenedResourceAudiencePolicies(t *testing.T) {

	userId1 := uuid.New()
	resourceId := uuid.New()

	tests := []struct {
		name           string
		filters        models.FlattenedResourceAudiencePoliciesFilters
		mockSetup      func(mock sqlmock.Sqlmock)
		expectedResult []models.FlattenedResourceAudiencePolicy
		expectedError  error
	}{
		{
			name: "should return policies with all filters",
			filters: models.FlattenedResourceAudiencePoliciesFilters{
				ResourceIds:   []uuid.UUID{resourceId},
				UserIds:       []uuid.UUID{userId1},
				ResourceTypes: []string{"type-1"},
				Privileges:    []models.ResourcePrivilege{models.PrivilegeDatasetViewer},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_id IN ($1) AND user_id IN ($2) AND resource_type IN ($3) AND privilege IN ($4)`)).
					WithArgs(resourceId, userId1, "type-1", models.PrivilegeDatasetViewer).
					WillReturnRows(sqlmock.NewRows([]string{
						"resource_audience_type",
						"user_id",
						"resource_id",
						"resource_type",
						"privilege",
						"created_at",
						"updated_at",
						"deleted_at",
					}).AddRow(
						"type-1",
						userId1.String(),
						resourceId.String(),
						"type-1",
						models.PrivilegeDatasetViewer,
						"2024-01-20",
						"2024-01-20",
						nil,
					))
			},
			expectedResult: []models.FlattenedResourceAudiencePolicy{
				{
					ResourceAudienceType: "type-1",
					UserId:               userId1,
					ResourceId:           resourceId,
					ResourceType:         "type-1",
					Privilege:            models.PrivilegeDatasetViewer,
					CreatedAt:            "2024-01-20",
					UpdatedAt:            "2024-01-20",
					DeletedAt:            "",
				},
			},
			expectedError: nil,
		},
		{
			name:    "should return all policies with no filters",
			filters: models.FlattenedResourceAudiencePoliciesFilters{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies"`)).
					WillReturnRows(sqlmock.NewRows([]string{
						"resource_audience_type",
						"user_id",
						"resource_id",
						"resource_type",
						"privilege",
						"created_at",
						"updated_at",
						"deleted_at",
					}).AddRow(
						"type-1",
						userId1,
						resourceId,
						"type-1",
						models.PrivilegeDatasetViewer,
						"2024-01-20",
						"2024-01-20",
						nil,
					))
			},
			expectedResult: []models.FlattenedResourceAudiencePolicy{
				{
					ResourceAudienceType: "type-1",
					UserId:               userId1,
					ResourceId:           resourceId,
					ResourceType:         "type-1",
					Privilege:            models.PrivilegeDatasetViewer,
					CreatedAt:            "2024-01-20",
					UpdatedAt:            "2024-01-20",
					DeletedAt:            "",
				},
			},
			expectedError: nil,
		},
		{
			name: "should handle partial filters",
			filters: models.FlattenedResourceAudiencePoliciesFilters{
				ResourceIds: []uuid.UUID{resourceId},
				UserIds:     []uuid.UUID{userId1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_id IN ($1) AND user_id IN ($2)`)).
					WithArgs(resourceId, userId1).
					WillReturnRows(sqlmock.NewRows([]string{
						"resource_audience_type",
						"user_id",
						"resource_id",
						"resource_type",
						"privilege",
						"created_at",
						"updated_at",
						"deleted_at",
					}))
			},
			expectedResult: []models.FlattenedResourceAudiencePolicy{},
			expectedError:  nil,
		},
		{
			name: "should handle database error",
			filters: models.FlattenedResourceAudiencePoliciesFilters{
				ResourceIds: []uuid.UUID{resourceId},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_id IN ($1)`)).
					WithArgs(resourceId).
					WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test database and mock
			db, mock := getMockDB(t)

			store := appStore{
				client: &pgclient.PostgresClient{DB: db},
			}

			// Setup mock expectations
			tt.mockSetup(mock)

			// Execute the function
			result, err := store.GetFlattenedResourceAudiencePolicies(context.Background(), tt.filters)

			// Verify results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
