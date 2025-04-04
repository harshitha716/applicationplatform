// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_service

import (
	context "context"

	actionsmodels "github.com/Zampfi/application-platform/services/api/core/datasets/actions/models"

	mock "github.com/stretchr/testify/mock"

	models "github.com/Zampfi/application-platform/services/api/db/models"

	uuid "github.com/google/uuid"
)

// MockDatasetActionService is an autogenerated mock type for the DatasetActionService type
type MockDatasetActionService struct {
	mock.Mock
}

type MockDatasetActionService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatasetActionService) EXPECT() *MockDatasetActionService_Expecter {
	return &MockDatasetActionService_Expecter{mock: &_m.Mock}
}

// CreateDatasetAction provides a mock function with given fields: ctx, organizationId, params
func (_m *MockDatasetActionService) CreateDatasetAction(ctx context.Context, organizationId uuid.UUID, params models.CreateDatasetActionParams) error {
	ret := _m.Called(ctx, organizationId, params)

	if len(ret) == 0 {
		panic("no return value specified for CreateDatasetAction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.CreateDatasetActionParams) error); ok {
		r0 = rf(ctx, organizationId, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatasetActionService_CreateDatasetAction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateDatasetAction'
type MockDatasetActionService_CreateDatasetAction_Call struct {
	*mock.Call
}

// CreateDatasetAction is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - params models.CreateDatasetActionParams
func (_e *MockDatasetActionService_Expecter) CreateDatasetAction(ctx interface{}, organizationId interface{}, params interface{}) *MockDatasetActionService_CreateDatasetAction_Call {
	return &MockDatasetActionService_CreateDatasetAction_Call{Call: _e.mock.On("CreateDatasetAction", ctx, organizationId, params)}
}

func (_c *MockDatasetActionService_CreateDatasetAction_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, params models.CreateDatasetActionParams)) *MockDatasetActionService_CreateDatasetAction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.CreateDatasetActionParams))
	})
	return _c
}

func (_c *MockDatasetActionService_CreateDatasetAction_Call) Return(_a0 error) *MockDatasetActionService_CreateDatasetAction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatasetActionService_CreateDatasetAction_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.CreateDatasetActionParams) error) *MockDatasetActionService_CreateDatasetAction_Call {
	_c.Call.Return(run)
	return _c
}

// GetDatasetActionFromActionId provides a mock function with given fields: ctx, actionId
func (_m *MockDatasetActionService) GetDatasetActionFromActionId(ctx context.Context, actionId string) (*actionsmodels.DatasetAction, error) {
	ret := _m.Called(ctx, actionId)

	if len(ret) == 0 {
		panic("no return value specified for GetDatasetActionFromActionId")
	}

	var r0 *actionsmodels.DatasetAction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*actionsmodels.DatasetAction, error)); ok {
		return rf(ctx, actionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *actionsmodels.DatasetAction); ok {
		r0 = rf(ctx, actionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*actionsmodels.DatasetAction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, actionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatasetActionService_GetDatasetActionFromActionId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDatasetActionFromActionId'
type MockDatasetActionService_GetDatasetActionFromActionId_Call struct {
	*mock.Call
}

// GetDatasetActionFromActionId is a helper method to define mock.On call
//   - ctx context.Context
//   - actionId string
func (_e *MockDatasetActionService_Expecter) GetDatasetActionFromActionId(ctx interface{}, actionId interface{}) *MockDatasetActionService_GetDatasetActionFromActionId_Call {
	return &MockDatasetActionService_GetDatasetActionFromActionId_Call{Call: _e.mock.On("GetDatasetActionFromActionId", ctx, actionId)}
}

func (_c *MockDatasetActionService_GetDatasetActionFromActionId_Call) Run(run func(ctx context.Context, actionId string)) *MockDatasetActionService_GetDatasetActionFromActionId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatasetActionService_GetDatasetActionFromActionId_Call) Return(_a0 *actionsmodels.DatasetAction, _a1 error) *MockDatasetActionService_GetDatasetActionFromActionId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatasetActionService_GetDatasetActionFromActionId_Call) RunAndReturn(run func(context.Context, string) (*actionsmodels.DatasetAction, error)) *MockDatasetActionService_GetDatasetActionFromActionId_Call {
	_c.Call.Return(run)
	return _c
}

// GetDatasetActions provides a mock function with given fields: ctx, organizationId, filters
func (_m *MockDatasetActionService) GetDatasetActions(ctx context.Context, organizationId uuid.UUID, filters models.DatasetActionFilters) ([]actionsmodels.DatasetAction, error) {
	ret := _m.Called(ctx, organizationId, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetDatasetActions")
	}

	var r0 []actionsmodels.DatasetAction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.DatasetActionFilters) ([]actionsmodels.DatasetAction, error)); ok {
		return rf(ctx, organizationId, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.DatasetActionFilters) []actionsmodels.DatasetAction); ok {
		r0 = rf(ctx, organizationId, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]actionsmodels.DatasetAction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.DatasetActionFilters) error); ok {
		r1 = rf(ctx, organizationId, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatasetActionService_GetDatasetActions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDatasetActions'
type MockDatasetActionService_GetDatasetActions_Call struct {
	*mock.Call
}

// GetDatasetActions is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - filters models.DatasetActionFilters
func (_e *MockDatasetActionService_Expecter) GetDatasetActions(ctx interface{}, organizationId interface{}, filters interface{}) *MockDatasetActionService_GetDatasetActions_Call {
	return &MockDatasetActionService_GetDatasetActions_Call{Call: _e.mock.On("GetDatasetActions", ctx, organizationId, filters)}
}

func (_c *MockDatasetActionService_GetDatasetActions_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, filters models.DatasetActionFilters)) *MockDatasetActionService_GetDatasetActions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.DatasetActionFilters))
	})
	return _c
}

