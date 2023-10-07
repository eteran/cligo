package cligo

import (
	"reflect"
	"strconv"
)

type Signed interface {
	int | int8 | int16 | int32 | int64
}

type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Option struct {
	positionalName   string
	longNames        []string
	shortNames       []string
	help             string
	value            any
	defaultValue     string
	exists           bool
	isFlag           bool
	isPositionalOnly bool
	isRequired       bool
	ignoreCase       bool
	group            string
	trigger          Trigger
	needs            []*Option
	excludes         []*Option

	// TODO(eteran):
	// envname
}

type Trigger func(opt *Option)

type Modifier func(opt *Option)

func Needs(dep *Option) Modifier {
	return func(opt *Option) {
		opt.needs = append(opt.needs, dep)
	}
}

func Excludes(dep *Option) Modifier {
	return func(opt *Option) {
		opt.excludes = append(opt.excludes, dep)
	}
}

func Required() Modifier {
	return func(opt *Option) {
		opt.isRequired = true
	}
}

func IgnoreCase() Modifier {
	return func(opt *Option) {
		opt.ignoreCase = true
	}
}

func SetGroup(group string) Modifier {
	return func(opt *Option) {
		opt.group = group
	}
}

func SetTrigger(trigger Trigger) Modifier {
	return func(opt *Option) {
		opt.trigger = trigger
	}
}

func SetDefault(value string) Modifier {
	return func(opt *Option) {
		if err := opt.setValue(value); err != nil {
			panic("failed to set default")
		}
		opt.defaultValue = value
	}
}

func IsNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func (opt Option) Exists() any {
	return opt.exists
}

func (opt Option) canonicalName() string {

	if len(opt.longNames) != 0 {
		return opt.longNames[0]
	}

	if len(opt.shortNames) != 0 {
		return opt.shortNames[0]
	}

	return opt.positionalName
}

func (opt *Option) setValue(value string) error {

	ptr := opt.value
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

	opt.exists = true
	if opt.trigger != nil {
		opt.trigger(opt)
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

func (opt *Option) incrementFlag(isNegated bool) error {

	if isNegated {
		return opt.zeroFlag()
	}

	ptr := opt.value
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

	opt.exists = true
	if opt.trigger != nil {
		opt.trigger(opt)
	}
	return nil
}

func (opt *Option) zeroFlag() error {

	ptr := opt.value
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

	opt.exists = true
	if opt.trigger != nil {
		opt.trigger(opt)
	}
	return nil
}
