// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package auth

import (
	"context"
	"htmx/model"
	"htmx/pkg/util"
	"sync"
)

// Ensure, that RepositoryMock does implement Repository.
// If this is not the case, regenerate this file with moq.
var _ Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of Repository.
//
//	func TestSomethingThatUsesRepository(t *testing.T) {
//
//		// make and configure a mocked Repository
//		mockedRepository := &RepositoryMock{
//			FindUserSessionFunc: func(ctx context.Context, userID model.UserID, sessionID model.SessionID) (util.Null[model.UserSession], error) {
//				panic("mock out the FindUserSession method")
//			},
//			GetUserFunc: func(ctx context.Context, userID model.UserID) (util.Null[model.User], error) {
//				panic("mock out the GetUser method")
//			},
//			InsertUserFunc: func(ctx context.Context, user model.User) (model.UserID, error) {
//				panic("mock out the InsertUser method")
//			},
//			InsertUserSessionFunc: func(ctx context.Context, sess model.UserSession) error {
//				panic("mock out the InsertUserSession method")
//			},
//		}
//
//		// use mockedRepository in code that requires Repository
//		// and then make assertions.
//
//	}
type RepositoryMock struct {
	// FindUserSessionFunc mocks the FindUserSession method.
	FindUserSessionFunc func(ctx context.Context, userID model.UserID, sessionID model.SessionID) (util.Null[model.UserSession], error)

	// GetUserFunc mocks the GetUser method.
	GetUserFunc func(ctx context.Context, userID model.UserID) (util.Null[model.User], error)

	// InsertUserFunc mocks the InsertUser method.
	InsertUserFunc func(ctx context.Context, user model.User) (model.UserID, error)

	// InsertUserSessionFunc mocks the InsertUserSession method.
	InsertUserSessionFunc func(ctx context.Context, sess model.UserSession) error

	// calls tracks calls to the methods.
	calls struct {
		// FindUserSession holds details about calls to the FindUserSession method.
		FindUserSession []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID model.UserID
			// SessionID is the sessionID argument value.
			SessionID model.SessionID
		}
		// GetUser holds details about calls to the GetUser method.
		GetUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID model.UserID
		}
		// InsertUser holds details about calls to the InsertUser method.
		InsertUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// User is the user argument value.
			User model.User
		}
		// InsertUserSession holds details about calls to the InsertUserSession method.
		InsertUserSession []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Sess is the sess argument value.
			Sess model.UserSession
		}
	}
	lockFindUserSession   sync.RWMutex
	lockGetUser           sync.RWMutex
	lockInsertUser        sync.RWMutex
	lockInsertUserSession sync.RWMutex
}

