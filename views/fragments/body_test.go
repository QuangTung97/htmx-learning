package fragments

import (
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	"htmx/views/viewtest"
)

func newGoldie(t *testing.T) *goldie.Goldie {
	return goldie.New(t,
		goldie.WithFixtureDir("testdata"),
		goldie.WithNameSuffix(".html"),
	)
}

func TestRenderBody_WithSampleContent(t *testing.T) {
	v := viewtest.New(t)

	err := RenderBodyWithSampleContent(v.Ctx, false, 0)
	assert.Equal(t, nil, err)

	v.Assert("full")
}

func TestRenderBody_WithSampleContent_Logged_In(t *testing.T) {
	v := viewtest.New(t)

	err := RenderBodyWithSampleContent(v.Ctx, true, 11)
	assert.Equal(t, nil, err)

	v.Assert("full-logged-in")
}

func TestRenderBody_WithSampleContent_With_HXRequest(t *testing.T) {
	t.Run("logged in", func(t *testing.T) {
		v := viewtest.New(t)
		v.Ctx.SetHXRequestHeader()

		err := RenderBodyWithSampleContent(v.Ctx, true, 11)
		assert.Equal(t, nil, err)

		v.Assert("body-logged-in")
	})
}
