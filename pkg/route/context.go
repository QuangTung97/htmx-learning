package route

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"htmx/pkg/util"
	"htmx/views"
)

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

func (c Context) isHxRequest() bool {
	req := c.Req.Header.Get("HX-Request")
	return len(req) > 0
}

// View ...
func (c Context) View(template string, data any) error {
	if c.isHxRequest() {
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

func (c Context) HXRedirect(redirectURL string) {
	c.Writer.Header().Set("HX-Redirect", redirectURL)
}
