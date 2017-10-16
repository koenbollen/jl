package main

import (
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/mattn/go-isatty"
)

var usage = `jl - JSON logs

Usage:
  jl [options] [FILE...]

Options:
  -h, --help    Show this screen.
  --version     Show version.

Output Options:
  --color           Force colorized output
  --no-color        Don't colorize output
  --skip-prefix     Skip printing truncated bytes before the JSON
  --skip-suffix     Skip printing truncated bytes after the JSON
  --skip-fields     Don't output misc json keys as fields

You can add any option to the JL_OPTS environment variable, ex:
  export JL_OPTS="--no-color"

`

var version = "<unknown_version>"

func cli() (files []string, color, showPrefix, showSuffix, showFields bool) {
	argv := append(os.Args[1:], strings.Split(os.Getenv("JL_OPTS"), " ")...)
	arguments, err := docopt.Parse(usage, argv, true, "jl "+version, false)
	if err != nil {
		panic(err)
	}
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	color = !arguments["--no-color"].(bool) && (arguments["--color"].(bool) || isTTY)
	showPrefix = !arguments["--skip-prefix"].(bool)
	showSuffix = !arguments["--skip-suffix"].(bool)
	showFields = !arguments["--skip-fields"].(bool)
	files = arguments["FILE"].([]string)
	return
}
