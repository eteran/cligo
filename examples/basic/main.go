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

	app := cligo.NewApp()
	app.AddOption("-f,--file", &filename, "filename", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	if err := app.ParseStrict(); err != nil {
		if !errors.Is(err, cligo.ErrHelpRequested) {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
