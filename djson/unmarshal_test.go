package djson_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/koenbollen/jl/djson"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	val := struct {
		Message string `djson:"msg,message"`
		NotUsed bool
	}{}
	djson.Unmarshal([]byte(`{"msg": "Hi"}`), &val)
	if val.Message != "Hi" {
		t.Error("failed to set .Message to 'Hi'")
	}

	_ = json.Unmarshal([]byte(`{"message": "Hi"}`), &val)
	if val.Message != "Hi" {
		t.Error("failed to set .Message to 'Hi'")
	}

	if val.NotUsed {
		t.Error("unused field .NotUsed somehow set to true")
	}
}

func TestInvalidJSON(t *testing.T) {
	t.Parallel()
	val := struct {
		Name string `djson:"test,name"`
	}{}
	djson.Unmarshal([]byte(`{"name`), &val)
	if val.Name != "" {
		t.Error(".Name was set")
	}
}

func TestPointerValue(t *testing.T) {
	t.Parallel()
	val := struct {
		Time *string `djson:"time,timestamp"`
	}{}
	djson.Unmarshal([]byte(`{"notused": "test"}`), &val)
	if val.Time != nil {
		t.Error(".Time somehow set")
	}
	djson.Unmarshal([]byte(`{"timestamp": "2011-12-18T20:46:00Z"}`), &val)
	if val.Time == nil {
		t.Error(".Time not set")
	}
}

func TestNested(t *testing.T) {
	t.Parallel()
	val := struct {
		Level1 string `djson:"lvl1,level1"`
		Level2 string `djson:"nested.lvl2,nested.level2,obj.lvl2"`
	}{}
	djson.Unmarshal([]byte(`{"lvl1": "1", "nested": {"level2": "2"}}`), &val)
	if val.Level1 != "1" {
		t.Error("failed to set .Level1")
	}
	if val.Level2 != "2" {
		t.Error("failed to set nested .Level2")
	}

	djson.Unmarshal([]byte(`{"obj": {"lvl2": "3"}}`), &val)
	if val.Level2 != "3" {
		t.Error("failed to set nested .Level2")
	}
}

func TestFloatToString(t *testing.T) {
	t.Parallel()
	val := struct {
		Level string `djson:"level"`
	}{}
	djson.Unmarshal([]byte(`{"level": "1"`), &val)
	if val.Level != "1" {
		t.Error("failed to set .Level from string")
	}
	djson.Unmarshal([]byte(`{"level": 2`), &val)
	if val.Level != "2" {
		t.Error("failed to set .Level from float64")
	}
	djson.Unmarshal([]byte(`{"level": false`), &val)
	if val.Level != "false" {
		t.Error("failed to set .Level from bool")
	}
}

func TestPointerCast(t *testing.T) {
	t.Parallel()
	val := struct {
		Level *string `djson:"level"`
	}{}
	djson.Unmarshal([]byte(`{"level": 1`), &val)
	if *val.Level != "1" {
		t.Error("failed to set .Level from float64")
	}
}

func TestWeirdTimestamp(t *testing.T) {
	t.Parallel()
	val := struct {
		Timestamp time.Time `djson:"time,t"`
	}{}

	testtime := time.Date(2017, 10, 16, 12, 05, 58, 00, time.UTC)
	djson.Unmarshal([]byte(`{"time": "2017-10-16T12:05:58Z"}`), &val)
	if val.Timestamp != testtime {
		t.Error("failed to set .Timestamp from RFC3339")
	}

	testtime = time.Date(2013, 01, 04, 18, 46, 23, 851000000, time.UTC)
	djson.Unmarshal([]byte(`{"time": "2013-01-04T18:46:23.851Z"}`), &val)
	if val.Timestamp != testtime {
		t.Error("failed to set .Timestamp from RFC3339Nano")
	}
}

func TestDoubleNestedType(t *testing.T) {
	t.Parallel()
	val := struct {
		Level string `djson:"log.level"`
	}{}

	djson.Unmarshal([]byte(`{"log.level": "info"}`), &val)
	if val.Level != "info" {
		t.Error("failed to set .Level from static flat path")
	}

	djson.Unmarshal([]byte(`{"log": {"level": "warning"}}`), &val)
	if val.Level != "warning" {
		t.Error("failed to set .Level from nested path")
	}
}

func TestWildcardNested(t *testing.T) {
	t.Parallel()
	val := struct {
		Message string `djson:"message,*.message"`
	}{}

	djson.Unmarshal([]byte(`{"fields.message": "Hello"}`), &val)
	if val.Message != "Hello" {
		t.Error("failed to set .Message from static flat path")
	}

	djson.Unmarshal([]byte(`{"fields": {"message": "Hello"}}`), &val)
	if val.Message != "Hello" {
		t.Error("failed to set .Message from nested path")
	}

	djson.Unmarshal([]byte(`{"message": "First", "fields.message": "Second"}`), &val)
	if val.Message != "First" {
		t.Error("failed to set .Message from first key")
	}
}
