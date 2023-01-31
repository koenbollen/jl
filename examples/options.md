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

Most JSON logging will include more fields then just the message. These fields are also printed through `jl` when the length of the value does not exceed a defined limit. Several flags are available to control the fields treatment, see the usage examples below.

By default `jl` will interpret misc json keys as fields and print them out:

    $ echo '{"level": "warning", "msg": "Login failed", "user_id": "42"}' | jl
    WARNING: Login failed [user_id=42]

It is possible to disable the fields processing all together with the --skip-fields flag:

    $ echo '{"level": "warning", "msg": "Login failed", "user_id": "42"}' | jl --skip-fields
    WARNING: Login failed

If the length of the field (key + value) exceeds the default limit of 30 characters, it will not be printed:

    $ echo '{"msg": "test", "ver": "1.0.0", "val": "Lorem ipsum dolor sit amet."}' | jl
    test [ver=1.0.0]

However it's possible to override the length limit with a --max-field-length flag:

    $ echo '{"msg": "test", "ver": "1.0.0", "val": "Lorem ipsum dolor sit amet."}' | jl --max-field-length 40
    test [val=Lorem ipsum dolor sit amet. ver=1.0.0]

You can also specify the fields, which should always be printed, no matter the length using the --include-fields flag:

    $ echo '{"msg": "test", "val": "Lorem ipsum dolor sit amet."}' | jl --include-fields val
    test [val=Lorem ipsum dolor sit amet.]

For the opposite use-case, certain fields may be excluded from the output using the --exclude-fields flag:

    $ echo '{"msg": "test", "ver": "1.0.0", "val": "42"}' | jl --exclude-fields ver
    test [val=42]

Note, --include-fields takes precedence over --exclude-fields