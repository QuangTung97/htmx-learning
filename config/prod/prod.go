package prod

import (
	"github.com/QuangTung97/svloc"

	"htmx/config"
)

func NewUniverse() *svloc.Universe {
	unv := svloc.NewUniverse()
	config.Loc.MustOverrideFunc(unv, func(unv *svloc.Universe) config.Config {
		return config.Load()
	})
	return unv
}
