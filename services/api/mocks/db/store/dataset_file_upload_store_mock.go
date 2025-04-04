// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockDatasetFileUploadStore is an autogenerated mock type for the DatasetFileUploadStore type
type MockDatasetFileUploadStore struct {
	mock.Mock
}

type MockDatasetFileUploadStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatasetFileUploadStore) EXPECT() *MockDatasetFileUploadStore_Expecter {
	return &MockDatasetFileUploadStore_Expecter{mock: &_m.Mock}
}

// CreateDatasetFileUpload provides a mock function with given fields: ctx, datasetFileUpload
func (_m *MockDatasetFileUploadStore) CreateDatasetFileUpload(ctx context.Context, datasetFileUpload *models.DatasetFileUpload) (*models.DatasetFileUpload, error) {
	ret := _m.Called(ctx, datasetFileUpload)

	if len(ret) == 0 {
		panic("no return value specified for CreateDatasetFileUpload")
	}

	var r0 *models.DatasetFileUpload
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.DatasetFileUpload) (*models.DatasetFileUpload, error)); ok {
		return rf(ctx, datasetFileUpload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.DatasetFileUpload) *models.DatasetFileUpload); ok {
		r0 = rf(ctx, datasetFileUpload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DatasetFileUpload)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.DatasetFileUpload) error); ok {
		r1 = rf(ctx, datasetFileUpload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatasetFileUploadStore_CreateDatasetFileUpload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateDatasetFileUpload'
type MockDatasetFileUploadStore_CreateDatasetFileUpload_Call struct {
	*mock.Call
}

// CreateDatasetFileUpload is a helper method to define mock.On call
//   - ctx context.Context
//   - datasetFileUpload *models.DatasetFileUpload
func (_e *MockDatasetFileUploadStore_Expecter) CreateDatasetFileUpload(ctx interface{}, datasetFileUpload interface{}) *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call {
	return &MockDatasetFileUploadStore_CreateDatasetFileUpload_Call{Call: _e.mock.On("CreateDatasetFileUpload", ctx, datasetFileUpload)}
}

func (_c *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call) Run(run func(ctx context.Context, datasetFileUpload *models.DatasetFileUpload)) *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.DatasetFileUpload))
	})
	return _c
}

func (_c *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call) Return(_a0 *models.DatasetFileUpload, _a1 error) *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call) RunAndReturn(run func(context.Context, *models.DatasetFileUpload) (*models.DatasetFileUpload, error)) *MockDatasetFileUploadStore_CreateDatasetFileUpload_Call {
	_c.Call.Return(run)
	return _c
}

// GetDatasetFileUploadByDatasetId provides a mock function with given fields: ctx, datasetId
func (_m *MockDatasetFileUploadStore) GetDatasetFileUploadByDatasetId(ctx context.Context, datasetId uuid.UUID) ([]models.DatasetFileUpload, error) {
	ret := _m.Called(ctx, datasetId)

	if len(ret) == 0 {
		panic("no return value specified for GetDatasetFileUploadByDatasetId")
	}

	var r0 []models.DatasetFileUpload
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.DatasetFileUpload, error)); ok {
		return rf(ctx, datasetId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.DatasetFileUpload); ok {
		r0 = rf(ctx, datasetId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.DatasetFileUpload)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, datasetId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDatasetFileUploadByDatasetId'
type MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call struct {
	*mock.Call
}

// GetDatasetFileUploadByDatasetId is a helper method to define mock.On call
//   - ctx context.Context
//   - datasetId uuid.UUID
func (_e *MockDatasetFileUploadStore_Expecter) GetDatasetFileUploadByDatasetId(ctx interface{}, datasetId interface{}) *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call {
	return &MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call{Call: _e.mock.On("GetDatasetFileUploadByDatasetId", ctx, datasetId)}
}

func (_c *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call) Run(run func(ctx context.Context, datasetId uuid.UUID)) *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call) Return(_a0 []models.DatasetFileUpload, _a1 error) *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.DatasetFileUpload, error)) *MockDatasetFileUploadStore_GetDatasetFileUploadByDatasetId_Call {
	_c.Call.Return(run)
	return _c
}

