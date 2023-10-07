package cligo

import "reflect"

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

	// TODO(eteran):
	// needs
	// excludes
	// envname
}

type Trigger func(opt *Option)

type Modifier func(opt *Option)

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
		if err := setValue(opt, value); err != nil {
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

func (opt Option) Value() any {
	return opt.value
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
