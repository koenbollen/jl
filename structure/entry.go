package structure

import "time"

// Entry represents a structured logline to be formatted.
type Entry struct {
	Timestamp      *time.Time `djson:"timestamp,@timestamp,time,date,ts"`
	RawTimestamp   string     `djson:"timestamp,@timestamp,time,date,ts"`
	FloatTimestamp float64    `djson:"timestamp,@timestamp,time,date,ts"`
	Severity       string     `djson:"severity,level,log.level"`
	Message        string     `djson:"message,msg,text,*.message"`

	Name string `djson:"app,name,service.name"`

	// IncludeFields is used by processors to indicate which fields should be included
	IncludeFields []string

	// ExcludeFields is used by processors to indicate which fields should be skipped
	ExcludeFields []string
}
