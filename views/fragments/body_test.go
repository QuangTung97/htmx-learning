package fragments

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	"htmx/pkg/route"
)

func newGoldie(t *testing.T) *goldie.Goldie {
	return goldie.New(t,
		goldie.WithFixtureDir("testdata"),
		goldie.WithNameSuffix(".html"),
	)
}

func newContext(buf *bytes.Buffer) route.Context {
	w := httptest.NewRecorder()
	w.Body = buf
	return route.NewContext(
		w,
		httptest.NewRequest(http.MethodGet, "/test", nil),
	)
}

func TestRenderBody_WithSampleContent(t *testing.T) {
	var buf bytes.Buffer
	ctx := newContext(&buf)

	err := RenderBodyWithSampleContent(ctx, false, 0)
	assert.Equal(t, nil, err)

	g := newGoldie(t)
	g.Assert(t, "full", buf.Bytes())
}

func TestRenderBody_WithSampleContent_Logged_In(t *testing.T) {
	var buf bytes.Buffer
	ctx := newContext(&buf)

	err := RenderBodyWithSampleContent(ctx, true, 11)
	assert.Equal(t, nil, err)

	g := newGoldie(t)
	g.Assert(t, "full-logged-in", buf.Bytes())
}

func TestRenderBody_WithSampleContent_With_HXRequest(t *testing.T) {
	t.Run("logged in", func(t *testing.T) {
		var buf bytes.Buffer
		ctx := newContext(&buf)
		ctx.SetHXRequestHeader()

		err := RenderBodyWithSampleContent(ctx, true, 11)
		assert.Equal(t, nil, err)

		g := newGoldie(t)
		g.Assert(t, "body-logged-in", buf.Bytes())
	})
}