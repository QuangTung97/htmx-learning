package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuangTung97/svloc"

	"htmx/config"
	"htmx/model"
	"htmx/pkg/route"
	"htmx/views/routes"
)

type Service interface {
	Handle(ctx route.Context) (continuing bool, err error)
}

type serviceImpl struct {
	repo Repository
	rand RandService

	isProd    bool
	secretKey string
}

func NewService(
	conf config.Auth,
	repo Repository,
	randSvc RandService,
) Service {
	if len(conf.CSRFHMACSecret) == 0 {
		panic("Missing csrf_hmac_secret")
	}
	return &serviceImpl{
		repo: repo,
		rand: randSvc,

		secretKey: conf.CSRFHMACSecret,
	}
}

var ServiceLoc = svloc.Register[Service](func(unv *svloc.Universe) Service {
	return NewService(
		config.Loc.Get(unv).Auth,
		RepoLoc.Get(unv),
		RandServiceLoc.Get(unv),
	)
})

const sessionIDCookie = "session_id"
const csrfTokenCookie = "csrf_token"
const sessionByteSize = 32

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

func (s *serviceImpl) setPreSession(ctx route.Context) error {
	sessID, err := s.rand.RandString(sessionByteSize)
	if err != nil {
		return err
	}
	sessID = preLoginSessionPrefix + ":" + sessID
	fmt.Println("GENERATED Session ID")

	const maxAge = 30 * 3600 * 24 // 30 days

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  sessionIDCookie,
		Value: sessID,

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

		MaxAge: maxAge,
		Secure: s.isProd,
	})

	return nil
}

func (s *serviceImpl) redirectToHome(ctx route.Context) (bool, error) {
	if ctx.IsHxRequest() {
		ctx.HXRedirect(routes.Home)
	} else {
		http.Redirect(ctx.Writer, ctx.Req, routes.Home, http.StatusTemporaryRedirect)
	}
	return false, s.setPreSession(ctx)
}

func (s *serviceImpl) checkCSRFToken(ctx route.Context, sessCookie *http.Cookie) (bool, error) {
	sessionID := sessCookie.Value

	method := ctx.Req.Method
	if method == http.MethodGet {
		return true, nil
	}

	token := ctx.Req.Header.Get("X-Csrf-Token")

	parts := strings.Split(token, "!")
	if len(parts) != 2 {
		log.Println("[WARN] Not Found CSRF Token", ctx.Req.URL.String(), token)
		return s.redirectToHome(ctx)
	}

	compareVal := s.generateHMACSig(sessionID, parts[1])
	if compareVal != parts[0] {
		log.Println("[WARN] Mismatch CSRF Token", ctx.Req.URL.String(), token)
		return s.redirectToHome(ctx)
	}

	return true, nil
}

func (s *serviceImpl) Handle(ctx route.Context) (bool, error) {
	sessCookie, err := ctx.Req.Cookie(sessionIDCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return true, s.setPreSession(ctx)
		}
		return false, err
	}

	parts := strings.Split(sessCookie.Value, ":")
	if parts[0] == preLoginSessionPrefix {
		// TODO Check parts >= 2
		return s.checkCSRFToken(ctx, sessCookie)
	}

	if len(parts) != 3 {
		return s.redirectToHome(ctx)
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return s.redirectToHome(ctx)
	}

	continuing, _ := s.checkCSRFToken(ctx, sessCookie)
	if !continuing {
		return false, nil
	}

	userSess, err := s.repo.FindUserSession(ctx.Ctx, model.UserID(userID), model.SessionID(parts[2]))
	if err != nil {
		return false, err
	}
	if !userSess.Valid {
		return s.redirectToHome(ctx)
	}

	return true, nil
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
