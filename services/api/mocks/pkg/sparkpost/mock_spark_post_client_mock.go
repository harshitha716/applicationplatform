// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_sparkpost

import (
	gosparkpost "github.com/SparkPost/gosparkpost"
	mock "github.com/stretchr/testify/mock"
)

// MockmockSparkPostClient is an autogenerated mock type for the mockSparkPostClient type
type MockmockSparkPostClient struct {
	mock.Mock
}

type MockmockSparkPostClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockmockSparkPostClient) EXPECT() *MockmockSparkPostClient_Expecter {
	return &MockmockSparkPostClient_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: transmission
func (_m *MockmockSparkPostClient) Send(transmission *gosparkpost.Transmission) (string, *gosparkpost.Response, error) {
	ret := _m.Called(transmission)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 string
	var r1 *gosparkpost.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(*gosparkpost.Transmission) (string, *gosparkpost.Response, error)); ok {
		return rf(transmission)
	}
	if rf, ok := ret.Get(0).(func(*gosparkpost.Transmission) string); ok {
		r0 = rf(transmission)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*gosparkpost.Transmission) *gosparkpost.Response); ok {
		r1 = rf(transmission)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gosparkpost.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(*gosparkpost.Transmission) error); ok {
		r2 = rf(transmission)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockmockSparkPostClient_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockmockSparkPostClient_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - transmission *gosparkpost.Transmission
func (_e *MockmockSparkPostClient_Expecter) Send(transmission interface{}) *MockmockSparkPostClient_Send_Call {
	return &MockmockSparkPostClient_Send_Call{Call: _e.mock.On("Send", transmission)}
}

func (_c *MockmockSparkPostClient_Send_Call) Run(run func(transmission *gosparkpost.Transmission)) *MockmockSparkPostClient_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gosparkpost.Transmission))
	})
	return _c
}

func (_c *MockmockSparkPostClient_Send_Call) Return(_a0 string, _a1 *gosparkpost.Response, _a2 error) *MockmockSparkPostClient_Send_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockmockSparkPostClient_Send_Call) RunAndReturn(run func(*gosparkpost.Transmission) (string, *gosparkpost.Response, error)) *MockmockSparkPostClient_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockmockSparkPostClient creates a new instance of MockmockSparkPostClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockmockSparkPostClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockmockSparkPostClient {
	mock := &MockmockSparkPostClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
