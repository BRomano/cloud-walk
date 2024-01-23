package repository

import (
	"cloud-walk/internal/domain"
	"cloud-walk/internal/domain/repository"
	"fmt"
	"strconv"
	"strings"
)

const (
	startGameEntry = "InitGame:"
	endGameEntry   = "ShutdownGame:"
	userInfoEntry  = "ClientUserinfoChanged:"
	item           = "Item:"
	killEntry      = "Kill:"
	logSeparator   = " "
	breakLog       = "------------------------------------------------------------"
)

func NewQuake3ArenaParser() repository.LogParser {
	lineHandler := make(map[string]lineHandlerFunc)
	lineHandler[startGameEntry] = parseInitGame
	lineHandler[userInfoEntry] = parseUserInfo
	lineHandler[killEntry] = parseKill
	lineHandler[endGameEntry] = parseEndGame

	return &quake3Arena{
		lineHandler: lineHandler,
	}
}

type lineHandlerFunc func(line string, statistics *domain.MatchData) error
type quake3Arena struct {
	lineHandler map[string]lineHandlerFunc
}

type aggregateMatchData struct {
	index   int
	matches map[string]*domain.MatchData
}

func (aggregate *aggregateMatchData) newItem() *domain.MatchData {
	newMatchData := domain.NewMatchData()
	aggregate.matches[fmt.Sprintf("game_%03d", aggregate.index)] = newMatchData
	aggregate.index++
	return newMatchData
}

func (aggregate *aggregateMatchData) getValidMatches() map[string]domain.MatchData {
	result := make(map[string]domain.MatchData)
	for key, value := range aggregate.matches {
		if value.HasStarGame || value.HasEndGame {
			result[key] = *value
		}
	}

	return result
}

func newAggregateMatchData() aggregateMatchData {
	return aggregateMatchData{index: 1, matches: make(map[string]*domain.MatchData)}
}

func (quake *quake3Arena) CollectStatisticsFromLog(logger []byte) (map[string]domain.MatchData, error) {
	matches := newAggregateMatchData()
	lines := strings.Split(string(logger), "\n")
	if len(lines) < 2 { //3 lines at least break and InitGame
		return nil, fmt.Errorf("there is no log to parse")
	}

	var matchData *domain.MatchData = nil
	for _, line := range lines {
		if matchData == nil {
			matchData = matches.newItem()
		}
		if strings.Contains(line, startGameEntry) && matchData.HasStarGame {
			matchData = matches.newItem()
		}

		quake.incMatchData(line, matchData)

		if strings.Contains(line, endGameEntry) && matchData.HasEndGame {
			matchData = matches.newItem()
		}

		if strings.Contains(line, breakLog) && (matchData.HasStarGame || matchData.HasEndGame) {
			matchData = matches.newItem()
		}
	}

	return matches.getValidMatches(), nil
}

func (quake *quake3Arena) incMatchData(logLine string, statistics *domain.MatchData) *domain.MatchData {
	logLine = strings.TrimSpace(logLine)
	tokens := strings.Split(logLine, logSeparator)
	if len(tokens) < 2 {
		return statistics
	}

	if lineHandler, ok := quake.lineHandler[tokens[1]]; ok {
		lineHandler(logLine, statistics)
	}

	return statistics
}

type matchProperties map[string]string

func (properties matchProperties) getAttributeValue(key string) string {
	if value, ok := properties[key]; ok {
		return value
	}
	return ""
}

type getPropertiesFromLogFunc func(line string) string

func getPropertiesInitGame(line string) string {
	return line[strings.Index(line, "\\")+1:]
}

func getPropertiesUserInfo(line string) string {
	propertyStr := line[strings.Index(line, userInfoEntry)+len(userInfoEntry):]
	return propertyStr[strings.Index(propertyStr, "n"):]
}

func newMatchProperties(line string, propertiesExtractor getPropertiesFromLogFunc) matchProperties {
	propertiesStr := propertiesExtractor(line)
	elements := strings.Split(propertiesStr, "\\")

	attributeMap := matchProperties{}
	for i := 0; i < len(elements); i += 2 {
		attributeMap[elements[i]] = elements[i+1]
	}
	return attributeMap
}

func parseInitGame(line string, statistics *domain.MatchData) error {
	statistics.HasStarGame = true
	attributeMap := newMatchProperties(line, getPropertiesInitGame)
	statistics.MapName = attributeMap.getAttributeValue("mapname")
	statistics.GameName = attributeMap.getAttributeValue("gamename")
	statistics.Hostname = attributeMap.getAttributeValue("sv_hostname")

	return nil
}

func getPlayerID(line string) (int, error) {
	propertyStr := line[strings.Index(line, userInfoEntry)+len(userInfoEntry):]
	player := propertyStr[:strings.Index(propertyStr, "n")]
	playerID, err := strconv.Atoi(strings.TrimSpace(player))
	if err != nil {
		return 0, err
	}
	return playerID, nil
}
func parseUserInfo(line string, statistics *domain.MatchData) error {
	attributeMap := newMatchProperties(line, getPropertiesUserInfo)
	playerID, err := getPlayerID(line)
	if err != nil {
		return fmt.Errorf("error on acquire players ID due to %w", err)
	}

	statistics.Players[playerID] = domain.PlayerInfo{
		Model: attributeMap.getAttributeValue("model"),
		Name:  attributeMap.getAttributeValue("n"),
	}

	return nil
}

func stringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
func parseKill(line string, statistics *domain.MatchData) error {
	propertyStr := line[strings.Index(line, killEntry)+len(killEntry):]
	propertyStr = strings.TrimSpace(propertyStr)
	propertyStr = propertyStr[:strings.Index(propertyStr, ":")]
	tokens := strings.Split(propertyStr, " ")
	killer, killed, cause := tokens[0], tokens[1], tokens[2]

	kill := domain.MatchKills{
		Killer:     stringToInt(killer),
		Killed:     stringToInt(killed),
		DeathCause: stringToInt(cause),
	}
	statistics.Kills = append(statistics.Kills, kill)
	return nil
}

func parseEndGame(line string, statistics *domain.MatchData) error {
	statistics.HasEndGame = true
	return nil
}
