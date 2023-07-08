package processors

import (
	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
)

type Processor interface {
	Detect(line *stream.Line, entry *structure.Entry) bool
	Process(line *stream.Line, entry *structure.Entry) error
}

var All = []Processor{
	&JournaldProcessor{},
}
