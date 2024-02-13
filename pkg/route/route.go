package route

import (
	"net/http"

	"github.com/QuangTung97/svloc"
	"github.com/go-chi/chi/v5"

	"htmx/views/routes"
)

type Router struct {
	router    chi.Router
	errorView ErrorView
}

type Mux struct {
	Router
	mux *chi.Mux
}

var MuxLoc = svloc.Register[*Mux](func(unv *svloc.Universe) *Mux {
	return NewMux(ErrorViewLoc.Get(unv))
})

type Handler func(ctx Context) error

func NewMux(errorView ErrorView) *Mux {
	r := chi.NewRouter()

	return &Mux{
		Router: Router{
			router:    r,
			errorView: errorView,
		},
		mux: r,
	}
}

func (m *Mux) Init() {
	m.Router.router.Get(routes.Error, func(writer http.ResponseWriter, request *http.Request) {
		m.errorView.Render(NewContext(writer, request))
	})
}

func (m *Mux) GetMux() *chi.Mux {
	return m.mux
}

func (r Router) Get(pattern string, handler Handler) {
	r.router.Get(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := NewContext(writer, request)
		err := handler(ctx)
		r.errorView.Redirect(ctx, err)
	})
}

func (r Router) Post(pattern string, handler Handler) {
	r.router.Post(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := NewContext(writer, request)

		err := handler(ctx)
		r.errorView.Redirect(ctx, err)
	})
}

// Route ...
func (r Router) Route(pattern string, fn func(router Router)) {
	r.router.Route(pattern, func(innerRouter chi.Router) {
		fn(Router{router: innerRouter, errorView: r.errorView})
	})
}
