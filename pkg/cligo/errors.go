package cligo

import (
	"errors"
	"fmt"
)

const (
	ErrorSuffix = "Run with --help for more information."
)

type ErrUnexpectedArgument struct {
	arg string
}

func (err ErrUnexpectedArgument) Error() string {
	return fmt.Sprintf("The following argument was not expected: %s\n%s", err.arg, ErrorSuffix)
}

type ErrUnexpectedArguments struct {
	args []string
}

func (err ErrUnexpectedArguments) Error() string {
	return fmt.Sprintf("The following arguments were not expected: %s\n%s", err.args, ErrorSuffix)
}

type ErrMissingRequiredArgument struct {
	name string
}

func (err ErrMissingRequiredArgument) Error() string {
	return fmt.Sprintf("%s is required\n%s", err.name, ErrorSuffix)
}

type ErrMissingRequiredOption struct {
	opt   string
	needs string
}

func (err ErrMissingRequiredOption) Error() string {
	return fmt.Sprintf("%s requires %s\n%s", err.opt, err.needs, ErrorSuffix)
}

type ErrConflictingOption struct {
	opt      string
	conflict string
}

func (err ErrConflictingOption) Error() string {
	return fmt.Sprintf("%s excludes %s\n%s", err.opt, err.conflict, ErrorSuffix)
}

var (
	ErrMissingParameter = errors.New("missing parameter")
	ErrUnsupportedType  = errors.New("unsupported bound variable type")
	ErrEndOfArguments   = errors.New("end of arguments")
	ErrDuplicateOption  = errors.New("duplicate option")
)
