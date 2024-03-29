package handlers

import (
	"net/url"

	"github.com/QuangTung97/svloc"

	"htmx/pkg/auth"
	"htmx/pkg/route"
	"htmx/pkg/util"
	"htmx/views"
	"htmx/views/routes"
)

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

func loginHandler(ctx route.Context) error {
	newURL := ctx.Req.URL
	newURL.RawQuery += "backUrl=" + url.QueryEscape(getCurrentURLPath(ctx))
	ctx.Req.URL = newURL

	return ctx.View(views.TemplateLogin, nil)
}

func mustNotLoggedIn(handler func(ctx route.Context) error) func(ctx route.Context) error {
	return func(ctx route.Context) error {
		if _, ok := auth.GetUserInfoNull(ctx.Ctx); ok {
			return auth.ErrUserAlreadyLoggedIn
		}
		return handler(ctx)
	}
}

func Register(unv *svloc.Universe) {
	mux := route.MuxLoc.Get(unv)

	mux.Get(routes.Login, mustNotLoggedIn(loginHandler))

	authSvc := auth.OAuthServiceLoc.Get(unv)
	loginSvc := auth.LoginServiceLoc.Get(unv)

	loginSubmitHandler := func(ctx route.Context) error {
		if _, ok := auth.GetUserInfoNull(ctx.Ctx); ok {
			return auth.ErrUserAlreadyLoggedIn
		}

		state, err := loginSvc.BuildOAuthState(ctx)
		if err != nil {
			return err
		}
		redirectURL := authSvc.AuthCodeURL(auth.ProviderGoogle, state)
		ctx.Redirect(redirectURL)
		return nil
	}
	mux.Post(routes.OAuthGoogleLogin, mustNotLoggedIn(loginSubmitHandler))

	mux.Get(routes.AuthCallback, mustNotLoggedIn(loginSvc.HandleCallback))

	mux.Post(routes.Logout, func(ctx route.Context) error {
		if _, ok := auth.GetUserInfoNull(ctx.Ctx); !ok {
			return auth.ErrUserNotYetLoggedIn
		}

		if err := loginSvc.HandleLogOut(ctx); err != nil {
			return err
		}
		ctx.Redirect(routes.Home)
		return nil
	})
}
