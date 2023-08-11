# Go's log/slog json handler

    $ slog_example
    {"time":"2023-08-11T16:35:31.162712808+02:00","level":"INFO","msg":"hello","count":3}
    {"time":"2023-08-11T16:35:31.162796194+02:00","level":"WARN","msg":"failed","err":"EOF"}

    $ slog_example | jl
    [2023-08-11 16:34:29]    INFO: hello [count=3]
    [2023-08-11 16:34:29] WARNING: failed [err=EOF]
