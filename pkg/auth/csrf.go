package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"

	"htmx/model"
	"htmx/pkg/route"
)

type CSRFService interface {
	Handle(ctx route.Context) (continuing bool, err error)
}

type csrfServiceImpl struct {
	repo Repository
	rand RandService

	isProd    bool
	secretKey string
}

func NewCSRFService(
	repo Repository,
) CSRFService {
	return &csrfServiceImpl{
		repo: repo,
	}
}

const sessionIDCookie = "session_id"
const csrfTokenCookie = "csrf_token"
const sessionByteSize = 32

func (r *csrfServiceImpl) computeCSRFToken(sessionID string) (string, error) {
	randomVal, err := r.rand.RandString(16)
	if err != nil {
		return "", err
	}

	msg := sessionID + "!" + randomVal
	hmacVal := hmac.New(sha256.New, []byte(r.secretKey)).Sum([]byte(msg))
	csrfToken := base64.StdEncoding.EncodeToString(hmacVal) + "!" + randomVal
	return csrfToken, nil
}

func (r *csrfServiceImpl) setPreSession(ctx route.Context) error {
	sessID, err := r.rand.RandString(sessionByteSize)
	if err != nil {
		return err
	}

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

	// TODO
	return nil
}

func (r *csrfServiceImpl) Handle(ctx route.Context) (bool, error) {
	sessCookie, err := ctx.Req.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return false, r.setPreSession(ctx)
		}
		return false, err
	}

	userSess, err := r.repo.FindUserSession(ctx.Ctx, model.SessionID(sessCookie.Value))
	if err != nil {
		return false, err
	}
	if !userSess.Valid {
		return false, r.setPreSession(ctx)
	}

	return true, nil
}

type RandService interface {
	RandString(size int) (string, error)
}

func NewRandService() RandService {
	return &randImpl{}
}

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
