package service

import (
	"context"
	"testing"

	"github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateConnection(t *testing.T) {
	type testCase struct {
		name          string
		setupContext  func() context.Context
		params        models.CreateConnectionParams
		setupMock     func(mockStore *mock_store.MockStore)
		expectedID    uuid.UUID
		expectedError error
	}

	connectorId := uuid.New()
	connectionId := uuid.New()
	userId := uuid.New()
	policyId := uuid.New()

	tests := []testCase{
		{
			name: "successful connection creation with policy",
			setupContext: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()})
			},
			params: models.CreateConnectionParams{
				ConnectorID: connectorId,
				Name:        "Test Connection",
				Status:      "active",
			},
			setupMock: func(mockStore *mock_store.MockStore) {

				// Setup CreateConnection mock
				mockStore.On("CreateConnection", mock.Anything, mock.MatchedBy(func(params *models.CreateConnectionParams) bool {
					return params.ConnectorID == connectorId &&
						params.Name == "Test Connection" &&
						params.Status == "active"
				})).Return(connectionId, nil)

				// Setup CreateConnectionPolicy mock
				mockStore.On("CreateConnectionPolicy",
					mock.Anything,
					connectionId,
					models.AudienceTypeUser,
					userId,
					models.PrivilegeConnectionAdmin,
				).Return(&models.ResourceAudiencePolicy{
					ID:                   policyId,
					ResourceType:         models.ResourceTypeConnection,
					ResourceID:           connectionId,
					ResourceAudienceType: models.AudienceTypeUser,
					ResourceAudienceID:   userId,
					Privilege:            models.PrivilegeConnectionAdmin,
				}, nil)
			},
			expectedID:    connectionId,
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock store
			mockStore := mock_store.NewMockStore(t)
			tc.setupMock(mockStore)

			// Create service with mock store
			service := NewConnectionService(mockStore)

			// Execute test with context
			ctx := tc.setupContext()
			id, err := service.CreateConnection(ctx, tc.params, nil)

			// Assert results
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedID, id)

			// Verify all mocks were called as expected
			mockStore.AssertExpectations(t)
		})
	}
}
