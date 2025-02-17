package cligo

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Validator func(str string) error

// ExistingFile checks if the string is a path to an existing file
func ExistingFile() Validator {
	return func(path string) error {
		st, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return err
		}

		if !st.Mode().IsRegular() {
			return fmt.Errorf("%s is not a file", path)
		}

		return nil
	}
}

// ExistingDirectory checks if the string is a path to an existing directory
func ExistingDirectory() Validator {
	return func(path string) error {
		st, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return err
		}

		if !st.Mode().IsDir() {
			return fmt.Errorf("%s is not a directory", path)
		}

		return nil
	}
}

// ExistingPath checks if the string is a path to an existing file or existing directory
func ExistingPath() Validator {
	return func(path string) error {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return err
		}

		return nil
	}
}

// NonexistentPath checks if the string is a path to an non-existing file or existing directory
func NonexistentPath() Validator {
	return func(path string) error {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("%s exists", path)
	}
}

// Range checks if integer value is withing the range [min-max].
func Range(min int64, max int64) Validator {
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
