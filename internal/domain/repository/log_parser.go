//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package repository

import "cloud-walk/internal/domain"

const (
	UnknowGame = iota
	Quake3Arena
)

type LogParser interface {
	CollectStatisticsFromLog(logger []byte) (map[string]domain.MatchData, error)
}
