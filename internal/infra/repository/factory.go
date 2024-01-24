package repository

import (
	"cloud-walk/internal/domain/repository"
	"fmt"
)

const (
	UnknowGame = iota
	Quake3Arena
)

func LogParserFactory(gameID int) (repository.LogParser, error) {
	switch gameID {
	case Quake3Arena:
		return NewQuake3ArenaParser(), nil
	default:
		return nil, fmt.Errorf("could not find parser %d", gameID)
	}
}
