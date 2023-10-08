package cligo

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Validator func(str string) string

// Check for existing file (returns error message if check fails)
func ExistingFile(path string) Validator {
	return func(str string) string {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return err.Error()
		}

		return ""
	}
}

// Check for an existing directory (returns error message if check fails)
func ExistingDirectory(path string) Validator {
	return func(str string) string {
		st, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return err.Error()
		}

		if !st.IsDir() {
			return fmt.Sprintf("%s is not a directory", path)
		}

		return ""
	}
}

// Produce a range (factory). Min and max are inclusive.
func Range(min, max int64) Validator {
	return func(str string) string {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err.Error()
		}

		if i < min || i > max {
			return fmt.Sprintf("%s is not in the range of [%d-%d]", str, min, max)
		}

		return ""
	}
}

/*
ExistingPath	Check for an existing path
NonexistentPath	Check for an non-existing path
*/
