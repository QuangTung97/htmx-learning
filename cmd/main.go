package main

import (
	"fmt"
	"net/http"

	"htmx/views"
)

func main() {
	http.HandleFunc("/reload", func(writer http.ResponseWriter, request *http.Request) {
		err := views.Execute(writer, "reload.html", nil)
		fmt.Println("ERROR:", err)
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_ = views.Execute(writer, "main.html", nil)
	})

	fmt.Println("Start HTTP on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
