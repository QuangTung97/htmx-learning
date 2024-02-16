package fragments

import (
	"html/template"

	"htmx/pkg/route"
)

type SampleContentData struct {
	LoggedIn bool
	Table    []string
	Reload   template.HTML
}

func (d SampleContentData) RenderHTML(ctx route.Context) (template.HTML, error) {
	return ctx.RenderHTML(TemplateSampleContent, d)
}
