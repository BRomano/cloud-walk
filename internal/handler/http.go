package handler

import (
	"cloud-walk/internal/domain/service"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LogParserHTTPHandler struct {
	LogParserService service.LogParserService
}

func Error(w http.ResponseWriter, msg error, statusCode int) {
	render(w, struct {
		Error string `json:"message"`
	}{
		Error: msg.Error(),
	}, statusCode)
}

func render(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

type loggerKey struct{}

func (handler *LogParserHTTPHandler) GetMatchesStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	startTime := time.Now()

	urlVars := mux.Vars(r)
	gameIDStr, exists := urlVars["gameID"]
	if !exists {
		err := fmt.Errorf("invalid gameID")
		Error(w, err, http.StatusBadRequest)
		slog.Error("error getting parameter gameID", "err", err)
		return
	}

	logger := slog.With("game_id", gameIDStr)
	ctx := context.WithValue(r.Context(), loggerKey{}, logger)

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		err := fmt.Errorf("could not parse gameID")
		Error(w, err, http.StatusBadRequest)
		logger.Error("error converting gameID", "err", err)
		return
	}

	statistics, err := handler.LogParserService.GetMatchesStatistics(ctx, gameID, nil)
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
