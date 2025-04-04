// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"
	json "encoding/json"

	mock "github.com/stretchr/testify/mock"

	models "github.com/Zampfi/application-platform/services/api/db/models"

	uuid "github.com/google/uuid"
)

// MockPaymentsConfigStore is an autogenerated mock type for the PaymentsConfigStore type
type MockPaymentsConfigStore struct {
	mock.Mock
}

type MockPaymentsConfigStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPaymentsConfigStore) EXPECT() *MockPaymentsConfigStore_Expecter {
	return &MockPaymentsConfigStore_Expecter{mock: &_m.Mock}
}

// CreatePaymentsConfig provides a mock function with given fields: ctx, paymentsConfig
func (_m *MockPaymentsConfigStore) CreatePaymentsConfig(ctx context.Context, paymentsConfig models.PaymentsConfig) (models.PaymentsConfig, error) {
	ret := _m.Called(ctx, paymentsConfig)

	if len(ret) == 0 {
		panic("no return value specified for CreatePaymentsConfig")
	}

	var r0 models.PaymentsConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PaymentsConfig) (models.PaymentsConfig, error)); ok {
		return rf(ctx, paymentsConfig)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PaymentsConfig) models.PaymentsConfig); ok {
		r0 = rf(ctx, paymentsConfig)
	} else {
		r0 = ret.Get(0).(models.PaymentsConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.PaymentsConfig) error); ok {
		r1 = rf(ctx, paymentsConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPaymentsConfigStore_CreatePaymentsConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePaymentsConfig'
type MockPaymentsConfigStore_CreatePaymentsConfig_Call struct {
	*mock.Call
}

// CreatePaymentsConfig is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentsConfig models.PaymentsConfig
func (_e *MockPaymentsConfigStore_Expecter) CreatePaymentsConfig(ctx interface{}, paymentsConfig interface{}) *MockPaymentsConfigStore_CreatePaymentsConfig_Call {
	return &MockPaymentsConfigStore_CreatePaymentsConfig_Call{Call: _e.mock.On("CreatePaymentsConfig", ctx, paymentsConfig)}
}

func (_c *MockPaymentsConfigStore_CreatePaymentsConfig_Call) Run(run func(ctx context.Context, paymentsConfig models.PaymentsConfig)) *MockPaymentsConfigStore_CreatePaymentsConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.PaymentsConfig))
	})
	return _c
}

func (_c *MockPaymentsConfigStore_CreatePaymentsConfig_Call) Return(_a0 models.PaymentsConfig, _a1 error) *MockPaymentsConfigStore_CreatePaymentsConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPaymentsConfigStore_CreatePaymentsConfig_Call) RunAndReturn(run func(context.Context, models.PaymentsConfig) (models.PaymentsConfig, error)) *MockPaymentsConfigStore_CreatePaymentsConfig_Call {
	_c.Call.Return(run)
	return _c
}

