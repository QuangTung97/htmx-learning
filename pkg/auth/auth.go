package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuangTung97/svloc"

	"htmx/config"
	"htmx/model"
	"htmx/pkg/dbtx"
	"htmx/pkg/route"
	"htmx/views/routes"
)

type Service interface {
	Handle(ctx *route.Context) (continuing bool, err error)

	VerifyCSRFToken(ctx route.Context, token string) bool
	SetSession(ctx route.Context, sess model.UserSession) error
	SetPreLoginSession(ctx route.Context) error
}

type serviceImpl struct {
	provider dbtx.Provider
	repo     Repository
	rand     RandService

	errorView route.ErrorView

	isProd    bool
	secretKey string

	logger *slog.Logger
}

func NewService(
	conf config.Auth,
	provider dbtx.Provider,
	repo Repository,
	randSvc RandService,
	errorView route.ErrorView,
) Service {
	if len(conf.CSRFHMACSecret) == 0 {
		panic("Missing csrf_hmac_secret")
	}
	return &serviceImpl{
		provider: provider,
		repo:     repo,
		rand:     randSvc,

		errorView: errorView,

		secretKey: conf.CSRFHMACSecret,
	}
}

var ServiceLoc = svloc.Register[Service](func(unv *svloc.Universe) Service {
	return NewService(
		config.Loc.Get(unv).Auth,
		dbtx.ProviderLoc.Get(unv),
		RepoLoc.Get(unv),
		RandServiceLoc.Get(unv),
		route.ErrorViewLoc.Get(unv),
	)
})

const sessionIDCookie = "session_id"
const csrfTokenCookie = "csrf_token"
const SessionByteSize = 32

func (s *serviceImpl) newHMAC() hash.Hash {
	return hmac.New(sha256.New, []byte(s.secretKey))
}

func (s *serviceImpl) generateHMACSig(sessionID string, randomVal string) string {
	msg := sessionID + "!" + randomVal
	hmacVal := s.newHMAC().Sum([]byte(msg))
	return base64.StdEncoding.EncodeToString(hmacVal)
}

func (s *serviceImpl) computeCSRFToken(sessionID string) (string, error) {
	randomVal, err := s.rand.RandString(16)
	if err != nil {
		return "", err
	}

	val := s.generateHMACSig(sessionID, randomVal)
	csrfToken := val + "!" + randomVal
	return csrfToken, nil
}

const preLoginSessionPrefix = "pre"
const sessionPrefix = "sess"

func (s *serviceImpl) setPreSession(ctx route.Context) error {
	sessID, err := s.rand.RandString(SessionByteSize)
	if err != nil {
		return err
	}
	sessID = preLoginSessionPrefix + ":" + sessID
	fmt.Println("Generated Pre-Login Session ID")

	return s.setSessionCookie(ctx, sessID)
}

func (s *serviceImpl) setSessionCookie(ctx route.Context, sessID string) error {
	const maxAge = 30 * 3600 * 24 // 30 days

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  sessionIDCookie,
		Value: sessID,

		Path:     cookiePath,
		MaxAge:   maxAge,
		Secure:   s.isProd,
		HttpOnly: true,
	})

	token, err := s.computeCSRFToken(sessID)
	if err != nil {
		return err
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  csrfTokenCookie,
		Value: token,

		Path:   cookiePath,
		MaxAge: maxAge,
		Secure: s.isProd,
	})

	return nil
}

func (s *serviceImpl) redirectToHome(ctx route.Context) (bool, error) {
	ctx.Redirect(routes.Home)
	return false, s.setPreSession(ctx)
}

func (s *serviceImpl) redirectToError(ctx route.Context, err error) bool {
	_ = s.setPreSession(ctx)
	s.errorView.Redirect(ctx, err)
	return false
}

func (s *serviceImpl) checkCSRFToken(ctx route.Context, sessCookie *http.Cookie) bool {
	method := ctx.Req.Method
	if method == http.MethodGet {
		return true
	}

	sessionID := sessCookie.Value
	token := ctx.Req.Header.Get("X-Csrf-Token")

	return s.verifyTokenWithSession(ctx, sessionID, token)
}

func (s *serviceImpl) verifyTokenWithSession(ctx route.Context, sessionID string, token string) bool {
	parts := strings.Split(token, "!")
	if len(parts) != 2 {
		return s.redirectToError(ctx, errors.New("invalid csrf token"))
	}

	compareVal := s.generateHMACSig(sessionID, parts[1])
	if compareVal != parts[0] {
		return s.redirectToError(ctx, errors.New("mismatch csrf token"))
	}

	return true
}

func (s *serviceImpl) handleError(ctx route.Context, err error) (bool, error) {
	s.redirectToError(ctx, err)
	return false, nil
}

const cookiePath = "/"

func getCookie(ctx route.Context, name string) (*http.Cookie, error) {
	return ctx.Req.Cookie(name)
}

func (s *serviceImpl) Handle(ctxPtr *route.Context) (bool, error) {
	ctx := *ctxPtr

	sessCookie, err := getCookie(ctx, sessionIDCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return true, s.setPreSession(ctx)
		}
		return false, err
	}

	parts := strings.Split(sessCookie.Value, ":")
	if parts[0] == preLoginSessionPrefix {
		return s.checkCSRFToken(ctx, sessCookie), nil
	}

	if len(parts) != 3 {
		return s.handleError(ctx, errors.New("invalid session id"))
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return s.handleError(ctx, errors.New("invalid session id"))
	}

	continuing := s.checkCSRFToken(ctx, sessCookie)
	if !continuing {
		return false, nil
	}

	readCtx := s.provider.Readonly(ctx.Ctx)
	userSess, err := s.repo.FindUserSession(readCtx, model.UserID(userID), model.SessionID(parts[2]))
	if err != nil {
		return false, err
	}
	if !userSess.Valid {
		return s.redirectToHome(ctx)
	}

	newCtx := SetUserInfo(ctx.Ctx, UserInfo{
		User: model.User{
			ID: userSess.Data.UserID,
		},
		Session: userSess.Data,
	})
	ctxPtr.Ctx = newCtx
	ctxPtr.Req = ctx.Req.WithContext(newCtx)

	return true, nil
}

func (s *serviceImpl) VerifyCSRFToken(ctx route.Context, token string) (ok bool) {
	sessCookie, err := getCookie(ctx, sessionIDCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return s.redirectToError(ctx, errors.New("no session id when verify csrf token"))
		}
		return s.redirectToError(ctx, err)
	}
	return s.verifyTokenWithSession(ctx, sessCookie.Value, token)
}

func (s *serviceImpl) SetSession(ctx route.Context, sess model.UserSession) error {
	sessID := fmt.Sprintf("%s:%d:%s", sessionPrefix, sess.UserID, sess.SessionID)
	return s.setSessionCookie(ctx, sessID)
}

func (s *serviceImpl) SetPreLoginSession(ctx route.Context) error {
	return s.setPreSession(ctx)
}

type RandService interface {
	RandString(size int) (string, error)
}

var RandServiceLoc = svloc.Register[RandService](func(unv *svloc.Universe) RandService {
	return &randImpl{}
})

type randImpl struct {
}

func (r *randImpl) RandString(size int) (string, error) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
