package stacktracers

import (
	"strings"

	"github.com/koenbollen/jl/structure"
)

type zap struct {
}

func init() {
	structure.RegisterStacktracer(&zap{})
}

func (b *zap) Detect(json map[string]interface{}) bool {

	_, ok := json["error"].(string)
	if !ok {
		return false
	}

	for _, value := range json {
		if str, ok := value.(string); ok {
			if strings.HasPrefix(str, "go.uber.org/zap.Stack") {
				return true
			}
		}
	}

	return false
}

func (b *zap) Format(json map[string]interface{}) string {
	err := json["error"].(string)

	stack := ""
	for _, value := range json {
		if str, ok := value.(string); ok {
			if strings.HasPrefix(str, "go.uber.org/zap.Stack") {
				stack = str
				break
			}
		}
	}

	stack = strings.TrimSpace(stack)
	stack = strings.Replace(stack, "\t", "  ", -1)
	stack = "    " + strings.Replace(stack, "\n", "\n    ", -1)
	return "\n    " + err + "\n" + stack
}
