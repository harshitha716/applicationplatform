package service

import (
	"context"
	"testing"

	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateDatasetActionConfig(t *testing.T) {
	t.Parallel()

	actionID := uuid.New().String()
	config := map[string]interface{}{
		"key":    "value",
		"number": 42,
	}

	tests := []struct {
		name      string
		actionID  string
		config    map[string]interface{}
		mockSetup func(*mock_store.MockDatasetActionStore)
		wantErr   bool
	}{
		{
			name:     "success",
			actionID: actionID,
			config:   config,
			mockSetup: func(m *mock_store.MockDatasetActionStore) {
				m.On("UpdateDatasetActionConfig", mock.Anything, actionID, config).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "error - store error",
			actionID: actionID,
			config:   config,
			mockSetup: func(m *mock_store.MockDatasetActionStore) {
				m.On("UpdateDatasetActionConfig", mock.Anything, actionID, config).
					Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockDatasetActionStore(t)
			tt.mockSetup(mockStore)

			service := NewDatasetActionService(mockStore)

			err := service.UpdateDatasetActionConfig(context.Background(), tt.actionID, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			mockStore.AssertExpectations(t)
		})
	}
}
