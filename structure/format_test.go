package structure_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/koenbollen/jl/djson"
	"github.com/koenbollen/jl/structure"
)

const example = `{
  "message": "Hello, world",
  "severity": "30",
  "timestamp": "2015-02-11T13:37:00Z",
  "lang": "fr",
  "labels": {"git_rev": "0992944"},
  "long_number": 22501438,
  "float_number": 2250.1438
}`

func TestHappypath(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	formatter, err := structure.NewFormatter(buf, "")
	if err != nil {
		t.Fatalf("failed to create new formatter: %v", err)
	}

	logline := []byte(strings.Replace(example, "\n", "", -1))
	var entry structure.Entry
	_ = json.Unmarshal(logline, &entry)

	err = formatter.Format(&entry, logline, nil, nil)
	if err != nil {
		t.Fatalf("failed to format entry: %v", err)
	}
	expect1 := "[2015-02-11 13:37:00]    INFO: Hello, world [float_number=2250.1438 git_rev=0992944 lang=fr long_number=22501438]\n"
	if buf.String() != expect1 {
		t.Errorf("\n\tnot match: %q\n\t   expect: %q\n", buf.String(), expect1)
	}

	buf.Reset()
	err = formatter.Format(&entry, logline, []byte("prefix: "), []byte(" suffix!"))
	if err != nil {
		t.Fatalf("failed to format entry: %v", err)
	}
	expect2 := "prefix: [2015-02-11 13:37:00]    INFO: Hello, world [float_number=2250.1438 git_rev=0992944 lang=fr long_number=22501438] suffix!\n"
	if buf.String() != expect2 {
		t.Errorf("\n\tnot match: %q\n\t   expect: %q\n", buf.String(), expect2)
	}
}

func TestNestedFieldPaths(t *testing.T) {
	t.Parallel()

	logline := []byte(`{"message": "Hi!", "severity": "info", "meta": {"count": 42}, "flat.root": "yes"}`)

	buf := &bytes.Buffer{}
	formatter, err := structure.NewFormatter(buf, "")
	if err != nil {
		t.Fatalf("failed to create new formatter: %v", err)
	}

	var entry structure.Entry
	err = json.Unmarshal(logline, &entry)
	if err != nil {
		t.Fatalf("failed to unmarshal logline: %v", err)
	}

	formatter.IncludeFields = []string{"flat.root", "meta.count"}

	err = formatter.Format(&entry, logline, nil, nil)
	if err != nil {
		t.Fatalf("failed to format entry: %v", err)
	}
	expect := "   INFO: Hi! [flat.root=yes meta.count=42]\n"
	if buf.String() != expect {
		t.Errorf("\n\tnot match: %q\n\t   expect: %q\n", buf.String(), expect)
	}
}

func TestIncludePaths(t *testing.T) {
	t.Parallel()

	logline := []byte(`{"service":{"name":"gunicorn"},"@timestamp":"2020-10-23T03:35:49.324754+00:00","message":"Served","time":1603424149.3247535,"log":{"level":"INFO"},"request":{"scheme":"https","path":"/users/users/notices/","method":"GET","customer":"test"},"event":{"duration":78518000}}`)

	buf := &bytes.Buffer{}
	formatter, err := structure.NewFormatter(buf, "")
	if err != nil {
		t.Fatalf("failed to create new formatter: %v", err)
	}

	var entry structure.Entry
	djson.Unmarshal(logline, &entry)

	formatter.IncludeFields = []string{"request.method", "request.path", "event.duration"}

	err = formatter.Format(&entry, logline, nil, nil)
	if err != nil {
		t.Fatalf("failed to format entry: %v", err)
	}
	expect := "[2020-10-23T03:35:49.324754+00:00]    INFO: Served [event.duration=78518000 request.method=GET request.path=/users/users/notices/]\n"
	if buf.String() != expect {
		t.Errorf("\n\tnot match: %q\n\t   expect: %q\n", buf.String(), expect)
	}
}
