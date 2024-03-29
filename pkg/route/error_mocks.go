// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package route

import (
	"sync"
)

// Ensure, that ErrorViewMock does implement ErrorView.
// If this is not the case, regenerate this file with moq.
var _ ErrorView = &ErrorViewMock{}

// ErrorViewMock is a mock implementation of ErrorView.
//
//	func TestSomethingThatUsesErrorView(t *testing.T) {
//
//		// make and configure a mocked ErrorView
//		mockedErrorView := &ErrorViewMock{
//			RedirectFunc: func(ctx Context, err error)  {
//				panic("mock out the Redirect method")
//			},
//			RenderFunc: func(ctx Context)  {
//				panic("mock out the Render method")
//			},
//		}
//
//		// use mockedErrorView in code that requires ErrorView
//		// and then make assertions.
//
//	}
type ErrorViewMock struct {
	// RedirectFunc mocks the Redirect method.
	RedirectFunc func(ctx Context, err error)

	// RenderFunc mocks the Render method.
	RenderFunc func(ctx Context)

	// calls tracks calls to the methods.
	calls struct {
		// Redirect holds details about calls to the Redirect method.
		Redirect []struct {
			// Ctx is the ctx argument value.
			Ctx Context
			// Err is the err argument value.
			Err error
		}
		// Render holds details about calls to the Render method.
		Render []struct {
			// Ctx is the ctx argument value.
			Ctx Context
		}
	}
	lockRedirect sync.RWMutex
	lockRender   sync.RWMutex
}

// Redirect calls RedirectFunc.
func (mock *ErrorViewMock) Redirect(ctx Context, err error) {
	if mock.RedirectFunc == nil {
		panic("ErrorViewMock.RedirectFunc: method is nil but ErrorView.Redirect was just called")
	}
	callInfo := struct {
		Ctx Context
		Err error
	}{
		Ctx: ctx,
		Err: err,
	}
	mock.lockRedirect.Lock()
	mock.calls.Redirect = append(mock.calls.Redirect, callInfo)
	mock.lockRedirect.Unlock()
	mock.RedirectFunc(ctx, err)
}

// RedirectCalls gets all the calls that were made to Redirect.
// Check the length with:
//
//	len(mockedErrorView.RedirectCalls())
func (mock *ErrorViewMock) RedirectCalls() []struct {
	Ctx Context
	Err error
} {
	var calls []struct {
		Ctx Context
		Err error
	}
	mock.lockRedirect.RLock()
	calls = mock.calls.Redirect
	mock.lockRedirect.RUnlock()
	return calls
}

// Render calls RenderFunc.
func (mock *ErrorViewMock) Render(ctx Context) {
	if mock.RenderFunc == nil {
		panic("ErrorViewMock.RenderFunc: method is nil but ErrorView.Render was just called")
	}
	callInfo := struct {
		Ctx Context
	}{
		Ctx: ctx,
	}
	mock.lockRender.Lock()
	mock.calls.Render = append(mock.calls.Render, callInfo)
	mock.lockRender.Unlock()
	mock.RenderFunc(ctx)
}

// RenderCalls gets all the calls that were made to Render.
// Check the length with:
//
//	len(mockedErrorView.RenderCalls())
func (mock *ErrorViewMock) RenderCalls() []struct {
	Ctx Context
} {
	var calls []struct {
		Ctx Context
	}
	mock.lockRender.RLock()
	calls = mock.calls.Render
	mock.lockRender.RUnlock()
	return calls
}
