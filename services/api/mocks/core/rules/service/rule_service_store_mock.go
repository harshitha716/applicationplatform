// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_service

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockRuleServiceStore is an autogenerated mock type for the RuleServiceStore type
type MockRuleServiceStore struct {
	mock.Mock
}

type MockRuleServiceStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRuleServiceStore) EXPECT() *MockRuleServiceStore_Expecter {
	return &MockRuleServiceStore_Expecter{mock: &_m.Mock}
}

// CreateRule provides a mock function with given fields: ctx, params
func (_m *MockRuleServiceStore) CreateRule(ctx context.Context, params models.CreateRuleParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for CreateRule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateRuleParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRuleServiceStore_CreateRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRule'
type MockRuleServiceStore_CreateRule_Call struct {
	*mock.Call
}

// CreateRule is a helper method to define mock.On call
//   - ctx context.Context
//   - params models.CreateRuleParams
func (_e *MockRuleServiceStore_Expecter) CreateRule(ctx interface{}, params interface{}) *MockRuleServiceStore_CreateRule_Call {
	return &MockRuleServiceStore_CreateRule_Call{Call: _e.mock.On("CreateRule", ctx, params)}
}

func (_c *MockRuleServiceStore_CreateRule_Call) Run(run func(ctx context.Context, params models.CreateRuleParams)) *MockRuleServiceStore_CreateRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.CreateRuleParams))
	})
	return _c
}

func (_c *MockRuleServiceStore_CreateRule_Call) Return(_a0 error) *MockRuleServiceStore_CreateRule_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRuleServiceStore_CreateRule_Call) RunAndReturn(run func(context.Context, models.CreateRuleParams) error) *MockRuleServiceStore_CreateRule_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteRule provides a mock function with given fields: ctx, params
func (_m *MockRuleServiceStore) DeleteRule(ctx context.Context, params models.DeleteRuleParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.DeleteRuleParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRuleServiceStore_DeleteRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRule'
type MockRuleServiceStore_DeleteRule_Call struct {
	*mock.Call
}

// DeleteRule is a helper method to define mock.On call
//   - ctx context.Context
//   - params models.DeleteRuleParams
func (_e *MockRuleServiceStore_Expecter) DeleteRule(ctx interface{}, params interface{}) *MockRuleServiceStore_DeleteRule_Call {
	return &MockRuleServiceStore_DeleteRule_Call{Call: _e.mock.On("DeleteRule", ctx, params)}
}

func (_c *MockRuleServiceStore_DeleteRule_Call) Run(run func(ctx context.Context, params models.DeleteRuleParams)) *MockRuleServiceStore_DeleteRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.DeleteRuleParams))
	})
	return _c
}

func (_c *MockRuleServiceStore_DeleteRule_Call) Return(_a0 error) *MockRuleServiceStore_DeleteRule_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRuleServiceStore_DeleteRule_Call) RunAndReturn(run func(context.Context, models.DeleteRuleParams) error) *MockRuleServiceStore_DeleteRule_Call {
	_c.Call.Return(run)
	return _c
}

// GetRuleById provides a mock function with given fields: ctx, ruleId
func (_m *MockRuleServiceStore) GetRuleById(ctx context.Context, ruleId uuid.UUID) (models.Rule, error) {
	ret := _m.Called(ctx, ruleId)

	if len(ret) == 0 {
		panic("no return value specified for GetRuleById")
	}

	var r0 models.Rule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (models.Rule, error)); ok {
		return rf(ctx, ruleId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Rule); ok {
		r0 = rf(ctx, ruleId)
	} else {
		r0 = ret.Get(0).(models.Rule)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, ruleId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRuleServiceStore_GetRuleById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRuleById'
type MockRuleServiceStore_GetRuleById_Call struct {
	*mock.Call
}

// GetRuleById is a helper method to define mock.On call
//   - ctx context.Context
//   - ruleId uuid.UUID
func (_e *MockRuleServiceStore_Expecter) GetRuleById(ctx interface{}, ruleId interface{}) *MockRuleServiceStore_GetRuleById_Call {
	return &MockRuleServiceStore_GetRuleById_Call{Call: _e.mock.On("GetRuleById", ctx, ruleId)}
}

func (_c *MockRuleServiceStore_GetRuleById_Call) Run(run func(ctx context.Context, ruleId uuid.UUID)) *MockRuleServiceStore_GetRuleById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRuleServiceStore_GetRuleById_Call) Return(_a0 models.Rule, _a1 error) *MockRuleServiceStore_GetRuleById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRuleServiceStore_GetRuleById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (models.Rule, error)) *MockRuleServiceStore_GetRuleById_Call {
	_c.Call.Return(run)
	return _c
}

// GetRuleByIds provides a mock function with given fields: ctx, ruleIds
func (_m *MockRuleServiceStore) GetRuleByIds(ctx context.Context, ruleIds []uuid.UUID) ([]models.Rule, error) {
	ret := _m.Called(ctx, ruleIds)

	if len(ret) == 0 {
		panic("no return value specified for GetRuleByIds")
	}

	var r0 []models.Rule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) ([]models.Rule, error)); ok {
		return rf(ctx, ruleIds)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) []models.Rule); ok {
		r0 = rf(ctx, ruleIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Rule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uuid.UUID) error); ok {
		r1 = rf(ctx, ruleIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRuleServiceStore_GetRuleByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRuleByIds'
type MockRuleServiceStore_GetRuleByIds_Call struct {
	*mock.Call
}

// GetRuleByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - ruleIds []uuid.UUID
func (_e *MockRuleServiceStore_Expecter) GetRuleByIds(ctx interface{}, ruleIds interface{}) *MockRuleServiceStore_GetRuleByIds_Call {
	return &MockRuleServiceStore_GetRuleByIds_Call{Call: _e.mock.On("GetRuleByIds", ctx, ruleIds)}
}

func (_c *MockRuleServiceStore_GetRuleByIds_Call) Run(run func(ctx context.Context, ruleIds []uuid.UUID)) *MockRuleServiceStore_GetRuleByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]uuid.UUID))
	})
	return _c
}

