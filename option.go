package cligo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Signed interface {
	int | int8 | int16 | int32 | int64
}

type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Option struct {

	// A positional name
	pName string

	// A list of the long names without the leading dashes
	lNames []string

	// A list of the short names without the leading dashes
	sNames []string

	// A list of the negated long names without the leading dashes
	lNamesNeg []string

	// A list of the negated short names without the leading dashes
	sNamesNeg []string

	owner         *App
	count         int
	ptr           any
	description   string
	defaultString string
	group         string
	isFlag        bool
	isRequired    bool
	ignoreCase    bool
	needs         []*Option
	excludes      []*Option
	validators    []Validator
	onSet         Callback
	setter        setterFunc

	// TODO(eteran):
	// envname
}

type setterFunc func(opt *Option, value string, isNegated bool) error

type Callback func(opt *Option)

func (opt Option) Exists() bool {
	return opt.count > 0
}

func (opt Option) Count() int {
	return opt.count
}

func (opt Option) IsPositional() bool {
	return opt.pName != ""
}

func (opt Option) IsPositionalOnly() bool {
	return opt.IsPositional() &&
		len(opt.lNames) == 0 &&
		len(opt.sNames) == 0 &&
		len(opt.lNamesNeg) == 0 &&
		len(opt.sNamesNeg) == 0
}

func (opt Option) canonicalName() string {

	if len(opt.lNames) != 0 {
		return opt.lNames[0]
	}

	if len(opt.lNamesNeg) != 0 {
		return opt.lNamesNeg[0]
	}

	if len(opt.sNames) != 0 {
		return opt.sNames[0]
	}

	if len(opt.sNamesNeg) != 0 {
		return opt.sNamesNeg[0]
	}

	return opt.pName
}

func getValue(ptr any) string {
	switch p := ptr.(type) {
	case *int:
		return strconv.FormatInt(int64(*p), 10)
	case *int8:
		return strconv.FormatInt(int64(*p), 10)
	case *int16:
		return strconv.FormatInt(int64(*p), 10)
	case *int32:
		return strconv.FormatInt(int64(*p), 10)
	case *int64:
		return strconv.FormatInt(int64(*p), 10)
	case *uint:
		return strconv.FormatUint(uint64(*p), 10)
	case *uint8:
		return strconv.FormatUint(uint64(*p), 10)
	case *uint16:
		return strconv.FormatUint(uint64(*p), 10)
	case *uint32:
		return strconv.FormatUint(uint64(*p), 10)
	case *uint64:
		return strconv.FormatUint(uint64(*p), 10)
	case *bool:
		return strconv.FormatBool(*p)
	case *string:
		return *p
	default:
		return ""
	}
}

func setValue(ptr any, value string) error {

	switch p := ptr.(type) {
	case *int:
		n, err := parseIntOrBool[int](value, 0)
		if err != nil {
			return err
		}
		*p = n
	case *int8:
		n, err := parseIntOrBool[int8](value, 8)
		if err != nil {
			return err
		}
		*p = n
	case *int16:
		n, err := parseIntOrBool[int16](value, 16)
		if err != nil {
			return err
		}
		*p = n
	case *int32:
		n, err := parseIntOrBool[int32](value, 32)
		if err != nil {
			return err
		}
		*p = n
	case *int64:
		n, err := parseIntOrBool[int64](value, 64)
		if err != nil {
			return err
		}
		*p = n
	case *uint:
		n, err := parseUintOrBool[uint](value, 0)
		if err != nil {
			return err
		}
		*p = n
	case *uint8:
		n, err := parseUintOrBool[uint8](value, 8)
		if err != nil {
			return err
		}
		*p = n
	case *uint16:
		n, err := parseUintOrBool[uint16](value, 16)
		if err != nil {
			return err
		}
		*p = n
	case *uint32:
		n, err := parseUintOrBool[uint32](value, 32)
		if err != nil {
			return err
		}
		*p = n
	case *uint64:
		n, err := parseUintOrBool[uint64](value, 64)
		if err != nil {
			return err
		}
		*p = n
	case *bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*p = b
	case *string:
		*p = value
	default:
		return ErrUnsupportedType
	}
	return nil
}

func parseUintOrBool[T Unsigned](s string, bitSize int) (T, error) {

	if b, err := strconv.ParseBool(s); err == nil {
		if b {
			return 1, nil
		}
		return 0, nil
	}

	i, err := strconv.ParseUint(s, 10, bitSize)
	return T(i), err
}

func parseIntOrBool[T Signed](s string, bitSize int) (T, error) {

	if b, err := strconv.ParseBool(s); err == nil {
		if b {
			return 1, nil
		}
		return 0, nil
	}

	i, err := strconv.ParseInt(s, 10, bitSize)
	return T(i), err
}

