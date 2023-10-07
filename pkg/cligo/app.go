package cligo

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type app struct {
	allOptions           []*Option
	lOptions             map[string]*Option
	sOptions             map[string]*Option
	lLowerOptions        map[string]*Option
	sLowerOptions        map[string]*Option
	lNegatedOptions      map[string]*Option
	sNegatedOptions      map[string]*Option
	lNegatedLowerOptions map[string]*Option
	sNegatedLowerOptions map[string]*Option
	pOptions             []*Option
	groups               map[string][]*Option
}

func NewApp() app {
	return app{
		lOptions:             make(map[string]*Option),
		sOptions:             make(map[string]*Option),
		lLowerOptions:        make(map[string]*Option),
		sLowerOptions:        make(map[string]*Option),
		lNegatedOptions:      make(map[string]*Option),
		sNegatedOptions:      make(map[string]*Option),
		lNegatedLowerOptions: make(map[string]*Option),
		sNegatedLowerOptions: make(map[string]*Option),
		groups:               make(map[string][]*Option),
	}
}

func valueTypeString(opt *Option) string {
	switch opt.value.(type) {
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		return " NUMBER"
	case *string:
		return " TEXT"
	default:
		return ""
	}
}

func (a app) Usage() {
	fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
	for _, opt := range a.pOptions {
		fmt.Printf(" %s", opt.positionalName)
	}

	fmt.Println("")
	if len(a.pOptions) != 0 {
		fmt.Println("")
		fmt.Println("Positionals:")
		for _, opt := range a.pOptions {
			name := opt.positionalName
			name = name + valueTypeString(opt)

			if opt.defaultValue != "" {
				name = name + fmt.Sprintf(" [%s]", opt.defaultValue)
			}

			if opt.isRequired {
				name = name + " REQUIRED"
			}
			fmt.Printf("  %-30s %s\n", name, opt.help)
		}
	}

	for groupName, group := range a.groups {
		fmt.Println("")
		fmt.Printf("%s:\n", groupName)
		for _, opt := range group {
			if !opt.isPositionalOnly {
				names := strings.Join(append(opt.shortNames, opt.longNames...), ",")
				names = names + valueTypeString(opt)

				if opt.defaultValue != "" {
					names = names + fmt.Sprintf(" [%s]", opt.defaultValue)
				}

				if opt.isRequired {
					names = names + " REQUIRED"
				}
				fmt.Printf("  %-30s %s\n", names, opt.help)
			}
		}

		if groupName == "Options" {
			fmt.Printf("  %-30s %s\n", "-h,--help", "Print this help message and exit")
		}
	}

	os.Exit(0)
}

func (a *app) addLongOption(name string, opt *Option, isNegated bool) error {

	opt.longNames = append(opt.longNames, "--"+name)
	opt.isPositionalOnly = false
	lName := strings.ToLower(name)

	if isNegated {
		if _, exists := a.lNegatedOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.lNegatedOptions[name] = opt

		if _, exists := a.lNegatedLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.lNegatedLowerOptions[lName] = opt

	} else {
		if _, exists := a.lOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.lOptions[name] = opt

		if _, exists := a.lLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.lLowerOptions[lName] = opt

	}
	return nil
}

func (a *app) addShortOption(name string, opt *Option, isNegated bool) error {

	opt.shortNames = append(opt.shortNames, "-"+name)
	opt.isPositionalOnly = false
	lName := strings.ToLower(name)

	if isNegated {
		if _, exists := a.sNegatedOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.sNegatedOptions[name] = opt

		if _, exists := a.sNegatedLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.sNegatedLowerOptions[lName] = opt
	} else {

		if _, exists := a.sOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.sOptions[name] = opt

		if _, exists := a.sLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.sLowerOptions[lName] = opt
	}
	return nil
}

func (a *app) AddOption(name string, value any, help string, modifiers ...Modifier) *Option {
	if value == nil {
		panic("bound variables cannot be nil")
	}

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		panic("bound variables must be pointers")
	}

	opt := &Option{
		help:             help,
		value:            value,
		isFlag:           false,
		isPositionalOnly: true,
		group:            "Options",
	}

	for _, mod := range modifiers {
		mod(opt)
	}

	name = strings.TrimSpace(name)
	names := strings.Split(name, ",")
	for _, optionName := range names {

		if optionName == "" {
			panic("argument has empty name, do you have a trailing comma?")
		}

		isLong := strings.HasPrefix(optionName, "--")
		isShort := !isLong && strings.HasPrefix(optionName, "-")
		isPositional := !isLong && !isShort

		if isLong {
			lName := optionName[2:]
			if err := a.addLongOption(lName, opt, false); err != nil {
				panic(err)
			}
		} else if isShort {
			sName := optionName[1:]
			if err := a.addShortOption(sName, opt, false); err != nil {
				panic(err)
			}
		} else if isPositional {
			opt.positionalName = optionName
			a.pOptions = append(a.pOptions, opt)
		}
	}

	a.allOptions = append(a.allOptions, opt)
	a.groups[opt.group] = append(a.groups[opt.group], opt)
	return opt
}

func ensureIntegralPointer(value any) {
	ty := reflect.TypeOf(value)
	if ty.Kind() != reflect.Ptr {
		panic("bound variables must be pointers")
	}

	switch ty.Elem().Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		break
	default:
		panic("flags may only be boolean or integral types")
	}
}

