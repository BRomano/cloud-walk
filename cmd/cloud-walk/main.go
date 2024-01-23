package main

import (
	"cloud-walk/internal/domain/service"
	"cloud-walk/internal/infra/repository"
	"context"
	"log/slog"
	"os"
	"sync"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("version", Version, "commit", Commit, "timestamp", Timestamp)
	slog.SetDefault(logger)
}

type app struct {
	settings *Settings
}

func newApp() app {
	return app{
		settings: &Settings{},
	}
}

func (app *app) NewService() service.LogParserService {
	return service.NewLogParser(repository.LogParserFactory)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := newApp()

	if err := app.settings.initSettings(); err != nil {
		slog.Error("error initializing settings", "err", err)
		panic("error initializing settings")
	}

	service := app.NewService()
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go app.startHTTPServer(ctx, &waitGroup, service)

	waitGroup.Wait()
	slog.Info("service finish gracefully")
}
