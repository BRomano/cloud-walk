//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package service

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/domain/repository"
	"context"
	"fmt"
)

type LogParserService interface {
	GetGamesStatistics(game int) (map[string]domain.MatchStatistics, error)
	GetKillsByMeans(game int) (map[string]domain.MatchDeathStatistics, error)
}

type logParserFactory func(game int) (repository.LogParser, error)

type logParserService struct {
	factory logParserFactory
}

func NewLogParser(ctx context.Context) LogParserService {
	return &logParserService{}
}

func (parserService *logParserService) GetGamesStatistics(game int) (map[string]domain.MatchStatistics, error) {
	parser, err := parserService.factory(game)
	if err != nil {
		return nil, fmt.Errorf("could not acquire log parser due to %w", err)
	}

	parser.CollectStatisticsFromLog(nil)
	return nil, nil
}

func (parserService *logParserService) GetKillsByMeans(game int) (map[string]domain.MatchDeathStatistics, error) {
	return nil, nil
}
