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

func TestFloatTimestamp(t *testing.T) {
	in := `{"level":"info","ts":1565361391.4279764,"caller":"ingress/main.go:109","msg":"Hi","environment":"production","artifact":"workflow","component":"ingress"}`
	e := structure.Entry{}
	djson.Unmarshal([]byte(in+"\n"), &e)
	if e.Timestamp != nil && !e.Timestamp.IsZero() {
		t.Errorf("expected .Timestamp to be nil: %v", e.Timestamp)
	}
	expectedRawTimestamp := "1.5653613914279764e+09"
	if e.RawTimestamp != expectedRawTimestamp {
		t.Errorf("expected .RawTimestamp to be %q, got %q", expectedRawTimestamp, e.RawTimestamp)
	}
	expectedFloatTimestamp := 1565361391.4279764
	if e.FloatTimestamp != expectedFloatTimestamp {
		t.Errorf("expected .RawTimestamp to be %v, got %v", expectedFloatTimestamp, e.FloatTimestamp)
	}
	if e.Message != "Hi" {
		t.Errorf("invalid message parsed: %v", e.Message)
	}
}
