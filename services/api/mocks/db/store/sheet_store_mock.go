// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockSheetStore is an autogenerated mock type for the SheetStore type
type MockSheetStore struct {
	mock.Mock
}

type MockSheetStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSheetStore) EXPECT() *MockSheetStore_Expecter {
	return &MockSheetStore_Expecter{mock: &_m.Mock}
}

// CreateSheet provides a mock function with given fields: ctx, sheet
func (_m *MockSheetStore) CreateSheet(ctx context.Context, sheet models.Sheet) (*models.Sheet, error) {
	ret := _m.Called(ctx, sheet)

	if len(ret) == 0 {
		panic("no return value specified for CreateSheet")
	}

	var r0 *models.Sheet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Sheet) (*models.Sheet, error)); ok {
		return rf(ctx, sheet)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Sheet) *models.Sheet); ok {
		r0 = rf(ctx, sheet)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Sheet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Sheet) error); ok {
		r1 = rf(ctx, sheet)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSheetStore_CreateSheet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSheet'
type MockSheetStore_CreateSheet_Call struct {
	*mock.Call
}

// CreateSheet is a helper method to define mock.On call
//   - ctx context.Context
//   - sheet models.Sheet
func (_e *MockSheetStore_Expecter) CreateSheet(ctx interface{}, sheet interface{}) *MockSheetStore_CreateSheet_Call {
	return &MockSheetStore_CreateSheet_Call{Call: _e.mock.On("CreateSheet", ctx, sheet)}
}

func (_c *MockSheetStore_CreateSheet_Call) Run(run func(ctx context.Context, sheet models.Sheet)) *MockSheetStore_CreateSheet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Sheet))
	})
	return _c
}

func (_c *MockSheetStore_CreateSheet_Call) Return(_a0 *models.Sheet, _a1 error) *MockSheetStore_CreateSheet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSheetStore_CreateSheet_Call) RunAndReturn(run func(context.Context, models.Sheet) (*models.Sheet, error)) *MockSheetStore_CreateSheet_Call {
	_c.Call.Return(run)
	return _c
}

// GetSheetById provides a mock function with given fields: ctx, sheetId
func (_m *MockSheetStore) GetSheetById(ctx context.Context, sheetId uuid.UUID) (*models.Sheet, error) {
	ret := _m.Called(ctx, sheetId)

	if len(ret) == 0 {
		panic("no return value specified for GetSheetById")
	}

	var r0 *models.Sheet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.Sheet, error)); ok {
		return rf(ctx, sheetId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.Sheet); ok {
		r0 = rf(ctx, sheetId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Sheet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, sheetId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSheetStore_GetSheetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSheetById'
type MockSheetStore_GetSheetById_Call struct {
	*mock.Call
}

// GetSheetById is a helper method to define mock.On call
//   - ctx context.Context
//   - sheetId uuid.UUID
func (_e *MockSheetStore_Expecter) GetSheetById(ctx interface{}, sheetId interface{}) *MockSheetStore_GetSheetById_Call {
	return &MockSheetStore_GetSheetById_Call{Call: _e.mock.On("GetSheetById", ctx, sheetId)}
}

func (_c *MockSheetStore_GetSheetById_Call) Run(run func(ctx context.Context, sheetId uuid.UUID)) *MockSheetStore_GetSheetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockSheetStore_GetSheetById_Call) Return(_a0 *models.Sheet, _a1 error) *MockSheetStore_GetSheetById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSheetStore_GetSheetById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.Sheet, error)) *MockSheetStore_GetSheetById_Call {
	_c.Call.Return(run)
	return _c
}

// GetSheetsAll provides a mock function with given fields: ctx, filters
func (_m *MockSheetStore) GetSheetsAll(ctx context.Context, filters models.SheetFilters) ([]models.Sheet, error) {
	ret := _m.Called(ctx, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetSheetsAll")
	}

	var r0 []models.Sheet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.SheetFilters) ([]models.Sheet, error)); ok {
		return rf(ctx, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.SheetFilters) []models.Sheet); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Sheet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.SheetFilters) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSheetStore_GetSheetsAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSheetsAll'
type MockSheetStore_GetSheetsAll_Call struct {
	*mock.Call
}

// GetSheetsAll is a helper method to define mock.On call
//   - ctx context.Context
//   - filters models.SheetFilters
func (_e *MockSheetStore_Expecter) GetSheetsAll(ctx interface{}, filters interface{}) *MockSheetStore_GetSheetsAll_Call {
	return &MockSheetStore_GetSheetsAll_Call{Call: _e.mock.On("GetSheetsAll", ctx, filters)}
}

func (_c *MockSheetStore_GetSheetsAll_Call) Run(run func(ctx context.Context, filters models.SheetFilters)) *MockSheetStore_GetSheetsAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.SheetFilters))
	})
	return _c
}

func (_c *MockSheetStore_GetSheetsAll_Call) Return(_a0 []models.Sheet, _a1 error) *MockSheetStore_GetSheetsAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSheetStore_GetSheetsAll_Call) RunAndReturn(run func(context.Context, models.SheetFilters) ([]models.Sheet, error)) *MockSheetStore_GetSheetsAll_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSheet provides a mock function with given fields: ctx, sheet
func (_m *MockSheetStore) UpdateSheet(ctx context.Context, sheet *models.Sheet) (*models.Sheet, error) {
	ret := _m.Called(ctx, sheet)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSheet")
	}

	var r0 *models.Sheet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Sheet) (*models.Sheet, error)); ok {
		return rf(ctx, sheet)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Sheet) *models.Sheet); ok {
		r0 = rf(ctx, sheet)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Sheet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Sheet) error); ok {
		r1 = rf(ctx, sheet)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSheetStore_UpdateSheet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSheet'
type MockSheetStore_UpdateSheet_Call struct {
	*mock.Call
}

// UpdateSheet is a helper method to define mock.On call
//   - ctx context.Context
//   - sheet *models.Sheet
func (_e *MockSheetStore_Expecter) UpdateSheet(ctx interface{}, sheet interface{}) *MockSheetStore_UpdateSheet_Call {
	return &MockSheetStore_UpdateSheet_Call{Call: _e.mock.On("UpdateSheet", ctx, sheet)}
}

func (_c *MockSheetStore_UpdateSheet_Call) Run(run func(ctx context.Context, sheet *models.Sheet)) *MockSheetStore_UpdateSheet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Sheet))
	})
	return _c
}

func (_c *MockSheetStore_UpdateSheet_Call) Return(_a0 *models.Sheet, _a1 error) *MockSheetStore_UpdateSheet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSheetStore_UpdateSheet_Call) RunAndReturn(run func(context.Context, *models.Sheet) (*models.Sheet, error)) *MockSheetStore_UpdateSheet_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSheetStore creates a new instance of MockSheetStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSheetStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSheetStore {
	mock := &MockSheetStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
