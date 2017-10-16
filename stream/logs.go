package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"text/scanner"
)

// Line represents a line from the given Reader of a Stream, containing the
// raw bytes and the RawMessage of JSON if present.
type Line struct {
	Raw  []byte
	JSON json.RawMessage

	Prefix []byte
	Suffix []byte
}

// Stream lets you scan through the lines of a io.Reader and return each line
// as a Line struct, containing the raw bytes and the JSON bytes if present.
// Lines parsed are exposed byt the Lines() method.
type Stream interface {
	Close()
	Lines() <-chan *Line
	Err() error
}

type stream struct {
	scanner *bufio.Scanner
	result  chan *Line
	stop    chan struct{}
}

// New will construct a new Stream and start it.
func New(r io.Reader) Stream {
	scanner := bufio.NewScanner(r)
	l := &stream{
		scanner: scanner,
		result:  make(chan *Line),
		stop:    make(chan struct{}),
	}
	go l.run()
	return l
}

func (l *stream) run() {
	for l.scanner.Scan() {
		raw := l.scanner.Bytes()
		json := l.parse(raw)
		prefix, suffix := split(raw, json)
		line := &Line{
			Raw:    raw,
			JSON:   json,
			Prefix: prefix,
			Suffix: suffix,
		}
		select {
		case <-l.stop:
			return
		case l.result <- line:
			continue
		}
	}
	close(l.result)
}

func (l *stream) parse(raw []byte) json.RawMessage {
	var s scanner.Scanner
	s.Init(bytes.NewReader(raw))
	s.Error = func(s *scanner.Scanner, msg string) {}
	depth := 0
	start := -1
	end := -1
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == '{' {
			if depth == 0 {
				start = s.Position.Offset
			}
			depth++
		}
		if tok == '}' {
			depth--
			if depth == 0 {
				end = s.Position.Offset + 1
				break
			}
		}
	}
	if start != -1 && end != -1 {
		slice := raw[start:end]
		var v interface{}
		err := json.Unmarshal(slice, &v)
		if err == nil {
			return slice
		}
	}
	return nil
}

func (l *stream) Close() {
	l.stop <- struct{}{}
	close(l.result)
}

func (l *stream) Lines() <-chan *Line {
	return l.result
}

func (l *stream) Err() error {
	return l.scanner.Err()
}

func split(raw, json []byte) (prefix, suffix []byte) {
	prefix, suffix = nil, nil
	if len(json) > 0 {
		parts := bytes.SplitN(raw, json, 2)
		if len(parts) >= 1 && len(parts[0]) > 0 {
			prefix = parts[0]
		}
		if len(parts) >= 2 && len(parts[1]) > 0 {
			suffix = parts[1]
		}
	}
	return
}
