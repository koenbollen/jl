# Nested Fields

By default `jl` will ignore nested fields unless you explicitly include them:

    $ echo '{"msg":"Booting...", "level": "INFO", "meta": {"app": "backend", "server": 6}}' | jl
       INFO: Booting...

    $ echo '{"msg":"Booting...", "level": "INFO", "meta": {"app": "backend", "server": 6}}' | jl -f meta.app
       INFO: Booting... [meta.app=backend]

    $ echo '{"msg":"Booting...", "level": "INFO", "meta": {"app": "backend", "server": 6}}' | jl -f meta
       INFO: Booting... [meta.app=backend meta.server=6]


Some logging formats have their message in a nested fields, like
[tokio/tracing](https://docs.rs/tracing-subscriber/latest/tracing_subscriber/fmt/format/struct.Json.html#example-output)
for example, have the log message in the `.fields.message` field:

    $ myprogram --nested-message
    {"timestamp":"2022-02-15T18:47:10.821422Z","level":"INFO","fields":{"message":"shaving yaks","yaks":7},"target":"fmt_json","spans":{"yaks":7,"name":"shaving_yaks"}}
    {"timestamp":"2022-02-15T18:47:10.821495Z","level":"TRACE","fields":{"message":"hello! Im gonna shave a yak","excitement":"yay!"},"target":"fmt_json","spans":{"yaks":7,"name":"shaving_yaks"}}

`jl` will also look for the json path `.*.message`, and when it finds `.fields.message` it will output that message and automatically include the `fields` object:

    $ myprogram --nested-message | jl
    [2022-02-15 18:47:10]    INFO: shaving yaks [fields.yaks=7 target=fmt_json]
    [2022-02-15 18:47:10]   TRACE: hello! Im gonna shave a yak [fields.excitement=yay! target=fmt_json]

You can also explicitly include other nested objects as well:

    $ myprogram --nested-message | jl -f spans
    [2022-02-15 18:47:10]    INFO: shaving yaks [fields.yaks=7 spans.name=shaving_yaks spans.yaks=7 target=fmt_json]
    [2022-02-15 18:47:10]   TRACE: hello! Im gonna shave a yak [fields.excitement=yay! spans.name=shaving_yaks spans.yaks=7 target=fmt_json]
