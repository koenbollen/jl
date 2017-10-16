# Bunyan compatible

[Bunyan](https://github.com/trentm/node-bunyan) is a simple and fast JSON
logging library for Node.JS services. `jl` aims to be compatible.

Here is an example bunyan output:

    $ fake_bunyan_app
    {"name":"myapp","hostname":"banana.local","pid":40161,"level":30,"msg":"hi","time":"2013-01-04T18:46:23.851Z","v":0}
    {"name":"myapp","hostname":"banana.local","pid":40161,"level":40,"lang":"fr","msg":"au revoir","time":"2013-01-04T18:46:23.853Z","v":0}

Passing the output of this command through `jl` will make it more readable.

    $ fake_bunyan_app | jl
    [2013-01-04 18:46:23] INFO: hi
    [2013-01-04 18:46:23] WARN: au revoir [lang=fr]

Yes, buynan supplies it's own CLI tool for development but `jl` should
be compatible with more kinds of structured logging.
