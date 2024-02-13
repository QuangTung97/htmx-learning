package route

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

type errorViewTest struct {
	view   ErrorView
	ctx    Context
	writer *httptest.ResponseRecorder
}

func newErrorViewTest() *errorViewTest {
	req := httptest.NewRequest(http.MethodGet, "/test-url", nil)
	writer := httptest.NewRecorder()

	return &errorViewTest{
		view: &errorViewImpl{},
		ctx: Context{
			Ctx:    context.Background(),
			Req:    req,
			Writer: writer,
		},
		writer: writer,
	}
}

func TestErrorView(t *testing.T) {
	t.Run("render", func(t *testing.T) {
		v := newErrorViewTest()
		v.ctx.Req.URL = &url.URL{
			RawQuery: "msg=" + url.QueryEscape("some error"),
		}

		v.view.Render(v.ctx)

		g := goldie.New(t,
			goldie.WithFixtureDir("testdata"),
			goldie.WithNameSuffix(".html"),
		)
		g.Assert(t, "error", v.writer.Body.Bytes())

		assert.Equal(t, http.Header{
			"Content-Type": {"text/html; charset=utf-8"},
		}, v.writer.Header())
	})

	t.Run("redirect normal", func(t *testing.T) {
		v := newErrorViewTest()

		v.view.Redirect(v.ctx, errors.New("another error"))

		assert.Equal(t, http.Header{
			"Content-Type": {"text/html; charset=utf-8"},
			"Location":     {"/error?msg=another+error"},
		}, v.writer.Header())
	})

	t.Run("redirect with hx request", func(t *testing.T) {
		v := newErrorViewTest()

		v.ctx.Req.Header.Add("HX-Request", "true")

		v.view.Redirect(v.ctx, errors.New("some error"))

		g := goldie.New(t,
			goldie.WithFixtureDir("testdata"),
			goldie.WithNameSuffix(".html"),
		)
		g.Assert(t, "error-hx", v.writer.Body.Bytes())

		assert.Equal(t, http.Header{
			"Content-Type": {"text/html; charset=utf-8"},
			"Hx-Push-Url":  {"/error?msg=some+error"},
			"Hx-Reswap":    {"innerHTML"},
			"Hx-Retarget":  {"#body"},
		}, v.writer.Header())
	})
}
