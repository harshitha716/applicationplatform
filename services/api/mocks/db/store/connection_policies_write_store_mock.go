// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockconnectionPoliciesWriteStore is an autogenerated mock type for the connectionPoliciesWriteStore type
type MockconnectionPoliciesWriteStore struct {
	mock.Mock
}

type MockconnectionPoliciesWriteStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockconnectionPoliciesWriteStore) EXPECT() *MockconnectionPoliciesWriteStore_Expecter {
	return &MockconnectionPoliciesWriteStore_Expecter{mock: &_m.Mock}
}

// CreateConnectionPolicy provides a mock function with given fields: ctx, connectionId, audienceType, audienceId, privilege
func (_m *MockconnectionPoliciesWriteStore) CreateConnectionPolicy(ctx context.Context, connectionId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, connectionId, audienceType, audienceId, privilege)

	if len(ret) == 0 {
		panic("no return value specified for CreateConnectionPolicy")
	}

	var r0 *models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, connectionId, audienceType, audienceId, privilege)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) *models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, connectionId, audienceType, audienceId, privilege)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) error); ok {
		r1 = rf(ctx, connectionId, audienceType, audienceId, privilege)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateConnectionPolicy'
type MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call struct {
	*mock.Call
}

// CreateConnectionPolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - connectionId uuid.UUID
//   - audienceType models.AudienceType
//   - audienceId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockconnectionPoliciesWriteStore_Expecter) CreateConnectionPolicy(ctx interface{}, connectionId interface{}, audienceType interface{}, audienceId interface{}, privilege interface{}) *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call {
	return &MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call{Call: _e.mock.On("CreateConnectionPolicy", ctx, connectionId, audienceType, audienceId, privilege)}
}

func (_c *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call) Run(run func(ctx context.Context, connectionId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege)) *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.AudienceType), args[3].(uuid.UUID), args[4].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockconnectionPoliciesWriteStore_CreateConnectionPolicy_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateConnectionPolicy provides a mock function with given fields: ctx, connectionId, audienceId, privilege
func (_m *MockconnectionPoliciesWriteStore) UpdateConnectionPolicy(ctx context.Context, connectionId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, connectionId, audienceId, privilege)

	if len(ret) == 0 {
		panic("no return value specified for UpdateConnectionPolicy")
	}

	var r0 *models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, connectionId, audienceId, privilege)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) *models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, connectionId, audienceId, privilege)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) error); ok {
		r1 = rf(ctx, connectionId, audienceId, privilege)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateConnectionPolicy'
type MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call struct {
	*mock.Call
}

// UpdateConnectionPolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - connectionId uuid.UUID
//   - audienceId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockconnectionPoliciesWriteStore_Expecter) UpdateConnectionPolicy(ctx interface{}, connectionId interface{}, audienceId interface{}, privilege interface{}) *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call {
	return &MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call{Call: _e.mock.On("UpdateConnectionPolicy", ctx, connectionId, audienceId, privilege)}
}

func (_c *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call) Run(run func(ctx context.Context, connectionId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege)) *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockconnectionPoliciesWriteStore_UpdateConnectionPolicy_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockconnectionPoliciesWriteStore creates a new instance of MockconnectionPoliciesWriteStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockconnectionPoliciesWriteStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockconnectionPoliciesWriteStore {
	mock := &MockconnectionPoliciesWriteStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
