package main

import (
	"fmt"
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
		ctx.Writer.Header().Set("HX-Redirect", redirectURL)
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
