package structure

import (
	"encoding/json"
	"io"
)

// StacktraceFormatter interfaces with the Formatter to format a possible
// stacktrace in a JSON log line. The Detect method returns true if it's
// compatible.
type StacktraceFormatter interface {
	Detect(json map[string]interface{}) bool
	Format(json map[string]interface{}) string
}

var stacktracers []StacktraceFormatter

// RegisterStacktracer adds a StacktraceFormatter to a thé list.
func RegisterStacktracer(tracer StacktraceFormatter) {
	stacktracers = append(stacktracers, tracer)
}

func stacktrace(w io.Writer, raw []byte) error {
	var root map[string]interface{}
	_ = json.Unmarshal(raw, &root)
	for _, tracer := range stacktracers {
		if tracer.Detect(root) {
			stack := tracer.Format(root)
			_, err := w.Write([]byte(stack))
			return err
		}
	}
	return nil
}
