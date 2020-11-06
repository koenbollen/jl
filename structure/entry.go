package structure

import "time"

// Entry represents a structured logline to be formatted.
type Entry struct {
	Timestamp      *time.Time `djson:"timestamp,@timestamp,time,date,ts"`
	RawTimestamp   string     `djson:"timestamp,@timestamp,time,date,ts"`
	FloatTimestamp float64    `djson:"timestamp,@timestamp,time,date,ts"`
	Severity       string     `djson:"severity,level,log.level"`
	Message        string     `djson:"message,msg,text"`

	Name string `djson:"app,name,service.name"`
}
