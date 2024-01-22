package main

import (
	"log/slog"
	"os"
)

var (
	Version   string
	Commit    string
	Timestamp string
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("version", Version, "commit", Commit, "timestamp", Timestamp)
	slog.SetDefault(logger)
}

func main() {
	slog.Info("hello world", "count", 1)

}
