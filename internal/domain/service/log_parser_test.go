package service_test

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/domain/repository"
	"cloud-walk/internal/domain/repository/mock"
	"cloud-walk/internal/domain/service"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestQuake3Arena_GetGamesStatistics(t *testing.T) {
	tests := map[string]struct {
		filename string
		gameID   int

		LogParserRepository *mock.MockLogParser
		initializeMock      func(f *mock.MockLogParser) *mock.MockLogParser

		wantResult func(t *testing.T, gotMatchStatistics map[string]domain.MatchStatistics, gotErr error)
	}{
		"should parse an empty file": {
			filename: "./testdata/q3agame_ok_game.json",
			initializeMock: func(f *mock.MockLogParser) *mock.MockLogParser {
				f.EXPECT().CollectStatisticsFromLog(gomock.Any()).
					DoAndReturn(func(log []byte) (map[string]domain.MatchData, error) {
						matchData := make(map[string]domain.MatchData)
						err := json.Unmarshal(log, &matchData)
						return matchData, err
					}).AnyTimes()
				return f
			},
			wantResult: func(t *testing.T, gotMatchStatistics map[string]domain.MatchStatistics, gotErr error) {
				assert.NoError(t, gotErr)
				wantMatchStatistics := map[string]domain.MatchStatistics{
					"game_1": {
						TotalKills:  105,
						PlayersName: []string{"Dono da Bola", "Isgalamido", "Zeh", "Assasinu Credi"},
						Kills:       map[string]int{"Isgalamido": 4, "Dono da Bola": -11, "Zeh": -5, "Assasinu Credi": -8},
					},
				}
				for key := range wantMatchStatistics {
					if _, exists := gotMatchStatistics[key]; assert.Truef(t, exists, "%#v does not exist on result", key) {
					}
					assert.Equalf(t, wantMatchStatistics[key].TotalKills, gotMatchStatistics[key].TotalKills, "Total kills of %#v does not match", key)
					assert.ElementsMatchf(t, wantMatchStatistics[key].PlayersName, gotMatchStatistics[key].PlayersName, "Players name of %#v does not match", key)
					assert.Equalf(t, wantMatchStatistics[key].Kills, gotMatchStatistics[key].Kills, "Players name of %#v does not match", key)
				}
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			//t.Parallel()
			ctrl := gomock.NewController(t)
			tt.LogParserRepository = tt.initializeMock(mock.NewMockLogParser(ctrl))
			factory := func(gameID int) (repository.LogParser, error) {
				return tt.LogParserRepository, nil
			}

			logParserService := service.NewLogParser(factory)
			loggerContent, err := os.ReadFile(tt.filename)
			if err != nil {
				panic(err)
			}

			statistics, err := logParserService.GetGamesStatistics(tt.gameID, loggerContent)
			tt.wantResult(t, statistics, err)
		})
	}
}
