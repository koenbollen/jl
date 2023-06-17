package processors

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/koenbollen/jl/djson"
	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
)

func TestNested(t *testing.T) {

	in := `{
		"level": "INFO",
		"nested": {
			"message": "Hi",
			"user": 56
		}
	}`

	raw := &bytes.Buffer{}
	_ = json.Compact(raw, []byte(in))

	line := &stream.Line{Raw: raw.Bytes(), JSON: raw.Bytes()}
	entry := &structure.Entry{Message: "Hi"}
	djson.Unmarshal(line.JSON, entry)

	p := &NestedProcessor{}

	detected := p.Detect(line, entry)
	if got, want := detected, true; got != want {
		t.Errorf("Detect() = %v, want %v", got, want)
	}

	err := p.Process(line, entry)
	if err != nil {
		t.Errorf("Process() = %v, want nil", err)
	}

	if got, want := entry.Message, "Hi"; got != want {
		t.Errorf("entry.Message = %v, want %v", got, want)
	}
	if got, want := entry.IncludeFields[0], "nested"; got != want {
		t.Errorf("entry.IncludeFields = %v, want %v", got, want)
	}
}
