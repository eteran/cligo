package cligo

type Modifier func(opt *Option)

func Needs(dep *Option) Modifier {
	return func(opt *Option) {
		opt.needs = append(opt.needs, dep)
	}
}

func Excludes(dep *Option) Modifier {
	return func(opt *Option) {
		opt.excludes = append(opt.excludes, dep)
		dep.excludes = append(dep.excludes, opt)
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

func Group(group string) Modifier {
	return func(opt *Option) {
		opt.group = group
	}
}

func Trigger(trigger Callback) Modifier {
	return func(opt *Option) {
		opt.onSet = trigger
	}
}

func DefaultString(value string) Modifier {
	return func(opt *Option) {
		opt.defaultString = value
	}
}

func CaptureDefault() Modifier {
	return func(opt *Option) {
		opt.defaultString = getValue(opt.ptr)
	}
}