func (a *app) AddFlag(name string, value any, help string, modifiers ...Modifier) *Option {

	if value == nil {
		panic("bound variables cannot be nil")
	}

	ensureIntegralPointer(value)

	opt := &Option{
		help:   help,
		value:  value,
		isFlag: true,
		group:  "Options",
	}

	for _, mod := range modifiers {
		mod(opt)
	}

	name = strings.TrimSpace(name)
	names := strings.Split(name, ",")
	for _, flagName := range names {

		if flagName == "" {
			panic("argument has empty name, do you have a trailing comma?")
		}

		isNegated := strings.HasPrefix(flagName, "!")
		if isNegated {
			flagName = flagName[1:]
		}

		isLong := strings.HasPrefix(flagName, "--")
		isShort := !isLong && strings.HasPrefix(flagName, "-")
		isPositional := !isLong && !isShort

		if isLong {
			lName := flagName[2:]
			if err := a.addLongOption(lName, opt, isNegated); err != nil {
				panic(err)
			}
		} else if isShort {
			sName := flagName[1:]
			if err := a.addShortOption(sName, opt, isNegated); err != nil {
				panic(err)
			}
		} else if isPositional {
			panic("positional arguments must not be flags")
		}
	}

	a.allOptions = append(a.allOptions, opt)
	a.groups[opt.group] = append(a.groups[opt.group], opt)
	return opt
}

func (a app) findLongOption(name string) (opt *Option, isNegated bool, exists bool) {
	if opt, exists := a.lOptions[name]; exists {
		return opt, false, true
	}

	if opt, exists := a.lNegatedOptions[name]; exists {
		return opt, true, true
	}

	lName := strings.ToLower(name)
	if opt, exists := a.lLowerOptions[lName]; exists && opt.ignoreCase {
		return opt, false, true
	}

	if opt, exists := a.lNegatedLowerOptions[lName]; exists && opt.ignoreCase {
		return opt, true, true
	}

	return nil, false, false
}

func (a app) parseOneLong(arg string, args []string) ([]string, error) {
	/*
		--file filename (space)
		--file=filename (equals)
		--long_flag=true (long flag with equals to override default value)
		--long (long flag)
	*/
	name := arg[2:]
	param := ""

	if strings.Contains(name, "=") {
		x := strings.SplitN(name, "=", 2)
		name = x[0]
		param = x[1]
	}

	opt, isNegated, exists := a.findLongOption(name)
	if !exists {
		return nil, fmt.Errorf("the following argument was not expected: %s\n%s", arg, ErrorSuffix)
	}

	if opt.isFlag {
		if param == "" {
			if err := opt.incrementFlag(isNegated); err != nil {
				return nil, err
			}
			return args, nil
		} else {
			if err := opt.setValue(param); err != nil {
				return nil, err
			}
			return args, nil
		}
	} else {
		if param == "" {
			if len(args) == 0 {
				return nil, ErrMissingParameter
			}
			param = args[0]
			args = args[1:]
		}

		if err := opt.setValue(param); err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (a app) findShortOption(name string) (opt *Option, isNegated bool, exists bool) {
	if opt, exists := a.sOptions[name]; exists {
		return opt, false, true
	}

	if opt, exists := a.sNegatedOptions[name]; exists {
		return opt, true, true
	}

	lName := strings.ToLower(name)
	if opt, exists := a.sLowerOptions[lName]; exists && opt.ignoreCase {
		return opt, false, true
	}

	if opt, exists := a.sNegatedLowerOptions[lName]; exists && opt.ignoreCase {
		return opt, true, true
	}

	return nil, false, false
}

func (a app) parseOneShort(arg string, args []string) ([]string, error) {
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
			if err := opt.incrementFlag(isNegated); err != nil {
				return nil, err
			}
		} else if isLast {
			if len(args) == 0 {
				return nil, ErrMissingParameter
			}

			param := args[0]
			args = args[1:]

			if err := opt.setValue(param); err != nil {
				return nil, err
			}
		} else {
			param := name[1+i:]
			if err := opt.setValue(param); err != nil {
				return nil, err
			}
			break
		}
	}

	return args, nil
}

func (a app) parseOne(args []string) ([]string, error) {
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

func (a app) parsePositional(args []string) ([]string, error) {

	for _, opt := range a.pOptions {

		if opt.exists {
			continue
		}

		if len(args) == 0 {
			break
		}

		if err := opt.setValue(args[0]); err != nil {
			return nil, err
		}
		args = args[1:]
	}

	return args, nil
}

func (a app) ParseStrict() error {
	return a.ParseArgsStrict(os.Args[1:])
}

func (a app) ParseArgsStrict(args []string) error {
	rest, err := a.ParseArgs(args)
	if err != nil {
		return err
	}

	if len(rest) != 0 {
		return fmt.Errorf("the following arguments were not expected: %s\n%s", rest, ErrorSuffix)
	}

	return nil
}

func (a app) Parse() ([]string, error) {
	return a.ParseArgs(os.Args[1:])
}

func (a app) ParseArgs(args []string) ([]string, error) {
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

	for _, opt := range a.allOptions {
		if opt.isRequired && !opt.exists {
			return nil, fmt.Errorf("%s is required\n%s", opt.canonicalName(), ErrorSuffix)
		}

		for _, need := range opt.needs {
			if !need.exists {
				return nil, fmt.Errorf("%s requires %s\n%s", opt.canonicalName(), need.canonicalName(), ErrorSuffix)
			}
		}

		for _, exclude := range opt.excludes {
			if exclude.exists {
				return nil, fmt.Errorf("%s excludes %s\n%s", opt.canonicalName(), exclude.canonicalName(), ErrorSuffix)
			}
		}
	}

	return args, nil
}
