# Formatted Stacktraces

Some structured logging libraries also support the output of stacktraces, `jl`
tried to detect this and format it neatly.

    $ fake_bunyan_app --error
    {.*"err":{"stack":".*","name":"TypeError","message":"boom"},"level":40,"msg":"operation went boom: TypeError: boom","time":"2012-02-02T04:42:53.206Z","v":0} (re)

Formatting buynan errors like it would:

    $ fake_bunyan_app --error | jl
    [2012-02-02 04:42:53] WARNING: operation went boom: TypeError: boom
        TypeError: boom
            at Object.<anonymous> (/Users/trentm/tm/node-bunyan/examples/err.js:15:9)
            at Module._compile (module.js:411:26)
            at Object..js (module.js:417:10)
            at Module.load (module.js:343:31)
            at Function._load (module.js:302:12)
            at Array.0 (module.js:430:10)
            at EventEmitter._tickCallback (node.js:126:26)

Zap's stacktrace feature has no set key but will always start with the zap pacakge:

    $ some_zap_program --error | jl
      ERROR: panic! [error=timeout]
        timeout
        go.uber.org/zap.Stack
          go.uber.org/zap/field.go:191
        main.somefunction
          main.go:11
        main.main
          main.go:15
