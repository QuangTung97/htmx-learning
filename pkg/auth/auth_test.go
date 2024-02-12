package auth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"htmx/model"
	"htmx/pkg/testhelper"
	"htmx/pkg/util"
)

type serviceTest struct {
	ht   *testhelper.HTTPTest
	repo *RepositoryMock
	rand *RandServiceMock
	svc  Service
}

func newServiceTest() *serviceTest {
	s := &serviceTest{
		ht:   testhelper.NewHTTPTest(),
		repo: &RepositoryMock{},
		rand: &RandServiceMock{},
	}

	s.svc = NewService(s.repo, s.rand)
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

		assert.Equal(t, http.Header{
			"Set-Cookie": []string{
				"session_id=pre:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict",
				"csrf_token=cHJlOnNvbWUtc2Vzc2lvbi1pZCFhYmNkthNnmggU2ex3L5XXeMNfxf8Wl8STcVZTxscSFEKSxa0=!abcd; Max-Age=2592000; SameSite=Strict",
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

		assert.Equal(t, http.Header{
			"Set-Cookie": {
				"session_id=pre:some-session-id; Max-Age=2592000; HttpOnly; SameSite=Strict",
				"csrf_token=cHJlOnNvbWUtc2Vzc2lvbi1pZCFhYmNkthNnmggU2ex3L5XXeMNfxf8Wl8STcVZTxscSFEKSxa0=!abcd; Max-Age=2592000; SameSite=Strict",
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
}

func TestRandService(t *testing.T) {
	r := &randImpl{}
	s, err := r.RandString(32)
	assert.Equal(t, nil, err)
	fmt.Println(len(s))
	fmt.Println(s)
}
