package errors

import (
	"fmt"

	"golang.org/x/text/language"
)

type TranslateFunc func(values map[string]any) string

type BaseError struct {
	code string
	msg  string // base message (for debugging)

	fields []string

	translate map[language.Base]TranslateFunc
}

func (e *BaseError) New(values map[string]any) error {
	return &DomainError{
		base:   e,
		values: values,
	}
}

type DomainError struct {
	base   *BaseError
	values map[string]any
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.base.code, e.base.msg)
}

func (e *DomainError) Translate(lang language.Tag) string {
	base, _ := lang.Base()
	return e.base.translate[base](e.values)
}

var _ error = &DomainError{}

func Register(code string, message string) *BaseError {
	return &BaseError{}
}

type JoinErrors struct {
	errors []*DomainError
	codes  map[string]struct{}
}

func (e *JoinErrors) Error() string {
	return ""
}

func Join(err error, newErrs ...error) error {
	return &JoinErrors{}
}

func Exist(err error, code string, filter map[string]any) (*DomainError, bool) {
	return nil, false
}