func (_c *MockRuleServiceStore_GetRuleByIds_Call) Return(_a0 []models.Rule, _a1 error) *MockRuleServiceStore_GetRuleByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRuleServiceStore_GetRuleByIds_Call) RunAndReturn(run func(context.Context, []uuid.UUID) ([]models.Rule, error)) *MockRuleServiceStore_GetRuleByIds_Call {
	_c.Call.Return(run)
	return _c
}

// GetRules provides a mock function with given fields: ctx, params
func (_m *MockRuleServiceStore) GetRules(ctx context.Context, params models.FilterRuleParams) (map[string]map[string][]models.Rule, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetRules")
	}

	var r0 map[string]map[string][]models.Rule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.FilterRuleParams) (map[string]map[string][]models.Rule, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.FilterRuleParams) map[string]map[string][]models.Rule); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]map[string][]models.Rule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.FilterRuleParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRuleServiceStore_GetRules_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRules'
type MockRuleServiceStore_GetRules_Call struct {
	*mock.Call
}

// GetRules is a helper method to define mock.On call
//   - ctx context.Context
//   - params models.FilterRuleParams
func (_e *MockRuleServiceStore_Expecter) GetRules(ctx interface{}, params interface{}) *MockRuleServiceStore_GetRules_Call {
	return &MockRuleServiceStore_GetRules_Call{Call: _e.mock.On("GetRules", ctx, params)}
}

func (_c *MockRuleServiceStore_GetRules_Call) Run(run func(ctx context.Context, params models.FilterRuleParams)) *MockRuleServiceStore_GetRules_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.FilterRuleParams))
	})
	return _c
}

func (_c *MockRuleServiceStore_GetRules_Call) Return(_a0 map[string]map[string][]models.Rule, _a1 error) *MockRuleServiceStore_GetRules_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRuleServiceStore_GetRules_Call) RunAndReturn(run func(context.Context, models.FilterRuleParams) (map[string]map[string][]models.Rule, error)) *MockRuleServiceStore_GetRules_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRule provides a mock function with given fields: ctx, ruleId, params
func (_m *MockRuleServiceStore) UpdateRule(ctx context.Context, ruleId uuid.UUID, params models.UpdateRuleParams) error {
	ret := _m.Called(ctx, ruleId, params)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.UpdateRuleParams) error); ok {
		r0 = rf(ctx, ruleId, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRuleServiceStore_UpdateRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRule'
type MockRuleServiceStore_UpdateRule_Call struct {
	*mock.Call
}

// UpdateRule is a helper method to define mock.On call
//   - ctx context.Context
//   - ruleId uuid.UUID
//   - params models.UpdateRuleParams
func (_e *MockRuleServiceStore_Expecter) UpdateRule(ctx interface{}, ruleId interface{}, params interface{}) *MockRuleServiceStore_UpdateRule_Call {
	return &MockRuleServiceStore_UpdateRule_Call{Call: _e.mock.On("UpdateRule", ctx, ruleId, params)}
}

func (_c *MockRuleServiceStore_UpdateRule_Call) Run(run func(ctx context.Context, ruleId uuid.UUID, params models.UpdateRuleParams)) *MockRuleServiceStore_UpdateRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.UpdateRuleParams))
	})
	return _c
}

func (_c *MockRuleServiceStore_UpdateRule_Call) Return(_a0 error) *MockRuleServiceStore_UpdateRule_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRuleServiceStore_UpdateRule_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.UpdateRuleParams) error) *MockRuleServiceStore_UpdateRule_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRulePriority provides a mock function with given fields: ctx, params
func (_m *MockRuleServiceStore) UpdateRulePriority(ctx context.Context, params models.UpdateRulePriorityParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRulePriority")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UpdateRulePriorityParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRuleServiceStore_UpdateRulePriority_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRulePriority'
type MockRuleServiceStore_UpdateRulePriority_Call struct {
	*mock.Call
}

// UpdateRulePriority is a helper method to define mock.On call
//   - ctx context.Context
//   - params models.UpdateRulePriorityParams
func (_e *MockRuleServiceStore_Expecter) UpdateRulePriority(ctx interface{}, params interface{}) *MockRuleServiceStore_UpdateRulePriority_Call {
	return &MockRuleServiceStore_UpdateRulePriority_Call{Call: _e.mock.On("UpdateRulePriority", ctx, params)}
}

func (_c *MockRuleServiceStore_UpdateRulePriority_Call) Run(run func(ctx context.Context, params models.UpdateRulePriorityParams)) *MockRuleServiceStore_UpdateRulePriority_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UpdateRulePriorityParams))
	})
	return _c
}

func (_c *MockRuleServiceStore_UpdateRulePriority_Call) Return(_a0 error) *MockRuleServiceStore_UpdateRulePriority_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRuleServiceStore_UpdateRulePriority_Call) RunAndReturn(run func(context.Context, models.UpdateRulePriorityParams) error) *MockRuleServiceStore_UpdateRulePriority_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRuleServiceStore creates a new instance of MockRuleServiceStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRuleServiceStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRuleServiceStore {
	mock := &MockRuleServiceStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
