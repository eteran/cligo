package cligo

import (
	"errors"
)

const (
	ErrorSuffix = "run with --help for more information."
)

var (
	ErrMissingParameter = errors.New("missing parameter")
	ErrUnsupportedType  = errors.New("unsupported bound variable type")
	ErrEndOfArguments   = errors.New("end of arguments")
	ErrDuplicateOption  = errors.New("duplicate option")
)
