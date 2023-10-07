package cligo

import (
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
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

func (a app) Usage() {
	fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
	for _, pOption := range a.pOptions {
		fmt.Printf(" %s", pOption.positionalName)
	}

	fmt.Printf("\n")
	if len(a.pOptions) != 0 {
		fmt.Println("")
		fmt.Println("Positionals:")
		for _, pOption := range a.pOptions {
			name := pOption.positionalName
			if pOption.isRequired {
				name = name + " REQUIRED"
			}
			fmt.Printf("  %-30s %s\n", name, pOption.help)
		}
	}

	for name, group := range a.groups {
		fmt.Println("")
		fmt.Printf("%s:\n", name)
		for _, option := range group {
			if !option.isPositionalOnly {
				names := strings.Join(append(option.shortNames, option.longNames...), ",")
				if option.isRequired {
					names = names + " REQUIRED"
				}
				fmt.Printf("  %-30s %s\n", names, option.help)
			}
		}

		if name == "Options" {
			fmt.Printf("  %-30s %s\n", "-h,--help", "Print this help message and exit")
		}
	}

	os.Exit(0)
}

func (a *app) addLongOption(name string, option *Option, isNegated bool) error {

	option.longNames = append(option.longNames, "--"+name)
	option.isPositionalOnly = false
	lName := strings.ToLower(name)

	if isNegated {
		if _, exists := a.lNegatedOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.lNegatedOptions[name] = option

		if _, exists := a.lNegatedLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.lNegatedLowerOptions[lName] = option

	} else {
		if _, exists := a.lOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.lOptions[name] = option

		if _, exists := a.lLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.lLowerOptions[lName] = option

	}
	return nil
}

func (a *app) addShortOption(name string, option *Option, isNegated bool) error {

	option.shortNames = append(option.shortNames, "-"+name)
	option.isPositionalOnly = false
	lName := strings.ToLower(name)

	if isNegated {
		if _, exists := a.sNegatedOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.sNegatedOptions[name] = option

		if _, exists := a.sNegatedLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.sNegatedLowerOptions[lName] = option
	} else {

		if _, exists := a.sOptions[name]; exists {
			return ErrDuplicateOption
		}
		a.sOptions[name] = option

		if _, exists := a.sLowerOptions[lName]; exists {
			return ErrDuplicateOption
		}
		a.sLowerOptions[lName] = option
	}
	return nil
}

func (a *app) AddOption(name string, value any, help string, modifiers ...Modifier) *Option {
	if value == nil {
		panic(ErrNilBoundVariable)
	}

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		panic(ErrNotPointer)
	}

	option := &Option{
		help:             help,
		value:            value,
		isFlag:           false,
		isPositionalOnly: true,
		group:            "Options",
	}

	for _, mod := range modifiers {
		mod(option)
	}

	name = strings.TrimSpace(name)
	names := strings.Split(name, ",")
	for _, optionName := range names {

		if optionName == "" {
			panic(ErrEmptyName)
		}

		isLong := strings.HasPrefix(optionName, "--")
		isShort := !isLong && strings.HasPrefix(optionName, "-")
		isPositional := !isLong && !isShort

		if isLong {
			lName := optionName[2:]
			if err := a.addLongOption(lName, option, false); err != nil {
				panic(err)
			}
		} else if isShort {
			sName := optionName[1:]
			if err := a.addShortOption(sName, option, false); err != nil {
				panic(err)
			}
		} else if isPositional {
			option.positionalName = optionName
			a.pOptions = append(a.pOptions, option)
		}
	}

	a.allOptions = append(a.allOptions, option)
	a.groups[option.group] = append(a.groups[option.group], option)
	return option
}

func ensureIntegralPointer(value any) {
	ty := reflect.TypeOf(value)
	if ty.Kind() != reflect.Ptr {
		panic(ErrNotPointer)
	}

	switch ty.Elem().Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		break
	default:
		panic(ErrNotFlagType)
	}
}

func (a *app) AddFlag(name string, value any, help string, modifiers ...Modifier) *Option {

	if value == nil {
		panic(ErrNilBoundVariable)
	}

	ensureIntegralPointer(value)

	option := &Option{
		help:   help,
		value:  value,
		isFlag: true,
		group:  "Options",
	}

	for _, mod := range modifiers {
		mod(option)
	}

	name = strings.TrimSpace(name)
	names := strings.Split(name, ",")
	for _, flagName := range names {

		if flagName == "" {
			panic(ErrEmptyName)
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
			if err := a.addLongOption(lName, option, isNegated); err != nil {
				panic(err)
			}
		} else if isShort {
			sName := flagName[1:]
			if err := a.addShortOption(sName, option, isNegated); err != nil {
				panic(err)
			}
		} else if isPositional {
			panic(ErrInvalidPositional)
		}
	}

	a.allOptions = append(a.allOptions, option)
	a.groups[option.group] = append(a.groups[option.group], option)
	return option
}

