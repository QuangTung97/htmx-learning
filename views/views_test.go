package views

import (
	"bytes"
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

func TestTemplate_Full(t *testing.T) {
	body, err := ExecuteHTML(TemplateBody, BodyData{})
	assert.Equal(t, nil, err)

	var buf bytes.Buffer
	err = View(&buf, body)

	assert.Equal(t, nil, err)
	g := newGoldie(t)
	g.Assert(t, "full", buf.Bytes())
}

func TestTemplate_Full_Logged_In(t *testing.T) {
	body, err := ExecuteHTML(TemplateBody, BodyData{
		LoggedIn: true,
	})
	assert.Equal(t, nil, err)

	var buf bytes.Buffer
	err = View(&buf, body)

	assert.Equal(t, nil, err)
	g := newGoldie(t)
	g.Assert(t, "full-logged-in", buf.Bytes())
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
