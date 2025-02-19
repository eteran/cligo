package main

import (
	"fmt"
	"os"

	"github.com/eteran/cligo"
)

func main() {

	var filename string
	var verbose bool
	var name string

	mux := cligo.NewMux()
	mux.CreateCommand("install", "install a package",
		func(app *cligo.App) {
			app.AddOption("-f,--file", &filename, "filename", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
			app.AddFlag("-v,--verbose", &verbose, "increase verbosity")
		})

	mux.CreateCommand("remove", "remove a package",
		func(app *cligo.App) {
			app.AddOption("name", &name, "name", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
			app.AddFlag("-v,--verbose", &verbose, "increase verbosity")
		})

	mux.CreateCommand("list", "list installed packages",
		func(app *cligo.App) {
			app.AddFlag("-v,--verbose", &verbose, "increase verbosity")
		})

	cmd, err := mux.ParseStrict()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println("Executing Command:", cmd)
}
