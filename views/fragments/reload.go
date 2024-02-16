package fragments

import (
	"html/template"

	"htmx/pkg/route"
)

type ReloadData struct {
	Count int64
}

func (d ReloadData) RenderHTML(ctx route.Context) (template.HTML, error) {
	return ctx.RenderHTML(TemplateReload, d)
}

func (d ReloadData) Render(ctx route.Context) error {
	return ctx.Render(TemplateReload, d)
}
