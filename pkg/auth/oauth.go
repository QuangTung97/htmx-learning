package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/QuangTung97/svloc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"htmx/config"
	"htmx/model"
	"htmx/pkg/dbtx"
	"htmx/pkg/route"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
)

type OAuthService interface {
	AuthCodeURL(provider Provider, state string) string
	Exchange(ctx context.Context, provider Provider, code string) (*oauth2.Token, error)
}

type LoginService interface {
	BuildOAuthState(ctx route.Context) (string, error)
	HandleCallback(ctx route.Context) error
	HandleLogOut(ctx route.Context) error
}

var OAuthServiceLoc = svloc.Register[OAuthService](func(unv *svloc.Universe) OAuthService {
	return NewOAuthService(
		config.Loc.Get(unv).Auth,
	)
})

type oauthServiceImpl struct {
	googleConf *oauth2.Config
}

func NewOAuthService(conf config.Auth) OAuthService {
	scopes := []string{
		"openid",
		"https://www.googleapis.com/auth/userinfo.email",
	}

	fmt.Println("GOOGLE CLIENT ID:", conf.GoogleClientID)

	googleConf := &oauth2.Config{
		ClientID:     conf.GoogleClientID,
		ClientSecret: conf.GoogleClientSecret,
		Scopes:       scopes,
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint:     google.Endpoint,
	}

	return &oauthServiceImpl{
		googleConf: googleConf,
	}
}

func (s *oauthServiceImpl) AuthCodeURL(_ Provider, state string) string {
	return s.googleConf.AuthCodeURL(state)
}

func (s *oauthServiceImpl) Exchange(ctx context.Context, _ Provider, code string) (*oauth2.Token, error) {
	return s.googleConf.Exchange(ctx, code)
}

// ===========================================
// Login Service
// ===========================================

type loginServiceImpl struct {
	provider dbtx.Provider
	oauth    OAuthService
	svc      Service
	rand     RandService
	repo     Repository
}

func NewLoginService(
	provider dbtx.Provider,
	oauthSvc OAuthService, svc Service,
	randSvc RandService,
	repo Repository,
) LoginService {
	return &loginServiceImpl{
		provider: provider,
		oauth:    oauthSvc,
		svc:      svc,
		rand:     randSvc,
		repo:     repo,
	}
}

var LoginServiceLoc = svloc.Register[LoginService](func(unv *svloc.Universe) LoginService {
	return NewLoginService(
		dbtx.ProviderLoc.Get(unv),
		OAuthServiceLoc.Get(unv),
		ServiceLoc.Get(unv),
		RandServiceLoc.Get(unv),
		RepoLoc.Get(unv),
	)
})

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func (s *loginServiceImpl) BuildOAuthState(ctx route.Context) (string, error) {
	cookie, err := getCookie(ctx, csrfTokenCookie)
	if err != nil {
		return "", errors.New("invalid csrf token")
	}
	return buildOAuthState(cookie.Value), nil
}

var ErrUserAlreadyLoggedIn = errors.New("user already logged in")

var ErrUserNotYetLoggedIn = errors.New("user not yet logged in")

func (s *loginServiceImpl) HandleCallback(ctx route.Context) error {
	if _, ok := GetUserInfoNull(ctx.Ctx); ok {
		return ErrUserAlreadyLoggedIn
	}

	state := ctx.Req.URL.Query().Get("state")
	csrfToken := getCSRFTokenFromState(state)
	if !s.svc.VerifyCSRFToken(ctx, csrfToken) {
		return nil
	}

	code := ctx.Req.URL.Query().Get("code")

	token, err := s.oauth.Exchange(ctx.Ctx, ProviderGoogle, code)
	if err != nil {
		return err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer func() { _ = response.Body.Close() }()

	var user googleUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return fmt.Errorf("failed reading response body: %s", err.Error())
	}

	if err := s.handleGoogleUser(ctx, user); err != nil {
		return err
	}

	http.Redirect(ctx.Writer, ctx.Req, "/", http.StatusTemporaryRedirect)
	return nil
}

func (s *loginServiceImpl) handleGoogleUser(ctx route.Context, user googleUser) error {
	autoCtx := s.provider.Autocommit(ctx.Ctx)

	nullUser, err := s.repo.FindUser(autoCtx, string(ProviderGoogle), user.ID)
	if err != nil {
		return err
	}

	loginUser := nullUser.Data
	if !nullUser.Valid {
		loginUser = model.User{
			Provider:    string(ProviderGoogle),
			OAuthUserID: user.ID,
			Email:       user.Email,
			ImageURL:    user.Picture,

			Status: model.UserStatusActive,
		}
		userID, err := s.repo.InsertUser(autoCtx, loginUser)
		if err != nil {
			return err
		}
		loginUser.ID = userID
	}

	sessionID, err := s.rand.RandString(SessionByteSize)
	if err != nil {
		return err
	}

	userSess := model.UserSession{
		UserID:    loginUser.ID,
		SessionID: model.SessionID(sessionID),
		Status:    model.UserSessionStatusActive,
	}
	if err := s.repo.InsertUserSession(autoCtx, userSess); err != nil {
		return err
	}

	return s.svc.SetSession(ctx, userSess)
}

func (s *loginServiceImpl) HandleLogOut(ctx route.Context) error {
	return s.svc.SetPreLoginSession(ctx)
}

const stateCSRFPrefix = "csrf="

func buildOAuthState(csrfToken string) string {
	return stateCSRFPrefix + csrfToken
}

func getCSRFTokenFromState(state string) string {
	return strings.TrimPrefix(state, stateCSRFPrefix)
}
