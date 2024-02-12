package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"
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

func (r *serviceImpl) newHMAC() hash.Hash {
	return hmac.New(sha256.New, []byte(r.secretKey))
}

func (r *serviceImpl) generateHMACSig(sessionID string, randomVal string) string {
	msg := sessionID + "!" + randomVal
	hmacVal := r.newHMAC().Sum([]byte(msg))
	return base64.StdEncoding.EncodeToString(hmacVal)
}

func (r *serviceImpl) computeCSRFToken(sessionID string) (string, error) {
	randomVal, err := r.rand.RandString(16)
	if err != nil {
		return "", err
	}

	val := r.generateHMACSig(sessionID, randomVal)
	csrfToken := val + "!" + randomVal
	return csrfToken, nil
}

const preLoginSessionPrefix = "pre"

func (r *serviceImpl) setPreSession(ctx route.Context) error {
	sessID, err := r.rand.RandString(sessionByteSize)
	if err != nil {
		return err
	}
	sessID = preLoginSessionPrefix + ":" + sessID

	const maxAge = 30 * 3600 * 24 // 30 days

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  sessionIDCookie,
		Value: sessID,

		MaxAge:   maxAge,
		Secure:   r.isProd,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	token, err := r.computeCSRFToken(sessID)
	if err != nil {
		return err
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  csrfTokenCookie,
		Value: token,

		MaxAge:   maxAge,
		Secure:   r.isProd,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func (r *serviceImpl) redirectToHome(ctx route.Context) (bool, error) {
	if ctx.IsHxRequest() {
		ctx.HXRedirect(routes.Home)
	} else {
		http.Redirect(ctx.Writer, ctx.Req, routes.Home, http.StatusTemporaryRedirect)
	}
	return false, r.setPreSession(ctx)
}

func (r *serviceImpl) checkCSRFToken(ctx route.Context, sessionID string) (bool, error) {
	method := ctx.Req.Method
	if method == http.MethodGet {
		return true, nil
	}

	token := ctx.Req.Header.Get("X-Csrf-Token")

	parts := strings.Split(token, "!")
	if len(parts) != 2 {
		return r.redirectToHome(ctx)
	}

	compareVal := r.generateHMACSig(sessionID, parts[1])
	if compareVal != parts[0] {
		return r.redirectToHome(ctx)
	}

	return true, nil
}

func (r *serviceImpl) Handle(ctx route.Context) (bool, error) {
	sessCookie, err := ctx.Req.Cookie(sessionIDCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return true, r.setPreSession(ctx)
		}
		return false, err
	}

	parts := strings.Split(sessCookie.Value, ":")
	if parts[0] == preLoginSessionPrefix {
		// TODO Check parts >= 2
		return r.checkCSRFToken(ctx, parts[1])
	}

	if len(parts) != 3 {
		return r.redirectToHome(ctx)
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return r.redirectToHome(ctx)
	}

	continuing, _ := r.checkCSRFToken(ctx, parts[2])
	if !continuing {
		return false, nil
	}

	userSess, err := r.repo.FindUserSession(ctx.Ctx, model.UserID(userID), model.SessionID(parts[2]))
	if err != nil {
		return false, err
	}
	if !userSess.Valid {
		return r.redirectToHome(ctx)
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
