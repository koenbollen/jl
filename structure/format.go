package structure

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

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

var defaultExcludes = []string{
	"@timestamp", "hostname", "level", "message", "msg", "name", "pid", "severity", "text", "time", "timestamp", "ts", "v",
}

// NewLine contains ['\n']
var NewLine = []byte("\n")

// Formatter is the system that outputs a sturctured log entry as a nice
// readable line using  a small go template (which could be given via the cli)
type Formatter struct {
	output   io.Writer
	template *template.Template

	Colorize      bool
	ShowFields    bool
	ShowPrefix    bool
	ShowSuffix    bool
	IncludeFields string
	ExcludeFields []string
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
		output:        w,
		template:      tmpl,
		Colorize:      false,
		ShowFields:    true,
		ShowPrefix:    true,
		ShowSuffix:    true,
		IncludeFields: "",
		ExcludeFields: defaultExcludes,
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

	if entry.Timestamp != nil && entry.Timestamp.Year() > 3000 { // timestamp was probably in milliseconds
		t := *entry.Timestamp
		t = time.Unix(t.Unix()/int64(time.Second/time.Millisecond), 0).UTC()
		entry.Timestamp = &t
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
		if padding > 0 {
			entry.Severity = strings.Repeat(" ", padding) + entry.Severity
		}
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
		delete(fields, "labels")
	}

	output := make([]string, 0)
	if err == nil {
		path := ""
		for key, value := range walkFields(fields, "") {
			if _, ok := value.(map[string]interface{}); ok {
				continue
			}
			if _, ok := value.([]interface{}); ok {
				continue
			}
			if !f.shouldSkipField(key, path+"."+key, value) {
				switch v := value.(type) {
				case float64:
					output = append(output, key+"="+strconv.FormatFloat(v, 'f', -1, 64))
				default:
					output = append(output, fmt.Sprintf("%s=%v", key, value))
				}
			}
		}
		if len(output) > 0 {
			sort.Strings(output)
			fmt.Fprintf(f.output, " %v", output)
		}
	}
}

func (f *Formatter) shouldSkipField(field, path string, value interface{}) bool {
	if strings.Contains(f.IncludeFields, field) || strings.Contains(f.IncludeFields, path) {
		return false
	}
	if strings.Count(path, ".") > 1 { // Only include nested fields when the are in the IncludeFields
		return true
	}
	if len(path+fmt.Sprintf("%v", value)) >= 30 {
		return true
	}

	return contains(f.ExcludeFields, field)
}

func contains(lst []string, val string) bool {
	for _, i := range lst {
		if strings.EqualFold(i, val) {
			return true
		}
	}
	return false
}

func walkFields(fields map[string]interface{}, path string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range fields {
		if path != "" {
			key = path + "." + key
		}
		if nested, ok := value.(map[string]interface{}); ok {
			for k, v := range walkFields(nested, key) {
				result[k] = v
			}
		} else {
			result[key] = value
		}
	}
	return result
}
