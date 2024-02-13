// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package dbtx

import (
	"context"
	"sync"
)

// Ensure, that ProviderMock does implement Provider.
// If this is not the case, regenerate this file with moq.
var _ Provider = &ProviderMock{}

// ProviderMock is a mock implementation of Provider.
//
//	func TestSomethingThatUsesProvider(t *testing.T) {
//
//		// make and configure a mocked Provider
//		mockedProvider := &ProviderMock{
//			AutocommitFunc: func(ctx context.Context) context.Context {
//				panic("mock out the Autocommit method")
//			},
//			ReadonlyFunc: func(ctx context.Context) context.Context {
//				panic("mock out the Readonly method")
//			},
//			TransactFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
//				panic("mock out the Transact method")
//			},
//		}
//
//		// use mockedProvider in code that requires Provider
//		// and then make assertions.
//
//	}
type ProviderMock struct {
	// AutocommitFunc mocks the Autocommit method.
	AutocommitFunc func(ctx context.Context) context.Context

	// ReadonlyFunc mocks the Readonly method.
	ReadonlyFunc func(ctx context.Context) context.Context

	// TransactFunc mocks the Transact method.
	TransactFunc func(ctx context.Context, fn func(ctx context.Context) error) error

	// calls tracks calls to the methods.
	calls struct {
		// Autocommit holds details about calls to the Autocommit method.
		Autocommit []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Readonly holds details about calls to the Readonly method.
		Readonly []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Transact holds details about calls to the Transact method.
		Transact []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Fn is the fn argument value.
			Fn func(ctx context.Context) error
		}
	}
	lockAutocommit sync.RWMutex
	lockReadonly   sync.RWMutex
	lockTransact   sync.RWMutex
}

// Autocommit calls AutocommitFunc.
func (mock *ProviderMock) Autocommit(ctx context.Context) context.Context {
	if mock.AutocommitFunc == nil {
		panic("ProviderMock.AutocommitFunc: method is nil but Provider.Autocommit was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockAutocommit.Lock()
	mock.calls.Autocommit = append(mock.calls.Autocommit, callInfo)
	mock.lockAutocommit.Unlock()
	return mock.AutocommitFunc(ctx)
}

// AutocommitCalls gets all the calls that were made to Autocommit.
// Check the length with:
//
//	len(mockedProvider.AutocommitCalls())
func (mock *ProviderMock) AutocommitCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockAutocommit.RLock()
	calls = mock.calls.Autocommit
	mock.lockAutocommit.RUnlock()
	return calls
}

// Readonly calls ReadonlyFunc.
func (mock *ProviderMock) Readonly(ctx context.Context) context.Context {
	if mock.ReadonlyFunc == nil {
		panic("ProviderMock.ReadonlyFunc: method is nil but Provider.Readonly was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockReadonly.Lock()
	mock.calls.Readonly = append(mock.calls.Readonly, callInfo)
	mock.lockReadonly.Unlock()
	return mock.ReadonlyFunc(ctx)
}

// ReadonlyCalls gets all the calls that were made to Readonly.
// Check the length with:
//
//	len(mockedProvider.ReadonlyCalls())
func (mock *ProviderMock) ReadonlyCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockReadonly.RLock()
	calls = mock.calls.Readonly
	mock.lockReadonly.RUnlock()
	return calls
}

// Transact calls TransactFunc.
func (mock *ProviderMock) Transact(ctx context.Context, fn func(ctx context.Context) error) error {
	if mock.TransactFunc == nil {
		panic("ProviderMock.TransactFunc: method is nil but Provider.Transact was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Fn  func(ctx context.Context) error
	}{
		Ctx: ctx,
		Fn:  fn,
	}
	mock.lockTransact.Lock()
	mock.calls.Transact = append(mock.calls.Transact, callInfo)
	mock.lockTransact.Unlock()
	return mock.TransactFunc(ctx, fn)
}

// TransactCalls gets all the calls that were made to Transact.
// Check the length with:
//
//	len(mockedProvider.TransactCalls())
func (mock *ProviderMock) TransactCalls() []struct {
	Ctx context.Context
	Fn  func(ctx context.Context) error
} {
	var calls []struct {
		Ctx context.Context
		Fn  func(ctx context.Context) error
	}
	mock.lockTransact.RLock()
	calls = mock.calls.Transact
	mock.lockTransact.RUnlock()
	return calls
}