// DeletePaymentsConfigById provides a mock function with given fields: ctx, paymentsConfigId
func (_m *MockPaymentsConfigStore) DeletePaymentsConfigById(ctx context.Context, paymentsConfigId uuid.UUID) error {
	ret := _m.Called(ctx, paymentsConfigId)

	if len(ret) == 0 {
		panic("no return value specified for DeletePaymentsConfigById")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, paymentsConfigId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPaymentsConfigStore_DeletePaymentsConfigById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeletePaymentsConfigById'
type MockPaymentsConfigStore_DeletePaymentsConfigById_Call struct {
	*mock.Call
}

// DeletePaymentsConfigById is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentsConfigId uuid.UUID
func (_e *MockPaymentsConfigStore_Expecter) DeletePaymentsConfigById(ctx interface{}, paymentsConfigId interface{}) *MockPaymentsConfigStore_DeletePaymentsConfigById_Call {
	return &MockPaymentsConfigStore_DeletePaymentsConfigById_Call{Call: _e.mock.On("DeletePaymentsConfigById", ctx, paymentsConfigId)}
}

func (_c *MockPaymentsConfigStore_DeletePaymentsConfigById_Call) Run(run func(ctx context.Context, paymentsConfigId uuid.UUID)) *MockPaymentsConfigStore_DeletePaymentsConfigById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockPaymentsConfigStore_DeletePaymentsConfigById_Call) Return(_a0 error) *MockPaymentsConfigStore_DeletePaymentsConfigById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPaymentsConfigStore_DeletePaymentsConfigById_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *MockPaymentsConfigStore_DeletePaymentsConfigById_Call {
	_c.Call.Return(run)
	return _c
}

// GetPaymentsConfigsByOrganizationId provides a mock function with given fields: ctx, organizationId
func (_m *MockPaymentsConfigStore) GetPaymentsConfigsByOrganizationId(ctx context.Context, organizationId string) (models.PaymentsConfig, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetPaymentsConfigsByOrganizationId")
	}

	var r0 models.PaymentsConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.PaymentsConfig, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.PaymentsConfig); ok {
		r0 = rf(ctx, organizationId)
	} else {
		r0 = ret.Get(0).(models.PaymentsConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPaymentsConfigsByOrganizationId'
type MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call struct {
	*mock.Call
}

// GetPaymentsConfigsByOrganizationId is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId string
func (_e *MockPaymentsConfigStore_Expecter) GetPaymentsConfigsByOrganizationId(ctx interface{}, organizationId interface{}) *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call {
	return &MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call{Call: _e.mock.On("GetPaymentsConfigsByOrganizationId", ctx, organizationId)}
}

func (_c *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call) Run(run func(ctx context.Context, organizationId string)) *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call) Return(_a0 models.PaymentsConfig, _a1 error) *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call) RunAndReturn(run func(context.Context, string) (models.PaymentsConfig, error)) *MockPaymentsConfigStore_GetPaymentsConfigsByOrganizationId_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePaymentsConfig provides a mock function with given fields: ctx, paymentsConfigId, config
func (_m *MockPaymentsConfigStore) UpdatePaymentsConfig(ctx context.Context, paymentsConfigId uuid.UUID, config json.RawMessage) (models.PaymentsConfig, error) {
	ret := _m.Called(ctx, paymentsConfigId, config)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePaymentsConfig")
	}

	var r0 models.PaymentsConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, json.RawMessage) (models.PaymentsConfig, error)); ok {
		return rf(ctx, paymentsConfigId, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, json.RawMessage) models.PaymentsConfig); ok {
		r0 = rf(ctx, paymentsConfigId, config)
	} else {
		r0 = ret.Get(0).(models.PaymentsConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, json.RawMessage) error); ok {
		r1 = rf(ctx, paymentsConfigId, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPaymentsConfigStore_UpdatePaymentsConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePaymentsConfig'
type MockPaymentsConfigStore_UpdatePaymentsConfig_Call struct {
	*mock.Call
}

// UpdatePaymentsConfig is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentsConfigId uuid.UUID
//   - config json.RawMessage
func (_e *MockPaymentsConfigStore_Expecter) UpdatePaymentsConfig(ctx interface{}, paymentsConfigId interface{}, config interface{}) *MockPaymentsConfigStore_UpdatePaymentsConfig_Call {
	return &MockPaymentsConfigStore_UpdatePaymentsConfig_Call{Call: _e.mock.On("UpdatePaymentsConfig", ctx, paymentsConfigId, config)}
}

func (_c *MockPaymentsConfigStore_UpdatePaymentsConfig_Call) Run(run func(ctx context.Context, paymentsConfigId uuid.UUID, config json.RawMessage)) *MockPaymentsConfigStore_UpdatePaymentsConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(json.RawMessage))
	})
	return _c
}

func (_c *MockPaymentsConfigStore_UpdatePaymentsConfig_Call) Return(_a0 models.PaymentsConfig, _a1 error) *MockPaymentsConfigStore_UpdatePaymentsConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPaymentsConfigStore_UpdatePaymentsConfig_Call) RunAndReturn(run func(context.Context, uuid.UUID, json.RawMessage) (models.PaymentsConfig, error)) *MockPaymentsConfigStore_UpdatePaymentsConfig_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePaymentsConfigStatus provides a mock function with given fields: ctx, paymentsConfigId, status
func (_m *MockPaymentsConfigStore) UpdatePaymentsConfigStatus(ctx context.Context, paymentsConfigId uuid.UUID, status models.PaymentsConfigStatus) (models.PaymentsConfig, error) {
	ret := _m.Called(ctx, paymentsConfigId, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePaymentsConfigStatus")
	}

	var r0 models.PaymentsConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.PaymentsConfigStatus) (models.PaymentsConfig, error)); ok {
		return rf(ctx, paymentsConfigId, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.PaymentsConfigStatus) models.PaymentsConfig); ok {
		r0 = rf(ctx, paymentsConfigId, status)
	} else {
		r0 = ret.Get(0).(models.PaymentsConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.PaymentsConfigStatus) error); ok {
		r1 = rf(ctx, paymentsConfigId, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePaymentsConfigStatus'
type MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call struct {
	*mock.Call
}

// UpdatePaymentsConfigStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentsConfigId uuid.UUID
//   - status models.PaymentsConfigStatus
func (_e *MockPaymentsConfigStore_Expecter) UpdatePaymentsConfigStatus(ctx interface{}, paymentsConfigId interface{}, status interface{}) *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call {
	return &MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call{Call: _e.mock.On("UpdatePaymentsConfigStatus", ctx, paymentsConfigId, status)}
}

func (_c *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call) Run(run func(ctx context.Context, paymentsConfigId uuid.UUID, status models.PaymentsConfigStatus)) *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.PaymentsConfigStatus))
	})
	return _c
}

func (_c *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call) Return(_a0 models.PaymentsConfig, _a1 error) *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.PaymentsConfigStatus) (models.PaymentsConfig, error)) *MockPaymentsConfigStore_UpdatePaymentsConfigStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPaymentsConfigStore creates a new instance of MockPaymentsConfigStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPaymentsConfigStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPaymentsConfigStore {
	mock := &MockPaymentsConfigStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
