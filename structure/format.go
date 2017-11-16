package structure

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

// DefaultTemplate is used when no template is given.
const DefaultTemplate = `{{if .Timestamp}}[{{.Timestamp.Format "2006-01-02 15:04:05"}}] {{else if .RawTimestamp}}[{{.RawTimestamp}}] {{end}}{{if .Severity}}{{.Severity}}: {{end}}{{.Message}}`

var severityMapping = map[string]string{
	"10":   "TRACE",
	"20":   "DEBUG",
	"30":   "INFO",
	"40":   "WARNING",
	"WARN": "WARNING",
	"50":   "ERROR",
	"60":   "FATAL",
}

var fieldsToSkip = []string{
	"@timestamp", "hostname", "level", "message", "msg", "name", "pid", "severity", "text", "time", "timestamp", "v",
}

// NewLine contains ['\n']
var NewLine = []byte("\n")

// Formatter is the system that outputs a sturctured log entry as a nice
// readable line using  a small go template (which could be given via the cli)
type Formatter struct {
	output   io.Writer
	template *template.Template

	Colorize   bool
	ShowFields bool
	ShowPrefix bool
	ShowSuffix bool
}

// NewFormatter compiles the given fmt as a go template and returns a Formatter
func NewFormatter(w io.Writer, fmt string) (*Formatter, error) {
	if fmt == "" {
		fmt = DefaultTemplate
	}
	tmpl, err := template.New("out").Parse(fmt)
	if err != nil {
		return nil, err
	}

	return &Formatter{
		output:     w,
		template:   tmpl,
		Colorize:   false,
		ShowFields: true,
		ShowPrefix: true,
		ShowSuffix: true,
	}, nil
}

// Format takes a structured log entry and formats it according the template.
func (f *Formatter) Format(entry *Entry, raw json.RawMessage, prefix, suffix []byte) error {
	color.NoColor = !f.Colorize
	f.enhance(entry)

	err := f.outputSimple(prefix, f.ShowPrefix)
	if err != nil {
		return err
	}

	err = f.template.Execute(f.output, entry)
	if err != nil {
		return err
	}

	f.outputFields(entry, raw)

	err = f.outputSimple(suffix, f.ShowSuffix)
	if err != nil {
		return err
	}

	err = stacktrace(f.output, raw)
	if err != nil {
		return err
	}

	_, err = f.output.Write(NewLine)
	return err
}

func (f *Formatter) enhance(entry *Entry) {
	if entry.Timestamp != nil && entry.Timestamp.IsZero() {
		entry.Timestamp = nil
	}

	entry.Severity = strings.ToUpper(entry.Severity)
	if level, ok := severityMapping[entry.Severity]; ok {
		entry.Severity = level
	}
	if entry.Severity != "" {
		padding := 7 - len(entry.Severity)
		if color, ok := severityColors[entry.Severity]; ok {
			entry.Severity = color(entry.Severity)
		}
		entry.Severity = strings.Repeat(" ", padding) + entry.Severity
	}

	entry.Message = messageColor(entry.Message)
}

func (f *Formatter) outputSimple(txt []byte, toggle bool) error {
	if toggle && txt != nil && len(txt) > 0 {
		_, err := f.output.Write(txt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) outputFields(entry *Entry, raw json.RawMessage) {
	if !f.ShowFields {
		return
	}
	fields := make(map[string]interface{})
	err := json.Unmarshal(raw, &fields)

	if labels, ok := fields["labels"]; ok {
		for k, v := range labels.(map[string]interface{}) {
			fields[k] = v
		}
	}

	output := make([]string, 0)
	if err == nil {
		for key, value := range fields {
			if _, ok := value.(map[string]interface{}); ok {
				continue
			}
			if _, ok := value.([]interface{}); ok {
				continue
			}
			if !shouldSkipField(key, value) {
				output = append(output, fmt.Sprintf("%s=%v", key, value))
			}
		}
		if len(output) > 0 {
			sort.Strings(output)
			fmt.Fprintf(f.output, " %v", output)
		}
	}
}

func shouldSkipField(field string, value interface{}) bool {
	if len(fmt.Sprintf("%v", value)) >= 24 {
		return true
	}
	ix := sort.SearchStrings(fieldsToSkip, field)
	return ix < len(fieldsToSkip) && fieldsToSkip[ix] == field
}
