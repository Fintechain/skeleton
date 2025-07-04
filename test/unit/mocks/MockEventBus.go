// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/event"
	mock "github.com/stretchr/testify/mock"
)

// NewMockEventBus creates a new instance of MockEventBus. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventBus(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventBus {
	mock := &MockEventBus{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockEventBus is an autogenerated mock type for the EventBus type
type MockEventBus struct {
	mock.Mock
}

type MockEventBus_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventBus) EXPECT() *MockEventBus_Expecter {
	return &MockEventBus_Expecter{mock: &_m.Mock}
}

// Publish provides a mock function for the type MockEventBus
func (_mock *MockEventBus) Publish(event1 *event.Event) error {
	ret := _mock.Called(event1)

	if len(ret) == 0 {
		panic("no return value specified for Publish")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(*event.Event) error); ok {
		r0 = returnFunc(event1)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockEventBus_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type MockEventBus_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - event1 *event.Event
func (_e *MockEventBus_Expecter) Publish(event1 interface{}) *MockEventBus_Publish_Call {
	return &MockEventBus_Publish_Call{Call: _e.mock.On("Publish", event1)}
}

func (_c *MockEventBus_Publish_Call) Run(run func(event1 *event.Event)) *MockEventBus_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 *event.Event
		if args[0] != nil {
			arg0 = args[0].(*event.Event)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockEventBus_Publish_Call) Return(err error) *MockEventBus_Publish_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockEventBus_Publish_Call) RunAndReturn(run func(event1 *event.Event) error) *MockEventBus_Publish_Call {
	_c.Call.Return(run)
	return _c
}

// PublishAsync provides a mock function for the type MockEventBus
func (_mock *MockEventBus) PublishAsync(event1 *event.Event) error {
	ret := _mock.Called(event1)

	if len(ret) == 0 {
		panic("no return value specified for PublishAsync")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(*event.Event) error); ok {
		r0 = returnFunc(event1)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockEventBus_PublishAsync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishAsync'
type MockEventBus_PublishAsync_Call struct {
	*mock.Call
}

// PublishAsync is a helper method to define mock.On call
//   - event1 *event.Event
func (_e *MockEventBus_Expecter) PublishAsync(event1 interface{}) *MockEventBus_PublishAsync_Call {
	return &MockEventBus_PublishAsync_Call{Call: _e.mock.On("PublishAsync", event1)}
}

func (_c *MockEventBus_PublishAsync_Call) Run(run func(event1 *event.Event)) *MockEventBus_PublishAsync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 *event.Event
		if args[0] != nil {
			arg0 = args[0].(*event.Event)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockEventBus_PublishAsync_Call) Return(err error) *MockEventBus_PublishAsync_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockEventBus_PublishAsync_Call) RunAndReturn(run func(event1 *event.Event) error) *MockEventBus_PublishAsync_Call {
	_c.Call.Return(run)
	return _c
}

// Subscribe provides a mock function for the type MockEventBus
func (_mock *MockEventBus) Subscribe(eventType string, handler event.EventHandler) event.Subscription {
	ret := _mock.Called(eventType, handler)

	if len(ret) == 0 {
		panic("no return value specified for Subscribe")
	}

	var r0 event.Subscription
	if returnFunc, ok := ret.Get(0).(func(string, event.EventHandler) event.Subscription); ok {
		r0 = returnFunc(eventType, handler)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}
	return r0
}

// MockEventBus_Subscribe_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Subscribe'
type MockEventBus_Subscribe_Call struct {
	*mock.Call
}

// Subscribe is a helper method to define mock.On call
//   - eventType string
//   - handler event.EventHandler
func (_e *MockEventBus_Expecter) Subscribe(eventType interface{}, handler interface{}) *MockEventBus_Subscribe_Call {
	return &MockEventBus_Subscribe_Call{Call: _e.mock.On("Subscribe", eventType, handler)}
}

func (_c *MockEventBus_Subscribe_Call) Run(run func(eventType string, handler event.EventHandler)) *MockEventBus_Subscribe_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 string
		if args[0] != nil {
			arg0 = args[0].(string)
		}
		var arg1 event.EventHandler
		if args[1] != nil {
			arg1 = args[1].(event.EventHandler)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockEventBus_Subscribe_Call) Return(subscription event.Subscription) *MockEventBus_Subscribe_Call {
	_c.Call.Return(subscription)
	return _c
}

func (_c *MockEventBus_Subscribe_Call) RunAndReturn(run func(eventType string, handler event.EventHandler) event.Subscription) *MockEventBus_Subscribe_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeAsync provides a mock function for the type MockEventBus
func (_mock *MockEventBus) SubscribeAsync(eventType string, handler event.EventHandler) event.Subscription {
	ret := _mock.Called(eventType, handler)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeAsync")
	}

	var r0 event.Subscription
	if returnFunc, ok := ret.Get(0).(func(string, event.EventHandler) event.Subscription); ok {
		r0 = returnFunc(eventType, handler)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}
	return r0
}

// MockEventBus_SubscribeAsync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeAsync'
type MockEventBus_SubscribeAsync_Call struct {
	*mock.Call
}

// SubscribeAsync is a helper method to define mock.On call
//   - eventType string
//   - handler event.EventHandler
func (_e *MockEventBus_Expecter) SubscribeAsync(eventType interface{}, handler interface{}) *MockEventBus_SubscribeAsync_Call {
	return &MockEventBus_SubscribeAsync_Call{Call: _e.mock.On("SubscribeAsync", eventType, handler)}
}

func (_c *MockEventBus_SubscribeAsync_Call) Run(run func(eventType string, handler event.EventHandler)) *MockEventBus_SubscribeAsync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 string
		if args[0] != nil {
			arg0 = args[0].(string)
		}
		var arg1 event.EventHandler
		if args[1] != nil {
			arg1 = args[1].(event.EventHandler)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockEventBus_SubscribeAsync_Call) Return(subscription event.Subscription) *MockEventBus_SubscribeAsync_Call {
	_c.Call.Return(subscription)
	return _c
}

func (_c *MockEventBus_SubscribeAsync_Call) RunAndReturn(run func(eventType string, handler event.EventHandler) event.Subscription) *MockEventBus_SubscribeAsync_Call {
	_c.Call.Return(run)
	return _c
}

// WaitAsync provides a mock function for the type MockEventBus
func (_mock *MockEventBus) WaitAsync() {
	_mock.Called()
	return
}

// MockEventBus_WaitAsync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WaitAsync'
type MockEventBus_WaitAsync_Call struct {
	*mock.Call
}

// WaitAsync is a helper method to define mock.On call
func (_e *MockEventBus_Expecter) WaitAsync() *MockEventBus_WaitAsync_Call {
	return &MockEventBus_WaitAsync_Call{Call: _e.mock.On("WaitAsync")}
}

func (_c *MockEventBus_WaitAsync_Call) Run(run func()) *MockEventBus_WaitAsync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEventBus_WaitAsync_Call) Return() *MockEventBus_WaitAsync_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEventBus_WaitAsync_Call) RunAndReturn(run func()) *MockEventBus_WaitAsync_Call {
	_c.Run(run)
	return _c
}
