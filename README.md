# jl — JSON Logs

`jl` is a development tool for working with structured JSON logging.

### Example usage:

    $ myprogram
    {"message": "Hello, world!!", "severity": "info"}
    {"message": "skipping file", "severity": "warn", "file": "empty.txt"}

    $ myprogram | jl
       INFO: Hello, world!!
    WARNING: skipping file [file=empty.txt]

### Installation

    $ go get -u github.com/koenbollen/jl
    $ echo '{"msg": "It works!"}' | jl
    It works!

### Project Status

The `jl` project is in early development

**Next steps:**

- [X] CI
- [X] CLI options/toggles/flags, also support environment defaults
- [X] Support stacktraces
- [X] Colorize output, like [bunyan](https://github.com/trentm/node-bunyan) does
- [ ] More examples/tests to make sure to be compatible with existing tooling
- [ ] Enhance this README and the --help
- [ ] Ship it! :shipit:
