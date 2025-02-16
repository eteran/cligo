# cligo

cligo is a library is designed to be clean, intuitive, but powerful. It library was inspired the very robust CLI11 library, and incorporates many of its user friendly features.


cligo implements classic GNU style parameters with both long and short forms of arguments. It supports validation of arguments, the ability to group them, and as well as the ability to group have arguments require that others do or do not exist. Because validators are just functions, they can be of course be expanded on by consumers of the library. The following is a classic example of an application which requires a filename, and has an optional verbosity flag.

```go
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/eteran/cligo"
)

func main() {
	app := cligo.NewApp()

	var filename string
	var verbose bool

	app.AddOption("-f,--file", &filename, "filename", cligo.Required(), cligo.AddValidator(cligo.ExistingFile()))
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	if err := app.ParseStrict(); err != nil {
		if !errors.Is(err, cligo.ErrHelpRequested) {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
```

