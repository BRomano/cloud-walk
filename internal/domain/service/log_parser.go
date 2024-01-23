//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package service

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/domain/repository"
	"context"
	"fmt"
)

const (
	worldKillerID = 1022
)

type LogParserService interface {
	GetMatchesStatistics(ctx context.Context, gameID int, logger []byte) (map[string]domain.MatchStatistics, error)
	GetKillsByMeans(gameID int, logger []byte) (map[string]domain.MatchDeathStatistics, error)
}

type logParserFactory func(game int) (repository.LogParser, error)

func NewLogParser(factory logParserFactory) LogParserService {
	return &logParserService{factory: factory}
}

type logParserService struct {
	factory logParserFactory
}

func (parserService *logParserService) GetMatchesStatistics(ctx context.Context, game int, logger []byte) (map[string]domain.MatchStatistics, error) {
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
		if m.Killer != m.Killed {
			playerID2Score.incScore(m.Killer, 1)
		}
		if m.Killer == worldKillerID || m.Killed == m.Killer {
			playerID2Score.incScore(m.Killed, -1)
		}
	}

	playerName2Score := make(map[string]int)
	for playerID, playerScore := range playerID2Score {
		if playerData, exists := matchData.Players[playerID]; exists {
			playerName2Score[playerData.Name] = playerScore
		}
	}

	return playerName2Score
}

func (parserService *logParserService) GetKillsByMeans(game int, logger []byte) (map[string]domain.MatchDeathStatistics, error) {
	parser, err := parserService.factory(game)
	if err != nil {
		return nil, fmt.Errorf("could not acquire log parser due to %w", err)
	}

	matches, err := parser.CollectStatisticsFromLog(logger)

	deathCausesStatistics := make(map[string]domain.MatchDeathStatistics)
	for matchID, match := range matches {
		deathCausesStatistics[matchID] = domain.MatchDeathStatistics{
			KillsByMeans: calcDeathCauses(match),
		}
	}

	return deathCausesStatistics, nil
}

type deathCause2Count map[int]int

func (deathCause deathCause2Count) incDeath(deathID int) {
	if kills, exists := deathCause[deathID]; exists {
		deathCause[deathID] = kills + 1
	} else {
		deathCause[deathID] = 1
	}
}

func calcDeathCauses(match domain.MatchData) map[string]int {
	deathCauses := make(deathCause2Count)
	for _, kill := range match.Kills {
		deathCauses.incDeath(kill.DeathCause)
	}

	deathCausesStrMap := make(map[string]int)
	for key, value := range deathCauses {
		deathCausesStrMap[getDeathCauseByID(key)] = value
	}

	return deathCausesStrMap
}
