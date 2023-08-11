package main

import (
	"io"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: FixedTime}))
	logger.Info("hello", "count", 3)
	logger.Warn("failed", "err", io.EOF)
}

func FixedTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		t, _ := time.Parse(time.DateTime, time.DateTime)
		a.Value = slog.TimeValue(t)
	}
	return a
}
