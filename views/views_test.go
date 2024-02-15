package views

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func newGoldie(t *testing.T) *goldie.Goldie {
	return goldie.New(t,
		goldie.WithFixtureDir("testdata"),
		goldie.WithNameSuffix(".html"),
	)
}

func TestTemplate(t *testing.T) {
	tmpl := getTemplates()

	for _, subTmpl := range tmpl.Templates() {
		fmt.Println(subTmpl.Name())
	}

	body, err := ExecuteHTML(TemplateBody, BodyData{})
	assert.Equal(t, nil, err)

	var buf bytes.Buffer
	err = View(&buf, body)

	assert.Equal(t, nil, err)
	g := newGoldie(t)
	g.Assert(t, "full", buf.Bytes())
}

func TestTemplate_RenderBody(t *testing.T) {
	t.Run("logged in", func(t *testing.T) {
		body, err := ExecuteHTML(TemplateBody, BodyData{
			LoggedIn: true,
		})

		assert.Equal(t, nil, err)
		g := newGoldie(t)
		g.Assert(t, "body-logged-in", []byte(body))
	})
}
