package handlers

import (
	"net/url"

	"github.com/QuangTung97/svloc"

	"htmx/pkg/auth"
	"htmx/pkg/route"
	"htmx/pkg/util"
	"htmx/views/routes"
)

const oauthState = "some-state"

func getCurrentURLPath(ctx route.Context) string {
	currentURL := ctx.Req.Header.Get("HX-Current-URL")
	if len(currentURL) == 0 {
		return ""
	}
	u, err := url.Parse(currentURL)
	if err != nil {
		return ""
	}
	return util.GetURLPathAndQuery(u)
}

func Register(unv *svloc.Universe) {
	mux := route.MuxLoc.Get(unv)

	mux.Get("/login", func(ctx route.Context) error {
		newURL := ctx.Req.URL
		newURL.RawQuery += "backUrl=" + url.QueryEscape(getCurrentURLPath(ctx))
		ctx.Req.URL = newURL

		return ctx.View("auth/google-login.html", nil)
	})

	authSvc := auth.OAuthServiceLoc.Get(unv)

	mux.Post(routes.OAuthGoogleLogin, func(ctx route.Context) error {
		redirectURL := authSvc.AuthCodeURL(auth.ProviderGoogle, oauthState)
		ctx.Redirect(redirectURL)
		return nil
	})

	loginSvc := auth.LoginServiceLoc.Get(unv)
	mux.Get(routes.AuthCallback, loginSvc.HandleCallback)
}
