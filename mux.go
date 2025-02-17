package cligo

import (
	"fmt"
	"os"
	"strings"
)

type CommandMux struct {
	commands map[string]*App
}

// CommandMux returns a new instance of the CommandMux type
func NewCommandMux() *CommandMux {
	return &CommandMux{
		commands: make(map[string]*App),
	}
}

func (mux *CommandMux) AddCommand(cmd string, app *App) error {

	if _, ok := mux.commands[cmd]; ok {
		return fmt.Errorf("command '%s' already registered", cmd)
	}

	mux.commands[cmd] = app
	return nil
}

func (mux *CommandMux) ParseStrict() error {
	return mux.ParseArgsStrict(os.Args[1:])
}

func (mux *CommandMux) ParseArgsStrict(args []string) error {

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

func (mux *CommandMux) Parse() ([]string, error) {
	return mux.ParseArgs(os.Args[1:])
}

func (mux *CommandMux) ParseArgs(args []string) ([]string, error) {
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

func (mux *CommandMux) subCommands() []string {

	commands := make([]string, 0)

	for name := range mux.commands {
		commands = append(commands, name)
	}

	return commands
}

func (mux *CommandMux) subCommandString() string {

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
