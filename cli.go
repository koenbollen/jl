package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/mattn/go-isatty"
)

var usage = `jl - JSON Logs

'jl' is a development tool for working with structured JSON logging

It will parse loglines from stdin and try to parse them as
structured logging entries. When such a message is detected it
will output the entry in a human readable way. Anything else
is forwarded as is.

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

Formatting Options:
  --skip-fields     Don't output misc json keys as fields
  --max-field-length <int>
                    Any field, exceeding the given length (including
                    field name) will be ommitted from output. Use 0
                    to remove the length limit [default: 30]
  --include-fields <fields>, -f <fields>
                    Always include these json keys as fields, no matter
                    the length (comma separated list)
  --exclude-fields <fields>
                    Always exclude these json keys (comma separated
                    list)

You can add any option to the JL_OPTS environment variable, ex:
  export JL_OPTS="--no-color"

Example:
  $ echo '{"level": "info", "msg": "Hello!", "size": 42}' | jl
  INFO: Hello! [size=42]
`

var version = "v1.4.0"

func cli() (files []string, color, showPrefix, showSuffix, showFields bool, includeFields string, excludeFields string, maxFieldLength int) {
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
	maxFieldLength, _ = strconv.Atoi(arguments["--max-field-length"].(string))
	includeFields, _ = arguments["--include-fields"].(string)
	excludeFields, _ = arguments["--exclude-fields"].(string)
	files = arguments["FILE"].([]string)
	return
}
