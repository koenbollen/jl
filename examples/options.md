# `jl` Options & Flags

## --help

Show the tools usage and help:

    $ jl --help
    jl - JSON Logs
    
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
      --include-fields <fields>, -f <fields>
                        Always include these json keys as fields (comma
                        separated list)
    
    You can add any option to the JL_OPTS environment variable, ex:
      export JL_OPTS="--no-color"
    
    Example:
      $ echo '{"level": "info", "msg": "Hello!", "size": 42}' | jl
      INFO: Hello! [size=42]

## Prefix & Suffix Data

`jl` will also output text/data around detected JSON on each line. The --skip-prefix and --skip-suffix flags allows you to disable this behaviour:

    $ echo 'prefix {"msg": "Hi", "level": "debug"} suffix' | jl
    prefix   DEBUG: Hi suffix

    $ echo 'prefix {"msg": "Hi", "level": "debug"} suffix' | jl --skip-prefix
      DEBUG: Hi suffix

    $ echo 'prefix {"msg": "Hi", "level": "debug"} suffix' | jl --skip-suffix
    prefix   DEBUG: Hi

    $ echo 'prefix {"msg": "Hi", "level": "debug"} suffix' | jl --skip-prefix --skip-suffix
      DEBUG: Hi

## Fields

Most JSON logging will include more fields then just the message. These fields are also printed through `jl` when the length of the value is not to long. You can influence this using the --skip-fields and --include-fields flags:

    $ echo '{"level": "warning", "msg": "Login failed", "user_id": "42"}' | jl
    WARNING: Login failed [user_id=42]

    $ echo '{"level": "warning", "msg": "Login failed", "user_id": "42"}' | jl --skip-fields
    WARNING: Login failed

    $ echo '{"msg": "test", "val": "Lorem ipsum dolor sit amet."}' | jl
    test

    $ echo '{"msg": "test", "val": "Lorem ipsum dolor sit amet."}' | jl --include-fields val
    test [val=Lorem ipsum dolor sit amet.]
