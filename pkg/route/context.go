package route

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"htmx/pkg/util"
	"htmx/views"
)

type Context struct {
	Ctx    context.Context
	Req    *http.Request
	Writer http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return Context{
		Ctx:    r.Context(),
		Req:    r,
		Writer: w,
	}
}

func (c Context) Render(template string, data any) error {
	return views.Execute(c.Writer, template, data)
}

func (c Context) HasHxRequestHeader() bool {
	req := c.Req.Header.Get("HX-Request")
	if len(req) == 0 {
		return false
	}
	return true
}

func (c Context) IsHxRequest() bool {
	if !c.HasHxRequestHeader() {
		return false
	}

	boosted := c.Req.Header.Get("HX-Boosted")
	if len(boosted) > 0 {
		return false
	}

	return true
}

// View ...
func (c Context) View(template string, data any) error {
	if c.IsHxRequest() {
		redirectURL := util.GetURLPathAndQuery(c.Req.URL)
		c.Writer.Header().Set("HX-Push-Url", redirectURL)
		return views.Execute(c.Writer, template, data)
	}

	body, err := views.ExecuteHTML(template, data)
	if err != nil {
		return err
	}
	return views.View(c.Writer, body)
}

func (c Context) GetParam(key string) string {
	return chi.URLParam(c.Req, key)
}

func (c Context) Redirect(toURL string) {
	if c.HasHxRequestHeader() {
		c.Writer.Header().Set("HX-Redirect", toURL)
		return
	}
	http.Redirect(c.Writer, c.Req, toURL, http.StatusTemporaryRedirect)
}
