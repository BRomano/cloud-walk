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
	endGame        = "ShutdownGame:"
	userInfoEntry  = "ClientUserinfoChanged:"
	item           = "Item:"
	killEntry      = "Kill:"
)

func NewQuake3ArenaParser() repository.LogParser {
	return &quake3Arena{}
}

type quake3Arena struct {
}

func (quake *quake3Arena) CollectStatisticsFromLog(logger []byte) (map[string]domain.MatchData, error) {
	games := make(map[string]domain.MatchData)

	lines := strings.Split(string(logger), "\n")
	gameLog := make([]string, 0)
	gameIndex := 1
	for _, line := range lines {
		if strings.Contains(line, startGameEntry) {
			gameLog = make([]string, 0)
		}

		gameLog = append(gameLog, line)

		if strings.Contains(line, endGame) {
			games[fmt.Sprintf("game_%d", gameIndex)] = parseGame(gameLog)
			gameIndex++
		}
	}

	return games, nil
}

func parseGame(log []string) domain.MatchData {
	statistics := domain.NewMatchData()
	for _, line := range log {
		if strings.Contains(line, startGameEntry) {
			parseStartGame(line, &statistics)
		} else if strings.Contains(line, userInfoEntry) {
			parseUserInfo(line, &statistics)
		} else if strings.Contains(line, killEntry) {
			parseKill(line, &statistics)
		}
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

func parseStartGame(line string, statistics *domain.MatchData) error {
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
