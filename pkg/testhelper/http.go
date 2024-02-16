package testhelper

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"htmx/pkg/route"
)

type HTTPTest struct {
	Req    *http.Request
	Writer *httptest.ResponseRecorder
}

func NewHTTPTest() *HTTPTest {
	return &HTTPTest{}
}

func (h *HTTPTest) NewGet(urlPath string) {
	h.Req = httptest.NewRequest(http.MethodGet, urlPath, nil)
}

func (h *HTTPTest) NewPost(urlPath string, body string) {
	var buf bytes.Buffer
	buf.WriteString(body)
	h.Req = httptest.NewRequest(http.MethodPost, urlPath, &buf)
}

func (h *HTTPTest) newResponse() *httptest.ResponseRecorder {
	h.Writer = httptest.NewRecorder()
	return h.Writer
}

func (h *HTTPTest) NewContext() route.Context {
	if h.Req == nil {
		panic("Not yet init request object")
	}
	return route.NewContext(h.newResponse(), h.Req)
}
