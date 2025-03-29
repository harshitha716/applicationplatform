// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockConnectionStore is an autogenerated mock type for the ConnectionStore type
type MockConnectionStore struct {
	mock.Mock
}

type MockConnectionStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConnectionStore) EXPECT() *MockConnectionStore_Expecter {
	return &MockConnectionStore_Expecter{mock: &_m.Mock}
}

// CreateConnection provides a mock function with given fields: ctx, connection
func (_m *MockConnectionStore) CreateConnection(ctx context.Context, connection *models.CreateConnectionParams) (uuid.UUID, error) {
	ret := _m.Called(ctx, connection)

	if len(ret) == 0 {
		panic("no return value specified for CreateConnection")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateConnectionParams) (uuid.UUID, error)); ok {
		return rf(ctx, connection)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateConnectionParams) uuid.UUID); ok {
		r0 = rf(ctx, connection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.CreateConnectionParams) error); ok {
		r1 = rf(ctx, connection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConnectionStore_CreateConnection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateConnection'
type MockConnectionStore_CreateConnection_Call struct {
	*mock.Call
}

// CreateConnection is a helper method to define mock.On call
//   - ctx context.Context
//   - connection *models.CreateConnectionParams
func (_e *MockConnectionStore_Expecter) CreateConnection(ctx interface{}, connection interface{}) *MockConnectionStore_CreateConnection_Call {
	return &MockConnectionStore_CreateConnection_Call{Call: _e.mock.On("CreateConnection", ctx, connection)}
}

func (_c *MockConnectionStore_CreateConnection_Call) Run(run func(ctx context.Context, connection *models.CreateConnectionParams)) *MockConnectionStore_CreateConnection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.CreateConnectionParams))
	})
	return _c
}

func (_c *MockConnectionStore_CreateConnection_Call) Return(_a0 uuid.UUID, _a1 error) *MockConnectionStore_CreateConnection_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConnectionStore_CreateConnection_Call) RunAndReturn(run func(context.Context, *models.CreateConnectionParams) (uuid.UUID, error)) *MockConnectionStore_CreateConnection_Call {
	_c.Call.Return(run)
	return _c
}

// GetConnectionByID provides a mock function with given fields: ctx, id
func (_m *MockConnectionStore) GetConnectionByID(ctx context.Context, id uuid.UUID) (*models.Connection, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetConnectionByID")
	}

	var r0 *models.Connection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.Connection, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.Connection); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Connection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConnectionStore_GetConnectionByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnectionByID'
type MockConnectionStore_GetConnectionByID_Call struct {
	*mock.Call
}

// GetConnectionByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MockConnectionStore_Expecter) GetConnectionByID(ctx interface{}, id interface{}) *MockConnectionStore_GetConnectionByID_Call {
	return &MockConnectionStore_GetConnectionByID_Call{Call: _e.mock.On("GetConnectionByID", ctx, id)}
}

func (_c *MockConnectionStore_GetConnectionByID_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MockConnectionStore_GetConnectionByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockConnectionStore_GetConnectionByID_Call) Return(_a0 *models.Connection, _a1 error) *MockConnectionStore_GetConnectionByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConnectionStore_GetConnectionByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.Connection, error)) *MockConnectionStore_GetConnectionByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetConnections provides a mock function with given fields: ctx
func (_m *MockConnectionStore) GetConnections(ctx context.Context) ([]models.Connection, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetConnections")
	}

	var r0 []models.Connection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Connection, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Connection); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Connection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConnectionStore_GetConnections_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnections'
type MockConnectionStore_GetConnections_Call struct {
	*mock.Call
}

// GetConnections is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockConnectionStore_Expecter) GetConnections(ctx interface{}) *MockConnectionStore_GetConnections_Call {
	return &MockConnectionStore_GetConnections_Call{Call: _e.mock.On("GetConnections", ctx)}
}

func (_c *MockConnectionStore_GetConnections_Call) Run(run func(ctx context.Context)) *MockConnectionStore_GetConnections_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockConnectionStore_GetConnections_Call) Return(_a0 []models.Connection, _a1 error) *MockConnectionStore_GetConnections_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConnectionStore_GetConnections_Call) RunAndReturn(run func(context.Context) ([]models.Connection, error)) *MockConnectionStore_GetConnections_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConnectionStore creates a new instance of MockConnectionStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConnectionStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConnectionStore {
	mock := &MockConnectionStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
