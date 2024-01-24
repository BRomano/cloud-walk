package domain

type MatchKills struct {
	Killer     int `json:"killer"`
	Killed     int `json:"killed"`
	DeathCause int `json:"death_cause"`
}

type PlayerInfo struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

func NewMatchData() *MatchData {
	return &MatchData{Players: map[int]PlayerInfo{}, Kills: []MatchKills{}}
}

type MatchData struct {
	GameName    string             `json:"game_name"`
	MapName     string             `json:"map_name"`
	Hostname    string             `json:"hostname"`
	Players     map[int]PlayerInfo `json:"players"`
	Kills       []MatchKills       `json:"kills"`
	Err         error
	HasStarGame bool
	HasEndGame  bool
}

func (match *MatchData) GetPlayersName() []string {
	players := make([]string, 0)
	for _, player := range match.Players {
		players = append(players, player.Name)
	}
	return players
}

func NewMatchStatistics() MatchStatistics {
	return MatchStatistics{
		PlayersName: []string{},
		Kills:       map[string]int{},
	}
}

type MatchStatistics struct {
	TotalKills  int            `json:"total_kills"`
	PlayersName []string       `json:"players"`
	Kills       map[string]int `json:"kills"`
}

type MatchDeathStatistics struct {
	KillsByMeans map[string]int `json:"kills_by_means"`
}
