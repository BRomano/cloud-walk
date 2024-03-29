package service_test

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/domain/repository"
	"cloud-walk/internal/domain/repository/mock"
	"cloud-walk/internal/domain/service"
	"context"
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
		"should parse an ok match for statistics": {
			filename: "./testdata/q3agame_2score.json",
			initializeMock: func(f *mock.MockLogParser) *mock.MockLogParser {
				f.EXPECT().CollectStatisticsFromLog(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, log []byte) (map[string]domain.MatchData, error) {
						matchData := make(map[string]domain.MatchData)
						err := json.Unmarshal(log, &matchData)
						return matchData, err
					}).Times(1)
				return f
			},
			wantResult: func(t *testing.T, gotMatchStatistics map[string]domain.MatchStatistics, gotErr error) {
				assert.NoError(t, gotErr)
				wantMatchStatistics := map[string]domain.MatchStatistics{
					"game_001": {
						TotalKills:  89,
						PlayersName: []string{"Dono da Bola", "Isgalamido", "Zeh", "Assasinu Credi", "Mal", "Oootsimo"},
						Kills:       map[string]int{"Isgalamido": 20, "Oootsimo": 14, "Zeh": 12, "Assasinu Credi": 8, "Dono da Bola": -1, "Mal": -4},
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
		"2 should parse an ok match for statistics": {
			filename: "./testdata/q3agame_score.json",
			initializeMock: func(f *mock.MockLogParser) *mock.MockLogParser {
				f.EXPECT().CollectStatisticsFromLog(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, log []byte) (map[string]domain.MatchData, error) {
						matchData := make(map[string]domain.MatchData)
						err := json.Unmarshal(log, &matchData)
						return matchData, err
					}).Times(1)
				return f
			},
			wantResult: func(t *testing.T, gotMatchStatistics map[string]domain.MatchStatistics, gotErr error) {
				assert.NoError(t, gotErr)
				wantMatchStatistics := map[string]domain.MatchStatistics{
					"game_001": {
						TotalKills:  105,
						PlayersName: []string{"Dono da Bola", "Isgalamido", "Zeh", "Assasinu Credi"},
						Kills:       map[string]int{"Isgalamido": 19, "Dono da Bola": 5, "Zeh": 20, "Assasinu Credi": 11},
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
		"2 matchs on the same log": {
			filename: "./testdata/q3agame_2match_socre.json",
			initializeMock: func(f *mock.MockLogParser) *mock.MockLogParser {
				f.EXPECT().CollectStatisticsFromLog(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, log []byte) (map[string]domain.MatchData, error) {
						matchData := make(map[string]domain.MatchData)
						err := json.Unmarshal(log, &matchData)
						return matchData, err
					}).Times(1)
				return f
			},
			wantResult: func(t *testing.T, gotMatchStatistics map[string]domain.MatchStatistics, gotErr error) {
				assert.NoError(t, gotErr)
				wantMatchStatistics := map[string]domain.MatchStatistics{
					"game_001": {
						TotalKills:  89,
						PlayersName: []string{"Dono da Bola", "Isgalamido", "Zeh", "Assasinu Credi", "Mal", "Oootsimo"},
						Kills:       map[string]int{"Isgalamido": 20, "Oootsimo": 14, "Zeh": 12, "Assasinu Credi": 8, "Dono da Bola": -1, "Mal": -4},
					},
					"game_002": {
						TotalKills:  105,
						PlayersName: []string{"Dono da Bola", "Isgalamido", "Zeh", "Assasinu Credi"},
						Kills:       map[string]int{"Isgalamido": 19, "Dono da Bola": 5, "Zeh": 20, "Assasinu Credi": 11},
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

			statistics, err := logParserService.GetMatchesStatistics(context.Background(), tt.gameID, loggerContent)
			tt.wantResult(t, statistics, err)
		})
	}
}

func TestQuake3Arena_GetKillsByMeans(t *testing.T) {
	tests := map[string]struct {
		filename string
		gameID   int

		LogParserRepository *mock.MockLogParser
		initializeMock      func(f *mock.MockLogParser) *mock.MockLogParser

		wantResult func(t *testing.T, gotMatchStatistics map[string]domain.MatchDeathStatistics, gotErr error)
	}{
		"should parse an ok match for death cause": {
			filename: "./testdata/q3agame_score.json",
			initializeMock: func(f *mock.MockLogParser) *mock.MockLogParser {
				f.EXPECT().CollectStatisticsFromLog(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, log []byte) (map[string]domain.MatchData, error) {
						matchData := make(map[string]domain.MatchData)
						err := json.Unmarshal(log, &matchData)
						return matchData, err
					}).Times(1)
				return f
			},
			wantResult: func(t *testing.T, gotDeathCauses map[string]domain.MatchDeathStatistics, gotErr error) {
				assert.NoError(t, gotErr)
				wantDeathCauses := map[string]domain.MatchDeathStatistics{
					"game_001": {
						KillsByMeans: map[string]int{"MOD_TRIGGER_HURT": 9, "MOD_FALLING": 11, "MOD_ROCKET": 20,
							"MOD_RAILGUN": 8, "MOD_ROCKET_SPLASH": 51, "MOD_MACHINEGUN": 4, "MOD_SHOTGUN": 2},
					},
				}
				assert.Equal(t, wantDeathCauses, gotDeathCauses)
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

			deathCauses, err := logParserService.GetKillsByMeans(context.Background(), tt.gameID, loggerContent)
			tt.wantResult(t, deathCauses, err)
		})
	}
}
