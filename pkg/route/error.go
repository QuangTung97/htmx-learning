package route

import (
	"fmt"
	"log"

	"github.com/QuangTung97/svloc"

	"htmx/views/routes"
)

type ErrorView interface {
	Render(ctx Context, err error)
}

var ErrorViewLoc = svloc.Register[ErrorView](func(unv *svloc.Universe) ErrorView {
	return &errorViewImpl{}
})

type errorViewImpl struct {
}

func (v *errorViewImpl) Render(ctx Context, err error) {
	type templateData struct {
		HomeURL string
		Text    string
	}
	newErr := ctx.View("common/error.html", templateData{
		HomeURL: routes.Home,
		Text:    fmt.Sprintf("Error: %s", err.Error()),
	})
	if newErr != nil {
		log.Println("[ERROR]", newErr)
	}

	if ctx.IsHxRequest() {
		ctx.Writer.Header().Set("Hx-Reswap", "innerHTML")
		ctx.Writer.Header().Set("Hx-Retarget", "#body")
	}
}
