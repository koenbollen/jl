# jl — JSON Logs

`jl` is a development tool for working with structured JSON logging.

### Example usage:

    $ myprogram
    {"message": "Hello, world!!", "severity": "info"}

    $ myprogram | jl
    INFO: Hello, world!!

### Installation

    $ go get -u github.com/koenbollen/jl
    $ echo '{"msg": "It works!"}' | jl
    It works!

### Project Status

The `jl` project is in early development

**Next steps:**

- [X] CI
- [X] CLI options/toggles/flags, also support environment defaults
- [ ] Support stacktraces
- [ ] Colorize output, like [bunyan](https://github.com/trentm/node-bunyan) does
- [ ] More examples/tests to make sure to be compatible with existing tooling
- [ ] Enhance this README and the --help
- [ ] Ship it! :shipit:
