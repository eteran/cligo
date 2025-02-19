package cligo

import (
	"fmt"
	"os"
)

type Mux struct {
	commands  map[string]*App
	usageFunc UsageFunc
}

// NewMux returns a new instance of the Mux type
func NewMux() *Mux {
	return &Mux{
		commands: make(map[string]*App),
	}
}

// SetUsageFunc sets the function to call in order to print the usage string
// setting a nil usage function (the default) will result in the default usage function
// being used
func (mux *Mux) SetUsageFunc(f UsageFunc) {
	mux.usageFunc = f
}

func (mux *Mux) AddCommand(cmd string, app *App) error {

	if _, ok := mux.commands[cmd]; ok {
		return fmt.Errorf("command '%s' already registered", cmd)
	}

	mux.commands[cmd] = app
	return nil
}

func (mux *Mux) CreateCommand(cmd string, description string, f func(app *App)) error {

	app := NewApp()
	app.SetName(cmd)
	app.SetDescription(description)
	f(app)
	return mux.AddCommand(cmd, app)
}

func (mux Mux) Usage() {
	if mux.usageFunc != nil {
		mux.usageFunc()
		return
	}

	fmt.Printf("Usage: %s [COMMAND] [OPTIONS]\n\n", os.Args[0])
	fmt.Println("Commands:")
	for name, cmd := range mux.commands {
		fmt.Printf("  %-15s - %s\n", name, cmd.description)
	}
}

func (mux *Mux) ParseStrict() (string, error) {
	return mux.ParseArgsStrict(os.Args[1:])
}

func (mux *Mux) ParseArgsStrict(args []string) (string, error) {

	if len(os.Args) < 2 {
		mux.Usage()
		os.Exit(-1)
	}

	cmd := os.Args[1]
	app, ok := mux.commands[cmd]
	if !ok {
		mux.Usage()
		os.Exit(-1)
	}

	if err := app.ParseArgsStrict(os.Args[2:]); err != nil {
		return "", err
	}

	return cmd, nil
}

func (mux *Mux) Parse() (string, []string, error) {
	return mux.ParseArgs(os.Args[1:])
}

func (mux *Mux) ParseArgs(args []string) (string, []string, error) {
	if len(os.Args) < 2 {
		mux.Usage()
		os.Exit(-1)
	}

	cmd := os.Args[1]
	app, ok := mux.commands[cmd]
	if !ok {
		mux.Usage()
		os.Exit(-1)
	}

	rest, err := app.ParseArgs(os.Args[2:])
	if err != nil {
		return "", nil, err
	}

	return cmd, rest, nil
}
