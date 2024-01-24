//go:generate mockgen -source $GOFILE -destination mock/$GOFILE -package=mock
package repository

import (
	"cloud-walk/internal/domain"
	"context"
)

type LogParser interface {
	CollectStatisticsFromLog(ctx context.Context, logger []byte) (map[string]domain.MatchData, error)
}
