package processors

import (
	"strconv"
	"strings"
	"time"

	"github.com/koenbollen/jl/stream"
	"github.com/koenbollen/jl/structure"
	"github.com/tidwall/gjson"
)

var priorityMapping = map[int64]string{
	0: "EMERGENCY",
	1: "ALERT",
	2: "CRITICAL",
	3: "ERROR",
	4: "WARNING",
	5: "NOTICE",
	6: "INFO",
	7: "DEBUG",
}

type JournaldProcessor struct {
}

func (p *JournaldProcessor) Detect(line *stream.Line, entry *structure.Entry) bool {
	return gjson.GetBytes(line.JSON, "SYSLOG_IDENTIFIER").Exists() && gjson.GetBytes(line.JSON, "__REALTIME_TIMESTAMP").Exists() && gjson.GetBytes(line.JSON, "PRIORITY").Exists()
}

func (p *JournaldProcessor) Process(line *stream.Line, entry *structure.Entry) error {
	entry.Message = gjson.GetBytes(line.JSON, "MESSAGE").String()
	entry.SkipFields["MESSAGE"] = true

	rawTimestamp := gjson.GetBytes(line.JSON, "__REALTIME_TIMESTAMP").String()
	micro, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		entry.RawTimestamp = rawTimestamp
	} else {
		t := time.Unix(micro/1_000_000, micro%1_000_000*1_000).UTC()
		entry.Timestamp = &t
	}
	entry.SkipFields["__REALTIME_TIMESTAMP"] = true

	prioField := gjson.GetBytes(line.JSON, "PRIORITY")
	switch prioField.Type {
	case gjson.Number:
		priority := prioField.Int()
		entry.Severity = priorityMapping[priority]
		entry.SkipFields["PRIORITY"] = true
	case gjson.String:
		priority := prioField.String()
		if i, err := strconv.ParseInt(priority, 10, 64); err == nil {
			entry.Severity = priorityMapping[i]
			entry.SkipFields["PRIORITY"] = true
			break
		} else if len(priority) < 12 {
			entry.Severity = strings.ToUpper(priority)
			entry.SkipFields["PRIORITY"] = true
			break
		}
		fallthrough
	default:
		entry.Severity = "UNKNOWN"
	}

	return nil
}
