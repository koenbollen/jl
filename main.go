package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/koenbollen/jl/djson"
	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
)

func main() {
	files, color, showPrefix, showSuffix, showFields := cli()
	formatter, err := structure.NewFormatter(os.Stdout, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid format: %v\n", err)
		os.Exit(1)
	}

	formatter.Colorize = color
	formatter.ShowPrefix = showPrefix
	formatter.ShowSuffix = showSuffix
	formatter.ShowFields = showFields

	r, err := openFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	s := stream.New(r)
	for line := range s.Lines() {
		var err error
		entry := &structure.Entry{}
		if line.JSON != nil && len(line.JSON) > 0 {
			err = json.Unmarshal(line.JSON, entry)
			djson.Unmarshal(line.JSON, entry)
		}

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
			fmt.Fprintf(os.Stderr, "broken pipe: %v\n", err)
			break
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "broken pipe: %v\n", err)
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

func openFiles(files []string) (io.Reader, error) {
	var filtered []string
	for _, file := range files {
		if file != "" {
			filtered = append(filtered, file)
		}
	}
	if len(filtered) == 0 {
		return os.Stdin, nil
	}
	readers := make([]io.Reader, 0)
	for _, file := range filtered {
		if file == "-" {
			readers = append(readers, os.Stdin)
		} else {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}
			readers = append(readers, f)
		}
	}
	return io.MultiReader(readers...), nil
}
