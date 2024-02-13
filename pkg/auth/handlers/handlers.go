package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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

	mux.Get(routes.AuthCallback, func(ctx route.Context) error {
		state := ctx.Req.URL.Query().Get("state")
		// TODO Check State

		code := ctx.Req.URL.Query().Get("code")

		if state != oauthState {
			return errors.New("some error")
		}

		token, err := authSvc.Exchange(ctx.Ctx, auth.ProviderGoogle, code)
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
	})
}
