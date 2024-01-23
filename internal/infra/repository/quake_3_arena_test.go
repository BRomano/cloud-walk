package repository_test

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func newMatchKill(killer, killed, deathCause int) domain.MatchKills {
	return domain.MatchKills{Killer: killer, Killed: killed, DeathCause: deathCause}
}

func TestQuake3Arena_CollectStatisticsFromLog(t *testing.T) {
	tests := map[string]struct {
		filename string

		wantResult func(t *testing.T, statistics map[string]domain.MatchData, gotErr error)
	}{
		"should parse a minimal log": {
			filename: "./testdata/q3agame_empty.log",
			wantResult: func(t *testing.T, gotStatistics map[string]domain.MatchData, gotErr error) {
				assert.NoError(t, gotErr)
				wantStatistics := map[string]domain.MatchData{
					"game_001": {
						GameName:    "baseq3",
						MapName:     "q3dm17",
						Hostname:    "Code Miner Server",
						Players:     map[int]domain.PlayerInfo{2: {Name: "Isgalamido", Model: "uriel/zael"}},
						Kills:       []domain.MatchKills{},
						HasStarGame: true,
						HasEndGame:  true,
					},
				}
				assert.Equal(t, wantStatistics, gotStatistics)
			},
		},
		"should parse correct with kills statistics": {
			filename: "./testdata/q3agame_ok.log",
			wantResult: func(t *testing.T, gotStatistics map[string]domain.MatchData, gotErr error) {
				assert.NoError(t, gotErr)
				wantStatistics := map[string]domain.MatchData{
					"game_001": {
						GameName: "baseq3",
						MapName:  "q3dm17",
						Hostname: "Code Miner Server",
						Players:  map[int]domain.PlayerInfo{2: {Name: "Isgalamido", Model: "uriel/zael"}, 3: {Name: "Oootsimo", Model: "razor/id"}, 4: {Name: "Dono da Bola", Model: "sarge"}},
						Kills: []domain.MatchKills{newMatchKill(3, 4, 7), newMatchKill(2, 5, 6), newMatchKill(1022, 7, 22),
							newMatchKill(4, 3, 7), newMatchKill(6, 2, 6), newMatchKill(6, 7, 6), newMatchKill(3, 4, 6)},
						HasStarGame: true,
						HasEndGame:  true,
					},
				}
				assert.Equal(t, wantStatistics, gotStatistics)
			},
		},
		"3 match statistics, only one is correct": {
			filename: "./testdata/q3agame_2invalid.log",
			wantResult: func(t *testing.T, gotStatistics map[string]domain.MatchData, gotErr error) {
				assert.NoError(t, gotErr)
				wantStatistics := map[string]domain.MatchData{
					"game_001": {
						GameName:    "game1",
						MapName:     "q3dm17",
						Hostname:    "Code Miner Server",
						Players:     map[int]domain.PlayerInfo{},
						Kills:       []domain.MatchKills{},
						HasStarGame: true,
					},
					"game_002": {
						GameName:    "game2",
						MapName:     "q3dm17",
						Hostname:    "Code Miner Server",
						Players:     map[int]domain.PlayerInfo{2: {Name: "Mocinha", Model: "sarge"}},
						Kills:       []domain.MatchKills{},
						HasStarGame: true, HasEndGame: true,
					},
					"game_003": {
						Players:    map[int]domain.PlayerInfo{},
						Kills:      []domain.MatchKills{{Killer: 6, Killed: 7, DeathCause: 6}, {Killer: 3, Killed: 4, DeathCause: 6}},
						HasEndGame: true,
					},
				}
				assert.Equal(t, wantStatistics, gotStatistics)
			},
		},
		"corrupt file": {
			filename: "./testdata/q3agame_corrupt.log",
			wantResult: func(t *testing.T, statistics map[string]domain.MatchData, gotErr error) {
				assert.Error(t, gotErr)
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			//t.Parallel()
			quakeParser := repository.NewQuake3ArenaParser()
			loggerContent, err := os.ReadFile(tt.filename)
			if err != nil {
				panic(err)
			}

			statistics, err := quakeParser.CollectStatisticsFromLog(loggerContent)
			tt.wantResult(t, statistics, err)
		})
	}
}
