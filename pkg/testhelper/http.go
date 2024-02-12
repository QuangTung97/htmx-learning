package testhelper

import (
	"context"
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

func (h *HTTPTest) NewResponse() *httptest.ResponseRecorder {
	h.Writer = httptest.NewRecorder()
	return h.Writer
}

func (h *HTTPTest) NewContext() route.Context {
	if h.Req == nil {
		panic("Not yet init request object")
	}
	return route.Context{
		Ctx:    context.Background(),
		Req:    h.Req,
		Writer: h.NewResponse(),
	}
}
