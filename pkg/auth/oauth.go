package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/QuangTung97/svloc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"htmx/config"
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
	HandleCallback(ctx route.Context) error
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
	oauth OAuthService
	svc   Service
}

func NewLoginService(
	oauthSvc OAuthService, svc Service,
) LoginService {
	return &loginServiceImpl{
		oauth: oauthSvc,
	}
}

const oauthState = "some-state"

func (s *loginServiceImpl) HandleCallback(ctx route.Context) error {
	state := ctx.Req.URL.Query().Get("state")
	// TODO Check State

	code := ctx.Req.URL.Query().Get("code")

	if state != oauthState {
		return errors.New("invalid oauth state")
	}

	token, err := s.oauth.Exchange(ctx.Ctx, ProviderGoogle, code)
	if err != nil {
		return err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer func() { _ = response.Body.Close() }()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed reading response body: %s", err.Error())
	}

	fmt.Println("CONTENTS:", string(contents))
	http.Redirect(ctx.Writer, ctx.Req, "/", http.StatusTemporaryRedirect)

	return nil
}
