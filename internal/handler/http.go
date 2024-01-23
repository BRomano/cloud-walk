package handler

import (
	"cloud-walk/internal/domain/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LogParserHTTPHandler struct {
	LogParserService service.LogParserService
}

type loggerKey struct{}

func (handler *LogParserHTTPHandler) GetMatchesStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	startTime := time.Now()

	urlVars := mux.Vars(r)
	gameIDStr, exists := urlVars["gameID"]
	if !exists {
		err := fmt.Errorf("invalid gameID")
		slog.Error("error getting parameter gameID", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger := slog.With("game_id", gameIDStr)
	ctx := context.WithValue(r.Context(), loggerKey{}, logger)

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		err := fmt.Errorf("could not parse gameID")
		logger.Error("error converting gameID", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Error("could not handle file due to", "err", err)
		http.Error(w, "error handling file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		logger.Error("could not read file, due to", "err", err)
		http.Error(w, "error on reading file", http.StatusInternalServerError)
		return
	}

	statistics, err := handler.LogParserService.GetMatchesStatistics(ctx, gameID, fileContent)
	if err != nil {
		logger.Error("could not get game statistics", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(statistics); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	logger.With("elapsed_time", time.Since(startTime))
	if err != nil {
		logger.Error("get matches statistics finished with error", "err", err)
	} else {
		logger.Info("get matches statistics finished")
	}
}
