// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_actions

import (
	context "context"

	dataplatformmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	mock "github.com/stretchr/testify/mock"

	models "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
)

// MockActionService is an autogenerated mock type for the ActionService type
type MockActionService struct {
	mock.Mock
}

type MockActionService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockActionService) EXPECT() *MockActionService_Expecter {
	return &MockActionService_Expecter{mock: &_m.Mock}
}

// CreateAction provides a mock function with given fields: ctx, payload
func (_m *MockActionService) CreateAction(ctx context.Context, payload models.CreateActionPayload) (models.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateAction")
	}

	var r0 models.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateActionPayload) (models.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateActionPayload) models.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(models.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.CreateActionPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockActionService_CreateAction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAction'
type MockActionService_CreateAction_Call struct {
	*mock.Call
}

// CreateAction is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.CreateActionPayload
func (_e *MockActionService_Expecter) CreateAction(ctx interface{}, payload interface{}) *MockActionService_CreateAction_Call {
	return &MockActionService_CreateAction_Call{Call: _e.mock.On("CreateAction", ctx, payload)}
}

func (_c *MockActionService_CreateAction_Call) Run(run func(ctx context.Context, payload models.CreateActionPayload)) *MockActionService_CreateAction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.CreateActionPayload))
	})
	return _c
}

func (_c *MockActionService_CreateAction_Call) Return(_a0 models.CreateActionResponse, _a1 error) *MockActionService_CreateAction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockActionService_CreateAction_Call) RunAndReturn(run func(context.Context, models.CreateActionPayload) (models.CreateActionResponse, error)) *MockActionService_CreateAction_Call {
	_c.Call.Return(run)
	return _c
}

// GetActionById provides a mock function with given fields: ctx, merchantId, actionId
func (_m *MockActionService) GetActionById(ctx context.Context, merchantId string, actionId string) (models.Action, error) {
	ret := _m.Called(ctx, merchantId, actionId)

	if len(ret) == 0 {
		panic("no return value specified for GetActionById")
	}

	var r0 models.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (models.Action, error)); ok {
		return rf(ctx, merchantId, actionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) models.Action); ok {
		r0 = rf(ctx, merchantId, actionId)
	} else {
		r0 = ret.Get(0).(models.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, merchantId, actionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockActionService_GetActionById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActionById'
type MockActionService_GetActionById_Call struct {
	*mock.Call
}

// GetActionById is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
//   - actionId string
func (_e *MockActionService_Expecter) GetActionById(ctx interface{}, merchantId interface{}, actionId interface{}) *MockActionService_GetActionById_Call {
	return &MockActionService_GetActionById_Call{Call: _e.mock.On("GetActionById", ctx, merchantId, actionId)}
}

func (_c *MockActionService_GetActionById_Call) Run(run func(ctx context.Context, merchantId string, actionId string)) *MockActionService_GetActionById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockActionService_GetActionById_Call) Return(_a0 models.Action, _a1 error) *MockActionService_GetActionById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockActionService_GetActionById_Call) RunAndReturn(run func(context.Context, string, string) (models.Action, error)) *MockActionService_GetActionById_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAction provides a mock function with given fields: ctx, jobStatusUpdate
func (_m *MockActionService) UpdateAction(ctx context.Context, jobStatusUpdate dataplatformmodels.DatabricksJobStatusUpdatePayload) (models.Action, error) {
	ret := _m.Called(ctx, jobStatusUpdate)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAction")
	}

	var r0 models.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dataplatformmodels.DatabricksJobStatusUpdatePayload) (models.Action, error)); ok {
		return rf(ctx, jobStatusUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dataplatformmodels.DatabricksJobStatusUpdatePayload) models.Action); ok {
		r0 = rf(ctx, jobStatusUpdate)
	} else {
		r0 = ret.Get(0).(models.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dataplatformmodels.DatabricksJobStatusUpdatePayload) error); ok {
		r1 = rf(ctx, jobStatusUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockActionService_UpdateAction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAction'
type MockActionService_UpdateAction_Call struct {
	*mock.Call
}

// UpdateAction is a helper method to define mock.On call
//   - ctx context.Context
//   - jobStatusUpdate dataplatformmodels.DatabricksJobStatusUpdatePayload
func (_e *MockActionService_Expecter) UpdateAction(ctx interface{}, jobStatusUpdate interface{}) *MockActionService_UpdateAction_Call {
	return &MockActionService_UpdateAction_Call{Call: _e.mock.On("UpdateAction", ctx, jobStatusUpdate)}
}

func (_c *MockActionService_UpdateAction_Call) Run(run func(ctx context.Context, jobStatusUpdate dataplatformmodels.DatabricksJobStatusUpdatePayload)) *MockActionService_UpdateAction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dataplatformmodels.DatabricksJobStatusUpdatePayload))
	})
	return _c
}

func (_c *MockActionService_UpdateAction_Call) Return(_a0 models.Action, _a1 error) *MockActionService_UpdateAction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockActionService_UpdateAction_Call) RunAndReturn(run func(context.Context, dataplatformmodels.DatabricksJobStatusUpdatePayload) (models.Action, error)) *MockActionService_UpdateAction_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockActionService creates a new instance of MockActionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockActionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockActionService {
	mock := &MockActionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
