# journald

Journald is a system service for collecting and storing log data, introduced with systemd. It tries to make it easier for system administrators to find interesting and relevant information among an ever-increasing amount of log messages. It also provides structured logging output, which `jl` can parse.

This is an short, slightly redacted, example of journald output:

    $ journald -xe -ojson
    {"_HOSTNAME":"example.org","_SYSTEMD_CGROUP":"/system.slice/sshd.service","_EXE":"/usr/sbin/sshd","__MONOTONIC_TIMESTAMP":"4231192657117","_CMDLINE":"sshd: unknown [priv]","_SYSTEMD_UNIT":"sshd.service","_MACHINE_ID":"be3292bb238d21a8de53f89d25ec97c4","_TRANSPORT":"stdout","PRIORITY":"5","__REALTIME_TIMESTAMP":"1686919896987169","_GID":"0","_CAP_EFFECTIVE":"1ffffffffff","__CURSOR":"s=11054c7dc82b4645a45da01c6bf62842","MESSAGE":"Invalid user hacker from 127.106.119.170 port 54520","SYSLOG_IDENTIFIER":"sshd","_UID":"0","_COMM":"sshd","SYSLOG_FACILITY":"3","_SYSTEMD_SLICE":"system.slice","_STREAM_ID":"08acce59fe1b44648b1d054f9a35156f","_PID":"1977203","_SYSTEMD_INVOCATION_ID":"9b199c04cfbe43afb339f73299c02a20","_BOOT_ID":"4cef257cf46b4818a75a0f463024e90d"}
    {"_UID":"0","__REALTIME_TIMESTAMP":"1686919897133605","_EXE":"/usr/sbin/sshd","_SYSTEMD_SLICE":"system.slice","_HOSTNAME":"example.org","_PID":"1977203","_STREAM_ID":"08acce59fe1b44648b1d054f9a35156f","_BOOT_ID":"4cef257cf46b4818a75a0f463024e90d","_CMDLINE":"sshd: unknown [priv]","_COMM":"sshd","PRIORITY":"6","_MACHINE_ID":"be3292bb238d21a8de53f89d25ec97c4","MESSAGE":"Disconnected from invalid user hacker 127.106.119.170 port 54520 [preauth]","_CAP_EFFECTIVE":"1ffffffffff","_SYSTEMD_CGROUP":"/system.slice/sshd.service","_SYSTEMD_UNIT":"sshd.service","__MONOTONIC_TIMESTAMP":"4231192803553","_TRANSPORT":"stdout","__CURSOR":"s=11054c7dc82b4645a45da01c6bf62842","_SYSTEMD_INVOCATION_ID":"9b199c04cfbe43afb339f73299c02a20","SYSLOG_IDENTIFIER":"sshd","_GID":"0","SYSLOG_FACILITY":"3"}

Passing this through `jl` will make it more readable:

    $ journald -xe -ojson | jl
    [2023-06-16 12:51:36]  NOTICE: Invalid user hacker from 127.106.119.170 port 54520 [SYSLOG_FACILITY=3 SYSLOG_IDENTIFIER=sshd _CAP_EFFECTIVE=1ffffffffff _CMDLINE=sshd: unknown [priv] _COMM=sshd _EXE=/usr/sbin/sshd _GID=0 _HOSTNAME=example.org _PID=1977203 _SYSTEMD_SLICE=system.slice _SYSTEMD_UNIT=sshd.service _TRANSPORT=stdout _UID=0]
    [2023-06-16 12:51:37]    INFO: Disconnected from invalid user hacker 127.106.119.170 port 54520 [preauth] [SYSLOG_FACILITY=3 SYSLOG_IDENTIFIER=sshd _CAP_EFFECTIVE=1ffffffffff _CMDLINE=sshd: unknown [priv] _COMM=sshd _EXE=/usr/sbin/sshd _GID=0 _HOSTNAME=example.org _PID=1977203 _SYSTEMD_SLICE=system.slice _SYSTEMD_UNIT=sshd.service _TRANSPORT=stdout _UID=0]

### Skipping fields:

`jl` will not output fields it used to parse the message, like __REALTIME_TIMESTAMP or PRIOIRTY. You can force `jl` to output these fields by explicitly including them:

    $ journald -xe -ojson | jl --include-field PRIORITY
    [2023-06-16 12:51:36]  NOTICE: Invalid user hacker from 127.106.119.170 port 54520 [PRIORITY=5 SYSLOG_FACILITY=3 SYSLOG_IDENTIFIER=sshd _CAP_EFFECTIVE=1ffffffffff _CMDLINE=sshd: unknown [priv] _COMM=sshd _EXE=/usr/sbin/sshd _GID=0 _HOSTNAME=example.org _PID=1977203 _SYSTEMD_SLICE=system.slice _SYSTEMD_UNIT=sshd.service _TRANSPORT=stdout _UID=0]
    [2023-06-16 12:51:37]    INFO: Disconnected from invalid user hacker 127.106.119.170 port 54520 [preauth] [PRIORITY=6 SYSLOG_FACILITY=3 SYSLOG_IDENTIFIER=sshd _CAP_EFFECTIVE=1ffffffffff _CMDLINE=sshd: unknown [priv] _COMM=sshd _EXE=/usr/sbin/sshd _GID=0 _HOSTNAME=example.org _PID=1977203 _SYSTEMD_SLICE=system.slice _SYSTEMD_UNIT=sshd.service _TRANSPORT=stdout _UID=0]