// FindUserSession calls FindUserSessionFunc.
func (mock *RepositoryMock) FindUserSession(ctx context.Context, userID model.UserID, sessionID model.SessionID) (util.Null[model.UserSession], error) {
	if mock.FindUserSessionFunc == nil {
		panic("RepositoryMock.FindUserSessionFunc: method is nil but Repository.FindUserSession was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		UserID    model.UserID
		SessionID model.SessionID
	}{
		Ctx:       ctx,
		UserID:    userID,
		SessionID: sessionID,
	}
	mock.lockFindUserSession.Lock()
	mock.calls.FindUserSession = append(mock.calls.FindUserSession, callInfo)
	mock.lockFindUserSession.Unlock()
	return mock.FindUserSessionFunc(ctx, userID, sessionID)
}

// FindUserSessionCalls gets all the calls that were made to FindUserSession.
// Check the length with:
//
//	len(mockedRepository.FindUserSessionCalls())
func (mock *RepositoryMock) FindUserSessionCalls() []struct {
	Ctx       context.Context
	UserID    model.UserID
	SessionID model.SessionID
} {
	var calls []struct {
		Ctx       context.Context
		UserID    model.UserID
		SessionID model.SessionID
	}
	mock.lockFindUserSession.RLock()
	calls = mock.calls.FindUserSession
	mock.lockFindUserSession.RUnlock()
	return calls
}

// GetUser calls GetUserFunc.
func (mock *RepositoryMock) GetUser(ctx context.Context, userID model.UserID) (util.Null[model.User], error) {
	if mock.GetUserFunc == nil {
		panic("RepositoryMock.GetUserFunc: method is nil but Repository.GetUser was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserID model.UserID
	}{
		Ctx:    ctx,
		UserID: userID,
	}
	mock.lockGetUser.Lock()
	mock.calls.GetUser = append(mock.calls.GetUser, callInfo)
	mock.lockGetUser.Unlock()
	return mock.GetUserFunc(ctx, userID)
}

// GetUserCalls gets all the calls that were made to GetUser.
// Check the length with:
//
//	len(mockedRepository.GetUserCalls())
func (mock *RepositoryMock) GetUserCalls() []struct {
	Ctx    context.Context
	UserID model.UserID
} {
	var calls []struct {
		Ctx    context.Context
		UserID model.UserID
	}
	mock.lockGetUser.RLock()
	calls = mock.calls.GetUser
	mock.lockGetUser.RUnlock()
	return calls
}

// InsertUser calls InsertUserFunc.
func (mock *RepositoryMock) InsertUser(ctx context.Context, user model.User) (model.UserID, error) {
	if mock.InsertUserFunc == nil {
		panic("RepositoryMock.InsertUserFunc: method is nil but Repository.InsertUser was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		User model.User
	}{
		Ctx:  ctx,
		User: user,
	}
	mock.lockInsertUser.Lock()
	mock.calls.InsertUser = append(mock.calls.InsertUser, callInfo)
	mock.lockInsertUser.Unlock()
	return mock.InsertUserFunc(ctx, user)
}

// InsertUserCalls gets all the calls that were made to InsertUser.
// Check the length with:
//
//	len(mockedRepository.InsertUserCalls())
func (mock *RepositoryMock) InsertUserCalls() []struct {
	Ctx  context.Context
	User model.User
} {
	var calls []struct {
		Ctx  context.Context
		User model.User
	}
	mock.lockInsertUser.RLock()
	calls = mock.calls.InsertUser
	mock.lockInsertUser.RUnlock()
	return calls
}

// InsertUserSession calls InsertUserSessionFunc.
func (mock *RepositoryMock) InsertUserSession(ctx context.Context, sess model.UserSession) error {
	if mock.InsertUserSessionFunc == nil {
		panic("RepositoryMock.InsertUserSessionFunc: method is nil but Repository.InsertUserSession was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Sess model.UserSession
	}{
		Ctx:  ctx,
		Sess: sess,
	}
	mock.lockInsertUserSession.Lock()
	mock.calls.InsertUserSession = append(mock.calls.InsertUserSession, callInfo)
	mock.lockInsertUserSession.Unlock()
	return mock.InsertUserSessionFunc(ctx, sess)
}

// InsertUserSessionCalls gets all the calls that were made to InsertUserSession.
// Check the length with:
//
//	len(mockedRepository.InsertUserSessionCalls())
func (mock *RepositoryMock) InsertUserSessionCalls() []struct {
	Ctx  context.Context
	Sess model.UserSession
} {
	var calls []struct {
		Ctx  context.Context
		Sess model.UserSession
	}
	mock.lockInsertUserSession.RLock()
	calls = mock.calls.InsertUserSession
	mock.lockInsertUserSession.RUnlock()
	return calls
}

// Ensure, that RandServiceMock does implement RandService.
// If this is not the case, regenerate this file with moq.
var _ RandService = &RandServiceMock{}

// RandServiceMock is a mock implementation of RandService.
//
//	func TestSomethingThatUsesRandService(t *testing.T) {
//
//		// make and configure a mocked RandService
//		mockedRandService := &RandServiceMock{
//			RandStringFunc: func(size int) (string, error) {
//				panic("mock out the RandString method")
//			},
//		}
//
//		// use mockedRandService in code that requires RandService
//		// and then make assertions.
//
//	}
type RandServiceMock struct {
	// RandStringFunc mocks the RandString method.
	RandStringFunc func(size int) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// RandString holds details about calls to the RandString method.
		RandString []struct {
			// Size is the size argument value.
			Size int
		}
	}
	lockRandString sync.RWMutex
}

// RandString calls RandStringFunc.
func (mock *RandServiceMock) RandString(size int) (string, error) {
	if mock.RandStringFunc == nil {
		panic("RandServiceMock.RandStringFunc: method is nil but RandService.RandString was just called")
	}
	callInfo := struct {
		Size int
	}{
		Size: size,
	}
	mock.lockRandString.Lock()
	mock.calls.RandString = append(mock.calls.RandString, callInfo)
	mock.lockRandString.Unlock()
	return mock.RandStringFunc(size)
}

// RandStringCalls gets all the calls that were made to RandString.
// Check the length with:
//
//	len(mockedRandService.RandStringCalls())
func (mock *RandServiceMock) RandStringCalls() []struct {
	Size int
} {
	var calls []struct {
		Size int
	}
	mock.lockRandString.RLock()
	calls = mock.calls.RandString
	mock.lockRandString.RUnlock()
	return calls
}
