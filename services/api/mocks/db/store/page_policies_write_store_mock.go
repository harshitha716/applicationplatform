// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockpagePoliciesWriteStore is an autogenerated mock type for the pagePoliciesWriteStore type
type MockpagePoliciesWriteStore struct {
	mock.Mock
}

type MockpagePoliciesWriteStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockpagePoliciesWriteStore) EXPECT() *MockpagePoliciesWriteStore_Expecter {
	return &MockpagePoliciesWriteStore_Expecter{mock: &_m.Mock}
}

// CreatePagePolicy provides a mock function with given fields: ctx, pageId, audienceType, audienceId, privilege
func (_m *MockpagePoliciesWriteStore) CreatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, pageId, audienceType, audienceId, privilege)

	if len(ret) == 0 {
		panic("no return value specified for CreatePagePolicy")
	}

	var r0 *models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, pageId, audienceType, audienceId, privilege)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) *models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, pageId, audienceType, audienceId, privilege)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) error); ok {
		r1 = rf(ctx, pageId, audienceType, audienceId, privilege)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockpagePoliciesWriteStore_CreatePagePolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePagePolicy'
type MockpagePoliciesWriteStore_CreatePagePolicy_Call struct {
	*mock.Call
}

// CreatePagePolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - audienceType models.AudienceType
//   - audienceId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockpagePoliciesWriteStore_Expecter) CreatePagePolicy(ctx interface{}, pageId interface{}, audienceType interface{}, audienceId interface{}, privilege interface{}) *MockpagePoliciesWriteStore_CreatePagePolicy_Call {
	return &MockpagePoliciesWriteStore_CreatePagePolicy_Call{Call: _e.mock.On("CreatePagePolicy", ctx, pageId, audienceType, audienceId, privilege)}
}

func (_c *MockpagePoliciesWriteStore_CreatePagePolicy_Call) Run(run func(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege)) *MockpagePoliciesWriteStore_CreatePagePolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.AudienceType), args[3].(uuid.UUID), args[4].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockpagePoliciesWriteStore_CreatePagePolicy_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockpagePoliciesWriteStore_CreatePagePolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockpagePoliciesWriteStore_CreatePagePolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockpagePoliciesWriteStore_CreatePagePolicy_Call {
	_c.Call.Return(run)
	return _c
}

// DeletePagePolicy provides a mock function with given fields: ctx, pageId, audienceType, audienceId
func (_m *MockpagePoliciesWriteStore) DeletePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {
	ret := _m.Called(ctx, pageId, audienceType, audienceId)

	if len(ret) == 0 {
		panic("no return value specified for DeletePagePolicy")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID) error); ok {
		r0 = rf(ctx, pageId, audienceType, audienceId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockpagePoliciesWriteStore_DeletePagePolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeletePagePolicy'
type MockpagePoliciesWriteStore_DeletePagePolicy_Call struct {
	*mock.Call
}

// DeletePagePolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - audienceType models.AudienceType
//   - audienceId uuid.UUID
func (_e *MockpagePoliciesWriteStore_Expecter) DeletePagePolicy(ctx interface{}, pageId interface{}, audienceType interface{}, audienceId interface{}) *MockpagePoliciesWriteStore_DeletePagePolicy_Call {
	return &MockpagePoliciesWriteStore_DeletePagePolicy_Call{Call: _e.mock.On("DeletePagePolicy", ctx, pageId, audienceType, audienceId)}
}

func (_c *MockpagePoliciesWriteStore_DeletePagePolicy_Call) Run(run func(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID)) *MockpagePoliciesWriteStore_DeletePagePolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.AudienceType), args[3].(uuid.UUID))
	})
	return _c
}

func (_c *MockpagePoliciesWriteStore_DeletePagePolicy_Call) Return(_a0 error) *MockpagePoliciesWriteStore_DeletePagePolicy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockpagePoliciesWriteStore_DeletePagePolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID) error) *MockpagePoliciesWriteStore_DeletePagePolicy_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePagePolicy provides a mock function with given fields: ctx, pageId, audienceId, privilege
func (_m *MockpagePoliciesWriteStore) UpdatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, pageId, audienceId, privilege)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePagePolicy")
	}

	var r0 *models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, pageId, audienceId, privilege)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) *models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, pageId, audienceId, privilege)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) error); ok {
		r1 = rf(ctx, pageId, audienceId, privilege)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockpagePoliciesWriteStore_UpdatePagePolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePagePolicy'
type MockpagePoliciesWriteStore_UpdatePagePolicy_Call struct {
	*mock.Call
}

// UpdatePagePolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - audienceId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockpagePoliciesWriteStore_Expecter) UpdatePagePolicy(ctx interface{}, pageId interface{}, audienceId interface{}, privilege interface{}) *MockpagePoliciesWriteStore_UpdatePagePolicy_Call {
	return &MockpagePoliciesWriteStore_UpdatePagePolicy_Call{Call: _e.mock.On("UpdatePagePolicy", ctx, pageId, audienceId, privilege)}
}

func (_c *MockpagePoliciesWriteStore_UpdatePagePolicy_Call) Run(run func(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege)) *MockpagePoliciesWriteStore_UpdatePagePolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockpagePoliciesWriteStore_UpdatePagePolicy_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockpagePoliciesWriteStore_UpdatePagePolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockpagePoliciesWriteStore_UpdatePagePolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockpagePoliciesWriteStore_UpdatePagePolicy_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockpagePoliciesWriteStore creates a new instance of MockpagePoliciesWriteStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockpagePoliciesWriteStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockpagePoliciesWriteStore {
	mock := &MockpagePoliciesWriteStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
