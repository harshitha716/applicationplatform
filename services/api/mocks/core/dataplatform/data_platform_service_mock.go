// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_dataplatform

import (
	context "context"

	actionsmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"

	datamodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"

	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"

	mock "github.com/stretchr/testify/mock"

	models "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
)

// MockDataPlatformService is an autogenerated mock type for the DataPlatformService type
type MockDataPlatformService struct {
	mock.Mock
}

type MockDataPlatformService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDataPlatformService) EXPECT() *MockDataPlatformService_Expecter {
	return &MockDataPlatformService_Expecter{mock: &_m.Mock}
}

// CopyDataset provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) CopyDataset(ctx context.Context, payload models.CopyDatasetPayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CopyDataset")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CopyDatasetPayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.CopyDatasetPayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.CopyDatasetPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_CopyDataset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CopyDataset'
type MockDataPlatformService_CopyDataset_Call struct {
	*mock.Call
}

// CopyDataset is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.CopyDatasetPayload
func (_e *MockDataPlatformService_Expecter) CopyDataset(ctx interface{}, payload interface{}) *MockDataPlatformService_CopyDataset_Call {
	return &MockDataPlatformService_CopyDataset_Call{Call: _e.mock.On("CopyDataset", ctx, payload)}
}

func (_c *MockDataPlatformService_CopyDataset_Call) Run(run func(ctx context.Context, payload models.CopyDatasetPayload)) *MockDataPlatformService_CopyDataset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.CopyDatasetPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_CopyDataset_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_CopyDataset_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_CopyDataset_Call) RunAndReturn(run func(context.Context, models.CopyDatasetPayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_CopyDataset_Call {
	_c.Call.Return(run)
	return _c
}

// CreateMV provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) CreateMV(ctx context.Context, payload models.CreateMVPayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateMV")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateMVPayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateMVPayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.CreateMVPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_CreateMV_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMV'
type MockDataPlatformService_CreateMV_Call struct {
	*mock.Call
}

// CreateMV is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.CreateMVPayload
func (_e *MockDataPlatformService_Expecter) CreateMV(ctx interface{}, payload interface{}) *MockDataPlatformService_CreateMV_Call {
	return &MockDataPlatformService_CreateMV_Call{Call: _e.mock.On("CreateMV", ctx, payload)}
}

func (_c *MockDataPlatformService_CreateMV_Call) Run(run func(ctx context.Context, payload models.CreateMVPayload)) *MockDataPlatformService_CreateMV_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.CreateMVPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_CreateMV_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_CreateMV_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_CreateMV_Call) RunAndReturn(run func(context.Context, models.CreateMVPayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_CreateMV_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteDataset provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) DeleteDataset(ctx context.Context, payload models.DeleteDatasetPayload) (string, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDataset")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.DeleteDatasetPayload) (string, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.DeleteDatasetPayload) string); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.DeleteDatasetPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_DeleteDataset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteDataset'
type MockDataPlatformService_DeleteDataset_Call struct {
	*mock.Call
}

// DeleteDataset is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.DeleteDatasetPayload
func (_e *MockDataPlatformService_Expecter) DeleteDataset(ctx interface{}, payload interface{}) *MockDataPlatformService_DeleteDataset_Call {
	return &MockDataPlatformService_DeleteDataset_Call{Call: _e.mock.On("DeleteDataset", ctx, payload)}
}

func (_c *MockDataPlatformService_DeleteDataset_Call) Run(run func(ctx context.Context, payload models.DeleteDatasetPayload)) *MockDataPlatformService_DeleteDataset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.DeleteDatasetPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_DeleteDataset_Call) Return(_a0 string, _a1 error) *MockDataPlatformService_DeleteDataset_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_DeleteDataset_Call) RunAndReturn(run func(context.Context, models.DeleteDatasetPayload) (string, error)) *MockDataPlatformService_DeleteDataset_Call {
	_c.Call.Return(run)
	return _c
}

// GetActionById provides a mock function with given fields: ctx, merchantId, actionId
func (_m *MockDataPlatformService) GetActionById(ctx context.Context, merchantId string, actionId string) (actionsmodels.Action, error) {
	ret := _m.Called(ctx, merchantId, actionId)

	if len(ret) == 0 {
		panic("no return value specified for GetActionById")
	}

	var r0 actionsmodels.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (actionsmodels.Action, error)); ok {
		return rf(ctx, merchantId, actionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) actionsmodels.Action); ok {
		r0 = rf(ctx, merchantId, actionId)
	} else {
		r0 = ret.Get(0).(actionsmodels.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, merchantId, actionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_GetActionById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActionById'
type MockDataPlatformService_GetActionById_Call struct {
	*mock.Call
}

// GetActionById is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
//   - actionId string
func (_e *MockDataPlatformService_Expecter) GetActionById(ctx interface{}, merchantId interface{}, actionId interface{}) *MockDataPlatformService_GetActionById_Call {
	return &MockDataPlatformService_GetActionById_Call{Call: _e.mock.On("GetActionById", ctx, merchantId, actionId)}
}

func (_c *MockDataPlatformService_GetActionById_Call) Run(run func(ctx context.Context, merchantId string, actionId string)) *MockDataPlatformService_GetActionById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDataPlatformService_GetActionById_Call) Return(_a0 actionsmodels.Action, _a1 error) *MockDataPlatformService_GetActionById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_GetActionById_Call) RunAndReturn(run func(context.Context, string, string) (actionsmodels.Action, error)) *MockDataPlatformService_GetActionById_Call {
	_c.Call.Return(run)
	return _c
}

// GetDags provides a mock function with given fields: ctx, merchantId
func (_m *MockDataPlatformService) GetDags(ctx context.Context, merchantId string) (map[string]*models.DAGNode, error) {
	ret := _m.Called(ctx, merchantId)

	if len(ret) == 0 {
		panic("no return value specified for GetDags")
	}

	var r0 map[string]*models.DAGNode
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (map[string]*models.DAGNode, error)); ok {
		return rf(ctx, merchantId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]*models.DAGNode); ok {
		r0 = rf(ctx, merchantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*models.DAGNode)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, merchantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_GetDags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDags'
type MockDataPlatformService_GetDags_Call struct {
	*mock.Call
}

// GetDags is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
func (_e *MockDataPlatformService_Expecter) GetDags(ctx interface{}, merchantId interface{}) *MockDataPlatformService_GetDags_Call {
	return &MockDataPlatformService_GetDags_Call{Call: _e.mock.On("GetDags", ctx, merchantId)}
}

func (_c *MockDataPlatformService_GetDags_Call) Run(run func(ctx context.Context, merchantId string)) *MockDataPlatformService_GetDags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDataPlatformService_GetDags_Call) Return(_a0 map[string]*models.DAGNode, _a1 error) *MockDataPlatformService_GetDags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_GetDags_Call) RunAndReturn(run func(context.Context, string) (map[string]*models.DAGNode, error)) *MockDataPlatformService_GetDags_Call {
	_c.Call.Return(run)
	return _c
}

