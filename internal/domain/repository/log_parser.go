//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package repository

import "cloud-walk/internal/domain"

const (
	Quake3Arena = iota
)

type LogParser interface {
	CollectStatisticsFromLog(logger []byte) (map[string]domain.MatchData, error)
}
