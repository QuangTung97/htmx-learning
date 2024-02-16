package route

import (
	"context"
	"html/template"
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

func (Context) RenderHTML(tmpl views.Template, data any) (template.HTML, error) {
	return views.ExecuteHTML(tmpl, data)
}

func (c Context) Render(template views.Template, data any) error {
	return views.Execute(c.Writer, template, data)
}

const hxRequestHeader = "HX-Request"

func (c Context) HasHxRequestHeader() bool {
	req := c.Req.Header.Get(hxRequestHeader)
	if len(req) == 0 {
		return false
	}
	return true
}

func (c Context) SetHXRequestHeader() {
	c.Req.Header.Set(hxRequestHeader, "true")
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

const hxPushURLHeader = "Hx-Push-Url"

func (c Context) View(template views.Template, data any) error {
	if c.IsHxRequest() {
		if len(c.Writer.Header().Get(hxPushURLHeader)) == 0 {
			redirectURL := util.GetURLPathAndQuery(c.Req.URL)
			c.Writer.Header().Set(hxPushURLHeader, redirectURL)
		}
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
