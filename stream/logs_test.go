package stream_test

import (
	"bytes"
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
		{[]byte("first line is normal"), nil, nil, nil},
		{[]byte(`{"second": "line is json"}`), json.RawMessage(`{"second": "line is json"}`), nil, nil},
		{[]byte(`{"third": "as well"}`), json.RawMessage(`{"third": "as well"}`), nil, nil},
		{[]byte(`forth line: {"is": "mixed"}`), json.RawMessage(`{"is": "mixed"}`), []byte(`forth line: `), nil},
		{[]byte(`fifth is {broken`), nil, nil, nil},
		{[]byte(`{"then": "we have"} trailing text`), json.RawMessage(`{"then": "we have"}`), nil, []byte(` trailing text`)},
		{[]byte(`json in {"the": "middle"} of the line`), json.RawMessage(`{"the": "middle"}`), []byte(`json in `), []byte(` of the line`)},
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
	if bytes.Compare(result.Raw, expected.Raw) != 0 || bytes.Compare(result.JSON, expected.JSON) != 0 {
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
		Raw:  []byte(`{"json": 1} Hello, world!!`),
		JSON: json.RawMessage(`{"json": 1}`),
	})
}

func TestLeadingText(t *testing.T) {
	test(t, `Sup? {"json": 2}`, &stream.Line{
		Raw:  []byte(`Sup? {"json": 2}`),
		JSON: json.RawMessage(`{"json": 2}`),
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
