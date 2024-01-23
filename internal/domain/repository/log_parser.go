//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package repository

import (
	"cloud-walk/internal/domain"
	"context"
)

const (
	UnknowGame = iota
	Quake3Arena
)

type LogParser interface {
	CollectStatisticsFromLog(ctx context.Context, logger []byte) (map[string]domain.MatchData, error)
}
