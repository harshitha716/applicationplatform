// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_databricks

import (
	context "context"

	jobs "github.com/databricks/databricks-sdk-go/service/jobs"

	mock "github.com/stretchr/testify/mock"
)

// MockDatabricksSDKProxy is an autogenerated mock type for the DatabricksSDKProxy type
type MockDatabricksSDKProxy struct {
	mock.Mock
}

type MockDatabricksSDKProxy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabricksSDKProxy) EXPECT() *MockDatabricksSDKProxy_Expecter {
	return &MockDatabricksSDKProxy_Expecter{mock: &_m.Mock}
}

// CreateJob provides a mock function with given fields: ctx, jobConfig
func (_m *MockDatabricksSDKProxy) CreateJob(ctx context.Context, jobConfig *jobs.CreateJob) (*jobs.CreateResponse, error) {
	ret := _m.Called(ctx, jobConfig)

	if len(ret) == 0 {
		panic("no return value specified for CreateJob")
	}

	var r0 *jobs.CreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *jobs.CreateJob) (*jobs.CreateResponse, error)); ok {
		return rf(ctx, jobConfig)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *jobs.CreateJob) *jobs.CreateResponse); ok {
		r0 = rf(ctx, jobConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jobs.CreateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *jobs.CreateJob) error); ok {
		r1 = rf(ctx, jobConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabricksSDKProxy_CreateJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateJob'
type MockDatabricksSDKProxy_CreateJob_Call struct {
	*mock.Call
}

// CreateJob is a helper method to define mock.On call
//   - ctx context.Context
//   - jobConfig *jobs.CreateJob
func (_e *MockDatabricksSDKProxy_Expecter) CreateJob(ctx interface{}, jobConfig interface{}) *MockDatabricksSDKProxy_CreateJob_Call {
	return &MockDatabricksSDKProxy_CreateJob_Call{Call: _e.mock.On("CreateJob", ctx, jobConfig)}
}

func (_c *MockDatabricksSDKProxy_CreateJob_Call) Run(run func(ctx context.Context, jobConfig *jobs.CreateJob)) *MockDatabricksSDKProxy_CreateJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*jobs.CreateJob))
	})
	return _c
}

func (_c *MockDatabricksSDKProxy_CreateJob_Call) Return(_a0 *jobs.CreateResponse, _a1 error) *MockDatabricksSDKProxy_CreateJob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabricksSDKProxy_CreateJob_Call) RunAndReturn(run func(context.Context, *jobs.CreateJob) (*jobs.CreateResponse, error)) *MockDatabricksSDKProxy_CreateJob_Call {
	_c.Call.Return(run)
	return _c
}

// GetRunDetails provides a mock function with given fields: ctx, runId
func (_m *MockDatabricksSDKProxy) GetRunDetails(ctx context.Context, runId int64) (*jobs.Run, error) {
	ret := _m.Called(ctx, runId)

	if len(ret) == 0 {
		panic("no return value specified for GetRunDetails")
	}

	var r0 *jobs.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*jobs.Run, error)); ok {
		return rf(ctx, runId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *jobs.Run); ok {
		r0 = rf(ctx, runId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jobs.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, runId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabricksSDKProxy_GetRunDetails_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRunDetails'
type MockDatabricksSDKProxy_GetRunDetails_Call struct {
	*mock.Call
}

// GetRunDetails is a helper method to define mock.On call
//   - ctx context.Context
//   - runId int64
func (_e *MockDatabricksSDKProxy_Expecter) GetRunDetails(ctx interface{}, runId interface{}) *MockDatabricksSDKProxy_GetRunDetails_Call {
	return &MockDatabricksSDKProxy_GetRunDetails_Call{Call: _e.mock.On("GetRunDetails", ctx, runId)}
}

func (_c *MockDatabricksSDKProxy_GetRunDetails_Call) Run(run func(ctx context.Context, runId int64)) *MockDatabricksSDKProxy_GetRunDetails_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockDatabricksSDKProxy_GetRunDetails_Call) Return(_a0 *jobs.Run, _a1 error) *MockDatabricksSDKProxy_GetRunDetails_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabricksSDKProxy_GetRunDetails_Call) RunAndReturn(run func(context.Context, int64) (*jobs.Run, error)) *MockDatabricksSDKProxy_GetRunDetails_Call {
	_c.Call.Return(run)
	return _c
}

// RunNow provides a mock function with given fields: ctx, params
func (_m *MockDatabricksSDKProxy) RunNow(ctx context.Context, params jobs.RunNow) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse], error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for RunNow")
	}

	var r0 *jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, jobs.RunNow) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse], error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, jobs.RunNow) *jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse]); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, jobs.RunNow) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabricksSDKProxy_RunNow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunNow'