func (_c *MockDatasetActionService_GetDatasetActions_Call) Return(_a0 []actionsmodels.DatasetAction, _a1 error) *MockDatasetActionService_GetDatasetActions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatasetActionService_GetDatasetActions_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.DatasetActionFilters) ([]actionsmodels.DatasetAction, error)) *MockDatasetActionService_GetDatasetActions_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDatasetActionConfig provides a mock function with given fields: ctx, actionId, config
func (_m *MockDatasetActionService) UpdateDatasetActionConfig(ctx context.Context, actionId string, config map[string]interface{}) error {
	ret := _m.Called(ctx, actionId, config)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDatasetActionConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) error); ok {
		r0 = rf(ctx, actionId, config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatasetActionService_UpdateDatasetActionConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDatasetActionConfig'
type MockDatasetActionService_UpdateDatasetActionConfig_Call struct {
	*mock.Call
}

// UpdateDatasetActionConfig is a helper method to define mock.On call
//   - ctx context.Context
//   - actionId string
//   - config map[string]interface{}
func (_e *MockDatasetActionService_Expecter) UpdateDatasetActionConfig(ctx interface{}, actionId interface{}, config interface{}) *MockDatasetActionService_UpdateDatasetActionConfig_Call {
	return &MockDatasetActionService_UpdateDatasetActionConfig_Call{Call: _e.mock.On("UpdateDatasetActionConfig", ctx, actionId, config)}
}

func (_c *MockDatasetActionService_UpdateDatasetActionConfig_Call) Run(run func(ctx context.Context, actionId string, config map[string]interface{})) *MockDatasetActionService_UpdateDatasetActionConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *MockDatasetActionService_UpdateDatasetActionConfig_Call) Return(_a0 error) *MockDatasetActionService_UpdateDatasetActionConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatasetActionService_UpdateDatasetActionConfig_Call) RunAndReturn(run func(context.Context, string, map[string]interface{}) error) *MockDatasetActionService_UpdateDatasetActionConfig_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDatasetActionStatus provides a mock function with given fields: ctx, actionId, status
func (_m *MockDatasetActionService) UpdateDatasetActionStatus(ctx context.Context, actionId string, status string) error {
	ret := _m.Called(ctx, actionId, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDatasetActionStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, actionId, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatasetActionService_UpdateDatasetActionStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDatasetActionStatus'
type MockDatasetActionService_UpdateDatasetActionStatus_Call struct {
	*mock.Call
}

// UpdateDatasetActionStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - actionId string
//   - status string
func (_e *MockDatasetActionService_Expecter) UpdateDatasetActionStatus(ctx interface{}, actionId interface{}, status interface{}) *MockDatasetActionService_UpdateDatasetActionStatus_Call {
	return &MockDatasetActionService_UpdateDatasetActionStatus_Call{Call: _e.mock.On("UpdateDatasetActionStatus", ctx, actionId, status)}
}

func (_c *MockDatasetActionService_UpdateDatasetActionStatus_Call) Run(run func(ctx context.Context, actionId string, status string)) *MockDatasetActionService_UpdateDatasetActionStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDatasetActionService_UpdateDatasetActionStatus_Call) Return(_a0 error) *MockDatasetActionService_UpdateDatasetActionStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatasetActionService_UpdateDatasetActionStatus_Call) RunAndReturn(run func(context.Context, string, string) error) *MockDatasetActionService_UpdateDatasetActionStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatasetActionService creates a new instance of MockDatasetActionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatasetActionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatasetActionService {
	mock := &MockDatasetActionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
