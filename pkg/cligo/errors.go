package cligo

import (
	"errors"
	"fmt"
)

type ErrUnexpectedArgument struct {
	arg string
}

func (err ErrUnexpectedArgument) Error() string {
	return fmt.Sprintf("The following argument was not expected: %s\nRun with --help for more information.", err.arg)
}

type ErrUnexpectedArguments struct {
	args []string
}

func (err ErrUnexpectedArguments) Error() string {
	return fmt.Sprintf("The following arguments were not expected: %s\nRun with --help for more information.", err.args)
}

type ErrMissingRequiredArgument struct {
	name string
}

func (err ErrMissingRequiredArgument) Error() string {
	return fmt.Sprintf("%s is required\nRun with --help for more information.", err.name)
}

var (
	ErrMissingParameter       = errors.New("missing parameter")
	ErrEmptyName              = errors.New("argument has empty name, do you have a trailing comma?")
	ErrNotPointer             = errors.New("bound variables must be pointers")
	ErrNilBoundVariable       = errors.New("bound variables cannot be nil")
	ErrNotFlagType            = errors.New("flags may only be boolean or integral types")
	ErrUnsupportedType        = errors.New("unsupported bound variable type")
	ErrUnexpectedIntegerWidth = errors.New("(u)int is not 64 or 32 bits)")
	ErrEndOfArguments         = errors.New("end of arguments")
	ErrInvalidPositional      = errors.New("positional arguments must not be flags")
	ErrDuplicateOption        = errors.New("duplicate option")
)
