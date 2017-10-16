package structure_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/koenbollen/jl/structure"
)

const example = `{
  "message": "Hello, world",
  "severity": "30",
  "timestamp": "2015-02-11T13:37:00Z",
  "lang": "fr",
  "labels": {"git_rev": "0992944"}
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
	if buf.String() != "[2015-02-11 13:37:00] INFO: Hello, world [git_rev=0992944 lang=fr]\n" {
		t.Errorf("invalid output: %q", buf.String())
	}

	buf.Reset()
	err = formatter.Format(&entry, logline, []byte("prefix: "), []byte(" suffix!"))
	if err != nil {
		t.Fatalf("failed to format entry: %v", err)
	}
	if buf.String() != "prefix: [2015-02-11 13:37:00] INFO: Hello, world [git_rev=0992944 lang=fr] suffix!\n" {
		t.Errorf("invalid output: %q", buf.String())
	}
}
