package cligo

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

type AppOption func(options *App)

func WithErrorOnHelp() AppOption {
	return func(app *App) {
		app.returnErrorOnHelp = true
	}
}

type UsageFunc func()

// An App serves as the main state for a cligo argument parser
type App struct {
	options           []*Option
	groups            map[string][]*Option
	usageFunc         UsageFunc
	returnErrorOnHelp bool
	name              string
	description       string
}

// NewApp returns a new instance of the App type
func NewApp(opts ...AppOption) *App {
	app := &App{
		groups: make(map[string][]*Option),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) SetName(name string) {
	a.name = name
}

func (a *App) SetDescription(description string) {
	a.description = description
}

// SetUsageFunc sets the function to call in order to print the usage string
// setting a nil usage function (the default) will result in the default usage function
// being used
func (a *App) SetUsageFunc(f UsageFunc) {
	a.usageFunc = f
}

// Usage prints the usage string for the application
func (a App) Usage() {

	if a.usageFunc != nil {
		a.usageFunc()
		return
	}

	if a.name != "" {
		fmt.Printf("Usage: %s %s [OPTIONS]", os.Args[0], a.name)
	} else {
		fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
	}

	positionals := filterFunc(a.options, func(opt *Option) bool {
		return opt.IsPositional()
	})

	for _, opt := range positionals {
		fmt.Printf(" %s", opt.pName)
	}

	fmt.Println("")
	if len(positionals) != 0 {
		fmt.Println("")
		fmt.Println("Positionals:")
		for _, opt := range positionals {
			fmt.Println(opt.formatPositional())
		}
	}

	for groupName, group := range a.groups {
		fmt.Println("")
		fmt.Printf("%s:\n", groupName)

		if groupName == "Options" {
			fmt.Printf("  %-30s %s\n", "-h,--help", "Print this help message and exit")
		}

		for _, opt := range group {
			if !opt.IsPositionalOnly() {
				fmt.Println(opt.format())
			}
		}
	}
}

func setOption(ptr any, v string, isNegated bool) error {
	// NOTE(eteran): isNegated is here for consistency of function definition,
	// but only flags can be negated
	_ = isNegated

	if ptr != nil {
		if err := setValue(ptr, v); err != nil {
			return err
		}
	}
	return nil
}

// AddOption adds a new option to the application and returns a pointer to an Option representing
// the new option so that it can be referred to by modifiers.
//
//   - name is a comma separated list of long and short options and can also include negations.
//     For example: "-a,--alpha,!--no-alpha" where --no-alpha will unset a previous --alpha or -a.
//   - ptr is a pointer to the variable which will receive the value of the option.
//   - help is the help string to use when printing the usage string
//   - modifiers is zero or more modifier functions, which can add additional rules to the parameter
func (a *App) AddOption(name string, ptr any, help string, modifiers ...Modifier) *Option {

	opt := NewOption(name, ptr, help, modifiers...)
	opt.owner = a
	a.options = append(a.options, opt)
	a.groups[opt.group] = append(a.groups[opt.group], opt)
	return opt
}

func setFlag(ptr any, v string, isNegated bool) error {
	if ptr != nil {
		if v == "" {
			if err := incrementFlag(ptr, isNegated); err != nil {
				return err
			}
		} else {
			if err := setValue(ptr, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddFlag adds a new flag to the application and returns a pointer to an Option representing
// the new flag so that it can be referred to by modifiers.
//
//   - name is a comma separated list of long and short options and can also include negations.
//     For example: "-a,--alpha,!--no-alpha" where --no-alpha will unset a previous --alpha or -a.
//   - ptr is a pointer to the variable which will receive the value of the option.
//   - help is the help string to use when printing the usage string
//   - modifiers is zero or more modifier functions, which can add additional rules to the parameter
//
// The main difference between a flag and an option is that flags are typically either an
// integer or a bool whose value "increases" which each successive usage. For example:
//
//	./my_app -v -v -v
//
// would result in the v flag having a value of 3 (assuming that it is bound to an integer)
func (a *App) AddFlag(name string, ptr any, help string, modifiers ...Modifier) *Option {

	opt := NewFlag(name, ptr, help, modifiers...)
	opt.owner = a
	a.options = append(a.options, opt)
	a.groups[opt.group] = append(a.groups[opt.group], opt)
	return opt
}

func (a App) findLongOption(name string) (opt *Option, isNegated bool, exists bool) {

	for _, opt := range a.options {
		if slices.Contains(opt.lNames, name) {
			return opt, false, true
		}

		if slices.Contains(opt.lNamesNeg, name) {
			return opt, true, true
		}

		if opt.ignoreCase {
			lName := strings.ToLower(name)
			if slices.ContainsFunc(opt.lNames, func(str string) bool {
				return lName == strings.ToLower(str)
			}) {
				return opt, false, true
			}

			if slices.ContainsFunc(opt.lNamesNeg, func(str string) bool {
				return lName == strings.ToLower(str)
			}) {
				return opt, true, true
			}
		}
	}

	return nil, false, false
}

func (a App) parseOneLong(arg string, args []string) ([]string, error) {
	/*
		--file filename (space)
		--file=filename (equals)
		--long_flag=true (long flag with equals to override default value)
		--long (long flag)
	*/
	name := arg[2:]
	param := ""

	if strings.Contains(name, "=") {
		parts := strings.SplitN(name, "=", 2)
		name = parts[0]
		param = parts[1]
	}

	opt, isNegated, exists := a.findLongOption(name)
	if !exists {
		return nil, fmt.Errorf("the following argument was not expected: %s\n%s", arg, ErrorSuffix)
	}

	if opt.isFlag {
		if err := opt.setter(opt, param, isNegated); err != nil {
			return nil, err
		}
	} else {
		if param == "" {
			if len(args) == 0 {
				return nil, ErrMissingParameter
			}
			param = args[0]
			args = args[1:]
		}

		if err := opt.setter(opt, param, isNegated); err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (a App) findShortOption(name string) (opt *Option, isNegated bool, exists bool) {

	for _, opt := range a.options {
		if slices.Contains(opt.sNames, name) {
			return opt, false, true
		}

		if slices.Contains(opt.sNamesNeg, name) {
			return opt, true, true
		}

		if opt.ignoreCase {
			lName := strings.ToLower(name)
			if slices.ContainsFunc(opt.sNames, func(str string) bool {
				return lName == strings.ToLower(str)
			}) {
				return opt, false, true
			}

			if slices.ContainsFunc(opt.sNamesNeg, func(str string) bool {
				return lName == strings.ToLower(str)
			}) {
				return opt, true, true
			}
		}
	}

	return nil, false, false
}

func (a App) parseOneShort(arg string, args []string) ([]string, error) {
	/*
		-a (flag)
		-f filename (option)
		-ffilename (no space required)
		-abc (flags can be combined)
		-abcf filename (flags and option can be combined)
	*/

	name := arg[1:]

	for i, ch := range name {
		shortName := string(ch)

		opt, isNegated, exists := a.findShortOption(shortName)
		if !exists {
			return nil, fmt.Errorf("the following argument was not expected: %s\n%s", arg, ErrorSuffix)
		}

		isLast := i == len(name)-1
		if opt.isFlag {
			if err := opt.setter(opt, "", isNegated); err != nil {
				return nil, err
			}
		} else if isLast {
			if len(args) == 0 {
				return nil, ErrMissingParameter
			}

			param := args[0]
			args = args[1:]

			if err := opt.setter(opt, param, isNegated); err != nil {
				return nil, err
			}
		} else {
			param := name[1+i:]
			if err := opt.setter(opt, param, isNegated); err != nil {
				return nil, err
			}
			break
		}
	}

	return args, nil
}

func (a App) parseOne(args []string) ([]string, error) {
	arg := args[0]

	var err error
	switch {
	case arg == "-h" || arg == "--help":
		a.Usage()
		return args, ErrHelpRequested
	case arg == "--":
		args = args[1:]
		return args, ErrEndOfArguments
	case strings.HasPrefix(arg, "--"):
		args = args[1:]
		args, err = a.parseOneLong(arg, args)
		if err != nil {
			return nil, err
		}
	case strings.HasPrefix(arg, "-"):
		args = args[1:]
		args, err = a.parseOneShort(arg, args)
		if err != nil {
			return nil, err
		}
	default:
		return args, ErrEndOfArguments
	}

	return args, nil
}

func (a App) parsePositional(args []string) ([]string, error) {

	for _, opt := range a.options {

		if !opt.IsPositional() {
			continue
		}

		if opt.Exists() {
			continue
		}

		if len(args) == 0 {
			break
		}

		if err := opt.setter(opt, args[0], false); err != nil {
			return nil, err
		}

		args = args[1:]
	}

	return args, nil
}

// ParseStrict will parse os.Args strictly. This means that unexpected positional arguments
// are considered an error. It equivalent to calling:
//
//	ParseArgsStrict(os.Args[1:])
func (a App) ParseStrict() error {
	return a.ParseArgsStrict(os.Args[1:])
}

// ParseArgsStrict will parse the string slice args strictly.
// This means that unexpected positional arguments are considered an error.
func (a App) ParseArgsStrict(args []string) error {
	rest, err := a.ParseArgs(args)
	if err != nil {
		return err
	}

	if len(rest) != 0 {
		return fmt.Errorf("the following arguments were not expected: %s\n%s", rest, ErrorSuffix)
	}

	return nil
}

// ParseStrict will parse os.Args.
// It equivalent to calling:
//
//	ParseArgs(os.Args[1:])
func (a App) Parse() ([]string, error) {
	return a.ParseArgs(os.Args[1:])
}

// ParseArgs will parse the string slice args and returns the unprocessed args as a new slice.
func (a App) ParseArgs(args []string) ([]string, error) {
	var err error

	for len(args) > 0 {
		args, err = a.parseOne(args)
		if err != nil {
			if errors.Is(err, ErrEndOfArguments) {
				break
			}

			if errors.Is(err, ErrHelpRequested) {
				if a.returnErrorOnHelp {
					return nil, err
				}
				os.Exit(0)
			}

			return nil, err
		}
	}

	args, err = a.parsePositional(args)
	if err != nil {
		return nil, err
	}

	if err := a.validateOptions(); err != nil {
		return nil, err
	}

	return args, nil
}

func (a App) validateOptions() error {
	for _, opt := range a.options {
		if opt.isRequired && !opt.Exists() {
			return fmt.Errorf("%s is required\n%s", opt.canonicalName(), ErrorSuffix)
		}

		if opt.Exists() {
			for _, need := range opt.needs {
				if !need.Exists() {
					return fmt.Errorf("%s requires %s\n%s", opt.canonicalName(), need.canonicalName(), ErrorSuffix)
				}
			}

			for _, exclude := range opt.excludes {
				if exclude.Exists() {
					return fmt.Errorf("%s excludes %s\n%s", opt.canonicalName(), exclude.canonicalName(), ErrorSuffix)
				}
			}
		}
	}

	return nil
}

func pointerType(ptr any) string {
	switch ptr.(type) {
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		return " NUMBER"
	case *float32, *float64:
		return " REAL"
	case *string:
		return " TEXT"
	default:
		return ""
	}
}

func filterFunc[T any](slice []T, f func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, value := range slice {
		if f(value) {
			result = append(result, value)
		}
	}

	return result
}
