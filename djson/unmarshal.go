package djson

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func Unmarshal(data []byte, val interface{}) {
	elem := reflect.TypeOf(val).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.FieldByIndex([]int{i})
		keys, ok := field.Tag.Lookup("djson")
		if !ok {
			continue
		}
		for _, key := range strings.Split(keys, ",") {
			result := gjson.GetBytes(data, strings.TrimSpace(key))
			if !result.Exists() {
				continue
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
