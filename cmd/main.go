package main

import (
	"fmt"
	"net/http"

	"htmx/pkg/route"
)

func main() {
	mux := route.NewMux()

	mux.Get("/", func(ctx route.Context) error {
		return ctx.Render("main.html", nil)
	})

	mux.Get("/reload", func(ctx route.Context) error {
		return ctx.Render("reload.html", nil)
	})

	mux.Route("/users", func(router route.Router) {
		router.Get("/{userId}", func(ctx route.Context) error {
			fmt.Println("USERID:", ctx.GetParam("userId"))
			fmt.Println("Another:", ctx.GetParam("another"))
			return nil
		})
	})

	fmt.Println("Start HTTP on :8080")
	if err := http.ListenAndServe(":8080", mux.GetMux()); err != nil {
		panic(err)
	}
}
