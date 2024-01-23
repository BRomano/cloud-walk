package main

import (
	"cloud-walk/internal/domain/service"
	"cloud-walk/internal/handler"
	"context"
	"errors"
	"fmt"
	"github.com/dimiro1/health"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func (app *app) startHTTPServer(ctx context.Context, waitGroup *sync.WaitGroup, service service.LogParserService) {
	defer waitGroup.Done()

	app.setupHTTPRoutes(ctx, service)

	slog.Info("Starting HTTP server", "port", app.settings.Server.Port, "context", app.settings.Server.Context)
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%s", app.settings.Server.Port),
		Handler:        nil,
		ReadTimeout:    app.settings.Server.ReadTimeout,
		WriteTimeout:   app.settings.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		sig := <-quit
		slog.Info("OS signal received, gracefully shutting down HTTP server", "signal", sig.String())

		shutdownTimeout := 15 * time.Second
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		httpServer.SetKeepAlivesEnabled(false)
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("error stopping HTTP server", "err", err)
		}
	}()

	slog.Info("starting HTTP server", "address", httpServer.Addr)

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server failed", "err", err)
		return
	}

	slog.Info("HTTP server stopped")
}

func (app *app) setupHTTPRoutes(ctx context.Context, service service.LogParserService) {
	router := mux.NewRouter()
	prefix := fmt.Sprintf("/%s", app.settings.Server.Context)

	appRouter := router.PathPrefix(prefix).Subrouter()
	logParserHandler := &handler.LogParserHTTPHandler{
		LogParserService: service,
	}

	appRouter.HandleFunc("/game/{gameID}", logParserHandler.GetMatchesStatistics).Methods(http.MethodPost)

	healthRouter := router.PathPrefix(prefix).Subrouter()
	healthRouter.Handle("/health", health.NewHandler()).Methods(http.MethodGet)

	http.Handle("/", router)
}
