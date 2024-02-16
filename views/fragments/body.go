package fragments

import (
	"html/template"

	"htmx/pkg/route"
)

type BodyData struct {
	LoggedIn bool
	Content  template.HTML
}

func RenderBody(ctx route.Context, data BodyData) error {
	return ctx.View(TemplateBody, data)
}

func RenderBodyWithSampleContent(ctx route.Context, loggedIn bool, count int64) error {
	reloadHTML, err := ReloadData{Count: count}.RenderHTML(ctx)
	if err != nil {
		return err
	}

	sampleData := SampleContentData{
		LoggedIn: loggedIn,
		Reload:   reloadHTML,
		Table: []string{
			"a", "b", "c", "d", "e",
		},
	}

	content, err := sampleData.RenderHTML(ctx)
	if err != nil {
		return err
	}

	return RenderBody(ctx, BodyData{
		LoggedIn: loggedIn,
		Content:  content,
	})
}
