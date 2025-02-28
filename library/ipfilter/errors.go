package ipfilter

import "errors"

var (
	ErrStartAndEndIPMismatchType = errors.New(`inconsistency between start and end ip types`)
	ErrParseIPAddress            = errors.New(`failed to parse ip`)
	ErrParseStartIPAddress       = errors.New(`failed to parse ip`)
	ErrParseEndIPAddress         = errors.New(`failed to parse ip`)
)
