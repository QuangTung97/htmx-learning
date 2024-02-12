package translate

import (
	"golang.org/x/text/language"
)

type Func func(values map[string]any) string

func SimpleString(s string) Func {
	return func(values map[string]any) string {
		return s
	}
}

type Text struct {
	field     []string
	defaultFn Func
	lang      map[language.Base]Func
}

func (t *Text) Translate(tag language.Tag, values map[string]any) string {
	return ""
}

type newConfig struct {
	lang map[language.Base]Func
}

type Option func(conf *newConfig)

func NewText(
	fields []string,
	defaultFunc Func,
	options ...Option,
) *Text {
	conf := newConfig{
		lang: map[language.Base]Func{},
	}
	return &Text{
		field:     fields,
		defaultFn: defaultFunc,
		lang:      conf.lang,
	}
}
