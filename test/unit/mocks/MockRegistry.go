// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	mock "github.com/stretchr/testify/mock"
)

// NewMockRegistry creates a new instance of MockRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRegistry {
	mock := &MockRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockRegistry is an autogenerated mock type for the Registry type
type MockRegistry struct {
	mock.Mock
}

type MockRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegistry) EXPECT() *MockRegistry_Expecter {
	return &MockRegistry_Expecter{mock: &_m.Mock}
}

// Clear provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Clear() error {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Clear")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func() error); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRegistry_Clear_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clear'
type MockRegistry_Clear_Call struct {
	*mock.Call
}

// Clear is a helper method to define mock.On call
func (_e *MockRegistry_Expecter) Clear() *MockRegistry_Clear_Call {
	return &MockRegistry_Clear_Call{Call: _e.mock.On("Clear")}
}

func (_c *MockRegistry_Clear_Call) Run(run func()) *MockRegistry_Clear_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRegistry_Clear_Call) Return(err error) *MockRegistry_Clear_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRegistry_Clear_Call) RunAndReturn(run func() error) *MockRegistry_Clear_Call {
	_c.Call.Return(run)
	return _c
}

// Count provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Count() int {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Count")
	}

	var r0 int
	if returnFunc, ok := ret.Get(0).(func() int); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(int)
	}
	return r0
}

// MockRegistry_Count_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Count'
type MockRegistry_Count_Call struct {
	*mock.Call
}

// Count is a helper method to define mock.On call
func (_e *MockRegistry_Expecter) Count() *MockRegistry_Count_Call {
	return &MockRegistry_Count_Call{Call: _e.mock.On("Count")}
}

func (_c *MockRegistry_Count_Call) Run(run func()) *MockRegistry_Count_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRegistry_Count_Call) Return(n int) *MockRegistry_Count_Call {
	_c.Call.Return(n)
	return _c
}

func (_c *MockRegistry_Count_Call) RunAndReturn(run func() int) *MockRegistry_Count_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Find(predicate func(component.Component) bool) ([]component.Component, error) {
	ret := _mock.Called(predicate)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 []component.Component
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(func(component.Component) bool) ([]component.Component, error)); ok {
		return returnFunc(predicate)
	}
	if returnFunc, ok := ret.Get(0).(func(func(component.Component) bool) []component.Component); ok {
		r0 = returnFunc(predicate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.Component)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(func(component.Component) bool) error); ok {
		r1 = returnFunc(predicate)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRegistry_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockRegistry_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - predicate func(component.Component) bool
func (_e *MockRegistry_Expecter) Find(predicate interface{}) *MockRegistry_Find_Call {
	return &MockRegistry_Find_Call{Call: _e.mock.On("Find", predicate)}
}

func (_c *MockRegistry_Find_Call) Run(run func(predicate func(component.Component) bool)) *MockRegistry_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 func(component.Component) bool
		if args[0] != nil {
			arg0 = args[0].(func(component.Component) bool)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRegistry_Find_Call) Return(components []component.Component, err error) *MockRegistry_Find_Call {
	_c.Call.Return(components, err)
	return _c
}

func (_c *MockRegistry_Find_Call) RunAndReturn(run func(predicate func(component.Component) bool) ([]component.Component, error)) *MockRegistry_Find_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Get(id component.ComponentID) (component.Component, error) {
	ret := _mock.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 component.Component
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(component.ComponentID) (component.Component, error)); ok {
		return returnFunc(id)
	}
	if returnFunc, ok := ret.Get(0).(func(component.ComponentID) component.Component); ok {
		r0 = returnFunc(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(component.Component)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(component.ComponentID) error); ok {
		r1 = returnFunc(id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRegistry_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockRegistry_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - id component.ComponentID
func (_e *MockRegistry_Expecter) Get(id interface{}) *MockRegistry_Get_Call {
	return &MockRegistry_Get_Call{Call: _e.mock.On("Get", id)}
}

func (_c *MockRegistry_Get_Call) Run(run func(id component.ComponentID)) *MockRegistry_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 component.ComponentID
		if args[0] != nil {
			arg0 = args[0].(component.ComponentID)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRegistry_Get_Call) Return(component1 component.Component, err error) *MockRegistry_Get_Call {
	_c.Call.Return(component1, err)
	return _c
}

func (_c *MockRegistry_Get_Call) RunAndReturn(run func(id component.ComponentID) (component.Component, error)) *MockRegistry_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByType provides a mock function for the type MockRegistry
func (_mock *MockRegistry) GetByType(typ component.ComponentType) ([]component.Component, error) {
	ret := _mock.Called(typ)

	if len(ret) == 0 {
		panic("no return value specified for GetByType")
	}

	var r0 []component.Component
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(component.ComponentType) ([]component.Component, error)); ok {
		return returnFunc(typ)
	}
	if returnFunc, ok := ret.Get(0).(func(component.ComponentType) []component.Component); ok {
		r0 = returnFunc(typ)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.Component)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(component.ComponentType) error); ok {
		r1 = returnFunc(typ)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRegistry_GetByType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByType'
type MockRegistry_GetByType_Call struct {
	*mock.Call
}

// GetByType is a helper method to define mock.On call
//   - typ component.ComponentType
func (_e *MockRegistry_Expecter) GetByType(typ interface{}) *MockRegistry_GetByType_Call {
	return &MockRegistry_GetByType_Call{Call: _e.mock.On("GetByType", typ)}
}

func (_c *MockRegistry_GetByType_Call) Run(run func(typ component.ComponentType)) *MockRegistry_GetByType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 component.ComponentType
		if args[0] != nil {
			arg0 = args[0].(component.ComponentType)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRegistry_GetByType_Call) Return(components []component.Component, err error) *MockRegistry_GetByType_Call {
	_c.Call.Return(components, err)
	return _c
}

func (_c *MockRegistry_GetByType_Call) RunAndReturn(run func(typ component.ComponentType) ([]component.Component, error)) *MockRegistry_GetByType_Call {
	_c.Call.Return(run)
	return _c
}

// Has provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Has(id component.ComponentID) bool {
	ret := _mock.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Has")
	}

	var r0 bool
	if returnFunc, ok := ret.Get(0).(func(component.ComponentID) bool); ok {
		r0 = returnFunc(id)
	} else {
		r0 = ret.Get(0).(bool)
	}
	return r0
}

// MockRegistry_Has_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Has'
type MockRegistry_Has_Call struct {
	*mock.Call
}

// Has is a helper method to define mock.On call
//   - id component.ComponentID
func (_e *MockRegistry_Expecter) Has(id interface{}) *MockRegistry_Has_Call {
	return &MockRegistry_Has_Call{Call: _e.mock.On("Has", id)}
}

func (_c *MockRegistry_Has_Call) Run(run func(id component.ComponentID)) *MockRegistry_Has_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 component.ComponentID
		if args[0] != nil {
			arg0 = args[0].(component.ComponentID)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRegistry_Has_Call) Return(b bool) *MockRegistry_Has_Call {
	_c.Call.Return(b)
	return _c
}

func (_c *MockRegistry_Has_Call) RunAndReturn(run func(id component.ComponentID) bool) *MockRegistry_Has_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function for the type MockRegistry
func (_mock *MockRegistry) List() []component.ComponentID {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []component.ComponentID
	if returnFunc, ok := ret.Get(0).(func() []component.ComponentID); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.ComponentID)
		}
	}
	return r0
}

