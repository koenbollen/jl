# Elastic Common Schema

`jl` understands basic log fields of the Elastic Common Schema (ECS), and 
supports to include fields that have a nested path.

    $ common_schema | jl -f request.method,request.path,event.duration
    [2020-10-23 03:35:49]    INFO: Served [customer=test event.duration=78518000 request.method=GET request.path=/users/users/notices/]
