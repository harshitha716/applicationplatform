package service

import (
	"context"
	"errors"
	"testing"
	"time"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"gorm.io/gorm"

	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(t *testing.T) (ConnectorService, *mock_store.MockStore, context.Context) {
	mockStore := mock_store.NewMockStore(t)
	service := NewConnectorService(mockStore)

	// Create a context with logger
	ctx := context.Background()

	return service, mockStore, ctx
}

func TestNewConnectorService(t *testing.T) {
	mockStore := mock_store.NewMockStore(t)
	service := NewConnectorService(mockStore)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
}

func TestConnectorsGetAll(t *testing.T) {
	connectorId := uuid.New()
	expectedConnectors := []dbmodels.ConnectorWithActiveConnectionsCount{
		{
			Connector: dbmodels.Connector{
				ID:          connectorId,
				Name:        "test",
				Description: "test",
				DisplayName: "test",
				LogoURL:     "test",
				Category:    "test",
				Status:      "test",
			},
			ActiveConnectionsCount: 1,
		},
	}

	expectedConnectorsDb := []dbmodels.ConnectorWithActiveConnectionsCount{
		{
			Connector: dbmodels.Connector{
				ID:          connectorId,
				Name:        "test",
				Description: "test",
				DisplayName: "test",
				LogoURL:     "test",
				Category:    "test",
				Status:      "test",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				DeletedAt:   gorm.DeletedAt{},
				IsDeleted:   false,
			},
			ActiveConnectionsCount: 1,
		},
	}
	tests := []struct {
		name               string
		setupMock          func(*mock_store.MockStore)
		expectedConnectors []dbmodels.ConnectorWithActiveConnectionsCount
		expectedErr        error
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetAllConnectors(mock.Anything).
					Return(expectedConnectorsDb, nil)
			},
			expectedConnectors: expectedConnectors,
			expectedErr:        nil,
		},
		{
			name: "store error",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetAllConnectors(mock.Anything).
					Return(nil, errors.New("store error"))
			},
			expectedErr: errors.New("store error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, mockStore, ctx := setupTest(t)

			test.setupMock(mockStore)
			connectors, err := service.ListConnectors(ctx)

			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedConnectors, connectors)
			}

			mockStore.AssertExpectations(t)
		})
	}
}
