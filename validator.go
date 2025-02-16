package cligo

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Validator func(str string) error

// Check for existing file
func ExistingFile() Validator {
	return func(path string) error {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return err
		}

		return nil
	}
}

// Check for an existing directory
func ExistingDirectory() Validator {
	return func(path string) error {
		st, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return err
		}

		if !st.IsDir() {
			return fmt.Errorf("%s is not a directory", path)
		}

		return nil
	}
}

// Produce a range (factory). Min and max are inclusive.
func Range(min, max int64) Validator {
	return func(str string) error {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}

		if i < min || i > max {
			return fmt.Errorf("%s is not in the range of [%d-%d]", str, min, max)
		}

		return nil
	}
}

/*
ExistingPath	Check for an existing path
NonexistentPath	Check for an non-existing path
*/
