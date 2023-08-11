# Go's log/slog json handler

    $ slog_example
    {"time":"2006-01-02T15:04:05Z","level":"INFO","msg":"hello","count":3}
    {"time":"2006-01-02T15:04:05Z","level":"WARN","msg":"failed","err":"EOF"}

    $ slog_example | jl
    [2006-01-02 15:04:05]    INFO: hello [count=3]
    [2006-01-02 15:04:05] WARNING: failed [err=EOF]
