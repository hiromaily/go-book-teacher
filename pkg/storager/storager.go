package storages

import (
	"github.com/hiromaily/go-book-teacher/pkg/config"
)

// Storager is Storager interface
type Storager interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}

// NewStorager is to return Storager interface
func NewStorager(conf *config.Config) (Storager, error) {
	if conf.ValidateRedis() {
		// redis mode
		rd, err := NewRedis(conf.Redis.URL)
		if err == nil {
			return rd, nil
		}
	}

	if conf.ValidateText() {
		// text mode
		return NewText(conf.Text.Path), nil
	}

	return NewDummy(), nil
}
