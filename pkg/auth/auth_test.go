package auth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"htmx/config"
	"htmx/model"
	"htmx/pkg/route"
	"htmx/pkg/testhelper"
	"htmx/pkg/util"
)

type serviceTest struct {
	ht   *testhelper.HTTPTest
	repo *RepositoryMock
	rand *RandServiceMock
	svc  Service

	errorView *route.ErrorViewMock
}

func newServiceTest() *serviceTest {
	s := &serviceTest{
		ht:   testhelper.NewHTTPTest(),
		repo: &RepositoryMock{},
		rand: &RandServiceMock{},

		errorView: &route.ErrorViewMock{},
	}

	s.svc = NewService(
		config.Auth{
			CSRFHMACSecret: "some-secret",
		},
		s.repo, s.rand, s.errorView,
	)
	s.rand.RandStringFunc = func(size int) (string, error) {
		return "random-string", nil
	}

	s.stubFindSess(model.NullUserSession{})

	return s
}

func (s *serviceTest) stubRand(values ...string) {
	s.rand.RandStringFunc = func(size int) (string, error) {
		index := len(s.rand.RandStringCalls()) - 1
		return values[index], nil
	}
}

func (s *serviceTest) stubFindSess(sess model.NullUserSession) {
	s.repo.FindUserSessionFunc = func(
		ctx context.Context, userID model.UserID, sessionID model.SessionID,
	) (util.Null[model.UserSession], error) {
		return sess, nil
	}
}

func TestService(t *testing.T) {
	t.Run("req no cookie, build pre session", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.stubRand(
			"some-session-id",
			"abcd",
		)

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, true, continuing)
		assert.Equal(t, nil, err)

		token := "cHJlOnNvbWUtc2Vzc2lvbi1pZCFhYmNk7H8KcrNzLLI07eG4eyWCAzHv75y8nWyihL0Tij11wZo="
		assert.Equal(t, http.Header{
			"Set-Cookie": []string{
				"session_id=pre:some-session-id; Max-Age=2592000; HttpOnly",
				fmt.Sprintf("csrf_token=%s!abcd; Max-Age=2592000", token),
			},
		}, s.ht.Writer.Header())

		calls := s.rand.RandStringCalls()
		assert.Equal(t, 2, len(calls))
		assert.Equal(t, 32, calls[0].Size)
		assert.Equal(t, 16, calls[1].Size)
	})

	t.Run("req already has pre session", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=pre:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, true, continuing)
		assert.Equal(t, nil, err)

		assert.Equal(t, http.Header{}, s.ht.Writer.Header())
	})

	t.Run("pre session is missing parts", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewPost("/users", "")
		s.ht.Req.Header.Add("Cookie", "session_id=pre; Max-Age=2592000; HttpOnly; SameSite=Strict")

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		headers := s.ht.Writer.Header()
		_, ok := headers["Set-Cookie"]
		assert.Equal(t, true, ok)
	})

	t.Run("req already has session, call find user session not found", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")

		s.stubRand(
			"some-session-id",
			"abcd",
		)

		s.stubFindSess(model.NullUserSession{})

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		calls := s.repo.FindUserSessionCalls()
		assert.Equal(t, 1, len(calls))
		assert.Equal(t, model.UserID(1234), calls[0].UserID)
		assert.Equal(t, model.SessionID("some-session-id"), calls[0].SessionID)

		token := "cHJlOnNvbWUtc2Vzc2lvbi1pZCFhYmNk7H8KcrNzLLI07eG4eyWCAzHv75y8nWyihL0Tij11wZo="
		assert.Equal(t, http.Header{
			"Set-Cookie": {
				"session_id=pre:some-session-id; Max-Age=2592000; HttpOnly",
				fmt.Sprintf("csrf_token=%s!abcd; Max-Age=2592000", token),
			},
			"Location":     {"/"},
			"Content-Type": {"text/html; charset=utf-8"},
		}, s.ht.Writer.Header())

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})

	t.Run("req already has session, session do not have enough parts", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234; Max-Age=2592000; HttpOnly; SameSite=Strict")

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})

	t.Run("req already has session, user id is not number", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234a:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		calls := s.repo.FindUserSessionCalls()
		assert.Equal(t, 0, len(calls))

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})

	t.Run("req already has session, user id is not number", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234a:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		calls := s.repo.FindUserSessionCalls()
		assert.Equal(t, 0, len(calls))

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})

	t.Run("req already has session, found user session", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")

		s.stubFindSess(model.NullUserSession{
			Valid: true,
			Data: model.UserSession{
				UserID:    1234,
				SessionID: "some-session-id",
			},
		})

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, true, continuing)
		assert.Equal(t, nil, err)

		assert.Equal(t, http.Header{}, s.ht.Writer.Header())
		assert.Equal(t, http.StatusOK, s.ht.Writer.Code)
	})

	t.Run("req already has session, not found, redirect hx", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewGet("/users")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")
		s.ht.Req.Header.Add("HX-Request", "true")

		s.stubFindSess(model.NullUserSession{})

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		headers := s.ht.Writer.Header()
		assert.Equal(t, 2, len(headers["Set-Cookie"]))
		delete(headers, "Set-Cookie")

		assert.Equal(t, http.Header{
			"Hx-Redirect": {"/"},
		}, headers)
		assert.Equal(t, http.StatusOK, s.ht.Writer.Code)
	})

	t.Run("post, has session, without csrf_token, redirect to home", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewPost("/users", "")
		s.ht.Req.Header.Add("Cookie", "session_id=sess:1234:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict")

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		calls := s.repo.FindUserSessionCalls()
		assert.Equal(t, 0, len(calls))

		headers := s.ht.Writer.Header()
		assert.Equal(t, 2, len(headers["Set-Cookie"]))
		delete(headers, "Set-Cookie")

		assert.Equal(t, http.Header{
			"Location": {"/"},
		}, headers)

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})

	t.Run("post, has session, with csrf_token, success", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewPost("/users", "")
		s.ht.Req.Header.Add(
			"Cookie",
			"session_id=sess:1234:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict",
		)

		token := "c2VzczoxMjM0OnNvbWUtc2Vzc2lvbi1pZCExMjM07H8KcrNzLLI07eG4eyWCAzHv75y8nWyihL0Tij11wZo="
		s.ht.Req.Header.Add(
			"X-Csrf-Token",
			token+"!1234",
		)

		s.stubFindSess(model.NullUserSession{
			Valid: true,
		})

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, true, continuing)
		assert.Equal(t, nil, err)

		calls := s.repo.FindUserSessionCalls()
		assert.Equal(t, 1, len(calls))

		headers := s.ht.Writer.Header()
		assert.Equal(t, http.Header{}, headers)
		assert.Equal(t, http.StatusOK, s.ht.Writer.Code)
	})

	t.Run("post, has session, with csrf_token, mismatch hmac", func(t *testing.T) {
		s := newServiceTest()

		s.ht.NewPost("/users", "")
		s.ht.Req.Header.Add(
			"Cookie",
			"session_id=sess:1234:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict",
		)

		token := "c2VzczoxMjM0OnNvbWUtc2Vzc2lvbi1pZCExMjM07H8KcrNzLLI07eG4eyWCAzHv75y8nWyihL0Tij11wZo="
		s.ht.Req.Header.Add(
			"X-Csrf-Token",
			token+"!12345",
		)

		continuing, err := s.svc.Handle(s.ht.NewContext())
		assert.Equal(t, false, continuing)
		assert.Equal(t, nil, err)

		calls := s.repo.FindUserSessionCalls()
		assert.Equal(t, 0, len(calls))

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})
}

