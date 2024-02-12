package auth

import (
	"net/http"
	"strings"

	"htmx/pkg/route"
)

func Middleware(s Service) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.URL.Path, "/public") {
				handler.ServeHTTP(writer, request)
				return
			}

			continuing, err := s.Handle(route.NewContext(writer, request))
			if err != nil {
				route.ResponseError(writer, err)
				return
			}

			if !continuing {
				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}
