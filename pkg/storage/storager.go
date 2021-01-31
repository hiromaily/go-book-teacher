package storage

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Storager interface
type Storager interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}

// NewStorager is to return Storager interface
func NewStorager(mode Mode, logger *zap.Logger, redisURL, textPath string) (Storager, error) {
	switch mode {
	case RedisMode:
		return NewRedis(logger, redisURL)
	case TextMode:
		return NewText(logger, textPath), nil
	case DummyMode:
		return NewDummy(logger), nil
	}
	return nil, errors.New("storage mode is not found")
}
