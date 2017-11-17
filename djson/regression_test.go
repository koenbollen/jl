package djson_test

import (
	"testing"

	"github.com/koenbollen/jl/djson"
	"github.com/koenbollen/jl/structure"
)

func TestInvalidTimestampFormat(t *testing.T) {
	e := structure.Entry{}
	in := `{"message":"Hi","timestamp":"2017-11-07T15:34:32+0100"}` // invalid timestamp
	djson.Unmarshal([]byte(in+"\n"), &e)
	if e.Timestamp != nil && !e.Timestamp.IsZero() {
		t.Errorf("expected .Timestamp to be nil: %v", e.Timestamp)
	}
	if e.RawTimestamp != "2017-11-07T15:34:32+0100" {
		t.Errorf("expected .RawTimestamp not set: %q", e.RawTimestamp)
	}
	if e.Message != "Hi" {
		t.Errorf("invalid message parsed: %v", e.Message)
	}
}
