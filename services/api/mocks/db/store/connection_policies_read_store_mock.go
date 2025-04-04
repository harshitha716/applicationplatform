// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockconnectionPoliciesReadStore is an autogenerated mock type for the connectionPoliciesReadStore type
type MockconnectionPoliciesReadStore struct {
	mock.Mock
}

type MockconnectionPoliciesReadStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockconnectionPoliciesReadStore) EXPECT() *MockconnectionPoliciesReadStore_Expecter {
	return &MockconnectionPoliciesReadStore_Expecter{mock: &_m.Mock}
}

// GetConnectionPolicies provides a mock function with given fields: ctx, connectionId
func (_m *MockconnectionPoliciesReadStore) GetConnectionPolicies(ctx context.Context, connectionId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, connectionId)

	if len(ret) == 0 {
		panic("no return value specified for GetConnectionPolicies")
	}

	var r0 []models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, connectionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, connectionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, connectionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockconnectionPoliciesReadStore_GetConnectionPolicies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnectionPolicies'
type MockconnectionPoliciesReadStore_GetConnectionPolicies_Call struct {
	*mock.Call
}

// GetConnectionPolicies is a helper method to define mock.On call
//   - ctx context.Context
//   - connectionId uuid.UUID
func (_e *MockconnectionPoliciesReadStore_Expecter) GetConnectionPolicies(ctx interface{}, connectionId interface{}) *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call {
	return &MockconnectionPoliciesReadStore_GetConnectionPolicies_Call{Call: _e.mock.On("GetConnectionPolicies", ctx, connectionId)}
}

func (_c *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call) Run(run func(ctx context.Context, connectionId uuid.UUID)) *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call) Return(_a0 []models.ResourceAudiencePolicy, _a1 error) *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)) *MockconnectionPoliciesReadStore_GetConnectionPolicies_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockconnectionPoliciesReadStore creates a new instance of MockconnectionPoliciesReadStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockconnectionPoliciesReadStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockconnectionPoliciesReadStore {
	mock := &MockconnectionPoliciesReadStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