func (a app) findLongOption(name string) (option *Option, isNegated bool, exists bool) {
	if option, exists := a.lOptions[name]; exists {
		return option, false, true
	}

	if option, exists := a.lNegatedOptions[name]; exists {
		return option, true, true
	}

	lName := strings.ToLower(name)
	if option, exists := a.lLowerOptions[lName]; exists && option.ignoreCase {
		return option, false, true
	}

	if option, exists := a.lNegatedLowerOptions[lName]; exists && option.ignoreCase {
		return option, true, true
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

	option, isNegated, exists := a.findLongOption(name)
	if !exists {
		return nil, ErrUnexpectedArgument{arg}
	}

	if option.isFlag {
		if param == "" {
			if err := incrementFlag(option, isNegated); err != nil {
				return nil, err
			}
			return args, nil
		} else {
			if err := setValue(option, param); err != nil {
				return nil, err
			}
			return args, nil
		}
	} else {
		if param == "" {
			if len(args) > 0 {
				param = args[0]
				args = args[1:]
			} else {
				return nil, ErrMissingParameter
			}
		}

		if err := setValue(option, param); err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (a app) findShortOption(name string) (option *Option, isNegated bool, exists bool) {
	if option, exists := a.sOptions[name]; exists {
		return option, false, true
	}

	if option, exists := a.sNegatedOptions[name]; exists {
		return option, true, true
	}

	lName := strings.ToLower(name)
	if option, exists := a.sLowerOptions[lName]; exists && option.ignoreCase {
		return option, false, true
	}

	if option, exists := a.sNegatedLowerOptions[lName]; exists && option.ignoreCase {
		return option, true, true
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

	for j := 0; j < len(name); j++ {
		shortName := string(name[j])

		option, isNegated, exists := a.findShortOption(shortName)
		if !exists {
			return nil, ErrUnexpectedArgument{arg}
		}

		isLast := j == len(name)-1
		if option.isFlag {
			if err := incrementFlag(option, isNegated); err != nil {
				return nil, err
			}
		} else if isLast {
			if len(args) > 0 {
				param := args[0]
				args = args[1:]

				if err := setValue(option, param); err != nil {
					return nil, err
				}
			} else {
				return nil, ErrMissingParameter
			}
		} else {
			param := name[1+j:]
			if err := setValue(option, param); err != nil {
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

	for _, pOption := range a.pOptions {

		if pOption.exists {
			continue
		}

		if len(args) == 0 {
			break
		}

		if err := setValue(pOption, args[0]); err != nil {
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
		return ErrUnexpectedArguments{rest}
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

	for _, option := range a.allOptions {
		if option.isRequired && !option.exists {
			return nil, ErrMissingRequiredArgument{option.canonicalName()}
		}
	}

	return args, nil
}

func b2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func b2u(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func setValue(option *Option, value string) error {

	ptr := option.value
	switch p := ptr.(type) {
	case *int:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = int(b2i(iValue))
		} else {
			if math.MaxInt == math.MaxInt32 {
				iValue, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					return err
				}
				*p = int(iValue)
			} else if math.MaxInt == math.MaxInt64 {
				iValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				*p = int(iValue)
			} else {
				return ErrUnexpectedIntegerWidth
			}
		}
	case *int8:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = b2i(iValue)
		} else {
			iValue, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return err
			}
			*p = int8(iValue)
		}
	case *int16:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = int16(b2i(iValue))
		} else {
			iValue, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}
			*p = int16(iValue)
		}
	case *int32:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = int32(b2i(iValue))
		} else {
			iValue, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			*p = int32(iValue)
		}
	case *int64:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = int64(b2i(iValue))
		} else {
			iValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			*p = iValue
		}
	case *uint:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = uint(b2u(iValue))
		} else {
			if math.MaxUint == math.MaxUint32 {
				iValue, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					return err
				}
				*p = uint(iValue)
			} else if math.MaxUint == math.MaxUint64 {
				iValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				*p = uint(iValue)
			} else {
				return ErrUnexpectedIntegerWidth
			}
		}
	case *uint8:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = b2u(iValue)
		} else {
			iValue, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return err
			}
			*p = uint8(iValue)
		}
	case *uint16:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = uint16(b2u(iValue))
		} else {
			iValue, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return err
			}
			*p = uint16(iValue)
		}
	case *uint32:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = uint32(b2u(iValue))
		} else {
			iValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return err
			}
			*p = uint32(iValue)
		}
	case *uint64:
		iValue, err := strconv.ParseBool(value)
		if err == nil {
			*p = uint64(b2u(iValue))
		} else {
			iValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			*p = iValue
		}
	case *bool:
		iValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*p = iValue
	case *string:
		*p = value
	default:
		return ErrUnsupportedType
	}

	option.exists = true
	if option.trigger != nil {
		option.trigger(option)
	}
	return nil
}

func incrementFlag(option *Option, isNegated bool) error {

	if isNegated {
		return zeroFlag(option)
	}

	ptr := option.value
	switch p := ptr.(type) {
	case *int:
		*p++
	case *int8:
		*p++
	case *int16:
		*p++
	case *int32:
		*p++
	case *int64:
		*p++
	case *uint:
		*p++
	case *uint8:
		*p++
	case *uint16:
		*p++
	case *uint32:
		*p++
	case *uint64:
		*p++
	case *bool:
		*p = true
	default:
		return ErrUnsupportedType
	}

	option.exists = true
	if option.trigger != nil {
		option.trigger(option)
	}
	return nil
}

func zeroFlag(option *Option) error {

	ptr := option.value
	switch p := ptr.(type) {
	case *int:
		*p = 0
	case *int8:
		*p = 0
	case *int16:
		*p = 0
	case *int32:
		*p = 0
	case *int64:
		*p = 0
	case *uint:
		*p = 0
	case *uint8:
		*p = 0
	case *uint16:
		*p = 0
	case *uint32:
		*p = 0
	case *uint64:
		*p = 0
	case *bool:
		*p = false
	default:
		return ErrUnsupportedType
	}

	option.exists = true
	if option.trigger != nil {
		option.trigger(option)
	}
	return nil
}
