package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/QuangTung97/svloc"

	"htmx/config"
	"htmx/pkg/auth"
	"htmx/pkg/route"
)

var counter atomic.Int64

func disableCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-store")
		handler.ServeHTTP(writer, request)
	})
}

const oauthState = "some-state"

func main() {
	unv := svloc.NewUniverse()
	config.Loc.MustOverrideFunc(unv, func(unv *svloc.Universe) config.Config {
		return config.Load()
	})

	mux := route.NewMux()

	mux.Get("/", func(ctx route.Context) error {
		return ctx.View("body.html", nil)
	})

	mux.Get("/reload", func(ctx route.Context) error {
		type tmplData struct {
			Count int64
		}
		return ctx.Render("reload.html", tmplData{
			Count: counter.Add(1),
		})
	})

	mux.Get("/login", func(ctx route.Context) error {
		return ctx.View("auth/login.html", nil)
	})

	authSvc := auth.ServiceLoc.Get(unv)
	mux.Post("/oauth/login/google", func(ctx route.Context) error {
		redirectURL := authSvc.AuthCodeURL(auth.ProviderGoogle, oauthState)
		ctx.HXRedirect(redirectURL)
		return nil
	})

	mux.Get("/callback", func(ctx route.Context) error {
		state := ctx.Req.URL.Query().Get("state")
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

	mux.Route("/users", func(router route.Router) {
		router.Get("/{userId}", func(ctx route.Context) error {
			fmt.Println("USERID:", ctx.GetParam("userId"))
			fmt.Println("Another:", ctx.GetParam("another"))
			return nil
		})
	})

	mux.GetMux().Handle(
		"/public/*",
		disableCache(
			http.StripPrefix(
				"/public/", http.FileServer(http.Dir("./public")),
			),
		),
	)

	fmt.Println("Start HTTP on :8080")
	if err := http.ListenAndServe(":8080", mux.GetMux()); err != nil {
		panic(err)
	}
}
