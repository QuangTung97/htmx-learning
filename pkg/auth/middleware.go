package auth

import (
	"net/http"
	"strings"

	"github.com/QuangTung97/svloc"

	"htmx/pkg/route"
)

func Middleware(s Service, errorView route.ErrorView) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.URL.Path, "/public/") {
				handler.ServeHTTP(writer, request)
				return
			}

			ctx := route.NewContext(writer, request)
			continuing, err := s.Handle(&ctx)
			if err != nil {
				errorView.Redirect(ctx, err)
				return
			}

			if !continuing {
				return
			}

			handler.ServeHTTP(writer, ctx.Req)
		})
	}
}

func InitMiddleware(unv *svloc.Universe) func(handler http.Handler) http.Handler {
	return Middleware(
		ServiceLoc.Get(unv),
		route.ErrorViewLoc.Get(unv),
	)
}
