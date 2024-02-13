package route

import (
	"encoding/json"
	"net/http"

	"github.com/QuangTung97/svloc"
	"github.com/go-chi/chi/v5"

	"htmx/views/routes"
)

// Router ...
type Router struct {
	router chi.Router
}

// Mux ...
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

	r.Get(routes.Error, func(writer http.ResponseWriter, request *http.Request) {
		errorView.Render(NewContext(writer, request))
	})

	return &Mux{
		Router: Router{
			router: r,
		},
		mux: r,
	}
}

func (r Router) Get(pattern string, handler Handler) {
	r.router.Get(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Context{
			Ctx:    request.Context(),
			Req:    request,
			Writer: writer,
		}

		err := handler(ctx)
		ResponseError(writer, err)
	})
}

func (r Router) Post(pattern string, handler Handler) {
	r.router.Post(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Context{
			Ctx:    request.Context(),
			Req:    request,
			Writer: writer,
		}

		err := handler(ctx)
		ResponseError(writer, err)
	})
}

// Route ...
func (r Router) Route(pattern string, fn func(router Router)) {
	r.router.Route(pattern, func(r chi.Router) {
		fn(Router{router: r})
	})
}

// GetMux ...
func (m *Mux) GetMux() *chi.Mux {
	return m.mux
}

func ResponseError(writer http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	type errorResponse struct {
		Message string `json:"message"`
	}
	_ = json.NewEncoder(writer).Encode(errorResponse{
		Message: err.Error(),
	})
}
