package cligo

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

type App struct {
	options []*Option
	groups  map[string][]*Option
}

func NewApp() App {
	return App{
		groups: make(map[string][]*Option),
	}
}

func pointerType(ptr any) string {
	switch ptr.(type) {
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		return " NUMBER"
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

func (a App) Usage() {
	fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])

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

	os.Exit(0)
}

func setOption(ptr any, v string, isNegated bool) error {
	if ptr != nil {
		if err := setValue(ptr, v); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) AddOption(name string, ptr any, help string, modifiers ...Modifier) *Option {

	opt := NewOption(name, ptr, help, modifiers...)
	opt.app = a
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

func (a *App) AddFlag(name string, ptr any, help string, modifiers ...Modifier) *Option {

	opt := NewFlag(name, ptr, help, modifiers...)
	opt.app = a
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
	if len(args) == 0 {
		return nil, errors.New("no arguments")
	}

	arg := args[0]

	var err error
	switch {
	case arg == "-h" || arg == "--help":
		a.Usage()
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

func (a App) ParseStrict() error {
	return a.ParseArgsStrict(os.Args[1:])
}

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

func (a App) Parse() ([]string, error) {
	return a.ParseArgs(os.Args[1:])
}

func (a App) ParseArgs(args []string) ([]string, error) {
	var err error

	for len(args) > 0 {
		args, err = a.parseOne(args)
		if err != nil {
			if errors.Is(err, ErrEndOfArguments) {
				break
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
