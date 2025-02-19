package main

import (
	"fmt"
	"os"

	"github.com/eteran/cligo"
)

func Run(filename string, verbose bool) error {
	fmt.Println("Running Program")
	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Verbose : %v\n", verbose)
	return nil
}

func main() {

	var filename string
	var verbose bool

	app := cligo.NewApp()
	app.AddOption("-f,--file", &filename, "filename", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	if err := app.ParseStrict(); err != nil {
		fmt.Println(err)
		os.Exit(0)

	}

	if err := Run(filename, verbose); err != nil {
		panic(err)
	}
}
