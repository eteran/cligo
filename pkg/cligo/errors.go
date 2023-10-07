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
	ErrMissingParameter = errors.New("missing parameter")
	ErrUnsupportedType  = errors.New("unsupported bound variable type")
	ErrEndOfArguments   = errors.New("end of arguments")
	ErrDuplicateOption  = errors.New("duplicate option")
)
