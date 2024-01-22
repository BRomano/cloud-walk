package repository_test

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestQuake3Arena_CollectStatisticsFromLog(t *testing.T) {
	tests := map[string]struct {
		filename string

		wantResult func(t *testing.T, statistics map[string]domain.MatchData, gotErr error)
	}{
		"should parse an empty file": {
			filename: "./testdata/q3agame_empty.log",
			wantResult: func(t *testing.T, gotStatistics map[string]domain.MatchData, gotErr error) {
				assert.NoError(t, gotErr)
				wantStatistics := map[string]domain.MatchData{
					"game_1": {
						GameName: "baseq3",
						MapName:  "q3dm17",
						Hostname: "Code Miner Server",
						Players:  map[int]domain.PlayerInfo{2: {Name: "Isgalamido", Model: "uriel/zael"}},
						Kills:    []domain.MatchKills{},
					},
				}
				assert.Equal(t, wantStatistics, gotStatistics)
			},
		},
		"should parse correct with kills statistics": {
			filename: "./testdata/q3agame_fullgame.log",
			wantResult: func(t *testing.T, statistics map[string]domain.MatchData, gotErr error) {
				assert.NoError(t, gotErr)
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
