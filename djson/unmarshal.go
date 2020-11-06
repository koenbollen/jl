package djson

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Unmarshal will try to load JSON from the given data into val. It'll read
// from the struct tags of the given val and look for the 'djson' tag which
// can supply multiple possible fields a JSON key can be. If a json key match
// with any of the tags it'll set the value.
func Unmarshal(data []byte, val interface{}) {
	elem := reflect.TypeOf(val).Elem()
	for i := 0; i < elem.NumField(); i++ {
		process(data, val, elem, i)
	}
}

func process(data []byte, val interface{}, elem reflect.Type, i int) {
	defer func() {
		_ = recover()
	}()
	field := elem.FieldByIndex([]int{i})
	keys, ok := field.Tag.Lookup("djson")
	if !ok {
		return
	}
	for _, key := range strings.Split(keys, ",") {
		result := gjson.GetBytes(data, strings.TrimSpace(key))
		if !result.Exists() {
			result = gjson.GetBytes(data, strings.ReplaceAll(strings.TrimSpace(key), ".", "\\."))
			if !result.Exists() {
				continue
			}
		}
		value := reflect.ValueOf(result.Value())
		fieldValue := reflect.ValueOf(val).Elem().FieldByIndex(field.Index)
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			fieldValue.Elem().Set(convert(value, fieldValue.Elem()))
		} else {
			fieldValue.Set(convert(value, fieldValue))
		}
	}
}

func convert(value, field reflect.Value) reflect.Value {
	if value.Type() == field.Type() {
		return value
	}
	switch field.Interface().(type) {
	case string:
		return reflect.ValueOf(fmt.Sprintf("%v", value.Interface()))
	case time.Time:
		if value.Kind() == reflect.String {
			t, err := time.Parse(time.RFC3339, value.String())
			if err == nil {
				return reflect.ValueOf(t)
			}
		}
	}
	return reflect.ValueOf(nil)
}
