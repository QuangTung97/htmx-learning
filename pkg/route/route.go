package route

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"htmx/views"
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

// Context ...
type Context struct {
	Ctx    context.Context
	Req    *http.Request
	Writer http.ResponseWriter
}

// Render ...
func (c Context) Render(template string, data any) error {
	return views.Execute(c.Writer, template, data)
}

// View ...
func (c Context) View(template string, data any) error {
	body, err := views.ExecuteHTML(template, data)
	if err != nil {
		return err
	}
	return views.View(c.Writer, body)
}

func (c Context) GetParam(key string) string {
	return chi.URLParam(c.Req, key)
}

type Handler func(ctx Context) error

func NewMux() *Mux {
	r := chi.NewRouter()
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
		if err != nil {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			type errorResponse struct {
				Message string `json:"message"`
			}
			_ = json.NewEncoder(writer).Encode(errorResponse{
				Message: err.Error(),
			})
		}
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
		if err != nil {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			type errorResponse struct {
				Message string `json:"message"`
			}
			_ = json.NewEncoder(writer).Encode(errorResponse{
				Message: err.Error(),
			})
		}
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
