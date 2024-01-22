//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package service

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/domain/repository"
	"fmt"
)

type LogParserService interface {
	GetGamesStatistics(game int, logger []byte) (map[string]domain.MatchStatistics, error)
	GetKillsByMeans(game int, logger []byte) (map[string]domain.MatchDeathStatistics, error)
}

type logParserFactory func(game int) (repository.LogParser, error)

func NewLogParser(factory logParserFactory) LogParserService {
	return &logParserService{factory: factory}
}

type logParserService struct {
	factory logParserFactory
}

func (parserService *logParserService) GetGamesStatistics(game int, logger []byte) (map[string]domain.MatchStatistics, error) {
	parser, err := parserService.factory(game)
	if err != nil {
		return nil, fmt.Errorf("could not acquire log parser due to %w", err)
	}

	matches, err := parser.CollectStatisticsFromLog(logger)
	if err != nil {
		return nil, fmt.Errorf("could not acquire matches statistics due to %w", err)
	}

	matchesStatistics := make(map[string]domain.MatchStatistics)
	for matchID, match := range matches {
		matchStatistics := domain.NewMatchStatistics()
		matchStatistics.PlayersName = match.GetPlayersName()
		matchStatistics.TotalKills = len(match.Kills)
		matchStatistics.Kills = calcKillScore(match)
		matchesStatistics[matchID] = matchStatistics
	}

	return matchesStatistics, nil
}

type playersScoreMap map[int]int

func (playersScore playersScoreMap) incScore(playerID int, point int) {
	score, exists := playersScore[playerID]
	if !exists {
		score = 0
	}

	score += point
	playersScore[playerID] = score
}

func calcKillScore(matchData domain.MatchData) map[string]int {
	playerID2Score := make(playersScoreMap)
	for _, m := range matchData.Kills {
		playerID2Score.incScore(m.Killer, 1)
		playerID2Score.incScore(m.Killed, -1)
	}

	playerName2Score := make(map[string]int)
	for playerID, playerScore := range playerID2Score {
		if playerData, exists := matchData.Players[playerID]; exists {
			playerName2Score[playerData.Name] = playerScore
		}
	}

	return playerName2Score
}

func (parserService *logParserService) GetKillsByMeans(game int) (map[string]domain.MatchDeathStatistics, error) {
	return nil, nil
}