// GetDatasetMetadata provides a mock function with given fields: ctx, merchantId, datasetId
func (_m *MockDataPlatformService) GetDatasetMetadata(ctx context.Context, merchantId string, datasetId string) (datamodels.DatasetMetadata, error) {
	ret := _m.Called(ctx, merchantId, datasetId)

	if len(ret) == 0 {
		panic("no return value specified for GetDatasetMetadata")
	}

	var r0 datamodels.DatasetMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (datamodels.DatasetMetadata, error)); ok {
		return rf(ctx, merchantId, datasetId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) datamodels.DatasetMetadata); ok {
		r0 = rf(ctx, merchantId, datasetId)
	} else {
		r0 = ret.Get(0).(datamodels.DatasetMetadata)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, merchantId, datasetId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_GetDatasetMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDatasetMetadata'
type MockDataPlatformService_GetDatasetMetadata_Call struct {
	*mock.Call
}

// GetDatasetMetadata is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
//   - datasetId string
func (_e *MockDataPlatformService_Expecter) GetDatasetMetadata(ctx interface{}, merchantId interface{}, datasetId interface{}) *MockDataPlatformService_GetDatasetMetadata_Call {
	return &MockDataPlatformService_GetDatasetMetadata_Call{Call: _e.mock.On("GetDatasetMetadata", ctx, merchantId, datasetId)}
}

func (_c *MockDataPlatformService_GetDatasetMetadata_Call) Run(run func(ctx context.Context, merchantId string, datasetId string)) *MockDataPlatformService_GetDatasetMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDataPlatformService_GetDatasetMetadata_Call) Return(_a0 datamodels.DatasetMetadata, _a1 error) *MockDataPlatformService_GetDatasetMetadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_GetDatasetMetadata_Call) RunAndReturn(run func(context.Context, string, string) (datamodels.DatasetMetadata, error)) *MockDataPlatformService_GetDatasetMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// GetDatasetParents provides a mock function with given fields: ctx, merchantId, datasetId
func (_m *MockDataPlatformService) GetDatasetParents(ctx context.Context, merchantId string, datasetId string) (datamodels.DatasetParents, error) {
	ret := _m.Called(ctx, merchantId, datasetId)

	if len(ret) == 0 {
		panic("no return value specified for GetDatasetParents")
	}

	var r0 datamodels.DatasetParents
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (datamodels.DatasetParents, error)); ok {
		return rf(ctx, merchantId, datasetId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) datamodels.DatasetParents); ok {
		r0 = rf(ctx, merchantId, datasetId)
	} else {
		r0 = ret.Get(0).(datamodels.DatasetParents)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, merchantId, datasetId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_GetDatasetParents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDatasetParents'
type MockDataPlatformService_GetDatasetParents_Call struct {
	*mock.Call
}

// GetDatasetParents is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
//   - datasetId string
func (_e *MockDataPlatformService_Expecter) GetDatasetParents(ctx interface{}, merchantId interface{}, datasetId interface{}) *MockDataPlatformService_GetDatasetParents_Call {
	return &MockDataPlatformService_GetDatasetParents_Call{Call: _e.mock.On("GetDatasetParents", ctx, merchantId, datasetId)}
}

func (_c *MockDataPlatformService_GetDatasetParents_Call) Run(run func(ctx context.Context, merchantId string, datasetId string)) *MockDataPlatformService_GetDatasetParents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDataPlatformService_GetDatasetParents_Call) Return(_a0 datamodels.DatasetParents, _a1 error) *MockDataPlatformService_GetDatasetParents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_GetDatasetParents_Call) RunAndReturn(run func(context.Context, string, string) (datamodels.DatasetParents, error)) *MockDataPlatformService_GetDatasetParents_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: ctx, merchantId, query, params, args
func (_m *MockDataPlatformService) Query(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (dataplatformmodels.QueryResult, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, merchantId, query, params)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 dataplatformmodels.QueryResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]string, ...interface{}) (dataplatformmodels.QueryResult, error)); ok {
		return rf(ctx, merchantId, query, params, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]string, ...interface{}) dataplatformmodels.QueryResult); ok {
		r0 = rf(ctx, merchantId, query, params, args...)
	} else {
		r0 = ret.Get(0).(dataplatformmodels.QueryResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, map[string]string, ...interface{}) error); ok {
		r1 = rf(ctx, merchantId, query, params, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type MockDataPlatformService_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
//   - query string
//   - params map[string]string
//   - args ...interface{}
func (_e *MockDataPlatformService_Expecter) Query(ctx interface{}, merchantId interface{}, query interface{}, params interface{}, args ...interface{}) *MockDataPlatformService_Query_Call {
	return &MockDataPlatformService_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{ctx, merchantId, query, params}, args...)...)}
}

