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

	installApp := cligo.NewApp()
	installApp.AddOption("-f,--file", &filename, "filename", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
	installApp.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	removeApp := cligo.NewApp()
	removeApp.AddOption("-n,--name", &name, "name", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
	removeApp.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	listApp := cligo.NewApp()
	listApp.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	mux := cligo.NewCommandMux()
	mux.AddCommand("install", installApp)
	mux.AddCommand("remove", removeApp)
	mux.AddCommand("list", listApp)

	if err := mux.ParseStrict(); err != nil {
		if !errors.Is(err, cligo.ErrHelpRequested) {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
