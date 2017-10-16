package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/koenbollen/jl/djson"
	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
)

func main() {
	formatter, err := structure.NewFormatter(os.Stdout, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid format: %v", err)
	}

	s := stream.New(os.Stdin)
	for line := range s.Lines() {
		entry := &structure.Entry{}
		err := json.Unmarshal(line.JSON, entry)
		djson.Unmarshal(line.JSON, entry)

		// unable to parse entry, outputting raw line:
		if err != nil || entry.Message == "" {
			os.Stdout.Write(line.Raw)
			os.Stdout.Write(structure.NewLine)
			continue
		}

		// Passing entry to formatter to output:
		prefix, suffix := split(line.Raw, line.JSON)
		err = formatter.Format(entry, line.JSON, prefix, suffix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "broken pipe: %v", err)
			break
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "broken pipe: %v", err)
	}
}

func split(raw, json []byte) (prefix, suffix []byte) {
	parts := bytes.SplitN(raw, json, 2)
	prefix = parts[0]
	if len(parts) >= 2 {
		suffix = parts[1]
	}
	return
}