func (_c *MockDataPlatformService_Query_Call) Run(run func(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{})) *MockDataPlatformService_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-4)
		for i, a := range args[4:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(map[string]string), variadicArgs...)
	})
	return _c
}

func (_c *MockDataPlatformService_Query_Call) Return(_a0 dataplatformmodels.QueryResult, _a1 error) *MockDataPlatformService_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_Query_Call) RunAndReturn(run func(context.Context, string, string, map[string]string, ...interface{}) (dataplatformmodels.QueryResult, error)) *MockDataPlatformService_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRealTime provides a mock function with given fields: ctx, merchantId, query, params, args
func (_m *MockDataPlatformService) QueryRealTime(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (dataplatformmodels.QueryResult, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, merchantId, query, params)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRealTime")
	}

	var r0 dataplatformmodels.QueryResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]string, ...interface{}) (dataplatformmodels.QueryResult, error)); ok {
		return rf(ctx, merchantId, query, params, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]string, ...interface{}) dataplatformmodels.QueryResult); ok {
		r0 = rf(ctx, merchantId, query, params, args...)
	} else {
		r0 = ret.Get(0).(dataplatformmodels.QueryResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, map[string]string, ...interface{}) error); ok {
		r1 = rf(ctx, merchantId, query, params, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_QueryRealTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRealTime'
type MockDataPlatformService_QueryRealTime_Call struct {
	*mock.Call
}

// QueryRealTime is a helper method to define mock.On call
//   - ctx context.Context
//   - merchantId string
//   - query string
//   - params map[string]string
//   - args ...interface{}
func (_e *MockDataPlatformService_Expecter) QueryRealTime(ctx interface{}, merchantId interface{}, query interface{}, params interface{}, args ...interface{}) *MockDataPlatformService_QueryRealTime_Call {
	return &MockDataPlatformService_QueryRealTime_Call{Call: _e.mock.On("QueryRealTime",
		append([]interface{}{ctx, merchantId, query, params}, args...)...)}
}

func (_c *MockDataPlatformService_QueryRealTime_Call) Run(run func(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{})) *MockDataPlatformService_QueryRealTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-4)
		for i, a := range args[4:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(map[string]string), variadicArgs...)
	})
	return _c
}