// MockRegistry_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockRegistry_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
func (_e *MockRegistry_Expecter) List() *MockRegistry_List_Call {
	return &MockRegistry_List_Call{Call: _e.mock.On("List")}
}

func (_c *MockRegistry_List_Call) Run(run func()) *MockRegistry_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRegistry_List_Call) Return(componentIDs []component.ComponentID) *MockRegistry_List_Call {
	_c.Call.Return(componentIDs)
	return _c
}

func (_c *MockRegistry_List_Call) RunAndReturn(run func() []component.ComponentID) *MockRegistry_List_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Register(component1 component.Component) error {
	ret := _mock.Called(component1)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(component.Component) error); ok {
		r0 = returnFunc(component1)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRegistry_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type MockRegistry_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - component1 component.Component
func (_e *MockRegistry_Expecter) Register(component1 interface{}) *MockRegistry_Register_Call {
	return &MockRegistry_Register_Call{Call: _e.mock.On("Register", component1)}
}

func (_c *MockRegistry_Register_Call) Run(run func(component1 component.Component)) *MockRegistry_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 component.Component
		if args[0] != nil {
			arg0 = args[0].(component.Component)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRegistry_Register_Call) Return(err error) *MockRegistry_Register_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRegistry_Register_Call) RunAndReturn(run func(component1 component.Component) error) *MockRegistry_Register_Call {
	_c.Call.Return(run)
	return _c
}

// Unregister provides a mock function for the type MockRegistry
func (_mock *MockRegistry) Unregister(id component.ComponentID) error {
	ret := _mock.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Unregister")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(component.ComponentID) error); ok {
		r0 = returnFunc(id)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRegistry_Unregister_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unregister'
type MockRegistry_Unregister_Call struct {
	*mock.Call
}

// Unregister is a helper method to define mock.On call
//   - id component.ComponentID
func (_e *MockRegistry_Expecter) Unregister(id interface{}) *MockRegistry_Unregister_Call {
	return &MockRegistry_Unregister_Call{Call: _e.mock.On("Unregister", id)}
}

func (_c *MockRegistry_Unregister_Call) Run(run func(id component.ComponentID)) *MockRegistry_Unregister_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 component.ComponentID
		if args[0] != nil {
			arg0 = args[0].(component.ComponentID)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRegistry_Unregister_Call) Return(err error) *MockRegistry_Unregister_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRegistry_Unregister_Call) RunAndReturn(run func(id component.ComponentID) error) *MockRegistry_Unregister_Call {
	_c.Call.Return(run)
	return _c
}
