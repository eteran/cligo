package cligo

import (
	"fmt"
	"os"
	"strings"
)

type Mux struct {
	commands map[string]Parser
}

// NewMux returns a new instance of the Mux type
func NewMux() *Mux {
	return &Mux{
		commands: make(map[string]Parser),
	}
}

func (mux *Mux) AddCommand(cmd string, parser Parser) error {

	if _, ok := mux.commands[cmd]; ok {
		return fmt.Errorf("command '%s' already registered", cmd)
	}

	mux.commands[cmd] = parser
	return nil
}

func (mux *Mux) CreateCommand(cmd string, f func(app *App)) error {

	app := NewApp()
	f(app)
	return mux.AddCommand(cmd, app)
}

func (mux *Mux) ParseStrict() error {
	return mux.ParseArgsStrict(os.Args[1:])
}

func (mux *Mux) ParseArgsStrict(args []string) error {

	if len(os.Args) < 2 {
		return fmt.Errorf("missing sub-command. expected to be one of: %s", mux.subCommandString())
	}

	cmd := os.Args[1]
	app, ok := mux.commands[cmd]
	if !ok {
		return fmt.Errorf("invalid sub-command '%s'. expected to be one of: %s", cmd, mux.subCommandString())
	}

	return app.ParseArgsStrict(os.Args[2:])
}

func (mux *Mux) Parse() ([]string, error) {
	return mux.ParseArgs(os.Args[1:])
}

func (mux *Mux) ParseArgs(args []string) ([]string, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("missing sub-command. expected to be one of: %s", mux.subCommandString())
	}

	cmd := os.Args[1]
	app, ok := mux.commands[cmd]
	if !ok {
		return nil, fmt.Errorf("invalid sub-command '%s'. expected to be one of: %s", cmd, mux.subCommandString())
	}

	return app.ParseArgs(os.Args[2:])
}

func (mux *Mux) subCommands() []string {

	commands := make([]string, 0)

	for name := range mux.commands {
		commands = append(commands, name)
	}

	return commands
}

func (mux *Mux) subCommandString() string {

	commands := mux.subCommands()

	switch len(commands) {
	case 0:
		return ""
	case 1:
		return commands[0]
	case 2:
		return fmt.Sprintf("%s or %s", commands[0], commands[1])
	default:
		return strings.Join(commands[:len(commands)-1], ", ") + " or " + commands[len(commands)-1]
	}
}