func (_c *MockDataPlatformService_QueryRealTime_Call) Return(_a0 dataplatformmodels.QueryResult, _a1 error) *MockDataPlatformService_QueryRealTime_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_QueryRealTime_Call) RunAndReturn(run func(context.Context, string, string, map[string]string, ...interface{}) (dataplatformmodels.QueryResult, error)) *MockDataPlatformService_QueryRealTime_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterDataset provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) RegisterDataset(ctx context.Context, payload models.RegisterDatasetPayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for RegisterDataset")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterDatasetPayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterDatasetPayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.RegisterDatasetPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_RegisterDataset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterDataset'
type MockDataPlatformService_RegisterDataset_Call struct {
	*mock.Call
}

// RegisterDataset is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.RegisterDatasetPayload
func (_e *MockDataPlatformService_Expecter) RegisterDataset(ctx interface{}, payload interface{}) *MockDataPlatformService_RegisterDataset_Call {
	return &MockDataPlatformService_RegisterDataset_Call{Call: _e.mock.On("RegisterDataset", ctx, payload)}
}

func (_c *MockDataPlatformService_RegisterDataset_Call) Run(run func(ctx context.Context, payload models.RegisterDatasetPayload)) *MockDataPlatformService_RegisterDataset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.RegisterDatasetPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_RegisterDataset_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_RegisterDataset_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_RegisterDataset_Call) RunAndReturn(run func(context.Context, models.RegisterDatasetPayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_RegisterDataset_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterJob provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) RegisterJob(ctx context.Context, payload models.RegisterJobPayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for RegisterJob")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterJobPayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterJobPayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.RegisterJobPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_RegisterJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterJob'
type MockDataPlatformService_RegisterJob_Call struct {
	*mock.Call
}

// RegisterJob is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.RegisterJobPayload
func (_e *MockDataPlatformService_Expecter) RegisterJob(ctx interface{}, payload interface{}) *MockDataPlatformService_RegisterJob_Call {
	return &MockDataPlatformService_RegisterJob_Call{Call: _e.mock.On("RegisterJob", ctx, payload)}
}

func (_c *MockDataPlatformService_RegisterJob_Call) Run(run func(ctx context.Context, payload models.RegisterJobPayload)) *MockDataPlatformService_RegisterJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.RegisterJobPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_RegisterJob_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_RegisterJob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_RegisterJob_Call) RunAndReturn(run func(context.Context, models.RegisterJobPayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_RegisterJob_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAction provides a mock function with given fields: ctx, jobStatusUpdate
func (_m *MockDataPlatformService) UpdateAction(ctx context.Context, jobStatusUpdate models.DatabricksJobStatusUpdatePayload) (actionsmodels.Action, error) {
	ret := _m.Called(ctx, jobStatusUpdate)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAction")
	}

	var r0 actionsmodels.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.DatabricksJobStatusUpdatePayload) (actionsmodels.Action, error)); ok {
		return rf(ctx, jobStatusUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.DatabricksJobStatusUpdatePayload) actionsmodels.Action); ok {
		r0 = rf(ctx, jobStatusUpdate)
	} else {
		r0 = ret.Get(0).(actionsmodels.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.DatabricksJobStatusUpdatePayload) error); ok {
		r1 = rf(ctx, jobStatusUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_UpdateAction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAction'
type MockDataPlatformService_UpdateAction_Call struct {
	*mock.Call
}

// UpdateAction is a helper method to define mock.On call
//   - ctx context.Context
//   - jobStatusUpdate models.DatabricksJobStatusUpdatePayload
func (_e *MockDataPlatformService_Expecter) UpdateAction(ctx interface{}, jobStatusUpdate interface{}) *MockDataPlatformService_UpdateAction_Call {
	return &MockDataPlatformService_UpdateAction_Call{Call: _e.mock.On("UpdateAction", ctx, jobStatusUpdate)}
}

func (_c *MockDataPlatformService_UpdateAction_Call) Run(run func(ctx context.Context, jobStatusUpdate models.DatabricksJobStatusUpdatePayload)) *MockDataPlatformService_UpdateAction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.DatabricksJobStatusUpdatePayload))
	})
	return _c
}

