# Handling malformed/unstructured logs

`jl` will handle the unstructured logs gracefully, printing the raw log entry when it cannot be parsed as a json object:

    $ unstructured_app | jl
    22:39:49: Executing clean...
    
    Task :clean
    {"message": "malformed json"

At the same time, any valid json will be interpreted as a structured log entry, even when some json keys are missing (e.g. no message):

    $ myprogram --no-message | jl
    [2017-09-28 06:43:13]   TRACE:  [user=john]
    [2017-09-28 06:43:14]   ERROR:  [@version=1.0.0]
