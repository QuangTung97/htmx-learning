package route

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/QuangTung97/svloc"

	"htmx/views/routes"
)

//go:generate moq -rm -out error_mocks.go . ErrorView

type ErrorView interface {
	Render(ctx Context)
	Redirect(ctx Context, err error)
}

var ErrorViewLoc = svloc.Register[ErrorView](func(unv *svloc.Universe) ErrorView {
	return &errorViewImpl{}
})

type errorViewImpl struct {
}

func (v *errorViewImpl) Redirect(ctx Context, err error) {
	if err == nil {
		return
	}

	errorURL := routes.Error + "?msg=" + url.QueryEscape(err.Error())

	if !ctx.HasHxRequestHeader() {
		http.Redirect(ctx.Writer, ctx.Req, errorURL, http.StatusTemporaryRedirect)
		return
	}

	ctx.Writer.Header().Set("Hx-Reswap", "innerHTML")
	ctx.Writer.Header().Set("Hx-Retarget", "#body")
	ctx.Writer.Header().Set(hxPushURLHeader, errorURL)

	v.renderWithMsg(ctx, err.Error())
}

func (v *errorViewImpl) renderWithMsg(ctx Context, msg string) {
	type templateData struct {
		HomeURL string
		Text    string
	}
	_ = ctx.View("common/error.html", templateData{
		HomeURL: routes.Home,
		Text:    fmt.Sprintf("Error: %s", msg),
	})
}

func (v *errorViewImpl) Render(ctx Context) {
	msg := ctx.Req.URL.Query().Get("msg")
	v.renderWithMsg(ctx, msg)
}
