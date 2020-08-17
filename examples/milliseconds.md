# Timestamp milliseconds detection

`jl` automatically guesses if your timestamps are in milliseconds or seconds
when a unix timestamp is used:

    $ echo '{"time":1597246404774, "app":"service-v1","severity":"INFO","message":"Initializing Servlet dispatcher"}' | jl
    [2020-08-12 15:33:24]    INFO: Initializing Servlet dispatcher [app=service-v1]

    $ echo '{"time":1597246404,    "app":"service-v1","severity":"INFO","message":"Initializing Servlet dispatcher"}' | jl
    [2020-08-12 15:33:24]    INFO: Initializing Servlet dispatcher [app=service-v1]