func (_c *MockDataPlatformService_UpdateAction_Call) Return(_a0 actionsmodels.Action, _a1 error) *MockDataPlatformService_UpdateAction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_UpdateAction_Call) RunAndReturn(run func(context.Context, models.DatabricksJobStatusUpdatePayload) (actionsmodels.Action, error)) *MockDataPlatformService_UpdateAction_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDataset provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) UpdateDataset(ctx context.Context, payload models.UpdateDatasetPayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDataset")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UpdateDatasetPayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UpdateDatasetPayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UpdateDatasetPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_UpdateDataset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDataset'
type MockDataPlatformService_UpdateDataset_Call struct {
	*mock.Call
}

// UpdateDataset is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.UpdateDatasetPayload
func (_e *MockDataPlatformService_Expecter) UpdateDataset(ctx interface{}, payload interface{}) *MockDataPlatformService_UpdateDataset_Call {
	return &MockDataPlatformService_UpdateDataset_Call{Call: _e.mock.On("UpdateDataset", ctx, payload)}
}

func (_c *MockDataPlatformService_UpdateDataset_Call) Run(run func(ctx context.Context, payload models.UpdateDatasetPayload)) *MockDataPlatformService_UpdateDataset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UpdateDatasetPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_UpdateDataset_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_UpdateDataset_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_UpdateDataset_Call) RunAndReturn(run func(context.Context, models.UpdateDatasetPayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_UpdateDataset_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDatasetData provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) UpdateDatasetData(ctx context.Context, payload models.UpdateDatasetDataPayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDatasetData")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UpdateDatasetDataPayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UpdateDatasetDataPayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UpdateDatasetDataPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_UpdateDatasetData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDatasetData'
type MockDataPlatformService_UpdateDatasetData_Call struct {
	*mock.Call
}

// UpdateDatasetData is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.UpdateDatasetDataPayload
func (_e *MockDataPlatformService_Expecter) UpdateDatasetData(ctx interface{}, payload interface{}) *MockDataPlatformService_UpdateDatasetData_Call {
	return &MockDataPlatformService_UpdateDatasetData_Call{Call: _e.mock.On("UpdateDatasetData", ctx, payload)}
}

func (_c *MockDataPlatformService_UpdateDatasetData_Call) Run(run func(ctx context.Context, payload models.UpdateDatasetDataPayload)) *MockDataPlatformService_UpdateDatasetData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UpdateDatasetDataPayload))
	})
	return _c
}

func (_c *MockDataPlatformService_UpdateDatasetData_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_UpdateDatasetData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_UpdateDatasetData_Call) RunAndReturn(run func(context.Context, models.UpdateDatasetDataPayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_UpdateDatasetData_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertTemplate provides a mock function with given fields: ctx, payload
func (_m *MockDataPlatformService) UpsertTemplate(ctx context.Context, payload models.UpsertTemplatePayload) (actionsmodels.CreateActionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for UpsertTemplate")
	}

	var r0 actionsmodels.CreateActionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UpsertTemplatePayload) (actionsmodels.CreateActionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UpsertTemplatePayload) actionsmodels.CreateActionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(actionsmodels.CreateActionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UpsertTemplatePayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataPlatformService_UpsertTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertTemplate'
type MockDataPlatformService_UpsertTemplate_Call struct {
	*mock.Call
}

// UpsertTemplate is a helper method to define mock.On call
//   - ctx context.Context
//   - payload models.UpsertTemplatePayload
func (_e *MockDataPlatformService_Expecter) UpsertTemplate(ctx interface{}, payload interface{}) *MockDataPlatformService_UpsertTemplate_Call {
	return &MockDataPlatformService_UpsertTemplate_Call{Call: _e.mock.On("UpsertTemplate", ctx, payload)}
}

func (_c *MockDataPlatformService_UpsertTemplate_Call) Run(run func(ctx context.Context, payload models.UpsertTemplatePayload)) *MockDataPlatformService_UpsertTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UpsertTemplatePayload))
	})
	return _c
}

func (_c *MockDataPlatformService_UpsertTemplate_Call) Return(_a0 actionsmodels.CreateActionResponse, _a1 error) *MockDataPlatformService_UpsertTemplate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataPlatformService_UpsertTemplate_Call) RunAndReturn(run func(context.Context, models.UpsertTemplatePayload) (actionsmodels.CreateActionResponse, error)) *MockDataPlatformService_UpsertTemplate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDataPlatformService creates a new instance of MockDataPlatformService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDataPlatformService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDataPlatformService {
	mock := &MockDataPlatformService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
