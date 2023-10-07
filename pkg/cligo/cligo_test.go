package cligo_test

import (
	"cligo/pkg/cligo"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -a (flag)
func TestFlagShort(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var verbose bool
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	args := []string{"-v"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, verbose, true)
	}
}

// --long (long flag)
func TestFlagLong(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var verbose bool
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	args := []string{"--verbose"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, verbose, true)
	}
}

// --long_flag=true (long flag with equals to override default value)
func TestFlagLongEqual(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var verbose bool
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	args := []string{"--verbose=true"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, verbose, true)
	}
}

// --long_flag=true (long flag with equals to override default value)
func TestFlagLongEqualInteger(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var count int
	app.AddFlag("-c,--count", &count, "increase verbosity")

	args := []string{"--count=9000"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, count, 9000)
	}
}

func TestFlagShortRepeated(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var verbose int
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	args := []string{"-vvvv"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, verbose, 4)
	}
}

func TestFlagLongRepeated(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var verbose int
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	args := []string{"--verbose", "--verbose", "--verbose"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, verbose, 3)
	}
}

func TestFlagsRepeatedMixed(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var verbose int
	app.AddFlag("-v,--verbose", &verbose, "increase verbosity")

	args := []string{"--verbose", "-v", "--verbose", "-v", "--verbose"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, verbose, 5)
	}
}

// -abc (flags can be combined)
func TestFlagCombined(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var v1 bool
	var v2 bool
	var v3 bool
	app.AddFlag("-a,--alpha", &v1, "v1")
	app.AddFlag("-b,--beta", &v2, "v2")
	app.AddFlag("-g,--gamma", &v3, "v3")

	args := []string{"-ag"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, v1, true)
		assert.Equal(t, v2, false)
		assert.Equal(t, v3, true)
	}
}

// -abcf filename (flags and option can be combined)
func TestFlagOptionCombined(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var v1 bool
	var v2 bool
	var v3 bool
	var filename string
	app.AddFlag("-a,--alpha", &v1, "v1")
	app.AddFlag("-b,--beta", &v2, "v2")
	app.AddFlag("-g,--gamma", &v3, "v3")
	app.AddOption("-f,--file", &filename, "filename")

	args := []string{"-agf", "document.txt"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, v1, true)
		assert.Equal(t, v2, false)
		assert.Equal(t, v3, true)
		assert.Equal(t, filename, "document.txt")
	}
}

// -ffilename (no space required)
func TestOptionShortStringNoSpace(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var filename string
	app.AddOption("-f,--file", &filename, "filename")

	args := []string{"-ftest.txt"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, filename, "test.txt")
	}
}

// -f filename (option)
func TestOptionShortStringSpace(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var filename string
	app.AddOption("-f,--file", &filename, "filename")

	args := []string{"-f", "test.txt"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, filename, "test.txt")
	}
}

// --file filename (space)
func TestOptionLongStringSpace(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var filename string
	app.AddOption("-f,--file", &filename, "filename")

	args := []string{"--file", "test.txt"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, filename, "test.txt")
	}
}

// --file=filename (equals)
func TestOptionLongStringEqual(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var filename string
	app.AddOption("-f,--file", &filename, "filename")

	args := []string{"--file=test.txt"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, filename, "test.txt")
	}
}

func TestOptionLongIntegerEqual(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var value int
	app.AddOption("-v,--value", &value, "filename")

	args := []string{"--value=42"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, value, 42)
	}
}

func TestOptionPositionalString(t *testing.T) {
	t.Parallel()
	app := cligo.NewApp()

	var value int
	var dest string
	app.AddOption("-v,--value", &value, "filename")
	app.AddOption("dest", &dest, "dest")

	args := []string{"--value=42", "my_destination_file"}
	err := app.ParseArgsStrict(args)
	if assert.NoError(t, err) {
		assert.Equal(t, value, 42)
		assert.Equal(t, dest, "my_destination_file")
	}
}
