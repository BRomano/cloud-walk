package domain

type MatchKills struct {
	Killer     int
	Killed     int
	DeathCause int
}

type PlayerInfo struct {
	Name  string
	Model string
}

func NewMatchData() MatchData {
	return MatchData{Players: map[int]PlayerInfo{}, Kills: []MatchKills{}}
}

type MatchData struct {
	GameName string
	MapName  string
	Hostname string
	Players  map[int]PlayerInfo
	Kills    []MatchKills
}

type MatchStatistics struct {
	TotalKills int            `json:"total_kills"`
	Players    map[int]string `json:"players"`
	Kills      map[string]int `json:"kills"`
}

type MatchDeathStatistics struct {
	KillsByMeans map[string]int `json:"kills_by_means"`
}
