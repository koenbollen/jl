package stacktracers

import (
	"strings"

	"github.com/koenbollen/jl/structure"
)

type bunyan struct {
}

func init() {
	structure.RegisterStacktracer(&bunyan{})
}

func (b *bunyan) Detect(json map[string]interface{}) bool {
	if _, ok := json["v"].(float64); !ok {
		return false
	}

	err, ok := json["err"].(map[string]interface{})
	if !ok {
		return false
	}

	_, ok = err["stack"].(string)
	return ok
}

func (b *bunyan) Format(json map[string]interface{}) string {
	err := json["err"].(map[string]interface{})
	stack := err["stack"].(string)
	stack = strings.TrimSpace(stack)
	stack = "\n    " + strings.Replace(stack, "\n", "\n    ", -1)
	return stack
}
