package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"

	"htmx/config/prod"
	"htmx/pkg/auth"
	auth_handlers "htmx/pkg/auth/handlers"
	"htmx/pkg/route"
	"htmx/views"
	"htmx/views/routes"
)

var counter atomic.Int64

func disableCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-store")
		handler.ServeHTTP(writer, request)
	})
}

func main() {
	unv := prod.NewUniverse()

	mux := route.MuxLoc.Get(unv)

	// ==========================
	// Setup Middlewares
	// ==========================
	mux.GetMux().Use(
		auth.InitMiddleware(unv),
	)

	mux.Init()

	mux.Get(routes.Home, func(ctx route.Context) error {
		_, ok := auth.GetUserInfoNull(ctx.Ctx)
		return ctx.View(views.TemplateBody, views.BodyData{
			LoggedIn: ok,
		})
	})

	mux.Post("/reload", func(ctx route.Context) error {
		if _, ok := auth.GetUserInfoNull(ctx.Ctx); !ok {
			return auth.ErrUserNotYetLoggedIn
		}
		type tmplData struct {
			Count int64
		}
		return ctx.Render("reload.html", tmplData{
			Count: counter.Add(1),
		})
	})

	auth_handlers.Register(unv)

	mux.Route("/users", func(router route.Router) {
		router.Get("/{userId}", func(ctx route.Context) error {
			fmt.Println("USERID:", ctx.GetParam("userId"))
			fmt.Println("Another:", ctx.GetParam("another"))
			return errors.New("user error")
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

	fmt.Println("Start HTTP on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux.GetMux()); err != nil {
		panic(err)
	}
}
