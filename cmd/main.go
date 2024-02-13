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

	mux.Get("/", func(ctx route.Context) error {
		return ctx.View("body.html", nil)
	})

	mux.Post("/reload", func(ctx route.Context) error {
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
