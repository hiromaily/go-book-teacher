package storage

import (
	"github.com/pkg/errors"
)

// Storager interface
type Storager interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}

// NewStorager is to return Storager interface
func NewStorager(mode Mode, redisURL, textPath string) (Storager, error) {
	switch mode {
	case RedisMode:
		return NewRedis(redisURL)
	case TextMode:
		return NewText(textPath), nil
	case DummyMode:
		return NewDummy(), nil
	}
	return nil, errors.New("storage mode is not found")
}
