package errorslice

import (
	"strings"
)

const (
	NL  = "\n"
	NLT = "\n\t"
)

func New(errs ...error) Errors {
	e := Errors{}
	if len(errs) > 0 {
		e.Join(errs...)
	}
	return e
}

// Errors 多个错误信息
type Errors []error

func (e Errors) Error() string {
	return e.Stringify(NL)
}

func (e Errors) ErrorTab() string {
	return e.Stringify(NLT)
}

func (e Errors) Stringify(separator string) string {
	s := make([]string, len(e))
	for k, v := range e {
		s[k] = v.Error()
	}
	return strings.Join(s, separator)
}

func (e Errors) String() string {
	return e.Error()
}

func (e Errors) IsEmpty() bool {
	return len(e) == 0
}

func (e *Errors) Join(errs ...error) *Errors {
	for _, err := range errs {
		if err != nil {
			*e = append(*e, err)
		}
	}
	return e
}

func (e *Errors) Add(err error) {
	if err != nil {
		*e = append(*e, err)
	}
}

func (e Errors) Unwrap() []error {
	return []error(e)
}

func (e Errors) ToError() error {
	if e.IsEmpty() {
		return nil
	}
	return e
}
