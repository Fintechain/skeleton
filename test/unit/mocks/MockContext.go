// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/context"
	mock "github.com/stretchr/testify/mock"
)

// NewMockContext creates a new instance of MockContext. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockContext(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockContext {
	mock := &MockContext{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockContext is an autogenerated mock type for the Context type
type MockContext struct {
	mock.Mock
}

type MockContext_Expecter struct {
	mock *mock.Mock
}

func (_m *MockContext) EXPECT() *MockContext_Expecter {
	return &MockContext_Expecter{mock: &_m.Mock}
}

// Deadline provides a mock function for the type MockContext
func (_mock *MockContext) Deadline() (time.Time, bool) {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Deadline")
	}

	var r0 time.Time
	var r1 bool
	if returnFunc, ok := ret.Get(0).(func() (time.Time, bool)); ok {
		return returnFunc()
	}
	if returnFunc, ok := ret.Get(0).(func() time.Time); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(time.Time)
	}
	if returnFunc, ok := ret.Get(1).(func() bool); ok {
		r1 = returnFunc()
	} else {
		r1 = ret.Get(1).(bool)
	}
	return r0, r1
}

// MockContext_Deadline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deadline'
type MockContext_Deadline_Call struct {
	*mock.Call
}

// Deadline is a helper method to define mock.On call
func (_e *MockContext_Expecter) Deadline() *MockContext_Deadline_Call {
	return &MockContext_Deadline_Call{Call: _e.mock.On("Deadline")}
}

func (_c *MockContext_Deadline_Call) Run(run func()) *MockContext_Deadline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockContext_Deadline_Call) Return(time1 time.Time, b bool) *MockContext_Deadline_Call {
	_c.Call.Return(time1, b)
	return _c
}

func (_c *MockContext_Deadline_Call) RunAndReturn(run func() (time.Time, bool)) *MockContext_Deadline_Call {
	_c.Call.Return(run)
	return _c
}

// Done provides a mock function for the type MockContext
func (_mock *MockContext) Done() <-chan struct{} {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Done")
	}

	var r0 <-chan struct{}
	if returnFunc, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}
	return r0
}

// MockContext_Done_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Done'
type MockContext_Done_Call struct {
	*mock.Call
}

// Done is a helper method to define mock.On call
func (_e *MockContext_Expecter) Done() *MockContext_Done_Call {
	return &MockContext_Done_Call{Call: _e.mock.On("Done")}
}

func (_c *MockContext_Done_Call) Run(run func()) *MockContext_Done_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockContext_Done_Call) Return(valCh <-chan struct{}) *MockContext_Done_Call {
	_c.Call.Return(valCh)
	return _c
}

func (_c *MockContext_Done_Call) RunAndReturn(run func() <-chan struct{}) *MockContext_Done_Call {
	_c.Call.Return(run)
	return _c
}

// Err provides a mock function for the type MockContext
func (_mock *MockContext) Err() error {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Err")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func() error); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockContext_Err_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Err'
type MockContext_Err_Call struct {
	*mock.Call
}

// Err is a helper method to define mock.On call
func (_e *MockContext_Expecter) Err() *MockContext_Err_Call {
	return &MockContext_Err_Call{Call: _e.mock.On("Err")}
}

func (_c *MockContext_Err_Call) Run(run func()) *MockContext_Err_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockContext_Err_Call) Return(err error) *MockContext_Err_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockContext_Err_Call) RunAndReturn(run func() error) *MockContext_Err_Call {
	_c.Call.Return(run)
	return _c
}

// Value provides a mock function for the type MockContext
func (_mock *MockContext) Value(key interface{}) interface{} {
	ret := _mock.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Value")
	}

	var r0 interface{}
	if returnFunc, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = returnFunc(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}
	return r0
}

// MockContext_Value_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Value'
type MockContext_Value_Call struct {
	*mock.Call
}

// Value is a helper method to define mock.On call
//   - key interface{}
func (_e *MockContext_Expecter) Value(key interface{}) *MockContext_Value_Call {
	return &MockContext_Value_Call{Call: _e.mock.On("Value", key)}
}

func (_c *MockContext_Value_Call) Run(run func(key interface{})) *MockContext_Value_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 interface{}
		if args[0] != nil {
			arg0 = args[0].(interface{})
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockContext_Value_Call) Return(ifaceVal interface{}) *MockContext_Value_Call {
	_c.Call.Return(ifaceVal)
	return _c
}

func (_c *MockContext_Value_Call) RunAndReturn(run func(key interface{}) interface{}) *MockContext_Value_Call {
	_c.Call.Return(run)
	return _c
}

// WithValue provides a mock function for the type MockContext
func (_mock *MockContext) WithValue(key interface{}, value interface{}) context.Context {
	ret := _mock.Called(key, value)

	if len(ret) == 0 {
		panic("no return value specified for WithValue")
	}

	var r0 context.Context
	if returnFunc, ok := ret.Get(0).(func(interface{}, interface{}) context.Context); ok {
		r0 = returnFunc(key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}
	return r0
}

// MockContext_WithValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithValue'
type MockContext_WithValue_Call struct {
	*mock.Call
}

// WithValue is a helper method to define mock.On call
//   - key interface{}
//   - value interface{}
func (_e *MockContext_Expecter) WithValue(key interface{}, value interface{}) *MockContext_WithValue_Call {
	return &MockContext_WithValue_Call{Call: _e.mock.On("WithValue", key, value)}
}

func (_c *MockContext_WithValue_Call) Run(run func(key interface{}, value interface{})) *MockContext_WithValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 interface{}
		if args[0] != nil {
			arg0 = args[0].(interface{})
		}
		var arg1 interface{}
		if args[1] != nil {
			arg1 = args[1].(interface{})
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockContext_WithValue_Call) Return(context1 context.Context) *MockContext_WithValue_Call {
	_c.Call.Return(context1)
	return _c
}

func (_c *MockContext_WithValue_Call) RunAndReturn(run func(key interface{}, value interface{}) context.Context) *MockContext_WithValue_Call {
	_c.Call.Return(run)
	return _c
}
