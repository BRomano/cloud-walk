package repository

import (
	"cloud-walk/internal/domain/repository"
	"fmt"
)

func LogParserFactory(game int) (repository.LogParser, error) {
	switch game {
	case repository.Quake3Arena:
		return NewQuake3ArenaParser(), nil
	default:
		return nil, fmt.Errorf("could not find parser %d", game)
	}
}
