package views

import (
	"bytes"
	"html/template"
	"io"
)

func Execute(w io.Writer, templateName string, data any) error {
	return getTemplates().ExecuteTemplate(w, templateName, data)
}

func ExecuteHTML(templateName string, data any) (template.HTML, error) {
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
