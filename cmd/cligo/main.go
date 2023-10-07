package main

import (
	. "cligo/pkg/cligo"
	"fmt"
	"os"
)

func main() {

	app := NewApp()

	var option1 string
	var option2 string

	o1 := app.AddOption("-a,--alpha", &option1, "Option1")
	o2 := app.AddOption("-b,--beta", &option2, "Option2", Excludes(o1))
	_ = o2

	err := app.ParseStrict()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println(option1, ":", option2)

}
