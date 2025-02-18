package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/eteran/cligo"
)

func main() {

	var filename string
	var verbose bool
	var name string

	mux := cligo.NewMux()
	mux.CreateCommand("install", func(app *cligo.App) {
		app.AddOption("-f,--file", &filename, "filename", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
		app.AddFlag("-v,--verbose", &verbose, "increase verbosity")
	})

	mux.CreateCommand("remove", func(app *cligo.App) {
		app.AddOption("name", &name, "name", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
		app.AddFlag("-v,--verbose", &verbose, "increase verbosity")
	})

	mux.CreateCommand("list", func(app *cligo.App) {
		app.AddFlag("-v,--verbose", &verbose, "increase verbosity")
	})

	err := mux.ParseStrict(func(cmd string) error {
		fmt.Println("Executing Command:", cmd)
		return nil

	})

	if err != nil {
		if !errors.Is(err, cligo.ErrHelpRequested) {
			fmt.Println(err)
			os.Exit(0)
		}
	}

}
