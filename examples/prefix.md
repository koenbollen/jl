# Prefixed Logs

[Stern](https://github.com/wercker/stern) is an example of logging that outputs
a string prefix before the sturctured JSON line. As you can see here:

    $ fake_stern mypod-549863668-1dqpx
    monitor-1331628357-qn2vp mailer-monitor-monitor {"labels": {"git_rev": "da161cc"}, "message": "context ended: NL/2017-09-28/digest", "origin": "EmailSentMonitor", "severity": "INFO", "timestamp": "2017-09-28T05:56:36Z"}

Passing this output through `jl` will include this prefix in the formatted
result. As shown here:

    $ fake_stern mypod-549863668-1dqpx | jl
    monitor-1331628357-qn2vp mailer-monitor-monitor [2017-09-28 05:56:36]    INFO: context ended: NL/2017-09-28/digest [git_rev=da161cc origin=EmailSentMonitor]

This behaviour can always be altered using `--skip-prefix`:

```
$ fake_stern mypod-549863668-1dqpx | jl --skip-prefix
[2017-09-28 05:56:36] INFO: context ended: NL/2017-09-28/digest [git_rev=da161cc origin=EmailSentMonitor]
```

(this feature is also available for suffixes)
