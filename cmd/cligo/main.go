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
	var option3 string
	var option4 int

	o1 := app.AddOption("-a,--alpha", &option1, "Option1")
	app.AddOption("-b,--beta", &option2, "Option2", Excludes(o1))
	app.AddOption("pos1", &option3, "option3")
	app.AddOption("pos2", &option4, "option4")

	err := app.ParseStrict()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println(option1, ":", option2, ":", option3, ":", option4)

}
