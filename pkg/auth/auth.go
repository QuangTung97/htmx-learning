package auth

import (
	"fmt"

	"github.com/QuangTung97/svloc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"htmx/config"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
)

type Service interface {
	AuthCodeURL(provider Provider, state string) string
}

var ServiceLoc = svloc.Register[Service](func(unv *svloc.Universe) Service {
	return NewService(
		config.Loc.Get(unv).Auth,
	)
})

type serviceImpl struct {
	googleConf *oauth2.Config
}

func NewService(conf config.Auth) Service {
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

	return &serviceImpl{
		googleConf: googleConf,
	}
}

func (s *serviceImpl) AuthCodeURL(_ Provider, state string) string {
	return s.googleConf.AuthCodeURL(state)
}
