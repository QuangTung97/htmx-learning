package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"htmx/pkg/route"
)

var counter atomic.Int64

func main() {
	mux := route.NewMux()

	mux.Get("/", func(ctx route.Context) error {
		return ctx.Render("main.html", nil)
	})

	mux.Get("/reload", func(ctx route.Context) error {
		type tmplData struct {
			Count int64
		}
		return ctx.Render("reload.html", tmplData{
			Count: counter.Add(1),
		})
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
		http.StripPrefix(
			"/public/", http.FileServer(http.Dir("./public")),
		),
	)

	fmt.Println("Start HTTP on :8080")
	if err := http.ListenAndServe(":8080", mux.GetMux()); err != nil {
		panic(err)
	}
}