// GetDatasetFileUploadById provides a mock function with given fields: ctx, fileUploadId
func (_m *MockDatasetFileUploadStore) GetDatasetFileUploadById(ctx context.Context, fileUploadId uuid.UUID) (models.DatasetFileUpload, error) {
	ret := _m.Called(ctx, fileUploadId)

	if len(ret) == 0 {
		panic("no return value specified for GetDatasetFileUploadById")
	}

	var r0 models.DatasetFileUpload
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (models.DatasetFileUpload, error)); ok {
		return rf(ctx, fileUploadId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.DatasetFileUpload); ok {
		r0 = rf(ctx, fileUploadId)
	} else {
		r0 = ret.Get(0).(models.DatasetFileUpload)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, fileUploadId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatasetFileUploadStore_GetDatasetFileUploadById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDatasetFileUploadById'
type MockDatasetFileUploadStore_GetDatasetFileUploadById_Call struct {
	*mock.Call
}

// GetDatasetFileUploadById is a helper method to define mock.On call
//   - ctx context.Context
//   - fileUploadId uuid.UUID
func (_e *MockDatasetFileUploadStore_Expecter) GetDatasetFileUploadById(ctx interface{}, fileUploadId interface{}) *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call {
	return &MockDatasetFileUploadStore_GetDatasetFileUploadById_Call{Call: _e.mock.On("GetDatasetFileUploadById", ctx, fileUploadId)}
}

func (_c *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call) Run(run func(ctx context.Context, fileUploadId uuid.UUID)) *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call) Return(_a0 models.DatasetFileUpload, _a1 error) *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (models.DatasetFileUpload, error)) *MockDatasetFileUploadStore_GetDatasetFileUploadById_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDatasetFileUploadStatus provides a mock function with given fields: ctx, id, fileAllignmentStatus, metadata
func (_m *MockDatasetFileUploadStore) UpdateDatasetFileUploadStatus(ctx context.Context, id uuid.UUID, fileAllignmentStatus models.DatasetFileAllignmentStatus, metadata models.DatasetFileUploadMetadata) (*models.DatasetFileUpload, error) {
	ret := _m.Called(ctx, id, fileAllignmentStatus, metadata)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDatasetFileUploadStatus")
	}

	var r0 *models.DatasetFileUpload
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.DatasetFileAllignmentStatus, models.DatasetFileUploadMetadata) (*models.DatasetFileUpload, error)); ok {
		return rf(ctx, id, fileAllignmentStatus, metadata)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.DatasetFileAllignmentStatus, models.DatasetFileUploadMetadata) *models.DatasetFileUpload); ok {
		r0 = rf(ctx, id, fileAllignmentStatus, metadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DatasetFileUpload)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.DatasetFileAllignmentStatus, models.DatasetFileUploadMetadata) error); ok {
		r1 = rf(ctx, id, fileAllignmentStatus, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDatasetFileUploadStatus'
type MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call struct {
	*mock.Call
}

// UpdateDatasetFileUploadStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - fileAllignmentStatus models.DatasetFileAllignmentStatus
//   - metadata models.DatasetFileUploadMetadata
func (_e *MockDatasetFileUploadStore_Expecter) UpdateDatasetFileUploadStatus(ctx interface{}, id interface{}, fileAllignmentStatus interface{}, metadata interface{}) *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call {
	return &MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call{Call: _e.mock.On("UpdateDatasetFileUploadStatus", ctx, id, fileAllignmentStatus, metadata)}
}

func (_c *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call) Run(run func(ctx context.Context, id uuid.UUID, fileAllignmentStatus models.DatasetFileAllignmentStatus, metadata models.DatasetFileUploadMetadata)) *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.DatasetFileAllignmentStatus), args[3].(models.DatasetFileUploadMetadata))
	})
	return _c
}

func (_c *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call) Return(_a0 *models.DatasetFileUpload, _a1 error) *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.DatasetFileAllignmentStatus, models.DatasetFileUploadMetadata) (*models.DatasetFileUpload, error)) *MockDatasetFileUploadStore_UpdateDatasetFileUploadStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatasetFileUploadStore creates a new instance of MockDatasetFileUploadStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatasetFileUploadStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatasetFileUploadStore {
	mock := &MockDatasetFileUploadStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
