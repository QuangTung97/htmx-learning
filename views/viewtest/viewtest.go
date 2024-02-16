package viewtest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"

	"htmx/pkg/route"
)

type ViewTest struct {
	Ctx route.Context

	t   *testing.T
	buf *bytes.Buffer
	g   *goldie.Goldie
}

func New(t *testing.T) *ViewTest {
	buf := &bytes.Buffer{}
	return &ViewTest{
		Ctx: newContext(buf),

		t:   t,
		buf: buf,
		g:   newGoldie(t),
	}
}

func (v *ViewTest) Assert(name string) {
	if os.Getenv("GOLDEN_UPDATE") != "" {
		err := v.g.Update(v.t, name, v.buf.Bytes())
		if err != nil {
			panic(err)
		}
	} else {
		v.g.Assert(v.t, name, v.buf.Bytes())
	}
}

func newContext(buf *bytes.Buffer) route.Context {
	w := httptest.NewRecorder()
	w.Body = buf
	return route.NewContext(
		w,
		httptest.NewRequest(http.MethodGet, "/test-url", nil),
	)
}

func newGoldie(t *testing.T) *goldie.Goldie {
	return goldie.New(t,
		goldie.WithFixtureDir("testdata"),
		goldie.WithNameSuffix(".html"),
	)
}
