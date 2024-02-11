package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/QuangTung97/svloc"

	"htmx/pkg/auth"
	"htmx/pkg/route"
	"htmx/views/routes"
)

const oauthState = "some-state"

func Register(unv *svloc.Universe, mux *route.Mux) {
	mux.Get("/login", func(ctx route.Context) error {
		return ctx.View("auth/google-login.html", nil)
	})

	authSvc := auth.ServiceLoc.Get(unv)

	mux.Post(routes.OAuthGoogleLogin, func(ctx route.Context) error {
		redirectURL := authSvc.AuthCodeURL(auth.ProviderGoogle, oauthState)
		ctx.HXRedirect(redirectURL)
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
