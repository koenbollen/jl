package main

import (
	"io"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("hello", "count", 3)
	logger.Warn("failed", "err", io.EOF)
}
