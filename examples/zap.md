
    $ some_zap_program
    {"level":"info","message":"logger construction succeeded","foo":"bar","ts":1565361391.4279764}

    $ some_zap_program | jl
    [2019-08-09 14:36:31]    INFO: logger construction succeeded [foo=bar]

    $ some_zap_program --error
    {"level":"error","msg":"panic!","error":"timeout","stack":"go.uber.org/zap.Stack\n\tgo.uber.org/zap/field.go:191\nmain.somefunction\n\tmain.go:11\nmain.main\n\tmain.go:15"} (no-eol)

    $ some_zap_program --error | jl
      ERROR: panic! [error=timeout]
        timeout
        go.uber.org/zap.Stack
          go.uber.org/zap/field.go:191
        main.somefunction
          main.go:11
        main.main
          main.go:15

    $ some_zap_program --complex-error | jl
    [2017-11-22 15:32:28]   ERROR: Kafka consumer received error. [caller=kafka/consumer.go:63 environment=development production=false tier=mailer version=5bb5b52]
        kafka server: The provided member is not known in the current generation.
        github.com/koenbollen/stream-processor-example/vendor/github.com/blendle/go-streamprocessor/streamclient/kafka.(*Client).NewConsumer.func2
          /home/jenkins/go/src/github.com/koenbollen/stream-processor-example/vendor/github.com/blendle/go-streamprocessor/streamclient/kafka/consumer.go:63
