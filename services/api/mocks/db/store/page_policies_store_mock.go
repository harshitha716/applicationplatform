// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockpagePoliciesStore is an autogenerated mock type for the pagePoliciesStore type
type MockpagePoliciesStore struct {
	mock.Mock
}

type MockpagePoliciesStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockpagePoliciesStore) EXPECT() *MockpagePoliciesStore_Expecter {
	return &MockpagePoliciesStore_Expecter{mock: &_m.Mock}
}

// CreatePagePolicy provides a mock function with given fields: ctx, pageId, audienceType, audienceId, privilege
func (_m *MockpagePoliciesStore) CreatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
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

// MockpagePoliciesStore_CreatePagePolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePagePolicy'
type MockpagePoliciesStore_CreatePagePolicy_Call struct {
	*mock.Call
}

// CreatePagePolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - audienceType models.AudienceType
//   - audienceId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockpagePoliciesStore_Expecter) CreatePagePolicy(ctx interface{}, pageId interface{}, audienceType interface{}, audienceId interface{}, privilege interface{}) *MockpagePoliciesStore_CreatePagePolicy_Call {
	return &MockpagePoliciesStore_CreatePagePolicy_Call{Call: _e.mock.On("CreatePagePolicy", ctx, pageId, audienceType, audienceId, privilege)}
}

func (_c *MockpagePoliciesStore_CreatePagePolicy_Call) Run(run func(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID, privilege models.ResourcePrivilege)) *MockpagePoliciesStore_CreatePagePolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.AudienceType), args[3].(uuid.UUID), args[4].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockpagePoliciesStore_CreatePagePolicy_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockpagePoliciesStore_CreatePagePolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockpagePoliciesStore_CreatePagePolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockpagePoliciesStore_CreatePagePolicy_Call {
	_c.Call.Return(run)
	return _c
}

// DeletePagePolicy provides a mock function with given fields: ctx, pageId, audienceType, audienceId
func (_m *MockpagePoliciesStore) DeletePagePolicy(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {
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

// MockpagePoliciesStore_DeletePagePolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeletePagePolicy'
type MockpagePoliciesStore_DeletePagePolicy_Call struct {
	*mock.Call
}

// DeletePagePolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - audienceType models.AudienceType
//   - audienceId uuid.UUID
func (_e *MockpagePoliciesStore_Expecter) DeletePagePolicy(ctx interface{}, pageId interface{}, audienceType interface{}, audienceId interface{}) *MockpagePoliciesStore_DeletePagePolicy_Call {
	return &MockpagePoliciesStore_DeletePagePolicy_Call{Call: _e.mock.On("DeletePagePolicy", ctx, pageId, audienceType, audienceId)}
}

func (_c *MockpagePoliciesStore_DeletePagePolicy_Call) Run(run func(ctx context.Context, pageId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID)) *MockpagePoliciesStore_DeletePagePolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.AudienceType), args[3].(uuid.UUID))
	})
	return _c
}

func (_c *MockpagePoliciesStore_DeletePagePolicy_Call) Return(_a0 error) *MockpagePoliciesStore_DeletePagePolicy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockpagePoliciesStore_DeletePagePolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID) error) *MockpagePoliciesStore_DeletePagePolicy_Call {
	_c.Call.Return(run)
	return _c
}

// GetPagePoliciesByEmail provides a mock function with given fields: ctx, pageId, email
func (_m *MockpagePoliciesStore) GetPagePoliciesByEmail(ctx context.Context, pageId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, pageId, email)

	if len(ret) == 0 {
		panic("no return value specified for GetPagePoliciesByEmail")
	}

	var r0 []models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) ([]models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, pageId, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) []models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, pageId, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, pageId, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockpagePoliciesStore_GetPagePoliciesByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPagePoliciesByEmail'
type MockpagePoliciesStore_GetPagePoliciesByEmail_Call struct {
	*mock.Call
}

// GetPagePoliciesByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - email string
func (_e *MockpagePoliciesStore_Expecter) GetPagePoliciesByEmail(ctx interface{}, pageId interface{}, email interface{}) *MockpagePoliciesStore_GetPagePoliciesByEmail_Call {
	return &MockpagePoliciesStore_GetPagePoliciesByEmail_Call{Call: _e.mock.On("GetPagePoliciesByEmail", ctx, pageId, email)}
}

func (_c *MockpagePoliciesStore_GetPagePoliciesByEmail_Call) Run(run func(ctx context.Context, pageId uuid.UUID, email string)) *MockpagePoliciesStore_GetPagePoliciesByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string))
	})
	return _c
}

func (_c *MockpagePoliciesStore_GetPagePoliciesByEmail_Call) Return(_a0 []models.ResourceAudiencePolicy, _a1 error) *MockpagePoliciesStore_GetPagePoliciesByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockpagePoliciesStore_GetPagePoliciesByEmail_Call) RunAndReturn(run func(context.Context, uuid.UUID, string) ([]models.ResourceAudiencePolicy, error)) *MockpagePoliciesStore_GetPagePoliciesByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetPagesPolicies provides a mock function with given fields: ctx, pageId
func (_m *MockpagePoliciesStore) GetPagesPolicies(ctx context.Context, pageId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, pageId)

	if len(ret) == 0 {
		panic("no return value specified for GetPagesPolicies")
	}

	var r0 []models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, pageId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, pageId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, pageId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockpagePoliciesStore_GetPagesPolicies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPagesPolicies'
type MockpagePoliciesStore_GetPagesPolicies_Call struct {
	*mock.Call
}

// GetPagesPolicies is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
func (_e *MockpagePoliciesStore_Expecter) GetPagesPolicies(ctx interface{}, pageId interface{}) *MockpagePoliciesStore_GetPagesPolicies_Call {
	return &MockpagePoliciesStore_GetPagesPolicies_Call{Call: _e.mock.On("GetPagesPolicies", ctx, pageId)}
}

func (_c *MockpagePoliciesStore_GetPagesPolicies_Call) Run(run func(ctx context.Context, pageId uuid.UUID)) *MockpagePoliciesStore_GetPagesPolicies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockpagePoliciesStore_GetPagesPolicies_Call) Return(_a0 []models.ResourceAudiencePolicy, _a1 error) *MockpagePoliciesStore_GetPagesPolicies_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockpagePoliciesStore_GetPagesPolicies_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)) *MockpagePoliciesStore_GetPagesPolicies_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePagePolicy provides a mock function with given fields: ctx, pageId, audienceId, privilege
func (_m *MockpagePoliciesStore) UpdatePagePolicy(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
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

// MockpagePoliciesStore_UpdatePagePolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePagePolicy'
type MockpagePoliciesStore_UpdatePagePolicy_Call struct {
	*mock.Call
}

// UpdatePagePolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - pageId uuid.UUID
//   - audienceId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockpagePoliciesStore_Expecter) UpdatePagePolicy(ctx interface{}, pageId interface{}, audienceId interface{}, privilege interface{}) *MockpagePoliciesStore_UpdatePagePolicy_Call {
	return &MockpagePoliciesStore_UpdatePagePolicy_Call{Call: _e.mock.On("UpdatePagePolicy", ctx, pageId, audienceId, privilege)}
}

func (_c *MockpagePoliciesStore_UpdatePagePolicy_Call) Run(run func(ctx context.Context, pageId uuid.UUID, audienceId uuid.UUID, privilege models.ResourcePrivilege)) *MockpagePoliciesStore_UpdatePagePolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockpagePoliciesStore_UpdatePagePolicy_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockpagePoliciesStore_UpdatePagePolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockpagePoliciesStore_UpdatePagePolicy_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockpagePoliciesStore_UpdatePagePolicy_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockpagePoliciesStore creates a new instance of MockpagePoliciesStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockpagePoliciesStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockpagePoliciesStore {
	mock := &MockpagePoliciesStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
