package processors

import (
	"strings"

	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
	"github.com/tidwall/gjson"
)

type NestedProcessor struct {
}

func (p *NestedProcessor) Detect(line *stream.Line, entry *structure.Entry) bool {
	result := gjson.GetBytes(line.JSON, "*.message")
	return result.Exists() && result.String() != "" && entry.Message == result.String()
}

func (p *NestedProcessor) Process(line *stream.Line, entry *structure.Entry) error {
	result := gjson.GetBytes(line.JSON, "*.message")
	path := result.Path(string(line.JSON))
	if field, _, found := strings.Cut(path, ".message"); found {
		entry.IncludeFields = append(entry.IncludeFields, field)
		entry.ExcludeFields = append(entry.ExcludeFields, path)
	}
	return nil
}