func TestService_Verify_Token(t *testing.T) {
	t.Run("without cookie", func(t *testing.T) {
		s := newServiceTest()
		s.ht.NewPost("/users", "")

		ok, err := s.svc.VerifyCSRFToken(s.ht.NewContext(), "")
		assert.Equal(t, false, ok)
		assert.Equal(t, nil, err)

		headers := s.ht.Writer.Header()
		_, exist := headers["Set-Cookie"]
		assert.Equal(t, true, exist)

		assert.Equal(t, http.StatusOK, s.ht.Writer.Code)
	})

	t.Run("with cookie and correct token", func(t *testing.T) {
		s := newServiceTest()
		s.ht.NewPost("/users", "")
		s.ht.Req.Header.Add("Cookie", "session_id=pre:some-session-id")

		token := "cHJlOnNvbWUtc2Vzc2lvbi1pZCFhYmNk7H8KcrNzLLI07eG4eyWCAzHv75y8nWyihL0Tij11wZo=!abcd"

		ok, err := s.svc.VerifyCSRFToken(s.ht.NewContext(), token)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)

		headers := s.ht.Writer.Header()
		assert.Equal(t, http.Header{}, headers)

		assert.Equal(t, http.StatusOK, s.ht.Writer.Code)
	})

	t.Run("with cookie and invalid token", func(t *testing.T) {
		s := newServiceTest()
		s.ht.NewPost("/users", "")
		s.ht.Req.Header.Add("Cookie", "session_id=pre:some-session-id")

		token := "cHJlOnNvbWUtc2Vzc2lvbi1pZCFhYmNk7H8KcrNzLLI07eG4eyWCAzHv75y8nWyihL0Tij11wZo=!abcde"

		ok, err := s.svc.VerifyCSRFToken(s.ht.NewContext(), token)
		assert.Equal(t, false, ok)
		assert.Equal(t, nil, err)

		headers := s.ht.Writer.Header()
		_, existed := headers["Set-Cookie"]
		assert.Equal(t, true, existed)

		assert.Equal(t, http.StatusTemporaryRedirect, s.ht.Writer.Code)
	})
}

func TestRandService(t *testing.T) {
	r := &randImpl{}
	s, err := r.RandString(32)
	assert.Equal(t, nil, err)
	fmt.Println("LEN:", len(s))
	fmt.Println("VAL:", s)
}