func incrementFlag(ptr any, isNegated bool) error {

	if isNegated {
		return decrementFlag(ptr)
	}

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
	return nil
}

func decrementFlag(ptr any) error {

	switch p := ptr.(type) {
	case *int:
		*p--
	case *int8:
		*p--
	case *int16:
		*p--
	case *int32:
		*p--
	case *int64:
		*p--
	case *uint:
		*p--
	case *uint8:
		*p--
	case *uint16:
		*p--
	case *uint32:
		*p--
	case *uint64:
		*p--
	case *bool:
		*p = false
	default:
		return ErrUnsupportedType
	}
	return nil
}

func (opt *Option) format() string {
	nameList := make([]string, 0, len(opt.sNames)+len(opt.lNames)+len(opt.sNamesNeg)+len(opt.lNamesNeg))
	for _, str := range opt.sNames {
		nameList = append(nameList, "-"+str)
	}

	for _, str := range opt.sNamesNeg {
		nameList = append(nameList, "-"+str)
	}

	for _, str := range opt.lNames {
		nameList = append(nameList, "--"+str)
	}

	for _, str := range opt.lNamesNeg {
		nameList = append(nameList, "--"+str)
	}

	names := strings.Join(nameList, ",")
	names = names + pointerType(opt.ptr)

	if opt.defaultString != "" {
		names = names + fmt.Sprintf(" [%s]", opt.defaultString)
	}

	if opt.isRequired {
		names = names + " REQUIRED"
	}
	return fmt.Sprintf("  %-30s %s", names, opt.description)
}

func (opt *Option) formatPositional() string {
	name := opt.pName
	name = name + pointerType(opt.ptr)

	if opt.defaultString != "" {
		name = name + fmt.Sprintf(" [%s]", opt.defaultString)
	}

	if opt.isRequired {
		name = name + " REQUIRED"
	}
	return fmt.Sprintf("  %-30s %s", name, opt.description)
}

func NewOption(name string, ptr any, help string, modifiers ...Modifier) *Option {
	if ptr != nil {
		rv := reflect.ValueOf(ptr)
		if rv.Kind() != reflect.Ptr {
			panic("bound variables must be pointers")
		}
	}

	opt := &Option{
		description: help,
		isFlag:      false,
		group:       "Options",
		ptr:         ptr,
		setter: func(opt *Option, v string, isNegated bool) error {

			for _, validator := range opt.validators {
				if err := validator(v); err != nil {
					return err
				}
			}

			if err := setOption(opt.ptr, v, isNegated); err != nil {
				return err
			}

			opt.count++
			if opt.onSet != nil {
				opt.onSet(opt)
			}
			return nil
		},
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

		isNegated := strings.HasPrefix(optionName, "!")
		if isNegated {
			panic("only flags can be negated")
		}

		isLong := strings.HasPrefix(optionName, "--")
		isShort := !isLong && strings.HasPrefix(optionName, "-")
		isPositional := !isLong && !isShort

		if isLong {
			lName := optionName[2:]
			if isNegated {
				opt.lNamesNeg = append(opt.lNamesNeg, lName)
			} else {
				opt.lNames = append(opt.lNames, lName)
			}
		} else if isShort {
			sName := optionName[1:]
			if isNegated {
				opt.sNamesNeg = append(opt.sNamesNeg, sName)
			} else {
				opt.sNames = append(opt.sNames, sName)
			}
		} else if isPositional {
			opt.pName = optionName
		}
	}

	return opt
}

func NewFlag(name string, ptr any, help string, modifiers ...Modifier) *Option {

	ensureIntegralPointer(ptr)

	opt := &Option{
		description: help,
		isFlag:      true,
		ptr:         ptr,
		group:       "Options",
		setter: func(opt *Option, v string, isNegated bool) error {

			for _, validator := range opt.validators {
				if err := validator(v); err != nil {
					return err
				}
			}

			if err := setFlag(opt.ptr, v, isNegated); err != nil {
				return err
			}

			opt.count++
			if opt.onSet != nil {
				opt.onSet(opt)
			}

			return nil
		},
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
			if isNegated {
				opt.lNamesNeg = append(opt.lNamesNeg, lName)
			} else {
				opt.lNames = append(opt.lNames, lName)
			}
		} else if isShort {
			sName := flagName[1:]
			if isNegated {
				opt.sNamesNeg = append(opt.sNamesNeg, sName)
			} else {
				opt.sNames = append(opt.sNames, sName)
			}
		} else if isPositional {
			panic("positional arguments must not be flags")
		}
	}

	return opt
}

func ensureIntegralPointer(ptr any) {
	if ptr != nil {
		rv := reflect.ValueOf(ptr)
		if rv.Kind() != reflect.Ptr {
			panic("bound variables must be pointers")
		}

		switch rv.Elem().Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			break
		default:
			panic("flags may only be boolean or integral types")
		}
	}
}
