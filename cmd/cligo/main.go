package main

import (
	. "cligo/pkg/cligo"
	"fmt"
	"os"
)

func main() {

	var filename string
	var verbose bool
	var mypos string
	var source string
	var flag uint16

	app := NewApp()

	app.AddOption("-f,--file", &filename, "A help string", SetGroup("Test"), SetDefault("monkey.txt"))
	app.AddFlag("-v,--verbose", &verbose, "Some other help string")
	app.AddFlag("--flag,!--no-flag", &flag, "help for flag")
	app.AddOption("-a,-b,--alpha,--beta,mypos", &mypos, "positional with an alias")
	app.AddOption("source", &source, "where to read the input from")

	err := app.ParseStrict()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println("Verbose:", verbose)
	fmt.Println("Filename is:", filename)
	fmt.Println("mypos is:", mypos)
	fmt.Println("source is:", source)
	fmt.Println("Flag is:", flag)

}
