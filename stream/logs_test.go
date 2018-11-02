package stream_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/koenbollen/jl/stream"
)

func TestStream(t *testing.T) {
	t.Parallel()
	in := `first line is normal
{"second": "line is json"}
{"third": "as well"}
forth line: {"is": "mixed"}
fifth is {broken
{"then": "we have"} trailing text
json in {"the": "middle"} of the line`
	expected := []*stream.Line{
		{Raw: []byte("first line is normal"), JSON: nil, Prefix: nil, Suffix: nil},
		{Raw: []byte(`{"second": "line is json"}`), JSON: json.RawMessage(`{"second": "line is json"}`), Prefix: nil, Suffix: nil},
		{Raw: []byte(`{"third": "as well"}`), JSON: json.RawMessage(`{"third": "as well"}`), Prefix: nil, Suffix: nil},
		{Raw: []byte(`forth line: {"is": "mixed"}`), JSON: json.RawMessage(`{"is": "mixed"}`), Prefix: []byte(`forth line: `), Suffix: nil},
		{Raw: []byte(`fifth is {broken`), JSON: nil, Prefix: nil, Suffix: nil},
		{Raw: []byte(`{"then": "we have"} trailing text`), JSON: json.RawMessage(`{"then": "we have"}`), Prefix: nil, Suffix: []byte(` trailing text`)},
		{Raw: []byte(`json in {"the": "middle"} of the line`), JSON: json.RawMessage(`{"the": "middle"}`), Prefix: []byte(`json in `), Suffix: []byte(` of the line`)},
	}
	s := stream.New(strings.NewReader(in))
	for i, line := range expected {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := <-s.Lines()
			if !reflect.DeepEqual(result, line) {
				t.Errorf("line %d didnt match, got %q expected %q", i, result, line)
			}
		})
	}
}

func test(t *testing.T, input string, expected *stream.Line) {
	t.Parallel()
	t.Helper()
	s := stream.New(strings.NewReader(input))
	result := <-s.Lines()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("line didnt match, got %q expected %q", result, expected)
	}
}

func TestFullJSON(t *testing.T) {
	test(t, `{"msg": "Hello", "key": "value"}`, &stream.Line{
		Raw:  []byte(`{"msg": "Hello", "key": "value"}`),
		JSON: json.RawMessage(`{"msg": "Hello", "key": "value"}`),
	})
}

func TestPlainText(t *testing.T) {
	test(t, `Hello, world!!`, &stream.Line{
		Raw:  []byte(`Hello, world!!`),
		JSON: nil,
	})
}

func TestTrailingText(t *testing.T) {
	test(t, `{"json": 1} Hello, world!!`, &stream.Line{
		Raw:    []byte(`{"json": 1} Hello, world!!`),
		JSON:   json.RawMessage(`{"json": 1}`),
		Suffix: []byte(` Hello, world!!`),
	})
}

func TestLeadingText(t *testing.T) {
	test(t, `Sup? {"json": 2}`, &stream.Line{
		Raw:    []byte(`Sup? {"json": 2}`),
		JSON:   json.RawMessage(`{"json": 2}`),
		Prefix: []byte(`Sup? `),
	})
}

func TestBacktick(t *testing.T) {
	test(t, "Single backtick: `", &stream.Line{
		Raw: []byte("Single backtick: `"),
	})
}

func TestClose(t *testing.T) {
	t.Parallel()
	s := stream.New(iotest.TimeoutReader(strings.NewReader("one\ntwo\n")))
	s.Close()
	line := <-s.Lines()
	if line != nil {
		t.Error("expected nil on a closed stream")
	}
}

func TestErrors(t *testing.T) {
	t.Parallel()
	s := stream.New(iotest.TimeoutReader(strings.NewReader("Hi!\n")))
	<-s.Lines()
	err := s.Err()
	if err != iotest.ErrTimeout {
		t.Errorf("expecting ErrTimeout, got %v", err)
	}
}
