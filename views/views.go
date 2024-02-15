package views

import (
	"bytes"
	"html/template"
	"io"
)

func Execute(w io.Writer, templateName Template, data any) error {
	return getTemplates().ExecuteTemplate(w, string(templateName), data)
}

func ExecuteHTML(templateName Template, data any) (template.HTML, error) {
	var buf bytes.Buffer
	if err := Execute(&buf, templateName, data); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func View(w io.Writer, body template.HTML) error {
	type templateData struct {
		Body template.HTML
	}
	return Execute(w, "main.html", templateData{
		Body: body,
	})
}

type Template string

const (
	TemplateBody Template = "body.html"

	TemplateLogin Template = "auth/google-login.html"

	TemplateError Template = "common/error.html"
)

type BodyData struct {
	LoggedIn bool
	Table    []string
}