type MockDatabricksSDKProxy_RunNow_Call struct {
	*mock.Call
}

// RunNow is a helper method to define mock.On call
//   - ctx context.Context
//   - params jobs.RunNow
func (_e *MockDatabricksSDKProxy_Expecter) RunNow(ctx interface{}, params interface{}) *MockDatabricksSDKProxy_RunNow_Call {
	return &MockDatabricksSDKProxy_RunNow_Call{Call: _e.mock.On("RunNow", ctx, params)}
}

func (_c *MockDatabricksSDKProxy_RunNow_Call) Run(run func(ctx context.Context, params jobs.RunNow)) *MockDatabricksSDKProxy_RunNow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(jobs.RunNow))
	})
	return _c
}

func (_c *MockDatabricksSDKProxy_RunNow_Call) Return(_a0 *jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse], _a1 error) *MockDatabricksSDKProxy_RunNow_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabricksSDKProxy_RunNow_Call) RunAndReturn(run func(context.Context, jobs.RunNow) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.RunNowResponse], error)) *MockDatabricksSDKProxy_RunNow_Call {
	_c.Call.Return(run)
	return _c
}

// SubmitOneTimeJob provides a mock function with given fields: ctx, jobConfig
func (_m *MockDatabricksSDKProxy) SubmitOneTimeJob(ctx context.Context, jobConfig *jobs.SubmitRun) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse], error) {
	ret := _m.Called(ctx, jobConfig)

	if len(ret) == 0 {
		panic("no return value specified for SubmitOneTimeJob")
	}

	var r0 *jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *jobs.SubmitRun) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse], error)); ok {
		return rf(ctx, jobConfig)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *jobs.SubmitRun) *jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse]); ok {
		r0 = rf(ctx, jobConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *jobs.SubmitRun) error); ok {
		r1 = rf(ctx, jobConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabricksSDKProxy_SubmitOneTimeJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubmitOneTimeJob'
type MockDatabricksSDKProxy_SubmitOneTimeJob_Call struct {
	*mock.Call
}

// SubmitOneTimeJob is a helper method to define mock.On call
//   - ctx context.Context
//   - jobConfig *jobs.SubmitRun
func (_e *MockDatabricksSDKProxy_Expecter) SubmitOneTimeJob(ctx interface{}, jobConfig interface{}) *MockDatabricksSDKProxy_SubmitOneTimeJob_Call {
	return &MockDatabricksSDKProxy_SubmitOneTimeJob_Call{Call: _e.mock.On("SubmitOneTimeJob", ctx, jobConfig)}
}

func (_c *MockDatabricksSDKProxy_SubmitOneTimeJob_Call) Run(run func(ctx context.Context, jobConfig *jobs.SubmitRun)) *MockDatabricksSDKProxy_SubmitOneTimeJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*jobs.SubmitRun))
	})
	return _c
}

func (_c *MockDatabricksSDKProxy_SubmitOneTimeJob_Call) Return(_a0 *jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse], _a1 error) *MockDatabricksSDKProxy_SubmitOneTimeJob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabricksSDKProxy_SubmitOneTimeJob_Call) RunAndReturn(run func(context.Context, *jobs.SubmitRun) (*jobs.WaitGetRunJobTerminatedOrSkipped[jobs.SubmitRunResponse], error)) *MockDatabricksSDKProxy_SubmitOneTimeJob_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatabricksSDKProxy creates a new instance of MockDatabricksSDKProxy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabricksSDKProxy(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabricksSDKProxy {
	mock := &MockDatabricksSDKProxy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
