package processors

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/koenbollen/jl/djson"
	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
)

func TestJournald_Happypath(t *testing.T) {
	t.Parallel()

	in := `{
		"_HOSTNAME": "example.org",
		"_SYSTEMD_CGROUP": "/system.slice/sshd.service",
		"_EXE": "/usr/sbin/sshd",
		"__MONOTONIC_TIMESTAMP": "4231192657117",
		"_CMDLINE": "sshd: unknown [priv]",
		"_SYSTEMD_UNIT": "sshd.service",
		"_MACHINE_ID": "be3292bb238d21a8de53f89d25ec97c4",
		"_TRANSPORT": "stdout",
		"PRIORITY": "6",
		"__REALTIME_TIMESTAMP": "1686919896987169",
		"_GID": "0",
		"_CAP_EFFECTIVE": "1ffffffffff",
		"__CURSOR": "s=11054c7dc82b4645a45da01c6bf62842",
		"MESSAGE": "Invalid user hacker from 127.106.119.170 port 54520",
		"SYSLOG_IDENTIFIER": "sshd",
		"_UID": "0",
		"_COMM": "sshd",
		"SYSLOG_FACILITY": "3",
		"_SYSTEMD_SLICE": "system.slice",
		"_STREAM_ID": "08acce59fe1b44648b1d054f9a35156f",
		"_PID": "1977203",
		"_SYSTEMD_INVOCATION_ID": "9b199c04cfbe43afb339f73299c02a20",
		"_BOOT_ID": "4cef257cf46b4818a75a0f463024e90d"
	}`
	raw := &bytes.Buffer{}
	_ = json.Compact(raw, []byte(in))

	line := &stream.Line{Raw: raw.Bytes(), JSON: raw.Bytes()}
	entry := &structure.Entry{SkipFields: make(map[string]bool)}
	djson.Unmarshal(line.JSON, entry)

	p := &JournaldProcessor{}

	detected := p.Detect(line, entry)
	if got, want := detected, true; got != want {
		t.Errorf("Detect() = %v, want %v", got, want)
	}

	err := p.Process(line, entry)
	if err != nil {
		t.Errorf("Process() = %v, want nil", err)
	}

	if got, want := entry.Message, "Invalid user hacker from 127.106.119.170 port 54520"; got != want {
		t.Errorf("entry.Message = %v, want %v", got, want)
	}

	if got, want := entry.Timestamp.String(), "2023-06-16 12:51:36.987169 +0000 UTC"; got != want {
		t.Errorf("entry.Timestamp = %v, want %v", got, want)
	}

	if got, want := entry.Severity, "INFO"; got != want {
		t.Errorf("entry.Severity = %v, want %v", got, want)
	}

	if got, want := entry.SkipFields["MESSAGE"], true; got != want {
		t.Errorf("entry.SkipFields['MESSAGE'] = %v, want %v", got, want)
	}

	if got, want := entry.SkipFields["__REALTIME_TIMESTAMP"], true; got != want {
		t.Errorf("entry.SkipFields['__REALTIME_TIMESTAMP'] = %v, want %v", got, want)
	}

	if got, want := entry.SkipFields["PRIORITY"], true; got != want {
		t.Errorf("entry.SkipFields['PRIORITY'] = %v, want %v", got, want)
	}
}

func TestJournald_InvalidTimestamp(t *testing.T) {
	t.Parallel()

	entry := process(t, `{
		"__REALTIME_TIMESTAMP": "yesterday",
		"PRIORITY": "6",
		"MESSAGE": "Invalid timestamp in __REALTIME_TIMESTAMP",
		"SYSLOG_IDENTIFIER": "sshd"
	}`)

	if got, want := entry.RawTimestamp, "yesterday"; got != want {
		t.Errorf("entry.RawTimestamp = %v, want %v", got, want)
	}

	if got, want := entry.Timestamp, (*time.Time)(nil); got != want {
		t.Errorf("entry.Timestamp = %v, want %v", got, want)
	}
}

func TestJournald_PriorityString(t *testing.T) {
	t.Parallel()

	entry := process(t, `{
		"__REALTIME_TIMESTAMP": "1686919896987169",
		"PRIORITY": "cake",
		"MESSAGE": "Invalid number in PRIORITY",
		"SYSLOG_IDENTIFIER": "sshd"
	}`)

	if got, want := entry.Severity, "CAKE"; got != want {
		t.Errorf("entry.Severity = %v, want %v", got, want)
	}

	if got, want := entry.SkipFields["PRIORITY"], true; got != want {
		t.Errorf("entry.SkipFields['PRIORITY'] = %v, want %v", got, want)
	}
}

func TestJournald_PriorityInt(t *testing.T) {
	t.Parallel()

	entry := process(t, `{
		"__REALTIME_TIMESTAMP": "1686919896987169",
		"PRIORITY": 2,
		"MESSAGE": "Actual number type in PRIORITY",
		"SYSLOG_IDENTIFIER": "sshd"
	}`)

	if got, want := entry.Severity, "CRITICAL"; got != want {
		t.Errorf("entry.Severity = %v, want %v", got, want)
	}

	if got, want := entry.SkipFields["PRIORITY"], true; got != want {
		t.Errorf("entry.SkipFields['PRIORITY'] = %v, want %v", got, want)
	}
}

func TestJournald_InvalidPriority(t *testing.T) {
	t.Parallel()

	entry := process(t, `{
		"__REALTIME_TIMESTAMP": "1686919896987169",
		"PRIORITY": {},
		"MESSAGE": "Object type in PRIORITY",
		"SYSLOG_IDENTIFIER": "sshd"
	}`)

	if got, want := entry.Severity, "UNKNOWN"; got != want {
		t.Errorf("entry.Severity = %v, want %v", got, want)
	}

	if got, want := entry.SkipFields["PRIORITY"], false; got != want {
		t.Errorf("entry.SkipFields['PRIORITY'] = %v, want %v", got, want)
	}
}

func process(t *testing.T, input string) *structure.Entry {
	t.Helper()

	raw := &bytes.Buffer{}
	_ = json.Compact(raw, []byte(input))

	line := &stream.Line{Raw: raw.Bytes(), JSON: raw.Bytes()}
	entry := &structure.Entry{SkipFields: make(map[string]bool)}
	djson.Unmarshal(line.JSON, entry)

	p := &JournaldProcessor{}

	detected := p.Detect(line, entry)
	if got, want := detected, true; got != want {
		t.Fatalf("Detect() = %v, want %v", got, want)
	}

	err := p.Process(line, entry)
	if err != nil {
		t.Fatalf("Process() = %v, want nil", err)
	}

	return entry
}
