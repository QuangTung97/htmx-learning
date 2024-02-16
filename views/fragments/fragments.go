package fragments

import (
	"htmx/views"
)

type Template = views.Template

const (
	TemplateBody Template = "body.html"

	TemplateSampleContent Template = "home-sample.html"
	TemplateReload        Template = "reload.html"
)
