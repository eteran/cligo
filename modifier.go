package cligo

type Modifier func(opt *Option)

// Needs specifies that the associated option requires that the option referred to by dep also be set.
func Needs(dep *Option) Modifier {
	return func(opt *Option) {
		opt.needs = append(opt.needs, dep)
	}
}

// Excludes specifies that the associated option requires that the option referred to by dep NOT be set.
func Excludes(dep *Option) Modifier {
	return func(opt *Option) {
		opt.excludes = append(opt.excludes, dep)
		dep.excludes = append(dep.excludes, opt)
	}
}

// Required specifies that the associated option MUST be set.
func Required() Modifier {
	return func(opt *Option) {
		opt.isRequired = true
	}
}

// IgnoreCase specifies that the associated option will be accepted in a case insensitive way.
// For example, --filename, and --FiLeNaMe would both be accepted.
func IgnoreCase() Modifier {
	return func(opt *Option) {
		opt.ignoreCase = true
	}
}

// Group places the associated option in a group specified by groupName and will be displayed
// grouped together will all other options in the same group.
func Group(groupName string) Modifier {
	return func(opt *Option) {
		opt.group = groupName
	}
}

// Trigger associates a callback function to trigger for each instance of a given option.
func Trigger(trigger Callback) Modifier {
	return func(opt *Option) {
		opt.onSet = trigger
	}
}

// DefaultString associates a default value in string form to print during usage statements.
// For example:
//
//	app := cligo.NewApp()
//	option1 := "hello world"
//	app.AddOption("-a,--alpha", &option1, "Option1", cligo.DefaultString("hello world"))
//	app.Usage()
//
// will print a usage statement reflecting that the default value for --alpha is "hello world".
func DefaultString(value string) Modifier {
	return func(opt *Option) {
		opt.defaultString = value
	}
}

// CaptureDefault associates the current value of the variable linked to this option as the
// default value string to use during usage statements.
// For example:
//
//	app := cligo.NewApp()
//	option1 := "hello world"
//	app.AddOption("-a,--alpha", &option1, "Option1", cligo.CaptureDefault())
//	app.Usage()
//
// will print a usage statement reflecting that the default value for --alpha is "hello world".
func CaptureDefault() Modifier {
	return func(opt *Option) {
		opt.defaultString = getValue(opt.ptr)
	}
}

// AddValidator adds a validator to a given option.
func AddValidator(v Validator) Modifier {
	return func(opt *Option) {
		opt.validators = append(opt.validators, v)
	}
}

// AddValidator is a convenience function which will add 1 or more validators to to a
// given option in the order that they are passed.
func AddValidators(v ...Validator) Modifier {
	return func(opt *Option) {
		opt.validators = append(opt.validators, v...)
	}
}
