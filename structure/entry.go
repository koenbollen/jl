package structure

import "time"

// Entry represents a structured logline to be formatted.
type Entry struct {
	Timestamp *time.Time `djson:"timestamp,@timestamp,time,date"`
	Severity  string     `djson:"severity,level"`
	Message   string     `djson:"message,msg,text"`

	Name string `djson:"app,name"`
}
