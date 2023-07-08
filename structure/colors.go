package structure

import "github.com/fatih/color"

var messageColor = color.New(color.FgHiCyan, color.Bold).SprintFunc()
var severityColors = map[string]func(a ...interface{}) string{
	"TRACE":   color.New(color.FgHiBlack).SprintFunc(),
	"DEBUG":   color.New(color.FgHiBlack).SprintFunc(),
	"INFO":    color.New(color.FgCyan).SprintFunc(),
	"WARNING": color.New(color.FgRed).SprintFunc(),
	"ERROR":   color.New(color.FgHiRed, color.Bold).SprintFunc(),
	"FATAL":   color.New(color.FgHiRed, color.Bold).SprintFunc(),
}
