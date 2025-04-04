// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_templates

import (
	template "html/template"
	io "io"

	mock "github.com/stretchr/testify/mock"

	templates "github.com/Zampfi/application-platform/services/api/server/routes/admin/templates"
)

// MockTemplateLoader is an autogenerated mock type for the TemplateLoader type
type MockTemplateLoader struct {
	mock.Mock
}

type MockTemplateLoader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTemplateLoader) EXPECT() *MockTemplateLoader_Expecter {
	return &MockTemplateLoader_Expecter{mock: &_m.Mock}
}

// ExecuteTemplate provides a mock function with given fields: w, name, data
func (_m *MockTemplateLoader) ExecuteTemplate(w io.Writer, name string, data templates.TemplateData) error {
	ret := _m.Called(w, name, data)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteTemplate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, string, templates.TemplateData) error); ok {
		r0 = rf(w, name, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTemplateLoader_ExecuteTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteTemplate'
type MockTemplateLoader_ExecuteTemplate_Call struct {
	*mock.Call
}

// ExecuteTemplate is a helper method to define mock.On call
//   - w io.Writer
//   - name string
//   - data templates.TemplateData
func (_e *MockTemplateLoader_Expecter) ExecuteTemplate(w interface{}, name interface{}, data interface{}) *MockTemplateLoader_ExecuteTemplate_Call {
	return &MockTemplateLoader_ExecuteTemplate_Call{Call: _e.mock.On("ExecuteTemplate", w, name, data)}
}

func (_c *MockTemplateLoader_ExecuteTemplate_Call) Run(run func(w io.Writer, name string, data templates.TemplateData)) *MockTemplateLoader_ExecuteTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(io.Writer), args[1].(string), args[2].(templates.TemplateData))
	})
	return _c
}

func (_c *MockTemplateLoader_ExecuteTemplate_Call) Return(_a0 error) *MockTemplateLoader_ExecuteTemplate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemplateLoader_ExecuteTemplate_Call) RunAndReturn(run func(io.Writer, string, templates.TemplateData) error) *MockTemplateLoader_ExecuteTemplate_Call {
	_c.Call.Return(run)
	return _c
}

// LoadTemplate provides a mock function with given fields: templateName, data
func (_m *MockTemplateLoader) LoadTemplate(templateName string, data templates.TemplateData) (*template.Template, error) {
	ret := _m.Called(templateName, data)

	if len(ret) == 0 {
		panic("no return value specified for LoadTemplate")
	}

	var r0 *template.Template
	var r1 error
	if rf, ok := ret.Get(0).(func(string, templates.TemplateData) (*template.Template, error)); ok {
		return rf(templateName, data)
	}
	if rf, ok := ret.Get(0).(func(string, templates.TemplateData) *template.Template); ok {
		r0 = rf(templateName, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*template.Template)
		}
	}

	if rf, ok := ret.Get(1).(func(string, templates.TemplateData) error); ok {
		r1 = rf(templateName, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTemplateLoader_LoadTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadTemplate'
type MockTemplateLoader_LoadTemplate_Call struct {
	*mock.Call
}

// LoadTemplate is a helper method to define mock.On call
//   - templateName string
//   - data templates.TemplateData
func (_e *MockTemplateLoader_Expecter) LoadTemplate(templateName interface{}, data interface{}) *MockTemplateLoader_LoadTemplate_Call {
	return &MockTemplateLoader_LoadTemplate_Call{Call: _e.mock.On("LoadTemplate", templateName, data)}
}

func (_c *MockTemplateLoader_LoadTemplate_Call) Run(run func(templateName string, data templates.TemplateData)) *MockTemplateLoader_LoadTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(templates.TemplateData))
	})
	return _c
}

func (_c *MockTemplateLoader_LoadTemplate_Call) Return(_a0 *template.Template, _a1 error) *MockTemplateLoader_LoadTemplate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTemplateLoader_LoadTemplate_Call) RunAndReturn(run func(string, templates.TemplateData) (*template.Template, error)) *MockTemplateLoader_LoadTemplate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTemplateLoader creates a new instance of MockTemplateLoader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTemplateLoader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTemplateLoader {
	mock := &MockTemplateLoader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
